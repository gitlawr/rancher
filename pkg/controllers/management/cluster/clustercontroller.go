package cluster

import (
	"context"
	"fmt"
	"github.com/rancher/kontainer-engine/service"
	"github.com/rancher/kontainer-engine/types"
	"github.com/rancher/rancher/pkg/controllers/management/clusterprovisioner"
	"github.com/rancher/rke/cloudprovider/aws"
	"github.com/rancher/rke/cloudprovider/azure"
	"github.com/rancher/types/apis/core/v1"
	"github.com/rancher/types/apis/management.cattle.io/v3"
	"github.com/rancher/types/config"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/labels"
	"reflect"
)

const (
	GoogleCloudLoadBalancer = "GCLB"
	ElasticLoadBalancer     = "ELB"
	AzureL4LB               = "Azure L4 LB"
	NginxIngressProvider    = "Nginx"
	DefaultNodePortRange    = "30000-32767"
)

type controller struct {
	clusterClient         v3.ClusterInterface
	clusterLister         v3.ClusterLister
	nodeLister            v3.NodeLister
	kontainerDriverLister v3.KontainerDriverLister
	namespaces            v1.NamespaceInterface
	coreV1                v1.Interface
}

func Register(management *config.ManagementContext) {
	c := controller{
		clusterClient:         management.Management.Clusters(""),
		clusterLister:         management.Management.Clusters("").Controller().Lister(),
		nodeLister:            management.Management.Nodes("").Controller().Lister(),
		kontainerDriverLister: management.Management.KontainerDrivers("").Controller().Lister(),
		namespaces:            management.Core.Namespaces(""),
		coreV1:                management.Core,
	}

	c.clusterClient.AddHandler("clusterCreateUpdate", c.capsSync)
}

func (c *controller) capsSync(key string, cluster *v3.Cluster) error {
	var err error
	if cluster == nil || cluster.DeletionTimestamp != nil {
		return nil
	}

	if cluster.Spec.ImportedConfig != nil {
		return nil
	}
	capabilities := v3.Capabilities{}
	capabilities.NodePortRange = DefaultNodePortRange
	if cluster.Spec.RancherKubernetesEngineConfig != nil {
		if capabilities, err = c.RKECapabilities(capabilities, *cluster.Spec.RancherKubernetesEngineConfig, cluster.Name); err != nil {
			return err
		}
	} else if cluster.Spec.GenericEngineConfig != nil {
		driverName, ok := (*cluster.Spec.GenericEngineConfig)["driverName"].(string)
		if !ok {
			logrus.Warnf("cluster %v had generic engine config but no driver name, k8s capabilities will "+
				"not be populated correctly", key)
			return nil
		}

		kontainerDriver, err := c.kontainerDriverLister.Get("", driverName)
		if err != nil {
			return fmt.Errorf("error getting kontainer driver: %v", err)
		}

		driver := service.NewEngineService(
			clusterprovisioner.NewPersistentStore(c.namespaces, c.coreV1),
		)
		k8sCapabilities, err := driver.GetK8sCapabilities(context.Background(), kontainerDriver.Name, kontainerDriver,
			cluster.Spec)
		if err != nil {
			return fmt.Errorf("error getting k8s capabilities: %v", err)
		}

		capabilities = toCapabilities(k8sCapabilities)
	} else {
		return nil
	}

	if !reflect.DeepEqual(capabilities, cluster.Status.Capabilities) {
		toUpdateCluster := cluster.DeepCopy()
		toUpdateCluster.Status.Capabilities = capabilities
		if _, err := c.clusterClient.Update(toUpdateCluster); err != nil {
			return err
		}
	}

	return nil
}

func (c *controller) RKECapabilities(capabilities v3.Capabilities, rkeConfig v3.RancherKubernetesEngineConfig, clusterName string) (v3.Capabilities, error) {
	switch rkeConfig.CloudProvider.Name {
	case aws.AWSCloudProviderName:
		capabilities.L4LoadBalancer = c.L4Capability(true, ElasticLoadBalancer, []string{"TCP"}, true)
	case azure.AzureCloudProviderName:
		capabilities.L4LoadBalancer = c.L4Capability(true, AzureL4LB, []string{"TCP", "UDP"}, true)
	}
	// only if not custom, non custom clusters have nodepools set
	nodes, err := c.nodeLister.List(clusterName, labels.Everything())
	if err != nil {
		return capabilities, err
	}

	if len(nodes) > 0 {
		if nodes[0].Spec.NodePoolName != "" {
			capabilities.NodePoolScalingSupported = true
		}
	}

	ingressController := c.IngressCapability(true, NginxIngressProvider, false)
	capabilities.IngressControllers = []v3.IngressController{ingressController}
	if rkeConfig.Services.KubeAPI.ExtraArgs["service-node-port-range"] != "" {
		capabilities.NodePortRange = rkeConfig.Services.KubeAPI.ExtraArgs["service-node-port-range"]
	}

	return capabilities, nil
}

func (c *controller) L4Capability(enabled bool, providerName string, protocols []string, healthCheck bool) v3.L4LoadBalancer {
	l4lb := v3.L4LoadBalancer{
		Enabled:              enabled,
		Provider:             providerName,
		ProtocolsSupported:   protocols,
		HealthCheckSupported: healthCheck,
	}
	return l4lb
}

func (c *controller) IngressCapability(httpLBEnabled bool, providerName string, customDefaultBackend bool) v3.IngressController {
	ing := v3.IngressController{
		HTTPLoadBalancingEnabled: httpLBEnabled,
		IngressProvider:          providerName,
		CustomDefaultBackend:     customDefaultBackend,
	}
	return ing
}

func toCapabilities(k8sCapabilities *types.K8SCapabilities) v3.Capabilities {
	var controllers []v3.IngressController

	for _, controller := range k8sCapabilities.IngressControllers {
		controllers = append(controllers, v3.IngressController{
			CustomDefaultBackend:     controller.CustomDefaultBackend,
			HTTPLoadBalancingEnabled: controller.HTTPLoadBalancingEnabled,
			IngressProvider:          controller.IngressProvider,
		})
	}

	return v3.Capabilities{
		IngressControllers: controllers,
		L4LoadBalancer: v3.L4LoadBalancer{
			Enabled:              k8sCapabilities.L4LoadBalancer.Enabled,
			HealthCheckSupported: k8sCapabilities.L4LoadBalancer.HealthCheckSupported,
			ProtocolsSupported:   k8sCapabilities.L4LoadBalancer.ProtocolsSupported,
			Provider:             k8sCapabilities.L4LoadBalancer.Provider,
		},
		NodePoolScalingSupported: k8sCapabilities.NodePoolScalingSupported,
		NodePortRange:            k8sCapabilities.NodePortRange,
	}
}
