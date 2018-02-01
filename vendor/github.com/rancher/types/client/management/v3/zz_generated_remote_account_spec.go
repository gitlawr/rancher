package client

const (
	RemoteAccountSpecType             = "remoteAccountSpec"
	RemoteAccountSpecFieldAccessToken = "accessToken"
	RemoteAccountSpecFieldAccountName = "accountName"
	RemoteAccountSpecFieldAvatarURL   = "avatarUrl"
	RemoteAccountSpecFieldHTMLURL     = "htmlUrl"
	RemoteAccountSpecFieldLogin       = "login"
	RemoteAccountSpecFieldUserID      = "userId"
)

type RemoteAccountSpec struct {
	AccessToken string `json:"accessToken,omitempty"`
	AccountName string `json:"accountName,omitempty"`
	AvatarURL   string `json:"avatarUrl,omitempty"`
	HTMLURL     string `json:"htmlUrl,omitempty"`
	Login       string `json:"login,omitempty"`
	UserID      string `json:"userId,omitempty"`
}
