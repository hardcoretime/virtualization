package controller

import (
	"context"
	"errors"
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
	"github.com/deckhouse/virtualization-controller/pkg/sdk/framework/helper"
	"github.com/deckhouse/virtualization-controller/pkg/sdk/framework/two_phase_reconciler"
	"github.com/deckhouse/virtualization-controller/pkg/util"
)

type VMDReconciler struct {
	*vmattachee.AttacheeReconciler[*virtv2.VirtualMachineDisk, virtv2.VirtualMachineDiskStatus]

	importerImage string
	uploaderImage string
	verbose       string
	pullPolicy    string
	dvcrSettings  *dvcr.Settings
}

func NewVMDReconciler(importerImage, uploaderImage, verbose, pullPolicy string, dvcrSettings *dvcr.Settings) *VMDReconciler {
	return &VMDReconciler{
		importerImage: importerImage,
		uploaderImage: uploaderImage,
		verbose:       verbose,
		pullPolicy:    pullPolicy,
		dvcrSettings:  dvcrSettings,
		AttacheeReconciler: vmattachee.NewAttacheeReconciler[
			*virtv2.VirtualMachineDisk,
			virtv2.VirtualMachineDiskStatus,
		]("vmd", true),
	}
}

func (r *VMDReconciler) SetupController(ctx context.Context, mgr manager.Manager, ctr controller.Controller) error {
	if err := ctr.Watch(source.Kind(mgr.GetCache(), &virtv2.VirtualMachineDisk{}), &handler.EnqueueRequestForObject{},
		predicate.Funcs{
			CreateFunc: func(e event.CreateEvent) bool { return true },
			DeleteFunc: func(e event.DeleteEvent) bool { return true },
			UpdateFunc: func(e event.UpdateEvent) bool { return true },
		},
	); err != nil {
		return fmt.Errorf("error setting watch on VMD: %w", err)
	}

	if err := ctr.Watch(
		source.Kind(mgr.GetCache(), &cdiv1.DataVolume{}),
		handler.EnqueueRequestForOwner(
			mgr.GetScheme(),
			mgr.GetRESTMapper(),
			&virtv2.VirtualMachineDisk{},
			handler.OnlyControllerOwner(),
		),
	); err != nil {
		return fmt.Errorf("error setting watch on DV: %w", err)
	}

	return r.AttacheeReconciler.SetupController(ctx, mgr, ctr)
}

// Sync starts an importer Pod and creates a DataVolume to import image into PVC.
func (r *VMDReconciler) Sync(ctx context.Context, _ reconcile.Request, state *VMDReconcilerState, opts two_phase_reconciler.ReconcilerOptions) error {
	log := opts.Log.WithValues("vmd.name", state.VMD.Current().GetName())

	log.V(2).Info("Sync VMD")

	if r.AttacheeReconciler.Sync(ctx, state.AttacheeState, opts) {
		return nil
	}

	switch {
	case state.IsDeletion():
		log.V(1).Info("Delete VMD, remove protective finalizers")
		return r.cleanupOnDeletion(ctx, state, opts)
	case !state.IsProtected():
		// Set protective finalizer atomically.
		if controllerutil.AddFinalizer(state.VMD.Changed(), virtv2.FinalizerVMDCleanup) {
			state.SetReconcilerResult(&reconcile.Result{Requeue: true})
			return nil
		}
	case state.IsReady():
		opts.Log.Info("VMD is ready: cleanup underlying resources")
		// Delete underlying importer/uploader Pod, Service and DataVolume and stop the reconcile process.
		if err := r.cleanup(ctx, state.VMD.Changed(), state.Client, state); err != nil {
			return err
		}

		if state.PVC == nil {
			return errors.New("pvc not found, please report a bug")
		}

		oldSize := state.PVC.Spec.Resources.Requests[corev1.ResourceStorage]
		newSize := state.VMD.Current().Spec.PersistentVolumeClaim.Size

		if newSize == nil || newSize.Cmp(oldSize) == -1 {
			return nil
		}

		if !newSize.Equal(oldSize) {
			opts.Log.Info("Increase PVC size", "oldPVCSize", oldSize.String(), "newPVCSize", newSize.String())

			state.PVC.Spec.Resources.Requests[corev1.ResourceStorage] = *newSize

			err := opts.Client.Update(ctx, state.PVC)
			if err != nil {
				return fmt.Errorf("failed to increase pvc size: %w", err)
			}
		}

		if !newSize.Equal(state.PVC.Status.Capacity[corev1.ResourceStorage]) {
			opts.Log.Info("PVC is in a process of increasing: wait for the PVC to be increased")
			state.SetReconcilerResult(&reconcile.Result{RequeueAfter: 2 * time.Second})
		}

		return nil
	// First phase: import to DVCR.
	case state.ShouldTrackPod() && !state.IsPodComplete():
		// Start and track importer/uploader Pod.
		switch {
		case !state.IsPodInited():
			log.V(1).Info("New VMI observed, update annotations with Pod name and namespace")
			r.initPodName(state)
			// Update annotations and status and restart reconcile to create an importer/uploader Pod.
			state.SetReconcilerResult(&reconcile.Result{Requeue: true})
			return nil
		case state.CanStartPod():
			// Create Pod using name and namespace from annotation.
			log.V(1).Info("Start new Pod for VMD")
			// Create importer/uploader pod, make sure the VMD owns it.
			if err := r.startPod(ctx, state, opts); err != nil {
				return err
			}
			// Requeue to wait until Pod become Running.
			state.SetReconcilerResult(&reconcile.Result{RequeueAfter: 2 * time.Second})
			return nil
		case state.Pod != nil:
			// Import is in progress, force a re-reconcile in 2 seconds to update status.
			log.V(2).Info("Requeue: wait until Pod is completed", "vmd.name", state.VMD.Current().Name)
			if err := r.ensurePodFinalizers(ctx, state, opts); err != nil {
				return err
			}
			state.SetReconcilerResult(&reconcile.Result{RequeueAfter: 2 * time.Second})
			return nil
		}
	// Second phase: import to PVC.
	case state.ShouldTrackDataVolume() && (!state.ShouldTrackPod() || state.IsPodComplete()):
		// Start and track DataVolume.
		switch {
		case !state.HasDataVolumeAnno():
			if state.ShouldTrackPod() {
				finalReport, err := monitoring.GetFinalReportFromPod(state.Pod)
				if err != nil {
					return err
				}

				if finalReport == nil {
					return errors.New("empty final report")
				}

				if finalReport.ErrMessage != "" {
					return nil
				}
			}

			log.V(1).Info("Update annotations with new DataVolume name")
			r.initDataVolumeName(state)
			// Update annotations and status and restart reconcile to create a DV.
			state.SetReconcilerResult(&reconcile.Result{Requeue: true})
			return nil
		case state.CanCreateDataVolume():
			log.V(1).Info("Create DataVolume for VMD")

			err := r.createDataVolume(ctx, state.VMD.Current(), state, opts)
			if err != nil {
				if !errors.Is(err, ErrDataSourceNotReady) {
					return err
				}

				log.V(1).Info("Wait for the data source to be ready", "err", err)
			}

			// Requeue to wait until Pod become Running.
			state.SetReconcilerResult(&reconcile.Result{RequeueAfter: 2 * time.Second})
			return nil
		case state.DV != nil:
			// Import is in progress, force a re-reconcile in 2 seconds to update status.
			log.V(2).Info("Requeue: wait until DataVolume is completed", "vmd.name", state.VMD.Current().Name)
			if err := r.ensureDVFinalizers(ctx, state, opts); err != nil {
				return err
			}
			state.SetReconcilerResult(&reconcile.Result{RequeueAfter: 2 * time.Second})
			return nil
		}
	}

	// Report unexpected state.
	details := fmt.Sprintf("vmd.Status.Phase='%s'", state.VMD.Current().Status.Phase)
	if state.Pod != nil {
		details += fmt.Sprintf(" pod.Name='%s' pod.Status.Phase='%s'", state.Pod.Name, state.Pod.Status.Phase)
	}
	if state.DV != nil {
		details += fmt.Sprintf(" dv.Name='%s' dv.Status.Phase='%s'", state.DV.Name, state.DV.Status.Phase)
	}
	if state.PVC != nil {
		details += fmt.Sprintf(" pvc.Name='%s' pvc.Status.Phase='%s'", state.PVC.Name, state.PVC.Status.Phase)
	}
	opts.Recorder.Event(state.VMD.Current(), corev1.EventTypeWarning, virtv2.ReasonErrUnknownState, fmt.Sprintf("VMD has unexpected state, recreate it to start import again. %s", details))

	return nil
}

func (r *VMDReconciler) UpdateStatus(_ context.Context, _ reconcile.Request, state *VMDReconcilerState, opts two_phase_reconciler.ReconcilerOptions) error {
	log := opts.Log.WithValues("vmd.name", state.VMD.Current().GetName())

	log.V(2).Info("Update VMD status")

	// Do nothing if object is being deleted as any update will lead to en error.
	if state.IsDeletion() {
		return nil
	}

	// Record event if importer/uploader Pod has error.
	// TODO set Failed status if Pod restarts are greater than some threshold?
	if state.Pod != nil && len(state.Pod.Status.ContainerStatuses) > 0 {
		if state.Pod.Status.ContainerStatuses[0].LastTerminationState.Terminated != nil &&
			state.Pod.Status.ContainerStatuses[0].LastTerminationState.Terminated.ExitCode > 0 {
			opts.Recorder.Event(state.VMD.Current(), corev1.EventTypeWarning, virtv2.ReasonErrImportFailed, fmt.Sprintf("importer pod phase '%s', message '%s'", state.Pod.Status.Phase, state.Pod.Status.ContainerStatuses[0].LastTerminationState.Terminated.Message))
		}
	}

	vmdStatus := state.VMD.Current().Status.DeepCopy()

	if vmdStatus.Phase != virtv2.DiskReady {
		vmdStatus.ImportDuration = time.Since(state.VMD.Current().CreationTimestamp.Time).Truncate(time.Second).String()
	}

	if vmdStatus.Progress == "" {
		vmdStatus.Progress = "0%"
	}

	switch {
	case vmdStatus.Phase == "":
		vmdStatus.Phase = virtv2.DiskPending
		state.SetReconcilerResult(&reconcile.Result{Requeue: true})
	case state.IsReady():
		vmdStatus.Capacity = util.GetPointer(state.PVC.Status.Capacity[corev1.ResourceStorage]).String()
	case state.ShouldTrackPod() && state.IsPodRunning():
		log.V(2).Info("Fetch progress from Pod")

		// Set statue UploadCommand if necessary.
		if state.VMD.Current().Spec.DataSource.Type == virtv2.DataSourceTypeUpload &&
			vmdStatus.UploadCommand == "" &&
			state.Service != nil &&
			len(state.Service.Spec.Ports) > 0 {
			vmdStatus.UploadCommand = fmt.Sprintf(
				"curl -X POST http://%s:%d/v1beta1/upload -T example.iso",
				state.Service.Spec.ClusterIP,
				state.Service.Spec.Ports[0].Port,
			)
		}

		progress, err := monitoring.GetImportProgressFromPod(string(state.VMD.Current().GetUID()), state.Pod)
		if err != nil {
			opts.Recorder.Event(state.VMD.Current(), corev1.EventTypeWarning, virtv2.ReasonErrGetProgressFailed, "Error fetching progress metrics from Pod "+err.Error())
			return err
		}
		if progress != nil {
			log.V(2).Info("Got Pod progress", "progress", progress.Progress(), "speed", progress.AvgSpeed(), "progress.raw", progress.ProgressRaw(), "speed.raw", progress.AvgSpeedRaw())
			// map 0-100% to 0-50%.
			progressPct := progress.Progress()
			if state.ShouldTrackDataVolume() {
				progressPct = common.ScalePercentage(progressPct, 0, 50.0)
			}
			vmdStatus.Progress = progressPct
			vmdStatus.DownloadSpeed.Avg = progress.AvgSpeed()
			vmdStatus.DownloadSpeed.AvgBytes = strconv.FormatUint(progress.AvgSpeedRaw(), 10)
			vmdStatus.DownloadSpeed.Current = progress.CurSpeed()
			vmdStatus.DownloadSpeed.CurrentBytes = strconv.FormatUint(progress.CurSpeedRaw(), 10)
		}

		// Set VMD phase.
		if state.VMD.Current().Spec.DataSource.Type == virtv2.DataSourceTypeUpload && (progress == nil || progress.ProgressRaw() == 0) {
			vmdStatus.Phase = virtv2.DiskWaitForUserUpload
		} else {
			vmdStatus.Phase = virtv2.DiskProvisioning
		}

	case state.IsPodComplete() && !state.HasDataVolumeAnno():
		finalReport, err := monitoring.GetFinalReportFromPod(state.Pod)
		if err != nil {
			return err
		}

		if finalReport == nil {
			err = errors.New("empty final report")
			log.Error(err, "Failed to process final report")
			return err
		}

		// Cleanup.
		vmdStatus.DownloadSpeed.Current = ""
		vmdStatus.DownloadSpeed.CurrentBytes = ""

		if finalReport.ErrMessage != "" {
			vmdStatus.Phase = virtv2.DiskFailed
			vmdStatus.FailureReason = virtv2.ReasonErrImportFailed
			vmdStatus.FailureMessage = finalReport.ErrMessage
			break
		}

		vmdStatus.DownloadSpeed.Avg = finalReport.GetAverageSpeed()
		vmdStatus.DownloadSpeed.AvgBytes = strconv.FormatUint(finalReport.GetAverageSpeedRaw(), 10)

	case state.ShouldTrackDataVolume() && state.IsDataVolumeInProgress():
		// Set phase from DataVolume resource.
		vmdStatus.Phase = MapDataVolumePhaseToVMDPhase(state.DV.Status.Phase)

		// Download speed is not available from DataVolume.
		vmdStatus.DownloadSpeed.Current = ""
		vmdStatus.DownloadSpeed.CurrentBytes = ""

		// Copy progress from DataVolume.
		// map 0-100% to 50%-100%.
		dvProgress := string(state.DV.Status.Progress)

		opts.Log.V(2).Info("Got DataVolume progress", "progress", dvProgress)

		if dvProgress != "N/A" && dvProgress != "" {
			vmdStatus.Progress = common.ScalePercentage(dvProgress, 50.0, 100.0)
		} else {
			vmdStatus.Progress = "50%"
		}

		// Copy capacity from PVC.
		if state.PVC != nil && state.PVC.Status.Phase == corev1.ClaimBound {
			vmdStatus.Capacity = util.GetPointer(state.PVC.Status.Capacity[corev1.ResourceStorage]).String()
		}
	case state.ShouldTrackDataVolume() && state.IsDataVolumeComplete():
		if state.PVC == nil {
			return errors.New("pvc not found, please report a bug")
		}

		if state.PVC.Status.Phase != corev1.ClaimBound {
			log.V(1).Info("Wait for the PVC to enter the Bound phase")
			state.SetReconcilerResult(&reconcile.Result{RequeueAfter: 2 * time.Second})
			break
		}

		log.V(1).Info("Import completed successfully")

		vmdStatus.Phase = virtv2.DiskReady
		vmdStatus.Progress = "100%"

		opts.Recorder.Event(state.VMD.Current(), corev1.EventTypeNormal, virtv2.ReasonImportSucceeded, "Successfully imported")

		// Cleanup download speed and set average from importer/uploader Pod if any.
		vmdStatus.DownloadSpeed.Current = ""
		vmdStatus.DownloadSpeed.CurrentBytes = ""
		vmdStatus.DownloadSpeed.Avg = ""
		vmdStatus.DownloadSpeed.AvgBytes = ""

		if state.Pod != nil {
			finalReport, err := monitoring.GetFinalReportFromPod(state.Pod)
			if err != nil {
				return err
			}

			if finalReport != nil {
				vmdStatus.DownloadSpeed.Avg = finalReport.GetAverageSpeed()
				vmdStatus.DownloadSpeed.AvgBytes = strconv.FormatUint(finalReport.GetAverageSpeedRaw(), 10)
			}
		}

		// PVC name is the same as the DataVolume name.
		vmdStatus.Target.PersistentVolumeClaimName = state.VMD.Current().Annotations[cc.AnnVMDDataVolume]

		// Copy capacity from PVC if IsDataVolumeInProgress was very quick.
		vmdStatus.Capacity = util.GetPointer(state.PVC.Status.Capacity[corev1.ResourceStorage]).String()
	}

	state.VMD.Changed().Status = *vmdStatus

	return nil
}

func MapDataVolumePhaseToVMDPhase(phase cdiv1.DataVolumePhase) virtv2.DiskPhase {
	switch phase {
	case cdiv1.PhaseUnset, cdiv1.Unknown, cdiv1.Pending:
		return virtv2.DiskPending
	case cdiv1.WaitForFirstConsumer, cdiv1.PVCBound,
		cdiv1.ImportScheduled, cdiv1.CloneScheduled, cdiv1.UploadScheduled,
		cdiv1.ImportInProgress, cdiv1.CloneInProgress,
		cdiv1.SnapshotForSmartCloneInProgress, cdiv1.SmartClonePVCInProgress,
		cdiv1.CSICloneInProgress,
		cdiv1.CloneFromSnapshotSourceInProgress,
		cdiv1.Paused:
		return virtv2.DiskProvisioning
	case cdiv1.Succeeded:
		return virtv2.DiskReady
	case cdiv1.Failed:
		return virtv2.DiskFailed
	default:
		panic(fmt.Sprintf("unexpected DataVolume phase %q, please report a bug", phase))
	}
}

// ensurePodFinalizers adds protective finalizers on importer/uploader Pod and Service dependencies.
func (r *VMDReconciler) ensurePodFinalizers(ctx context.Context, state *VMDReconcilerState, opts two_phase_reconciler.ReconcilerOptions) error {
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
func (r *VMDReconciler) ensureDVFinalizers(ctx context.Context, state *VMDReconcilerState, opts two_phase_reconciler.ReconcilerOptions) error {
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

func (r *VMDReconciler) ShouldDeleteChildResources(state *VMDReconcilerState) bool {
	return state.Pod != nil || state.Service != nil || state.PV != nil || state.PVC != nil || state.DV != nil
}

// removeFinalizerChildResources removes protective finalizers on Pod, Service, DataVolume, PersistentVolumeClaim and PersistentVolume dependencies.
func (r *VMDReconciler) removeFinalizerChildResources(ctx context.Context, state *VMDReconcilerState, opts two_phase_reconciler.ReconcilerOptions) error {
	if state.Pod != nil && controllerutil.RemoveFinalizer(state.Pod, virtv2.FinalizerPodProtection) {
		if err := opts.Client.Update(ctx, state.Pod); err != nil {
			return fmt.Errorf("unable to remove Pod %q finalizer %q: %w", state.Pod.Name, virtv2.FinalizerPodProtection, err)
		}
	}
	if state.Service != nil && controllerutil.RemoveFinalizer(state.Service, virtv2.FinalizerServiceProtection) {
		if err := opts.Client.Update(ctx, state.Service); err != nil {
			return fmt.Errorf("unable to remove Service %q finalizer %q: %w", state.Service.Name, virtv2.FinalizerServiceProtection, err)
		}
	}
	if state.PV != nil && controllerutil.RemoveFinalizer(state.PV, virtv2.FinalizerPVProtection) {
		if err := opts.Client.Update(ctx, state.PV); err != nil {
			return fmt.Errorf("unable to remove PV %q finalizer %q: %w", state.PV.Name, virtv2.FinalizerPVProtection, err)
		}
	}
	if state.PVC != nil && controllerutil.RemoveFinalizer(state.PVC, virtv2.FinalizerPVCProtection) {
		if err := opts.Client.Update(ctx, state.PVC); err != nil {
			return fmt.Errorf("unable to remove PVC %q finalizer %q: %w", state.PVC.Name, virtv2.FinalizerPVCProtection, err)
		}
	}
	if state.DV != nil && controllerutil.RemoveFinalizer(state.DV, virtv2.FinalizerDVProtection) {
		if err := opts.Client.Update(ctx, state.DV); err != nil {
			return fmt.Errorf("unable to remove DV %q finalizer %q: %w", state.DV.Name, virtv2.FinalizerDVProtection, err)
		}
	}
	return nil
}

// initPodName saves the Pod name in the annotation.
func (r *VMDReconciler) initPodName(state *VMDReconcilerState) {
	vmd := state.VMD.Changed()

	// Should not happen, but check it anyway.
	if vmd.Spec.DataSource == nil {
		return
	}

	switch vmd.Spec.DataSource.Type {
	case virtv2.DataSourceTypeUpload:
		uploaderPod := state.Supplements.UploaderPod()
		cc.AddAnnotation(vmd, cc.AnnUploadPodName, uploaderPod.Name)
	default:
		importerPod := state.Supplements.ImporterPod()
		cc.AddAnnotation(vmd, cc.AnnImportPodName, importerPod.Name)
	}
}

func (r *VMDReconciler) startPod(
	ctx context.Context,
	state *VMDReconcilerState,
	opts two_phase_reconciler.ReconcilerOptions,
) error {
	vmd := state.VMD.Current()

	// Should not happen, but check it anyway.
	if vmd.Spec.DataSource == nil {
		return nil
	}

	switch vmd.Spec.DataSource.Type {
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

// initDataVolumeName creates new DV name and update it in the annotation.
func (r *VMDReconciler) initDataVolumeName(state *VMDReconcilerState) {
	vmd := state.VMD.Changed()

	// Prevent DataVolume name regeneration.
	if _, hasKey := vmd.Annotations[cc.AnnVMDDataVolume]; hasKey {
		return
	}

	dv := state.Supplements.DataVolume()
	cc.AddAnnotation(vmd, cc.AnnVMDDataVolume, dv.Name)
}

func (r *VMDReconciler) cleanup(ctx context.Context, vmd *virtv2.VirtualMachineDisk, client client.Client, state *VMDReconcilerState) error {
	if state.DV != nil {
		err := supplements.CleanupForDataVolume(ctx, client, state.Supplements, r.dvcrSettings)
		if err != nil {
			return fmt.Errorf("cleanup supplements for DataVolume: %w", err)
		}
		// TODO(future): take ownership on PVC and delete DataVolume.
		// if err := client.Delete(ctx, state.DV); err != nil {
		//	return fmt.Errorf("cleanup DataVolume: %w", err)
		// }
	}

	if state.Pod != nil && cc.ShouldDeletePod(state.VMD.Current()) {
		if vmd.Spec.DataSource == nil {
			return fmt.Errorf("unexpected nil spec.dataSource to cleanup")
		}

		switch vmd.Spec.DataSource.Type {
		case virtv2.DataSourceTypeUpload:
			if err := uploader.CleanupService(ctx, client, state.Service); err != nil {
				return err
			}

			if err := uploader.CleanupPod(ctx, client, state.Pod); err != nil {
				return err
			}
		default:
			if err := importer.CleanupPod(ctx, client, state.Pod); err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *VMDReconciler) cleanupOnDeletion(ctx context.Context, state *VMDReconcilerState, opts two_phase_reconciler.ReconcilerOptions) error {
	if err := r.removeFinalizerChildResources(ctx, state, opts); err != nil {
		return err
	}
	if r.ShouldDeleteChildResources(state) {
		if err := r.cleanup(ctx, state.VMD.Current(), opts.Client, state); err != nil {
			return err
		}
		if err := helper.DeleteObject(ctx, opts.Client, state.DV); err != nil {
			return err
		}
		state.SetReconcilerResult(&reconcile.Result{RequeueAfter: 2 * time.Second})
		return nil
	}
	controllerutil.RemoveFinalizer(state.VMD.Changed(), virtv2.FinalizerVMDCleanup)
	return nil
}
