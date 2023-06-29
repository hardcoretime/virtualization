package controller

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"

	virtv2alpha1 "github.com/deckhouse/virtualization-controller/api/v2alpha1"
	common "github.com/deckhouse/virtualization-controller/pkg/common"
	cc "github.com/deckhouse/virtualization-controller/pkg/controller/common"
)

const (
	// CertVolName is the name of the volume containing certs
	CertVolName = "cdi-cert-vol"

	// SecretVolName is the name of the volume containing gcs key
	SecretVolName = "cdi-secret-vol"

	// AnnOwnerRef is used when owner is in a different namespace
	AnnOwnerRef = cc.AnnAPIGroup + "/storage.ownerRef"

	// PodRunningReason is const that defines the pod was started as a reason
	PodRunningReason = "Pod is running"

	// ProxyCertVolName is the name of the volumecontaining certs
	ProxyCertVolName = "cdi-proxy-cert-vol"

	// secretExtraHeadersVolumeName is the format string that specifies where extra HTTP header secrets will be mounted
	secretExtraHeadersVolumeName = "import-extra-headers-vol-%d"

	// DestinationAuthVol is the name of the volume containing DVCR docker auth config.
	DestinationAuthVol = "dvcr-secret-vol"
)

type importerPodArgs struct {
	image       string
	importImage string
	verbose     string
	pullPolicy  string
	podEnvVar   *cc.ImportPodEnvVar
	namespace   string
	//pvc                     *corev1.PersistentVolumeClaim
	cvmi                    *virtv2alpha1.ClusterVirtualMachineImage
	vmi                     *virtv2alpha1.VirtualMachineImage
	scratchPvcName          *string
	podResourceRequirements *corev1.ResourceRequirements
	imagePullSecrets        []corev1.LocalObjectReference
	//workloadNodePlacement   *sdkapi.NodePlacement
	vddkImageName     *string
	priorityClassName string
}

//// returns the import image part of the endpoint string
//func getRegistryImportImage(cvmi *virtv2alpha1.ClusterVirtualMachineImage) (string, error) {
//	ep, err := cc.GetEndpoint(cvmi)
//	if err != nil {
//		return "", nil
//	}
//	if cc.IsImageStream(cvmi) {
//		return ep, nil
//	}
//	url, err := url.Parse(ep)
//	if err != nil {
//		return "", errors.Errorf("illegal registry endpoint %s", ep)
//	}
//	return url.Host + url.Path, nil
//}

// getValueFromAnnotation returns the value of an annotation
// cvmi *v1alpha1.ClusterVirtualMachineImage
func getValueFromAnnotation(obj metav1.Object, annotation string) string {
	return obj.GetAnnotations()[annotation]
}

// If this pod is going to transfer one checkpoint in a multi-stage import, attach the checkpoint name to the pod name so
// that each checkpoint gets a unique pod. That way each pod can be inspected using the retainAfterCompletion annotation.
func podNameWithCheckpoint(pvc *corev1.PersistentVolumeClaim) string {
	if checkpoint := pvc.Annotations[cc.AnnCurrentCheckpoint]; checkpoint != "" {
		return pvc.Name + "-checkpoint-" + checkpoint
	}
	return pvc.Name
}

//func getImportPodNameFromPvc(pvc *corev1.PersistentVolumeClaim) string {
//	podName, ok := pvc.Annotations[cc.AnnImportPod]
//	if ok {
//		return podName
//	}
//	// fallback to legacy naming, in fact the following function is fully compatible with legacy
//	// name concatenation "importer-{pvc.Name}" if the name length is under the size limits,
//	return naming.GetResourceName(common.ImporterPodName, podNameWithCheckpoint(pvc))
//}
//
//func createImportPodNameFromPvc(pvc *corev1.PersistentVolumeClaim) string {
//	return naming.GetResourceName(common.ImporterPodName, podNameWithCheckpoint(pvc))
//}

// createImporterPod creates and returns a pointer to a pod which is created based on the passed-in endpoint, secret
// name, and pvc. A nil secret means the endpoint credentials are not passed to the
// importer pod.
func createImporterPod(ctx context.Context, log logr.Logger, client client.Client, args *importerPodArgs, installerLabels map[string]string) (*corev1.Pod, error) {
	var err error
	//args.podResourceRequirements, err = cc.GetDefaultPodResourceRequirements(client)
	//if err != nil {
	//	return nil, err
	//}

	//args.imagePullSecrets, err = cc.GetImagePullSecrets(client)
	//if err != nil {
	//	return nil, err
	//}

	//args.workloadNodePlacement, err = cc.GetWorkloadNodePlacement(ctx, client)
	//if err != nil {
	//	return nil, err
	//}

	var pod *corev1.Pod
	//if cc.GetSource(args.pvc) == cc.SourceRegistry && args.pvc.Annotations[cc.AnnRegistryImportMethod] == string(cdiv1.RegistryPullNode) {
	//	args.importImage, err = getRegistryImportImage(args.pvc)
	//	if err != nil {
	//		return nil, err
	//	}
	//	pod = makeNodeImporterPodSpec(args)
	//} else {
	//	pod = makeImporterPodSpec(args)
	//}
	pod = makeImporterPodSpec(args)

	SetRecommendedLabels(pod, installerLabels, importControllerName)

	if err = client.Create(context.TODO(), pod); err != nil {
		return nil, err
	}

	log.V(3).Info("importer pod created\n", "pod.Name", pod.Name, "pod.Namespace", pod.Namespace, "image name", args.image)
	return pod, nil
}

// makeImporterPodSpec creates and return the importer pod spec based on the passed-in endpoint, secret and pvc.
func makeImporterPodSpec(args *importerPodArgs) *corev1.Pod {
	// importer pod name contains the pvc name
	podName := args.cvmi.Annotations[cc.AnnImportPod]

	blockOwnerDeletion := true
	isController := true

	volumes := []corev1.Volume{
		{
			// For test only
			Name: "emptydir",
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{
					Medium:    "",
					SizeLimit: nil,
				},
			},
		},
		{
			Name: "dvcr-auth",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: "ghcr-io",
				},
			},
		},
	}

	importerContainer := makeImporterContainerSpec(args.image, args.verbose, args.pullPolicy)

	pod := &corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      podName,
			Namespace: args.namespace,
			Annotations: map[string]string{
				cc.AnnCreatedBy: "yes",
			},
			//Labels: map[string]string{
			//	common.CDILabelKey:        common.CDILabelValue,
			//	common.CDIComponentLabel:  common.ImporterPodName,
			//	common.PrometheusLabelKey: common.PrometheusLabelValue,
			//},
			OwnerReferences: []metav1.OwnerReference{
				// Set CVMI as a controller for this Pod.
				{
					APIVersion:         "v2alpha2",
					Kind:               "ClusterVirtualMachineImage",
					Name:               args.cvmi.Name,
					UID:                args.cvmi.GetUID(),
					BlockOwnerDeletion: &blockOwnerDeletion,
					Controller:         &isController,
				},
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				*importerContainer,
			},
			RestartPolicy: corev1.RestartPolicyOnFailure,
			Volumes:       volumes,
			//NodeSelector:      args.workloadNodePlacement.NodeSelector,
			//Tolerations:       args.workloadNodePlacement.Tolerations,
			//Affinity:          args.workloadNodePlacement.Affinity,
			PriorityClassName: args.priorityClassName,
			ImagePullSecrets:  args.imagePullSecrets,
		},
	}

	setImporterPodCommons(pod, args.podEnvVar, args.cvmi, args.podResourceRequirements, args.imagePullSecrets)

	if args.podEnvVar.CertConfigMap != "" {
		vm := corev1.VolumeMount{
			Name:      CertVolName,
			MountPath: common.ImporterCertDir,
		}
		pod.Spec.Containers[0].VolumeMounts = append(pod.Spec.Containers[0].VolumeMounts, vm)
		pod.Spec.Volumes = append(pod.Spec.Volumes, createConfigMapVolume(CertVolName, args.podEnvVar.CertConfigMap))
	}

	// DVCR auth secret.
	if args.podEnvVar.DestinationAuthSecret != "" {
		vm := corev1.VolumeMount{
			Name:      DestinationAuthVol,
			MountPath: common.ImporterDestinationAuthConfigDir,
		}

		pod.Spec.Containers[0].VolumeMounts = append(pod.Spec.Containers[0].VolumeMounts, vm)
		pod.Spec.Volumes = append(pod.Spec.Volumes, createSecretVolume(DestinationAuthVol, args.podEnvVar.DestinationAuthSecret))
	}

	//if args.podEnvVar.certConfigMapProxy != "" {
	//	vm := corev1.VolumeMount{
	//		Name:      ProxyCertVolName,
	//		MountPath: common.ImporterProxyCertDir,
	//	}
	//	pod.Spec.Containers[0].VolumeMounts = append(pod.Spec.Containers[0].VolumeMounts, vm)
	//	pod.Spec.Volumes = append(pod.Spec.Volumes, createConfigMapVolume(ProxyCertVolName, GetImportProxyConfigMapName(args.cvmi.Name)))
	//}

	//if args.podEnvVar.source == cc.SourceGCS && args.podEnvVar.secretName != "" {
	//	vm := corev1.VolumeMount{
	//		Name:      SecretVolName,
	//		MountPath: common.ImporterGoogleCredentialDir,
	//	}
	//	pod.Spec.Containers[0].VolumeMounts = append(pod.Spec.Containers[0].VolumeMounts, vm)
	//	pod.Spec.Volumes = append(pod.Spec.Volumes, createSecretVolume(SecretVolName, args.podEnvVar.secretName))
	//}

	//for index, header := range args.podEnvVar.secretExtraHeaders {
	//	vm := corev1.VolumeMount{
	//		Name:      fmt.Sprintf(secretExtraHeadersVolumeName, index),
	//		MountPath: path.Join(common.ImporterSecretExtraHeadersDir, fmt.Sprint(index)),
	//	}
	//
	//	vol := corev1.Volume{
	//		Name: fmt.Sprintf(secretExtraHeadersVolumeName, index),
	//		VolumeSource: corev1.VolumeSource{
	//			Secret: &corev1.SecretVolumeSource{
	//				SecretName: header,
	//			},
	//		},
	//	}
	//
	//	pod.Spec.Containers[0].VolumeMounts = append(pod.Spec.Containers[0].VolumeMounts, vm)
	//	pod.Spec.Volumes = append(pod.Spec.Volumes, vol)
	//}

	cc.SetRestrictedSecurityContext(&pod.Spec)
	//// We explicitly define a NodeName for dynamically provisioned PVCs
	//// when the PVC is being handled by a populator (PVC')
	//cc.SetNodeNameIfPopulator(args.pvc, &pod.Spec)

	return pod
}

func setImporterPodCommons(pod *corev1.Pod, podEnvVar *cc.ImportPodEnvVar, cvmi *virtv2alpha1.ClusterVirtualMachineImage, podResourceRequirements *corev1.ResourceRequirements, imagePullSecrets []corev1.LocalObjectReference) {
	if podResourceRequirements != nil {
		for i := range pod.Spec.Containers {
			pod.Spec.Containers[i].Resources = *podResourceRequirements
		}
	}
	pod.Spec.ImagePullSecrets = imagePullSecrets

	ownerUID := cvmi.UID
	if len(cvmi.OwnerReferences) == 1 {
		ownerUID = cvmi.OwnerReferences[0].UID
	}

	pod.Spec.Containers[0].Env = makeImportEnv(podEnvVar, ownerUID)

	//setPodPvcAnnotations(pod, pvc)
}

func makeImporterContainerSpec(image, verbose, pullPolicy string) *corev1.Container {
	return &corev1.Container{
		Name:            common.ImporterPodName,
		Image:           image,
		ImagePullPolicy: corev1.PullPolicy(pullPolicy),
		Command:         []string{"sh"},
		Args:            []string{"/entrypoint.sh", "-v=" + verbose},
		Ports: []corev1.ContainerPort{
			{
				Name:          "metrics",
				ContainerPort: 8443,
				Protocol:      corev1.ProtocolTCP,
			},
		},
	}
}

func createConfigMapVolume(certVolName, objRef string) corev1.Volume {
	return corev1.Volume{
		Name: certVolName,
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: objRef,
				},
			},
		},
	}
}

func createSecretVolume(thisVolName, objRef string) corev1.Volume {
	return corev1.Volume{
		Name: thisVolName,
		VolumeSource: corev1.VolumeSource{
			Secret: &corev1.SecretVolumeSource{
				SecretName: objRef,
			},
		},
	}
}

// return the Env portion for the importer container.
func makeImportEnv(podEnvVar *cc.ImportPodEnvVar, uid types.UID) []corev1.EnvVar {
	env := []corev1.EnvVar{
		{
			Name:  common.ImporterSource,
			Value: podEnvVar.Source,
		},
		{
			Name:  common.ImporterEndpoint,
			Value: podEnvVar.Endpoint,
		},
		{
			Name:  common.ImporterContentType,
			Value: podEnvVar.ContentType,
		},
		{
			Name:  common.ImporterImageSize,
			Value: podEnvVar.ImageSize,
		},
		{
			Name:  common.OwnerUID,
			Value: string(uid),
		},
		{
			Name:  common.FilesystemOverheadVar,
			Value: podEnvVar.FilesystemOverhead,
		},
		{
			Name:  common.InsecureTLSVar,
			Value: strconv.FormatBool(podEnvVar.InsecureTLS),
		},
		{
			Name:  common.ImporterDiskID,
			Value: podEnvVar.DiskID,
		},
		{
			Name:  common.ImporterUUID,
			Value: podEnvVar.UUID,
		},
		{
			Name:  common.ImporterReadyFile,
			Value: podEnvVar.ReadyFile,
		},
		{
			Name:  common.ImporterDoneFile,
			Value: podEnvVar.DoneFile,
		},
		{
			Name:  common.ImporterBackingFile,
			Value: podEnvVar.BackingFile,
		},
		{
			Name:  common.ImporterThumbprint,
			Value: podEnvVar.Thumbprint,
		},
		{
			Name:  common.ImportProxyHTTP,
			Value: podEnvVar.HttpProxy,
		},
		{
			Name:  common.ImportProxyHTTPS,
			Value: podEnvVar.HttpsProxy,
		},
		{
			Name:  common.ImportProxyNoProxy,
			Value: podEnvVar.NoProxy,
		},
		{
			Name:  common.ImporterDestinationEndpoint,
			Value: podEnvVar.DestinationEndpoint,
		},
		{
			Name:  common.DestinationInsecureTLSVar,
			Value: podEnvVar.DestinationInsecureTLS,
		},
		// DVCR settings.
		{
			Name:  common.ImporterDestinationAuthConfigVar,
			Value: common.ImporterDestinationAuthConfigFile,
		},
	}

	//if podEnvVar.secretName != "" && podEnvVar.source != cc.SourceGCS {
	if podEnvVar.SecretName != "" {
		env = append(env, corev1.EnvVar{
			Name: common.ImporterAccessKeyID,
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: podEnvVar.SecretName,
					},
					Key: common.KeyAccess,
				},
			},
		}, corev1.EnvVar{
			Name: common.ImporterSecretKey,
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: podEnvVar.SecretName,
					},
					Key: common.KeySecret,
				},
			},
		})

	}
	if podEnvVar.CertConfigMap != "" {
		env = append(env, corev1.EnvVar{
			Name:  common.ImporterCertDirVar,
			Value: common.ImporterCertDir,
		})
	}
	if podEnvVar.CertConfigMapProxy != "" {
		env = append(env, corev1.EnvVar{
			Name:  common.ImporterProxyCertDirVar,
			Value: common.ImporterProxyCertDir,
		})
	}
	for index, header := range podEnvVar.ExtraHeaders {
		env = append(env, corev1.EnvVar{
			Name:  fmt.Sprintf("%s%d", common.ImporterExtraHeader, index),
			Value: header,
		})
	}
	return env
}

// setPodPvcAnnotations applies PVC annotations on the pod
func setPodPvcAnnotations(pod *corev1.Pod, pvc *corev1.PersistentVolumeClaim) {
	allowedAnnotations := map[string]string{
		cc.AnnPodNetwork:              "",
		cc.AnnPodSidecarInjection:     cc.AnnPodSidecarInjectionDefault,
		cc.AnnPodMultusDefaultNetwork: ""}
	for ann, def := range allowedAnnotations {
		val, ok := pvc.Annotations[ann]
		if !ok && def != "" {
			val = def
		}
		if val != "" {
			klog.V(1).Info("Applying PVC annotation on the pod", ann, val)
			if pod.Annotations == nil {
				pod.Annotations = map[string]string{}
			}
			pod.Annotations[ann] = val
		}
	}
}

// GetImportProxyConfigMapName return prefixed name.
// TODO add validation against name limitations
func GetImportProxyConfigMapName(suffix string) string {
	return fmt.Sprintf("import-proxy-cm-%s", suffix)
}
