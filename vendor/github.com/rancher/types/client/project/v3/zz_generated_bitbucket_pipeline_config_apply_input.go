package client

const (
	BitbucketPipelineConfigApplyInputType                 = "bitbucketPipelineConfigApplyInput"
	BitbucketPipelineConfigApplyInputFieldBitbucketConfig = "bitbucketConfig"
	BitbucketPipelineConfigApplyInputFieldCode            = "code"
)

type BitbucketPipelineConfigApplyInput struct {
	BitbucketConfig *BitbucketPipelineConfig `json:"bitbucketConfig,omitempty" yaml:"bitbucketConfig,omitempty"`
	Code            string                   `json:"code,omitempty" yaml:"code,omitempty"`
}
