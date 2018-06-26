package client

const (
	SourceCodeRepositorySpecType                        = "sourceCodeRepositorySpec"
	SourceCodeRepositorySpecFieldDefaultBranch          = "defaultBranch"
	SourceCodeRepositorySpecFieldLanguage               = "language"
	SourceCodeRepositorySpecFieldPermissions            = "permissions"
	SourceCodeRepositorySpecFieldProjectId              = "projectId"
	SourceCodeRepositorySpecFieldSourceCodeCredentialId = "sourceCodeCredentialId"
	SourceCodeRepositorySpecFieldSourceCodeType         = "sourceCodeType"
	SourceCodeRepositorySpecFieldURL                    = "url"
	SourceCodeRepositorySpecFieldUserId                 = "userId"
)

type SourceCodeRepositorySpec struct {
	DefaultBranch          string    `json:"defaultBranch,omitempty" yaml:"defaultBranch,omitempty"`
	Language               string    `json:"language,omitempty" yaml:"language,omitempty"`
	Permissions            *RepoPerm `json:"permissions,omitempty" yaml:"permissions,omitempty"`
	ProjectId              string    `json:"projectId,omitempty" yaml:"projectId,omitempty"`
	SourceCodeCredentialId string    `json:"sourceCodeCredentialId,omitempty" yaml:"sourceCodeCredentialId,omitempty"`
	SourceCodeType         string    `json:"sourceCodeType,omitempty" yaml:"sourceCodeType,omitempty"`
	URL                    string    `json:"url,omitempty" yaml:"url,omitempty"`
	UserId                 string    `json:"userId,omitempty" yaml:"userId,omitempty"`
}
