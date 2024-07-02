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

package internal

import (
	"context"

	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/deckhouse/virtualization-controller/pkg/controller/vmip/internal/state"
	virtv2 "github.com/deckhouse/virtualization/api/core/v1alpha2"
)

type ProtectionHandler struct {
	client client.Client
	logger logr.Logger
}

func NewProtectionHandler(client client.Client, logger logr.Logger) *ProtectionHandler {
	return &ProtectionHandler{
		client: client,
		logger: logger.WithValues("handler", "ProtectionHandler"),
	}
}

func (h *ProtectionHandler) Handle(ctx context.Context, state state.VMIPState) (reconcile.Result, error) {
	vm, err := state.VirtualMachine(ctx)
	if err != nil {
		return reconcile.Result{}, err
	}

	shouldUnbound := vm == nil

	switch {
	case shouldUnbound:
		h.logger.Info("The VirtualMachineIP is no longer used by the VM: unbound", "name", state.VirtualMachineIP().Name())
		controllerutil.RemoveFinalizer(state.VirtualMachineIP().Changed(), virtv2.FinalizerIPAddressCleanup)

	case controllerutil.AddFinalizer(state.VirtualMachineIP().Changed(), virtv2.FinalizerIPAddressCleanup):
		return reconcile.Result{Requeue: true}, nil
	}

	return reconcile.Result{}, nil
}
