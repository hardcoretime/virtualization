package controller

import (
	"context"
	"errors"
	"fmt"

	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/types"
	cdiv1 "kubevirt.io/containerized-data-importer-api/pkg/apis/core/v1beta1"

	virtv2 "github.com/deckhouse/virtualization-controller/api/v2alpha1"
	vmiutil "github.com/deckhouse/virtualization-controller/pkg/common/vmi"
	"github.com/deckhouse/virtualization-controller/pkg/controller/kvbuilder"
	"github.com/deckhouse/virtualization-controller/pkg/controller/monitoring"
	"github.com/deckhouse/virtualization-controller/pkg/controller/supplements"
	"github.com/deckhouse/virtualization-controller/pkg/sdk/framework/two_phase_reconciler"
)

// getPVCSize retrieves PVC size from importer Pod final report after import is done.
func (r *VMIReconciler) getPVCSize(state *VMIReconcilerState) (resource.Quantity, error) {
	finalReport, err := monitoring.GetFinalReportFromPod(state.Pod)
	if err != nil {
		return resource.Quantity{}, fmt.Errorf("importer Pod final report missing: %w", err)
	}

	unpackedSize := *resource.NewQuantity(int64(finalReport.UnpackedSizeBytes), resource.BinarySI)
	if unpackedSize.IsZero() {
		return resource.Quantity{}, errors.New("no unpacked size in final report")
	}

	return unpackedSize, nil
}

func (r *VMIReconciler) createDataVolume(ctx context.Context, vmi *virtv2.VirtualMachineImage, state *VMIReconcilerState, opts two_phase_reconciler.ReconcilerOptions) error {
	// Retrieve PVC size.
	pvcSize, err := r.getPVCSize(state)
	if err != nil {
		return err
	}

	dvName := state.Supplements.DataVolume()

	dv, err := r.makeDataVolumeFromVMI(state, dvName, pvcSize)
	if err != nil {
		return err
	}

	if err = opts.Client.Create(ctx, dv); err != nil {
		opts.Log.V(2).Info("Error create new DV spec", "dv.spec", dv.Spec)
		return fmt.Errorf("create DataVolume/%s for VMI/%s: %w", dv.GetName(), vmi.GetName(), err)
	}
	opts.Log.Info("Created new DV", "dv.name", dv.GetName())
	opts.Log.V(2).Info("Created new DV spec", "dv.spec", dv.Spec)

	if vmiutil.IsTwoPhaseImport(vmi) {
		// Copy auth credentials and ca bundle to access DVCR as 'registry' data source.
		// Set DV as an ownerRef to auto-cleanup these copies.
		return supplements.EnsureForDataVolume(ctx, opts.Client, state.Supplements, dv, r.dvcrSettings)
	}

	return nil
}

// makeDataVolumeFromVMD makes DataVolume with 'registry' dataSource to import
// DVCR image onto PVC.
func (r *VMIReconciler) makeDataVolumeFromVMI(state *VMIReconcilerState, dvName types.NamespacedName, pvcSize resource.Quantity) (*cdiv1.DataVolume, error) {
	dvBuilder := kvbuilder.NewDV(dvName)
	vmi := state.VMI.Current()
	ds := vmi.Spec.DataSource

	authSecretName := r.dvcrSettings.AuthSecret
	if supplements.ShouldCopyDVCRAuthSecretForDataVolume(r.dvcrSettings, state.Supplements) {
		authSecret := state.Supplements.DVCRAuthSecretForDV()
		authSecretName = authSecret.Name
	}
	caBundleName := state.Supplements.DVCRCABundleConfigMapForDV().Name

	// Set datasource:
	// 'registry' if import is two phased.
	switch {
	case vmiutil.IsTwoPhaseImport(vmi):
		// The image was preloaded from source into dvcr.
		// We can't use the same data source a second time, but we can set dvcr as the data source.
		// Use DV name for the Secret with DVCR auth and the ConfigMap with DVCR CA Bundle.
		dvcrSourceImageName := r.dvcrSettings.RegistryImageForVMI(vmi.Name, vmi.Namespace)
		dvBuilder.SetRegistryDataSource(dvcrSourceImageName, authSecretName, caBundleName)
	case ds.Type == virtv2.DataSourceTypeClusterVirtualMachineImage:
		dvcrSourceImageName := r.dvcrSettings.RegistryImageForCVMI(ds.ClusterVirtualMachineImage.Name)
		dvBuilder.SetRegistryDataSource(dvcrSourceImageName, authSecretName, caBundleName)
	case ds.Type == virtv2.DataSourceTypeVirtualMachineImage:
		vmiRef := ds.VirtualMachineImage
		// NOTE: use namespace from current VMI.
		dvcrSourceImageName := r.dvcrSettings.RegistryImageForVMI(vmiRef.Name, vmi.Namespace)
		dvBuilder.SetRegistryDataSource(dvcrSourceImageName, authSecretName, caBundleName)
	default:
		return nil, fmt.Errorf("unsupported dataSource type %q", vmiutil.GetDataSourceType(vmi))
	}

	dvBuilder.SetPVC(vmi.Spec.PersistentVolumeClaim.StorageClassName, pvcSize)

	dvBuilder.SetOwnerRef(vmi, vmi.GetObjectKind().GroupVersionKind())
	dvBuilder.AddFinalizer(virtv2.FinalizerDVProtection)

	return dvBuilder.GetResource(), nil
}