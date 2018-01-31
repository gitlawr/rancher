package client

const (
	GithubConfigType              = "githubConfig"
	GithubConfigFieldClientId     = "clientId"
	GithubConfigFieldClientSecret = "clientSecret"
	GithubConfigFieldHost         = "host"
	GithubConfigFieldScheme       = "githubConfig"
)

type GithubConfig struct {
	ClientId     string `json:"clientId,omitempty"`
	ClientSecret string `json:"clientSecret,omitempty"`
	Host         string `json:"host,omitempty"`
	Scheme       string `json:"githubConfig,omitempty"`
}
