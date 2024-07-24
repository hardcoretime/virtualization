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

package common

import (
	"net"
	"strings"

	virtv2 "github.com/deckhouse/virtualization/api/core/v1alpha2"
)

const ipPrefix = "ip-"

type AllocatedIPs map[string]*virtv2.VirtualMachineIPAddressLease

// IpToLeaseName generate the Virtual Machine IP Address Lease's name from the ip address
func IpToLeaseName(ip string) string {
	addr := net.ParseIP(ip)
	if addr.To4() != nil {
		// IPv4 address
		return ipPrefix + strings.ReplaceAll(addr.String(), ".", "-")
	}

	return ""
}

// LeaseNameToIP generate the ip address from the Virtual Machine IP Address Lease's name
func LeaseNameToIP(leaseName string) string {
	if strings.HasPrefix(leaseName, ipPrefix) && len(leaseName) > len(ipPrefix) {
		return strings.ReplaceAll(leaseName[len(ipPrefix):], "-", ".")
	}

	return ""
}
