package globalmonitoring

import (
	"testing"
	"time"

	"github.com/rancher/rancher/pkg/monitoring"

	"github.com/rancher/rancher/pkg/settings"

	"github.com/stretchr/testify/assert"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/rancher/norman/types"
	v3 "github.com/rancher/types/apis/management.cattle.io/v3"
	pv3 "github.com/rancher/types/apis/project.cattle.io/v3"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"

	cfakes "github.com/rancher/types/apis/core/v1/fakes"

	"github.com/rancher/types/apis/management.cattle.io/v3/fakes"
	pfakes "github.com/rancher/types/apis/project.cattle.io/v3/fakes"
	apitypes "k8s.io/apimachinery/pkg/types"
)

type userManagerMock struct {
}

func (um userManagerMock) SetPrincipalOnCurrentUser(apiContext *types.APIContext, principal v3.Principal) (*v3.User, error) {
	return nil, nil
}
func (um userManagerMock) GetUser(apiContext *types.APIContext) string {
	return ""
}
func (um userManagerMock) EnsureToken(tokenName, description, kind, userName string) (string, error) {
	return "", nil
}
func (um userManagerMock) EnsureClusterToken(clusterName, tokenName, description, kind, userName string) (string, error) {
	return "", nil
}
func (um userManagerMock) EnsureUser(principalName, displayName string) (*v3.User, error) {
	return nil, nil
}
func (um userManagerMock) CheckAccess(accessMode string, allowedPrincipalIDs []string, userPrincipalID string, groups []v3.Principal) (bool, error) {
	return false, nil
}
func (um userManagerMock) SetPrincipalOnCurrentUserByUserID(userID string, principal v3.Principal) (*v3.User, error) {
	return nil, nil
}
func (um userManagerMock) CreateNewUserClusterRoleBinding(userName string, userUID apitypes.UID) error {
	return nil
}
func (um userManagerMock) GetUserByPrincipalID(principalName string) (*v3.User, error) {
	return nil, nil
}

func newMockAppHandler(clusters map[string]*v3.Cluster, apps map[string]*pv3.App) appHandler {
	ah := appHandler{
		clusterClient: &fakes.ClusterInterfaceMock{
			UpdateFunc: func(in1 *v3.Cluster) (*v3.Cluster, error) {
				clusters[in1.Name] = in1
				return in1, nil
			},
			GetFunc: func(name string, opts metav1.GetOptions) (*v3.Cluster, error) {
				return clusters[name], nil
			},
		},
		clusterLister: &fakes.ClusterListerMock{
			ListFunc: func(namespace string, selector labels.Selector) ([]*v3.Cluster, error) {
				var items []*v3.Cluster
				for _, c := range clusters {
					items = append(items, c)
				}
				return items, nil
			},
		},
		appClient: &pfakes.AppInterfaceMock{
			UpdateFunc: func(in1 *pv3.App) (*pv3.App, error) {
				apps[in1.Name] = in1
				return in1, nil
			},
			GetFunc: func(name string, opts metav1.GetOptions) (*pv3.App, error) {
				return apps[name], nil
			},
		},
		secretLister: &cfakes.SecretListerMock{
			GetFunc: func(namespace string, name string) (*v1.Secret, error) {
				return nil, errors.NewNotFound(schema.GroupResource{}, "")
			},
		},
		secretClient: &cfakes.SecretInterfaceMock{
			CreateFunc: func(in1 *v1.Secret) (*v1.Secret, error) {
				return nil, nil
			},
		},
		userLister: &fakes.UserListerMock{
			ListFunc: func(namespace string, selector labels.Selector) ([]*v3.User, error) {
				users := []*v3.User{
					{
						DisplayName: "Default Admin",
					},
				}
				return users, nil
			},
		},
		userManager: userManagerMock{},
	}
	return ah
}

func TestGlobalMonitoringAppInstall(t *testing.T) {
	settings.ServerURL.Set("https://test.example.com")
	settings.GlobalMonitoringClusterID.Set("cluster1")
	assert := assert.New(t)
	clusters := map[string]*v3.Cluster{
		"cluster1": {
			ObjectMeta: metav1.ObjectMeta{
				Name: "cluster1",
				Annotations: map[string]string{
					"field.cattle.io/overwriteAppAnswers": "{\"answers\":{\"operator-init.enabled\":\"true\"},\"version\":\"0.0.5\"}",
				},
			},
			Status: v3.ClusterStatus{
				Conditions: []v3.ClusterCondition{
					{
						Type:   "MonitoringEnabled",
						Status: v1.ConditionTrue,
					},
				},
			},
		},
		"cluster2": {
			ObjectMeta: metav1.ObjectMeta{
				Name: "cluster2",
			},
		},
	}
	apps := map[string]*pv3.App{
		"cluster-monitoring": {
			ObjectMeta: metav1.ObjectMeta{
				Name:      "cluster-monitoring",
				Namespace: "cluster1",
			},
		},
	}

	ah := newMockAppHandler(clusters, apps)

	app := pv3.App{
		ObjectMeta: metav1.ObjectMeta{
			Name: globalMonitoringAppName,
		},
		Spec: pv3.AppSpec{
			TargetNamespace: globalDataNamespace,
		},
	}
	_, err := ah.sync(&app)
	assert.Empty(err)
	updatedApp, err := ah.appClient.Get(globalMonitoringAppName, metav1.GetOptions{})
	assert.Empty(err)
	assert.Equal("test.example.com", updatedApp.Spec.Answers[rancherHostnameKey])
	assert.Equal("cluster1", updatedApp.Spec.Answers[clusterIDAnswerKey])

	_, err = ah.sync(updatedApp)
	assert.Empty(err)
	updatedCluster, err := ah.clusterClient.Get("cluster1", metav1.GetOptions{})
	updatedClusterMonitoringAnswers, _ := monitoring.GetOverwroteAppAnswersAndVersion(updatedCluster.Annotations)
	assert.Equal("true", updatedClusterMonitoringAnswers[ThanosEnabledAnswerKey])
}

func TestGlobalMonitoringAppRemove(t *testing.T) {
	settings.ServerURL.Set("https://test.example.com")
	settings.GlobalMonitoringClusterID.Set("cluster1")
	assert := assert.New(t)
	clusters := map[string]*v3.Cluster{
		"cluster1": {
			ObjectMeta: metav1.ObjectMeta{
				Name: "cluster1",
				Annotations: map[string]string{
					"field.cattle.io/overwriteAppAnswers": "{\"answers\":{\"prometheus.thanos.enabled\":\"true\"},\"version\":\"0.0.5\"}",
				},
			},
			Status: v3.ClusterStatus{
				Conditions: []v3.ClusterCondition{
					{
						Type:   "MonitoringEnabled",
						Status: v1.ConditionTrue,
					},
				},
			},
		},
		"cluster2": {
			ObjectMeta: metav1.ObjectMeta{
				Name: "cluster2",
			},
		},
	}
	apps := map[string]*pv3.App{
		"cluster-monitoring": {
			ObjectMeta: metav1.ObjectMeta{
				Name:      "cluster-monitoring",
				Namespace: "cluster1",
			},
		},
	}

	ah := newMockAppHandler(clusters, apps)

	app := pv3.App{
		ObjectMeta: metav1.ObjectMeta{
			Name:              globalMonitoringAppName,
			DeletionTimestamp: &metav1.Time{Time: time.Now()},
		},
		Spec: pv3.AppSpec{
			TargetNamespace: globalDataNamespace,
			Answers: map[string]string{
				rancherHostnameKey: "test.example.com",
				clusterIDAnswerKey: "cluster1",
			},
		},
	}
	_, err := ah.sync(&app)
	assert.Empty(err)
	updatedCluster, err := ah.clusterClient.Get("cluster1", metav1.GetOptions{})
	updatedClusterMonitoringAnswers, _ := monitoring.GetOverwroteAppAnswersAndVersion(updatedCluster.Annotations)
	assert.Equal("", updatedClusterMonitoringAnswers[ThanosEnabledAnswerKey])
}
