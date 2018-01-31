package v3

import (
	"github.com/rancher/norman/lifecycle"
	"k8s.io/apimachinery/pkg/runtime"
)

type RemoteAccountLifecycle interface {
	Create(obj *RemoteAccount) (*RemoteAccount, error)
	Remove(obj *RemoteAccount) (*RemoteAccount, error)
	Updated(obj *RemoteAccount) (*RemoteAccount, error)
}

type remoteAccountLifecycleAdapter struct {
	lifecycle RemoteAccountLifecycle
}

func (w *remoteAccountLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Create(obj.(*RemoteAccount))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *remoteAccountLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Remove(obj.(*RemoteAccount))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *remoteAccountLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Updated(obj.(*RemoteAccount))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewRemoteAccountLifecycleAdapter(name string, clusterScoped bool, client RemoteAccountInterface, l RemoteAccountLifecycle) RemoteAccountHandlerFunc {
	adapter := &remoteAccountLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *RemoteAccount) error {
		if obj == nil {
			return syncFn(key, nil)
		}
		return syncFn(key, obj)
	}
}
