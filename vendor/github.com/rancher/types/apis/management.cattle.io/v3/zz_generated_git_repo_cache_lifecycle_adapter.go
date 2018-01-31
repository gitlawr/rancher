package v3

import (
	"github.com/rancher/norman/lifecycle"
	"k8s.io/apimachinery/pkg/runtime"
)

type GitRepoCacheLifecycle interface {
	Create(obj *GitRepoCache) (*GitRepoCache, error)
	Remove(obj *GitRepoCache) (*GitRepoCache, error)
	Updated(obj *GitRepoCache) (*GitRepoCache, error)
}

type gitRepoCacheLifecycleAdapter struct {
	lifecycle GitRepoCacheLifecycle
}

func (w *gitRepoCacheLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Create(obj.(*GitRepoCache))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *gitRepoCacheLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Remove(obj.(*GitRepoCache))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *gitRepoCacheLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Updated(obj.(*GitRepoCache))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewGitRepoCacheLifecycleAdapter(name string, clusterScoped bool, client GitRepoCacheInterface, l GitRepoCacheLifecycle) GitRepoCacheHandlerFunc {
	adapter := &gitRepoCacheLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *GitRepoCache) error {
		if obj == nil {
			return syncFn(key, nil)
		}
		return syncFn(key, obj)
	}
}
