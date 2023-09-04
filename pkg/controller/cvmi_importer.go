package controller

import (
	"context"
	"fmt"

	virtv2alpha1 "github.com/deckhouse/virtualization-controller/api/v2alpha1"
	cvmiutil "github.com/deckhouse/virtualization-controller/pkg/common/cvmi"
	cc "github.com/deckhouse/virtualization-controller/pkg/controller/common"
	"github.com/deckhouse/virtualization-controller/pkg/controller/importer"
	"github.com/deckhouse/virtualization-controller/pkg/sdk/framework/two_phase_reconciler"
)

func (r *CVMIReconciler) startImporterPod(ctx context.Context, cvmi *virtv2alpha1.ClusterVirtualMachineImage, opts two_phase_reconciler.ReconcilerOptions) error {
	opts.Log.V(1).Info("Creating importer POD for PVC", "pvc.Name", cvmi.Name)

	importerSettings, err := r.createImporterSettings(cvmi)
	if err != nil {
		return err
	}

	// all checks passed, let's create the importer pod!
	podSettings := r.createImporterPodSettings(cvmi)

	caBundleSettings := importer.NewCABundleSettings(cvmiutil.GetCABundle(cvmi), cvmi.Annotations[cc.AnnCABundleConfigMap])

	imp := importer.NewImporter(podSettings, importerSettings, caBundleSettings)
	pod, err := imp.CreatePod(ctx, opts.Client)
	if err != nil {
		err = cc.PublishPodErr(err, cvmi.Annotations[cc.AnnImportPodName], cvmi, opts.Recorder, opts.Client)
		if err != nil {
			return err
		}
	}

	opts.Log.V(1).Info("Created importer POD", "pod.Name", pod.Name)

	if caBundleSettings != nil {
		if err := imp.EnsureCABundleConfigMap(ctx, opts.Client, pod); err != nil {
			return fmt.Errorf("create ConfigMap with certs from caBundle: %w", err)
		}
		opts.Log.V(1).Info("Created ConfigMap with caBundle", "cm.Name", caBundleSettings.ConfigMapName)
	}

	return nil
}

// createImporterSettings fills settings for the dvcr-importer binary.
func (r *CVMIReconciler) createImporterSettings(cvmi *virtv2alpha1.ClusterVirtualMachineImage) (*importer.Settings, error) {
	settings := &importer.Settings{
		Verbose: r.verbose,
		Source:  cc.GetSource(cvmi.Spec.DataSource),
	}

	switch settings.Source {
	case cc.SourceHTTP:
		if http := cvmi.Spec.DataSource.HTTP; http != nil {
			importer.UpdateHTTPSettings(settings, http)
		}
	case cc.SourceNone:
	default:
		return nil, fmt.Errorf("unknown settings source: %s", settings.Source)
	}

	// Set DVCR settings.
	importer.UpdateDVCRSettings(settings, r.dvcrSettings, cc.PrepareDVCREndpointFromCVMI(cvmi, r.dvcrSettings))

	// TODO Update proxy settings.

	return settings, nil
}

func (r *CVMIReconciler) createImporterPodSettings(cvmi *virtv2alpha1.ClusterVirtualMachineImage) *importer.PodSettings {
	return &importer.PodSettings{
		Name:            cvmi.Annotations[cc.AnnImportPodName],
		Image:           r.importerImage,
		PullPolicy:      r.pullPolicy,
		Namespace:       r.namespace,
		OwnerReference:  cvmiutil.MakeOwnerReference(cvmi),
		ControllerName:  cvmiControllerName,
		InstallerLabels: r.installerLabels,
	}
}
