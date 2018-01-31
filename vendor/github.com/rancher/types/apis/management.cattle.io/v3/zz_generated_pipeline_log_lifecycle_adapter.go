package v3

import (
	"github.com/rancher/norman/lifecycle"
	"k8s.io/apimachinery/pkg/runtime"
)

type PipelineLogLifecycle interface {
	Create(obj *PipelineLog) (*PipelineLog, error)
	Remove(obj *PipelineLog) (*PipelineLog, error)
	Updated(obj *PipelineLog) (*PipelineLog, error)
}

type pipelineLogLifecycleAdapter struct {
	lifecycle PipelineLogLifecycle
}

func (w *pipelineLogLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Create(obj.(*PipelineLog))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *pipelineLogLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Remove(obj.(*PipelineLog))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *pipelineLogLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Updated(obj.(*PipelineLog))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewPipelineLogLifecycleAdapter(name string, clusterScoped bool, client PipelineLogInterface, l PipelineLogLifecycle) PipelineLogHandlerFunc {
	adapter := &pipelineLogLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *PipelineLog) error {
		if obj == nil {
			return syncFn(key, nil)
		}
		return syncFn(key, obj)
	}
}
