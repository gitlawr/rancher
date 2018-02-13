package clusterpipeline

import (
	"github.com/pkg/errors"
	"github.com/rancher/rancher/pkg/controllers/user/pipeline/utils"
	"github.com/rancher/types/apis/core/v1"
	"github.com/rancher/types/apis/management.cattle.io/v3"
	"github.com/rancher/types/config"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Syncer struct {
	cluster    *config.UserContext
	namespaces v1.NamespaceInterface
}

func Register(cluster *config.UserContext) {
	clusterPipelines := cluster.Management.Management.ClusterPipelines("")
	clusterPipelineSyncer := &Syncer{
		cluster: cluster,
	}
	clusterPipelines.AddHandler("cluster-pipeline-syncer", clusterPipelineSyncer.Sync)
}

func (s *Syncer) Sync(key string, obj *v3.ClusterPipeline) error {
	//ensure clusterpipeline singleton in the cluster
	utils.InitClusterPipeline(s.cluster)

	if obj.Spec.Deploy {
		if err := s.deploy(); err != nil {
			return err
		}
	} else {
		if err := s.destroy(); err != nil {
			return err
		}
	}

	return nil
}

func (s *Syncer) destroy() error {
	if err := s.cluster.Core.Namespaces("").Delete(utils.PipelineNamespace, &metav1.DeleteOptions{}); err != nil && !apierrors.IsNotFound(err) {

		logrus.Errorf("Error occured while removing ns: %v", err)
		return err

	}

	return nil
}

func (s *Syncer) deploy() error {
	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: utils.PipelineNamespace,
		},
	}
	if _, err := s.cluster.Core.Namespaces("").Create(ns); err != nil && !apierrors.IsAlreadyExists(err) {
		logrus.Errorf("Error occured while create ns: %v", err)
		return errors.Wrapf(err, "Creating ns")
	}

	secret := getSecret()
	if _, err := s.cluster.Core.Secrets("").Create(secret); err != nil && !apierrors.IsAlreadyExists(err) {
		logrus.Errorf("Error occured while create secret: %v", err)
		return errors.Wrapf(err, "Creating secret")
	}

	configmap := getConfigMap()
	if _, err := s.cluster.K8sClient.CoreV1().ConfigMaps(utils.PipelineNamespace).Create(configmap); err != nil && !apierrors.IsAlreadyExists(err) {
		logrus.Errorf("Error occured while create configmap: %v", err)
		return errors.Wrapf(err, "Creating configmap")
	}

	service := getJenkinsService()
	if _, err := s.cluster.Core.Services("").Create(service); err != nil && !apierrors.IsAlreadyExists(err) {
		logrus.Errorf("Error occured while create service: %v", err)
		return errors.Wrapf(err, "Creating service")
	}

	agentservice := getJenkinsAgentService()
	if _, err := s.cluster.Core.Services("").Create(agentservice); err != nil && !apierrors.IsAlreadyExists(err) {
		logrus.Errorf("Error occured while create service: %v", err)
		return errors.Wrapf(err, "Creating service")
	}

	sa := getServiceAccount()
	if _, err := s.cluster.K8sClient.CoreV1().ServiceAccounts(utils.PipelineNamespace).Create(sa); err != nil && !apierrors.IsAlreadyExists(err) {
		logrus.Errorf("Error occured while create service account: %v", err)
		return errors.Wrapf(err, "Creating service account")
	}
	rb := getRoleBindings()
	if _, err := s.cluster.K8sClient.RbacV1().ClusterRoleBindings().Create(rb); err != nil && !apierrors.IsAlreadyExists(err) {
		logrus.Errorf("Error occured while create role binding: %v", err)
		return errors.Wrapf(err, "Creating role binding")
	}
	deployment := getJenkinsDeployment()
	if _, err := s.cluster.Apps.Deployments("").Create(deployment); err != nil && !apierrors.IsAlreadyExists(err) {
		logrus.Errorf("Error occured while create deployment: %v", err)
		return errors.Wrapf(err, "Creating deployment")
	}

	return nil
}
