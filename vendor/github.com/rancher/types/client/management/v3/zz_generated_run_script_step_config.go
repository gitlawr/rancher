package client

const (
	RunScriptStepConfigType             = "runScriptStepConfig"
	RunScriptStepConfigFieldArgs        = "args"
	RunScriptStepConfigFieldEntrypoint  = "entrypoint"
	RunScriptStepConfigFieldEnv         = "env"
	RunScriptStepConfigFieldImage       = "image"
	RunScriptStepConfigFieldShellScript = "shellScript"
)

type RunScriptStepConfig struct {
	Args        string   `json:"args,omitempty"`
	Entrypoint  string   `json:"entrypoint,omitempty"`
	Env         []string `json:"env,omitempty"`
	Image       string   `json:"image,omitempty"`
	ShellScript string   `json:"shellScript,omitempty"`
}
