package clusterstatus

import (
	"fmt"
	"strconv"

	"github.com/rancher/types/apis/management.cattle.io/v3"
	"github.com/rancher/types/config"
)

const temporaryCredentialsAnnotationKey = "clusterstatus.management.cattle.io/temporary-security-credentials"

func Register(management *config.ManagementContext) {
	c := &clusterAnnotations{
		clusters: management.Management.Clusters(""),
	}

	management.Management.Clusters("").AddHandler("temporary-credentials", c.sync)
}

type clusterAnnotations struct {
	clusters v3.ClusterInterface
}

func (cd *clusterAnnotations) sync(key string, cluster *v3.Cluster) error {
	if key == "" || cluster == nil || cluster.DeletionTimestamp != nil {
		return nil
	}

	if genericConfig := cluster.Spec.GenericEngineConfig; genericConfig != nil {
		if eksConfig, ok := (*genericConfig)["amazonElasticContainerServiceConfig"].(map[string]string); ok {
			newValue := strconv.FormatBool(eksConfig["sessionToken"] != "")

			if newValue != cluster.Annotations[temporaryCredentialsAnnotationKey] {
				original := cluster
				cluster = original.DeepCopy()

				if cluster.Annotations == nil {
					cluster.Annotations = make(map[string]string)
				}

				cluster.Annotations[temporaryCredentialsAnnotationKey] = newValue
				_, err := cd.clusters.Update(cluster)
				if err != nil {
					return fmt.Errorf("error updating temporary credentials annotation: %v", err)
				}
			}
		}
	}

	return nil
}
