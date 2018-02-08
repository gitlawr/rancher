package client

const (
	SourceCodeRepositorySpecType                          = "sourceCodeRepositorySpec"
	SourceCodeRepositorySpecFieldClusterId                = "clusterId"
	SourceCodeRepositorySpecFieldLanguage                 = "language"
	SourceCodeRepositorySpecFieldPermissions              = "permissions"
	SourceCodeRepositorySpecFieldSourceCodeCredentialName = "sourceCodeCredentialName"
	SourceCodeRepositorySpecFieldSourceCodeType           = "sourceCodeType"
	SourceCodeRepositorySpecFieldUrl                      = "url"
	SourceCodeRepositorySpecFieldUserId                   = "userId"
)

type SourceCodeRepositorySpec struct {
	ClusterId                string    `json:"clusterId,omitempty"`
	Language                 string    `json:"language,omitempty"`
	Permissions              *RepoPerm `json:"permissions,omitempty"`
	SourceCodeCredentialName string    `json:"sourceCodeCredentialName,omitempty"`
	SourceCodeType           string    `json:"sourceCodeType,omitempty"`
	Url                      string    `json:"url,omitempty"`
	UserId                   string    `json:"userId,omitempty"`
}
