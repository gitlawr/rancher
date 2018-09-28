package client

const (
	BitbucketServerLoginInputType               = "bitbucketServerLoginInput"
	BitbucketServerLoginInputFieldOAuthToken    = "oauthToken"
	BitbucketServerLoginInputFieldOAuthVerifier = "oauthVerifier"
)

type BitbucketServerLoginInput struct {
	OAuthToken    string `json:"oauthToken,omitempty" yaml:"oauthToken,omitempty"`
	OAuthVerifier string `json:"oauthVerifier,omitempty" yaml:"oauthVerifier,omitempty"`
}
