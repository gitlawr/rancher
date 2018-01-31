package controller

import (
	"github.com/pkg/errors"
	"github.com/rancher/types/apis/management.cattle.io/v3"
	"github.com/rancher/types/config"
	"github.com/sirupsen/logrus"
	appsv1beta2 "k8s.io/api/apps/v1beta2"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

const (
	PipelineConfig = ""
)

func getSecret() *corev1.Secret {
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "cattle-pipeline",
			Name:      "jenkins",
		},
		Data: map[string][]byte{
			"config.yml": []byte(PipelineConfig),
		},
	}
}

func getService() *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "cattle-pipeline",
			Name:      "jenkins",
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeNodePort,
			Selector: map[string]string{
				"app": "jenkins",
			},
			Ports: []corev1.ServicePort{
				{
					Name: "management",
					Port: 8080,
				},
				{
					Name: "jnlp",
					Port: 50000,
				},
			},
		},
	}
}

func getDeployment() *appsv1beta2.Deployment {
	replicas := int32(1)
	return &appsv1beta2.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "cattle-pipeline",
			Name:      "jenkins",
		},
		Spec: appsv1beta2.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": "jenkins"},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": "jenkins"},
					Name:   "jenkins",
				},
				Spec: corev1.PodSpec{
					InitContainers: []corev1.Container{
						{
							Name:  "jenkins-boot",
							Image: "rancher/pipeline-jenkins-boot:v1.0.0",
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "jenkins-home",
									MountPath: "/var/jenkins_home",
								},
							},
						},
					},
					Containers: []corev1.Container{
						{
							Name:  "jenkins",
							Image: "jenkins/jenkins:2.60.2-alpine",
							Env: []corev1.EnvVar{
								{
									Name:  "JENKINS_SLAVE_AGENT_PORT",
									Value: "50000",
								},
								{
									Name:  "JENKINS_HOME",
									Value: "/var/jenkins_home",
								},
							},
							Ports: []corev1.ContainerPort{
								{
									Name:          "management",
									ContainerPort: 8080,
								},
								{
									Name:          "jnlp",
									ContainerPort: 50000,
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "jenkins-home",
									MountPath: "/var/jenkins_home",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "jenkins-home",
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{},
							},
						},
					},
				},
			},
		},
	}
}

type ClusterPipelineLifecycle struct {
	cluster *config.ClusterContext
}

func (l *ClusterPipelineLifecycle) Sync(key string, obj *v3.ClusterPipeline) error {
	clusterPipeline := &v3.ClusterPipeline{
		ObjectMeta: metav1.ObjectMeta{
			Name: l.cluster.ClusterName,
		},
		Spec: v3.ClusterPipelineSpec{
			ClusterName: l.cluster.ClusterName,
		},
	}
	v3.ClusterPipelineConditionInitialized.False(clusterPipeline)
	_, err := l.cluster.Management.Management.ClusterPipelines(l.cluster.ClusterName).Create(clusterPipeline)
	if err != nil && !apierrors.IsAlreadyExists(err) {
		return err
	}
	return nil
}

func (l *ClusterPipelineLifecycle) Create(obj *v3.ClusterPipeline) (*v3.ClusterPipeline, error) {
	if obj.Spec.Deploy {
		newObj, err := v3.ClusterPipelineConditionInitialized.Once(obj, func() (runtime.Object, error) {
			err := l.deploy()
			return obj, err
		})
		return newObj.(*v3.ClusterPipeline), err
	}

	return obj, nil
}

func (l *ClusterPipelineLifecycle) Updated(obj *v3.ClusterPipeline) (*v3.ClusterPipeline, error) {
	if obj.Spec.Deploy {
		if err := l.deploy(); err != nil {
			return obj, err
		}
	} else {
		if err := l.destroy(); err != nil {
			return obj, err
		}
	}
	return obj, nil
}

func (l *ClusterPipelineLifecycle) Remove(obj *v3.ClusterPipeline) (*v3.ClusterPipeline, error) {
	if err := l.cluster.Core.Namespaces("").Delete("cattle-pipeline", &metav1.DeleteOptions{}); err != nil && !apierrors.IsAlreadyExists(err) {
		if err != nil && !apierrors.IsNotFound(err) {
			logrus.Errorf("Error occured while removing ns: %v", err)
			return obj, err
		}
	}
	return obj, nil
}

func (l *ClusterPipelineLifecycle) deploy() error {
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

	deployment := getDeployment()
	if _, err := l.cluster.Apps.Deployments("").Create(deployment); err != nil && !apierrors.IsAlreadyExists(err) {
		logrus.Errorf("Error occured while create deployment: %v", err)
		return errors.Wrapf(err, "Creating deployment")
	}

	service := getService()
	if _, err := l.cluster.Core.Services("").Create(service); err != nil && !apierrors.IsAlreadyExists(err) {
		logrus.Errorf("Error occured while create service: %v", err)
		return errors.Wrapf(err, "Creating service")
	}

	return nil
}

func (l *ClusterPipelineLifecycle) destroy() error {
	if err := l.cluster.Core.Namespaces("").Delete("cattle-pipeline", &metav1.DeleteOptions{}); err != nil && !apierrors.IsAlreadyExists(err) {
		if err != nil && !apierrors.IsNotFound(err) {
			logrus.Errorf("Error occured while removing ns: %v", err)
			return err
		}
	}
	return nil
}
