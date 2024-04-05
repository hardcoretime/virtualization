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
// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v1alpha2 "github.com/deckhouse/virtualization/api/core/v1alpha2"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeVirtualMachineIPAddressLeases implements VirtualMachineIPAddressLeaseInterface
type FakeVirtualMachineIPAddressLeases struct {
	Fake *FakeVirtualizationV1alpha2
}

var virtualmachineipaddressleasesResource = v1alpha2.SchemeGroupVersion.WithResource("virtualmachineipaddressleases")

var virtualmachineipaddressleasesKind = v1alpha2.SchemeGroupVersion.WithKind("VirtualMachineIPAddressLease")

// Get takes name of the virtualMachineIPAddressLease, and returns the corresponding virtualMachineIPAddressLease object, and an error if there is any.
func (c *FakeVirtualMachineIPAddressLeases) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha2.VirtualMachineIPAddressLease, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(virtualmachineipaddressleasesResource, name), &v1alpha2.VirtualMachineIPAddressLease{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha2.VirtualMachineIPAddressLease), err
}

// List takes label and field selectors, and returns the list of VirtualMachineIPAddressLeases that match those selectors.
func (c *FakeVirtualMachineIPAddressLeases) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha2.VirtualMachineIPAddressLeaseList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(virtualmachineipaddressleasesResource, virtualmachineipaddressleasesKind, opts), &v1alpha2.VirtualMachineIPAddressLeaseList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha2.VirtualMachineIPAddressLeaseList{ListMeta: obj.(*v1alpha2.VirtualMachineIPAddressLeaseList).ListMeta}
	for _, item := range obj.(*v1alpha2.VirtualMachineIPAddressLeaseList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested virtualMachineIPAddressLeases.
func (c *FakeVirtualMachineIPAddressLeases) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(virtualmachineipaddressleasesResource, opts))
}

// Create takes the representation of a virtualMachineIPAddressLease and creates it.  Returns the server's representation of the virtualMachineIPAddressLease, and an error, if there is any.
func (c *FakeVirtualMachineIPAddressLeases) Create(ctx context.Context, virtualMachineIPAddressLease *v1alpha2.VirtualMachineIPAddressLease, opts v1.CreateOptions) (result *v1alpha2.VirtualMachineIPAddressLease, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(virtualmachineipaddressleasesResource, virtualMachineIPAddressLease), &v1alpha2.VirtualMachineIPAddressLease{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha2.VirtualMachineIPAddressLease), err
}

// Update takes the representation of a virtualMachineIPAddressLease and updates it. Returns the server's representation of the virtualMachineIPAddressLease, and an error, if there is any.
func (c *FakeVirtualMachineIPAddressLeases) Update(ctx context.Context, virtualMachineIPAddressLease *v1alpha2.VirtualMachineIPAddressLease, opts v1.UpdateOptions) (result *v1alpha2.VirtualMachineIPAddressLease, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(virtualmachineipaddressleasesResource, virtualMachineIPAddressLease), &v1alpha2.VirtualMachineIPAddressLease{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha2.VirtualMachineIPAddressLease), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeVirtualMachineIPAddressLeases) UpdateStatus(ctx context.Context, virtualMachineIPAddressLease *v1alpha2.VirtualMachineIPAddressLease, opts v1.UpdateOptions) (*v1alpha2.VirtualMachineIPAddressLease, error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateSubresourceAction(virtualmachineipaddressleasesResource, "status", virtualMachineIPAddressLease), &v1alpha2.VirtualMachineIPAddressLease{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha2.VirtualMachineIPAddressLease), err
}

// Delete takes name of the virtualMachineIPAddressLease and deletes it. Returns an error if one occurs.
func (c *FakeVirtualMachineIPAddressLeases) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteActionWithOptions(virtualmachineipaddressleasesResource, name, opts), &v1alpha2.VirtualMachineIPAddressLease{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeVirtualMachineIPAddressLeases) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(virtualmachineipaddressleasesResource, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha2.VirtualMachineIPAddressLeaseList{})
	return err
}

// Patch applies the patch and returns the patched virtualMachineIPAddressLease.
func (c *FakeVirtualMachineIPAddressLeases) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha2.VirtualMachineIPAddressLease, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(virtualmachineipaddressleasesResource, name, pt, data, subresources...), &v1alpha2.VirtualMachineIPAddressLease{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha2.VirtualMachineIPAddressLease), err
}