package v3

import (
	"github.com/rancher/norman/lifecycle"
	"k8s.io/apimachinery/pkg/runtime"
)

type KontainerDriverLifecycle interface {
	Create(obj *KontainerDriver) (*KontainerDriver, error)
	Remove(obj *KontainerDriver) (*KontainerDriver, error)
	Updated(obj *KontainerDriver) (*KontainerDriver, error)
}

type kontainerDriverLifecycleAdapter struct {
	lifecycle KontainerDriverLifecycle
}

func (w *kontainerDriverLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Create(obj.(*KontainerDriver))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *kontainerDriverLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Remove(obj.(*KontainerDriver))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *kontainerDriverLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Updated(obj.(*KontainerDriver))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewKontainerDriverLifecycleAdapter(name string, clusterScoped bool, client KontainerDriverInterface, l KontainerDriverLifecycle) KontainerDriverHandlerFunc {
	adapter := &kontainerDriverLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *KontainerDriver) error {
		if obj == nil {
			return syncFn(key, nil)
		}
		return syncFn(key, obj)
	}
}
