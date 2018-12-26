package client

const (
	WebhookTemplateConfigType              = "webhookTemplateConfig"
	WebhookTemplateConfigFieldTemplate     = "template"
	WebhookTemplateConfigFieldTemplateKind = "templateKind"
	WebhookTemplateConfigFieldURL          = "url"
)

type WebhookTemplateConfig struct {
	Template     string `json:"template,omitempty" yaml:"template,omitempty"`
	TemplateKind string `json:"templateKind,omitempty" yaml:"templateKind,omitempty"`
	URL          string `json:"url,omitempty" yaml:"url,omitempty"`
}
