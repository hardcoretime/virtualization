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
	"net/http"

	"github.com/deckhouse/virtualization/api/client/generated/clientset/versioned/scheme"
	v1alpha2 "github.com/deckhouse/virtualization/api/core/v1alpha2"
	rest "k8s.io/client-go/rest"
)

type VirtualizationV1alpha2Interface interface {
	RESTClient() rest.Interface
	ClusterVirtualImagesGetter
	VirtualDisksGetter
	VirtualImagesGetter
	VirtualMachinesGetter
	VirtualMachineBlockDeviceAttachmentsGetter
	VirtualMachineCPUModelsGetter
	VirtualMachineIPAddressClaimsGetter
	VirtualMachineIPAddressLeasesGetter
	VirtualMachineOperationsGetter
}

// VirtualizationV1alpha2Client is used to interact with features provided by the virtualization.deckhouse.io group.
type VirtualizationV1alpha2Client struct {
	restClient rest.Interface
}

func (c *VirtualizationV1alpha2Client) ClusterVirtualImages() ClusterVirtualImageInterface {
	return newClusterVirtualImages(c)
}

func (c *VirtualizationV1alpha2Client) VirtualDisks(namespace string) VirtualDiskInterface {
	return newVirtualDisks(c, namespace)
}

func (c *VirtualizationV1alpha2Client) VirtualImages(namespace string) VirtualImageInterface {
	return newVirtualImages(c, namespace)
}

func (c *VirtualizationV1alpha2Client) VirtualMachines(namespace string) VirtualMachineInterface {
	return newVirtualMachines(c, namespace)
}

func (c *VirtualizationV1alpha2Client) VirtualMachineBlockDeviceAttachments(namespace string) VirtualMachineBlockDeviceAttachmentInterface {
	return newVirtualMachineBlockDeviceAttachments(c, namespace)
}

func (c *VirtualizationV1alpha2Client) VirtualMachineCPUModels() VirtualMachineCPUModelInterface {
	return newVirtualMachineCPUModels(c)
}

func (c *VirtualizationV1alpha2Client) VirtualMachineIPAddressClaims(namespace string) VirtualMachineIPAddressClaimInterface {
	return newVirtualMachineIPAddressClaims(c, namespace)
}

func (c *VirtualizationV1alpha2Client) VirtualMachineIPAddressLeases() VirtualMachineIPAddressLeaseInterface {
	return newVirtualMachineIPAddressLeases(c)
}

func (c *VirtualizationV1alpha2Client) VirtualMachineOperations(namespace string) VirtualMachineOperationInterface {
	return newVirtualMachineOperations(c, namespace)
}

// NewForConfig creates a new VirtualizationV1alpha2Client for the given config.
// NewForConfig is equivalent to NewForConfigAndClient(c, httpClient),
// where httpClient was generated with rest.HTTPClientFor(c).
func NewForConfig(c *rest.Config) (*VirtualizationV1alpha2Client, error) {
	config := *c
	if err := setConfigDefaults(&config); err != nil {
		return nil, err
	}
	httpClient, err := rest.HTTPClientFor(&config)
	if err != nil {
		return nil, err
	}
	return NewForConfigAndClient(&config, httpClient)
}

// NewForConfigAndClient creates a new VirtualizationV1alpha2Client for the given config and http client.
// Note the http client provided takes precedence over the configured transport values.
func NewForConfigAndClient(c *rest.Config, h *http.Client) (*VirtualizationV1alpha2Client, error) {
	config := *c
	if err := setConfigDefaults(&config); err != nil {
		return nil, err
	}
	client, err := rest.RESTClientForConfigAndClient(&config, h)
	if err != nil {
		return nil, err
	}
	return &VirtualizationV1alpha2Client{client}, nil
}

// NewForConfigOrDie creates a new VirtualizationV1alpha2Client for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *VirtualizationV1alpha2Client {
	client, err := NewForConfig(c)
	if err != nil {
		panic(err)
	}
	return client
}

// New creates a new VirtualizationV1alpha2Client for the given RESTClient.
func New(c rest.Interface) *VirtualizationV1alpha2Client {
	return &VirtualizationV1alpha2Client{c}
}

func setConfigDefaults(config *rest.Config) error {
	gv := v1alpha2.SchemeGroupVersion
	config.GroupVersion = &gv
	config.APIPath = "/apis"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()

	if config.UserAgent == "" {
		config.UserAgent = rest.DefaultKubernetesUserAgent()
	}

	return nil
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *VirtualizationV1alpha2Client) RESTClient() rest.Interface {
	if c == nil {
		return nil
	}
	return c.restClient
}
