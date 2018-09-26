package client

const (
	SourceCodeCredentialSpecType                = "sourceCodeCredentialSpec"
	SourceCodeCredentialSpecFieldAvatarURL      = "avatarUrl"
	SourceCodeCredentialSpecFieldDisplayName    = "displayName"
	SourceCodeCredentialSpecFieldExpiry         = "expiry"
	SourceCodeCredentialSpecFieldGitLoginName   = "gitLoginName"
	SourceCodeCredentialSpecFieldHTMLURL        = "htmlUrl"
	SourceCodeCredentialSpecFieldLoginName      = "loginName"
	SourceCodeCredentialSpecFieldProjectID      = "projectId"
	SourceCodeCredentialSpecFieldRefreshToken   = "accessToken"
	SourceCodeCredentialSpecFieldSourceCodeType = "sourceCodeType"
	SourceCodeCredentialSpecFieldUserID         = "userId"
)

type SourceCodeCredentialSpec struct {
	AvatarURL      string `json:"avatarUrl,omitempty" yaml:"avatarUrl,omitempty"`
	DisplayName    string `json:"displayName,omitempty" yaml:"displayName,omitempty"`
	Expiry         string `json:"expiry,omitempty" yaml:"expiry,omitempty"`
	GitLoginName   string `json:"gitLoginName,omitempty" yaml:"gitLoginName,omitempty"`
	HTMLURL        string `json:"htmlUrl,omitempty" yaml:"htmlUrl,omitempty"`
	LoginName      string `json:"loginName,omitempty" yaml:"loginName,omitempty"`
	ProjectID      string `json:"projectId,omitempty" yaml:"projectId,omitempty"`
	RefreshToken   string `json:"accessToken,omitempty" yaml:"accessToken,omitempty"`
	SourceCodeType string `json:"sourceCodeType,omitempty" yaml:"sourceCodeType,omitempty"`
	UserID         string `json:"userId,omitempty" yaml:"userId,omitempty"`
}
