package client

const (
	PipelineStatusType               = "pipelineStatus"
	PipelineStatusFieldLastRunId     = "lastRunId"
	PipelineStatusFieldLastRunState  = "lastRunState"
	PipelineStatusFieldLastRunTime   = "lastRunTime"
	PipelineStatusFieldNextRunNumber = "nextRunNumber"
	PipelineStatusFieldNextRunTime   = "nextRunTime"
	PipelineStatusFieldWebHookId     = "webhookId"
	PipelineStatusFieldWebHookToken  = "webhookToken"
)

type PipelineStatus struct {
	LastRunId     string `json:"lastRunId,omitempty"`
	LastRunState  string `json:"lastRunState,omitempty"`
	LastRunTime   *int64 `json:"lastRunTime,omitempty"`
	NextRunNumber *int64 `json:"nextRunNumber,omitempty"`
	NextRunTime   *int64 `json:"nextRunTime,omitempty"`
	WebHookId     string `json:"webhookId,omitempty"`
	WebHookToken  string `json:"webhookToken,omitempty"`
}
