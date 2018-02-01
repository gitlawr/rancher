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
	RemoteAccountGroupVersionKind = schema.GroupVersionKind{
		Version: Version,
		Group:   GroupName,
		Kind:    "RemoteAccount",
	}
	RemoteAccountResource = metav1.APIResource{
		Name:         "remoteaccounts",
		SingularName: "remoteaccount",
		Namespaced:   false,
		Kind:         RemoteAccountGroupVersionKind.Kind,
	}
)

type RemoteAccountList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RemoteAccount
}

type RemoteAccountHandlerFunc func(key string, obj *RemoteAccount) error

type RemoteAccountLister interface {
	List(namespace string, selector labels.Selector) (ret []*RemoteAccount, err error)
	Get(namespace, name string) (*RemoteAccount, error)
}

type RemoteAccountController interface {
	Informer() cache.SharedIndexInformer
	Lister() RemoteAccountLister
	AddHandler(name string, handler RemoteAccountHandlerFunc)
	AddClusterScopedHandler(name, clusterName string, handler RemoteAccountHandlerFunc)
	Enqueue(namespace, name string)
	Sync(ctx context.Context) error
	Start(ctx context.Context, threadiness int) error
}

type RemoteAccountInterface interface {
	ObjectClient() *clientbase.ObjectClient
	Create(*RemoteAccount) (*RemoteAccount, error)
	GetNamespaced(namespace, name string, opts metav1.GetOptions) (*RemoteAccount, error)
	Get(name string, opts metav1.GetOptions) (*RemoteAccount, error)
	Update(*RemoteAccount) (*RemoteAccount, error)
	Delete(name string, options *metav1.DeleteOptions) error
	DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error
	List(opts metav1.ListOptions) (*RemoteAccountList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Controller() RemoteAccountController
	AddHandler(name string, sync RemoteAccountHandlerFunc)
	AddLifecycle(name string, lifecycle RemoteAccountLifecycle)
	AddClusterScopedHandler(name, clusterName string, sync RemoteAccountHandlerFunc)
	AddClusterScopedLifecycle(name, clusterName string, lifecycle RemoteAccountLifecycle)
}

type remoteAccountLister struct {
	controller *remoteAccountController
}

func (l *remoteAccountLister) List(namespace string, selector labels.Selector) (ret []*RemoteAccount, err error) {
	err = cache.ListAllByNamespace(l.controller.Informer().GetIndexer(), namespace, selector, func(obj interface{}) {
		ret = append(ret, obj.(*RemoteAccount))
	})
	return
}

func (l *remoteAccountLister) Get(namespace, name string) (*RemoteAccount, error) {
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
			Group:    RemoteAccountGroupVersionKind.Group,
			Resource: "remoteAccount",
		}, name)
	}
	return obj.(*RemoteAccount), nil
}

type remoteAccountController struct {
	controller.GenericController
}

func (c *remoteAccountController) Lister() RemoteAccountLister {
	return &remoteAccountLister{
		controller: c,
	}
}

func (c *remoteAccountController) AddHandler(name string, handler RemoteAccountHandlerFunc) {
	c.GenericController.AddHandler(name, func(key string) error {
		obj, exists, err := c.Informer().GetStore().GetByKey(key)
		if err != nil {
			return err
		}
		if !exists {
			return handler(key, nil)
		}
		return handler(key, obj.(*RemoteAccount))
	})
}

func (c *remoteAccountController) AddClusterScopedHandler(name, cluster string, handler RemoteAccountHandlerFunc) {
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

		return handler(key, obj.(*RemoteAccount))
	})
}

type remoteAccountFactory struct {
}

func (c remoteAccountFactory) Object() runtime.Object {
	return &RemoteAccount{}
}

func (c remoteAccountFactory) List() runtime.Object {
	return &RemoteAccountList{}
}

func (s *remoteAccountClient) Controller() RemoteAccountController {
	s.client.Lock()
	defer s.client.Unlock()

	c, ok := s.client.remoteAccountControllers[s.ns]
	if ok {
		return c
	}

	genericController := controller.NewGenericController(RemoteAccountGroupVersionKind.Kind+"Controller",
		s.objectClient)

	c = &remoteAccountController{
		GenericController: genericController,
	}

	s.client.remoteAccountControllers[s.ns] = c
	s.client.starters = append(s.client.starters, c)

	return c
}

type remoteAccountClient struct {
	client       *Client
	ns           string
	objectClient *clientbase.ObjectClient
	controller   RemoteAccountController
}

func (s *remoteAccountClient) ObjectClient() *clientbase.ObjectClient {
	return s.objectClient
}

func (s *remoteAccountClient) Create(o *RemoteAccount) (*RemoteAccount, error) {
	obj, err := s.objectClient.Create(o)
	return obj.(*RemoteAccount), err
}

func (s *remoteAccountClient) Get(name string, opts metav1.GetOptions) (*RemoteAccount, error) {
	obj, err := s.objectClient.Get(name, opts)
	return obj.(*RemoteAccount), err
}

func (s *remoteAccountClient) GetNamespaced(namespace, name string, opts metav1.GetOptions) (*RemoteAccount, error) {
	obj, err := s.objectClient.GetNamespaced(namespace, name, opts)
	return obj.(*RemoteAccount), err
}

func (s *remoteAccountClient) Update(o *RemoteAccount) (*RemoteAccount, error) {
	obj, err := s.objectClient.Update(o.Name, o)
	return obj.(*RemoteAccount), err
}

func (s *remoteAccountClient) Delete(name string, options *metav1.DeleteOptions) error {
	return s.objectClient.Delete(name, options)
}

func (s *remoteAccountClient) DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error {
	return s.objectClient.DeleteNamespaced(namespace, name, options)
}

func (s *remoteAccountClient) List(opts metav1.ListOptions) (*RemoteAccountList, error) {
	obj, err := s.objectClient.List(opts)
	return obj.(*RemoteAccountList), err
}

func (s *remoteAccountClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return s.objectClient.Watch(opts)
}

// Patch applies the patch and returns the patched deployment.
func (s *remoteAccountClient) Patch(o *RemoteAccount, data []byte, subresources ...string) (*RemoteAccount, error) {
	obj, err := s.objectClient.Patch(o.Name, o, data, subresources...)
	return obj.(*RemoteAccount), err
}

func (s *remoteAccountClient) DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	return s.objectClient.DeleteCollection(deleteOpts, listOpts)
}

func (s *remoteAccountClient) AddHandler(name string, sync RemoteAccountHandlerFunc) {
	s.Controller().AddHandler(name, sync)
}

func (s *remoteAccountClient) AddLifecycle(name string, lifecycle RemoteAccountLifecycle) {
	sync := NewRemoteAccountLifecycleAdapter(name, false, s, lifecycle)
	s.AddHandler(name, sync)
}

func (s *remoteAccountClient) AddClusterScopedHandler(name, clusterName string, sync RemoteAccountHandlerFunc) {
	s.Controller().AddClusterScopedHandler(name, clusterName, sync)
}

func (s *remoteAccountClient) AddClusterScopedLifecycle(name, clusterName string, lifecycle RemoteAccountLifecycle) {
	sync := NewRemoteAccountLifecycleAdapter(name+"_"+clusterName, true, s, lifecycle)
	s.AddClusterScopedHandler(name, clusterName, sync)
}
