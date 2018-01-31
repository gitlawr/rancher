package client

const (
	CronTriggerType          = "cronTrigger"
	CronTriggerFieldActive   = "active"
	CronTriggerFieldSpec     = "spec"
	CronTriggerFieldTimezone = "timezone"
)

type CronTrigger struct {
	Active   *bool  `json:"active,omitempty"`
	Spec     string `json:"spec,omitempty"`
	Timezone string `json:"timezone,omitempty"`
}
