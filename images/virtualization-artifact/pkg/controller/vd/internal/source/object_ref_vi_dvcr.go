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
	cdiv1 "kubevirt.io/containerized-data-importer-api/pkg/apis/core/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	common2 "github.com/deckhouse/virtualization-controller/pkg/common"
	"github.com/deckhouse/virtualization-controller/pkg/controller/common"
	"github.com/deckhouse/virtualization-controller/pkg/controller/service"
	"github.com/deckhouse/virtualization-controller/pkg/controller/supplements"
	"github.com/deckhouse/virtualization-controller/pkg/imageformat"
	"github.com/deckhouse/virtualization-controller/pkg/logger"
	"github.com/deckhouse/virtualization-controller/pkg/util"
	virtv2 "github.com/deckhouse/virtualization/api/core/v1alpha2"
	"github.com/deckhouse/virtualization/api/core/v1alpha2/vdcondition"
)

type ObjectRefVirtualImageDVCR struct {
	statService *service.StatService
	diskService *service.DiskService
	client      client.Client
}

func NewObjectRefVirtualImageDVCR(
	statService *service.StatService,
	diskService *service.DiskService,
	client client.Client,
) *ObjectRefVirtualImageDVCR {
	return &ObjectRefVirtualImageDVCR{
		statService: statService,
		diskService: diskService,
		client:      client,
	}
}

func (ds ObjectRefVirtualImageDVCR) Sync(ctx context.Context, vd *virtv2.VirtualDisk) (reconcile.Result, error) {
	if vd.Spec.DataSource == nil || vd.Spec.DataSource.ObjectRef == nil {
		return reconcile.Result{}, errors.New("object ref missed for data source")
	}

	log, ctx := logger.GetDataSourceContext(ctx, objectRefDataSource)

	condition, _ := service.GetCondition(vdcondition.ReadyType, vd.Status.Conditions)
	defer func() { service.SetCondition(condition, &vd.Status.Conditions) }()

	supgen := supplements.NewGenerator(common.VDShortName, vd.Name, vd.Namespace, vd.UID)
	dv, err := ds.diskService.GetDataVolume(ctx, supgen)
	if err != nil {
		return reconcile.Result{}, err
	}
	pvc, err := ds.diskService.GetPersistentVolumeClaim(ctx, supgen)
	if err != nil {
		return reconcile.Result{}, err
	}
	vi, err := ds.diskService.GetVirtualImage(ctx, vd.Spec.DataSource.ObjectRef.Name, vd.Namespace)
	if err != nil {
		return reconcile.Result{}, err
	}
	if vi == nil {
		return reconcile.Result{}, errors.New("the source virtual image not found")
	}

	switch {
	case isDiskProvisioningFinished(condition):
		log.Debug("Disk provisioning finished: clean up")

		setPhaseConditionForFinishedDisk(pvc, &condition, &vd.Status.Phase, supgen)

		// Protect Ready Disk and underlying PVC.
		err = ds.diskService.Protect(ctx, vd, nil, pvc)
		if err != nil {
			return reconcile.Result{}, err
		}

		err = ds.diskService.Unprotect(ctx, dv)
		if err != nil {
			return reconcile.Result{}, err
		}

		return CleanUpSupplements(ctx, vd, ds)
	case common.AnyTerminating(dv, pvc):
		log.Info("Waiting for supplements to be terminated")
	case dv == nil:
		log.Info("Start import to PVC")

		vd.Status.Progress = "0%"
		vd.Status.SourceUID = util.GetPointer(vi.GetUID())

		if imageformat.IsISO(vi.Status.Format) {
			setPhaseConditionToFailed(&condition, &vd.Status.Phase, ErrISOSourceNotSupported)
			return reconcile.Result{}, nil
		}

		var diskSize resource.Quantity
		diskSize, err = ds.getPVCSize(vd, vi.Status.Size)
		if err != nil {
			setPhaseConditionToFailed(&condition, &vd.Status.Phase, err)

			if errors.Is(err, service.ErrInsufficientPVCSize) {
				return reconcile.Result{}, nil
			}

			return reconcile.Result{}, err
		}

		source := ds.getSource(supgen, vi)

		err = ds.diskService.Start(ctx, diskSize, vd.Spec.PersistentVolumeClaim.StorageClass, source, vd, supgen)
		if updated, err := setPhaseConditionFromStorageError(err, vd, &condition); err != nil || updated {
			return reconcile.Result{}, err
		}

		vd.Status.Phase = virtv2.DiskProvisioning
		condition.Status = metav1.ConditionFalse
		condition.Reason = vdcondition.Provisioning
		condition.Message = "PVC Provisioner not found: create the new one."

		return reconcile.Result{Requeue: true}, nil
	case pvc == nil:
		vd.Status.Phase = virtv2.DiskProvisioning
		condition.Status = metav1.ConditionFalse
		condition.Reason = vdcondition.Provisioning
		condition.Message = "PVC not found: waiting for creation."
		return reconcile.Result{Requeue: true}, nil
	case ds.diskService.IsImportDone(dv, pvc):
		log.Info("Import has completed", "dvProgress", dv.Status.Progress, "dvPhase", dv.Status.Phase, "pvcPhase", pvc.Status.Phase)

		vd.Status.Phase = virtv2.DiskReady
		condition.Status = metav1.ConditionTrue
		condition.Reason = vdcondition.Ready
		condition.Message = ""

		vd.Status.Progress = "100%"
		vd.Status.Capacity = ds.diskService.GetCapacity(pvc)
		vd.Status.Target.PersistentVolumeClaim = dv.Status.ClaimName
	default:
		log.Info("Provisioning to PVC is in progress", "dvProgress", dv.Status.Progress, "dvPhase", dv.Status.Phase, "pvcPhase", pvc.Status.Phase)

		vd.Status.Progress = ds.diskService.GetProgress(dv, vd.Status.Progress, service.NewScaleOption(0, 100))
		vd.Status.Capacity = ds.diskService.GetCapacity(pvc)
		vd.Status.Target.PersistentVolumeClaim = dv.Status.ClaimName

		err = ds.diskService.Protect(ctx, vd, dv, pvc)
		if err != nil {
			return reconcile.Result{}, err
		}
		sc, err := ds.diskService.GetStorageClass(ctx, pvc.Spec.StorageClassName)
		if updated, err := setPhaseConditionFromStorageError(err, vd, &condition); err != nil || updated {
			return reconcile.Result{}, err
		}
		if err = setPhaseConditionForPVCProvisioningDisk(ctx, dv, vd, pvc, sc, &condition, ds.diskService); err != nil {
			return reconcile.Result{}, err
		}

		return reconcile.Result{}, nil
	}

	return reconcile.Result{Requeue: true}, nil
}

func (ds ObjectRefVirtualImageDVCR) Validate(ctx context.Context, vd *virtv2.VirtualDisk) error {
	if vd.Spec.DataSource == nil || vd.Spec.DataSource.ObjectRef == nil {
		return errors.New("object ref missed for data source")
	}

	vi, err := ds.diskService.GetVirtualImage(ctx, vd.Spec.DataSource.ObjectRef.Name, vd.Namespace)
	if err != nil {
		return err
	}

	if vi == nil || vi.Status.Phase != virtv2.ImageReady || vi.Status.Target.RegistryURL == "" {
		return NewImageNotReadyError(vd.Spec.DataSource.ObjectRef.Name)
	}

	return nil
}

func (ds ObjectRefVirtualImageDVCR) CleanUpSupplements(ctx context.Context, vd *virtv2.VirtualDisk) (reconcile.Result, error) {
	supgen := supplements.NewGenerator(common.VDShortName, vd.Name, vd.Namespace, vd.UID)

	requeue, err := ds.diskService.CleanUpSupplements(ctx, supgen)
	if err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{Requeue: requeue}, nil
}

func (ds ObjectRefVirtualImageDVCR) getSource(sup *supplements.Generator, vi *virtv2.VirtualImage) *cdiv1.DataVolumeSource {
	url := common2.DockerRegistrySchemePrefix + vi.Status.Target.RegistryURL
	secretName := sup.DVCRAuthSecretForDV().Name
	certConfigMapName := sup.DVCRCABundleConfigMapForDV().Name

	return &cdiv1.DataVolumeSource{
		Registry: &cdiv1.DataVolumeSourceRegistry{
			URL:           &url,
			SecretRef:     &secretName,
			CertConfigMap: &certConfigMapName,
		},
	}
}

func (ds ObjectRefVirtualImageDVCR) getPVCSize(vd *virtv2.VirtualDisk, imageSize virtv2.ImageStatusSize) (resource.Quantity, error) {
	unpackedSize, err := resource.ParseQuantity(imageSize.UnpackedBytes)
	if err != nil {
		return resource.Quantity{}, fmt.Errorf("failed to parse unpacked bytes %s: %w", imageSize.UnpackedBytes, err)
	}

	return service.GetValidatedPVCSize(vd.Spec.PersistentVolumeClaim.Size, unpackedSize)
}