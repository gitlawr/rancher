package client

const (
	ReloadPipelineInputType        = "reloadPipelineInput"
	ReloadPipelineInputFieldBranch = "branch"
)

type ReloadPipelineInput struct {
	Branch string `json:"branch,omitempty" yaml:"branch,omitempty"`
}
