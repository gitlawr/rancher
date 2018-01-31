package client

const (
	SourceCodeStepConfigType                 = "sourceCodeStepConfig"
	SourceCodeStepConfigFieldBranch          = "branch"
	SourceCodeStepConfigFieldRemoteAccountId = "remoteAccountId"
	SourceCodeStepConfigFieldRepository      = "repository"
)

type SourceCodeStepConfig struct {
	Branch          string `json:"branch,omitempty"`
	RemoteAccountId string `json:"remoteAccountId,omitempty"`
	Repository      string `json:"repository,omitempty"`
}
