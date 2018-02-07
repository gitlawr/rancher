package client

const (
	AuthUserInputType             = "authUserInput"
	AuthUserInputFieldCode        = "code"
	AuthUserInputFieldRedirectUrl = "redirectUrl"
	AuthUserInputFieldType        = "type"
)

type AuthUserInput struct {
	Code        string `json:"code,omitempty"`
	RedirectUrl string `json:"redirectUrl,omitempty"`
	Type        string `json:"type,omitempty"`
}
