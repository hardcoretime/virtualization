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

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	cdiv1 "kubevirt.io/containerized-data-importer-api/pkg/apis/core/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	cc "github.com/deckhouse/virtualization-controller/pkg/common"
	"github.com/deckhouse/virtualization-controller/pkg/common/datasource"
	"github.com/deckhouse/virtualization-controller/pkg/controller/common"
	"github.com/deckhouse/virtualization-controller/pkg/controller/service"
	"github.com/deckhouse/virtualization-controller/pkg/controller/supplements"
	"github.com/deckhouse/virtualization-controller/pkg/controller/uploader"
	"github.com/deckhouse/virtualization-controller/pkg/dvcr"
	"github.com/deckhouse/virtualization-controller/pkg/logger"
	virtv2 "github.com/deckhouse/virtualization/api/core/v1alpha2"
	"github.com/deckhouse/virtualization/api/core/v1alpha2/vicondition"
)

const uploadDataSource = "upload"

type UploadDataSource struct {
	statService        Stat
	uploaderService    Uploader
	dvcrSettings       *dvcr.Settings
	diskService        *service.DiskService
	storageClassForPVC string
}

func NewUploadDataSource(
	statService Stat,
	uploaderService Uploader,
	dvcrSettings *dvcr.Settings,
	diskService *service.DiskService,
	storageClassForPVC string,
) *UploadDataSource {
	return &UploadDataSource{
		statService:        statService,
		uploaderService:    uploaderService,
		dvcrSettings:       dvcrSettings,
		diskService:        diskService,
		storageClassForPVC: storageClassForPVC,
	}
}

func (ds UploadDataSource) StoreToPVC(ctx context.Context, vi *virtv2.VirtualImage) (reconcile.Result, error) {
	log, ctx := logger.GetDataSourceContext(ctx, uploadDataSource)

	condition, _ := service.GetCondition(vicondition.ReadyType, vi.Status.Conditions)
	defer func() { service.SetCondition(condition, &vi.Status.Conditions) }()

	supgen := supplements.NewGenerator(common.VIShortName, vi.Name, vi.Namespace, vi.UID)
	pod, err := ds.uploaderService.GetPod(ctx, supgen)
	if err != nil {
		return reconcile.Result{}, err
	}
	svc, err := ds.uploaderService.GetService(ctx, supgen)
	if err != nil {
		return reconcile.Result{}, err
	}
	ing, err := ds.uploaderService.GetIngress(ctx, supgen)
	if err != nil {
		return reconcile.Result{}, err
	}
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

		// Unprotect upload time supplements to delete them later.
		err = ds.uploaderService.Unprotect(ctx, pod, svc, ing)
		if err != nil {
			return reconcile.Result{}, err
		}

		err = ds.diskService.Unprotect(ctx, dv)
		if err != nil {
			return reconcile.Result{}, err
		}

		return CleanUpSupplements(ctx, vi, ds)
	case common.AnyTerminating(pod, svc, ing, dv, pvc):
		log.Info("Waiting for supplements to be terminated")
	case pod == nil || svc == nil || ing == nil:
		log.Info("Start import to DVCR")

		vi.Status.Progress = "0%"

		envSettings := ds.getEnvSettings(vi, supgen)
		err = ds.uploaderService.Start(ctx, envSettings, vi, supgen, datasource.NewCABundleForVMI(vi.Spec.DataSource))
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

		return reconcile.Result{Requeue: true}, nil
	case !common.IsPodComplete(pod):
		log.Info("Provisioning to DVCR is in progress", "podPhase", pod.Status.Phase)

		err = ds.statService.CheckPod(pod)
		if err != nil {
			return reconcile.Result{}, setPhaseConditionFromPodError(&condition, vi, err)
		}

		if !ds.statService.IsUploadStarted(vi.GetUID(), pod) {
			if ds.statService.IsUploaderReady(pod, svc, ing) {
				log.Info("Waiting for the user upload", "pod.phase", pod.Status.Phase)

				vi.Status.Phase = virtv2.ImageWaitForUserUpload
				condition.Status = metav1.ConditionFalse
				condition.Reason = vicondition.WaitForUserUpload
				condition.Message = "Waiting for the user upload."

				vi.Status.ImageUploadURLs = &virtv2.ImageUploadURLs{
					External:  ds.uploaderService.GetExternalURL(ctx, ing),
					InCluster: ds.uploaderService.GetInClusterURL(ctx, svc),
				}
			} else {
				log.Info("Waiting for the uploader to be ready to process the user's upload", "pod.phase", pod.Status.Phase)

				vi.Status.Phase = virtv2.ImagePending
				condition.Status = metav1.ConditionFalse
				condition.Reason = vicondition.ProvisioningNotStarted
				condition.Message = fmt.Sprintf("Waiting for the uploader %q to be ready to process the user's upload.", pod.Name)
			}

			return reconcile.Result{Requeue: true}, nil
		}

		vi.Status.Phase = virtv2.ImageProvisioning
		condition.Status = metav1.ConditionFalse
		condition.Reason = vicondition.Provisioning
		condition.Message = "Import is in the process of provisioning to DVCR."

		vi.Status.Progress = ds.statService.GetProgress(vi.GetUID(), pod, vi.Status.Progress, service.NewScaleOption(0, 50))
		vi.Status.DownloadSpeed = ds.statService.GetDownloadSpeed(vi.GetUID(), pod)

		err = ds.uploaderService.Protect(ctx, pod, svc, ing)
		if err != nil {
			return reconcile.Result{}, err
		}
	case dv == nil:
		log.Info("Start import to PVC")

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

		vi.Status.Progress = "50.0%"
		vi.Status.DownloadSpeed = ds.statService.GetDownloadSpeed(vi.GetUID(), pod)

		var diskSize resource.Quantity
		diskSize, err = ds.getPVCSize(pod)
		if err != nil {
			setPhaseConditionToFailed(&condition, &vi.Status.Phase, err)

			if errors.Is(err, service.ErrInsufficientPVCSize) {
				return reconcile.Result{}, nil
			}

			return reconcile.Result{}, err
		}

		source := ds.getSource(supgen, ds.statService.GetDVCRImageName(pod))

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

		vi.Status.Progress = "100%"
		vi.Status.Size = ds.statService.GetSize(pod)
		vi.Status.DownloadSpeed = ds.statService.GetDownloadSpeed(vi.GetUID(), pod)
		vi.Status.Target.PersistentVolumeClaim = dv.Status.ClaimName

		log.Info("Ready", "vi", vi.Name, "progress", vi.Status.Progress, "dv.phase", dv.Status.Phase)
	default:
		log.Info("Provisioning to PVC is in progress", "dvProgress", dv.Status.Progress, "dvPhase", dv.Status.Phase, "pvcPhase", pvc.Status.Phase)

		vi.Status.Progress = ds.diskService.GetProgress(dv, vi.Status.Progress, service.NewScaleOption(50, 100))
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

func (ds UploadDataSource) StoreToDVCR(ctx context.Context, vi *virtv2.VirtualImage) (reconcile.Result, error) {
	log, ctx := logger.GetDataSourceContext(ctx, "upload")

	condition, _ := service.GetCondition(vicondition.ReadyType, vi.Status.Conditions)
	defer func() { service.SetCondition(condition, &vi.Status.Conditions) }()

	supgen := supplements.NewGenerator(common.VIShortName, vi.Name, vi.Namespace, vi.UID)
	pod, err := ds.uploaderService.GetPod(ctx, supgen)
	if err != nil {
		return reconcile.Result{}, err
	}
	svc, err := ds.uploaderService.GetService(ctx, supgen)
	if err != nil {
		return reconcile.Result{}, err
	}
	ing, err := ds.uploaderService.GetIngress(ctx, supgen)
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

		err = ds.uploaderService.Unprotect(ctx, pod, svc, ing)
		if err != nil {
			return reconcile.Result{}, err
		}

		return CleanUpSupplements(ctx, vi, ds)
	case common.AnyTerminating(pod, svc, ing):
		vi.Status.Phase = virtv2.ImagePending

		log.Info("Cleaning up...")
	case pod == nil || svc == nil || ing == nil:
		envSettings := ds.getEnvSettings(vi, supgen)
		err = ds.uploaderService.Start(ctx, envSettings, vi, supgen, datasource.NewCABundleForVMI(vi.Spec.DataSource))
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

		log.Info("Create uploader pod...", "progress", vi.Status.Progress, "pod.phase", nil)

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

		condition.Status = metav1.ConditionTrue
		condition.Reason = vicondition.Ready
		condition.Message = ""

		vi.Status.Phase = virtv2.ImageReady
		vi.Status.Size = ds.statService.GetSize(pod)
		vi.Status.CDROM = ds.statService.GetCDROM(pod)
		vi.Status.Format = ds.statService.GetFormat(pod)
		vi.Status.Progress = "100%"
		vi.Status.Target.RegistryURL = ds.statService.GetDVCRImageName(pod)
		vi.Status.DownloadSpeed = ds.statService.GetDownloadSpeed(vi.GetUID(), pod)

		log.Info("Ready", "progress", vi.Status.Progress, "pod.phase", pod.Status.Phase)
	case ds.statService.IsUploadStarted(vi.GetUID(), pod):
		err = ds.statService.CheckPod(pod)
		if err != nil {
			return reconcile.Result{}, setPhaseConditionFromPodError(&condition, vi, err)
		}

		condition.Status = metav1.ConditionFalse
		condition.Reason = vicondition.Provisioning
		condition.Message = "Import is in the process of provisioning to DVCR."

		vi.Status.Phase = virtv2.ImageProvisioning
		vi.Status.Progress = ds.statService.GetProgress(vi.GetUID(), pod, vi.Status.Progress)
		vi.Status.Target.RegistryURL = ds.statService.GetDVCRImageName(pod)
		vi.Status.DownloadSpeed = ds.statService.GetDownloadSpeed(vi.GetUID(), pod)

		err = ds.uploaderService.Protect(ctx, pod, svc, ing)
		if err != nil {
			return reconcile.Result{}, err
		}

		log.Info("Provisioning...", "pod.phase", pod.Status.Phase)
	case ds.statService.IsUploaderReady(pod, svc, ing):
		condition.Status = metav1.ConditionFalse
		condition.Reason = vicondition.WaitForUserUpload
		condition.Message = "Waiting for the user upload."

		vi.Status.Phase = virtv2.ImageWaitForUserUpload
		vi.Status.Target.RegistryURL = ds.statService.GetDVCRImageName(pod)
		vi.Status.ImageUploadURLs = &virtv2.ImageUploadURLs{
			External:  ds.uploaderService.GetExternalURL(ctx, ing),
			InCluster: ds.uploaderService.GetInClusterURL(ctx, svc),
		}

		log.Info("Waiting for the user upload", "pod.phase", pod.Status.Phase)
	default:
		condition.Status = metav1.ConditionFalse
		condition.Reason = vicondition.ProvisioningNotStarted
		condition.Message = fmt.Sprintf("Waiting for the uploader %q to be ready to process the user's upload.", pod.Name)

		vi.Status.Phase = virtv2.ImagePending

		log.Info("Waiting for the uploader to be ready to process the user's upload", "pod.phase", pod.Status.Phase)
	}

	return reconcile.Result{Requeue: true}, nil
}

func (ds UploadDataSource) CleanUp(ctx context.Context, vi *virtv2.VirtualImage) (bool, error) {
	supgen := supplements.NewGenerator(common.VIShortName, vi.Name, vi.Namespace, vi.UID)

	importerRequeue, err := ds.uploaderService.CleanUp(ctx, supgen)
	if err != nil {
		return false, err
	}

	diskRequeue, err := ds.diskService.CleanUp(ctx, supgen)
	if err != nil {
		return false, err
	}

	return importerRequeue || diskRequeue, nil
}

func (ds UploadDataSource) Validate(_ context.Context, _ *virtv2.VirtualImage) error {
	return nil
}

func (ds UploadDataSource) getEnvSettings(vi *virtv2.VirtualImage, supgen *supplements.Generator) *uploader.Settings {
	var settings uploader.Settings

	uploader.ApplyDVCRDestinationSettings(
		&settings,
		ds.dvcrSettings,
		supgen,
		ds.dvcrSettings.RegistryImageForVI(vi),
	)

	return &settings
}

func (ds UploadDataSource) CleanUpSupplements(ctx context.Context, vi *virtv2.VirtualImage) (reconcile.Result, error) {
	supgen := supplements.NewGenerator(common.VIShortName, vi.Name, vi.Namespace, vi.UID)

	uploaderRequeue, err := ds.uploaderService.CleanUpSupplements(ctx, supgen)
	if err != nil {
		return reconcile.Result{}, err
	}

	diskRequeue, err := ds.diskService.CleanUpSupplements(ctx, supgen)
	if err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{Requeue: uploaderRequeue || diskRequeue}, nil
}

func (ds UploadDataSource) getPVCSize(pod *corev1.Pod) (resource.Quantity, error) {
	// Get size from the importer Pod to detect if specified PVC size is enough.
	unpackedSize, err := resource.ParseQuantity(ds.statService.GetSize(pod).UnpackedBytes)
	if err != nil {
		return resource.Quantity{}, fmt.Errorf("failed to parse unpacked bytes %s: %w", ds.statService.GetSize(pod).UnpackedBytes, err)
	}

	if unpackedSize.IsZero() {
		return resource.Quantity{}, errors.New("got zero unpacked size from data source")
	}

	return service.GetValidatedPVCSize(&unpackedSize, unpackedSize)
}

func (ds UploadDataSource) getSource(sup *supplements.Generator, dvcrSourceImageName string) *cdiv1.DataVolumeSource {
	// The image was preloaded from source into dvcr.
	// We can't use the same data source a second time, but we can set dvcr as the data source.
	// Use DV name for the Secret with DVCR auth and the ConfigMap with DVCR CA Bundle.
	url := cc.DockerRegistrySchemePrefix + dvcrSourceImageName
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
