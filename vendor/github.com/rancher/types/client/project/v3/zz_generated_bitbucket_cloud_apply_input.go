package client

const (
	BitbucketCloudApplyInputType              = "bitbucketCloudApplyInput"
	BitbucketCloudApplyInputFieldClientID     = "clientId"
	BitbucketCloudApplyInputFieldClientSecret = "clientSecret"
	BitbucketCloudApplyInputFieldCode         = "code"
	BitbucketCloudApplyInputFieldRedirectURL  = "redirectUrl"
)

type BitbucketCloudApplyInput struct {
	ClientID     string `json:"clientId,omitempty" yaml:"clientId,omitempty"`
	ClientSecret string `json:"clientSecret,omitempty" yaml:"clientSecret,omitempty"`
	Code         string `json:"code,omitempty" yaml:"code,omitempty"`
	RedirectURL  string `json:"redirectUrl,omitempty" yaml:"redirectUrl,omitempty"`
}
