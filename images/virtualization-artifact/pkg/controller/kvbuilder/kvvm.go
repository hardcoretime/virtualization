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

package kvbuilder

import (
	"fmt"
	"maps"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	virtv1 "kubevirt.io/api/core/v1"

	"github.com/deckhouse/virtualization-controller/pkg/controller/common"
	"github.com/deckhouse/virtualization-controller/pkg/sdk/framework/helper"
	"github.com/deckhouse/virtualization-controller/pkg/util"
	virtv2 "github.com/deckhouse/virtualization/api/core/v1alpha2"
)

// TODO(VM): Implement at this level some mechanics supporting "effectiveSpec" logic
// TODO(VM): KVVM builder should know which fields are allowed to be changed on-fly, and what params need a new KVVM instance.
// TODO(VM): Somehow report from this layer that "restart is needed" and controller will do other "effectiveSpec"-related stuff.

const (
	CloudInitDiskName = "cloudinit"
	SysprepDiskName   = "sysprep"
)

type KVVMOptions struct {
	EnableParavirtualization bool
	OsType                   virtv2.OsType

	// These options are for local development mode
	DisableHypervSyNIC bool
}

type KVVM struct {
	helper.ResourceBuilder[*virtv1.VirtualMachine]
	opts KVVMOptions
}

func NewKVVM(currentKVVM *virtv1.VirtualMachine, opts KVVMOptions) *KVVM {
	return &KVVM{
		ResourceBuilder: helper.NewResourceBuilder(currentKVVM, helper.ResourceBuilderOptions{ResourceExists: true}),
		opts:            opts,
	}
}

func NewEmptyKVVM(name types.NamespacedName, opts KVVMOptions) *KVVM {
	return &KVVM{
		opts: opts,
		ResourceBuilder: helper.NewResourceBuilder(
			&virtv1.VirtualMachine{
				TypeMeta: metav1.TypeMeta{
					Kind:       virtv1.VirtualMachineGroupVersionKind.Kind,
					APIVersion: virtv1.VirtualMachineGroupVersionKind.GroupVersion().String(),
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      name.Name,
					Namespace: name.Namespace,
				},
				Spec: virtv1.VirtualMachineSpec{
					Template: &virtv1.VirtualMachineInstanceTemplateSpec{},
				},
			}, helper.ResourceBuilderOptions{},
		),
	}
}

func (b *KVVM) SetKVVMIAnnotation(annoKey, annoValue string) {
	anno := b.Resource.Spec.Template.ObjectMeta.GetAnnotations()
	if anno == nil {
		anno = make(map[string]string)
	}

	anno[annoKey] = annoValue

	b.Resource.Spec.Template.ObjectMeta.SetAnnotations(anno)
}

func (b *KVVM) SetCPUModel(class *virtv2.VirtualMachineClass) error {
	var cpu virtv1.CPU

	switch class.Spec.CPU.Type {
	case virtv2.CPUTypeHost:
		cpu.Model = virtv1.CPUModeHostModel
	case virtv2.CPUTypeHostPassthrough:
		cpu.Model = virtv1.CPUModeHostPassthrough
	case virtv2.CPUTypeModel:
		cpu.Model = class.Spec.CPU.Model
	case virtv2.CPUTypeFeatures, virtv2.CPUTypeDiscovery:
		cpu.Features = make([]virtv1.CPUFeature, len(class.Status.CpuFeatures.Enabled))
		for i, feature := range class.Status.CpuFeatures.Enabled {
			cpu.Features[i] = virtv1.CPUFeature{
				Name:   feature,
				Policy: "require",
			}
		}
	default:
		return fmt.Errorf("unexpected cpu type: %q", class.Spec.CPU.Type)
	}

	b.Resource.Spec.Template.Spec.Domain.CPU = &cpu

	return nil
}

func (b *KVVM) SetRunPolicy(runPolicy virtv2.RunPolicy) error {
	switch runPolicy {
	case virtv2.AlwaysOnPolicy:
		b.Resource.Spec.RunStrategy = util.GetPointer(virtv1.RunStrategyAlways)
	case virtv2.AlwaysOffPolicy:
		b.Resource.Spec.RunStrategy = util.GetPointer(virtv1.RunStrategyHalted)
	case virtv2.ManualPolicy:
		if !b.ResourceExists {
			// initialize only
			b.Resource.Spec.RunStrategy = util.GetPointer(virtv1.RunStrategyManual)
		}
	case virtv2.AlwaysOnUnlessStoppedManually:
		if !b.ResourceExists {
			// initialize only
			b.Resource.Spec.RunStrategy = util.GetPointer(virtv1.RunStrategyAlways)
		}
	default:
		return fmt.Errorf("unexpected runPolicy %s. %w", runPolicy, common.ErrUnknownValue)
	}
	return nil
}

func (b *KVVM) SetNodeSelector(vmNodeSelector, classNodeSelector map[string]string) {
	if len(vmNodeSelector) == 0 && len(classNodeSelector) == 0 {
		return
	}
	selector := make(map[string]string, len(vmNodeSelector)+len(classNodeSelector))
	maps.Copy(selector, vmNodeSelector)
	maps.Copy(selector, classNodeSelector)
	b.Resource.Spec.Template.Spec.NodeSelector = selector
}

func (b *KVVM) SetTolerations(vmTolerations, classTolerations []corev1.Toleration) {
	tolerationsMap := make(map[string]corev1.Toleration)
	for _, toleration := range classTolerations {
		tolerationsMap[toleration.Key] = toleration
	}
	for _, toleration := range vmTolerations {
		tolerationsMap[toleration.Key] = toleration
	}
	resultTolerations := make([]corev1.Toleration, 0, len(tolerationsMap))
	for _, toleration := range tolerationsMap {
		resultTolerations = append(resultTolerations, toleration)
	}
	b.Resource.Spec.Template.Spec.Tolerations = resultTolerations
}

func (b *KVVM) SetPriorityClassName(priorityClassName string) {
	b.Resource.Spec.Template.Spec.PriorityClassName = priorityClassName
}

func (b *KVVM) SetAffinity(vmAffinity *corev1.Affinity, classMatchExpressions []corev1.NodeSelectorRequirement) {
	if len(classMatchExpressions) == 0 {
		b.Resource.Spec.Template.Spec.Affinity = vmAffinity
		return
	}
	if vmAffinity == nil {
		vmAffinity = &corev1.Affinity{}
	}
	if vmAffinity.NodeAffinity == nil {
		vmAffinity.NodeAffinity = &corev1.NodeAffinity{}
	}
	if vmAffinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution == nil {
		vmAffinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution = &corev1.NodeSelector{}
	}
	if vmAffinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms == nil {
		vmAffinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms = []corev1.NodeSelectorTerm{}
	}

	vmAffinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms = append(
		vmAffinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms,
		corev1.NodeSelectorTerm{MatchExpressions: classMatchExpressions},
	)

	b.Resource.Spec.Template.Spec.Affinity = vmAffinity
}

func (b *KVVM) SetTerminationGracePeriod(period *int64) {
	b.Resource.Spec.Template.Spec.TerminationGracePeriodSeconds = period
}

func (b *KVVM) SetTopologySpreadConstraint(topology []corev1.TopologySpreadConstraint) {
	b.Resource.Spec.Template.Spec.TopologySpreadConstraints = topology
}

func (b *KVVM) SetResourceRequirements(cores int, coreFraction string, memorySize resource.Quantity) error {
	cpuRequest, err := GetCPURequest(cores, coreFraction)
	if err != nil {
		return err
	}
	b.Resource.Spec.Template.Spec.Domain.Resources = virtv1.ResourceRequirements{
		Requests: corev1.ResourceList{
			corev1.ResourceCPU:    *cpuRequest,
			corev1.ResourceMemory: memorySize,
		},
		Limits: corev1.ResourceList{
			corev1.ResourceCPU:    *GetCPULimit(cores),
			corev1.ResourceMemory: memorySize,
		},
	}
	return nil
}

func GetCPURequest(cores int, coreFraction string) (*resource.Quantity, error) {
	if coreFraction == "" {
		return GetCPULimit(cores), nil
	}
	fraction := intstr.FromString(coreFraction)
	req, err := intstr.GetScaledValueFromIntOrPercent(&fraction, cores*1000, true)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate coreFraction. %w", err)
	}
	if req == 0 {
		return GetCPULimit(cores), nil
	}
	return resource.NewMilliQuantity(int64(req), resource.DecimalSI), nil
}

func GetCPULimit(cores int) *resource.Quantity {
	return resource.NewQuantity(int64(cores), resource.DecimalSI)
}

type SetDiskOptions struct {
	Provisioning *virtv2.Provisioning

	ContainerDisk         *string
	PersistentVolumeClaim *string

	IsHotplugged bool
	IsCdrom      bool
	IsEphemeral  bool

	Serial string

	BootOrder uint
}

func (b *KVVM) ClearDisks() {
	b.Resource.Spec.Template.Spec.Domain.Devices.Disks = nil
	b.Resource.Spec.Template.Spec.Volumes = nil
}

func (b *KVVM) SetDisk(name string, opts SetDiskOptions) error {
	devPreset := DeviceOptionsPresets.Find(b.opts.EnableParavirtualization)

	var dd virtv1.DiskDevice
	if opts.IsCdrom {
		dd.CDRom = &virtv1.CDRomTarget{
			Bus: devPreset.CdromBus,
		}
	} else {
		dd.Disk = &virtv1.DiskTarget{
			Bus: devPreset.DiskBus,
		}
	}

	disk := virtv1.Disk{
		Name:       name,
		DiskDevice: dd,
		Serial:     opts.Serial,
	}

	if opts.BootOrder > 0 {
		disk.BootOrder = &opts.BootOrder
	}

	b.Resource.Spec.Template.Spec.Domain.Devices.Disks = util.SetArrayElem(
		b.Resource.Spec.Template.Spec.Domain.Devices.Disks, disk,
		func(v1, v2 virtv1.Disk) bool {
			return v1.Name == v2.Name
		}, true,
	)

	var vs virtv1.VolumeSource
	switch {
	case opts.PersistentVolumeClaim != nil && !opts.IsEphemeral:
		vs.PersistentVolumeClaim = &virtv1.PersistentVolumeClaimVolumeSource{
			PersistentVolumeClaimVolumeSource: corev1.PersistentVolumeClaimVolumeSource{
				ClaimName: *opts.PersistentVolumeClaim,
			},
			Hotpluggable: opts.IsHotplugged,
		}

	case opts.PersistentVolumeClaim != nil && opts.IsEphemeral:
		vs.Ephemeral = &virtv1.EphemeralVolumeSource{
			PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
				ClaimName: *opts.PersistentVolumeClaim,
			},
		}

	case opts.ContainerDisk != nil:
		vs.ContainerDisk = &virtv1.ContainerDiskSource{
			Image: *opts.ContainerDisk,
		}

	case opts.Provisioning != nil:
		switch opts.Provisioning.Type {
		case virtv2.ProvisioningTypeSysprepRef:
			if opts.Provisioning.SysprepRef == nil {
				return fmt.Errorf("nil sysprep ref: %s", opts.Provisioning.Type)
			}

			switch opts.Provisioning.SysprepRef.Kind {
			case virtv2.SysprepRefKindSecret:
				vs.Sysprep = &virtv1.SysprepSource{
					Secret: &corev1.LocalObjectReference{
						Name: opts.Provisioning.SysprepRef.Name,
					},
				}
			default:
				return fmt.Errorf("unexpected sysprep ref kind: %s", opts.Provisioning.SysprepRef.Kind)
			}
		case virtv2.ProvisioningTypeUserData:
			vs.CloudInitNoCloud = &virtv1.CloudInitNoCloudSource{
				UserData: opts.Provisioning.UserData,
			}
		case virtv2.ProvisioningTypeUserDataRef:
			if opts.Provisioning.UserDataRef == nil {
				return fmt.Errorf("nil user data ref: %s", opts.Provisioning.Type)
			}

			switch opts.Provisioning.UserDataRef.Kind {
			case virtv2.UserDataRefKindSecret:
				vs.CloudInitNoCloud = &virtv1.CloudInitNoCloudSource{
					UserDataSecretRef: &corev1.LocalObjectReference{
						Name: opts.Provisioning.UserDataRef.Name,
					},
				}
			default:
				return fmt.Errorf("unexpected user data ref kind: %s", opts.Provisioning.UserDataRef.Kind)
			}
		default:
			return fmt.Errorf("unexpected provisioning type %s. %w", opts.Provisioning.Type, common.ErrUnknownType)
		}

	default:
		return fmt.Errorf("expected either opts.PersistentVolumeClaim or opts.ContainerDisk to be set, please report a bug")
	}

	volume := virtv1.Volume{
		Name:         name,
		VolumeSource: vs,
	}
	b.Resource.Spec.Template.Spec.Volumes = util.SetArrayElem(
		b.Resource.Spec.Template.Spec.Volumes, volume,
		func(v1, v2 virtv1.Volume) bool {
			return v1.Name == v2.Name
		}, true,
	)
	return nil
}

func (b *KVVM) SetTablet(name string) {
	i := virtv1.Input{
		Name: name,
		Bus:  virtv1.InputBusUSB,
		Type: virtv1.InputTypeTablet,
	}

	b.Resource.Spec.Template.Spec.Domain.Devices.Inputs = util.SetArrayElem(
		b.Resource.Spec.Template.Spec.Domain.Devices.Inputs, i,
		func(v1, v2 virtv1.Input) bool {
			return v1.Name == v2.Name
		}, true,
	)
}

// HasTablet checks tablet presence by its name.
func (b *KVVM) HasTablet(name string) bool {
	for _, input := range b.Resource.Spec.Template.Spec.Domain.Devices.Inputs {
		if input.Name == name && input.Type == virtv1.InputTypeTablet {
			return true
		}
	}
	return false
}

func (b *KVVM) SetProvisioning(p *virtv2.Provisioning) error {
	if p == nil {
		return nil
	}

	switch p.Type {
	case virtv2.ProvisioningTypeSysprepRef:
		return b.SetDisk(SysprepDiskName, SetDiskOptions{Provisioning: p, IsCdrom: true})
	case virtv2.ProvisioningTypeUserData, virtv2.ProvisioningTypeUserDataRef:
		return b.SetDisk(CloudInitDiskName, SetDiskOptions{Provisioning: p})
	default:
		return fmt.Errorf("unexpected provisioning type %s. %w", p.Type, common.ErrUnknownType)
	}
}

func (b *KVVM) SetOsType(osType virtv2.OsType) error {
	switch osType {
	case virtv2.Windows:
		b.Resource.Spec.Template.Spec.Domain.Machine = &virtv1.Machine{
			Type: "q35",
		}
		b.Resource.Spec.Template.Spec.Domain.Devices.AutoattachInputDevice = util.GetPointer(true)
		b.Resource.Spec.Template.Spec.Domain.Devices.TPM = &virtv1.TPMDevice{}
		b.Resource.Spec.Template.Spec.Domain.Features = &virtv1.Features{
			ACPI: virtv1.FeatureState{Enabled: util.GetPointer(true)},
			APIC: &virtv1.FeatureAPIC{Enabled: util.GetPointer(true)},
			SMM:  &virtv1.FeatureState{Enabled: util.GetPointer(true)},
			Hyperv: &virtv1.FeatureHyperv{
				Frequencies:     &virtv1.FeatureState{Enabled: util.GetPointer(true)},
				IPI:             &virtv1.FeatureState{Enabled: util.GetPointer(true)},
				Reenlightenment: &virtv1.FeatureState{Enabled: util.GetPointer(true)},
				Relaxed:         &virtv1.FeatureState{Enabled: util.GetPointer(true)},
				Reset:           &virtv1.FeatureState{Enabled: util.GetPointer(true)},
				Runtime:         &virtv1.FeatureState{Enabled: util.GetPointer(true)},
				Spinlocks: &virtv1.FeatureSpinlocks{
					Enabled: util.GetPointer(true),
					Retries: util.GetPointer[uint32](8191),
				},
				TLBFlush: &virtv1.FeatureState{Enabled: util.GetPointer(true)},
				VAPIC:    &virtv1.FeatureState{Enabled: util.GetPointer(true)},
				VPIndex:  &virtv1.FeatureState{Enabled: util.GetPointer(true)},
			},
		}

		if !b.opts.DisableHypervSyNIC {
			b.Resource.Spec.Template.Spec.Domain.Features.Hyperv.SyNIC = &virtv1.FeatureState{Enabled: util.GetPointer(true)}
			b.Resource.Spec.Template.Spec.Domain.Features.Hyperv.SyNICTimer = &virtv1.SyNICTimer{
				Enabled: util.GetPointer(true),
				Direct:  &virtv1.FeatureState{Enabled: util.GetPointer(true)},
			}
		}

	case virtv2.GenericOs:
		b.Resource.Spec.Template.Spec.Domain.Machine = &virtv1.Machine{
			Type: "q35",
		}
		b.Resource.Spec.Template.Spec.Domain.Devices.AutoattachInputDevice = util.GetPointer(true)
		b.Resource.Spec.Template.Spec.Domain.Devices.Rng = &virtv1.Rng{}
		b.Resource.Spec.Template.Spec.Domain.Features = &virtv1.Features{
			ACPI: virtv1.FeatureState{Enabled: util.GetPointer(true)},
			SMM:  &virtv1.FeatureState{Enabled: util.GetPointer(true)},
		}
	default:
		return fmt.Errorf("unexpected os type %q. %w", osType, common.ErrUnknownType)
	}
	return nil
}

// GetOSSettings returns a portion of devices and features related to d8 VM osType.
func (b *KVVM) GetOSSettings() map[string]interface{} {
	return map[string]interface{}{
		"machine": b.Resource.Spec.Template.Spec.Domain.Machine,
		"devices": map[string]interface{}{
			"autoattach": b.Resource.Spec.Template.Spec.Domain.Devices.AutoattachInputDevice,
			"tpm":        b.Resource.Spec.Template.Spec.Domain.Devices.TPM,
			"rng":        b.Resource.Spec.Template.Spec.Domain.Devices.Rng,
		},
		"features": map[string]interface{}{
			"acpi":   b.Resource.Spec.Template.Spec.Domain.Features.ACPI,
			"apic":   b.Resource.Spec.Template.Spec.Domain.Features.APIC,
			"smm":    b.Resource.Spec.Template.Spec.Domain.Features.SMM,
			"hyperv": b.Resource.Spec.Template.Spec.Domain.Features.Hyperv,
		},
	}
}

func (b *KVVM) SetNetworkInterface(name string) {
	devPreset := DeviceOptionsPresets.Find(b.opts.EnableParavirtualization)

	net := virtv1.Network{
		Name: name,
		NetworkSource: virtv1.NetworkSource{
			Pod: &virtv1.PodNetwork{},
		},
	}
	b.Resource.Spec.Template.Spec.Networks = util.SetArrayElem(
		b.Resource.Spec.Template.Spec.Networks, net,
		func(v1, v2 virtv1.Network) bool {
			return v1.Name == v2.Name
		}, true,
	)

	iface := virtv1.Interface{
		Name:  name,
		Model: devPreset.InterfaceModel,
	}
	iface.InterfaceBindingMethod.Bridge = &virtv1.InterfaceBridge{}
	b.Resource.Spec.Template.Spec.Domain.Devices.Interfaces = util.SetArrayElem(
		b.Resource.Spec.Template.Spec.Domain.Devices.Interfaces, iface,
		func(v1, v2 virtv1.Interface) bool {
			return v1.Name == v2.Name
		}, true,
	)
}

func (b *KVVM) SetBootloader(bootloader virtv2.BootloaderType) error {
	if b.Resource.Spec.Template.Spec.Domain.Firmware == nil {
		b.Resource.Spec.Template.Spec.Domain.Firmware = &virtv1.Firmware{}
	}

	switch bootloader {
	case "", virtv2.BIOS:
		b.Resource.Spec.Template.Spec.Domain.Firmware.Bootloader = nil
	case virtv2.EFI:
		b.Resource.Spec.Template.Spec.Domain.Firmware.Bootloader = &virtv1.Bootloader{
			EFI: &virtv1.EFI{
				SecureBoot: util.GetPointer(false),
			},
		}
	case virtv2.EFIWithSecureBoot:
		if b.Resource.Spec.Template.Spec.Domain.Features == nil {
			b.Resource.Spec.Template.Spec.Domain.Features = &virtv1.Features{}
		}
		b.Resource.Spec.Template.Spec.Domain.Features.SMM = &virtv1.FeatureState{
			Enabled: util.GetPointer(true),
		}
		b.Resource.Spec.Template.Spec.Domain.Firmware.Bootloader = &virtv1.Bootloader{
			EFI: &virtv1.EFI{SecureBoot: util.GetPointer(true)},
		}
	default:
		return fmt.Errorf("unexpected bootloader type %q. %w", bootloader, common.ErrUnknownType)
	}
	return nil
}

// GetBootloaderSettings returns a portion of features related to d8 VM bootloader.
func (b *KVVM) GetBootloaderSettings() map[string]interface{} {
	return map[string]interface{}{
		"firmare": b.Resource.Spec.Template.Spec.Domain.Firmware,
		"features": map[string]interface{}{
			"smm": b.Resource.Spec.Template.Spec.Domain.Features.SMM,
		},
	}
}
