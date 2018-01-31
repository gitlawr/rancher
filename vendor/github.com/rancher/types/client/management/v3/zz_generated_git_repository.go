package client

const (
	GitRepositoryType             = "gitRepository"
	GitRepositoryFieldCloneURL    = "cloneUrl"
	GitRepositoryFieldLanguage    = "language"
	GitRepositoryFieldName        = "name"
	GitRepositoryFieldPermissions = "permissions"
)

type GitRepository struct {
	CloneURL    string    `json:"cloneUrl,omitempty"`
	Language    string    `json:"language,omitempty"`
	Name        string    `json:"name,omitempty"`
	Permissions *RepoPerm `json:"permissions,omitempty"`
}
