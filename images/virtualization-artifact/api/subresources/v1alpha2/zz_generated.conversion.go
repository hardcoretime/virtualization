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

// Code generated by conversion-gen. DO NOT EDIT.

package v1alpha2

import (
	url "net/url"

	subresources "github.com/deckhouse/virtualization-controller/api/subresources"
	conversion "k8s.io/apimachinery/pkg/conversion"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

func init() {
	localSchemeBuilder.Register(RegisterConversions)
}

// RegisterConversions adds conversion functions to the given scheme.
// Public to allow building arbitrary schemes.
func RegisterConversions(s *runtime.Scheme) error {
	if err := s.AddGeneratedConversionFunc((*VirtualMachineConsole)(nil), (*subresources.VirtualMachineConsole)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha2_VirtualMachineConsole_To_subresources_VirtualMachineConsole(a.(*VirtualMachineConsole), b.(*subresources.VirtualMachineConsole), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*subresources.VirtualMachineConsole)(nil), (*VirtualMachineConsole)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_subresources_VirtualMachineConsole_To_v1alpha2_VirtualMachineConsole(a.(*subresources.VirtualMachineConsole), b.(*VirtualMachineConsole), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*VirtualMachinePortForward)(nil), (*subresources.VirtualMachinePortForward)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha2_VirtualMachinePortForward_To_subresources_VirtualMachinePortForward(a.(*VirtualMachinePortForward), b.(*subresources.VirtualMachinePortForward), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*subresources.VirtualMachinePortForward)(nil), (*VirtualMachinePortForward)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_subresources_VirtualMachinePortForward_To_v1alpha2_VirtualMachinePortForward(a.(*subresources.VirtualMachinePortForward), b.(*VirtualMachinePortForward), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*VirtualMachineVNC)(nil), (*subresources.VirtualMachineVNC)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha2_VirtualMachineVNC_To_subresources_VirtualMachineVNC(a.(*VirtualMachineVNC), b.(*subresources.VirtualMachineVNC), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*subresources.VirtualMachineVNC)(nil), (*VirtualMachineVNC)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_subresources_VirtualMachineVNC_To_v1alpha2_VirtualMachineVNC(a.(*subresources.VirtualMachineVNC), b.(*VirtualMachineVNC), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*url.Values)(nil), (*VirtualMachineConsole)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_url_Values_To_v1alpha2_VirtualMachineConsole(a.(*url.Values), b.(*VirtualMachineConsole), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*url.Values)(nil), (*VirtualMachinePortForward)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_url_Values_To_v1alpha2_VirtualMachinePortForward(a.(*url.Values), b.(*VirtualMachinePortForward), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*url.Values)(nil), (*VirtualMachineVNC)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_url_Values_To_v1alpha2_VirtualMachineVNC(a.(*url.Values), b.(*VirtualMachineVNC), scope)
	}); err != nil {
		return err
	}
	return nil
}

func autoConvert_v1alpha2_VirtualMachineConsole_To_subresources_VirtualMachineConsole(in *VirtualMachineConsole, out *subresources.VirtualMachineConsole, s conversion.Scope) error {
	return nil
}

// Convert_v1alpha2_VirtualMachineConsole_To_subresources_VirtualMachineConsole is an autogenerated conversion function.
func Convert_v1alpha2_VirtualMachineConsole_To_subresources_VirtualMachineConsole(in *VirtualMachineConsole, out *subresources.VirtualMachineConsole, s conversion.Scope) error {
	return autoConvert_v1alpha2_VirtualMachineConsole_To_subresources_VirtualMachineConsole(in, out, s)
}

func autoConvert_subresources_VirtualMachineConsole_To_v1alpha2_VirtualMachineConsole(in *subresources.VirtualMachineConsole, out *VirtualMachineConsole, s conversion.Scope) error {
	return nil
}

// Convert_subresources_VirtualMachineConsole_To_v1alpha2_VirtualMachineConsole is an autogenerated conversion function.
func Convert_subresources_VirtualMachineConsole_To_v1alpha2_VirtualMachineConsole(in *subresources.VirtualMachineConsole, out *VirtualMachineConsole, s conversion.Scope) error {
	return autoConvert_subresources_VirtualMachineConsole_To_v1alpha2_VirtualMachineConsole(in, out, s)
}

func autoConvert_url_Values_To_v1alpha2_VirtualMachineConsole(in *url.Values, out *VirtualMachineConsole, s conversion.Scope) error {
	// WARNING: Field TypeMeta does not have json tag, skipping.

	return nil
}

// Convert_url_Values_To_v1alpha2_VirtualMachineConsole is an autogenerated conversion function.
func Convert_url_Values_To_v1alpha2_VirtualMachineConsole(in *url.Values, out *VirtualMachineConsole, s conversion.Scope) error {
	return autoConvert_url_Values_To_v1alpha2_VirtualMachineConsole(in, out, s)
}

func autoConvert_v1alpha2_VirtualMachinePortForward_To_subresources_VirtualMachinePortForward(in *VirtualMachinePortForward, out *subresources.VirtualMachinePortForward, s conversion.Scope) error {
	out.Protocol = in.Protocol
	out.Port = in.Port
	return nil
}

// Convert_v1alpha2_VirtualMachinePortForward_To_subresources_VirtualMachinePortForward is an autogenerated conversion function.
func Convert_v1alpha2_VirtualMachinePortForward_To_subresources_VirtualMachinePortForward(in *VirtualMachinePortForward, out *subresources.VirtualMachinePortForward, s conversion.Scope) error {
	return autoConvert_v1alpha2_VirtualMachinePortForward_To_subresources_VirtualMachinePortForward(in, out, s)
}

func autoConvert_subresources_VirtualMachinePortForward_To_v1alpha2_VirtualMachinePortForward(in *subresources.VirtualMachinePortForward, out *VirtualMachinePortForward, s conversion.Scope) error {
	out.Protocol = in.Protocol
	out.Port = in.Port
	return nil
}

// Convert_subresources_VirtualMachinePortForward_To_v1alpha2_VirtualMachinePortForward is an autogenerated conversion function.
func Convert_subresources_VirtualMachinePortForward_To_v1alpha2_VirtualMachinePortForward(in *subresources.VirtualMachinePortForward, out *VirtualMachinePortForward, s conversion.Scope) error {
	return autoConvert_subresources_VirtualMachinePortForward_To_v1alpha2_VirtualMachinePortForward(in, out, s)
}

func autoConvert_url_Values_To_v1alpha2_VirtualMachinePortForward(in *url.Values, out *VirtualMachinePortForward, s conversion.Scope) error {
	// WARNING: Field TypeMeta does not have json tag, skipping.

	if values, ok := map[string][]string(*in)["protocol"]; ok && len(values) > 0 {
		if err := runtime.Convert_Slice_string_To_string(&values, &out.Protocol, s); err != nil {
			return err
		}
	} else {
		out.Protocol = ""
	}
	if values, ok := map[string][]string(*in)["port"]; ok && len(values) > 0 {
		if err := runtime.Convert_Slice_string_To_int(&values, &out.Port, s); err != nil {
			return err
		}
	} else {
		out.Port = 0
	}
	return nil
}

// Convert_url_Values_To_v1alpha2_VirtualMachinePortForward is an autogenerated conversion function.
func Convert_url_Values_To_v1alpha2_VirtualMachinePortForward(in *url.Values, out *VirtualMachinePortForward, s conversion.Scope) error {
	return autoConvert_url_Values_To_v1alpha2_VirtualMachinePortForward(in, out, s)
}

func autoConvert_v1alpha2_VirtualMachineVNC_To_subresources_VirtualMachineVNC(in *VirtualMachineVNC, out *subresources.VirtualMachineVNC, s conversion.Scope) error {
	return nil
}

// Convert_v1alpha2_VirtualMachineVNC_To_subresources_VirtualMachineVNC is an autogenerated conversion function.
func Convert_v1alpha2_VirtualMachineVNC_To_subresources_VirtualMachineVNC(in *VirtualMachineVNC, out *subresources.VirtualMachineVNC, s conversion.Scope) error {
	return autoConvert_v1alpha2_VirtualMachineVNC_To_subresources_VirtualMachineVNC(in, out, s)
}

func autoConvert_subresources_VirtualMachineVNC_To_v1alpha2_VirtualMachineVNC(in *subresources.VirtualMachineVNC, out *VirtualMachineVNC, s conversion.Scope) error {
	return nil
}

// Convert_subresources_VirtualMachineVNC_To_v1alpha2_VirtualMachineVNC is an autogenerated conversion function.
func Convert_subresources_VirtualMachineVNC_To_v1alpha2_VirtualMachineVNC(in *subresources.VirtualMachineVNC, out *VirtualMachineVNC, s conversion.Scope) error {
	return autoConvert_subresources_VirtualMachineVNC_To_v1alpha2_VirtualMachineVNC(in, out, s)
}

func autoConvert_url_Values_To_v1alpha2_VirtualMachineVNC(in *url.Values, out *VirtualMachineVNC, s conversion.Scope) error {
	// WARNING: Field TypeMeta does not have json tag, skipping.

	return nil
}

// Convert_url_Values_To_v1alpha2_VirtualMachineVNC is an autogenerated conversion function.
func Convert_url_Values_To_v1alpha2_VirtualMachineVNC(in *url.Values, out *VirtualMachineVNC, s conversion.Scope) error {
	return autoConvert_url_Values_To_v1alpha2_VirtualMachineVNC(in, out, s)
}
