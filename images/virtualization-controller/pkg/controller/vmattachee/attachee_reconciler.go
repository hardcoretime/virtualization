package vmattachee

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	virtv2 "github.com/deckhouse/virtualization-controller/api/v2alpha1"
	"github.com/deckhouse/virtualization-controller/pkg/controller/common"
	"github.com/deckhouse/virtualization-controller/pkg/sdk/framework/helper"
	"github.com/deckhouse/virtualization-controller/pkg/sdk/framework/two_phase_reconciler"
)

// AttacheeReconciler struct aimed to be included into the image or disk, which is attached to the VM
type AttacheeReconciler[T helper.Object[T, ST], ST any] struct {
	Kind         string
	IsNamespaced bool
}

func NewAttacheeReconciler[T helper.Object[T, ST], ST any](kind string, isNamespaced bool) *AttacheeReconciler[T, ST] {
	return &AttacheeReconciler[T, ST]{
		Kind:         kind,
		IsNamespaced: isNamespaced,
	}
}

func (r *AttacheeReconciler[T, ST]) SetupController(_ context.Context, mgr manager.Manager, ctr controller.Controller) error {
	matchAttacheeKindFunc := func(k, _ string) bool {
		_, found := ExtractAttachedResourceName(r.Kind, k)
		if found {
			return true
		}

		_, found = ExtractHotpluggedResourceName(r.Kind, k)

		return found
	}

	if err := ctr.Watch(
		source.Kind(mgr.GetCache(), &virtv2.VirtualMachine{}),
		handler.EnqueueRequestsFromMapFunc(r.enqueueAttacheeRequestsFromVM),
		predicate.Funcs{
			CreateFunc: func(e event.CreateEvent) bool {
				return common.HasLabel(e.Object.GetLabels(), matchAttacheeKindFunc)
			},
			DeleteFunc: func(e event.DeleteEvent) bool {
				return common.HasLabel(e.Object.GetLabels(), matchAttacheeKindFunc)
			},
			UpdateFunc: func(e event.UpdateEvent) bool {
				return common.HasLabel(e.ObjectOld.GetLabels(), matchAttacheeKindFunc) ||
					common.HasLabel(e.ObjectNew.GetLabels(), matchAttacheeKindFunc)
			},
		},
	); err != nil {
		return fmt.Errorf("error setting watch on VirtualMachineInstance: %w", err)
	}

	return nil
}

func (r *AttacheeReconciler[T, ST]) enqueueAttacheeRequestsFromVM(_ context.Context, obj client.Object) (res []reconcile.Request) {
	for k := range obj.GetLabels() {
		name, found := ExtractAttachedResourceName(r.Kind, k)
		if !found {
			name, found = ExtractHotpluggedResourceName(r.Kind, k)
		}

		if found {
			targetName := types.NamespacedName{Name: name}
			if r.IsNamespaced {
				if obj.GetNamespace() == "" {
					targetName.Namespace = "default"
				} else {
					targetName.Namespace = obj.GetNamespace()
				}
			}
			res = append(res, reconcile.Request{NamespacedName: targetName})
		}
	}
	return
}

func (r *AttacheeReconciler[T, ST]) Sync(_ context.Context, state *AttacheeState[T, ST], opts two_phase_reconciler.ReconcilerOptions) bool {
	opts.Log.V(2).Info("AttacheeReconciler Sync", "ShouldRemoveProtectionFinalizer", state.ShouldRemoveProtectionFinalizer())
	if state.ShouldRemoveProtectionFinalizer() {
		state.RemoveProtectionFinalizer()
		state.SetReconcilerResult(&reconcile.Result{Requeue: true})
		return true
	}
	return false
}