// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package service

import (
	"context"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sync"
)

// Ensure, that ClientMock does implement Client.
// If this is not the case, regenerate this file with moq.
var _ Client = &ClientMock{}

// ClientMock is a mock implementation of Client.
//
//	func TestSomethingThatUsesClient(t *testing.T) {
//
//		// make and configure a mocked Client
//		mockedClient := &ClientMock{
//			CreateFunc: func(ctx context.Context, obj client.Object, opts ...client.CreateOption) error {
//				panic("mock out the Create method")
//			},
//			DeleteFunc: func(ctx context.Context, obj client.Object, opts ...client.DeleteOption) error {
//				panic("mock out the Delete method")
//			},
//			DeleteAllOfFunc: func(ctx context.Context, obj client.Object, opts ...client.DeleteAllOfOption) error {
//				panic("mock out the DeleteAllOf method")
//			},
//			GetFunc: func(ctx context.Context, key types.NamespacedName, obj client.Object, opts ...client.GetOption) error {
//				panic("mock out the Get method")
//			},
//			GroupVersionKindForFunc: func(obj runtime.Object) (schema.GroupVersionKind, error) {
//				panic("mock out the GroupVersionKindFor method")
//			},
//			IsObjectNamespacedFunc: func(obj runtime.Object) (bool, error) {
//				panic("mock out the IsObjectNamespaced method")
//			},
//			ListFunc: func(ctx context.Context, list client.ObjectList, opts ...client.ListOption) error {
//				panic("mock out the List method")
//			},
//			PatchFunc: func(ctx context.Context, obj client.Object, patch client.Patch, opts ...client.PatchOption) error {
//				panic("mock out the Patch method")
//			},
//			RESTMapperFunc: func() meta.RESTMapper {
//				panic("mock out the RESTMapper method")
//			},
//			SchemeFunc: func() *runtime.Scheme {
//				panic("mock out the Scheme method")
//			},
//			StatusFunc: func() client.SubResourceWriter {
//				panic("mock out the Status method")
//			},
//			SubResourceFunc: func(subResource string) client.SubResourceClient {
//				panic("mock out the SubResource method")
//			},
//			UpdateFunc: func(ctx context.Context, obj client.Object, opts ...client.UpdateOption) error {
//				panic("mock out the Update method")
//			},
//		}
//
//		// use mockedClient in code that requires Client
//		// and then make assertions.
//
//	}
type ClientMock struct {
	// CreateFunc mocks the Create method.
	CreateFunc func(ctx context.Context, obj client.Object, opts ...client.CreateOption) error

	// DeleteFunc mocks the Delete method.
	DeleteFunc func(ctx context.Context, obj client.Object, opts ...client.DeleteOption) error

	// DeleteAllOfFunc mocks the DeleteAllOf method.
	DeleteAllOfFunc func(ctx context.Context, obj client.Object, opts ...client.DeleteAllOfOption) error

	// GetFunc mocks the Get method.
	GetFunc func(ctx context.Context, key types.NamespacedName, obj client.Object, opts ...client.GetOption) error

	// GroupVersionKindForFunc mocks the GroupVersionKindFor method.
	GroupVersionKindForFunc func(obj runtime.Object) (schema.GroupVersionKind, error)

	// IsObjectNamespacedFunc mocks the IsObjectNamespaced method.
	IsObjectNamespacedFunc func(obj runtime.Object) (bool, error)

	// ListFunc mocks the List method.
	ListFunc func(ctx context.Context, list client.ObjectList, opts ...client.ListOption) error

	// PatchFunc mocks the Patch method.
	PatchFunc func(ctx context.Context, obj client.Object, patch client.Patch, opts ...client.PatchOption) error

	// RESTMapperFunc mocks the RESTMapper method.
	RESTMapperFunc func() meta.RESTMapper

	// SchemeFunc mocks the Scheme method.
	SchemeFunc func() *runtime.Scheme

	// StatusFunc mocks the Status method.
	StatusFunc func() client.SubResourceWriter

	// SubResourceFunc mocks the SubResource method.
	SubResourceFunc func(subResource string) client.SubResourceClient

	// UpdateFunc mocks the Update method.
	UpdateFunc func(ctx context.Context, obj client.Object, opts ...client.UpdateOption) error

	// calls tracks calls to the methods.
	calls struct {
		// Create holds details about calls to the Create method.
		Create []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Obj is the obj argument value.
			Obj client.Object
			// Opts is the opts argument value.
			Opts []client.CreateOption
		}
		// Delete holds details about calls to the Delete method.
		Delete []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Obj is the obj argument value.
			Obj client.Object
			// Opts is the opts argument value.
			Opts []client.DeleteOption
		}
		// DeleteAllOf holds details about calls to the DeleteAllOf method.
		DeleteAllOf []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Obj is the obj argument value.
			Obj client.Object
			// Opts is the opts argument value.
			Opts []client.DeleteAllOfOption
		}
		// Get holds details about calls to the Get method.
		Get []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Key is the key argument value.
			Key types.NamespacedName
			// Obj is the obj argument value.
			Obj client.Object
			// Opts is the opts argument value.
			Opts []client.GetOption
		}
		// GroupVersionKindFor holds details about calls to the GroupVersionKindFor method.
		GroupVersionKindFor []struct {
			// Obj is the obj argument value.
			Obj runtime.Object
		}
		// IsObjectNamespaced holds details about calls to the IsObjectNamespaced method.
		IsObjectNamespaced []struct {
			// Obj is the obj argument value.
			Obj runtime.Object
		}
		// List holds details about calls to the List method.
		List []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// List is the list argument value.
			List client.ObjectList
			// Opts is the opts argument value.
			Opts []client.ListOption
		}
		// Patch holds details about calls to the Patch method.
		Patch []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Obj is the obj argument value.
			Obj client.Object
			// Patch is the patch argument value.
			Patch client.Patch
			// Opts is the opts argument value.
			Opts []client.PatchOption
		}
		// RESTMapper holds details about calls to the RESTMapper method.
		RESTMapper []struct {
		}
		// Scheme holds details about calls to the Scheme method.
		Scheme []struct {
		}
		// Status holds details about calls to the Status method.
		Status []struct {
		}
		// SubResource holds details about calls to the SubResource method.
		SubResource []struct {
			// SubResource is the subResource argument value.
			SubResource string
		}
		// Update holds details about calls to the Update method.
		Update []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Obj is the obj argument value.
			Obj client.Object
			// Opts is the opts argument value.
			Opts []client.UpdateOption
		}
	}
	lockCreate              sync.RWMutex
	lockDelete              sync.RWMutex
	lockDeleteAllOf         sync.RWMutex
	lockGet                 sync.RWMutex
	lockGroupVersionKindFor sync.RWMutex
	lockIsObjectNamespaced  sync.RWMutex
	lockList                sync.RWMutex
	lockPatch               sync.RWMutex
	lockRESTMapper          sync.RWMutex
	lockScheme              sync.RWMutex
	lockStatus              sync.RWMutex
	lockSubResource         sync.RWMutex
	lockUpdate              sync.RWMutex
}

// Create calls CreateFunc.
func (mock *ClientMock) Create(ctx context.Context, obj client.Object, opts ...client.CreateOption) error {
	if mock.CreateFunc == nil {
		panic("ClientMock.CreateFunc: method is nil but Client.Create was just called")
	}
	callInfo := struct {
		Ctx  context.Context
		Obj  client.Object
		Opts []client.CreateOption
	}{
		Ctx:  ctx,
		Obj:  obj,
		Opts: opts,
	}
	mock.lockCreate.Lock()
	mock.calls.Create = append(mock.calls.Create, callInfo)
	mock.lockCreate.Unlock()
	return mock.CreateFunc(ctx, obj, opts...)
}

// CreateCalls gets all the calls that were made to Create.
// Check the length with:
//
//	len(mockedClient.CreateCalls())
func (mock *ClientMock) CreateCalls() []struct {
	Ctx  context.Context
	Obj  client.Object
	Opts []client.CreateOption
} {
	var calls []struct {
		Ctx  context.Context
		Obj  client.Object
		Opts []client.CreateOption
	}
	mock.lockCreate.RLock()
	calls = mock.calls.Create
	mock.lockCreate.RUnlock()
	return calls
}

// Delete calls DeleteFunc.
func (mock *ClientMock) Delete(ctx context.Context, obj client.Object, opts ...client.DeleteOption) error {
	if mock.DeleteFunc == nil {
		panic("ClientMock.DeleteFunc: method is nil but Client.Delete was just called")
	}
	callInfo := struct {
		Ctx  context.Context
		Obj  client.Object
		Opts []client.DeleteOption
	}{
		Ctx:  ctx,
		Obj:  obj,
		Opts: opts,
	}
	mock.lockDelete.Lock()
	mock.calls.Delete = append(mock.calls.Delete, callInfo)
	mock.lockDelete.Unlock()
	return mock.DeleteFunc(ctx, obj, opts...)
}

// DeleteCalls gets all the calls that were made to Delete.
// Check the length with:
//
//	len(mockedClient.DeleteCalls())
func (mock *ClientMock) DeleteCalls() []struct {
	Ctx  context.Context
	Obj  client.Object
	Opts []client.DeleteOption
} {
	var calls []struct {
		Ctx  context.Context
		Obj  client.Object
		Opts []client.DeleteOption
	}
	mock.lockDelete.RLock()
	calls = mock.calls.Delete
	mock.lockDelete.RUnlock()
	return calls
}

// DeleteAllOf calls DeleteAllOfFunc.
func (mock *ClientMock) DeleteAllOf(ctx context.Context, obj client.Object, opts ...client.DeleteAllOfOption) error {
	if mock.DeleteAllOfFunc == nil {
		panic("ClientMock.DeleteAllOfFunc: method is nil but Client.DeleteAllOf was just called")
	}
	callInfo := struct {
		Ctx  context.Context
		Obj  client.Object
		Opts []client.DeleteAllOfOption
	}{
		Ctx:  ctx,
		Obj:  obj,
		Opts: opts,
	}
	mock.lockDeleteAllOf.Lock()
	mock.calls.DeleteAllOf = append(mock.calls.DeleteAllOf, callInfo)
	mock.lockDeleteAllOf.Unlock()
	return mock.DeleteAllOfFunc(ctx, obj, opts...)
}

// DeleteAllOfCalls gets all the calls that were made to DeleteAllOf.
// Check the length with:
//
//	len(mockedClient.DeleteAllOfCalls())
func (mock *ClientMock) DeleteAllOfCalls() []struct {
	Ctx  context.Context
	Obj  client.Object
	Opts []client.DeleteAllOfOption
} {
	var calls []struct {
		Ctx  context.Context
		Obj  client.Object
		Opts []client.DeleteAllOfOption
	}
	mock.lockDeleteAllOf.RLock()
	calls = mock.calls.DeleteAllOf
	mock.lockDeleteAllOf.RUnlock()
	return calls
}

// Get calls GetFunc.
func (mock *ClientMock) Get(ctx context.Context, key types.NamespacedName, obj client.Object, opts ...client.GetOption) error {
	if mock.GetFunc == nil {
		panic("ClientMock.GetFunc: method is nil but Client.Get was just called")
	}
	callInfo := struct {
		Ctx  context.Context
		Key  types.NamespacedName
		Obj  client.Object
		Opts []client.GetOption
	}{
		Ctx:  ctx,
		Key:  key,
		Obj:  obj,
		Opts: opts,
	}
	mock.lockGet.Lock()
	mock.calls.Get = append(mock.calls.Get, callInfo)
	mock.lockGet.Unlock()
	return mock.GetFunc(ctx, key, obj, opts...)
}

// GetCalls gets all the calls that were made to Get.
// Check the length with:
//
//	len(mockedClient.GetCalls())
func (mock *ClientMock) GetCalls() []struct {
	Ctx  context.Context
	Key  types.NamespacedName
	Obj  client.Object
	Opts []client.GetOption
} {
	var calls []struct {
		Ctx  context.Context
		Key  types.NamespacedName
		Obj  client.Object
		Opts []client.GetOption
	}
	mock.lockGet.RLock()
	calls = mock.calls.Get
	mock.lockGet.RUnlock()
	return calls
}

// GroupVersionKindFor calls GroupVersionKindForFunc.
func (mock *ClientMock) GroupVersionKindFor(obj runtime.Object) (schema.GroupVersionKind, error) {
	if mock.GroupVersionKindForFunc == nil {
		panic("ClientMock.GroupVersionKindForFunc: method is nil but Client.GroupVersionKindFor was just called")
	}
	callInfo := struct {
		Obj runtime.Object
	}{
		Obj: obj,
	}
	mock.lockGroupVersionKindFor.Lock()
	mock.calls.GroupVersionKindFor = append(mock.calls.GroupVersionKindFor, callInfo)
	mock.lockGroupVersionKindFor.Unlock()
	return mock.GroupVersionKindForFunc(obj)
}

// GroupVersionKindForCalls gets all the calls that were made to GroupVersionKindFor.
// Check the length with:
//
//	len(mockedClient.GroupVersionKindForCalls())
func (mock *ClientMock) GroupVersionKindForCalls() []struct {
	Obj runtime.Object
} {
	var calls []struct {
		Obj runtime.Object
	}
	mock.lockGroupVersionKindFor.RLock()
	calls = mock.calls.GroupVersionKindFor
	mock.lockGroupVersionKindFor.RUnlock()
	return calls
}

// IsObjectNamespaced calls IsObjectNamespacedFunc.
func (mock *ClientMock) IsObjectNamespaced(obj runtime.Object) (bool, error) {
	if mock.IsObjectNamespacedFunc == nil {
		panic("ClientMock.IsObjectNamespacedFunc: method is nil but Client.IsObjectNamespaced was just called")
	}
	callInfo := struct {
		Obj runtime.Object
	}{
		Obj: obj,
	}
	mock.lockIsObjectNamespaced.Lock()
	mock.calls.IsObjectNamespaced = append(mock.calls.IsObjectNamespaced, callInfo)
	mock.lockIsObjectNamespaced.Unlock()
	return mock.IsObjectNamespacedFunc(obj)
}

// IsObjectNamespacedCalls gets all the calls that were made to IsObjectNamespaced.
// Check the length with:
//
//	len(mockedClient.IsObjectNamespacedCalls())
func (mock *ClientMock) IsObjectNamespacedCalls() []struct {
	Obj runtime.Object
} {
	var calls []struct {
		Obj runtime.Object
	}
	mock.lockIsObjectNamespaced.RLock()
	calls = mock.calls.IsObjectNamespaced
	mock.lockIsObjectNamespaced.RUnlock()
	return calls
}

// List calls ListFunc.
func (mock *ClientMock) List(ctx context.Context, list client.ObjectList, opts ...client.ListOption) error {
	if mock.ListFunc == nil {
		panic("ClientMock.ListFunc: method is nil but Client.List was just called")
	}
	callInfo := struct {
		Ctx  context.Context
		List client.ObjectList
		Opts []client.ListOption
	}{
		Ctx:  ctx,
		List: list,
		Opts: opts,
	}
	mock.lockList.Lock()
	mock.calls.List = append(mock.calls.List, callInfo)
	mock.lockList.Unlock()
	return mock.ListFunc(ctx, list, opts...)
}

// ListCalls gets all the calls that were made to List.
// Check the length with:
//
//	len(mockedClient.ListCalls())
func (mock *ClientMock) ListCalls() []struct {
	Ctx  context.Context
	List client.ObjectList
	Opts []client.ListOption
} {
	var calls []struct {
		Ctx  context.Context
		List client.ObjectList
		Opts []client.ListOption
	}
	mock.lockList.RLock()
	calls = mock.calls.List
	mock.lockList.RUnlock()
	return calls
}

// Patch calls PatchFunc.
func (mock *ClientMock) Patch(ctx context.Context, obj client.Object, patch client.Patch, opts ...client.PatchOption) error {
	if mock.PatchFunc == nil {
		panic("ClientMock.PatchFunc: method is nil but Client.Patch was just called")
	}
	callInfo := struct {
		Ctx   context.Context
		Obj   client.Object
		Patch client.Patch
		Opts  []client.PatchOption
	}{
		Ctx:   ctx,
		Obj:   obj,
		Patch: patch,
		Opts:  opts,
	}
	mock.lockPatch.Lock()
	mock.calls.Patch = append(mock.calls.Patch, callInfo)
	mock.lockPatch.Unlock()
	return mock.PatchFunc(ctx, obj, patch, opts...)
}

// PatchCalls gets all the calls that were made to Patch.
// Check the length with:
//
//	len(mockedClient.PatchCalls())
func (mock *ClientMock) PatchCalls() []struct {
	Ctx   context.Context
	Obj   client.Object
	Patch client.Patch
	Opts  []client.PatchOption
} {
	var calls []struct {
		Ctx   context.Context
		Obj   client.Object
		Patch client.Patch
		Opts  []client.PatchOption
	}
	mock.lockPatch.RLock()
	calls = mock.calls.Patch
	mock.lockPatch.RUnlock()
	return calls
}

// RESTMapper calls RESTMapperFunc.
func (mock *ClientMock) RESTMapper() meta.RESTMapper {
	if mock.RESTMapperFunc == nil {
		panic("ClientMock.RESTMapperFunc: method is nil but Client.RESTMapper was just called")
	}
	callInfo := struct {
	}{}
	mock.lockRESTMapper.Lock()
	mock.calls.RESTMapper = append(mock.calls.RESTMapper, callInfo)
	mock.lockRESTMapper.Unlock()
	return mock.RESTMapperFunc()
}

// RESTMapperCalls gets all the calls that were made to RESTMapper.
// Check the length with:
//
//	len(mockedClient.RESTMapperCalls())
func (mock *ClientMock) RESTMapperCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockRESTMapper.RLock()
	calls = mock.calls.RESTMapper
	mock.lockRESTMapper.RUnlock()
	return calls
}

// Scheme calls SchemeFunc.
func (mock *ClientMock) Scheme() *runtime.Scheme {
	if mock.SchemeFunc == nil {
		panic("ClientMock.SchemeFunc: method is nil but Client.Scheme was just called")
	}
	callInfo := struct {
	}{}
	mock.lockScheme.Lock()
	mock.calls.Scheme = append(mock.calls.Scheme, callInfo)
	mock.lockScheme.Unlock()
	return mock.SchemeFunc()
}

// SchemeCalls gets all the calls that were made to Scheme.
// Check the length with:
//
//	len(mockedClient.SchemeCalls())
func (mock *ClientMock) SchemeCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockScheme.RLock()
	calls = mock.calls.Scheme
	mock.lockScheme.RUnlock()
	return calls
}

// Status calls StatusFunc.
func (mock *ClientMock) Status() client.SubResourceWriter {
	if mock.StatusFunc == nil {
		panic("ClientMock.StatusFunc: method is nil but Client.Status was just called")
	}
	callInfo := struct {
	}{}
	mock.lockStatus.Lock()
	mock.calls.Status = append(mock.calls.Status, callInfo)
	mock.lockStatus.Unlock()
	return mock.StatusFunc()
}

// StatusCalls gets all the calls that were made to Status.
// Check the length with:
//
//	len(mockedClient.StatusCalls())
func (mock *ClientMock) StatusCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockStatus.RLock()
	calls = mock.calls.Status
	mock.lockStatus.RUnlock()
	return calls
}

// SubResource calls SubResourceFunc.
func (mock *ClientMock) SubResource(subResource string) client.SubResourceClient {
	if mock.SubResourceFunc == nil {
		panic("ClientMock.SubResourceFunc: method is nil but Client.SubResource was just called")
	}
	callInfo := struct {
		SubResource string
	}{
		SubResource: subResource,
	}
	mock.lockSubResource.Lock()
	mock.calls.SubResource = append(mock.calls.SubResource, callInfo)
	mock.lockSubResource.Unlock()
	return mock.SubResourceFunc(subResource)
}

// SubResourceCalls gets all the calls that were made to SubResource.
// Check the length with:
//
//	len(mockedClient.SubResourceCalls())
func (mock *ClientMock) SubResourceCalls() []struct {
	SubResource string
} {
	var calls []struct {
		SubResource string
	}
	mock.lockSubResource.RLock()
	calls = mock.calls.SubResource
	mock.lockSubResource.RUnlock()
	return calls
}

// Update calls UpdateFunc.
func (mock *ClientMock) Update(ctx context.Context, obj client.Object, opts ...client.UpdateOption) error {
	if mock.UpdateFunc == nil {
		panic("ClientMock.UpdateFunc: method is nil but Client.Update was just called")
	}
	callInfo := struct {
		Ctx  context.Context
		Obj  client.Object
		Opts []client.UpdateOption
	}{
		Ctx:  ctx,
		Obj:  obj,
		Opts: opts,
	}
	mock.lockUpdate.Lock()
	mock.calls.Update = append(mock.calls.Update, callInfo)
	mock.lockUpdate.Unlock()
	return mock.UpdateFunc(ctx, obj, opts...)
}

// UpdateCalls gets all the calls that were made to Update.
// Check the length with:
//
//	len(mockedClient.UpdateCalls())
func (mock *ClientMock) UpdateCalls() []struct {
	Ctx  context.Context
	Obj  client.Object
	Opts []client.UpdateOption
} {
	var calls []struct {
		Ctx  context.Context
		Obj  client.Object
		Opts []client.UpdateOption
	}
	mock.lockUpdate.RLock()
	calls = mock.calls.Update
	mock.lockUpdate.RUnlock()
	return calls
}