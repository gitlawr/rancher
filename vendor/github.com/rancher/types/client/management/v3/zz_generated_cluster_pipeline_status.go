package client

const (
	ClusterPipelineStatusType            = "clusterPipelineStatus"
	ClusterPipelineStatusFieldConditions = "conditions"
)

type ClusterPipelineStatus struct {
	Conditions []PipelineCondition `json:"conditions,omitempty"`
}
