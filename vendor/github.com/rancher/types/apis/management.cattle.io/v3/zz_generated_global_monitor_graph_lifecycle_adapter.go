package v3

import (
	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"k8s.io/apimachinery/pkg/runtime"
)

type GlobalMonitorGraphLifecycle interface {
	Create(obj *GlobalMonitorGraph) (runtime.Object, error)
	Remove(obj *GlobalMonitorGraph) (runtime.Object, error)
	Updated(obj *GlobalMonitorGraph) (runtime.Object, error)
}

type globalMonitorGraphLifecycleAdapter struct {
	lifecycle GlobalMonitorGraphLifecycle
}

func (w *globalMonitorGraphLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *globalMonitorGraphLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *globalMonitorGraphLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Create(obj.(*GlobalMonitorGraph))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *globalMonitorGraphLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Remove(obj.(*GlobalMonitorGraph))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *globalMonitorGraphLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Updated(obj.(*GlobalMonitorGraph))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewGlobalMonitorGraphLifecycleAdapter(name string, clusterScoped bool, client GlobalMonitorGraphInterface, l GlobalMonitorGraphLifecycle) GlobalMonitorGraphHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(GlobalMonitorGraphGroupVersionResource)
	}
	adapter := &globalMonitorGraphLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *GlobalMonitorGraph) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}
