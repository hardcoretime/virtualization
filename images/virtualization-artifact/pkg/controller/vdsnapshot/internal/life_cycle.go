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

package internal

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/deckhouse/virtualization-controller/pkg/controller/service"
	"github.com/deckhouse/virtualization-controller/pkg/logger"
	virtv2 "github.com/deckhouse/virtualization/api/core/v1alpha2"
	"github.com/deckhouse/virtualization/api/core/v1alpha2/vdscondition"
)

type LifeCycleHandler struct {
	snapshotter LifeCycleSnapshotter
}

func NewLifeCycleHandler(snapshotter LifeCycleSnapshotter) *LifeCycleHandler {
	return &LifeCycleHandler{
		snapshotter: snapshotter,
	}
}

func (h LifeCycleHandler) Handle(ctx context.Context, vdSnapshot *virtv2.VirtualDiskSnapshot) (reconcile.Result, error) {
	log := logger.FromContext(ctx).With(logger.SlogHandler("lifecycle"))

	condition, ok := service.GetCondition(vdscondition.VirtualDiskSnapshotReadyType, vdSnapshot.Status.Conditions)
	if !ok {
		condition = metav1.Condition{
			Type:   vdscondition.VirtualDiskSnapshotReadyType,
			Status: metav1.ConditionUnknown,
		}
	}

	defer func() { service.SetCondition(condition, &vdSnapshot.Status.Conditions) }()

	vs, err := h.snapshotter.GetVolumeSnapshot(ctx, vdSnapshot.Name, vdSnapshot.Namespace)
	if err != nil {
		setPhaseConditionToFailed(&condition, &vdSnapshot.Status.Phase, err)
		return reconcile.Result{}, err
	}

	vd, err := h.snapshotter.GetVirtualDisk(ctx, vdSnapshot.Spec.VirtualDiskName, vdSnapshot.Namespace)
	if err != nil {
		setPhaseConditionToFailed(&condition, &vdSnapshot.Status.Phase, err)
		return reconcile.Result{}, err
	}

	vm, err := getVirtualMachine(ctx, vd, h.snapshotter)
	if err != nil {
		setPhaseConditionToFailed(&condition, &vdSnapshot.Status.Phase, err)
		return reconcile.Result{}, err
	}

	if vdSnapshot.DeletionTimestamp != nil {
		vdSnapshot.Status.Phase = virtv2.VirtualDiskSnapshotPhaseTerminating
		condition.Status = metav1.ConditionUnknown
		condition.Reason = ""
		condition.Message = ""

		return reconcile.Result{}, nil
	}

	switch vdSnapshot.Status.Phase {
	case "":
		vdSnapshot.Status.Phase = virtv2.VirtualDiskSnapshotPhasePending
	case virtv2.VirtualDiskSnapshotPhaseReady:
		if vs == nil || vs.Status == nil || vs.Status.ReadyToUse == nil || !*vs.Status.ReadyToUse {
			vdSnapshot.Status.Phase = virtv2.VirtualDiskSnapshotPhaseFailed
			condition.Status = metav1.ConditionFalse
			condition.Reason = vdscondition.VolumeSnapshotLost
			condition.Message = fmt.Sprintf("The underlieng volume snapshot %q is not ready to use.", vdSnapshot.Status.VolumeSnapshotName)
			return reconcile.Result{Requeue: true}, nil
		}

		vdSnapshot.Status.Phase = virtv2.VirtualDiskSnapshotPhaseReady
		condition.Status = metav1.ConditionTrue
		condition.Reason = vdscondition.VirtualDiskSnapshotReady
		condition.Message = ""
		return reconcile.Result{}, nil
	}

	virtualDiskReadyCondition, _ := service.GetCondition(vdscondition.VirtualDiskReadyType, vdSnapshot.Status.Conditions)
	if vd == nil || virtualDiskReadyCondition.Status != metav1.ConditionTrue {
		vdSnapshot.Status.Phase = virtv2.VirtualDiskSnapshotPhasePending
		condition.Status = metav1.ConditionFalse
		condition.Reason = vdscondition.WaitingForTheVirtualDisk
		condition.Message = fmt.Sprintf("Waiting for the virtual disk %q to be ready for snapshotting.", vdSnapshot.Spec.VirtualDiskName)
		return reconcile.Result{}, nil
	}

	var pvc *corev1.PersistentVolumeClaim
	if vd.Status.Target.PersistentVolumeClaim != "" {
		pvc, err = h.snapshotter.GetPersistentVolumeClaim(ctx, vd.Status.Target.PersistentVolumeClaim, vd.Namespace)
		if err != nil {
			setPhaseConditionToFailed(&condition, &vdSnapshot.Status.Phase, err)
			return reconcile.Result{}, err
		}
	}

	if pvc == nil || pvc.Status.Phase != corev1.ClaimBound {
		vdSnapshot.Status.Phase = virtv2.VirtualDiskSnapshotPhasePending
		condition.Status = metav1.ConditionFalse
		condition.Reason = vdscondition.WaitingForTheVirtualDisk
		condition.Message = "Waiting for the virtual disk's pvc to be in phase Bound."
		return reconcile.Result{}, nil
	}

	switch {
	case vs == nil:
		if vm != nil && vm.Status.Phase != virtv2.MachineStopped && !h.snapshotter.IsFrozen(vm) {
			if h.snapshotter.CanFreeze(vm) {
				log.Debug("Freeze the virtual machine to take a snapshot")

				err = h.snapshotter.Freeze(ctx, vm.Name, vm.Namespace)
				if err != nil {
					setPhaseConditionToFailed(&condition, &vdSnapshot.Status.Phase, err)
					return reconcile.Result{}, err
				}

				vdSnapshot.Status.Phase = virtv2.VirtualDiskSnapshotPhaseInProgress
				condition.Status = metav1.ConditionFalse
				condition.Reason = vdscondition.FileSystemFreezing
				condition.Message = fmt.Sprintf(
					"The virtual machine %q with an attached virtual disk %q is in the process of being frozen for taking a snapshot.",
					vm.Name, vdSnapshot.Spec.VirtualDiskName,
				)
				return reconcile.Result{}, nil
			}

			if vdSnapshot.Spec.RequiredConsistency {
				vdSnapshot.Status.Phase = virtv2.VirtualDiskSnapshotPhasePending
				condition.Status = metav1.ConditionFalse
				condition.Reason = vdscondition.PotentiallyInconsistent
				condition.Message = fmt.Sprintf(
					"The virtual machine %q with an attached virtual disk %q is %s: "+
						"the snapshotting of virtual disk might result in an inconsistent snapshot: "+
						"waiting for the virtual machine to be %s or the disk to be detached",
					vm.Name, vd.Name, vm.Status.Phase, virtv2.MachineStopped,
				)
				return reconcile.Result{}, nil
			}
		}

		log.Debug("The corresponding volume snapshot not found: create the new one")

		vs, err = h.snapshotter.CreateVolumeSnapshot(ctx, vdSnapshot, pvc)
		if err != nil {
			setPhaseConditionToFailed(&condition, &vdSnapshot.Status.Phase, err)
			return reconcile.Result{}, err
		}

		vdSnapshot.Status.Phase = virtv2.VirtualDiskSnapshotPhaseInProgress
		vdSnapshot.Status.VolumeSnapshotName = vs.Name
		condition.Status = metav1.ConditionFalse
		condition.Reason = vdscondition.Snapshotting
		condition.Message = fmt.Sprintf("The snapshotting process for virtual disk %q has started.", vdSnapshot.Spec.VirtualDiskName)
		return reconcile.Result{}, nil
	case vs.Status != nil && vs.Status.Error != nil && vs.Status.Error.Message != nil:
		log.Debug("The volume snapshot has an error")

		vdSnapshot.Status.Phase = virtv2.VirtualDiskSnapshotPhaseFailed
		condition.Status = metav1.ConditionFalse
		condition.Reason = vdscondition.VirtualDiskSnapshotFailed
		condition.Message = fmt.Sprintf("VolumeSnapshot %q has an error: %s.", vs.Name, *vs.Status.Error.Message)
		return reconcile.Result{}, nil
	case vs.Status == nil || vs.Status.ReadyToUse == nil || !*vs.Status.ReadyToUse:
		log.Debug("Waiting for the volume snapshot to be ready to use")

		vdSnapshot.Status.Phase = virtv2.VirtualDiskSnapshotPhaseInProgress
		condition.Status = metav1.ConditionFalse
		condition.Reason = vdscondition.Snapshotting
		condition.Message = fmt.Sprintf("Waiting fot the volume snapshot %q to be ready to use.", vdSnapshot.Name)
		return reconcile.Result{}, nil
	default:
		log.Debug("The volume snapshot is ready to use")

		switch {
		case vm == nil, vm.Status.Phase == virtv2.MachineStopped:
			vdSnapshot.Status.Consistent = ptr.To(true)
		case h.snapshotter.IsFrozen(vm):
			vdSnapshot.Status.Consistent = ptr.To(true)

			var canUnfreeze bool
			canUnfreeze, err = h.snapshotter.CanUnfreeze(ctx, vdSnapshot.Name, vm)
			if err != nil {
				setPhaseConditionToFailed(&condition, &vdSnapshot.Status.Phase, err)
				return reconcile.Result{}, err
			}

			if canUnfreeze {
				log.Debug("Unfreeze the virtual machine after taking a snapshot")

				err = h.snapshotter.Unfreeze(ctx, vm.Name, vm.Namespace)
				if err != nil {
					setPhaseConditionToFailed(&condition, &vdSnapshot.Status.Phase, err)
					return reconcile.Result{}, err
				}
			}
		}

		vdSnapshot.Status.Phase = virtv2.VirtualDiskSnapshotPhaseReady
		condition.Status = metav1.ConditionTrue
		condition.Reason = vdscondition.VirtualDiskSnapshotReady
		condition.Message = ""

		return reconcile.Result{}, nil
	}
}

func getVirtualMachine(ctx context.Context, vd *virtv2.VirtualDisk, snapshotter LifeCycleSnapshotter) (*virtv2.VirtualMachine, error) {
	if vd == nil {
		return nil, nil
	}

	// TODO: ensure vd.Status.AttachedToVirtualMachines is in the actual state.
	switch len(vd.Status.AttachedToVirtualMachines) {
	case 0:
		return nil, nil
	case 1:
		vm, err := snapshotter.GetVirtualMachine(ctx, vd.Status.AttachedToVirtualMachines[0].Name, vd.Namespace)
		if err != nil {
			return nil, err
		}

		return vm, nil
	default:
		return nil, fmt.Errorf("the virtual disk %q is attached to multiple virtual machines", vd.Name)
	}
}

func setPhaseConditionToFailed(cond *metav1.Condition, phase *virtv2.VirtualDiskSnapshotPhase, err error) {
	*phase = virtv2.VirtualDiskSnapshotPhaseFailed
	cond.Status = metav1.ConditionFalse
	cond.Reason = vdscondition.VirtualDiskSnapshotFailed
	cond.Message = service.CapitalizeFirstLetter(err.Error())
}