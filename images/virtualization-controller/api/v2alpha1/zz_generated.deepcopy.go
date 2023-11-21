//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright 2022 Flant JSC

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

// Code generated by deepcopy-gen. DO NOT EDIT.

package v2alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BlockDeviceAttachmentBlockDevice) DeepCopyInto(out *BlockDeviceAttachmentBlockDevice) {
	*out = *in
	if in.VirtualMachineDisk != nil {
		in, out := &in.VirtualMachineDisk, &out.VirtualMachineDisk
		*out = new(BlockDeviceAttachmentVirtualMachineDisk)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BlockDeviceAttachmentBlockDevice.
func (in *BlockDeviceAttachmentBlockDevice) DeepCopy() *BlockDeviceAttachmentBlockDevice {
	if in == nil {
		return nil
	}
	out := new(BlockDeviceAttachmentBlockDevice)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BlockDeviceAttachmentVirtualMachineDisk) DeepCopyInto(out *BlockDeviceAttachmentVirtualMachineDisk) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BlockDeviceAttachmentVirtualMachineDisk.
func (in *BlockDeviceAttachmentVirtualMachineDisk) DeepCopy() *BlockDeviceAttachmentVirtualMachineDisk {
	if in == nil {
		return nil
	}
	out := new(BlockDeviceAttachmentVirtualMachineDisk)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BlockDeviceSpec) DeepCopyInto(out *BlockDeviceSpec) {
	*out = *in
	if in.VirtualMachineImage != nil {
		in, out := &in.VirtualMachineImage, &out.VirtualMachineImage
		*out = new(ImageDeviceSpec)
		**out = **in
	}
	if in.ClusterVirtualMachineImage != nil {
		in, out := &in.ClusterVirtualMachineImage, &out.ClusterVirtualMachineImage
		*out = new(ClusterImageDeviceSpec)
		**out = **in
	}
	if in.VirtualMachineDisk != nil {
		in, out := &in.VirtualMachineDisk, &out.VirtualMachineDisk
		*out = new(DiskDeviceSpec)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BlockDeviceSpec.
func (in *BlockDeviceSpec) DeepCopy() *BlockDeviceSpec {
	if in == nil {
		return nil
	}
	out := new(BlockDeviceSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BlockDeviceStatus) DeepCopyInto(out *BlockDeviceStatus) {
	*out = *in
	if in.VirtualMachineImage != nil {
		in, out := &in.VirtualMachineImage, &out.VirtualMachineImage
		*out = new(ImageDeviceSpec)
		**out = **in
	}
	if in.ClusterVirtualMachineImage != nil {
		in, out := &in.ClusterVirtualMachineImage, &out.ClusterVirtualMachineImage
		*out = new(ClusterImageDeviceSpec)
		**out = **in
	}
	if in.VirtualMachineDisk != nil {
		in, out := &in.VirtualMachineDisk, &out.VirtualMachineDisk
		*out = new(DiskDeviceSpec)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BlockDeviceStatus.
func (in *BlockDeviceStatus) DeepCopy() *BlockDeviceStatus {
	if in == nil {
		return nil
	}
	out := new(BlockDeviceStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CPUSpec) DeepCopyInto(out *CPUSpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CPUSpec.
func (in *CPUSpec) DeepCopy() *CPUSpec {
	if in == nil {
		return nil
	}
	out := new(CPUSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Checksum) DeepCopyInto(out *Checksum) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Checksum.
func (in *Checksum) DeepCopy() *Checksum {
	if in == nil {
		return nil
	}
	out := new(Checksum)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterImageDeviceSpec) DeepCopyInto(out *ClusterImageDeviceSpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterImageDeviceSpec.
func (in *ClusterImageDeviceSpec) DeepCopy() *ClusterImageDeviceSpec {
	if in == nil {
		return nil
	}
	out := new(ClusterImageDeviceSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterVirtualMachineImage) DeepCopyInto(out *ClusterVirtualMachineImage) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterVirtualMachineImage.
func (in *ClusterVirtualMachineImage) DeepCopy() *ClusterVirtualMachineImage {
	if in == nil {
		return nil
	}
	out := new(ClusterVirtualMachineImage)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ClusterVirtualMachineImage) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterVirtualMachineImageList) DeepCopyInto(out *ClusterVirtualMachineImageList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ClusterVirtualMachineImage, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterVirtualMachineImageList.
func (in *ClusterVirtualMachineImageList) DeepCopy() *ClusterVirtualMachineImageList {
	if in == nil {
		return nil
	}
	out := new(ClusterVirtualMachineImageList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ClusterVirtualMachineImageList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterVirtualMachineImageSpec) DeepCopyInto(out *ClusterVirtualMachineImageSpec) {
	*out = *in
	in.DataSource.DeepCopyInto(&out.DataSource)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterVirtualMachineImageSpec.
func (in *ClusterVirtualMachineImageSpec) DeepCopy() *ClusterVirtualMachineImageSpec {
	if in == nil {
		return nil
	}
	out := new(ClusterVirtualMachineImageSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterVirtualMachineImageStatus) DeepCopyInto(out *ClusterVirtualMachineImageStatus) {
	*out = *in
	out.ImageStatus = in.ImageStatus
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterVirtualMachineImageStatus.
func (in *ClusterVirtualMachineImageStatus) DeepCopy() *ClusterVirtualMachineImageStatus {
	if in == nil {
		return nil
	}
	out := new(ClusterVirtualMachineImageStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DataSource) DeepCopyInto(out *DataSource) {
	*out = *in
	if in.HTTP != nil {
		in, out := &in.HTTP, &out.HTTP
		*out = new(DataSourceHTTP)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DataSource.
func (in *DataSource) DeepCopy() *DataSource {
	if in == nil {
		return nil
	}
	out := new(DataSource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DataSourceHTTP) DeepCopyInto(out *DataSourceHTTP) {
	*out = *in
	if in.CABundle != nil {
		in, out := &in.CABundle, &out.CABundle
		*out = make([]byte, len(*in))
		copy(*out, *in)
	}
	if in.Checksum != nil {
		in, out := &in.Checksum, &out.Checksum
		*out = new(Checksum)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DataSourceHTTP.
func (in *DataSourceHTTP) DeepCopy() *DataSourceHTTP {
	if in == nil {
		return nil
	}
	out := new(DataSourceHTTP)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DiskDeviceSpec) DeepCopyInto(out *DiskDeviceSpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DiskDeviceSpec.
func (in *DiskDeviceSpec) DeepCopy() *DiskDeviceSpec {
	if in == nil {
		return nil
	}
	out := new(DiskDeviceSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DiskTarget) DeepCopyInto(out *DiskTarget) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DiskTarget.
func (in *DiskTarget) DeepCopy() *DiskTarget {
	if in == nil {
		return nil
	}
	out := new(DiskTarget)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ImageDeviceSpec) DeepCopyInto(out *ImageDeviceSpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ImageDeviceSpec.
func (in *ImageDeviceSpec) DeepCopy() *ImageDeviceSpec {
	if in == nil {
		return nil
	}
	out := new(ImageDeviceSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ImageStatus) DeepCopyInto(out *ImageStatus) {
	*out = *in
	out.DownloadSpeed = in.DownloadSpeed
	out.Size = in.Size
	out.Target = in.Target
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ImageStatus.
func (in *ImageStatus) DeepCopy() *ImageStatus {
	if in == nil {
		return nil
	}
	out := new(ImageStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ImageStatusSize) DeepCopyInto(out *ImageStatusSize) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ImageStatusSize.
func (in *ImageStatusSize) DeepCopy() *ImageStatusSize {
	if in == nil {
		return nil
	}
	out := new(ImageStatusSize)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ImageStatusSpeed) DeepCopyInto(out *ImageStatusSpeed) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ImageStatusSpeed.
func (in *ImageStatusSpeed) DeepCopy() *ImageStatusSpeed {
	if in == nil {
		return nil
	}
	out := new(ImageStatusSpeed)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ImageStatusTarget) DeepCopyInto(out *ImageStatusTarget) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ImageStatusTarget.
func (in *ImageStatusTarget) DeepCopy() *ImageStatusTarget {
	if in == nil {
		return nil
	}
	out := new(ImageStatusTarget)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MemorySpec) DeepCopyInto(out *MemorySpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MemorySpec.
func (in *MemorySpec) DeepCopy() *MemorySpec {
	if in == nil {
		return nil
	}
	out := new(MemorySpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VMDDownloadSpeed) DeepCopyInto(out *VMDDownloadSpeed) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VMDDownloadSpeed.
func (in *VMDDownloadSpeed) DeepCopy() *VMDDownloadSpeed {
	if in == nil {
		return nil
	}
	out := new(VMDDownloadSpeed)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VMDPersistentVolumeClaim) DeepCopyInto(out *VMDPersistentVolumeClaim) {
	*out = *in
	out.Size = in.Size.DeepCopy()
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VMDPersistentVolumeClaim.
func (in *VMDPersistentVolumeClaim) DeepCopy() *VMDPersistentVolumeClaim {
	if in == nil {
		return nil
	}
	out := new(VMDPersistentVolumeClaim)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VMIPersistentVolumeClaim) DeepCopyInto(out *VMIPersistentVolumeClaim) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VMIPersistentVolumeClaim.
func (in *VMIPersistentVolumeClaim) DeepCopy() *VMIPersistentVolumeClaim {
	if in == nil {
		return nil
	}
	out := new(VMIPersistentVolumeClaim)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VirtualMachine) DeepCopyInto(out *VirtualMachine) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VirtualMachine.
func (in *VirtualMachine) DeepCopy() *VirtualMachine {
	if in == nil {
		return nil
	}
	out := new(VirtualMachine)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *VirtualMachine) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VirtualMachineBlockDeviceAttachment) DeepCopyInto(out *VirtualMachineBlockDeviceAttachment) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VirtualMachineBlockDeviceAttachment.
func (in *VirtualMachineBlockDeviceAttachment) DeepCopy() *VirtualMachineBlockDeviceAttachment {
	if in == nil {
		return nil
	}
	out := new(VirtualMachineBlockDeviceAttachment)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *VirtualMachineBlockDeviceAttachment) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VirtualMachineBlockDeviceAttachmentList) DeepCopyInto(out *VirtualMachineBlockDeviceAttachmentList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]VirtualMachineBlockDeviceAttachment, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VirtualMachineBlockDeviceAttachmentList.
func (in *VirtualMachineBlockDeviceAttachmentList) DeepCopy() *VirtualMachineBlockDeviceAttachmentList {
	if in == nil {
		return nil
	}
	out := new(VirtualMachineBlockDeviceAttachmentList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *VirtualMachineBlockDeviceAttachmentList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VirtualMachineBlockDeviceAttachmentSpec) DeepCopyInto(out *VirtualMachineBlockDeviceAttachmentSpec) {
	*out = *in
	in.BlockDevice.DeepCopyInto(&out.BlockDevice)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VirtualMachineBlockDeviceAttachmentSpec.
func (in *VirtualMachineBlockDeviceAttachmentSpec) DeepCopy() *VirtualMachineBlockDeviceAttachmentSpec {
	if in == nil {
		return nil
	}
	out := new(VirtualMachineBlockDeviceAttachmentSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VirtualMachineBlockDeviceAttachmentStatus) DeepCopyInto(out *VirtualMachineBlockDeviceAttachmentStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VirtualMachineBlockDeviceAttachmentStatus.
func (in *VirtualMachineBlockDeviceAttachmentStatus) DeepCopy() *VirtualMachineBlockDeviceAttachmentStatus {
	if in == nil {
		return nil
	}
	out := new(VirtualMachineBlockDeviceAttachmentStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VirtualMachineDisk) DeepCopyInto(out *VirtualMachineDisk) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VirtualMachineDisk.
func (in *VirtualMachineDisk) DeepCopy() *VirtualMachineDisk {
	if in == nil {
		return nil
	}
	out := new(VirtualMachineDisk)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *VirtualMachineDisk) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VirtualMachineDiskList) DeepCopyInto(out *VirtualMachineDiskList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]VirtualMachineDisk, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VirtualMachineDiskList.
func (in *VirtualMachineDiskList) DeepCopy() *VirtualMachineDiskList {
	if in == nil {
		return nil
	}
	out := new(VirtualMachineDiskList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *VirtualMachineDiskList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VirtualMachineDiskSpec) DeepCopyInto(out *VirtualMachineDiskSpec) {
	*out = *in
	if in.DataSource != nil {
		in, out := &in.DataSource, &out.DataSource
		*out = new(DataSource)
		(*in).DeepCopyInto(*out)
	}
	in.PersistentVolumeClaim.DeepCopyInto(&out.PersistentVolumeClaim)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VirtualMachineDiskSpec.
func (in *VirtualMachineDiskSpec) DeepCopy() *VirtualMachineDiskSpec {
	if in == nil {
		return nil
	}
	out := new(VirtualMachineDiskSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VirtualMachineDiskStatus) DeepCopyInto(out *VirtualMachineDiskStatus) {
	*out = *in
	out.DownloadSpeed = in.DownloadSpeed
	out.Target = in.Target
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VirtualMachineDiskStatus.
func (in *VirtualMachineDiskStatus) DeepCopy() *VirtualMachineDiskStatus {
	if in == nil {
		return nil
	}
	out := new(VirtualMachineDiskStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VirtualMachineIPAddressClaim) DeepCopyInto(out *VirtualMachineIPAddressClaim) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VirtualMachineIPAddressClaim.
func (in *VirtualMachineIPAddressClaim) DeepCopy() *VirtualMachineIPAddressClaim {
	if in == nil {
		return nil
	}
	out := new(VirtualMachineIPAddressClaim)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *VirtualMachineIPAddressClaim) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VirtualMachineIPAddressClaimList) DeepCopyInto(out *VirtualMachineIPAddressClaimList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]VirtualMachineIPAddressClaim, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VirtualMachineIPAddressClaimList.
func (in *VirtualMachineIPAddressClaimList) DeepCopy() *VirtualMachineIPAddressClaimList {
	if in == nil {
		return nil
	}
	out := new(VirtualMachineIPAddressClaimList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *VirtualMachineIPAddressClaimList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VirtualMachineIPAddressClaimSpec) DeepCopyInto(out *VirtualMachineIPAddressClaimSpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VirtualMachineIPAddressClaimSpec.
func (in *VirtualMachineIPAddressClaimSpec) DeepCopy() *VirtualMachineIPAddressClaimSpec {
	if in == nil {
		return nil
	}
	out := new(VirtualMachineIPAddressClaimSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VirtualMachineIPAddressClaimStatus) DeepCopyInto(out *VirtualMachineIPAddressClaimStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VirtualMachineIPAddressClaimStatus.
func (in *VirtualMachineIPAddressClaimStatus) DeepCopy() *VirtualMachineIPAddressClaimStatus {
	if in == nil {
		return nil
	}
	out := new(VirtualMachineIPAddressClaimStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VirtualMachineIPAddressLease) DeepCopyInto(out *VirtualMachineIPAddressLease) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VirtualMachineIPAddressLease.
func (in *VirtualMachineIPAddressLease) DeepCopy() *VirtualMachineIPAddressLease {
	if in == nil {
		return nil
	}
	out := new(VirtualMachineIPAddressLease)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *VirtualMachineIPAddressLease) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VirtualMachineIPAddressLeaseClaimRef) DeepCopyInto(out *VirtualMachineIPAddressLeaseClaimRef) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VirtualMachineIPAddressLeaseClaimRef.
func (in *VirtualMachineIPAddressLeaseClaimRef) DeepCopy() *VirtualMachineIPAddressLeaseClaimRef {
	if in == nil {
		return nil
	}
	out := new(VirtualMachineIPAddressLeaseClaimRef)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VirtualMachineIPAddressLeaseList) DeepCopyInto(out *VirtualMachineIPAddressLeaseList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]VirtualMachineIPAddressLease, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VirtualMachineIPAddressLeaseList.
func (in *VirtualMachineIPAddressLeaseList) DeepCopy() *VirtualMachineIPAddressLeaseList {
	if in == nil {
		return nil
	}
	out := new(VirtualMachineIPAddressLeaseList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *VirtualMachineIPAddressLeaseList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VirtualMachineIPAddressLeaseSpec) DeepCopyInto(out *VirtualMachineIPAddressLeaseSpec) {
	*out = *in
	if in.ClaimRef != nil {
		in, out := &in.ClaimRef, &out.ClaimRef
		*out = new(VirtualMachineIPAddressLeaseClaimRef)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VirtualMachineIPAddressLeaseSpec.
func (in *VirtualMachineIPAddressLeaseSpec) DeepCopy() *VirtualMachineIPAddressLeaseSpec {
	if in == nil {
		return nil
	}
	out := new(VirtualMachineIPAddressLeaseSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VirtualMachineIPAddressLeaseStatus) DeepCopyInto(out *VirtualMachineIPAddressLeaseStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VirtualMachineIPAddressLeaseStatus.
func (in *VirtualMachineIPAddressLeaseStatus) DeepCopy() *VirtualMachineIPAddressLeaseStatus {
	if in == nil {
		return nil
	}
	out := new(VirtualMachineIPAddressLeaseStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VirtualMachineImage) DeepCopyInto(out *VirtualMachineImage) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VirtualMachineImage.
func (in *VirtualMachineImage) DeepCopy() *VirtualMachineImage {
	if in == nil {
		return nil
	}
	out := new(VirtualMachineImage)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *VirtualMachineImage) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VirtualMachineImageList) DeepCopyInto(out *VirtualMachineImageList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]VirtualMachineImage, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VirtualMachineImageList.
func (in *VirtualMachineImageList) DeepCopy() *VirtualMachineImageList {
	if in == nil {
		return nil
	}
	out := new(VirtualMachineImageList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *VirtualMachineImageList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VirtualMachineImageSpec) DeepCopyInto(out *VirtualMachineImageSpec) {
	*out = *in
	out.PersistentVolumeClaim = in.PersistentVolumeClaim
	in.DataSource.DeepCopyInto(&out.DataSource)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VirtualMachineImageSpec.
func (in *VirtualMachineImageSpec) DeepCopy() *VirtualMachineImageSpec {
	if in == nil {
		return nil
	}
	out := new(VirtualMachineImageSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VirtualMachineImageStatus) DeepCopyInto(out *VirtualMachineImageStatus) {
	*out = *in
	out.ImageStatus = in.ImageStatus
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VirtualMachineImageStatus.
func (in *VirtualMachineImageStatus) DeepCopy() *VirtualMachineImageStatus {
	if in == nil {
		return nil
	}
	out := new(VirtualMachineImageStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VirtualMachineList) DeepCopyInto(out *VirtualMachineList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]VirtualMachine, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VirtualMachineList.
func (in *VirtualMachineList) DeepCopy() *VirtualMachineList {
	if in == nil {
		return nil
	}
	out := new(VirtualMachineList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *VirtualMachineList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VirtualMachineSpec) DeepCopyInto(out *VirtualMachineSpec) {
	*out = *in
	out.CPU = in.CPU
	out.Memory = in.Memory
	if in.BlockDevices != nil {
		in, out := &in.BlockDevices, &out.BlockDevices
		*out = make([]BlockDeviceSpec, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VirtualMachineSpec.
func (in *VirtualMachineSpec) DeepCopy() *VirtualMachineSpec {
	if in == nil {
		return nil
	}
	out := new(VirtualMachineSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VirtualMachineStatus) DeepCopyInto(out *VirtualMachineStatus) {
	*out = *in
	if in.BlockDevicesAttached != nil {
		in, out := &in.BlockDevicesAttached, &out.BlockDevicesAttached
		*out = make([]BlockDeviceStatus, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	out.GuestOSInfo = in.GuestOSInfo
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VirtualMachineStatus.
func (in *VirtualMachineStatus) DeepCopy() *VirtualMachineStatus {
	if in == nil {
		return nil
	}
	out := new(VirtualMachineStatus)
	in.DeepCopyInto(out)
	return out
}
