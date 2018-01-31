package client

const (
	PipelineHistorySpecType             = "pipelineHistorySpec"
	PipelineHistorySpecFieldDisplayName = "displayName"
	PipelineHistorySpecFieldPipeline    = "pipeline"
	PipelineHistorySpecFieldProjectId   = "projectId"
	PipelineHistorySpecFieldRunNumber   = "runNumber"
	PipelineHistorySpecFieldTriggerType = "triggerType"
)

type PipelineHistorySpec struct {
	DisplayName string    `json:"displayName,omitempty"`
	Pipeline    *Pipeline `json:"pipeline,omitempty"`
	ProjectId   string    `json:"projectId,omitempty"`
	RunNumber   *int64    `json:"runNumber,omitempty"`
	TriggerType string    `json:"triggerType,omitempty"`
}
