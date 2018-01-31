package client

const (
	PipelineHistoryStatusType             = "pipelineHistoryStatus"
	PipelineHistoryStatusFieldCommitInfo  = "commitInfo"
	PipelineHistoryStatusFieldEndTime     = "endTime"
	PipelineHistoryStatusFieldEnvVars     = "envVars"
	PipelineHistoryStatusFieldStageStatus = "stageStatus"
	PipelineHistoryStatusFieldStartTime   = "startTime"
	PipelineHistoryStatusFieldState       = "state"
)

type PipelineHistoryStatus struct {
	CommitInfo  string            `json:"commitInfo,omitempty"`
	EndTime     *int64            `json:"endTime,omitempty"`
	EnvVars     map[string]string `json:"envVars,omitempty"`
	StageStatus []StageStatus     `json:"stageStatus,omitempty"`
	StartTime   *int64            `json:"startTime,omitempty"`
	State       string            `json:"state,omitempty"`
}
