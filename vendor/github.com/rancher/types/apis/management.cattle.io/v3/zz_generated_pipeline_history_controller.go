package v3

import (
	"context"

	"github.com/rancher/norman/clientbase"
	"github.com/rancher/norman/controller"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

var (
	PipelineHistoryGroupVersionKind = schema.GroupVersionKind{
		Version: Version,
		Group:   GroupName,
		Kind:    "PipelineHistory",
	}
	PipelineHistoryResource = metav1.APIResource{
		Name:         "pipelinehistories",
		SingularName: "pipelinehistory",
		Namespaced:   true,

		Kind: PipelineHistoryGroupVersionKind.Kind,
	}
)

type PipelineHistoryList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PipelineHistory
}

type PipelineHistoryHandlerFunc func(key string, obj *PipelineHistory) error

type PipelineHistoryLister interface {
	List(namespace string, selector labels.Selector) (ret []*PipelineHistory, err error)
	Get(namespace, name string) (*PipelineHistory, error)
}

type PipelineHistoryController interface {
	Informer() cache.SharedIndexInformer
	Lister() PipelineHistoryLister
	AddHandler(name string, handler PipelineHistoryHandlerFunc)
	AddClusterScopedHandler(name, clusterName string, handler PipelineHistoryHandlerFunc)
	Enqueue(namespace, name string)
	Sync(ctx context.Context) error
	Start(ctx context.Context, threadiness int) error
}

type PipelineHistoryInterface interface {
	ObjectClient() *clientbase.ObjectClient
	Create(*PipelineHistory) (*PipelineHistory, error)
	GetNamespaced(namespace, name string, opts metav1.GetOptions) (*PipelineHistory, error)
	Get(name string, opts metav1.GetOptions) (*PipelineHistory, error)
	Update(*PipelineHistory) (*PipelineHistory, error)
	Delete(name string, options *metav1.DeleteOptions) error
	DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error
	List(opts metav1.ListOptions) (*PipelineHistoryList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Controller() PipelineHistoryController
	AddHandler(name string, sync PipelineHistoryHandlerFunc)
	AddLifecycle(name string, lifecycle PipelineHistoryLifecycle)
	AddClusterScopedHandler(name, clusterName string, sync PipelineHistoryHandlerFunc)
	AddClusterScopedLifecycle(name, clusterName string, lifecycle PipelineHistoryLifecycle)
}

type pipelineHistoryLister struct {
	controller *pipelineHistoryController
}

func (l *pipelineHistoryLister) List(namespace string, selector labels.Selector) (ret []*PipelineHistory, err error) {
	err = cache.ListAllByNamespace(l.controller.Informer().GetIndexer(), namespace, selector, func(obj interface{}) {
		ret = append(ret, obj.(*PipelineHistory))
	})
	return
}

func (l *pipelineHistoryLister) Get(namespace, name string) (*PipelineHistory, error) {
	var key string
	if namespace != "" {
		key = namespace + "/" + name
	} else {
		key = name
	}
	obj, exists, err := l.controller.Informer().GetIndexer().GetByKey(key)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(schema.GroupResource{
			Group:    PipelineHistoryGroupVersionKind.Group,
			Resource: "pipelineHistory",
		}, name)
	}
	return obj.(*PipelineHistory), nil
}

type pipelineHistoryController struct {
	controller.GenericController
}

func (c *pipelineHistoryController) Lister() PipelineHistoryLister {
	return &pipelineHistoryLister{
		controller: c,
	}
}

func (c *pipelineHistoryController) AddHandler(name string, handler PipelineHistoryHandlerFunc) {
	c.GenericController.AddHandler(name, func(key string) error {
		obj, exists, err := c.Informer().GetStore().GetByKey(key)
		if err != nil {
			return err
		}
		if !exists {
			return handler(key, nil)
		}
		return handler(key, obj.(*PipelineHistory))
	})
}

func (c *pipelineHistoryController) AddClusterScopedHandler(name, cluster string, handler PipelineHistoryHandlerFunc) {
	c.GenericController.AddHandler(name, func(key string) error {
		obj, exists, err := c.Informer().GetStore().GetByKey(key)
		if err != nil {
			return err
		}
		if !exists {
			return handler(key, nil)
		}

		if !controller.ObjectInCluster(cluster, obj) {
			return nil
		}

		return handler(key, obj.(*PipelineHistory))
	})
}

type pipelineHistoryFactory struct {
}

func (c pipelineHistoryFactory) Object() runtime.Object {
	return &PipelineHistory{}
}

func (c pipelineHistoryFactory) List() runtime.Object {
	return &PipelineHistoryList{}
}

func (s *pipelineHistoryClient) Controller() PipelineHistoryController {
	s.client.Lock()
	defer s.client.Unlock()

	c, ok := s.client.pipelineHistoryControllers[s.ns]
	if ok {
		return c
	}

	genericController := controller.NewGenericController(PipelineHistoryGroupVersionKind.Kind+"Controller",
		s.objectClient)

	c = &pipelineHistoryController{
		GenericController: genericController,
	}

	s.client.pipelineHistoryControllers[s.ns] = c
	s.client.starters = append(s.client.starters, c)

	return c
}

type pipelineHistoryClient struct {
	client       *Client
	ns           string
	objectClient *clientbase.ObjectClient
	controller   PipelineHistoryController
}

func (s *pipelineHistoryClient) ObjectClient() *clientbase.ObjectClient {
	return s.objectClient
}

func (s *pipelineHistoryClient) Create(o *PipelineHistory) (*PipelineHistory, error) {
	obj, err := s.objectClient.Create(o)
	return obj.(*PipelineHistory), err
}

func (s *pipelineHistoryClient) Get(name string, opts metav1.GetOptions) (*PipelineHistory, error) {
	obj, err := s.objectClient.Get(name, opts)
	return obj.(*PipelineHistory), err
}

func (s *pipelineHistoryClient) GetNamespaced(namespace, name string, opts metav1.GetOptions) (*PipelineHistory, error) {
	obj, err := s.objectClient.GetNamespaced(namespace, name, opts)
	return obj.(*PipelineHistory), err
}

func (s *pipelineHistoryClient) Update(o *PipelineHistory) (*PipelineHistory, error) {
	obj, err := s.objectClient.Update(o.Name, o)
	return obj.(*PipelineHistory), err
}

func (s *pipelineHistoryClient) Delete(name string, options *metav1.DeleteOptions) error {
	return s.objectClient.Delete(name, options)
}

func (s *pipelineHistoryClient) DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error {
	return s.objectClient.DeleteNamespaced(namespace, name, options)
}

func (s *pipelineHistoryClient) List(opts metav1.ListOptions) (*PipelineHistoryList, error) {
	obj, err := s.objectClient.List(opts)
	return obj.(*PipelineHistoryList), err
}

func (s *pipelineHistoryClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return s.objectClient.Watch(opts)
}

// Patch applies the patch and returns the patched deployment.
func (s *pipelineHistoryClient) Patch(o *PipelineHistory, data []byte, subresources ...string) (*PipelineHistory, error) {
	obj, err := s.objectClient.Patch(o.Name, o, data, subresources...)
	return obj.(*PipelineHistory), err
}

func (s *pipelineHistoryClient) DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	return s.objectClient.DeleteCollection(deleteOpts, listOpts)
}

func (s *pipelineHistoryClient) AddHandler(name string, sync PipelineHistoryHandlerFunc) {
	s.Controller().AddHandler(name, sync)
}

func (s *pipelineHistoryClient) AddLifecycle(name string, lifecycle PipelineHistoryLifecycle) {
	sync := NewPipelineHistoryLifecycleAdapter(name, false, s, lifecycle)
	s.AddHandler(name, sync)
}

func (s *pipelineHistoryClient) AddClusterScopedHandler(name, clusterName string, sync PipelineHistoryHandlerFunc) {
	s.Controller().AddClusterScopedHandler(name, clusterName, sync)
}

func (s *pipelineHistoryClient) AddClusterScopedLifecycle(name, clusterName string, lifecycle PipelineHistoryLifecycle) {
	sync := NewPipelineHistoryLifecycleAdapter(name+"_"+clusterName, true, s, lifecycle)
	s.AddClusterScopedHandler(name, clusterName, sync)
}
