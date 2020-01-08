package globalmonitoring

import (
	"encoding/json"
	"strings"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"

	managementv3 "github.com/rancher/types/apis/management.cattle.io/v3"
	"sigs.k8s.io/yaml"
)

type deployer struct {
	graphClient  managementv3.GlobalMonitorGraphInterface
	graphLister  managementv3.GlobalMonitorGraphLister
	metricLister managementv3.MonitorMetricLister
	metricClient managementv3.MonitorMetricInterface
}

func (d deployer) deploy() error {
	graphs, err := d.graphLister.List(globalDataNamespace, labels.NewSelector())
	if err != nil {
		return err
	}
	if len(graphs) != len(preDefinedClusterGraph) {
		for _, graph := range preDefinedClusterGraph {
			graph.Namespace = globalDataNamespace
			_, err := d.graphClient.Create(graph)
			if err != nil && !errors.IsAlreadyExists(err) {
				return err
			}
		}
	}
	metrics, err := d.metricLister.List(globalDataNamespace, labels.NewSelector())
	if err != nil {
		return err
	}
	if len(metrics) != len(preDefinedClusterMetrics) {
		for _, metric := range preDefinedClusterMetrics {
			metric.Namespace = globalDataNamespace
			_, err := d.metricClient.Create(metric)
			if err != nil && !errors.IsAlreadyExists(err) {
				return err
			}
		}
	}
	return nil
}

var (
	preDefinedClusterMetrics = getPredefinedClusterMetrics()
	preDefinedClusterGraph   = getPredefinedClusterGraph()
)

func getPredefinedClusterMetrics() []*managementv3.MonitorMetric {
	yamls := strings.Split(MonitorMetricsTemplate, "\n---\n")
	var rtn []*managementv3.MonitorMetric
	for _, yml := range yamls {
		var tmp managementv3.MonitorMetric
		if err := yamlToObject(yml, &tmp); err != nil {
			panic(err)
		}
		if tmp.Name == "" {
			continue
		}
		rtn = append(rtn, &tmp)
	}

	return rtn
}

func getPredefinedClusterGraph() []*managementv3.GlobalMonitorGraph {
	yamls := strings.Split(GlobalMetricExpression, "\n---\n")
	var rtn []*managementv3.GlobalMonitorGraph
	for _, yml := range yamls {
		var tmp managementv3.GlobalMonitorGraph
		if err := yamlToObject(yml, &tmp); err != nil {
			panic(err)
		}
		if tmp.Name == "" {
			continue
		}
		rtn = append(rtn, &tmp)
	}

	return rtn
}

func yamlToObject(yml string, obj interface{}) error {
	jsondata, err := yaml.YAMLToJSON([]byte(yml))
	if err != nil {
		return err
	}
	if err := json.Unmarshal(jsondata, obj); err != nil {
		return err
	}
	return nil
}

var (
	GlobalMetricExpression = `
---
# Source: metric-expression-cluster/templates/graphcluster.yaml
apiVersion: management.cattle.io/v3
kind: GlobalMonitorGraph
metadata:
  labels:
    app: metric-expression
    source: rancher-monitoring
    level: cluster
    component: cluster
  name: cluster-cpu-usage
spec:
  resourceType: cluster
  priority: 100
  title: cluster-cpu-usage
  metricsSelector:
    details: "false"
    component: cluster
    metric: cpu-usage-percent
  detailsMetricsSelector:
    details: "true"
    component: cluster
    metric: cpu-usage-percent
  yAxis:
    unit: percent
---
apiVersion: management.cattle.io/v3
kind: GlobalMonitorGraph
metadata:
  labels:
    app: metric-expression
    source: rancher-monitoring
    level: cluster
    component: cluster
  name: cluster-memory-usage
spec:
  resourceType: cluster
  priority: 102
  title: cluster-memory-usage
  metricsSelector:
    details: "false"
    component: cluster
    metric: memory-usage-percent
  detailsMetricsSelector:
    details: "true"
    component: cluster
    metric: memory-usage-percent
  yAxis:
    unit: percent
`
	MonitorMetricsTemplate = `
---
kind: MonitorMetric
apiVersion: management.cattle.io/v3
metadata:
  name: cluster-cpu-usage-percent
  labels:
    app: metric-expression
    component: cluster
    details: "false"
    level: cluster
    metric: cpu-usage-percent
    source: rancher-monitoring
spec:
  expression: 1 - (avg(irate(node_cpu_seconds_total{mode="idle"}[5m])) by (prometheus_from))
  legendFormat: '[[prometheus_from]]'
  description: cluster cpu usage percent
---
kind: MonitorMetric
apiVersion: management.cattle.io/v3
metadata:
  name: cluster-cpu-usage-percent-details
  labels:
    app: metric-expression
    component: cluster
    details: "true"
    level: cluster
    metric: cpu-usage-percent
    source: rancher-monitoring
spec:
  expression: 1 - (avg(irate(node_cpu_seconds_total{mode="idle"}[5m])) by (prometheus_from))
  legendFormat: '[[prometheus_from]]'
  description: cluster cpu usage percent details
---
kind: MonitorMetric
apiVersion: management.cattle.io/v3
metadata:
  name: cluster-memory-usage-percent
  labels:
    app: metric-expression
    component: cluster
    details: "false"
    level: cluster
    metric: memory-usage-percent
    source: rancher-monitoring
spec:
  expression: 100 * (1 - sum(node_memory_MemAvailable_bytes) by (prometheus_from) / sum(node_memory_MemTotal_bytes) by (prometheus_from))
  legendFormat: Memory usage
  description: cluster memory usage percent
---
kind: MonitorMetric
apiVersion: management.cattle.io/v3
metadata:
  name: cluster-memory-usage-percent-details
  labels:
    app: metric-expression
    component: cluster
    details: "true"
    level: cluster
    metric: memory-usage-percent
    source: rancher-monitoring
spec:
  expression: 100 * (1 - sum(node_memory_MemAvailable_bytes) by (prometheus_from) / sum(node_memory_MemTotal_bytes) by (prometheus_from))
  legendFormat: '[[prometheus_from]]'
  description: cluster memory usage percent
`
)
