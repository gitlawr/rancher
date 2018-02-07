package client

const (
	GithubConfigType              = "githubConfig"
	GithubConfigFieldClientId     = "clientId"
	GithubConfigFieldClientSecret = "clientSecret"
	GithubConfigFieldHost         = "host"
	GithubConfigFieldRedirectUrl  = "redirectUrl"
	GithubConfigFieldTLS          = "tls"
)

type GithubConfig struct {
	ClientId     string `json:"clientId,omitempty"`
	ClientSecret string `json:"clientSecret,omitempty"`
	Host         string `json:"host,omitempty"`
	RedirectUrl  string `json:"redirectUrl,omitempty"`
	TLS          *bool  `json:"tls,omitempty"`
}
