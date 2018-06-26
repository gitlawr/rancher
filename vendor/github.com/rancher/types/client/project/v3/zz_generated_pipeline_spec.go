package client

const (
	PipelineSpecType                        = "pipelineSpec"
	PipelineSpecFieldDisplayName            = "displayName"
	PipelineSpecFieldProjectId              = "projectId"
	PipelineSpecFieldRepositoryURL          = "repositoryUrl"
	PipelineSpecFieldSourceCodeCredentialId = "sourceCodeCredentialId"
	PipelineSpecFieldTriggerWebhookPr       = "triggerWebhookPr"
	PipelineSpecFieldTriggerWebhookPush     = "triggerWebhookPush"
	PipelineSpecFieldTriggerWebhookTag      = "triggerWebhookTag"
	PipelineSpecFieldUnSyncConfigs          = "unSyncConfigs"
)

type PipelineSpec struct {
	DisplayName            string                    `json:"displayName,omitempty" yaml:"displayName,omitempty"`
	ProjectId              string                    `json:"projectId,omitempty" yaml:"projectId,omitempty"`
	RepositoryURL          string                    `json:"repositoryUrl,omitempty" yaml:"repositoryUrl,omitempty"`
	SourceCodeCredentialId string                    `json:"sourceCodeCredentialId,omitempty" yaml:"sourceCodeCredentialId,omitempty"`
	TriggerWebhookPr       bool                      `json:"triggerWebhookPr,omitempty" yaml:"triggerWebhookPr,omitempty"`
	TriggerWebhookPush     bool                      `json:"triggerWebhookPush,omitempty" yaml:"triggerWebhookPush,omitempty"`
	TriggerWebhookTag      bool                      `json:"triggerWebhookTag,omitempty" yaml:"triggerWebhookTag,omitempty"`
	UnSyncConfigs          map[string]PipelineConfig `json:"unSyncConfigs,omitempty" yaml:"unSyncConfigs,omitempty"`
}
