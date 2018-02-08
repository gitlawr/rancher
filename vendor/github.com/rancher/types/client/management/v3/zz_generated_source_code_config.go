package client

const (
	SourceCodeConfigType                        = "sourceCodeConfig"
	SourceCodeConfigFieldBranch                 = "branch"
	SourceCodeConfigFieldSourceCodeCredentialId = "sourceCodeCredentialId"
	SourceCodeConfigFieldUrl                    = "url"
)

type SourceCodeConfig struct {
	Branch                 string `json:"branch,omitempty"`
	SourceCodeCredentialId string `json:"sourceCodeCredentialId,omitempty"`
	Url                    string `json:"url,omitempty"`
}
