package client

const (
	PublishImageStepConfigType                = "publishImageStepConfig"
	PublishImageStepConfigFieldBuildContext   = "buildContext"
	PublishImageStepConfigFieldDockerfilePath = "dockerFilePath"
	PublishImageStepConfigFieldImageTag       = "imageTag"
)

type PublishImageStepConfig struct {
	BuildContext   string `json:"buildContext,omitempty"`
	DockerfilePath string `json:"dockerFilePath,omitempty"`
	ImageTag       string `json:"imageTag,omitempty"`
}
