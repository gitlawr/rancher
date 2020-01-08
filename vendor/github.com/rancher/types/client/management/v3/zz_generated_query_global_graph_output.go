package client

const (
	QueryGlobalGraphOutputType      = "queryGlobalGraphOutput"
	QueryGlobalGraphOutputFieldData = "data"
	QueryGlobalGraphOutputFieldType = "type"
)

type QueryGlobalGraphOutput struct {
	Data []QueryGlobalGraph `json:"data,omitempty" yaml:"data,omitempty"`
	Type string             `json:"type,omitempty" yaml:"type,omitempty"`
}
