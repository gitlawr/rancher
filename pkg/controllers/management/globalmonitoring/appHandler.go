package globalmonitoring

import (
	"encoding/json"
	"net/url"
	"strings"

	"github.com/pkg/errors"
	"github.com/rancher/rancher/pkg/monitoring"
	"github.com/rancher/rancher/pkg/randomtoken"
	"github.com/rancher/rancher/pkg/settings"
	v1 "github.com/rancher/types/apis/core/v1"
	mgmtv3 "github.com/rancher/types/apis/management.cattle.io/v3"
	v3 "github.com/rancher/types/apis/project.cattle.io/v3"
	"github.com/rancher/types/user"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
)

const (
	ThanosEnabledAnswerKey   = "prometheus.thanos.enabled"
	rancherHostnameKey       = "rancherHostname"
	objectConfigAnswerPrefix = "thanos.objectConfig"
	authTokenKey             = "token"
	apiTokenKey              = "apiToken"
)

type appHandler struct {
	clusterClient mgmtv3.ClusterInterface
	appClient     v3.AppInterface
	clusterLister mgmtv3.ClusterLister
	secretLister  v1.SecretLister
	secretClient  v1.SecretInterface
	userLister    mgmtv3.UserLister
	userManager   user.Manager
}

func (ah *appHandler) Create(obj *v3.App) (runtime.Object, error) {
	return obj, nil
}

func (ah *appHandler) Updated(obj *v3.App) (runtime.Object, error) {
	return ah.sync(obj)
}

func (ah *appHandler) Remove(obj *v3.App) (runtime.Object, error) {
	return ah.sync(obj)
}

func (ah *appHandler) sync(app *v3.App) (runtime.Object, error) {
	if app == nil {
		return app, nil
	}
	if (app.Name != globalMonitoringAppName) || (app.Spec.TargetNamespace != globalDataNamespace) {
		return app, nil
	}

	if app.Spec.Answers == nil || app.Spec.Answers[rancherHostnameKey] == "" {
		return app, ah.initGlobalMonitoringApp(app)
	}

	if err := ah.UpdateAllClusterMonitoringThanosEnabled(app); err != nil {
		return app, err
	}
	return app, nil
}

func (ah *appHandler) initGlobalMonitoringApp(app *v3.App) error {
	clusters, err := ah.clusterLister.List("", labels.NewSelector())
	if err != nil {
		return err
	}
	var monitoringEnabledClusters []string
	for _, c := range clusters {
		if !mgmtv3.ClusterConditionMonitoringEnabled.IsTrue(c) {
			continue
		}
		monitoringEnabledClusters = append(monitoringEnabledClusters, c.Name)
	}
	toUpdateApp := app.DeepCopy()
	if toUpdateApp.Spec.Answers == nil {
		toUpdateApp.Spec.Answers = make(map[string]string)
	}
	serverURL := settings.ServerURL.Get()
	u, err := url.Parse(serverURL)
	if err != nil {
		return err
	}
	token, err := ah.GetOrCreateGlobalMonitoringToken()
	if err != nil {
		return err
	}
	apiToken, err := ah.getAPIToken()
	if err != nil {
		return err
	}
	toUpdateApp.Spec.Answers[rancherHostnameKey] = u.Hostname()
	toUpdateApp.Spec.Answers[clusterIDAnswerKey] = strings.Join(monitoringEnabledClusters, ":")
	toUpdateApp.Spec.Answers[authTokenKey] = token
	toUpdateApp.Spec.Answers[apiTokenKey] = apiToken
	_, err = ah.appClient.Update(toUpdateApp)
	return err
}

func (ah *appHandler) GetOrCreateGlobalMonitoringToken() (string, error) {
	secret, err := ah.secretLister.Get(globalDataNamespace, globalMonitoringSecretName)
	if err == nil {
		return string(secret.Data[authTokenKey]), nil
	} else if !apierrors.IsNotFound(err) {
		return "", err
	}

	token, err := randomtoken.Generate()
	if err != nil {
		return "", err
	}
	secret = &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      globalMonitoringSecretName,
			Namespace: globalDataNamespace,
		},
		StringData: map[string]string{
			authTokenKey: token,
		},
	}
	if _, err := ah.secretClient.Create(secret); err != nil {
		return "", err
	}
	return token, nil
}

func (ah *appHandler) UpdateAllClusterMonitoringThanosEnabled(app *v3.App) error {
	clusters, err := ah.clusterLister.List("", labels.NewSelector())
	if err != nil {
		return err
	}
	for _, c := range clusters {
		if !mgmtv3.ClusterConditionMonitoringEnabled.IsTrue(c) {
			continue
		}
		if err := UpdateClusterMonitoringAnswers(ah.clusterClient, c, app); err != nil {
			return err
		}
	}
	return nil
}

func (ah *appHandler) getAPIToken() (string, error) {
	users, err := ah.userLister.List("", labels.SelectorFromSet(
		map[string]string{
			"authz.management.cattle.io/bootstrapping": "admin-user",
		},
	))
	if err != nil {
		return "", errors.Wrapf(err, "Can't list admin-user for telemetry. Err: %v. Retry after 5 seconds", err)
	}
	for _, user := range users {
		if user.DisplayName == "Default Admin" {
			token, err := ah.userManager.EnsureToken("global-monitoring", "global monitoring token", "global-monitoring", user.Name)
			if err != nil {
				return "", errors.Wrapf(err, "Can't create token for global monitoring. Err: %v", err)
			}
			return token, nil
		}
	}
	return "", nil
}

func UpdateClusterMonitoringAnswers(clusterClient mgmtv3.ClusterInterface, cluster *mgmtv3.Cluster, app *v3.App) error {
	answers := map[string]string{}
	if app == nil || app.DeletionTimestamp != nil {
		answers[ThanosEnabledAnswerKey] = ""
	} else {
		answers[ThanosEnabledAnswerKey] = "true"

		for k, v := range app.Spec.Answers {
			if strings.HasPrefix(k, objectConfigAnswerPrefix) {
				answers["prometheus."+k] = v
			}
		}
	}

	toUpdateAnswers, _ := monitoring.GetOverwroteAppAnswersAndVersion(cluster.Annotations)

	thanosAnswer := toUpdateAnswers[ThanosEnabledAnswerKey]
	expected := answers[ThanosEnabledAnswerKey]
	if thanosAnswer == expected {
		return nil
	}
	toUpdateCluster := cluster.DeepCopy()
	for k, v := range answers {
		toUpdateAnswers[k] = v
	}

	data := map[string]interface{}{}
	data["answers"] = answers
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	toUpdateCluster.Annotations = monitoring.AppendAppOverwritingAnswers(toUpdateCluster.Annotations, string(b))
	if _, err := clusterClient.Update(toUpdateCluster); err != nil {
		return err
	}

	return nil
}
