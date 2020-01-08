package globalmonitoring

import (
	"strings"

	"github.com/rancher/rancher/pkg/settings"

	"github.com/sirupsen/logrus"

	"github.com/rancher/norman/types/slice"

	"github.com/rancher/rancher/pkg/project"

	mgmtv3 "github.com/rancher/types/apis/management.cattle.io/v3"
	v3 "github.com/rancher/types/apis/project.cattle.io/v3"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
)

const (
	globalMonitoringAppName    = "global-monitoring"
	globalDataNamespace        = "cattle-global-data"
	globalMonitoringSecretName = "global-monitoring"
	clusterIDAnswerKey         = "clusterIds"
)

type clusterHandler struct {
	appLister     v3.AppLister
	appClient     v3.AppInterface
	projectLister mgmtv3.ProjectLister
	deployer      deployer
}

func (ch *clusterHandler) sync(key string, cluster *mgmtv3.Cluster) (runtime.Object, error) {
	if cluster == nil || cluster.DeletionTimestamp != nil {
		return cluster, nil
	}

	// create pre-defined metric and graph resources
	if err := ch.deployer.deploy(); err != nil {
		return cluster, err
	}
	adminClusterID := settings.AdminClusterID.Get()
	if adminClusterID == "" {
		return cluster, nil
	}

	systemProject, err := project.GetSystemProject(adminClusterID, ch.projectLister)
	if err != nil {
		return cluster, err
	}
	app, err := ch.appLister.Get(systemProject.Name, globalMonitoringAppName)
	if k8serrors.IsNotFound(err) {
		//global monitoring is not enabled
		return cluster, nil
	} else if err != nil {
		return cluster, err
	}

	if app.Spec.Answers == nil || app.Spec.Answers[rancherHostnameKey] == "" {
		//app is not initialized
		return cluster, nil
	}
	var appliedClusterIDs []string
	appliedClusterIDAnswer := app.Spec.Answers[clusterIDAnswerKey]
	if appliedClusterIDAnswer != "" {
		appliedClusterIDs = strings.Split(appliedClusterIDAnswer, ":")
	}
	if cluster.Spec.EnableClusterMonitoring {
		if !slice.ContainsString(appliedClusterIDs, cluster.Name) {
			toUpdateClusterIDs := append(appliedClusterIDs, cluster.Name)
			toUpdateApp := app.DeepCopy()
			logrus.Infof("enabled, going to update: %v,%v", appliedClusterIDs, cluster.Name)
			toUpdateApp.Spec.Answers[clusterIDAnswerKey] = strings.Join(toUpdateClusterIDs, ":")
			if _, err := ch.appClient.Update(toUpdateApp); err != nil {
				return cluster, err
			}
		}
		return cluster, nil
	} else {
		if slice.ContainsString(appliedClusterIDs, cluster.Name) {
			toUpdateClusterIDs := removeElement(appliedClusterIDs, cluster.Name)
			toUpdateApp := app.DeepCopy()
			logrus.Infof("is not enabled, going to update: %v,%v", appliedClusterIDs, cluster.Name)
			toUpdateApp.Spec.Answers[clusterIDAnswerKey] = strings.Join(toUpdateClusterIDs, ":")
			if _, err := ch.appClient.Update(toUpdateApp); err != nil {
				return cluster, err
			}
		}
	}
	return cluster, nil

}

func removeElement(slice []string, item string) []string {
	for i, j := range slice {
		if j == item {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}
