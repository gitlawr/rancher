package client

const (
	PipelineHistorySpecType             = "pipelineHistorySpec"
	PipelineHistorySpecFieldPipeline    = "pipeline"
	PipelineHistorySpecFieldProjectId   = "projectId"
	PipelineHistorySpecFieldRunNumber   = "runNumber"
	PipelineHistorySpecFieldTriggerType = "triggerType"
)

type PipelineHistorySpec struct {
	Pipeline    *Pipeline `json:"pipeline,omitempty"`
	ProjectId   string    `json:"projectId,omitempty"`
	RunNumber   *int64    `json:"runNumber,omitempty"`
	TriggerType string    `json:"triggerType,omitempty"`
}
