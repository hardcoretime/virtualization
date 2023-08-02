package controller

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	virtv2 "github.com/deckhouse/virtualization-controller/api/v2alpha1"
	cc "github.com/deckhouse/virtualization-controller/pkg/controller/common"
	"github.com/deckhouse/virtualization-controller/pkg/sdk/framework/helper"
	"github.com/deckhouse/virtualization-controller/pkg/util"
)

type CVMIReconcilerState struct {
	Client      client.Client
	CVMI        *helper.Resource[*virtv2.ClusterVirtualMachineImage, virtv2.ClusterVirtualMachineImageStatus]
	Pod         *corev1.Pod
	AttachedVMs []*virtv2.VirtualMachine
	Result      *reconcile.Result
}

func NewCVMIReconcilerState(name types.NamespacedName, log logr.Logger, client client.Client, cache cache.Cache) *CVMIReconcilerState {
	return &CVMIReconcilerState{
		Client: client,
		CVMI: helper.NewResource(
			name, log, client, cache,
			func() *virtv2.ClusterVirtualMachineImage { return &virtv2.ClusterVirtualMachineImage{} },
			func(obj *virtv2.ClusterVirtualMachineImage) virtv2.ClusterVirtualMachineImageStatus {
				return obj.Status
			},
		),
	}
}

func (state *CVMIReconcilerState) ApplySync(ctx context.Context, _ logr.Logger) error {
	if err := state.CVMI.UpdateMeta(ctx); err != nil {
		return fmt.Errorf("unable to update CVMI %q meta: %w", state.CVMI.Name(), err)
	}
	return nil
}

func (state *CVMIReconcilerState) ApplyUpdateStatus(ctx context.Context, _ logr.Logger) error {
	return state.CVMI.UpdateStatus(ctx)
}

func (state *CVMIReconcilerState) SetReconcilerResult(result *reconcile.Result) {
	state.Result = result
}

func (state *CVMIReconcilerState) GetReconcilerResult() *reconcile.Result {
	return state.Result
}

func (state *CVMIReconcilerState) ShouldApplyUpdateStatus() bool {
	return state.CVMI.IsStatusChanged()
}

func (state *CVMIReconcilerState) Reload(ctx context.Context, req reconcile.Request, log logr.Logger, client client.Client) error {
	if err := state.CVMI.Fetch(ctx); err != nil {
		return fmt.Errorf("unable to get %q: %w", req.NamespacedName, err)
	}
	if state.CVMI.IsEmpty() {
		log.Info("Reconcile observe an absent CVMI: it may be deleted", "CVMI", req.NamespacedName)
		return nil
	}

	pod, err := state.findImporterPod(ctx, client)
	if err != nil {
		return err
	}
	state.Pod = pod

	attachedVMs, err := state.findAttachedVMs(ctx)
	if err != nil {
		return err
	}
	state.AttachedVMs = attachedVMs

	log.Info("CVMI Reload", "AttachedVMs", attachedVMs)

	return nil
}

// ShouldReconcile tells if Sync and UpdateStatus should run.
// CVMI should not be reconciled if phase is Ready and no importer Pod found.
func (state *CVMIReconcilerState) ShouldReconcile(log logr.Logger) bool {
	// CVMI was not found. E.g. CVMI was deleted, but requeue task was triggered.
	if state.CVMI.IsEmpty() {
		return false
	}
	log.V(2).Info("CVMI ShouldReconcile", "ShouldRemoveProtectionFinalizer", state.ShouldRemoveProtectionFinalizer())
	if state.ShouldRemoveProtectionFinalizer() {
		return true
	}
	return !(cc.IsCVMIComplete(state.CVMI.Current()) && state.Pod == nil)
}

func (state *CVMIReconcilerState) findAttachedVMs(ctx context.Context) ([]*virtv2.VirtualMachine, error) {
	vml := &virtv2.VirtualMachineList{}

	req, err := labels.NewRequirement(
		MakeAttachedCVMILabelKey(state.CVMI.Name().Name),
		selection.Equals,
		[]string{"true"},
	)
	if err != nil {
		return nil, fmt.Errorf("unable to create label requirement: %w", err)
	}

	sel := labels.NewSelector()
	sel.Add(*req)

	if err := state.Client.List(ctx, vml, &client.ListOptions{LabelSelector: sel}); err != nil {
		return nil, fmt.Errorf("error getting VM by selector %v: %w", sel, err)
	}
	return util.ToPointersArray(vml.Items), nil
}

func (state *CVMIReconcilerState) findImporterPod(ctx context.Context, client client.Client) (*corev1.Pod, error) {
	// Extract namespace and name of the importer Pod from annotations.
	podName := state.CVMI.Current().Annotations[cc.AnnImportPodName]
	if podName == "" {
		return nil, nil
	}
	podNS := state.CVMI.Current().Annotations[cc.AnnImportPodNamespace]
	if podNS == "" {
		return nil, nil
	}

	objName := types.NamespacedName{Name: podName, Namespace: podNS}

	return helper.FetchObject(ctx, objName, client, &corev1.Pod{})
}

// HasImporterPodAnno returns whether CVMI is just created and has no annotations with importer Pod coordinates.
func (state *CVMIReconcilerState) HasImporterPodAnno() bool {
	if state.CVMI.IsEmpty() {
		return false
	}
	anno := state.CVMI.Current().GetAnnotations()
	if _, ok := anno[cc.AnnImportPodName]; !ok {
		return false
	}
	if _, ok := anno[cc.AnnImportPodNamespace]; !ok {
		return false
	}
	return true
}

// CanStartImporterPod returns whether CVMI is just created and has no annotations with importer Pod coordinates.
func (state *CVMIReconcilerState) CanStartImporterPod() bool {
	return !state.IsReady() && state.HasImporterPodAnno() && state.Pod == nil
}

// ShouldRemoveProtectionFinalizer returns whether CVMI protection finalizer should be removed
func (state *CVMIReconcilerState) ShouldRemoveProtectionFinalizer() bool {
	return controllerutil.ContainsFinalizer(state.CVMI.Current(), virtv2.FinalizerCVMIProtection) && (len(state.AttachedVMs) == 0)
}

func (state *CVMIReconcilerState) RemoveProtectionFinalizer() {
	controllerutil.RemoveFinalizer(state.CVMI.Changed(), virtv2.FinalizerCVMIProtection)
}

func (state *CVMIReconcilerState) HasImporterPod() bool {
	return state.HasImporterPodAnno() && state.Pod != nil
}

func (state *CVMIReconcilerState) IsReady() bool {
	if state.CVMI.IsEmpty() {
		return false
	}
	if !state.HasImporterPodAnno() {
		return false
	}
	return state.CVMI.Current().Status.Phase == virtv2.ImageReady
}

func (state *CVMIReconcilerState) IsDeletion() bool {
	if state.CVMI.IsEmpty() {
		return false
	}
	return state.CVMI.Current().DeletionTimestamp != nil
}

func (state *CVMIReconcilerState) Is() bool {
	return false
}

func (state *CVMIReconcilerState) IsImportComplete() bool {
	return state.Pod != nil && cc.IsPodComplete(state.Pod)
}

func (state *CVMIReconcilerState) IsImportInProgress() bool {
	return state.Pod != nil && state.Pod.Status.Phase == corev1.PodRunning
}