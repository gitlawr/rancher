package client

const (
	SourceCodeStepConfigType                   = "sourceCodeStepConfig"
	SourceCodeStepConfigFieldBranch            = "branch"
	SourceCodeStepConfigFieldRemoteAccountName = "remoteAccountName"
	SourceCodeStepConfigFieldRepository        = "repository"
)

type SourceCodeStepConfig struct {
	Branch            string `json:"branch,omitempty"`
	RemoteAccountName string `json:"remoteAccountName,omitempty"`
	Repository        string `json:"repository,omitempty"`
}
