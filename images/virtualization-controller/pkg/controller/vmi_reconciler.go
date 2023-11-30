package controller

import (
	"context"
	"fmt"
	"strconv"
	"time"

	corev1 "k8s.io/api/core/v1"
	cdiv1 "kubevirt.io/containerized-data-importer-api/pkg/apis/core/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	virtv2 "github.com/deckhouse/virtualization-controller/api/v2alpha1"
	"github.com/deckhouse/virtualization-controller/pkg/common"
	cc "github.com/deckhouse/virtualization-controller/pkg/controller/common"
	"github.com/deckhouse/virtualization-controller/pkg/controller/importer"
	"github.com/deckhouse/virtualization-controller/pkg/controller/monitoring"
	"github.com/deckhouse/virtualization-controller/pkg/controller/supplements"
	"github.com/deckhouse/virtualization-controller/pkg/controller/uploader"
	"github.com/deckhouse/virtualization-controller/pkg/controller/vmattachee"
	"github.com/deckhouse/virtualization-controller/pkg/dvcr"
	"github.com/deckhouse/virtualization-controller/pkg/sdk/framework/two_phase_reconciler"
	"github.com/deckhouse/virtualization-controller/pkg/util"
)

type VMIReconciler struct {
	*vmattachee.AttacheeReconciler[*virtv2.VirtualMachineImage, virtv2.VirtualMachineImageStatus]

	importerImage string
	uploaderImage string
	verbose       string
	pullPolicy    string
	dvcrSettings  *dvcr.Settings
}

func NewVMIReconciler(importerImage, uploaderImage, verbose, pullPolicy string, dvcrSettings *dvcr.Settings) *VMIReconciler {
	return &VMIReconciler{
		importerImage: importerImage,
		uploaderImage: uploaderImage,
		verbose:       verbose,
		pullPolicy:    pullPolicy,
		dvcrSettings:  dvcrSettings,
		AttacheeReconciler: vmattachee.NewAttacheeReconciler[
			*virtv2.VirtualMachineImage,
			virtv2.VirtualMachineImageStatus,
		]("vmi", true),
	}
}

func (r *VMIReconciler) SetupController(ctx context.Context, mgr manager.Manager, ctr controller.Controller) error {
	if err := ctr.Watch(source.Kind(mgr.GetCache(), &virtv2.VirtualMachineImage{}), &handler.EnqueueRequestForObject{},
		predicate.Funcs{
			CreateFunc: func(e event.CreateEvent) bool { return true },
			DeleteFunc: func(e event.DeleteEvent) bool { return true },
			UpdateFunc: func(e event.UpdateEvent) bool { return true },
		},
	); err != nil {
		return fmt.Errorf("error setting watch on VMI: %w", err)
	}

	if err := ctr.Watch(
		source.Kind(mgr.GetCache(), &cdiv1.DataVolume{}),
		handler.EnqueueRequestForOwner(
			mgr.GetScheme(),
			mgr.GetRESTMapper(),
			&virtv2.VirtualMachineImage{},
			handler.OnlyControllerOwner(),
		),
	); err != nil {
		return fmt.Errorf("error setting watch on DV: %w", err)
	}

	return r.AttacheeReconciler.SetupController(ctx, mgr, ctr)
}

// Sync starts an importer/uploader Pod or creates a DataVolume to import image into DVCR or into PVC.
// There are 3 modes of import:
// - Start and track importer/uploader Pod only (e.g. dataSource is HTTP and storage is ContainerRegistry).
// - Start importer/uploader Pod first and then create DataVolume (e.g. target size is unknown: dataSource is HTTP and storage is Kubernetes without specified size for PVC).
// - Create and track DataVolume only (e.g. dataSource is ClusterVirtualMachineImage and storage is Kubernetes).
func (r *VMIReconciler) Sync(ctx context.Context, _ reconcile.Request, state *VMIReconcilerState, opts two_phase_reconciler.ReconcilerOptions) error {
	if r.AttacheeReconciler.Sync(ctx, state.AttacheeState, opts) {
		return nil
	}

	switch {
	case state.IsDeletion():
		opts.Log.V(1).Info("Delete VMI, remove protective finalizers")
		return r.removeFinalizers(ctx, state, opts)
	case !state.IsProtected():
		// Set protective finalizer atomically.
		if controllerutil.AddFinalizer(state.VMI.Changed(), virtv2.FinalizerVMICleanup) {
			state.SetReconcilerResult(&reconcile.Result{Requeue: true})
			return nil
		}
	case state.IsReady():
		// Delete underlying importer/uploader Pod, Service and DataVolume and stop the reconcile process.
		return r.cleanup(ctx, state.VMI.Changed(), state.Client, state, opts)

	case state.ShouldTrackPod() && !state.IsPodComplete():
		// Start and track importer/uploader Pod.
		switch {
		case !state.IsPodInited():
			if err := r.verifyDataSource(state.VMI.Current(), state); err != nil {
				return err
			}
			opts.Log.V(1).Info("Update annotations with Pod name and namespace")
			r.initPodName(state)
			// Update annotations and status and restart reconcile to create an importer Pod.
			state.SetReconcilerResult(&reconcile.Result{Requeue: true})
			return nil
		case state.CanStartPod():
			// Create Pod using name and namespace from annotation.
			opts.Log.V(1).Info("Start new Pod for VMI")
			// Create importer/uploader pod, make sure the VMI owns it.
			if err := r.startPod(ctx, state, opts); err != nil {
				return err
			}
			// Requeue to wait until Pod become Running.
			state.SetReconcilerResult(&reconcile.Result{RequeueAfter: 2 * time.Second})
			return nil
		case state.Pod != nil:
			// Import is in progress, force a re-reconcile in 2 seconds to update status.
			opts.Log.V(2).Info("Requeue: wait until Pod is completed", "vmi.name", state.VMI.Current().Name)
			if err := r.ensurePodFinalizers(ctx, state, opts); err != nil {
				return err
			}
			state.SetReconcilerResult(&reconcile.Result{RequeueAfter: 2 * time.Second})
			return nil
		}
	case !state.ShouldTrackDataVolume() && state.ShouldTrackPod() && state.IsPodComplete():
		// Proceed to UpdateStatus and requeue to handle Ready state.
		state.SetReconcilerResult(&reconcile.Result{RequeueAfter: time.Second})
		return nil
	case state.ShouldTrackDataVolume() && (!state.ShouldTrackPod() || state.IsPodComplete()):
		// Start and track DataVolume.
		switch {
		case !state.HasDataVolumeAnno():
			opts.Log.V(1).Info("Update annotations with new DataVolume name")
			r.initDataVolumeName(state)
			// Update annotations and status and restart reconcile to create an importer/uploader Pod.
			state.SetReconcilerResult(&reconcile.Result{Requeue: true})
			return nil
		case state.CanCreateDataVolume():
			opts.Log.V(1).Info("Create DataVolume for VMI")

			if err := r.createDataVolume(ctx, state.VMI.Current(), state, opts); err != nil {
				return err
			}
			// Requeue to wait until Pod become Running.
			state.SetReconcilerResult(&reconcile.Result{RequeueAfter: 2 * time.Second})
			return nil
		case state.DV != nil:
			// Import is in progress, force a re-reconcile in 2 seconds to update status.
			opts.Log.V(2).Info("Requeue: wait until DataVolume is completed", "vmi.name", state.VMI.Current().Name)
			if err := r.ensureDVFinalizers(ctx, state, opts); err != nil {
				return err
			}
			state.SetReconcilerResult(&reconcile.Result{RequeueAfter: 2 * time.Second})
			return nil
		}
	}

	// Report unexpected state.
	details := fmt.Sprintf("vmi.Status.Phase='%s'", state.VMI.Current().Status.Phase)
	if state.Pod != nil {
		details += fmt.Sprintf(" pod.Name='%s' pod.Status.Phase='%s'", state.Pod.Name, state.Pod.Status.Phase)
	}
	if state.DV != nil {
		details += fmt.Sprintf(" dv.Name='%s' dv.Status.Phase='%s'", state.DV.Name, state.DV.Status.Phase)
	}
	if state.PVC != nil {
		details += fmt.Sprintf(" pvc.Name='%s' pvc.Status.Phase='%s'", state.PVC.Name, state.PVC.Status.Phase)
	}
	opts.Recorder.Event(state.VMI.Current(), corev1.EventTypeWarning, virtv2.ReasonErrUnknownState, fmt.Sprintf("VMI has unexpected state, recreate it to start import again. %s", details))

	return nil
}

func (r *VMIReconciler) UpdateStatus(_ context.Context, _ reconcile.Request, state *VMIReconcilerState, opts two_phase_reconciler.ReconcilerOptions) error {
	opts.Log.V(2).Info("Update VMI status", "vmi.name", state.VMI.Current().GetName())

	// Do nothing if object is being deleted as any update will lead to en error.
	if state.IsDeletion() {
		return nil
	}

	// Record event if importer/uploader Pod has error.
	// TODO set Failed status if Pod restarts are greater than some threshold?
	if state.Pod != nil && len(state.Pod.Status.ContainerStatuses) > 0 {
		if state.Pod.Status.ContainerStatuses[0].LastTerminationState.Terminated != nil &&
			state.Pod.Status.ContainerStatuses[0].LastTerminationState.Terminated.ExitCode > 0 {
			opts.Recorder.Event(state.VMI.Current(), corev1.EventTypeWarning, virtv2.ReasonErrImportFailed, fmt.Sprintf("pod phase '%s', message '%s'", state.Pod.Status.Phase, state.Pod.Status.ContainerStatuses[0].LastTerminationState.Terminated.Message))
		}
	}

	vmiStatus := state.VMI.Current().Status.DeepCopy()

	if vmiStatus.Phase != virtv2.ImageReady {
		vmiStatus.ImportDuration = time.Since(state.VMI.Current().CreationTimestamp.Time).Truncate(time.Second).String()
	}

	switch {
	case vmiStatus.Phase == "":
		vmiStatus.Phase = virtv2.ImagePending
		state.SetReconcilerResult(&reconcile.Result{Requeue: true})
		if err := r.verifyDataSource(state.VMI.Current(), state); err != nil {
			vmiStatus.FailureReason = FailureReasonCannotBeProcessed
			vmiStatus.FailureMessage = fmt.Sprintf("DataSource is invalid. %s", err)
		}
	case state.IsReady():
		// No need to update status.
		break
	case state.ShouldTrackPod() && state.IsPodInProgress():
		opts.Log.V(2).Info("Fetch progress from Pod", "vmi.name", state.VMI.Current().GetName())

		vmiStatus.Phase = virtv2.ImageProvisioning
		t := state.VMI.Current().Spec.DataSource.Type
		if t == virtv2.DataSourceTypeUpload &&
			vmiStatus.UploadCommand == "" &&
			state.Service != nil &&
			len(state.Service.Spec.Ports) > 0 {
			vmiStatus.UploadCommand = fmt.Sprintf(
				"curl -X POST --data-binary @example.iso http://%s:%d/v1beta1/upload",
				state.Service.Spec.ClusterIP,
				state.Service.Spec.Ports[0].Port,
			)
		}
		var progress *monitoring.ImportProgress
		if t != virtv2.DataSourceTypeClusterVirtualMachineImage &&
			t != virtv2.DataSourceTypeVirtualMachineImage {
			progress, err := monitoring.GetImportProgressFromPod(string(state.VMI.Current().GetUID()), state.Pod)
			if err != nil {
				opts.Recorder.Event(state.VMI.Current(), corev1.EventTypeWarning, virtv2.ReasonErrGetProgressFailed, "Error fetching progress metrics from Pod "+err.Error())
				return err
			}
			if progress != nil {
				opts.Log.V(2).Info("Got progress", "vmi.name", state.VMI.Current().Name, "progress", progress.Progress(), "speed", progress.AvgSpeed(), "progress.raw", progress.ProgressRaw(), "speed.raw", progress.AvgSpeedRaw())
				// map 0-100% to 0-50%.
				progressPct := progress.Progress()
				if state.ShouldTrackDataVolume() {
					progressPct = common.ScalePercentage(progressPct, 0, 50.0)
				}
				vmiStatus.Progress = progressPct
				vmiStatus.DownloadSpeed.Avg = progress.AvgSpeed()
				vmiStatus.DownloadSpeed.AvgBytes = strconv.FormatUint(progress.AvgSpeedRaw(), 10)
				vmiStatus.DownloadSpeed.Current = progress.CurSpeed()
				vmiStatus.DownloadSpeed.CurrentBytes = strconv.FormatUint(progress.CurSpeedRaw(), 10)
			}
		}

		// Set VMI phase.
		if state.VMI.Current().Spec.DataSource.Type == virtv2.DataSourceTypeUpload && (progress == nil || progress.ProgressRaw() == 0) {
			vmiStatus.Phase = virtv2.ImageWaitForUserUpload
		} else {
			vmiStatus.Phase = virtv2.ImageProvisioning
		}
	case !state.ShouldTrackDataVolume() && state.ShouldTrackPod() && state.IsPodComplete():
		vmiStatus.Phase = virtv2.ImageReady

		opts.Log.V(1).Info("Import completed successfully")

		opts.Recorder.Event(state.VMI.Current(), corev1.EventTypeNormal, virtv2.ReasonImportSucceeded, "Import Successful")

		// Cleanup
		vmiStatus.DownloadSpeed.Current = ""
		vmiStatus.DownloadSpeed.CurrentBytes = ""

		finalReport, err := monitoring.GetFinalReportFromPod(state.Pod)
		if err != nil {
			opts.Log.Error(err, "parsing final report", "vmi.name", state.VMI.Current().Name)
		}
		if finalReport != nil {
			vmiStatus.DownloadSpeed.Avg = finalReport.GetAverageSpeed()
			vmiStatus.DownloadSpeed.AvgBytes = strconv.FormatUint(finalReport.GetAverageSpeedRaw(), 10)
			vmiStatus.Size.Stored = finalReport.StoredSize()
			vmiStatus.Size.StoredBytes = strconv.FormatUint(finalReport.StoredSizeBytes, 10)
			vmiStatus.Size.Unpacked = finalReport.UnpackedSize()
			vmiStatus.Size.UnpackedBytes = strconv.FormatUint(finalReport.UnpackedSizeBytes, 10)

			switch state.VMI.Current().Spec.DataSource.Type {
			case virtv2.DataSourceTypeClusterVirtualMachineImage:
				vmiStatus.Format = state.DVCRDataSource.CVMI.Status.Format
			case virtv2.DataSourceTypeVirtualMachineImage:
				vmiStatus.Format = state.DVCRDataSource.VMI.Status.Format
			default:
				vmiStatus.Format = finalReport.Format
			}
		}
		// Set target image name the same way as for the importer/uploader Pod.
		dvcrDestImageName := r.dvcrSettings.RegistryImageForVMI(state.VMI.Current().Name, state.VMI.Current().Namespace)
		vmiStatus.Target.RegistryURL = dvcrDestImageName
	case state.ShouldTrackDataVolume() && state.IsDataVolumeInProgress():
		// Set phase from DataVolume resource.
		vmiStatus.Phase = MapDataVolumePhaseToVMIPhase(state.DV.Status.Phase)

		// Download speed is not available from DataVolume.
		vmiStatus.DownloadSpeed.Current = ""
		vmiStatus.DownloadSpeed.CurrentBytes = ""

		// Copy progress from DataVolume.
		// map 0-100% to 50%-100%.
		dvProgress := string(state.DV.Status.Progress)

		opts.Log.V(2).Info("Got DataVolume progress", "progress", dvProgress)

		if dvProgress != "N/A" && dvProgress != "" {
			vmiStatus.Progress = common.ScalePercentage(dvProgress, 50.0, 100.0)
		}

		// Copy capacity from PVC.
		if state.PVC != nil && state.PVC.Status.Phase == corev1.ClaimBound {
			vmiStatus.Capacity = util.GetPointer(state.PVC.Status.Capacity[corev1.ResourceStorage]).String()
		}
	case state.ShouldTrackDataVolume() && state.IsDataVolumeComplete():
		opts.Recorder.Event(state.VMI.Current(), corev1.EventTypeNormal, virtv2.ReasonImportSucceededToPVC, "Import Successful")
		opts.Log.V(1).Info("Import completed successfully")
		vmiStatus.Phase = virtv2.ImageReady

		// Cleanup.
		vmiStatus.DownloadSpeed.Current = ""
		vmiStatus.DownloadSpeed.CurrentBytes = ""
		// PVC name equals to the DataVolume name.
		dv := state.Supplements.DataVolume()
		vmiStatus.Target.PersistentVolumeClaimName = dv.Name
	}

	state.VMI.Changed().Status = *vmiStatus

	return nil
}

func MapDataVolumePhaseToVMIPhase(phase cdiv1.DataVolumePhase) virtv2.ImagePhase {
	switch phase {
	case cdiv1.PhaseUnset, cdiv1.Unknown, cdiv1.Pending:
		return virtv2.ImagePending
	case cdiv1.WaitForFirstConsumer, cdiv1.PVCBound,
		cdiv1.ImportScheduled, cdiv1.CloneScheduled, cdiv1.UploadScheduled,
		cdiv1.ImportInProgress, cdiv1.CloneInProgress,
		cdiv1.SnapshotForSmartCloneInProgress, cdiv1.SmartClonePVCInProgress,
		cdiv1.CSICloneInProgress,
		cdiv1.CloneFromSnapshotSourceInProgress,
		cdiv1.Paused:
		return virtv2.ImageProvisioning
	case cdiv1.Succeeded:
		return virtv2.ImageReady
	case cdiv1.Failed:
		return virtv2.ImageFailed
	default:
		panic(fmt.Sprintf("unexpected DataVolume phase %q, please report a bug", phase))
	}
}

// ensurePodFinalizers adds protective finalizers on importer/uploader Pod and Service dependencies.
func (r *VMIReconciler) ensurePodFinalizers(ctx context.Context, state *VMIReconcilerState, opts two_phase_reconciler.ReconcilerOptions) error {
	if state.Pod != nil {
		if controllerutil.AddFinalizer(state.Pod, virtv2.FinalizerPodProtection) {
			if err := opts.Client.Update(ctx, state.Pod); err != nil {
				return fmt.Errorf("error setting finalizer on a Pod %q: %w", state.Pod.Name, err)
			}
		}
	}
	if state.Service != nil {
		if controllerutil.AddFinalizer(state.Service, virtv2.FinalizerServiceProtection) {
			if err := opts.Client.Update(ctx, state.Service); err != nil {
				return fmt.Errorf("error setting finalizer on a Service %q: %w", state.Service.Name, err)
			}
		}
	}

	return nil
}

// ensureDVFinalizers adds protective finalizers on DataVolume, PersistentVolumeClaim and PersistentVolume dependencies.
func (r *VMIReconciler) ensureDVFinalizers(ctx context.Context, state *VMIReconcilerState, opts two_phase_reconciler.ReconcilerOptions) error {
	if state.DV != nil {
		// Ensure DV finalizer is set in case DV was created manually (take ownership of already existing object)
		if controllerutil.AddFinalizer(state.DV, virtv2.FinalizerDVProtection) {
			if err := opts.Client.Update(ctx, state.DV); err != nil {
				return fmt.Errorf("error setting finalizer on a DV %q: %w", state.DV.Name, err)
			}
		}
	}
	if state.PVC != nil {
		if controllerutil.AddFinalizer(state.PVC, virtv2.FinalizerPVCProtection) {
			if err := opts.Client.Update(ctx, state.PVC); err != nil {
				return fmt.Errorf("error setting finalizer on a PVC %q: %w", state.PVC.Name, err)
			}
		}
	}
	if state.PV != nil {
		if controllerutil.AddFinalizer(state.PV, virtv2.FinalizerPVProtection) {
			if err := opts.Client.Update(ctx, state.PV); err != nil {
				return fmt.Errorf("error setting finalizer on a PV %q: %w", state.PV.Name, err)
			}
		}
	}

	return nil
}

// removeFinalizers removes protective finalizers on Pod, Service, DataVolume, PersistentVolumeClaim and PersistentVolume dependencies.
func (r *VMIReconciler) removeFinalizers(ctx context.Context, state *VMIReconcilerState, opts two_phase_reconciler.ReconcilerOptions) error {
	if state.Pod != nil {
		if controllerutil.RemoveFinalizer(state.Pod, virtv2.FinalizerPodProtection) {
			if err := opts.Client.Update(ctx, state.Pod); err != nil {
				return fmt.Errorf("unable to remove Pod %q finalizer %q: %w", state.Pod.Name, virtv2.FinalizerPodProtection, err)
			}
		}
	}
	if state.Service != nil {
		if controllerutil.RemoveFinalizer(state.Service, virtv2.FinalizerServiceProtection) {
			if err := opts.Client.Update(ctx, state.Service); err != nil {
				return fmt.Errorf("unable to remove Service %q finalizer %q: %w", state.Service.Name, virtv2.FinalizerServiceProtection, err)
			}
		}
	}
	if state.PV != nil {
		if controllerutil.RemoveFinalizer(state.PV, virtv2.FinalizerPVProtection) {
			if err := opts.Client.Update(ctx, state.PV); err != nil {
				return fmt.Errorf("unable to remove PV %q finalizer %q: %w", state.PV.Name, virtv2.FinalizerPVProtection, err)
			}
		}
	}
	if state.PVC != nil {
		if controllerutil.RemoveFinalizer(state.PVC, virtv2.FinalizerPVCProtection) {
			if err := opts.Client.Update(ctx, state.PVC); err != nil {
				return fmt.Errorf("unable to remove PVC %q finalizer %q: %w", state.PVC.Name, virtv2.FinalizerPVCProtection, err)
			}
		}
	}
	if state.DV != nil {
		if controllerutil.RemoveFinalizer(state.DV, virtv2.FinalizerDVProtection) {
			if err := opts.Client.Update(ctx, state.DV); err != nil {
				return fmt.Errorf("unable to remove DV %q finalizer %q: %w", state.DV.Name, virtv2.FinalizerDVProtection, err)
			}
		}
	}
	controllerutil.RemoveFinalizer(state.VMI.Changed(), virtv2.FinalizerVMICleanup)

	return nil
}

func (r *VMIReconciler) verifyDataSource(vmi *virtv2.VirtualMachineImage, state *VMIReconcilerState) error {
	return VerifyDVCRDataSources(vmi.Spec.DataSource, state.DVCRDataSource)
}

// initPodName creates new name and update it in the annotation.
func (r *VMIReconciler) initPodName(state *VMIReconcilerState) {
	vmi := state.VMI.Changed()

	switch vmi.Spec.DataSource.Type {
	case virtv2.DataSourceTypeUpload:
		uploaderPod := state.Supplements.UploaderPod()
		cc.AddAnnotation(vmi, cc.AnnUploadPodName, uploaderPod.Name)
	default:
		importerPod := state.Supplements.ImporterPod()
		cc.AddAnnotation(vmi, cc.AnnImportPodName, importerPod.Name)
	}
}

func (r *VMIReconciler) startPod(ctx context.Context, state *VMIReconcilerState, opts two_phase_reconciler.ReconcilerOptions) error {
	switch state.VMI.Current().Spec.DataSource.Type {
	case virtv2.DataSourceTypeUpload:
		if err := r.startUploaderPod(ctx, state, opts); err != nil {
			return err
		}

		if err := r.startUploaderService(ctx, state, opts); err != nil {
			return err
		}
	default:
		if err := r.startImporterPod(ctx, state, opts); err != nil {
			return err
		}
	}

	return nil
}

// initDataVolumeName creates new DataVolume name and update it in the annotation.
func (r *VMIReconciler) initDataVolumeName(state *VMIReconcilerState) {
	vmi := state.VMI.Changed()

	// Prevent DataVolume name regeneration.
	if _, hasKey := vmi.Annotations[cc.AnnVMIDataVolume]; hasKey {
		return
	}

	dv := state.Supplements.DataVolume()
	cc.AddAnnotation(vmi, cc.AnnVMIDataVolume, dv.Name)
}

func (r *VMIReconciler) cleanup(
	ctx context.Context,
	vmi *virtv2.VirtualMachineImage,
	client client.Client,
	state *VMIReconcilerState,
	opts two_phase_reconciler.ReconcilerOptions,
) error {
	opts.Log.V(1).Info("Import done, cleanup")
	if state.DV != nil {
		err := supplements.CleanupForDataVolume(ctx, client, state.Supplements, r.dvcrSettings)
		if err != nil {
			return fmt.Errorf("cleanup supplements for DataVolume: %w", err)
		}
		// TODO(future): take ownership on PVC and delete DataVolume.
		// if err := client.Delete(ctx, state.DV); err != nil {
		// 	return fmt.Errorf("cleanup DataVolume: %w", err)
		// }
	}

	if state.Pod != nil && cc.ShouldDeletePod(state.VMI.Current()) {
		switch vmi.Spec.DataSource.Type {
		case virtv2.DataSourceTypeUpload:
			if err := uploader.CleanupService(ctx, client, state.Service); err != nil {
				return err
			}

			if err := uploader.CleanupPod(ctx, client, state.Pod); err != nil {
				return err
			}
		default:
			if err := r.removeFinalizers(ctx, state, opts); err != nil {
				return err
			}
			if err := importer.CleanupPod(ctx, client, state.Pod); err != nil {
				return err
			}
		}
	}

	return nil
}