// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package internal

import (
	"context"
	virtv2 "github.com/deckhouse/virtualization/api/core/v1alpha2"
	corev1 "k8s.io/api/core/v1"
	"sync"
)

// Ensure, that StorerMock does implement Storer.
// If this is not the case, regenerate this file with moq.
var _ Storer = &StorerMock{}

// StorerMock is a mock implementation of Storer.
//
//	func TestSomethingThatUsesStorer(t *testing.T) {
//
//		// make and configure a mocked Storer
//		mockedStorer := &StorerMock{
//			StoreFunc: func(ctx context.Context, vm *virtv2.VirtualMachine, vmSnapshot *virtv2.VirtualMachineSnapshot) (*corev1.Secret, error) {
//				panic("mock out the Store method")
//			},
//		}
//
//		// use mockedStorer in code that requires Storer
//		// and then make assertions.
//
//	}
type StorerMock struct {
	// StoreFunc mocks the Store method.
	StoreFunc func(ctx context.Context, vm *virtv2.VirtualMachine, vmSnapshot *virtv2.VirtualMachineSnapshot) (*corev1.Secret, error)

	// calls tracks calls to the methods.
	calls struct {
		// Store holds details about calls to the Store method.
		Store []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// VM is the vm argument value.
			VM *virtv2.VirtualMachine
			// VmSnapshot is the vmSnapshot argument value.
			VmSnapshot *virtv2.VirtualMachineSnapshot
		}
	}
	lockStore sync.RWMutex
}

// Store calls StoreFunc.
func (mock *StorerMock) Store(ctx context.Context, vm *virtv2.VirtualMachine, vmSnapshot *virtv2.VirtualMachineSnapshot) (*corev1.Secret, error) {
	if mock.StoreFunc == nil {
		panic("StorerMock.StoreFunc: method is nil but Storer.Store was just called")
	}
	callInfo := struct {
		Ctx        context.Context
		VM         *virtv2.VirtualMachine
		VmSnapshot *virtv2.VirtualMachineSnapshot
	}{
		Ctx:        ctx,
		VM:         vm,
		VmSnapshot: vmSnapshot,
	}
	mock.lockStore.Lock()
	mock.calls.Store = append(mock.calls.Store, callInfo)
	mock.lockStore.Unlock()
	return mock.StoreFunc(ctx, vm, vmSnapshot)
}

// StoreCalls gets all the calls that were made to Store.
// Check the length with:
//
//	len(mockedStorer.StoreCalls())
func (mock *StorerMock) StoreCalls() []struct {
	Ctx        context.Context
	VM         *virtv2.VirtualMachine
	VmSnapshot *virtv2.VirtualMachineSnapshot
} {
	var calls []struct {
		Ctx        context.Context
		VM         *virtv2.VirtualMachine
		VmSnapshot *virtv2.VirtualMachineSnapshot
	}
	mock.lockStore.RLock()
	calls = mock.calls.Store
	mock.lockStore.RUnlock()
	return calls
}

// Ensure, that SnapshotterMock does implement Snapshotter.
// If this is not the case, regenerate this file with moq.
var _ Snapshotter = &SnapshotterMock{}

// SnapshotterMock is a mock implementation of Snapshotter.
//
//	func TestSomethingThatUsesSnapshotter(t *testing.T) {
//
//		// make and configure a mocked Snapshotter
//		mockedSnapshotter := &SnapshotterMock{
//			CanFreezeFunc: func(vm *virtv2.VirtualMachine) bool {
//				panic("mock out the CanFreeze method")
//			},
//			CanUnfreezeFunc: func(ctx context.Context, vdSnapshotName string, vm *virtv2.VirtualMachine) (bool, error) {
//				panic("mock out the CanUnfreeze method")
//			},
//			CreateVirtualDiskSnapshotFunc: func(ctx context.Context, vdSnapshot *virtv2.VirtualDiskSnapshot) (*virtv2.VirtualDiskSnapshot, error) {
//				panic("mock out the CreateVirtualDiskSnapshot method")
//			},
//			FreezeFunc: func(ctx context.Context, name string, namespace string) error {
//				panic("mock out the Freeze method")
//			},
//			GetPersistentVolumeClaimFunc: func(ctx context.Context, name string, namespace string) (*corev1.PersistentVolumeClaim, error) {
//				panic("mock out the GetPersistentVolumeClaim method")
//			},
//			GetSecretFunc: func(ctx context.Context, name string, namespace string) (*corev1.Secret, error) {
//				panic("mock out the GetSecret method")
//			},
//			GetVirtualDiskFunc: func(ctx context.Context, name string, namespace string) (*virtv2.VirtualDisk, error) {
//				panic("mock out the GetVirtualDisk method")
//			},
//			GetVirtualDiskSnapshotFunc: func(ctx context.Context, name string, namespace string) (*virtv2.VirtualDiskSnapshot, error) {
//				panic("mock out the GetVirtualDiskSnapshot method")
//			},
//			GetVirtualMachineFunc: func(ctx context.Context, name string, namespace string) (*virtv2.VirtualMachine, error) {
//				panic("mock out the GetVirtualMachine method")
//			},
//			IsFrozenFunc: func(vm *virtv2.VirtualMachine) bool {
//				panic("mock out the IsFrozen method")
//			},
//			UnfreezeFunc: func(ctx context.Context, name string, namespace string) error {
//				panic("mock out the Unfreeze method")
//			},
//		}
//
//		// use mockedSnapshotter in code that requires Snapshotter
//		// and then make assertions.
//
//	}
type SnapshotterMock struct {
	// CanFreezeFunc mocks the CanFreeze method.
	CanFreezeFunc func(vm *virtv2.VirtualMachine) bool

	// CanUnfreezeFunc mocks the CanUnfreeze method.
	CanUnfreezeFunc func(ctx context.Context, vdSnapshotName string, vm *virtv2.VirtualMachine) (bool, error)

	// CreateVirtualDiskSnapshotFunc mocks the CreateVirtualDiskSnapshot method.
	CreateVirtualDiskSnapshotFunc func(ctx context.Context, vdSnapshot *virtv2.VirtualDiskSnapshot) (*virtv2.VirtualDiskSnapshot, error)

	// FreezeFunc mocks the Freeze method.
	FreezeFunc func(ctx context.Context, name string, namespace string) error

	// GetPersistentVolumeClaimFunc mocks the GetPersistentVolumeClaim method.
	GetPersistentVolumeClaimFunc func(ctx context.Context, name string, namespace string) (*corev1.PersistentVolumeClaim, error)

	// GetSecretFunc mocks the GetSecret method.
	GetSecretFunc func(ctx context.Context, name string, namespace string) (*corev1.Secret, error)

	// GetVirtualDiskFunc mocks the GetVirtualDisk method.
	GetVirtualDiskFunc func(ctx context.Context, name string, namespace string) (*virtv2.VirtualDisk, error)

	// GetVirtualDiskSnapshotFunc mocks the GetVirtualDiskSnapshot method.
	GetVirtualDiskSnapshotFunc func(ctx context.Context, name string, namespace string) (*virtv2.VirtualDiskSnapshot, error)

	// GetVirtualMachineFunc mocks the GetVirtualMachine method.
	GetVirtualMachineFunc func(ctx context.Context, name string, namespace string) (*virtv2.VirtualMachine, error)

	// IsFrozenFunc mocks the IsFrozen method.
	IsFrozenFunc func(vm *virtv2.VirtualMachine) bool

	// UnfreezeFunc mocks the Unfreeze method.
	UnfreezeFunc func(ctx context.Context, name string, namespace string) error

	// calls tracks calls to the methods.
	calls struct {
		// CanFreeze holds details about calls to the CanFreeze method.
		CanFreeze []struct {
			// VM is the vm argument value.
			VM *virtv2.VirtualMachine
		}
		// CanUnfreeze holds details about calls to the CanUnfreeze method.
		CanUnfreeze []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// VdSnapshotName is the vdSnapshotName argument value.
			VdSnapshotName string
			// VM is the vm argument value.
			VM *virtv2.VirtualMachine
		}
		// CreateVirtualDiskSnapshot holds details about calls to the CreateVirtualDiskSnapshot method.
		CreateVirtualDiskSnapshot []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// VdSnapshot is the vdSnapshot argument value.
			VdSnapshot *virtv2.VirtualDiskSnapshot
		}
		// Freeze holds details about calls to the Freeze method.
		Freeze []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Name is the name argument value.
			Name string
			// Namespace is the namespace argument value.
			Namespace string
		}
		// GetPersistentVolumeClaim holds details about calls to the GetPersistentVolumeClaim method.
		GetPersistentVolumeClaim []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Name is the name argument value.
			Name string
			// Namespace is the namespace argument value.
			Namespace string
		}
		// GetSecret holds details about calls to the GetSecret method.
		GetSecret []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Name is the name argument value.
			Name string
			// Namespace is the namespace argument value.
			Namespace string
		}
		// GetVirtualDisk holds details about calls to the GetVirtualDisk method.
		GetVirtualDisk []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Name is the name argument value.
			Name string
			// Namespace is the namespace argument value.
			Namespace string
		}
		// GetVirtualDiskSnapshot holds details about calls to the GetVirtualDiskSnapshot method.
		GetVirtualDiskSnapshot []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Name is the name argument value.
			Name string
			// Namespace is the namespace argument value.
			Namespace string
		}
		// GetVirtualMachine holds details about calls to the GetVirtualMachine method.
		GetVirtualMachine []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Name is the name argument value.
			Name string
			// Namespace is the namespace argument value.
			Namespace string
		}
		// IsFrozen holds details about calls to the IsFrozen method.
		IsFrozen []struct {
			// VM is the vm argument value.
			VM *virtv2.VirtualMachine
		}
		// Unfreeze holds details about calls to the Unfreeze method.
		Unfreeze []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Name is the name argument value.
			Name string
			// Namespace is the namespace argument value.
			Namespace string
		}
	}
	lockCanFreeze                 sync.RWMutex
	lockCanUnfreeze               sync.RWMutex
	lockCreateVirtualDiskSnapshot sync.RWMutex
	lockFreeze                    sync.RWMutex
	lockGetPersistentVolumeClaim  sync.RWMutex
	lockGetSecret                 sync.RWMutex
	lockGetVirtualDisk            sync.RWMutex
	lockGetVirtualDiskSnapshot    sync.RWMutex
	lockGetVirtualMachine         sync.RWMutex
	lockIsFrozen                  sync.RWMutex
	lockUnfreeze                  sync.RWMutex
}

// CanFreeze calls CanFreezeFunc.
func (mock *SnapshotterMock) CanFreeze(vm *virtv2.VirtualMachine) bool {
	if mock.CanFreezeFunc == nil {
		panic("SnapshotterMock.CanFreezeFunc: method is nil but Snapshotter.CanFreeze was just called")
	}
	callInfo := struct {
		VM *virtv2.VirtualMachine
	}{
		VM: vm,
	}
	mock.lockCanFreeze.Lock()
	mock.calls.CanFreeze = append(mock.calls.CanFreeze, callInfo)
	mock.lockCanFreeze.Unlock()
	return mock.CanFreezeFunc(vm)
}

// CanFreezeCalls gets all the calls that were made to CanFreeze.
// Check the length with:
//
//	len(mockedSnapshotter.CanFreezeCalls())
func (mock *SnapshotterMock) CanFreezeCalls() []struct {
	VM *virtv2.VirtualMachine
} {
	var calls []struct {
		VM *virtv2.VirtualMachine
	}
	mock.lockCanFreeze.RLock()
	calls = mock.calls.CanFreeze
	mock.lockCanFreeze.RUnlock()
	return calls
}

// CanUnfreeze calls CanUnfreezeFunc.
func (mock *SnapshotterMock) CanUnfreeze(ctx context.Context, vdSnapshotName string, vm *virtv2.VirtualMachine) (bool, error) {
	if mock.CanUnfreezeFunc == nil {
		panic("SnapshotterMock.CanUnfreezeFunc: method is nil but Snapshotter.CanUnfreeze was just called")
	}
	callInfo := struct {
		Ctx            context.Context
		VdSnapshotName string
		VM             *virtv2.VirtualMachine
	}{
		Ctx:            ctx,
		VdSnapshotName: vdSnapshotName,
		VM:             vm,
	}
	mock.lockCanUnfreeze.Lock()
	mock.calls.CanUnfreeze = append(mock.calls.CanUnfreeze, callInfo)
	mock.lockCanUnfreeze.Unlock()
	return mock.CanUnfreezeFunc(ctx, vdSnapshotName, vm)
}

// CanUnfreezeCalls gets all the calls that were made to CanUnfreeze.
// Check the length with:
//
//	len(mockedSnapshotter.CanUnfreezeCalls())
func (mock *SnapshotterMock) CanUnfreezeCalls() []struct {
	Ctx            context.Context
	VdSnapshotName string
	VM             *virtv2.VirtualMachine
} {
	var calls []struct {
		Ctx            context.Context
		VdSnapshotName string
		VM             *virtv2.VirtualMachine
	}
	mock.lockCanUnfreeze.RLock()
	calls = mock.calls.CanUnfreeze
	mock.lockCanUnfreeze.RUnlock()
	return calls
}

// CreateVirtualDiskSnapshot calls CreateVirtualDiskSnapshotFunc.
func (mock *SnapshotterMock) CreateVirtualDiskSnapshot(ctx context.Context, vdSnapshot *virtv2.VirtualDiskSnapshot) (*virtv2.VirtualDiskSnapshot, error) {
	if mock.CreateVirtualDiskSnapshotFunc == nil {
		panic("SnapshotterMock.CreateVirtualDiskSnapshotFunc: method is nil but Snapshotter.CreateVirtualDiskSnapshot was just called")
	}
	callInfo := struct {
		Ctx        context.Context
		VdSnapshot *virtv2.VirtualDiskSnapshot
	}{
		Ctx:        ctx,
		VdSnapshot: vdSnapshot,
	}
	mock.lockCreateVirtualDiskSnapshot.Lock()
	mock.calls.CreateVirtualDiskSnapshot = append(mock.calls.CreateVirtualDiskSnapshot, callInfo)
	mock.lockCreateVirtualDiskSnapshot.Unlock()
	return mock.CreateVirtualDiskSnapshotFunc(ctx, vdSnapshot)
}

// CreateVirtualDiskSnapshotCalls gets all the calls that were made to CreateVirtualDiskSnapshot.
// Check the length with:
//
//	len(mockedSnapshotter.CreateVirtualDiskSnapshotCalls())
func (mock *SnapshotterMock) CreateVirtualDiskSnapshotCalls() []struct {
	Ctx        context.Context
	VdSnapshot *virtv2.VirtualDiskSnapshot
} {
	var calls []struct {
		Ctx        context.Context
		VdSnapshot *virtv2.VirtualDiskSnapshot
	}
	mock.lockCreateVirtualDiskSnapshot.RLock()
	calls = mock.calls.CreateVirtualDiskSnapshot
	mock.lockCreateVirtualDiskSnapshot.RUnlock()
	return calls
}

// Freeze calls FreezeFunc.
func (mock *SnapshotterMock) Freeze(ctx context.Context, name string, namespace string) error {
	if mock.FreezeFunc == nil {
		panic("SnapshotterMock.FreezeFunc: method is nil but Snapshotter.Freeze was just called")
	}
	callInfo := struct {
		Ctx       context.Context
		Name      string
		Namespace string
	}{
		Ctx:       ctx,
		Name:      name,
		Namespace: namespace,
	}
	mock.lockFreeze.Lock()
	mock.calls.Freeze = append(mock.calls.Freeze, callInfo)
	mock.lockFreeze.Unlock()
	return mock.FreezeFunc(ctx, name, namespace)
}

// FreezeCalls gets all the calls that were made to Freeze.
// Check the length with:
//
//	len(mockedSnapshotter.FreezeCalls())
func (mock *SnapshotterMock) FreezeCalls() []struct {
	Ctx       context.Context
	Name      string
	Namespace string
} {
	var calls []struct {
		Ctx       context.Context
		Name      string
		Namespace string
	}
	mock.lockFreeze.RLock()
	calls = mock.calls.Freeze
	mock.lockFreeze.RUnlock()
	return calls
}

// GetPersistentVolumeClaim calls GetPersistentVolumeClaimFunc.
func (mock *SnapshotterMock) GetPersistentVolumeClaim(ctx context.Context, name string, namespace string) (*corev1.PersistentVolumeClaim, error) {
	if mock.GetPersistentVolumeClaimFunc == nil {
		panic("SnapshotterMock.GetPersistentVolumeClaimFunc: method is nil but Snapshotter.GetPersistentVolumeClaim was just called")
	}
	callInfo := struct {
		Ctx       context.Context
		Name      string
		Namespace string
	}{
		Ctx:       ctx,
		Name:      name,
		Namespace: namespace,
	}
	mock.lockGetPersistentVolumeClaim.Lock()
	mock.calls.GetPersistentVolumeClaim = append(mock.calls.GetPersistentVolumeClaim, callInfo)
	mock.lockGetPersistentVolumeClaim.Unlock()
	return mock.GetPersistentVolumeClaimFunc(ctx, name, namespace)
}

// GetPersistentVolumeClaimCalls gets all the calls that were made to GetPersistentVolumeClaim.
// Check the length with:
//
//	len(mockedSnapshotter.GetPersistentVolumeClaimCalls())
func (mock *SnapshotterMock) GetPersistentVolumeClaimCalls() []struct {
	Ctx       context.Context
	Name      string
	Namespace string
} {
	var calls []struct {
		Ctx       context.Context
		Name      string
		Namespace string
	}
	mock.lockGetPersistentVolumeClaim.RLock()
	calls = mock.calls.GetPersistentVolumeClaim
	mock.lockGetPersistentVolumeClaim.RUnlock()
	return calls
}

// GetSecret calls GetSecretFunc.
func (mock *SnapshotterMock) GetSecret(ctx context.Context, name string, namespace string) (*corev1.Secret, error) {
	if mock.GetSecretFunc == nil {
		panic("SnapshotterMock.GetSecretFunc: method is nil but Snapshotter.GetSecret was just called")
	}
	callInfo := struct {
		Ctx       context.Context
		Name      string
		Namespace string
	}{
		Ctx:       ctx,
		Name:      name,
		Namespace: namespace,
	}
	mock.lockGetSecret.Lock()
	mock.calls.GetSecret = append(mock.calls.GetSecret, callInfo)
	mock.lockGetSecret.Unlock()
	return mock.GetSecretFunc(ctx, name, namespace)
}

// GetSecretCalls gets all the calls that were made to GetSecret.
// Check the length with:
//
//	len(mockedSnapshotter.GetSecretCalls())
func (mock *SnapshotterMock) GetSecretCalls() []struct {
	Ctx       context.Context
	Name      string
	Namespace string
} {
	var calls []struct {
		Ctx       context.Context
		Name      string
		Namespace string
	}
	mock.lockGetSecret.RLock()
	calls = mock.calls.GetSecret
	mock.lockGetSecret.RUnlock()
	return calls
}

// GetVirtualDisk calls GetVirtualDiskFunc.
func (mock *SnapshotterMock) GetVirtualDisk(ctx context.Context, name string, namespace string) (*virtv2.VirtualDisk, error) {
	if mock.GetVirtualDiskFunc == nil {
		panic("SnapshotterMock.GetVirtualDiskFunc: method is nil but Snapshotter.GetVirtualDisk was just called")
	}
	callInfo := struct {
		Ctx       context.Context
		Name      string
		Namespace string
	}{
		Ctx:       ctx,
		Name:      name,
		Namespace: namespace,
	}
	mock.lockGetVirtualDisk.Lock()
	mock.calls.GetVirtualDisk = append(mock.calls.GetVirtualDisk, callInfo)
	mock.lockGetVirtualDisk.Unlock()
	return mock.GetVirtualDiskFunc(ctx, name, namespace)
}

// GetVirtualDiskCalls gets all the calls that were made to GetVirtualDisk.
// Check the length with:
//
//	len(mockedSnapshotter.GetVirtualDiskCalls())
func (mock *SnapshotterMock) GetVirtualDiskCalls() []struct {
	Ctx       context.Context
	Name      string
	Namespace string
} {
	var calls []struct {
		Ctx       context.Context
		Name      string
		Namespace string
	}
	mock.lockGetVirtualDisk.RLock()
	calls = mock.calls.GetVirtualDisk
	mock.lockGetVirtualDisk.RUnlock()
	return calls
}

// GetVirtualDiskSnapshot calls GetVirtualDiskSnapshotFunc.
func (mock *SnapshotterMock) GetVirtualDiskSnapshot(ctx context.Context, name string, namespace string) (*virtv2.VirtualDiskSnapshot, error) {
	if mock.GetVirtualDiskSnapshotFunc == nil {
		panic("SnapshotterMock.GetVirtualDiskSnapshotFunc: method is nil but Snapshotter.GetVirtualDiskSnapshot was just called")
	}
	callInfo := struct {
		Ctx       context.Context
		Name      string
		Namespace string
	}{
		Ctx:       ctx,
		Name:      name,
		Namespace: namespace,
	}
	mock.lockGetVirtualDiskSnapshot.Lock()
	mock.calls.GetVirtualDiskSnapshot = append(mock.calls.GetVirtualDiskSnapshot, callInfo)
	mock.lockGetVirtualDiskSnapshot.Unlock()
	return mock.GetVirtualDiskSnapshotFunc(ctx, name, namespace)
}

// GetVirtualDiskSnapshotCalls gets all the calls that were made to GetVirtualDiskSnapshot.
// Check the length with:
//
//	len(mockedSnapshotter.GetVirtualDiskSnapshotCalls())
func (mock *SnapshotterMock) GetVirtualDiskSnapshotCalls() []struct {
	Ctx       context.Context
	Name      string
	Namespace string
} {
	var calls []struct {
		Ctx       context.Context
		Name      string
		Namespace string
	}
	mock.lockGetVirtualDiskSnapshot.RLock()
	calls = mock.calls.GetVirtualDiskSnapshot
	mock.lockGetVirtualDiskSnapshot.RUnlock()
	return calls
}

// GetVirtualMachine calls GetVirtualMachineFunc.
func (mock *SnapshotterMock) GetVirtualMachine(ctx context.Context, name string, namespace string) (*virtv2.VirtualMachine, error) {
	if mock.GetVirtualMachineFunc == nil {
		panic("SnapshotterMock.GetVirtualMachineFunc: method is nil but Snapshotter.GetVirtualMachine was just called")
	}
	callInfo := struct {
		Ctx       context.Context
		Name      string
		Namespace string
	}{
		Ctx:       ctx,
		Name:      name,
		Namespace: namespace,
	}
	mock.lockGetVirtualMachine.Lock()
	mock.calls.GetVirtualMachine = append(mock.calls.GetVirtualMachine, callInfo)
	mock.lockGetVirtualMachine.Unlock()
	return mock.GetVirtualMachineFunc(ctx, name, namespace)
}

// GetVirtualMachineCalls gets all the calls that were made to GetVirtualMachine.
// Check the length with:
//
//	len(mockedSnapshotter.GetVirtualMachineCalls())
func (mock *SnapshotterMock) GetVirtualMachineCalls() []struct {
	Ctx       context.Context
	Name      string
	Namespace string
} {
	var calls []struct {
		Ctx       context.Context
		Name      string
		Namespace string
	}
	mock.lockGetVirtualMachine.RLock()
	calls = mock.calls.GetVirtualMachine
	mock.lockGetVirtualMachine.RUnlock()
	return calls
}

// IsFrozen calls IsFrozenFunc.
func (mock *SnapshotterMock) IsFrozen(vm *virtv2.VirtualMachine) bool {
	if mock.IsFrozenFunc == nil {
		panic("SnapshotterMock.IsFrozenFunc: method is nil but Snapshotter.IsFrozen was just called")
	}
	callInfo := struct {
		VM *virtv2.VirtualMachine
	}{
		VM: vm,
	}
	mock.lockIsFrozen.Lock()
	mock.calls.IsFrozen = append(mock.calls.IsFrozen, callInfo)
	mock.lockIsFrozen.Unlock()
	return mock.IsFrozenFunc(vm)
}

// IsFrozenCalls gets all the calls that were made to IsFrozen.
// Check the length with:
//
//	len(mockedSnapshotter.IsFrozenCalls())
func (mock *SnapshotterMock) IsFrozenCalls() []struct {
	VM *virtv2.VirtualMachine
} {
	var calls []struct {
		VM *virtv2.VirtualMachine
	}
	mock.lockIsFrozen.RLock()
	calls = mock.calls.IsFrozen
	mock.lockIsFrozen.RUnlock()
	return calls
}

// Unfreeze calls UnfreezeFunc.
func (mock *SnapshotterMock) Unfreeze(ctx context.Context, name string, namespace string) error {
	if mock.UnfreezeFunc == nil {
		panic("SnapshotterMock.UnfreezeFunc: method is nil but Snapshotter.Unfreeze was just called")
	}
	callInfo := struct {
		Ctx       context.Context
		Name      string
		Namespace string
	}{
		Ctx:       ctx,
		Name:      name,
		Namespace: namespace,
	}
	mock.lockUnfreeze.Lock()
	mock.calls.Unfreeze = append(mock.calls.Unfreeze, callInfo)
	mock.lockUnfreeze.Unlock()
	return mock.UnfreezeFunc(ctx, name, namespace)
}

// UnfreezeCalls gets all the calls that were made to Unfreeze.
// Check the length with:
//
//	len(mockedSnapshotter.UnfreezeCalls())
func (mock *SnapshotterMock) UnfreezeCalls() []struct {
	Ctx       context.Context
	Name      string
	Namespace string
} {
	var calls []struct {
		Ctx       context.Context
		Name      string
		Namespace string
	}
	mock.lockUnfreeze.RLock()
	calls = mock.calls.Unfreeze
	mock.lockUnfreeze.RUnlock()
	return calls
}