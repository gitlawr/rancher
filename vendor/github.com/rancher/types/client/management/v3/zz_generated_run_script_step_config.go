package client

const (
	RunScriptStepConfigType             = "runScriptStepConfig"
	RunScriptStepConfigFieldCommand     = "command"
	RunScriptStepConfigFieldEntrypoint  = "entrypoint"
	RunScriptStepConfigFieldEnv         = "env"
	RunScriptStepConfigFieldImage       = "image"
	RunScriptStepConfigFieldIsShell     = "isShell"
	RunScriptStepConfigFieldShellScript = "shellScript"
)

type RunScriptStepConfig struct {
	Command     string   `json:"command,omitempty"`
	Entrypoint  string   `json:"entrypoint,omitempty"`
	Env         []string `json:"env,omitempty"`
	Image       string   `json:"image,omitempty"`
	IsShell     *bool    `json:"isShell,omitempty"`
	ShellScript string   `json:"shellScript,omitempty"`
}
