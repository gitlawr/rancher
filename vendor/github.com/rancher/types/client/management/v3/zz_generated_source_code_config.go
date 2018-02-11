package client

const (
	SourceCodeConfigType                        = "sourceCodeConfig"
	SourceCodeConfigFieldBranch                 = "branch"
	SourceCodeConfigFieldBranchCondition        = "branchCondition"
	SourceCodeConfigFieldSourceCodeCredentialId = "sourceCodeCredentialId"
	SourceCodeConfigFieldUrl                    = "url"
)

type SourceCodeConfig struct {
	Branch                 string `json:"branch,omitempty"`
	BranchCondition        string `json:"branchCondition,omitempty"`
	SourceCodeCredentialId string `json:"sourceCodeCredentialId,omitempty"`
	Url                    string `json:"url,omitempty"`
}
