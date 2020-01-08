package client

const (
	QueryGlobalGraphType           = "queryGlobalGraph"
	QueryGlobalGraphFieldGraphName = "graphID"
	QueryGlobalGraphFieldSeries    = "series"
)

type QueryGlobalGraph struct {
	GraphName string   `json:"graphID,omitempty" yaml:"graphID,omitempty"`
	Series    []string `json:"series,omitempty" yaml:"series,omitempty"`
}
