package v3

import (
	"context"

	"github.com/rancher/norman/controller"
	"github.com/rancher/norman/objectclient"
	"github.com/rancher/norman/resource"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

var (
	GlobalMonitorGraphGroupVersionKind = schema.GroupVersionKind{
		Version: Version,
		Group:   GroupName,
		Kind:    "GlobalMonitorGraph",
	}
	GlobalMonitorGraphResource = metav1.APIResource{
		Name:         "globalmonitorgraphs",
		SingularName: "globalmonitorgraph",
		Namespaced:   true,

		Kind: GlobalMonitorGraphGroupVersionKind.Kind,
	}

	GlobalMonitorGraphGroupVersionResource = schema.GroupVersionResource{
		Group:    GroupName,
		Version:  Version,
		Resource: "globalmonitorgraphs",
	}
)

func init() {
	resource.Put(GlobalMonitorGraphGroupVersionResource)
}

func NewGlobalMonitorGraph(namespace, name string, obj GlobalMonitorGraph) *GlobalMonitorGraph {
	obj.APIVersion, obj.Kind = GlobalMonitorGraphGroupVersionKind.ToAPIVersionAndKind()
	obj.Name = name
	obj.Namespace = namespace
	return &obj
}

type GlobalMonitorGraphList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GlobalMonitorGraph `json:"items"`
}

type GlobalMonitorGraphHandlerFunc func(key string, obj *GlobalMonitorGraph) (runtime.Object, error)

type GlobalMonitorGraphChangeHandlerFunc func(obj *GlobalMonitorGraph) (runtime.Object, error)

type GlobalMonitorGraphLister interface {
	List(namespace string, selector labels.Selector) (ret []*GlobalMonitorGraph, err error)
	Get(namespace, name string) (*GlobalMonitorGraph, error)
}

type GlobalMonitorGraphController interface {
	Generic() controller.GenericController
	Informer() cache.SharedIndexInformer
	Lister() GlobalMonitorGraphLister
	AddHandler(ctx context.Context, name string, handler GlobalMonitorGraphHandlerFunc)
	AddFeatureHandler(ctx context.Context, enabled func() bool, name string, sync GlobalMonitorGraphHandlerFunc)
	AddClusterScopedHandler(ctx context.Context, name, clusterName string, handler GlobalMonitorGraphHandlerFunc)
	AddClusterScopedFeatureHandler(ctx context.Context, enabled func() bool, name, clusterName string, handler GlobalMonitorGraphHandlerFunc)
	Enqueue(namespace, name string)
	Sync(ctx context.Context) error
	Start(ctx context.Context, threadiness int) error
}

type GlobalMonitorGraphInterface interface {
	ObjectClient() *objectclient.ObjectClient
	Create(*GlobalMonitorGraph) (*GlobalMonitorGraph, error)
	GetNamespaced(namespace, name string, opts metav1.GetOptions) (*GlobalMonitorGraph, error)
	Get(name string, opts metav1.GetOptions) (*GlobalMonitorGraph, error)
	Update(*GlobalMonitorGraph) (*GlobalMonitorGraph, error)
	Delete(name string, options *metav1.DeleteOptions) error
	DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error
	List(opts metav1.ListOptions) (*GlobalMonitorGraphList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Controller() GlobalMonitorGraphController
	AddHandler(ctx context.Context, name string, sync GlobalMonitorGraphHandlerFunc)
	AddFeatureHandler(ctx context.Context, enabled func() bool, name string, sync GlobalMonitorGraphHandlerFunc)
	AddLifecycle(ctx context.Context, name string, lifecycle GlobalMonitorGraphLifecycle)
	AddFeatureLifecycle(ctx context.Context, enabled func() bool, name string, lifecycle GlobalMonitorGraphLifecycle)
	AddClusterScopedHandler(ctx context.Context, name, clusterName string, sync GlobalMonitorGraphHandlerFunc)
	AddClusterScopedFeatureHandler(ctx context.Context, enabled func() bool, name, clusterName string, sync GlobalMonitorGraphHandlerFunc)
	AddClusterScopedLifecycle(ctx context.Context, name, clusterName string, lifecycle GlobalMonitorGraphLifecycle)
	AddClusterScopedFeatureLifecycle(ctx context.Context, enabled func() bool, name, clusterName string, lifecycle GlobalMonitorGraphLifecycle)
}

type globalMonitorGraphLister struct {
	controller *globalMonitorGraphController
}

func (l *globalMonitorGraphLister) List(namespace string, selector labels.Selector) (ret []*GlobalMonitorGraph, err error) {
	err = cache.ListAllByNamespace(l.controller.Informer().GetIndexer(), namespace, selector, func(obj interface{}) {
		ret = append(ret, obj.(*GlobalMonitorGraph))
	})
	return
}

func (l *globalMonitorGraphLister) Get(namespace, name string) (*GlobalMonitorGraph, error) {
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
			Group:    GlobalMonitorGraphGroupVersionKind.Group,
			Resource: "globalMonitorGraph",
		}, key)
	}
	return obj.(*GlobalMonitorGraph), nil
}

type globalMonitorGraphController struct {
	controller.GenericController
}

func (c *globalMonitorGraphController) Generic() controller.GenericController {
	return c.GenericController
}

func (c *globalMonitorGraphController) Lister() GlobalMonitorGraphLister {
	return &globalMonitorGraphLister{
		controller: c,
	}
}

func (c *globalMonitorGraphController) AddHandler(ctx context.Context, name string, handler GlobalMonitorGraphHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*GlobalMonitorGraph); ok {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

func (c *globalMonitorGraphController) AddFeatureHandler(ctx context.Context, enabled func() bool, name string, handler GlobalMonitorGraphHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if !enabled() {
			return nil, nil
		} else if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*GlobalMonitorGraph); ok {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

func (c *globalMonitorGraphController) AddClusterScopedHandler(ctx context.Context, name, cluster string, handler GlobalMonitorGraphHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*GlobalMonitorGraph); ok && controller.ObjectInCluster(cluster, obj) {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

func (c *globalMonitorGraphController) AddClusterScopedFeatureHandler(ctx context.Context, enabled func() bool, name, cluster string, handler GlobalMonitorGraphHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if !enabled() {
			return nil, nil
		} else if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*GlobalMonitorGraph); ok && controller.ObjectInCluster(cluster, obj) {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

type globalMonitorGraphFactory struct {
}

func (c globalMonitorGraphFactory) Object() runtime.Object {
	return &GlobalMonitorGraph{}
}

func (c globalMonitorGraphFactory) List() runtime.Object {
	return &GlobalMonitorGraphList{}
}

func (s *globalMonitorGraphClient) Controller() GlobalMonitorGraphController {
	s.client.Lock()
	defer s.client.Unlock()

	c, ok := s.client.globalMonitorGraphControllers[s.ns]
	if ok {
		return c
	}

	genericController := controller.NewGenericController(GlobalMonitorGraphGroupVersionKind.Kind+"Controller",
		s.objectClient)

	c = &globalMonitorGraphController{
		GenericController: genericController,
	}

	s.client.globalMonitorGraphControllers[s.ns] = c
	s.client.starters = append(s.client.starters, c)

	return c
}

type globalMonitorGraphClient struct {
	client       *Client
	ns           string
	objectClient *objectclient.ObjectClient
	controller   GlobalMonitorGraphController
}

func (s *globalMonitorGraphClient) ObjectClient() *objectclient.ObjectClient {
	return s.objectClient
}

func (s *globalMonitorGraphClient) Create(o *GlobalMonitorGraph) (*GlobalMonitorGraph, error) {
	obj, err := s.objectClient.Create(o)
	return obj.(*GlobalMonitorGraph), err
}

func (s *globalMonitorGraphClient) Get(name string, opts metav1.GetOptions) (*GlobalMonitorGraph, error) {
	obj, err := s.objectClient.Get(name, opts)
	return obj.(*GlobalMonitorGraph), err
}

func (s *globalMonitorGraphClient) GetNamespaced(namespace, name string, opts metav1.GetOptions) (*GlobalMonitorGraph, error) {
	obj, err := s.objectClient.GetNamespaced(namespace, name, opts)
	return obj.(*GlobalMonitorGraph), err
}

func (s *globalMonitorGraphClient) Update(o *GlobalMonitorGraph) (*GlobalMonitorGraph, error) {
	obj, err := s.objectClient.Update(o.Name, o)
	return obj.(*GlobalMonitorGraph), err
}

func (s *globalMonitorGraphClient) Delete(name string, options *metav1.DeleteOptions) error {
	return s.objectClient.Delete(name, options)
}

func (s *globalMonitorGraphClient) DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error {
	return s.objectClient.DeleteNamespaced(namespace, name, options)
}

func (s *globalMonitorGraphClient) List(opts metav1.ListOptions) (*GlobalMonitorGraphList, error) {
	obj, err := s.objectClient.List(opts)
	return obj.(*GlobalMonitorGraphList), err
}

func (s *globalMonitorGraphClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return s.objectClient.Watch(opts)
}

// Patch applies the patch and returns the patched deployment.
func (s *globalMonitorGraphClient) Patch(o *GlobalMonitorGraph, patchType types.PatchType, data []byte, subresources ...string) (*GlobalMonitorGraph, error) {
	obj, err := s.objectClient.Patch(o.Name, o, patchType, data, subresources...)
	return obj.(*GlobalMonitorGraph), err
}

func (s *globalMonitorGraphClient) DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	return s.objectClient.DeleteCollection(deleteOpts, listOpts)
}

func (s *globalMonitorGraphClient) AddHandler(ctx context.Context, name string, sync GlobalMonitorGraphHandlerFunc) {
	s.Controller().AddHandler(ctx, name, sync)
}

func (s *globalMonitorGraphClient) AddFeatureHandler(ctx context.Context, enabled func() bool, name string, sync GlobalMonitorGraphHandlerFunc) {
	s.Controller().AddFeatureHandler(ctx, enabled, name, sync)
}

func (s *globalMonitorGraphClient) AddLifecycle(ctx context.Context, name string, lifecycle GlobalMonitorGraphLifecycle) {
	sync := NewGlobalMonitorGraphLifecycleAdapter(name, false, s, lifecycle)
	s.Controller().AddHandler(ctx, name, sync)
}

func (s *globalMonitorGraphClient) AddFeatureLifecycle(ctx context.Context, enabled func() bool, name string, lifecycle GlobalMonitorGraphLifecycle) {
	sync := NewGlobalMonitorGraphLifecycleAdapter(name, false, s, lifecycle)
	s.Controller().AddFeatureHandler(ctx, enabled, name, sync)
}

func (s *globalMonitorGraphClient) AddClusterScopedHandler(ctx context.Context, name, clusterName string, sync GlobalMonitorGraphHandlerFunc) {
	s.Controller().AddClusterScopedHandler(ctx, name, clusterName, sync)
}

func (s *globalMonitorGraphClient) AddClusterScopedFeatureHandler(ctx context.Context, enabled func() bool, name, clusterName string, sync GlobalMonitorGraphHandlerFunc) {
	s.Controller().AddClusterScopedFeatureHandler(ctx, enabled, name, clusterName, sync)
}

func (s *globalMonitorGraphClient) AddClusterScopedLifecycle(ctx context.Context, name, clusterName string, lifecycle GlobalMonitorGraphLifecycle) {
	sync := NewGlobalMonitorGraphLifecycleAdapter(name+"_"+clusterName, true, s, lifecycle)
	s.Controller().AddClusterScopedHandler(ctx, name, clusterName, sync)
}

func (s *globalMonitorGraphClient) AddClusterScopedFeatureLifecycle(ctx context.Context, enabled func() bool, name, clusterName string, lifecycle GlobalMonitorGraphLifecycle) {
	sync := NewGlobalMonitorGraphLifecycleAdapter(name+"_"+clusterName, true, s, lifecycle)
	s.Controller().AddClusterScopedFeatureHandler(ctx, enabled, name, clusterName, sync)
}

type GlobalMonitorGraphIndexer func(obj *GlobalMonitorGraph) ([]string, error)

type GlobalMonitorGraphClientCache interface {
	Get(namespace, name string) (*GlobalMonitorGraph, error)
	List(namespace string, selector labels.Selector) ([]*GlobalMonitorGraph, error)

	Index(name string, indexer GlobalMonitorGraphIndexer)
	GetIndexed(name, key string) ([]*GlobalMonitorGraph, error)
}

type GlobalMonitorGraphClient interface {
	Create(*GlobalMonitorGraph) (*GlobalMonitorGraph, error)
	Get(namespace, name string, opts metav1.GetOptions) (*GlobalMonitorGraph, error)
	Update(*GlobalMonitorGraph) (*GlobalMonitorGraph, error)
	Delete(namespace, name string, options *metav1.DeleteOptions) error
	List(namespace string, opts metav1.ListOptions) (*GlobalMonitorGraphList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)

	Cache() GlobalMonitorGraphClientCache

	OnCreate(ctx context.Context, name string, sync GlobalMonitorGraphChangeHandlerFunc)
	OnChange(ctx context.Context, name string, sync GlobalMonitorGraphChangeHandlerFunc)
	OnRemove(ctx context.Context, name string, sync GlobalMonitorGraphChangeHandlerFunc)
	Enqueue(namespace, name string)

	Generic() controller.GenericController
	ObjectClient() *objectclient.ObjectClient
	Interface() GlobalMonitorGraphInterface
}

type globalMonitorGraphClientCache struct {
	client *globalMonitorGraphClient2
}

type globalMonitorGraphClient2 struct {
	iface      GlobalMonitorGraphInterface
	controller GlobalMonitorGraphController
}

func (n *globalMonitorGraphClient2) Interface() GlobalMonitorGraphInterface {
	return n.iface
}

func (n *globalMonitorGraphClient2) Generic() controller.GenericController {
	return n.iface.Controller().Generic()
}

func (n *globalMonitorGraphClient2) ObjectClient() *objectclient.ObjectClient {
	return n.Interface().ObjectClient()
}

func (n *globalMonitorGraphClient2) Enqueue(namespace, name string) {
	n.iface.Controller().Enqueue(namespace, name)
}

func (n *globalMonitorGraphClient2) Create(obj *GlobalMonitorGraph) (*GlobalMonitorGraph, error) {
	return n.iface.Create(obj)
}

func (n *globalMonitorGraphClient2) Get(namespace, name string, opts metav1.GetOptions) (*GlobalMonitorGraph, error) {
	return n.iface.GetNamespaced(namespace, name, opts)
}

func (n *globalMonitorGraphClient2) Update(obj *GlobalMonitorGraph) (*GlobalMonitorGraph, error) {
	return n.iface.Update(obj)
}

func (n *globalMonitorGraphClient2) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	return n.iface.DeleteNamespaced(namespace, name, options)
}

func (n *globalMonitorGraphClient2) List(namespace string, opts metav1.ListOptions) (*GlobalMonitorGraphList, error) {
	return n.iface.List(opts)
}

func (n *globalMonitorGraphClient2) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return n.iface.Watch(opts)
}

func (n *globalMonitorGraphClientCache) Get(namespace, name string) (*GlobalMonitorGraph, error) {
	return n.client.controller.Lister().Get(namespace, name)
}

func (n *globalMonitorGraphClientCache) List(namespace string, selector labels.Selector) ([]*GlobalMonitorGraph, error) {
	return n.client.controller.Lister().List(namespace, selector)
}

func (n *globalMonitorGraphClient2) Cache() GlobalMonitorGraphClientCache {
	n.loadController()
	return &globalMonitorGraphClientCache{
		client: n,
	}
}

func (n *globalMonitorGraphClient2) OnCreate(ctx context.Context, name string, sync GlobalMonitorGraphChangeHandlerFunc) {
	n.loadController()
	n.iface.AddLifecycle(ctx, name+"-create", &globalMonitorGraphLifecycleDelegate{create: sync})
}

func (n *globalMonitorGraphClient2) OnChange(ctx context.Context, name string, sync GlobalMonitorGraphChangeHandlerFunc) {
	n.loadController()
	n.iface.AddLifecycle(ctx, name+"-change", &globalMonitorGraphLifecycleDelegate{update: sync})
}

func (n *globalMonitorGraphClient2) OnRemove(ctx context.Context, name string, sync GlobalMonitorGraphChangeHandlerFunc) {
	n.loadController()
	n.iface.AddLifecycle(ctx, name, &globalMonitorGraphLifecycleDelegate{remove: sync})
}

func (n *globalMonitorGraphClientCache) Index(name string, indexer GlobalMonitorGraphIndexer) {
	err := n.client.controller.Informer().GetIndexer().AddIndexers(map[string]cache.IndexFunc{
		name: func(obj interface{}) ([]string, error) {
			if v, ok := obj.(*GlobalMonitorGraph); ok {
				return indexer(v)
			}
			return nil, nil
		},
	})

	if err != nil {
		panic(err)
	}
}

func (n *globalMonitorGraphClientCache) GetIndexed(name, key string) ([]*GlobalMonitorGraph, error) {
	var result []*GlobalMonitorGraph
	objs, err := n.client.controller.Informer().GetIndexer().ByIndex(name, key)
	if err != nil {
		return nil, err
	}
	for _, obj := range objs {
		if v, ok := obj.(*GlobalMonitorGraph); ok {
			result = append(result, v)
		}
	}

	return result, nil
}

func (n *globalMonitorGraphClient2) loadController() {
	if n.controller == nil {
		n.controller = n.iface.Controller()
	}
}

type globalMonitorGraphLifecycleDelegate struct {
	create GlobalMonitorGraphChangeHandlerFunc
	update GlobalMonitorGraphChangeHandlerFunc
	remove GlobalMonitorGraphChangeHandlerFunc
}

func (n *globalMonitorGraphLifecycleDelegate) HasCreate() bool {
	return n.create != nil
}

func (n *globalMonitorGraphLifecycleDelegate) Create(obj *GlobalMonitorGraph) (runtime.Object, error) {
	if n.create == nil {
		return obj, nil
	}
	return n.create(obj)
}

func (n *globalMonitorGraphLifecycleDelegate) HasFinalize() bool {
	return n.remove != nil
}

func (n *globalMonitorGraphLifecycleDelegate) Remove(obj *GlobalMonitorGraph) (runtime.Object, error) {
	if n.remove == nil {
		return obj, nil
	}
	return n.remove(obj)
}

func (n *globalMonitorGraphLifecycleDelegate) Updated(obj *GlobalMonitorGraph) (runtime.Object, error) {
	if n.update == nil {
		return obj, nil
	}
	return n.update(obj)
}
