package client

const (
	GitRepoCacheSpecType                   = "gitRepoCacheSpec"
	GitRepoCacheSpecFieldRemoteAccountName = "remoteAccountName"
	GitRepoCacheSpecFieldRepositories      = "repositories"
	GitRepoCacheSpecFieldType              = "type"
	GitRepoCacheSpecFieldUserId            = "userId"
)

type GitRepoCacheSpec struct {
	RemoteAccountName string          `json:"remoteAccountName,omitempty"`
	Repositories      []GitRepository `json:"repositories,omitempty"`
	Type              string          `json:"type,omitempty"`
	UserId            string          `json:"userId,omitempty"`
}
