package v3

import (
	"github.com/rancher/norman/lifecycle"
	"k8s.io/apimachinery/pkg/runtime"
)

type BitbucketCloudApplyInputLifecycle interface {
	Create(obj *BitbucketCloudApplyInput) (*BitbucketCloudApplyInput, error)
	Remove(obj *BitbucketCloudApplyInput) (*BitbucketCloudApplyInput, error)
	Updated(obj *BitbucketCloudApplyInput) (*BitbucketCloudApplyInput, error)
}

type bitbucketCloudApplyInputLifecycleAdapter struct {
	lifecycle BitbucketCloudApplyInputLifecycle
}

func (w *bitbucketCloudApplyInputLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Create(obj.(*BitbucketCloudApplyInput))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *bitbucketCloudApplyInputLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Remove(obj.(*BitbucketCloudApplyInput))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *bitbucketCloudApplyInputLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Updated(obj.(*BitbucketCloudApplyInput))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewBitbucketCloudApplyInputLifecycleAdapter(name string, clusterScoped bool, client BitbucketCloudApplyInputInterface, l BitbucketCloudApplyInputLifecycle) BitbucketCloudApplyInputHandlerFunc {
	adapter := &bitbucketCloudApplyInputLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *BitbucketCloudApplyInput) error {
		if obj == nil {
			return syncFn(key, nil)
		}
		return syncFn(key, obj)
	}
}
