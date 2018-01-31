package client

const (
	StepType                      = "step"
	StepFieldBuildImageStepConfig = "buildImageStepConfig"
	StepFieldPushImageStepConfig  = "pushImageStepConfig"
	StepFieldRunScriptStepConfig  = "runScriptStepConfig"
	StepFieldSourceCodeStepConfig = "sourceCodeStepConfig"
	StepFieldTimeout              = "timeout"
	StepFieldType                 = "type"
)

type Step struct {
	BuildImageStepConfig *BuildImageStepConfig `json:"buildImageStepConfig,omitempty"`
	PushImageStepConfig  *PushImageStepConfig  `json:"pushImageStepConfig,omitempty"`
	RunScriptStepConfig  *RunScriptStepConfig  `json:"runScriptStepConfig,omitempty"`
	SourceCodeStepConfig *SourceCodeStepConfig `json:"sourceCodeStepConfig,omitempty"`
	Timeout              *int64                `json:"timeout,omitempty"`
	Type                 string                `json:"type,omitempty"`
}
