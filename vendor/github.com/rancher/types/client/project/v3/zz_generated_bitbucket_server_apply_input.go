package client

const (
	BitbucketServerApplyInputType             = "bitbucketServerApplyInput"
	BitbucketServerApplyInputFieldHostname    = "hostname"
	BitbucketServerApplyInputFieldPassword    = "password"
	BitbucketServerApplyInputFieldRedirectURL = "redirectUrl"
	BitbucketServerApplyInputFieldTLS         = "tls"
	BitbucketServerApplyInputFieldUserName    = "username"
)

type BitbucketServerApplyInput struct {
	Hostname    string `json:"hostname,omitempty" yaml:"hostname,omitempty"`
	Password    string `json:"password,omitempty" yaml:"password,omitempty"`
	RedirectURL string `json:"redirectUrl,omitempty" yaml:"redirectUrl,omitempty"`
	TLS         bool   `json:"tls,omitempty" yaml:"tls,omitempty"`
	UserName    string `json:"username,omitempty" yaml:"username,omitempty"`
}
