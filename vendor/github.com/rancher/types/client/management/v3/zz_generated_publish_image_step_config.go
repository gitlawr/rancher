package client

const (
	PublishImageStepConfigType                = "publishImageStepConfig"
	PublishImageStepConfigFieldBuildContext   = "buildContext"
	PublishImageStepConfigFieldDockerfilePath = "dockerfilePath"
	PublishImageStepConfigFieldTag            = "tag"
)

type PublishImageStepConfig struct {
	BuildContext   string `json:"buildContext,omitempty"`
	DockerfilePath string `json:"dockerfilePath,omitempty"`
	Tag            string `json:"tag,omitempty"`
}
