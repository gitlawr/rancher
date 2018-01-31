package client

const (
	StepStatusType           = "stepStatus"
	StepStatusFieldEndTime   = "endTime"
	StepStatusFieldStartTime = "startTime"
	StepStatusFieldState     = "state"
)

type StepStatus struct {
	EndTime   *int64 `json:"endTime,omitempty"`
	StartTime *int64 `json:"startTime,omitempty"`
	State     string `json:"state,omitempty"`
}
