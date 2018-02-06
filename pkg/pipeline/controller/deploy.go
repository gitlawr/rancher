package controller

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (l *ClusterPipelineLifecycle) deployK() error {
	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "cattle-pipeline",
		},
	}
	if _, err := l.cluster.Core.Namespaces("").Create(ns); err != nil && !apierrors.IsAlreadyExists(err) {
		logrus.Errorf("Error occured while create ns: %v", err)
		return errors.Wrapf(err, "Creating ns")
	}

	secret := getSecret()
	if _, err := l.cluster.Core.Secrets("").Create(secret); err != nil && !apierrors.IsAlreadyExists(err) {
		logrus.Errorf("Error occured while create secret: %v", err)
		return errors.Wrapf(err, "Creating secret")
	}

	configmap := getConfigMap()
	if _, err := l.cluster.K8sClient.CoreV1().ConfigMaps("cattle-pipeline").Create(configmap); err != nil && !apierrors.IsAlreadyExists(err) {
		logrus.Errorf("Error occured while create configmap: %v", err)
		return errors.Wrapf(err, "Creating configmap")
	}

	service := getJenkinsService()
	if _, err := l.cluster.Core.Services("").Create(service); err != nil && !apierrors.IsAlreadyExists(err) {
		logrus.Errorf("Error occured while create service: %v", err)
		return errors.Wrapf(err, "Creating service")
	}

	agentservice := getJenkinsAgentService()
	if _, err := l.cluster.Core.Services("").Create(agentservice); err != nil && !apierrors.IsAlreadyExists(err) {
		logrus.Errorf("Error occured while create service: %v", err)
		return errors.Wrapf(err, "Creating service")
	}

	sa := getServiceAccount()
	if _, err := l.cluster.K8sClient.CoreV1().ServiceAccounts("cattle-pipeline").Create(sa); err != nil && !apierrors.IsAlreadyExists(err) {
		logrus.Errorf("Error occured while create service account: %v", err)
		return errors.Wrapf(err, "Creating service account")
	}
	rb := getRoleBindings()
	if _, err := l.cluster.K8sClient.RbacV1().ClusterRoleBindings().Create(rb); err != nil && !apierrors.IsAlreadyExists(err) {
		logrus.Errorf("Error occured while create role binding: %v", err)
		return errors.Wrapf(err, "Creating role binding")
	}
	deployment := getJenkinsDeployment()
	if _, err := l.cluster.Apps.Deployments("").Create(deployment); err != nil && !apierrors.IsAlreadyExists(err) {
		logrus.Errorf("Error occured while create deployment: %v", err)
		return errors.Wrapf(err, "Creating deployment")
	}

	return nil
}
