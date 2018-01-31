package client

const (
	StageStatusType            = "stageStatus"
	StageStatusFieldEndTime    = "endTime"
	StageStatusFieldStartTime  = "startTime"
	StageStatusFieldState      = "state"
	StageStatusFieldStepStatus = "stepStatus"
)

type StageStatus struct {
	EndTime    *int64       `json:"endTime,omitempty"`
	StartTime  *int64       `json:"startTime,omitempty"`
	State      string       `json:"state,omitempty"`
	StepStatus []StepStatus `json:"stepStatus,omitempty"`
}
