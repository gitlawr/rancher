package client

const (
	WebhookTriggerType        = "webhookTrigger"
	WebhookTriggerFieldActive = "active"
)

type WebhookTrigger struct {
	Active *bool `json:"active,omitempty"`
}
