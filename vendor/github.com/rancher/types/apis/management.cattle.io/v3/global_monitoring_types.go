package v3

import (
	"github.com/rancher/norman/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type GlobalMonitorGraph struct {
	types.Namespaced

	metav1.TypeMeta `json:",inline"`
	// Standard object’s metadata. More info:
	// https://github.com/kubernetes/community/blob/master/contributors/devel/api-conventions.md#metadata
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec GlobalMonitorGraphSpec `json:"spec"`
}

type GlobalMonitorGraphSpec struct {
	ResourceType        string `json:"resourceType,omitempty"  norman:"type=enum,options=node|cluster|etcd|apiserver|scheduler|controllermanager|fluentd|istiocluster|istioproject"`
	DisplayResourceType string `json:"displayResourceType,omitempty" norman:"type=enum,options=node|cluster|etcd|kube-component|rancher-component"`
	CommonMonitorGraphSpec
}

type GlobalAlertRule struct {
	types.Namespaced

	metav1.TypeMeta `json:",inline"`
	// Standard object’s metadata. More info:
	// https://github.com/kubernetes/community/blob/master/contributors/devel/api-conventions.md#metadata
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec GlobalAlertRuleSpec `json:"spec"`
	// Most recent observed status of the alert. More info:
	// https://github.com/kubernetes/community/blob/master/contributors/devel/api-conventions.md#spec-and-status
	Status AlertStatus `json:"status"`
}

type GlobalAlertRuleSpec struct {
	CommonRuleField
	GroupName  string      `json:"groupName" norman:"type=reference[clusterAlertGroup]"`
	MetricRule *MetricRule `json:"metricRule,omitempty"`
}

type GlobalAlertGroup struct {
	types.Namespaced

	metav1.TypeMeta `json:",inline"`
	// Standard object’s metadata. More info:
	// https://github.com/kubernetes/community/blob/master/contributors/devel/api-conventions.md#metadata
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec GlobalGroupSpec `json:"spec"`
	// Most recent observed status of the alert. More info:
	// https://github.com/kubernetes/community/blob/master/contributors/devel/api-conventions.md#spec-and-status
	Status AlertStatus `json:"status"`
}

type GlobalGroupSpec struct {
	Recipients []Recipient `json:"recipients,omitempty"`
	CommonGroupField
}
