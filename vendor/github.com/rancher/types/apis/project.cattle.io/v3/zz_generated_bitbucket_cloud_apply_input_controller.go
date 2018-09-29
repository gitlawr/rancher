package v3

import (
	"context"

	"github.com/rancher/norman/controller"
	"github.com/rancher/norman/objectclient"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

var (
	BitbucketCloudApplyInputGroupVersionKind = schema.GroupVersionKind{
		Version: Version,
		Group:   GroupName,
		Kind:    "BitbucketCloudApplyInput",
	}
	BitbucketCloudApplyInputResource = metav1.APIResource{
		Name:         "bitbucketcloudapplyinputs",
		SingularName: "bitbucketcloudapplyinput",
		Namespaced:   true,

		Kind: BitbucketCloudApplyInputGroupVersionKind.Kind,
	}
)

type BitbucketCloudApplyInputList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []BitbucketCloudApplyInput
}

type BitbucketCloudApplyInputHandlerFunc func(key string, obj *BitbucketCloudApplyInput) error

type BitbucketCloudApplyInputLister interface {
	List(namespace string, selector labels.Selector) (ret []*BitbucketCloudApplyInput, err error)
	Get(namespace, name string) (*BitbucketCloudApplyInput, error)
}

type BitbucketCloudApplyInputController interface {
	Informer() cache.SharedIndexInformer
	Lister() BitbucketCloudApplyInputLister
	AddHandler(name string, handler BitbucketCloudApplyInputHandlerFunc)
	AddClusterScopedHandler(name, clusterName string, handler BitbucketCloudApplyInputHandlerFunc)
	Enqueue(namespace, name string)
	Sync(ctx context.Context) error
	Start(ctx context.Context, threadiness int) error
}

type BitbucketCloudApplyInputInterface interface {
	ObjectClient() *objectclient.ObjectClient
	Create(*BitbucketCloudApplyInput) (*BitbucketCloudApplyInput, error)
	GetNamespaced(namespace, name string, opts metav1.GetOptions) (*BitbucketCloudApplyInput, error)
	Get(name string, opts metav1.GetOptions) (*BitbucketCloudApplyInput, error)
	Update(*BitbucketCloudApplyInput) (*BitbucketCloudApplyInput, error)
	Delete(name string, options *metav1.DeleteOptions) error
	DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error
	List(opts metav1.ListOptions) (*BitbucketCloudApplyInputList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Controller() BitbucketCloudApplyInputController
	AddHandler(name string, sync BitbucketCloudApplyInputHandlerFunc)
	AddLifecycle(name string, lifecycle BitbucketCloudApplyInputLifecycle)
	AddClusterScopedHandler(name, clusterName string, sync BitbucketCloudApplyInputHandlerFunc)
	AddClusterScopedLifecycle(name, clusterName string, lifecycle BitbucketCloudApplyInputLifecycle)
}

type bitbucketCloudApplyInputLister struct {
	controller *bitbucketCloudApplyInputController
}

func (l *bitbucketCloudApplyInputLister) List(namespace string, selector labels.Selector) (ret []*BitbucketCloudApplyInput, err error) {
	err = cache.ListAllByNamespace(l.controller.Informer().GetIndexer(), namespace, selector, func(obj interface{}) {
		ret = append(ret, obj.(*BitbucketCloudApplyInput))
	})
	return
}

func (l *bitbucketCloudApplyInputLister) Get(namespace, name string) (*BitbucketCloudApplyInput, error) {
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
			Group:    BitbucketCloudApplyInputGroupVersionKind.Group,
			Resource: "bitbucketCloudApplyInput",
		}, key)
	}
	return obj.(*BitbucketCloudApplyInput), nil
}

type bitbucketCloudApplyInputController struct {
	controller.GenericController
}

func (c *bitbucketCloudApplyInputController) Lister() BitbucketCloudApplyInputLister {
	return &bitbucketCloudApplyInputLister{
		controller: c,
	}
}

func (c *bitbucketCloudApplyInputController) AddHandler(name string, handler BitbucketCloudApplyInputHandlerFunc) {
	c.GenericController.AddHandler(name, func(key string) error {
		obj, exists, err := c.Informer().GetStore().GetByKey(key)
		if err != nil {
			return err
		}
		if !exists {
			return handler(key, nil)
		}
		return handler(key, obj.(*BitbucketCloudApplyInput))
	})
}

func (c *bitbucketCloudApplyInputController) AddClusterScopedHandler(name, cluster string, handler BitbucketCloudApplyInputHandlerFunc) {
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

		return handler(key, obj.(*BitbucketCloudApplyInput))
	})
}

type bitbucketCloudApplyInputFactory struct {
}

func (c bitbucketCloudApplyInputFactory) Object() runtime.Object {
	return &BitbucketCloudApplyInput{}
}

func (c bitbucketCloudApplyInputFactory) List() runtime.Object {
	return &BitbucketCloudApplyInputList{}
}

func (s *bitbucketCloudApplyInputClient) Controller() BitbucketCloudApplyInputController {
	s.client.Lock()
	defer s.client.Unlock()

	c, ok := s.client.bitbucketCloudApplyInputControllers[s.ns]
	if ok {
		return c
	}

	genericController := controller.NewGenericController(BitbucketCloudApplyInputGroupVersionKind.Kind+"Controller",
		s.objectClient)

	c = &bitbucketCloudApplyInputController{
		GenericController: genericController,
	}

	s.client.bitbucketCloudApplyInputControllers[s.ns] = c
	s.client.starters = append(s.client.starters, c)

	return c
}

type bitbucketCloudApplyInputClient struct {
	client       *Client
	ns           string
	objectClient *objectclient.ObjectClient
	controller   BitbucketCloudApplyInputController
}

func (s *bitbucketCloudApplyInputClient) ObjectClient() *objectclient.ObjectClient {
	return s.objectClient
}

func (s *bitbucketCloudApplyInputClient) Create(o *BitbucketCloudApplyInput) (*BitbucketCloudApplyInput, error) {
	obj, err := s.objectClient.Create(o)
	return obj.(*BitbucketCloudApplyInput), err
}

func (s *bitbucketCloudApplyInputClient) Get(name string, opts metav1.GetOptions) (*BitbucketCloudApplyInput, error) {
	obj, err := s.objectClient.Get(name, opts)
	return obj.(*BitbucketCloudApplyInput), err
}

func (s *bitbucketCloudApplyInputClient) GetNamespaced(namespace, name string, opts metav1.GetOptions) (*BitbucketCloudApplyInput, error) {
	obj, err := s.objectClient.GetNamespaced(namespace, name, opts)
	return obj.(*BitbucketCloudApplyInput), err
}

func (s *bitbucketCloudApplyInputClient) Update(o *BitbucketCloudApplyInput) (*BitbucketCloudApplyInput, error) {
	obj, err := s.objectClient.Update(o.Name, o)
	return obj.(*BitbucketCloudApplyInput), err
}

func (s *bitbucketCloudApplyInputClient) Delete(name string, options *metav1.DeleteOptions) error {
	return s.objectClient.Delete(name, options)
}

func (s *bitbucketCloudApplyInputClient) DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error {
	return s.objectClient.DeleteNamespaced(namespace, name, options)
}

func (s *bitbucketCloudApplyInputClient) List(opts metav1.ListOptions) (*BitbucketCloudApplyInputList, error) {
	obj, err := s.objectClient.List(opts)
	return obj.(*BitbucketCloudApplyInputList), err
}

func (s *bitbucketCloudApplyInputClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return s.objectClient.Watch(opts)
}

// Patch applies the patch and returns the patched deployment.
func (s *bitbucketCloudApplyInputClient) Patch(o *BitbucketCloudApplyInput, data []byte, subresources ...string) (*BitbucketCloudApplyInput, error) {
	obj, err := s.objectClient.Patch(o.Name, o, data, subresources...)
	return obj.(*BitbucketCloudApplyInput), err
}

func (s *bitbucketCloudApplyInputClient) DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	return s.objectClient.DeleteCollection(deleteOpts, listOpts)
}

func (s *bitbucketCloudApplyInputClient) AddHandler(name string, sync BitbucketCloudApplyInputHandlerFunc) {
	s.Controller().AddHandler(name, sync)
}

func (s *bitbucketCloudApplyInputClient) AddLifecycle(name string, lifecycle BitbucketCloudApplyInputLifecycle) {
	sync := NewBitbucketCloudApplyInputLifecycleAdapter(name, false, s, lifecycle)
	s.AddHandler(name, sync)
}

func (s *bitbucketCloudApplyInputClient) AddClusterScopedHandler(name, clusterName string, sync BitbucketCloudApplyInputHandlerFunc) {
	s.Controller().AddClusterScopedHandler(name, clusterName, sync)
}

func (s *bitbucketCloudApplyInputClient) AddClusterScopedLifecycle(name, clusterName string, lifecycle BitbucketCloudApplyInputLifecycle) {
	sync := NewBitbucketCloudApplyInputLifecycleAdapter(name+"_"+clusterName, true, s, lifecycle)
	s.AddClusterScopedHandler(name, clusterName, sync)
}
