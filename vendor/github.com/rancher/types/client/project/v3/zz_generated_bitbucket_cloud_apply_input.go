package client

const (
	BitbucketCloudApplyInputType                 = "bitbucketCloudApplyInput"
	BitbucketCloudApplyInputFieldBitbucketConfig = "bitbucketConfig"
	BitbucketCloudApplyInputFieldCode            = "code"
)

type BitbucketCloudApplyInput struct {
	BitbucketConfig *BitbucketCloudPipelineConfig `json:"bitbucketConfig,omitempty" yaml:"bitbucketConfig,omitempty"`
	Code            string                        `json:"code,omitempty" yaml:"code,omitempty"`
}
