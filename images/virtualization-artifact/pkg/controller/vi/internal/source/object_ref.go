/*
Copyright 2024 Flant JSC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package source

import (
	"context"
	"errors"
	"fmt"

	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/ptr"
	cdiv1 "kubevirt.io/containerized-data-importer-api/pkg/apis/core/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	cc "github.com/deckhouse/virtualization-controller/pkg/common"
	"github.com/deckhouse/virtualization-controller/pkg/common/datasource"
	"github.com/deckhouse/virtualization-controller/pkg/controller"
	"github.com/deckhouse/virtualization-controller/pkg/controller/common"
	"github.com/deckhouse/virtualization-controller/pkg/controller/importer"
	"github.com/deckhouse/virtualization-controller/pkg/controller/service"
	"github.com/deckhouse/virtualization-controller/pkg/controller/supplements"
	"github.com/deckhouse/virtualization-controller/pkg/dvcr"
	"github.com/deckhouse/virtualization-controller/pkg/logger"
	"github.com/deckhouse/virtualization-controller/pkg/sdk/framework/helper"
	"github.com/deckhouse/virtualization-controller/pkg/util"
	virtv2 "github.com/deckhouse/virtualization/api/core/v1alpha2"
	"github.com/deckhouse/virtualization/api/core/v1alpha2/vicondition"
)

const objectRefDataSource = "objectref"

type ObjectRefDataSource struct {
	statService        Stat
	importerService    Importer
	dvcrSettings       *dvcr.Settings
	client             client.Client
	diskService        *service.DiskService
	storageClassForPVC string

	viObjectRefOnPvc *ObjectRefDataVirtualImageOnPVC
	vdSyncer         *ObjectRefVirtualDisk
}

func NewObjectRefDataSource(
	statService Stat,
	importerService Importer,
	dvcrSettings *dvcr.Settings,
	client client.Client,
	diskService *service.DiskService,
	storageClassForPVC string,
) *ObjectRefDataSource {
	return &ObjectRefDataSource{
		statService:        statService,
		importerService:    importerService,
		dvcrSettings:       dvcrSettings,
		client:             client,
		diskService:        diskService,
		storageClassForPVC: storageClassForPVC,
		viObjectRefOnPvc:   NewObjectRefDataVirtualImageOnPVC(statService, importerService, dvcrSettings, client, diskService, storageClassForPVC),
		vdSyncer:           NewObjectRefVirtualDisk(importerService, client, diskService, dvcrSettings, statService, storageClassForPVC),
	}
}

func (ds ObjectRefDataSource) StoreToPVC(ctx context.Context, vi *virtv2.VirtualImage) (reconcile.Result, error) {
	log, ctx := logger.GetDataSourceContext(ctx, objectRefDataSource)

	condition, _ := service.GetCondition(vicondition.ReadyType, vi.Status.Conditions)
	defer func() { service.SetCondition(condition, &vi.Status.Conditions) }()

	switch vi.Spec.DataSource.ObjectRef.Kind {
	case virtv2.VirtualImageKind:
		viKey := types.NamespacedName{Name: vi.Spec.DataSource.ObjectRef.Name, Namespace: vi.Namespace}
		viRef, err := helper.FetchObject(ctx, viKey, ds.client, &virtv2.VirtualImage{})
		if err != nil {
			return reconcile.Result{}, fmt.Errorf("unable to get VI %s: %w", viKey, err)
		}

		if viRef == nil {
			return reconcile.Result{}, fmt.Errorf("VI object ref %s is nil", viKey)
		}

		if viRef.Spec.Storage == virtv2.StorageKubernetes {
			return ds.viObjectRefOnPvc.StoreToPVC(ctx, vi, viRef, &condition)
		}
	case virtv2.VirtualDiskKind:
		viKey := types.NamespacedName{Name: vi.Spec.DataSource.ObjectRef.Name, Namespace: vi.Namespace}
		vd, err := helper.FetchObject(ctx, viKey, ds.client, &virtv2.VirtualDisk{})
		if err != nil {
			return reconcile.Result{}, fmt.Errorf("unable to get VI %s: %w", viKey, err)
		}

		if vd == nil {
			return reconcile.Result{}, fmt.Errorf("VD object ref %s is nil", viKey)
		}

		return ds.vdSyncer.StoreToPVC(ctx, vi, vd, &condition)
	}

	supgen := supplements.NewGenerator(common.VIShortName, vi.Name, vi.Namespace, vi.UID)
	dv, err := ds.diskService.GetDataVolume(ctx, supgen)
	if err != nil {
		return reconcile.Result{}, err
	}
	pvc, err := ds.diskService.GetPersistentVolumeClaim(ctx, supgen)
	if err != nil {
		return reconcile.Result{}, err
	}

	switch {
	case isDiskProvisioningFinished(condition):
		log.Info("Disk provisioning finished: clean up")

		setPhaseConditionForFinishedImage(pvc, &condition, &vi.Status.Phase, supgen)

		// Protect Ready Disk and underlying PVC.
		err = ds.diskService.Protect(ctx, vi, nil, pvc)
		if err != nil {
			return reconcile.Result{}, err
		}

		err = ds.diskService.Unprotect(ctx, dv)
		if err != nil {
			return reconcile.Result{}, err
		}

		return CleanUpSupplements(ctx, vi, ds)
	case common.AnyTerminating(dv, pvc):
		log.Info("Waiting for supplements to be terminated")
	case dv == nil:
		log.Info("Start import to PVC")
		var dvcrDataSource controller.DVCRDataSource
		dvcrDataSource, err = controller.NewDVCRDataSourcesForVMI(ctx, vi.Spec.DataSource, vi, ds.client)
		if err != nil {
			return reconcile.Result{}, err
		}

		if !dvcrDataSource.IsReady() {
			condition.Status = metav1.ConditionFalse
			condition.Reason = vicondition.ProvisioningFailed
			condition.Message = "Failed to get stats from non-ready datasource: waiting for the DataSource to be ready."
			return reconcile.Result{}, nil
		}

		vi.Status.Progress = "0%"
		vi.Status.SourceUID = util.GetPointer(dvcrDataSource.GetUID())

		var diskSize resource.Quantity
		diskSize, err = ds.getPVCSize(dvcrDataSource)
		if err != nil {
			setPhaseConditionToFailed(&condition, &vi.Status.Phase, err)

			if errors.Is(err, service.ErrInsufficientPVCSize) {
				return reconcile.Result{}, nil
			}

			return reconcile.Result{}, err
		}

		var source *cdiv1.DataVolumeSource
		source, err = ds.getSource(supgen, dvcrDataSource)
		if err != nil {
			return reconcile.Result{}, err
		}

		err = ds.diskService.StartImmediate(ctx, diskSize, ptr.To(ds.storageClassForPVC), source, vi, supgen)
		if updated, err := setPhaseConditionFromStorageError(err, vi, &condition); err != nil || updated {
			return reconcile.Result{}, err
		}

		vi.Status.Phase = virtv2.ImageProvisioning
		condition.Status = metav1.ConditionFalse
		condition.Reason = vicondition.Provisioning
		condition.Message = "PVC Provisioner not found: create the new one."

		return reconcile.Result{Requeue: true}, nil
	case pvc == nil:
		vi.Status.Phase = virtv2.ImageProvisioning
		condition.Status = metav1.ConditionFalse
		condition.Reason = vicondition.Provisioning
		condition.Message = "PVC not found: waiting for creation."
		return reconcile.Result{Requeue: true}, nil
	case ds.diskService.IsImportDone(dv, pvc):
		log.Info("Import has completed", "dvProgress", dv.Status.Progress, "dvPhase", dv.Status.Phase, "pvcPhase", pvc.Status.Phase)

		vi.Status.Phase = virtv2.ImageReady
		condition.Status = metav1.ConditionTrue
		condition.Reason = vicondition.Ready
		condition.Message = ""

		var dvcrDataSource controller.DVCRDataSource
		dvcrDataSource, err = controller.NewDVCRDataSourcesForVMI(ctx, vi.Spec.DataSource, vi, ds.client)
		if err != nil {
			return reconcile.Result{}, err
		}

		vi.Status.Size = dvcrDataSource.GetSize()
		vi.Status.CDROM = dvcrDataSource.IsCDROM()
		vi.Status.Format = dvcrDataSource.GetFormat()
		vi.Status.Progress = "100%"
		vi.Status.Target.PersistentVolumeClaim = dv.Status.ClaimName
	default:
		log.Info("Provisioning to PVC is in progress", "dvProgress", dv.Status.Progress, "dvPhase", dv.Status.Phase, "pvcPhase", pvc.Status.Phase)

		vi.Status.Progress = ds.diskService.GetProgress(dv, vi.Status.Progress, service.NewScaleOption(0, 100))
		vi.Status.Target.PersistentVolumeClaim = dv.Status.ClaimName

		err = ds.diskService.Protect(ctx, vi, dv, pvc)
		if err != nil {
			return reconcile.Result{}, err
		}

		err = setPhaseConditionForPVCProvisioningImage(ctx, dv, vi, pvc, &condition, ds.diskService)
		if err != nil {
			return reconcile.Result{}, err
		}

		return reconcile.Result{}, nil
	}

	return reconcile.Result{Requeue: true}, nil
}

func (ds ObjectRefDataSource) StoreToDVCR(ctx context.Context, vi *virtv2.VirtualImage) (reconcile.Result, error) {
	log, ctx := logger.GetDataSourceContext(ctx, "objectref")

	condition, _ := service.GetCondition(vicondition.ReadyType, vi.Status.Conditions)
	defer func() { service.SetCondition(condition, &vi.Status.Conditions) }()

	switch vi.Spec.DataSource.ObjectRef.Kind {
	case virtv2.VirtualImageKind:
		viKey := types.NamespacedName{Name: vi.Spec.DataSource.ObjectRef.Name, Namespace: vi.Namespace}
		viRef, err := helper.FetchObject(ctx, viKey, ds.client, &virtv2.VirtualImage{})
		if err != nil {
			return reconcile.Result{}, fmt.Errorf("unable to get VI %s: %w", viKey, err)
		}

		if viRef == nil {
			return reconcile.Result{}, fmt.Errorf("VI object ref source %s is nil", vi.Spec.DataSource.ObjectRef.Name)
		}

		if viRef.Spec.Storage == virtv2.StorageKubernetes {
			return ds.viObjectRefOnPvc.StoreToDVCR(ctx, vi, viRef, &condition)
		}
	case virtv2.VirtualDiskKind:
		viKey := types.NamespacedName{Name: vi.Spec.DataSource.ObjectRef.Name, Namespace: vi.Namespace}
		vd, err := helper.FetchObject(ctx, viKey, ds.client, &virtv2.VirtualDisk{})
		if err != nil {
			return reconcile.Result{}, fmt.Errorf("unable to get VD %s: %w", viKey, err)
		}

		if vd == nil {
			return reconcile.Result{}, fmt.Errorf("VD object ref %s is nil", viKey)
		}

		return ds.vdSyncer.StoreToDVCR(ctx, vi, vd, &condition)
	}

	supgen := supplements.NewGenerator(common.VIShortName, vi.Name, vi.Namespace, vi.UID)
	pod, err := ds.importerService.GetPod(ctx, supgen)
	if err != nil {
		return reconcile.Result{}, err
	}

	switch {
	case isDiskProvisioningFinished(condition):
		log.Info("Virtual image provisioning finished: clean up")

		condition.Status = metav1.ConditionTrue
		condition.Reason = vicondition.Ready
		condition.Message = ""

		vi.Status.Phase = virtv2.ImageReady

		err = ds.importerService.Unprotect(ctx, pod)
		if err != nil {
			return reconcile.Result{}, err
		}

		return CleanUpSupplements(ctx, vi, ds)
	case common.IsTerminating(pod):
		vi.Status.Phase = virtv2.ImagePending

		log.Info("Cleaning up...")
	case pod == nil:
		vi.Status.Progress = "0%"

		var dvcrDataSource controller.DVCRDataSource
		dvcrDataSource, err = controller.NewDVCRDataSourcesForVMI(ctx, vi.Spec.DataSource, vi, ds.client)
		if err != nil {
			return reconcile.Result{}, err
		}

		vi.Status.SourceUID = util.GetPointer(dvcrDataSource.GetUID())

		var envSettings *importer.Settings
		envSettings, err = ds.getEnvSettings(vi, supgen, dvcrDataSource)
		if err != nil {
			return reconcile.Result{}, err
		}

		err = ds.importerService.Start(ctx, envSettings, vi, supgen, datasource.NewCABundleForVMI(vi.Spec.DataSource))
		switch {
		case err == nil:
			// OK.
		case common.ErrQuotaExceeded(err):
			return setQuotaExceededPhaseCondition(&condition, &vi.Status.Phase, err, vi.CreationTimestamp), nil
		default:
			setPhaseConditionToFailed(&condition, &vi.Status.Phase, fmt.Errorf("unexpected error: %w", err))
			return reconcile.Result{}, err
		}

		vi.Status.Phase = virtv2.ImageProvisioning
		condition.Status = metav1.ConditionFalse
		condition.Reason = vicondition.Provisioning
		condition.Message = "DVCR Provisioner not found: create the new one."

		log.Info("Ready", "progress", vi.Status.Progress, "pod.phase", "nil")

		return reconcile.Result{Requeue: true}, nil
	case common.IsPodComplete(pod):
		err = ds.statService.CheckPod(pod)
		if err != nil {
			vi.Status.Phase = virtv2.ImageFailed

			switch {
			case errors.Is(err, service.ErrProvisioningFailed):
				condition.Status = metav1.ConditionFalse
				condition.Reason = vicondition.ProvisioningFailed
				condition.Message = service.CapitalizeFirstLetter(err.Error() + ".")
				return reconcile.Result{}, nil
			default:
				return reconcile.Result{}, err
			}
		}
		var dvcrDataSource controller.DVCRDataSource
		dvcrDataSource, err = controller.NewDVCRDataSourcesForVMI(ctx, vi.Spec.DataSource, vi, ds.client)
		if err != nil {
			return reconcile.Result{}, err
		}

		if !dvcrDataSource.IsReady() {
			condition.Status = metav1.ConditionFalse
			condition.Reason = vicondition.ProvisioningFailed
			condition.Message = "Failed to get stats from non-ready datasource: waiting for the DataSource to be ready."
			return reconcile.Result{}, nil
		}

		condition.Status = metav1.ConditionTrue
		condition.Reason = vicondition.Ready
		condition.Message = ""

		vi.Status.Phase = virtv2.ImageReady
		vi.Status.Size = dvcrDataSource.GetSize()
		vi.Status.CDROM = dvcrDataSource.IsCDROM()
		vi.Status.Format = dvcrDataSource.GetFormat()
		vi.Status.Progress = "100%"
		vi.Status.Target.RegistryURL = ds.statService.GetDVCRImageName(pod)

		log.Info("Ready", "progress", vi.Status.Progress, "pod.phase", pod.Status.Phase)
	default:
		err = ds.statService.CheckPod(pod)
		if err != nil {
			return reconcile.Result{}, setPhaseConditionFromPodError(&condition, vi, err)
		}

		condition.Status = metav1.ConditionFalse
		condition.Reason = vicondition.Provisioning
		condition.Message = "Import is in the process of provisioning to DVCR."

		vi.Status.Phase = virtv2.ImageProvisioning
		vi.Status.Target.RegistryURL = ds.statService.GetDVCRImageName(pod)

		log.Info("Ready", "progress", vi.Status.Progress, "pod.phase", pod.Status.Phase)
	}

	return reconcile.Result{Requeue: true}, nil
}

func (ds ObjectRefDataSource) CleanUp(ctx context.Context, vi *virtv2.VirtualImage) (bool, error) {
	supgen := supplements.NewGenerator(common.VIShortName, vi.Name, vi.Namespace, vi.UID)

	importerRequeue, err := ds.importerService.CleanUp(ctx, supgen)
	if err != nil {
		return false, err
	}

	diskRequeue, err := ds.diskService.CleanUp(ctx, supgen)
	if err != nil {
		return false, err
	}

	return importerRequeue || diskRequeue, nil
}

func (ds ObjectRefDataSource) Validate(ctx context.Context, vi *virtv2.VirtualImage) error {
	if vi.Spec.DataSource.ObjectRef == nil {
		return fmt.Errorf("nil object ref: %s", vi.Spec.DataSource.Type)
	}

	switch vi.Spec.DataSource.ObjectRef.Kind {
	case virtv2.VirtualImageObjectRefKindVirtualImage:
		viKey := types.NamespacedName{Name: vi.Spec.DataSource.ObjectRef.Name, Namespace: vi.Namespace}
		viRef, err := helper.FetchObject(ctx, viKey, ds.client, &virtv2.VirtualImage{})
		if err != nil {
			return fmt.Errorf("unable to get VI %s: %w", viKey, err)
		}

		if viRef == nil {
			return fmt.Errorf("VI object ref source %s is nil", vi.Spec.DataSource.ObjectRef.Name)
		}

		if viRef.Spec.Storage == virtv2.StorageKubernetes {
			if viRef.Status.Phase != virtv2.ImageReady {
				return NewImageNotReadyError(vi.Spec.DataSource.ObjectRef.Name)
			}
			return nil
		}

		dvcrDataSource, err := controller.NewDVCRDataSourcesForVMI(ctx, vi.Spec.DataSource, vi, ds.client)
		if err != nil {
			return err
		}

		if dvcrDataSource.IsReady() {
			return nil
		}

		return NewImageNotReadyError(vi.Spec.DataSource.ObjectRef.Name)
	case virtv2.VirtualImageObjectRefKindClusterVirtualImage:
		dvcrDataSource, err := controller.NewDVCRDataSourcesForVMI(ctx, vi.Spec.DataSource, vi, ds.client)
		if err != nil {
			return err
		}

		if dvcrDataSource.IsReady() {
			return nil
		}

		return NewClusterImageNotReadyError(vi.Spec.DataSource.ObjectRef.Name)
	case virtv2.VirtualImageObjectRefKindVirtualDisk:
		return ds.vdSyncer.Validate(ctx, vi)
	default:
		return fmt.Errorf("unexpected object ref kind: %s", vi.Spec.DataSource.ObjectRef.Kind)
	}
}

func (ds ObjectRefDataSource) getEnvSettings(vi *virtv2.VirtualImage, sup *supplements.Generator, dvcrDataSource controller.DVCRDataSource) (*importer.Settings, error) {
	if !dvcrDataSource.IsReady() {
		return nil, errors.New("dvcr data source is not ready")
	}

	var settings importer.Settings
	importer.ApplyDVCRSourceSettings(&settings, dvcrDataSource.GetTarget())
	importer.ApplyDVCRDestinationSettings(
		&settings,
		ds.dvcrSettings,
		sup,
		ds.dvcrSettings.RegistryImageForVI(vi),
	)

	return &settings, nil
}

func (ds ObjectRefDataSource) CleanUpSupplements(ctx context.Context, vi *virtv2.VirtualImage) (reconcile.Result, error) {
	supgen := supplements.NewGenerator(common.VIShortName, vi.Name, vi.Namespace, vi.UID)

	importerRequeue, err := ds.importerService.CleanUpSupplements(ctx, supgen)
	if err != nil {
		return reconcile.Result{}, err
	}

	diskRequeue, err := ds.diskService.CleanUpSupplements(ctx, supgen)
	if err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{Requeue: importerRequeue || diskRequeue}, nil
}

func (ds ObjectRefDataSource) getPVCSize(dvcrDataSource controller.DVCRDataSource) (resource.Quantity, error) {
	if !dvcrDataSource.IsReady() {
		return resource.Quantity{}, errors.New("dvcr data source is not ready")
	}

	unpackedSize, err := resource.ParseQuantity(dvcrDataSource.GetSize().UnpackedBytes)
	if err != nil {
		return resource.Quantity{}, fmt.Errorf("failed to parse unpacked bytes %s: %w", dvcrDataSource.GetSize().UnpackedBytes, err)
	}

	if unpackedSize.IsZero() {
		return resource.Quantity{}, errors.New("got zero unpacked size from data source")
	}

	return service.GetValidatedPVCSize(&unpackedSize, unpackedSize)
}

func (ds ObjectRefDataSource) getSource(sup *supplements.Generator, dvcrDataSource controller.DVCRDataSource) (*cdiv1.DataVolumeSource, error) {
	if !dvcrDataSource.IsReady() {
		return nil, errors.New("dvcr data source is not ready")
	}

	url := cc.DockerRegistrySchemePrefix + dvcrDataSource.GetTarget()
	secretName := sup.DVCRAuthSecretForDV().Name
	certConfigMapName := sup.DVCRCABundleConfigMapForDV().Name

	return &cdiv1.DataVolumeSource{
		Registry: &cdiv1.DataVolumeSourceRegistry{
			URL:           &url,
			SecretRef:     &secretName,
			CertConfigMap: &certConfigMapName,
		},
	}, nil
}
