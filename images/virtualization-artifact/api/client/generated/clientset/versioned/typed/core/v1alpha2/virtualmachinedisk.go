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

package v1alpha2

import (
	"context"
	"time"

	scheme "github.com/deckhouse/virtualization-controller/api/client/generated/clientset/versioned/scheme"
	v1alpha2 "github.com/deckhouse/virtualization-controller/api/core/v1alpha2"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// VirtualMachineDisksGetter has a method to return a VirtualMachineDiskInterface.
// A group's client should implement this interface.
type VirtualMachineDisksGetter interface {
	VirtualMachineDisks(namespace string) VirtualMachineDiskInterface
}

// VirtualMachineDiskInterface has methods to work with VirtualMachineDisk resources.
type VirtualMachineDiskInterface interface {
	Create(ctx context.Context, virtualMachineDisk *v1alpha2.VirtualMachineDisk, opts v1.CreateOptions) (*v1alpha2.VirtualMachineDisk, error)
	Update(ctx context.Context, virtualMachineDisk *v1alpha2.VirtualMachineDisk, opts v1.UpdateOptions) (*v1alpha2.VirtualMachineDisk, error)
	UpdateStatus(ctx context.Context, virtualMachineDisk *v1alpha2.VirtualMachineDisk, opts v1.UpdateOptions) (*v1alpha2.VirtualMachineDisk, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1alpha2.VirtualMachineDisk, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1alpha2.VirtualMachineDiskList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha2.VirtualMachineDisk, err error)
	VirtualMachineDiskExpansion
}

// virtualMachineDisks implements VirtualMachineDiskInterface
type virtualMachineDisks struct {
	client rest.Interface
	ns     string
}

// newVirtualMachineDisks returns a VirtualMachineDisks
func newVirtualMachineDisks(c *VirtualizationV1alpha2Client, namespace string) *virtualMachineDisks {
	return &virtualMachineDisks{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the virtualMachineDisk, and returns the corresponding virtualMachineDisk object, and an error if there is any.
func (c *virtualMachineDisks) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha2.VirtualMachineDisk, err error) {
	result = &v1alpha2.VirtualMachineDisk{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("virtualmachinedisks").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of VirtualMachineDisks that match those selectors.
func (c *virtualMachineDisks) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha2.VirtualMachineDiskList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha2.VirtualMachineDiskList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("virtualmachinedisks").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested virtualMachineDisks.
func (c *virtualMachineDisks) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("virtualmachinedisks").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a virtualMachineDisk and creates it.  Returns the server's representation of the virtualMachineDisk, and an error, if there is any.
func (c *virtualMachineDisks) Create(ctx context.Context, virtualMachineDisk *v1alpha2.VirtualMachineDisk, opts v1.CreateOptions) (result *v1alpha2.VirtualMachineDisk, err error) {
	result = &v1alpha2.VirtualMachineDisk{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("virtualmachinedisks").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(virtualMachineDisk).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a virtualMachineDisk and updates it. Returns the server's representation of the virtualMachineDisk, and an error, if there is any.
func (c *virtualMachineDisks) Update(ctx context.Context, virtualMachineDisk *v1alpha2.VirtualMachineDisk, opts v1.UpdateOptions) (result *v1alpha2.VirtualMachineDisk, err error) {
	result = &v1alpha2.VirtualMachineDisk{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("virtualmachinedisks").
		Name(virtualMachineDisk.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(virtualMachineDisk).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *virtualMachineDisks) UpdateStatus(ctx context.Context, virtualMachineDisk *v1alpha2.VirtualMachineDisk, opts v1.UpdateOptions) (result *v1alpha2.VirtualMachineDisk, err error) {
	result = &v1alpha2.VirtualMachineDisk{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("virtualmachinedisks").
		Name(virtualMachineDisk.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(virtualMachineDisk).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the virtualMachineDisk and deletes it. Returns an error if one occurs.
func (c *virtualMachineDisks) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("virtualmachinedisks").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *virtualMachineDisks) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("virtualmachinedisks").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched virtualMachineDisk.
func (c *virtualMachineDisks) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha2.VirtualMachineDisk, err error) {
	result = &v1alpha2.VirtualMachineDisk{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("virtualmachinedisks").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
