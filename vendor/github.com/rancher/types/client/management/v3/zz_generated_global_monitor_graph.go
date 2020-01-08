package client

import (
	"github.com/rancher/norman/types"
)

const (
	GlobalMonitorGraphType                        = "globalMonitorGraph"
	GlobalMonitorGraphFieldAnnotations            = "annotations"
	GlobalMonitorGraphFieldCreated                = "created"
	GlobalMonitorGraphFieldCreatorID              = "creatorId"
	GlobalMonitorGraphFieldDescription            = "description"
	GlobalMonitorGraphFieldDetailsMetricsSelector = "detailsMetricsSelector"
	GlobalMonitorGraphFieldDisplayResourceType    = "displayResourceType"
	GlobalMonitorGraphFieldGraphType              = "graphType"
	GlobalMonitorGraphFieldLabels                 = "labels"
	GlobalMonitorGraphFieldMetricsSelector        = "metricsSelector"
	GlobalMonitorGraphFieldName                   = "name"
	GlobalMonitorGraphFieldNamespaceId            = "namespaceId"
	GlobalMonitorGraphFieldOwnerReferences        = "ownerReferences"
	GlobalMonitorGraphFieldPriority               = "priority"
	GlobalMonitorGraphFieldRemoved                = "removed"
	GlobalMonitorGraphFieldResourceType           = "resourceType"
	GlobalMonitorGraphFieldUUID                   = "uuid"
	GlobalMonitorGraphFieldYAxis                  = "yAxis"
)

type GlobalMonitorGraph struct {
	types.Resource
	Annotations            map[string]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`
	Created                string            `json:"created,omitempty" yaml:"created,omitempty"`
	CreatorID              string            `json:"creatorId,omitempty" yaml:"creatorId,omitempty"`
	Description            string            `json:"description,omitempty" yaml:"description,omitempty"`
	DetailsMetricsSelector map[string]string `json:"detailsMetricsSelector,omitempty" yaml:"detailsMetricsSelector,omitempty"`
	DisplayResourceType    string            `json:"displayResourceType,omitempty" yaml:"displayResourceType,omitempty"`
	GraphType              string            `json:"graphType,omitempty" yaml:"graphType,omitempty"`
	Labels                 map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	MetricsSelector        map[string]string `json:"metricsSelector,omitempty" yaml:"metricsSelector,omitempty"`
	Name                   string            `json:"name,omitempty" yaml:"name,omitempty"`
	NamespaceId            string            `json:"namespaceId,omitempty" yaml:"namespaceId,omitempty"`
	OwnerReferences        []OwnerReference  `json:"ownerReferences,omitempty" yaml:"ownerReferences,omitempty"`
	Priority               int64             `json:"priority,omitempty" yaml:"priority,omitempty"`
	Removed                string            `json:"removed,omitempty" yaml:"removed,omitempty"`
	ResourceType           string            `json:"resourceType,omitempty" yaml:"resourceType,omitempty"`
	UUID                   string            `json:"uuid,omitempty" yaml:"uuid,omitempty"`
	YAxis                  *YAxis            `json:"yAxis,omitempty" yaml:"yAxis,omitempty"`
}

type GlobalMonitorGraphCollection struct {
	types.Collection
	Data   []GlobalMonitorGraph `json:"data,omitempty"`
	client *GlobalMonitorGraphClient
}

type GlobalMonitorGraphClient struct {
	apiClient *Client
}

type GlobalMonitorGraphOperations interface {
	List(opts *types.ListOpts) (*GlobalMonitorGraphCollection, error)
	Create(opts *GlobalMonitorGraph) (*GlobalMonitorGraph, error)
	Update(existing *GlobalMonitorGraph, updates interface{}) (*GlobalMonitorGraph, error)
	Replace(existing *GlobalMonitorGraph) (*GlobalMonitorGraph, error)
	ByID(id string) (*GlobalMonitorGraph, error)
	Delete(container *GlobalMonitorGraph) error

	CollectionActionQuery(resource *GlobalMonitorGraphCollection, input *QueryGraphInput) (*QueryGlobalGraphOutput, error)
}

func newGlobalMonitorGraphClient(apiClient *Client) *GlobalMonitorGraphClient {
	return &GlobalMonitorGraphClient{
		apiClient: apiClient,
	}
}

func (c *GlobalMonitorGraphClient) Create(container *GlobalMonitorGraph) (*GlobalMonitorGraph, error) {
	resp := &GlobalMonitorGraph{}
	err := c.apiClient.Ops.DoCreate(GlobalMonitorGraphType, container, resp)
	return resp, err
}

func (c *GlobalMonitorGraphClient) Update(existing *GlobalMonitorGraph, updates interface{}) (*GlobalMonitorGraph, error) {
	resp := &GlobalMonitorGraph{}
	err := c.apiClient.Ops.DoUpdate(GlobalMonitorGraphType, &existing.Resource, updates, resp)
	return resp, err
}

func (c *GlobalMonitorGraphClient) Replace(obj *GlobalMonitorGraph) (*GlobalMonitorGraph, error) {
	resp := &GlobalMonitorGraph{}
	err := c.apiClient.Ops.DoReplace(GlobalMonitorGraphType, &obj.Resource, obj, resp)
	return resp, err
}

func (c *GlobalMonitorGraphClient) List(opts *types.ListOpts) (*GlobalMonitorGraphCollection, error) {
	resp := &GlobalMonitorGraphCollection{}
	err := c.apiClient.Ops.DoList(GlobalMonitorGraphType, opts, resp)
	resp.client = c
	return resp, err
}

func (cc *GlobalMonitorGraphCollection) Next() (*GlobalMonitorGraphCollection, error) {
	if cc != nil && cc.Pagination != nil && cc.Pagination.Next != "" {
		resp := &GlobalMonitorGraphCollection{}
		err := cc.client.apiClient.Ops.DoNext(cc.Pagination.Next, resp)
		resp.client = cc.client
		return resp, err
	}
	return nil, nil
}

func (c *GlobalMonitorGraphClient) ByID(id string) (*GlobalMonitorGraph, error) {
	resp := &GlobalMonitorGraph{}
	err := c.apiClient.Ops.DoByID(GlobalMonitorGraphType, id, resp)
	return resp, err
}

func (c *GlobalMonitorGraphClient) Delete(container *GlobalMonitorGraph) error {
	return c.apiClient.Ops.DoResourceDelete(GlobalMonitorGraphType, &container.Resource)
}

func (c *GlobalMonitorGraphClient) CollectionActionQuery(resource *GlobalMonitorGraphCollection, input *QueryGraphInput) (*QueryGlobalGraphOutput, error) {
	resp := &QueryGlobalGraphOutput{}
	err := c.apiClient.Ops.DoCollectionAction(GlobalMonitorGraphType, "query", &resource.Collection, input, resp)
	return resp, err
}
