package client

const (
	SourceCodeCredentialSpecType             = "sourceCodeCredentialSpec"
	SourceCodeCredentialSpecFieldAccessToken = "accessToken"
	SourceCodeCredentialSpecFieldAvatarURL   = "avatarUrl"
	SourceCodeCredentialSpecFieldClusterId   = "clusterId"
	SourceCodeCredentialSpecFieldDisplayName = "displayName"
	SourceCodeCredentialSpecFieldHTMLURL     = "htmlUrl"
	SourceCodeCredentialSpecFieldLoginName   = "loginName"
	SourceCodeCredentialSpecFieldType        = "type"
	SourceCodeCredentialSpecFieldUserId      = "UserId"
)

type SourceCodeCredentialSpec struct {
	AccessToken string `json:"accessToken,omitempty"`
	AvatarURL   string `json:"avatarUrl,omitempty"`
	ClusterId   string `json:"clusterId,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
	HTMLURL     string `json:"htmlUrl,omitempty"`
	LoginName   string `json:"loginName,omitempty"`
	Type        string `json:"type,omitempty"`
	UserId      string `json:"UserId,omitempty"`
}
