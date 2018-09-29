package client

import (
	"github.com/rancher/norman/types"
)

const (
	BitbucketCloudApplyInputType                 = "bitbucketCloudApplyInput"
	BitbucketCloudApplyInputFieldAnnotations     = "annotations"
	BitbucketCloudApplyInputFieldClientID        = "clientId"
	BitbucketCloudApplyInputFieldClientSecret    = "clientSecret"
	BitbucketCloudApplyInputFieldCode            = "code"
	BitbucketCloudApplyInputFieldCreated         = "created"
	BitbucketCloudApplyInputFieldCreatorID       = "creatorId"
	BitbucketCloudApplyInputFieldEnabled         = "enabled"
	BitbucketCloudApplyInputFieldLabels          = "labels"
	BitbucketCloudApplyInputFieldName            = "name"
	BitbucketCloudApplyInputFieldNamespaceId     = "namespaceId"
	BitbucketCloudApplyInputFieldOwnerReferences = "ownerReferences"
	BitbucketCloudApplyInputFieldProjectID       = "projectId"
	BitbucketCloudApplyInputFieldRedirectURL     = "redirectUrl"
	BitbucketCloudApplyInputFieldRemoved         = "removed"
	BitbucketCloudApplyInputFieldType            = "type"
	BitbucketCloudApplyInputFieldUUID            = "uuid"
)

type BitbucketCloudApplyInput struct {
	types.Resource
	Annotations     map[string]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`
	ClientID        string            `json:"clientId,omitempty" yaml:"clientId,omitempty"`
	ClientSecret    string            `json:"clientSecret,omitempty" yaml:"clientSecret,omitempty"`
	Code            string            `json:"code,omitempty" yaml:"code,omitempty"`
	Created         string            `json:"created,omitempty" yaml:"created,omitempty"`
	CreatorID       string            `json:"creatorId,omitempty" yaml:"creatorId,omitempty"`
	Enabled         bool              `json:"enabled,omitempty" yaml:"enabled,omitempty"`
	Labels          map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	Name            string            `json:"name,omitempty" yaml:"name,omitempty"`
	NamespaceId     string            `json:"namespaceId,omitempty" yaml:"namespaceId,omitempty"`
	OwnerReferences []OwnerReference  `json:"ownerReferences,omitempty" yaml:"ownerReferences,omitempty"`
	ProjectID       string            `json:"projectId,omitempty" yaml:"projectId,omitempty"`
	RedirectURL     string            `json:"redirectUrl,omitempty" yaml:"redirectUrl,omitempty"`
	Removed         string            `json:"removed,omitempty" yaml:"removed,omitempty"`
	Type            string            `json:"type,omitempty" yaml:"type,omitempty"`
	UUID            string            `json:"uuid,omitempty" yaml:"uuid,omitempty"`
}

type BitbucketCloudApplyInputCollection struct {
	types.Collection
	Data   []BitbucketCloudApplyInput `json:"data,omitempty"`
	client *BitbucketCloudApplyInputClient
}

type BitbucketCloudApplyInputClient struct {
	apiClient *Client
}

type BitbucketCloudApplyInputOperations interface {
	List(opts *types.ListOpts) (*BitbucketCloudApplyInputCollection, error)
	Create(opts *BitbucketCloudApplyInput) (*BitbucketCloudApplyInput, error)
	Update(existing *BitbucketCloudApplyInput, updates interface{}) (*BitbucketCloudApplyInput, error)
	Replace(existing *BitbucketCloudApplyInput) (*BitbucketCloudApplyInput, error)
	ByID(id string) (*BitbucketCloudApplyInput, error)
	Delete(container *BitbucketCloudApplyInput) error
}

func newBitbucketCloudApplyInputClient(apiClient *Client) *BitbucketCloudApplyInputClient {
	return &BitbucketCloudApplyInputClient{
		apiClient: apiClient,
	}
}

func (c *BitbucketCloudApplyInputClient) Create(container *BitbucketCloudApplyInput) (*BitbucketCloudApplyInput, error) {
	resp := &BitbucketCloudApplyInput{}
	err := c.apiClient.Ops.DoCreate(BitbucketCloudApplyInputType, container, resp)
	return resp, err
}

func (c *BitbucketCloudApplyInputClient) Update(existing *BitbucketCloudApplyInput, updates interface{}) (*BitbucketCloudApplyInput, error) {
	resp := &BitbucketCloudApplyInput{}
	err := c.apiClient.Ops.DoUpdate(BitbucketCloudApplyInputType, &existing.Resource, updates, resp)
	return resp, err
}

func (c *BitbucketCloudApplyInputClient) Replace(obj *BitbucketCloudApplyInput) (*BitbucketCloudApplyInput, error) {
	resp := &BitbucketCloudApplyInput{}
	err := c.apiClient.Ops.DoReplace(BitbucketCloudApplyInputType, &obj.Resource, obj, resp)
	return resp, err
}

func (c *BitbucketCloudApplyInputClient) List(opts *types.ListOpts) (*BitbucketCloudApplyInputCollection, error) {
	resp := &BitbucketCloudApplyInputCollection{}
	err := c.apiClient.Ops.DoList(BitbucketCloudApplyInputType, opts, resp)
	resp.client = c
	return resp, err
}

func (cc *BitbucketCloudApplyInputCollection) Next() (*BitbucketCloudApplyInputCollection, error) {
	if cc != nil && cc.Pagination != nil && cc.Pagination.Next != "" {
		resp := &BitbucketCloudApplyInputCollection{}
		err := cc.client.apiClient.Ops.DoNext(cc.Pagination.Next, resp)
		resp.client = cc.client
		return resp, err
	}
	return nil, nil
}

func (c *BitbucketCloudApplyInputClient) ByID(id string) (*BitbucketCloudApplyInput, error) {
	resp := &BitbucketCloudApplyInput{}
	err := c.apiClient.Ops.DoByID(BitbucketCloudApplyInputType, id, resp)
	return resp, err
}

func (c *BitbucketCloudApplyInputClient) Delete(container *BitbucketCloudApplyInput) error {
	return c.apiClient.Ops.DoResourceDelete(BitbucketCloudApplyInputType, &container.Resource)
}
