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

package util

import (
	"context"
	"fmt"
	"net"
	"strings"

	k8snet "k8s.io/utils/net"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/deckhouse/virtualization-controller/pkg/controller/service"
	virtv2 "github.com/deckhouse/virtualization/api/core/v1alpha2"
)

type AllocatedIPs map[string]*virtv2.VirtualMachineIPAddressLease

func GetAllocatedIPs(ctx context.Context, apiClient client.Client) (AllocatedIPs, error) {
	var leases virtv2.VirtualMachineIPAddressLeaseList

	err := apiClient.List(ctx, &leases)
	if err != nil {
		return nil, fmt.Errorf("error getting leases: %w", err)
	}

	allocatedIPs := make(AllocatedIPs, len(leases.Items))
	for _, lease := range leases.Items {
		l := lease
		allocatedIPs[LeaseNameToIP(lease.Name)] = &l
	}

	return allocatedIPs, nil
}

const ipPrefix = "ip-"

func LeaseNameToIP(leaseName string) string {

	if strings.HasPrefix(leaseName, ipPrefix) && len(leaseName) > len(ipPrefix) {
		return strings.ReplaceAll(leaseName[len(ipPrefix):], "-", ".")
	}

	return ""
}

func IpToLeaseName(ip string) string {
	addr := net.ParseIP(ip)
	if addr.To4() != nil {
		// IPv4 address
		return ipPrefix + strings.ReplaceAll(addr.String(), ".", "-")
	}

	return ""
}

func IsFirstLastIP(ip net.IP, cidr *net.IPNet) (bool, error) {
	size := int(k8snet.RangeSize(cidr))

	first, err := k8snet.GetIndexedIP(cidr, 0)
	if err != nil {
		return false, err
	}

	if first.Equal(ip) {
		return true, nil
	}

	last, err := k8snet.GetIndexedIP(cidr, size-1)
	if err != nil {
		return false, err
	}

	return last.Equal(ip), nil
}

func IsBoundLease(vmipl *virtv2.VirtualMachineIPAddressLease, vmip *service.Resource[*virtv2.VirtualMachineIPAddress, virtv2.VirtualMachineIPAddressStatus]) bool {
	if vmipl.Status.Phase != virtv2.VirtualMachineIPAddressLeasePhaseBound {
		return false
	}

	if vmipl.Spec.IpAddressRef == nil {
		return false
	}

	if vmipl.Spec.IpAddressRef.Namespace != vmip.Name().Namespace || vmipl.Spec.IpAddressRef.Name != vmip.Name().Name {
		return false
	}

	return true
}
