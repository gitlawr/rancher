package client

const (
	TriggersType                = "triggers"
	TriggersFieldCronTrigger    = "cronTrigger"
	TriggersFieldWebhookTrigger = "webhookTrigger"
)

type Triggers struct {
	CronTrigger    *CronTrigger    `json:"cronTrigger,omitempty"`
	WebhookTrigger *WebhookTrigger `json:"webhookTrigger,omitempty"`
}
