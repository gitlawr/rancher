package globalmonitoring

import (
	"net/url"
	"strings"

	"github.com/pkg/errors"
	"github.com/rancher/types/user"

	"github.com/sirupsen/logrus"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"

	"github.com/rancher/rancher/pkg/randomtoken"
	v1 "github.com/rancher/types/apis/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/rancher/rancher/pkg/app/utils"
	"github.com/rancher/rancher/pkg/settings"
	mgmtv3 "github.com/rancher/types/apis/management.cattle.io/v3"
	v3 "github.com/rancher/types/apis/project.cattle.io/v3"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
)

const (
	ThanosEnabledAnswerKey = "prometheus.thanos.enabled"
	rancherHostnameKey     = "rancherHostname"
	delegateAnswerPrefix   = "delegate."
)

type appHandler struct {
	appLister     v3.AppLister
	appClient     v3.AppInterface
	projectLister mgmtv3.ProjectLister
	clusterLister mgmtv3.ClusterLister
	secretLister  v1.SecretLister
	secretClient  v1.SecretInterface
	userLister    mgmtv3.UserLister
	userManager   user.Manager
}

func (ah *appHandler) sync(key string, app *v3.App) (runtime.Object, error) {
	if app == nil {
		return app, nil
	}
	if (app.Name != globalMonitoringAppName) || (app.Spec.TargetNamespace != globalDataNamespace) {
		return app, nil
	}

	if app.Spec.Answers == nil || app.Spec.Answers[rancherHostnameKey] == "" {
		return app, ah.initGlobalMonitoringApp(app)
	}

	delegateAnswers := map[string]string{}
	for k, v := range app.Spec.Answers {
		if strings.HasPrefix(k, delegateAnswerPrefix) {
			delegateAnswers[strings.TrimPrefix(k, delegateAnswerPrefix)] = v
		}
	}

	if app.DeletionTimestamp == nil {
		delegateAnswers[ThanosEnabledAnswerKey] = "true"
	}

	if err := ah.UpdateClusterMonitoringThanosEnabled(delegateAnswers); err != nil {
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
	toUpdateApp.Spec.Answers["token"] = token
	toUpdateApp.Spec.Answers["apiToken"] = apiToken
	_, err = ah.appClient.Update(toUpdateApp)
	return err
}

func (ah *appHandler) GetOrCreateGlobalMonitoringToken() (string, error) {
	secret, err := ah.secretLister.Get(globalDataNamespace, globalMonitoringSecretName)
	if err == nil {
		logrus.Infof("getting %s", string(secret.Data["token"]))
		return string(secret.Data["token"]), nil
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
			"token": token,
		},
	}
	if _, err := ah.secretClient.Create(secret); err != nil {
		return "", err
	}
	return token, nil
}

func (ah *appHandler) UpdateClusterMonitoringThanosEnabled(delegateAnswers map[string]string) error {
	clusters, err := ah.clusterLister.List("", labels.NewSelector())
	if err != nil {
		return err
	}
	for _, c := range clusters {
		if !mgmtv3.ClusterConditionMonitoringEnabled.IsTrue(c) {
			continue
		}
		projectID, err := utils.GetSystemProjectID(c.Name, ah.projectLister)
		if err != nil {
			return err
		}
		clusterMonitoringApp, err := ah.appLister.Get(projectID, "cluster-monitoring")
		if err != nil {
			return err
		}
		if err := ah.UpdateAppThanosAnswers(clusterMonitoringApp, delegateAnswers); err != nil {
			return err
		}
	}
	return nil
}

func (ah *appHandler) UpdateAppThanosAnswers(app *v3.App, answers map[string]string) error {
	thanosAnswer := app.Spec.Answers[ThanosEnabledAnswerKey]
	expected := answers[ThanosEnabledAnswerKey]
	if thanosAnswer == expected {
		return nil
	}
	toUpdate := app.DeepCopy()
	for k, v := range answers {
		toUpdate.Spec.Answers[k] = v
	}
	if _, err := ah.appClient.Update(toUpdate); err != nil {
		return err
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
