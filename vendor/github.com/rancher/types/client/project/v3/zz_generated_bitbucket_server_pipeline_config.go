package client

const (
	BitbucketServerPipelineConfigType                 = "bitbucketServerPipelineConfig"
	BitbucketServerPipelineConfigFieldAnnotations     = "annotations"
	BitbucketServerPipelineConfigFieldClientID        = "clientId"
	BitbucketServerPipelineConfigFieldClientSecret    = "clientSecret"
	BitbucketServerPipelineConfigFieldCreated         = "created"
	BitbucketServerPipelineConfigFieldCreatorID       = "creatorId"
	BitbucketServerPipelineConfigFieldEnabled         = "enabled"
	BitbucketServerPipelineConfigFieldHostname        = "hostname"
	BitbucketServerPipelineConfigFieldLabels          = "labels"
	BitbucketServerPipelineConfigFieldName            = "name"
	BitbucketServerPipelineConfigFieldNamespaceId     = "namespaceId"
	BitbucketServerPipelineConfigFieldOwnerReferences = "ownerReferences"
	BitbucketServerPipelineConfigFieldPassword        = "password"
	BitbucketServerPipelineConfigFieldProjectID       = "projectId"
	BitbucketServerPipelineConfigFieldRedirectURL     = "redirectUrl"
	BitbucketServerPipelineConfigFieldRemoved         = "removed"
	BitbucketServerPipelineConfigFieldTLS             = "tls"
	BitbucketServerPipelineConfigFieldType            = "type"
	BitbucketServerPipelineConfigFieldUUID            = "uuid"
	BitbucketServerPipelineConfigFieldUserName        = "username"
)

type BitbucketServerPipelineConfig struct {
	Annotations     map[string]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`
	ClientID        string            `json:"clientId,omitempty" yaml:"clientId,omitempty"`
	ClientSecret    string            `json:"clientSecret,omitempty" yaml:"clientSecret,omitempty"`
	Created         string            `json:"created,omitempty" yaml:"created,omitempty"`
	CreatorID       string            `json:"creatorId,omitempty" yaml:"creatorId,omitempty"`
	Enabled         bool              `json:"enabled,omitempty" yaml:"enabled,omitempty"`
	Hostname        string            `json:"hostname,omitempty" yaml:"hostname,omitempty"`
	Labels          map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	Name            string            `json:"name,omitempty" yaml:"name,omitempty"`
	NamespaceId     string            `json:"namespaceId,omitempty" yaml:"namespaceId,omitempty"`
	OwnerReferences []OwnerReference  `json:"ownerReferences,omitempty" yaml:"ownerReferences,omitempty"`
	Password        string            `json:"password,omitempty" yaml:"password,omitempty"`
	ProjectID       string            `json:"projectId,omitempty" yaml:"projectId,omitempty"`
	RedirectURL     string            `json:"redirectUrl,omitempty" yaml:"redirectUrl,omitempty"`
	Removed         string            `json:"removed,omitempty" yaml:"removed,omitempty"`
	TLS             bool              `json:"tls,omitempty" yaml:"tls,omitempty"`
	Type            string            `json:"type,omitempty" yaml:"type,omitempty"`
	UUID            string            `json:"uuid,omitempty" yaml:"uuid,omitempty"`
	UserName        string            `json:"username,omitempty" yaml:"username,omitempty"`
}
