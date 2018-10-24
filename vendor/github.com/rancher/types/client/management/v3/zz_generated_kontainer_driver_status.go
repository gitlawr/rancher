package client

const (
	KontainerDriverStatusType                = "kontainerDriverStatus"
	KontainerDriverStatusFieldActualURL      = "actualUrl"
	KontainerDriverStatusFieldConditions     = "conditions"
	KontainerDriverStatusFieldExecutablePath = "executablePath"
)

type KontainerDriverStatus struct {
	ActualURL      string      `json:"actualUrl,omitempty" yaml:"actualUrl,omitempty"`
	Conditions     []Condition `json:"conditions,omitempty" yaml:"conditions,omitempty"`
	ExecutablePath string      `json:"executablePath,omitempty" yaml:"executablePath,omitempty"`
}
