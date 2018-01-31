package client

const (
	PushImageStepConfigType          = "pushImageStepConfig"
	PushImageStepConfigFieldImageTag = "imageTag"
)

type PushImageStepConfig struct {
	ImageTag string `json:"imageTag,omitempty"`
}
