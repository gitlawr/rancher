package client

const (
	BuildImageStepConfigType                = "buildImageStepConfig"
	BuildImageStepConfigFieldBuildPath      = "buildPath"
	BuildImageStepConfigFieldDockerfilePath = "dockerFilePath"
	BuildImageStepConfigFieldImageTag       = "imageTag"
)

type BuildImageStepConfig struct {
	BuildPath      string `json:"buildPath,omitempty"`
	DockerfilePath string `json:"dockerFilePath,omitempty"`
	ImageTag       string `json:"imageTag,omitempty"`
}
