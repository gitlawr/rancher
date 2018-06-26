package client

const (
	PipelineConfigType        = "pipelineConfig"
	PipelineConfigFieldBranch = "branch"
	PipelineConfigFieldStages = "stages"
)

type PipelineConfig struct {
	Branch *Constraint `json:"branch,omitempty" yaml:"branch,omitempty"`
	Stages []Stage     `json:"stages,omitempty" yaml:"stages,omitempty"`
}
