package client

const (
	KontainerDriverSpecType               = "kontainerDriverSpec"
	KontainerDriverSpecFieldActive        = "active"
	KontainerDriverSpecFieldBuiltIn       = "builtIn"
	KontainerDriverSpecFieldChecksum      = "checksum"
	KontainerDriverSpecFieldDesiredURL    = "desirdUrl"
	KontainerDriverSpecFieldDisplayName   = "displayName"
	KontainerDriverSpecFieldDynamicSchema = "dynamicSchema"
)

type KontainerDriverSpec struct {
	Active        bool   `json:"active,omitempty" yaml:"active,omitempty"`
	BuiltIn       bool   `json:"builtIn,omitempty" yaml:"builtIn,omitempty"`
	Checksum      string `json:"checksum,omitempty" yaml:"checksum,omitempty"`
	DesiredURL    string `json:"desirdUrl,omitempty" yaml:"desirdUrl,omitempty"`
	DisplayName   string `json:"displayName,omitempty" yaml:"displayName,omitempty"`
	DynamicSchema bool   `json:"dynamicSchema,omitempty" yaml:"dynamicSchema,omitempty"`
}
