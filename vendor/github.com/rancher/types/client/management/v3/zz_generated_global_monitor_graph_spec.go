package client

const (
	GlobalMonitorGraphSpecType                        = "globalMonitorGraphSpec"
	GlobalMonitorGraphSpecFieldDescription            = "description"
	GlobalMonitorGraphSpecFieldDetailsMetricsSelector = "detailsMetricsSelector"
	GlobalMonitorGraphSpecFieldDisplayResourceType    = "displayResourceType"
	GlobalMonitorGraphSpecFieldGraphType              = "graphType"
	GlobalMonitorGraphSpecFieldMetricsSelector        = "metricsSelector"
	GlobalMonitorGraphSpecFieldPriority               = "priority"
	GlobalMonitorGraphSpecFieldResourceType           = "resourceType"
	GlobalMonitorGraphSpecFieldYAxis                  = "yAxis"
)

type GlobalMonitorGraphSpec struct {
	Description            string            `json:"description,omitempty" yaml:"description,omitempty"`
	DetailsMetricsSelector map[string]string `json:"detailsMetricsSelector,omitempty" yaml:"detailsMetricsSelector,omitempty"`
	DisplayResourceType    string            `json:"displayResourceType,omitempty" yaml:"displayResourceType,omitempty"`
	GraphType              string            `json:"graphType,omitempty" yaml:"graphType,omitempty"`
	MetricsSelector        map[string]string `json:"metricsSelector,omitempty" yaml:"metricsSelector,omitempty"`
	Priority               int64             `json:"priority,omitempty" yaml:"priority,omitempty"`
	ResourceType           string            `json:"resourceType,omitempty" yaml:"resourceType,omitempty"`
	YAxis                  *YAxis            `json:"yAxis,omitempty" yaml:"yAxis,omitempty"`
}
