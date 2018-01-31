package client

import (
	"github.com/rancher/norman/types"
)

const (
	RemoteAccountType                      = "remoteAccount"
	RemoteAccountFieldAccessToken          = "accessToken"
	RemoteAccountFieldAccountName          = "accountName"
	RemoteAccountFieldAnnotations          = "annotations"
	RemoteAccountFieldAvatarURL            = "avatarUrl"
	RemoteAccountFieldCreated              = "created"
	RemoteAccountFieldCreatorID            = "creatorId"
	RemoteAccountFieldDisplayName          = "displayName"
	RemoteAccountFieldHTMLURL              = "htmlUrl"
	RemoteAccountFieldLabels               = "labels"
	RemoteAccountFieldName                 = "name"
	RemoteAccountFieldNamespaceId          = "namespaceId"
	RemoteAccountFieldOwnerReferences      = "ownerReferences"
	RemoteAccountFieldRemoved              = "removed"
	RemoteAccountFieldState                = "state"
	RemoteAccountFieldStatus               = "status"
	RemoteAccountFieldTransitioning        = "transitioning"
	RemoteAccountFieldTransitioningMessage = "transitioningMessage"
	RemoteAccountFieldType                 = "type"
	RemoteAccountFieldUserId               = "userId"
	RemoteAccountFieldUuid                 = "uuid"
)

type RemoteAccount struct {
	types.Resource
	AccessToken          string               `json:"accessToken,omitempty"`
	AccountName          string               `json:"accountName,omitempty"`
	Annotations          map[string]string    `json:"annotations,omitempty"`
	AvatarURL            string               `json:"avatarUrl,omitempty"`
	Created              string               `json:"created,omitempty"`
	CreatorID            string               `json:"creatorId,omitempty"`
	DisplayName          string               `json:"displayName,omitempty"`
	HTMLURL              string               `json:"htmlUrl,omitempty"`
	Labels               map[string]string    `json:"labels,omitempty"`
	Name                 string               `json:"name,omitempty"`
	NamespaceId          string               `json:"namespaceId,omitempty"`
	OwnerReferences      []OwnerReference     `json:"ownerReferences,omitempty"`
	Removed              string               `json:"removed,omitempty"`
	State                string               `json:"state,omitempty"`
	Status               *RemoteAccountStatus `json:"status,omitempty"`
	Transitioning        string               `json:"transitioning,omitempty"`
	TransitioningMessage string               `json:"transitioningMessage,omitempty"`
	Type                 string               `json:"type,omitempty"`
	UserId               string               `json:"userId,omitempty"`
	Uuid                 string               `json:"uuid,omitempty"`
}
type RemoteAccountCollection struct {
	types.Collection
	Data   []RemoteAccount `json:"data,omitempty"`
	client *RemoteAccountClient
}

type RemoteAccountClient struct {
	apiClient *Client
}

type RemoteAccountOperations interface {
	List(opts *types.ListOpts) (*RemoteAccountCollection, error)
	Create(opts *RemoteAccount) (*RemoteAccount, error)
	Update(existing *RemoteAccount, updates interface{}) (*RemoteAccount, error)
	ByID(id string) (*RemoteAccount, error)
	Delete(container *RemoteAccount) error
}

func newRemoteAccountClient(apiClient *Client) *RemoteAccountClient {
	return &RemoteAccountClient{
		apiClient: apiClient,
	}
}

func (c *RemoteAccountClient) Create(container *RemoteAccount) (*RemoteAccount, error) {
	resp := &RemoteAccount{}
	err := c.apiClient.Ops.DoCreate(RemoteAccountType, container, resp)
	return resp, err
}

func (c *RemoteAccountClient) Update(existing *RemoteAccount, updates interface{}) (*RemoteAccount, error) {
	resp := &RemoteAccount{}
	err := c.apiClient.Ops.DoUpdate(RemoteAccountType, &existing.Resource, updates, resp)
	return resp, err
}

func (c *RemoteAccountClient) List(opts *types.ListOpts) (*RemoteAccountCollection, error) {
	resp := &RemoteAccountCollection{}
	err := c.apiClient.Ops.DoList(RemoteAccountType, opts, resp)
	resp.client = c
	return resp, err
}

func (cc *RemoteAccountCollection) Next() (*RemoteAccountCollection, error) {
	if cc != nil && cc.Pagination != nil && cc.Pagination.Next != "" {
		resp := &RemoteAccountCollection{}
		err := cc.client.apiClient.Ops.DoNext(cc.Pagination.Next, resp)
		resp.client = cc.client
		return resp, err
	}
	return nil, nil
}

func (c *RemoteAccountClient) ByID(id string) (*RemoteAccount, error) {
	resp := &RemoteAccount{}
	err := c.apiClient.Ops.DoByID(RemoteAccountType, id, resp)
	return resp, err
}

func (c *RemoteAccountClient) Delete(container *RemoteAccount) error {
	return c.apiClient.Ops.DoResourceDelete(RemoteAccountType, &container.Resource)
}
