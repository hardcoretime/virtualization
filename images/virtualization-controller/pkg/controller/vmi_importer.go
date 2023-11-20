package controller

import (
	"context"
	"fmt"

	virtv2alpha1 "github.com/deckhouse/virtualization-controller/api/v2alpha1"
	vmiutil "github.com/deckhouse/virtualization-controller/pkg/common/vmi"
	cc "github.com/deckhouse/virtualization-controller/pkg/controller/common"
	"github.com/deckhouse/virtualization-controller/pkg/controller/importer"
	"github.com/deckhouse/virtualization-controller/pkg/dvcr"
	"github.com/deckhouse/virtualization-controller/pkg/sdk/framework/two_phase_reconciler"
)

func (r *VMIReconciler) startImporterPod(ctx context.Context, vmi *virtv2alpha1.VirtualMachineImage, opts two_phase_reconciler.ReconcilerOptions) error {
	opts.Log.V(1).Info("Creating importer POD for PVC", "pvc.Name", vmi.Name)

	importerSettings, err := r.createImporterSettings(vmi)
	if err != nil {
		return err
	}

	// all checks passed, let's create the importer pod!
	podSettings := r.createImporterPodSettings(vmi)

	caBundleSettings := importer.NewCABundleSettings(vmiutil.GetCABundle(vmi), vmi.Annotations[cc.AnnCABundleConfigMap])

	imp := importer.NewImporter(podSettings, importerSettings, caBundleSettings)
	pod, err := imp.CreatePod(ctx, opts.Client)
	if err != nil {
		err = cc.PublishPodErr(err, vmi.Annotations[cc.AnnImportPodName], vmi, opts.Recorder, opts.Client)
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
func (r *VMIReconciler) createImporterSettings(vmi *virtv2alpha1.VirtualMachineImage) (*importer.Settings, error) {
	settings := &importer.Settings{
		Verbose: r.verbose,
		Source:  cc.GetSource(vmi.Spec.DataSource),
	}

	switch settings.Source {
	case cc.SourceHTTP:
		if http := vmi.Spec.DataSource.HTTP; http != nil {
			importer.UpdateHTTPSettings(settings, http)
		}
	case cc.SourceRegistry:
		if secret := vmi.Spec.DataSource.ContainerImage.ImagePullSecret.Name; secret != "" {
			settings.AuthSecret = secret
		}
		if ctrImg := vmi.Spec.DataSource.ContainerImage; ctrImg != nil {
			importer.UpdateContainerImageSettings(settings, ctrImg)
		}
	case cc.SourceDVCR:
		switch vmi.Spec.DataSource.Type {
		case virtv2alpha1.DataSourceTypeClusterVirtualMachineImage:
			if cvmiImg := vmi.Spec.DataSource.ClusterVirtualMachineImage; cvmiImg != nil {
				importer.UpdateClusterVirtualMachineImageSettings(settings, cvmiImg, r.dvcrSettings.RegistryURL)
			}
		case virtv2alpha1.DataSourceTypeVirtualMachineImage:
			if vmiImg := vmi.Spec.DataSource.VirtualMachineImage; vmiImg != nil {
				vi := &virtv2alpha1.DataSourceVirtualMachineImage{
					Name:      vmiImg.Name,
					Namespace: vmi.Namespace,
				}
				importer.UpdateVirtualMachineImageSettings(settings, vi, r.dvcrSettings.RegistryURL)
			}
		default:
			return nil, fmt.Errorf("unknown dvcr settings source type: %s", vmi.Spec.DataSource.Type)
		}
	default:
		return nil, fmt.Errorf("unknown settings source: %s", settings.Source)
	}

	// Set DVCR settings.
	importer.UpdateDVCRSettings(settings, r.dvcrSettings, dvcr.RegistryImageName(r.dvcrSettings, dvcr.ImagePathForVMI(vmi)))

	// TODO Update proxy settings.

	return settings, nil
}

func (r *VMIReconciler) createImporterPodSettings(vmi *virtv2alpha1.VirtualMachineImage) *importer.PodSettings {
	return &importer.PodSettings{
		Name:            vmi.Annotations[cc.AnnImportPodName],
		Image:           r.importerImage,
		PullPolicy:      r.pullPolicy,
		Namespace:       vmi.GetNamespace(),
		OwnerReference:  vmiutil.MakeOwnerReference(vmi),
		ControllerName:  vmiControllerName,
		InstallerLabels: map[string]string{},
	}
}
