package client

import (
	"github.com/rancher/norman/types"
)

const (
	PipelineHistoryType                      = "pipelineHistory"
	PipelineHistoryFieldAnnotations          = "annotations"
	PipelineHistoryFieldCreated              = "created"
	PipelineHistoryFieldCreatorID            = "creatorId"
	PipelineHistoryFieldDisplayName          = "displayName"
	PipelineHistoryFieldLabels               = "labels"
	PipelineHistoryFieldName                 = "name"
	PipelineHistoryFieldNamespaceId          = "namespaceId"
	PipelineHistoryFieldOwnerReferences      = "ownerReferences"
	PipelineHistoryFieldPipeline             = "pipeline"
	PipelineHistoryFieldProjectId            = "projectId"
	PipelineHistoryFieldRemoved              = "removed"
	PipelineHistoryFieldRunNumber            = "runNumber"
	PipelineHistoryFieldState                = "state"
	PipelineHistoryFieldStatus               = "status"
	PipelineHistoryFieldTransitioning        = "transitioning"
	PipelineHistoryFieldTransitioningMessage = "transitioningMessage"
	PipelineHistoryFieldTriggerType          = "triggerType"
	PipelineHistoryFieldUuid                 = "uuid"
)

type PipelineHistory struct {
	types.Resource
	Annotations          map[string]string      `json:"annotations,omitempty"`
	Created              string                 `json:"created,omitempty"`
	CreatorID            string                 `json:"creatorId,omitempty"`
	DisplayName          string                 `json:"displayName,omitempty"`
	Labels               map[string]string      `json:"labels,omitempty"`
	Name                 string                 `json:"name,omitempty"`
	NamespaceId          string                 `json:"namespaceId,omitempty"`
	OwnerReferences      []OwnerReference       `json:"ownerReferences,omitempty"`
	Pipeline             *Pipeline              `json:"pipeline,omitempty"`
	ProjectId            string                 `json:"projectId,omitempty"`
	Removed              string                 `json:"removed,omitempty"`
	RunNumber            *int64                 `json:"runNumber,omitempty"`
	State                string                 `json:"state,omitempty"`
	Status               *PipelineHistoryStatus `json:"status,omitempty"`
	Transitioning        string                 `json:"transitioning,omitempty"`
	TransitioningMessage string                 `json:"transitioningMessage,omitempty"`
	TriggerType          string                 `json:"triggerType,omitempty"`
	Uuid                 string                 `json:"uuid,omitempty"`
}
type PipelineHistoryCollection struct {
	types.Collection
	Data   []PipelineHistory `json:"data,omitempty"`
	client *PipelineHistoryClient
}

type PipelineHistoryClient struct {
	apiClient *Client
}

type PipelineHistoryOperations interface {
	List(opts *types.ListOpts) (*PipelineHistoryCollection, error)
	Create(opts *PipelineHistory) (*PipelineHistory, error)
	Update(existing *PipelineHistory, updates interface{}) (*PipelineHistory, error)
	ByID(id string) (*PipelineHistory, error)
	Delete(container *PipelineHistory) error
}

func newPipelineHistoryClient(apiClient *Client) *PipelineHistoryClient {
	return &PipelineHistoryClient{
		apiClient: apiClient,
	}
}

func (c *PipelineHistoryClient) Create(container *PipelineHistory) (*PipelineHistory, error) {
	resp := &PipelineHistory{}
	err := c.apiClient.Ops.DoCreate(PipelineHistoryType, container, resp)
	return resp, err
}

func (c *PipelineHistoryClient) Update(existing *PipelineHistory, updates interface{}) (*PipelineHistory, error) {
	resp := &PipelineHistory{}
	err := c.apiClient.Ops.DoUpdate(PipelineHistoryType, &existing.Resource, updates, resp)
	return resp, err
}

func (c *PipelineHistoryClient) List(opts *types.ListOpts) (*PipelineHistoryCollection, error) {
	resp := &PipelineHistoryCollection{}
	err := c.apiClient.Ops.DoList(PipelineHistoryType, opts, resp)
	resp.client = c
	return resp, err
}

func (cc *PipelineHistoryCollection) Next() (*PipelineHistoryCollection, error) {
	if cc != nil && cc.Pagination != nil && cc.Pagination.Next != "" {
		resp := &PipelineHistoryCollection{}
		err := cc.client.apiClient.Ops.DoNext(cc.Pagination.Next, resp)
		resp.client = cc.client
		return resp, err
	}
	return nil, nil
}

func (c *PipelineHistoryClient) ByID(id string) (*PipelineHistory, error) {
	resp := &PipelineHistory{}
	err := c.apiClient.Ops.DoByID(PipelineHistoryType, id, resp)
	return resp, err
}

func (c *PipelineHistoryClient) Delete(container *PipelineHistory) error {
	return c.apiClient.Ops.DoResourceDelete(PipelineHistoryType, &container.Resource)
}
