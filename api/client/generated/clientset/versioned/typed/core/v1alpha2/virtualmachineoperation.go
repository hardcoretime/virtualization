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

	scheme "github.com/deckhouse/virtualization/api/client/generated/clientset/versioned/scheme"
	v1alpha2 "github.com/deckhouse/virtualization/api/core/v1alpha2"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// VirtualMachineOperationsGetter has a method to return a VirtualMachineOperationInterface.
// A group's client should implement this interface.
type VirtualMachineOperationsGetter interface {
	VirtualMachineOperations(namespace string) VirtualMachineOperationInterface
}

// VirtualMachineOperationInterface has methods to work with VirtualMachineOperation resources.
type VirtualMachineOperationInterface interface {
	Create(ctx context.Context, virtualMachineOperation *v1alpha2.VirtualMachineOperation, opts v1.CreateOptions) (*v1alpha2.VirtualMachineOperation, error)
	Update(ctx context.Context, virtualMachineOperation *v1alpha2.VirtualMachineOperation, opts v1.UpdateOptions) (*v1alpha2.VirtualMachineOperation, error)
	UpdateStatus(ctx context.Context, virtualMachineOperation *v1alpha2.VirtualMachineOperation, opts v1.UpdateOptions) (*v1alpha2.VirtualMachineOperation, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1alpha2.VirtualMachineOperation, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1alpha2.VirtualMachineOperationList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha2.VirtualMachineOperation, err error)
	VirtualMachineOperationExpansion
}

// virtualMachineOperations implements VirtualMachineOperationInterface
type virtualMachineOperations struct {
	client rest.Interface
	ns     string
}

// newVirtualMachineOperations returns a VirtualMachineOperations
func newVirtualMachineOperations(c *VirtualizationV1alpha2Client, namespace string) *virtualMachineOperations {
	return &virtualMachineOperations{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the virtualMachineOperation, and returns the corresponding virtualMachineOperation object, and an error if there is any.
func (c *virtualMachineOperations) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha2.VirtualMachineOperation, err error) {
	result = &v1alpha2.VirtualMachineOperation{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("virtualmachineoperations").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of VirtualMachineOperations that match those selectors.
func (c *virtualMachineOperations) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha2.VirtualMachineOperationList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha2.VirtualMachineOperationList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("virtualmachineoperations").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested virtualMachineOperations.
func (c *virtualMachineOperations) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("virtualmachineoperations").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a virtualMachineOperation and creates it.  Returns the server's representation of the virtualMachineOperation, and an error, if there is any.
func (c *virtualMachineOperations) Create(ctx context.Context, virtualMachineOperation *v1alpha2.VirtualMachineOperation, opts v1.CreateOptions) (result *v1alpha2.VirtualMachineOperation, err error) {
	result = &v1alpha2.VirtualMachineOperation{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("virtualmachineoperations").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(virtualMachineOperation).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a virtualMachineOperation and updates it. Returns the server's representation of the virtualMachineOperation, and an error, if there is any.
func (c *virtualMachineOperations) Update(ctx context.Context, virtualMachineOperation *v1alpha2.VirtualMachineOperation, opts v1.UpdateOptions) (result *v1alpha2.VirtualMachineOperation, err error) {
	result = &v1alpha2.VirtualMachineOperation{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("virtualmachineoperations").
		Name(virtualMachineOperation.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(virtualMachineOperation).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *virtualMachineOperations) UpdateStatus(ctx context.Context, virtualMachineOperation *v1alpha2.VirtualMachineOperation, opts v1.UpdateOptions) (result *v1alpha2.VirtualMachineOperation, err error) {
	result = &v1alpha2.VirtualMachineOperation{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("virtualmachineoperations").
		Name(virtualMachineOperation.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(virtualMachineOperation).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the virtualMachineOperation and deletes it. Returns an error if one occurs.
func (c *virtualMachineOperations) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("virtualmachineoperations").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *virtualMachineOperations) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("virtualmachineoperations").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched virtualMachineOperation.
func (c *virtualMachineOperations) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha2.VirtualMachineOperation, err error) {
	result = &v1alpha2.VirtualMachineOperation{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("virtualmachineoperations").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}