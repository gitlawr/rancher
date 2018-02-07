package client

const (
	StepType                        = "step"
	StepFieldPublishImageStepConfig = "publishImageConfig"
	StepFieldRunScriptStepConfig    = "runScriptConfig"
	StepFieldSourceCodeStepConfig   = "sourceCodeConfig"
	StepFieldTimeout                = "timeout"
)

type Step struct {
	PublishImageStepConfig *PublishImageStepConfig `json:"publishImageConfig,omitempty"`
	RunScriptStepConfig    *RunScriptStepConfig    `json:"runScriptConfig,omitempty"`
	SourceCodeStepConfig   *SourceCodeStepConfig   `json:"sourceCodeConfig,omitempty"`
	Timeout                *int64                  `json:"timeout,omitempty"`
}
