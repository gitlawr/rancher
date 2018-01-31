package client

const (
	PipelineConditionType                    = "pipelineCondition"
	PipelineConditionFieldLastTransitionTime = "lastTransitionTime"
	PipelineConditionFieldLastUpdateTime     = "lastUpdateTime"
	PipelineConditionFieldMessage            = "message"
	PipelineConditionFieldReason             = "reason"
	PipelineConditionFieldStatus             = "status"
	PipelineConditionFieldType               = "type"
)

type PipelineCondition struct {
	LastTransitionTime string `json:"lastTransitionTime,omitempty"`
	LastUpdateTime     string `json:"lastUpdateTime,omitempty"`
	Message            string `json:"message,omitempty"`
	Reason             string `json:"reason,omitempty"`
	Status             string `json:"status,omitempty"`
	Type               string `json:"type,omitempty"`
}
