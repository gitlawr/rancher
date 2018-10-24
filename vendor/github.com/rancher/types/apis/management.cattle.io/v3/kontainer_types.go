package v3

import (
	"github.com/rancher/norman/condition"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type KontainerDriver struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object’s metadata. More info:
	// https://github.com/kubernetes/community/blob/master/contributors/devel/api-conventions.md#metadata
	metav1.ObjectMeta `json:"metadata,omitempty"`
	// Specification of the desired behavior of the the cluster. More info:
	// https://github.com/kubernetes/community/blob/master/contributors/devel/api-conventions.md#spec-and-status
	Spec KontainerDriverSpec `json:"spec"`
	// Most recent observed status of the cluster. More info:
	// https://github.com/kubernetes/community/blob/master/contributors/devel/api-conventions.md#spec-and-status
	Status KontainerDriverStatus `json:"status"`
}

type KontainerDriverStatus struct {
	ActualURL      string      `json:"actualUrl"`
	ExecutablePath string      `json:"executablePath"`
	Conditions     []Condition `json:"conditions"`
}

type KontainerDriverSpec struct {
	DisplayName   string `json:"displayName"`
	DesiredURL    string `json:"desirdUrl" norman:"required"`
	Checksum      string `json:"checksum"`
	BuiltIn       bool   `json:"builtIn"`
	DynamicSchema bool   `json:"dynamicSchema"`
	Active        bool   `json:"active"`
}

var KontainerDriverConditionDownloaded condition.Cond = "Downloaded"
