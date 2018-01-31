package client

import (
	"github.com/rancher/norman/types"
)

const (
	GitRepoCacheType                      = "gitRepoCache"
	GitRepoCacheFieldAnnotations          = "annotations"
	GitRepoCacheFieldCreated              = "created"
	GitRepoCacheFieldCreatorID            = "creatorId"
	GitRepoCacheFieldLabels               = "labels"
	GitRepoCacheFieldName                 = "name"
	GitRepoCacheFieldNamespaceId          = "namespaceId"
	GitRepoCacheFieldOwnerReferences      = "ownerReferences"
	GitRepoCacheFieldRemoteAccountName    = "remoteAccountName"
	GitRepoCacheFieldRemoved              = "removed"
	GitRepoCacheFieldRepositories         = "repositories"
	GitRepoCacheFieldState                = "state"
	GitRepoCacheFieldStatus               = "status"
	GitRepoCacheFieldTransitioning        = "transitioning"
	GitRepoCacheFieldTransitioningMessage = "transitioningMessage"
	GitRepoCacheFieldType                 = "type"
	GitRepoCacheFieldUserId               = "userId"
	GitRepoCacheFieldUuid                 = "uuid"
)

type GitRepoCache struct {
	types.Resource
	Annotations          map[string]string   `json:"annotations,omitempty"`
	Created              string              `json:"created,omitempty"`
	CreatorID            string              `json:"creatorId,omitempty"`
	Labels               map[string]string   `json:"labels,omitempty"`
	Name                 string              `json:"name,omitempty"`
	NamespaceId          string              `json:"namespaceId,omitempty"`
	OwnerReferences      []OwnerReference    `json:"ownerReferences,omitempty"`
	RemoteAccountName    string              `json:"remoteAccountName,omitempty"`
	Removed              string              `json:"removed,omitempty"`
	Repositories         []GitRepository     `json:"repositories,omitempty"`
	State                string              `json:"state,omitempty"`
	Status               *GitRepoCacheStatus `json:"status,omitempty"`
	Transitioning        string              `json:"transitioning,omitempty"`
	TransitioningMessage string              `json:"transitioningMessage,omitempty"`
	Type                 string              `json:"type,omitempty"`
	UserId               string              `json:"userId,omitempty"`
	Uuid                 string              `json:"uuid,omitempty"`
}
type GitRepoCacheCollection struct {
	types.Collection
	Data   []GitRepoCache `json:"data,omitempty"`
	client *GitRepoCacheClient
}

type GitRepoCacheClient struct {
	apiClient *Client
}

type GitRepoCacheOperations interface {
	List(opts *types.ListOpts) (*GitRepoCacheCollection, error)
	Create(opts *GitRepoCache) (*GitRepoCache, error)
	Update(existing *GitRepoCache, updates interface{}) (*GitRepoCache, error)
	ByID(id string) (*GitRepoCache, error)
	Delete(container *GitRepoCache) error
}

func newGitRepoCacheClient(apiClient *Client) *GitRepoCacheClient {
	return &GitRepoCacheClient{
		apiClient: apiClient,
	}
}

func (c *GitRepoCacheClient) Create(container *GitRepoCache) (*GitRepoCache, error) {
	resp := &GitRepoCache{}
	err := c.apiClient.Ops.DoCreate(GitRepoCacheType, container, resp)
	return resp, err
}

func (c *GitRepoCacheClient) Update(existing *GitRepoCache, updates interface{}) (*GitRepoCache, error) {
	resp := &GitRepoCache{}
	err := c.apiClient.Ops.DoUpdate(GitRepoCacheType, &existing.Resource, updates, resp)
	return resp, err
}

func (c *GitRepoCacheClient) List(opts *types.ListOpts) (*GitRepoCacheCollection, error) {
	resp := &GitRepoCacheCollection{}
	err := c.apiClient.Ops.DoList(GitRepoCacheType, opts, resp)
	resp.client = c
	return resp, err
}

func (cc *GitRepoCacheCollection) Next() (*GitRepoCacheCollection, error) {
	if cc != nil && cc.Pagination != nil && cc.Pagination.Next != "" {
		resp := &GitRepoCacheCollection{}
		err := cc.client.apiClient.Ops.DoNext(cc.Pagination.Next, resp)
		resp.client = cc.client
		return resp, err
	}
	return nil, nil
}

func (c *GitRepoCacheClient) ByID(id string) (*GitRepoCache, error) {
	resp := &GitRepoCache{}
	err := c.apiClient.Ops.DoByID(GitRepoCacheType, id, resp)
	return resp, err
}

func (c *GitRepoCacheClient) Delete(container *GitRepoCache) error {
	return c.apiClient.Ops.DoResourceDelete(GitRepoCacheType, &container.Resource)
}
