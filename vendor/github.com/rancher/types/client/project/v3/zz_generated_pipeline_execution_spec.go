package client

const (
	PipelineExecutionSpecType                = "pipelineExecutionSpec"
	PipelineExecutionSpecFieldAuthor         = "author"
	PipelineExecutionSpecFieldAvatarURL      = "avatarUrl"
	PipelineExecutionSpecFieldBranch         = "branch"
	PipelineExecutionSpecFieldCommit         = "commit"
	PipelineExecutionSpecFieldEmail          = "email"
	PipelineExecutionSpecFieldEvent          = "event"
	PipelineExecutionSpecFieldHTMLLink       = "htmlLink"
	PipelineExecutionSpecFieldMessage        = "message"
	PipelineExecutionSpecFieldPipelineConfig = "pipelineConfig"
	PipelineExecutionSpecFieldPipelineId     = "pipelineId"
	PipelineExecutionSpecFieldProjectId      = "projectId"
	PipelineExecutionSpecFieldRef            = "ref"
	PipelineExecutionSpecFieldRepositoryURL  = "repositoryUrl"
	PipelineExecutionSpecFieldRun            = "run"
	PipelineExecutionSpecFieldSharedEngine   = "sharedEngine"
	PipelineExecutionSpecFieldTitle          = "title"
	PipelineExecutionSpecFieldTriggerUserId  = "triggerUserId"
	PipelineExecutionSpecFieldTriggeredBy    = "triggeredBy"
)

type PipelineExecutionSpec struct {
	Author         string          `json:"author,omitempty" yaml:"author,omitempty"`
	AvatarURL      string          `json:"avatarUrl,omitempty" yaml:"avatarUrl,omitempty"`
	Branch         string          `json:"branch,omitempty" yaml:"branch,omitempty"`
	Commit         string          `json:"commit,omitempty" yaml:"commit,omitempty"`
	Email          string          `json:"email,omitempty" yaml:"email,omitempty"`
	Event          string          `json:"event,omitempty" yaml:"event,omitempty"`
	HTMLLink       string          `json:"htmlLink,omitempty" yaml:"htmlLink,omitempty"`
	Message        string          `json:"message,omitempty" yaml:"message,omitempty"`
	PipelineConfig *PipelineConfig `json:"pipelineConfig,omitempty" yaml:"pipelineConfig,omitempty"`
	PipelineId     string          `json:"pipelineId,omitempty" yaml:"pipelineId,omitempty"`
	ProjectId      string          `json:"projectId,omitempty" yaml:"projectId,omitempty"`
	Ref            string          `json:"ref,omitempty" yaml:"ref,omitempty"`
	RepositoryURL  string          `json:"repositoryUrl,omitempty" yaml:"repositoryUrl,omitempty"`
	Run            int64           `json:"run,omitempty" yaml:"run,omitempty"`
	SharedEngine   bool            `json:"sharedEngine,omitempty" yaml:"sharedEngine,omitempty"`
	Title          string          `json:"title,omitempty" yaml:"title,omitempty"`
	TriggerUserId  string          `json:"triggerUserId,omitempty" yaml:"triggerUserId,omitempty"`
	TriggeredBy    string          `json:"triggeredBy,omitempty" yaml:"triggeredBy,omitempty"`
}
