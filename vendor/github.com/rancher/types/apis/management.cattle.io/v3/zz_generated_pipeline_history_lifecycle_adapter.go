package v3

import (
	"github.com/rancher/norman/lifecycle"
	"k8s.io/apimachinery/pkg/runtime"
)

type PipelineHistoryLifecycle interface {
	Create(obj *PipelineHistory) (*PipelineHistory, error)
	Remove(obj *PipelineHistory) (*PipelineHistory, error)
	Updated(obj *PipelineHistory) (*PipelineHistory, error)
}

type pipelineHistoryLifecycleAdapter struct {
	lifecycle PipelineHistoryLifecycle
}

func (w *pipelineHistoryLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Create(obj.(*PipelineHistory))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *pipelineHistoryLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Remove(obj.(*PipelineHistory))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *pipelineHistoryLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Updated(obj.(*PipelineHistory))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewPipelineHistoryLifecycleAdapter(name string, clusterScoped bool, client PipelineHistoryInterface, l PipelineHistoryLifecycle) PipelineHistoryHandlerFunc {
	adapter := &pipelineHistoryLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *PipelineHistory) error {
		if obj == nil {
			return syncFn(key, nil)
		}
		return syncFn(key, obj)
	}
}
