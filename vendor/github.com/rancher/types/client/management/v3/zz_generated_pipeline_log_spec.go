package client

const (
	PipelineLogSpecType                     = "pipelineLogSpec"
	PipelineLogSpecFieldMessage             = "message"
	PipelineLogSpecFieldPipelineHistoryName = "pipelineHistoryName"
	PipelineLogSpecFieldProjectId           = "projectId"
	PipelineLogSpecFieldStageOrdinal        = "stageOrdinal"
	PipelineLogSpecFieldStepOrdinal         = "stepOrdinal"
)

type PipelineLogSpec struct {
	Message             string `json:"message,omitempty"`
	PipelineHistoryName string `json:"pipelineHistoryName,omitempty"`
	ProjectId           string `json:"projectId,omitempty"`
	StageOrdinal        *int64 `json:"stageOrdinal,omitempty"`
	StepOrdinal         *int64 `json:"stepOrdinal,omitempty"`
}
