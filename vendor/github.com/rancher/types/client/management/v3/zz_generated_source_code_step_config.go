package client

const (
	SourceCodeStepConfigType                        = "sourceCodeStepConfig"
	SourceCodeStepConfigFieldBranch                 = "branch"
	SourceCodeStepConfigFieldSourceCodeCredentialId = "sourceCodeCredentialId"
	SourceCodeStepConfigFieldUrl                    = "url"
)

type SourceCodeStepConfig struct {
	Branch                 string `json:"branch,omitempty"`
	SourceCodeCredentialId string `json:"sourceCodeCredentialId,omitempty"`
	Url                    string `json:"url,omitempty"`
}
