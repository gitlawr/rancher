package client

const (
	PipelineSpecType             = "pipelineSpec"
	PipelineSpecFieldActive      = "active"
	PipelineSpecFieldDisplayName = "displayName"
	PipelineSpecFieldProjectId   = "projectId"
	PipelineSpecFieldStages      = "stages"
	PipelineSpecFieldTriggers    = "triggers"
)

type PipelineSpec struct {
	Active      *bool     `json:"active,omitempty"`
	DisplayName string    `json:"displayName,omitempty"`
	ProjectId   string    `json:"projectId,omitempty"`
	Stages      []Stage   `json:"stages,omitempty"`
	Triggers    *Triggers `json:"triggers,omitempty"`
}
