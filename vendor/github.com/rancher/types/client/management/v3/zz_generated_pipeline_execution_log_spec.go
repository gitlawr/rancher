package client

const (
	PipelineExecutionLogSpecType                       = "pipelineExecutionLogSpec"
	PipelineExecutionLogSpecFieldMessage               = "message"
	PipelineExecutionLogSpecFieldPipelineExecutionName = "pipelineExecutionName"
	PipelineExecutionLogSpecFieldProjectId             = "projectId"
	PipelineExecutionLogSpecFieldStage                 = "stage"
	PipelineExecutionLogSpecFieldStep                  = "step"
)

type PipelineExecutionLogSpec struct {
	Message               string `json:"message,omitempty"`
	PipelineExecutionName string `json:"pipelineExecutionName,omitempty"`
	ProjectId             string `json:"projectId,omitempty"`
	Stage                 *int64 `json:"stage,omitempty"`
	Step                  *int64 `json:"step,omitempty"`
}
