package client

const (
	BitbucketServerApplyInputType                 = "bitbucketServerApplyInput"
	BitbucketServerApplyInputFieldBitbucketConfig = "bitbucketConfig"
	BitbucketServerApplyInputFieldCode            = "code"
)

type BitbucketServerApplyInput struct {
	BitbucketConfig *BitbucketServerPipelineConfig `json:"bitbucketConfig,omitempty" yaml:"bitbucketConfig,omitempty"`
	Code            string                         `json:"code,omitempty" yaml:"code,omitempty"`
}
