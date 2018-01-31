package client

import (
	"github.com/rancher/norman/types"
)

const (
	PipelineLogType                     = "pipelineLog"
	PipelineLogFieldAnnotations         = "annotations"
	PipelineLogFieldCreated             = "created"
	PipelineLogFieldCreatorID           = "creatorId"
	PipelineLogFieldLabels              = "labels"
	PipelineLogFieldMessage             = "message"
	PipelineLogFieldName                = "name"
	PipelineLogFieldNamespaceId         = "namespaceId"
	PipelineLogFieldOwnerReferences     = "ownerReferences"
	PipelineLogFieldPipelineHistoryName = "pipelineHistoryName"
	PipelineLogFieldProjectId           = "projectId"
	PipelineLogFieldRemoved             = "removed"
	PipelineLogFieldStageOrdinal        = "stageOrdinal"
	PipelineLogFieldStepOrdinal         = "stepOrdinal"
	PipelineLogFieldUuid                = "uuid"
)

type PipelineLog struct {
	types.Resource
	Annotations         map[string]string `json:"annotations,omitempty"`
	Created             string            `json:"created,omitempty"`
	CreatorID           string            `json:"creatorId,omitempty"`
	Labels              map[string]string `json:"labels,omitempty"`
	Message             string            `json:"message,omitempty"`
	Name                string            `json:"name,omitempty"`
	NamespaceId         string            `json:"namespaceId,omitempty"`
	OwnerReferences     []OwnerReference  `json:"ownerReferences,omitempty"`
	PipelineHistoryName string            `json:"pipelineHistoryName,omitempty"`
	ProjectId           string            `json:"projectId,omitempty"`
	Removed             string            `json:"removed,omitempty"`
	StageOrdinal        *int64            `json:"stageOrdinal,omitempty"`
	StepOrdinal         *int64            `json:"stepOrdinal,omitempty"`
	Uuid                string            `json:"uuid,omitempty"`
}
type PipelineLogCollection struct {
	types.Collection
	Data   []PipelineLog `json:"data,omitempty"`
	client *PipelineLogClient
}

type PipelineLogClient struct {
	apiClient *Client
}

type PipelineLogOperations interface {
	List(opts *types.ListOpts) (*PipelineLogCollection, error)
	Create(opts *PipelineLog) (*PipelineLog, error)
	Update(existing *PipelineLog, updates interface{}) (*PipelineLog, error)
	ByID(id string) (*PipelineLog, error)
	Delete(container *PipelineLog) error
}

func newPipelineLogClient(apiClient *Client) *PipelineLogClient {
	return &PipelineLogClient{
		apiClient: apiClient,
	}
}

func (c *PipelineLogClient) Create(container *PipelineLog) (*PipelineLog, error) {
	resp := &PipelineLog{}
	err := c.apiClient.Ops.DoCreate(PipelineLogType, container, resp)
	return resp, err
}

func (c *PipelineLogClient) Update(existing *PipelineLog, updates interface{}) (*PipelineLog, error) {
	resp := &PipelineLog{}
	err := c.apiClient.Ops.DoUpdate(PipelineLogType, &existing.Resource, updates, resp)
	return resp, err
}

func (c *PipelineLogClient) List(opts *types.ListOpts) (*PipelineLogCollection, error) {
	resp := &PipelineLogCollection{}
	err := c.apiClient.Ops.DoList(PipelineLogType, opts, resp)
	resp.client = c
	return resp, err
}

func (cc *PipelineLogCollection) Next() (*PipelineLogCollection, error) {
	if cc != nil && cc.Pagination != nil && cc.Pagination.Next != "" {
		resp := &PipelineLogCollection{}
		err := cc.client.apiClient.Ops.DoNext(cc.Pagination.Next, resp)
		resp.client = cc.client
		return resp, err
	}
	return nil, nil
}

func (c *PipelineLogClient) ByID(id string) (*PipelineLog, error) {
	resp := &PipelineLog{}
	err := c.apiClient.Ops.DoByID(PipelineLogType, id, resp)
	return resp, err
}

func (c *PipelineLogClient) Delete(container *PipelineLog) error {
	return c.apiClient.Ops.DoResourceDelete(PipelineLogType, &container.Resource)
}
