package pipelineexecution

import (
	"github.com/rancher/rancher/pkg/pipeline/utils"

	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/pkg/errors"
	"github.com/rancher/rancher/pkg/controllers/user/nslabels"
	images "github.com/rancher/rancher/pkg/image"
	"github.com/rancher/rancher/pkg/randomtoken"
	"github.com/rancher/rancher/pkg/ref"
	"github.com/rancher/rke/pki"
	mv3 "github.com/rancher/types/apis/management.cattle.io/v3"
	"github.com/rancher/types/apis/project.cattle.io/v3"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	appsv1beta2 "k8s.io/api/apps/v1beta2"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/api/networking/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/util/cert"
	"k8s.io/kubernetes/pkg/credentialprovider"
	"math/rand"
	"strconv"
	"time"
)

const projectIDFieldLabel = "field.cattle.io/projectId"

func (l *Lifecycle) deploy(obj *v3.PipelineExecution) error {
	logrus.Debug("deploy pipeline workloads and services")

	token, err := randomtoken.Generate()
	if err != nil {
		logrus.Warningf("warning generate random token got - %v, use default instead", err)
		token = utils.PipelineSecretDefaultToken
	}

	cns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: utils.PipelineNamespace,
		},
	}
	if _, err := l.namespaces.Create(cns); err != nil && !apierrors.IsAlreadyExists(err) {
		return errors.Wrapf(err, "Error create ns")
	}

	nsName := utils.GetPipelineCommonName(obj)
	_, pname := ref.Parse(obj.Spec.ProjectName)
	pns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:        nsName,
			Labels:      labels.Set(map[string]string{nslabels.ProjectIDFieldLabel: pname}),
			Annotations: map[string]string{nslabels.ProjectIDFieldLabel: obj.Spec.ProjectName},
		},
	}
	if _, err := l.namespaces.Create(pns); err != nil && !apierrors.IsAlreadyExists(err) {
		return errors.Wrapf(err, "Error create ns")
	}

	secret := getPipelineSecret(nsName, token)
	if _, err := l.secrets.Create(secret); err != nil && !apierrors.IsAlreadyExists(err) {
		return errors.Wrapf(err, "Error create secret")
	}

	secret = getRegistryCrtSecret(nsName)
	if _, err := l.secrets.Create(secret); err != nil && !apierrors.IsAlreadyExists(err) {
		return errors.Wrapf(err, "Error create secret")
	}

	sa := getServiceAccount(nsName)
	if _, err := l.serviceAccounts.Create(sa); err != nil && !apierrors.IsAlreadyExists(err) {
		return errors.Wrapf(err, "Error create service account")
	}

	np := getNetworkPolicy(nsName)
	if _, err := l.networkPolicies.Create(np); err != nil && !apierrors.IsAlreadyExists(err) {
		return errors.Wrapf(err, "Error create networkpolicy")
	}

	jenkinsService := getJenkinsService(nsName)
	if _, err := l.services.Create(jenkinsService); err != nil && !apierrors.IsAlreadyExists(err) {
		return errors.Wrapf(err, "Error create jenkins service")
	}
	jenkinsDeployment := getJenkinsDeployment(nsName)
	if _, err := l.deployments.Create(jenkinsDeployment); err != nil && !apierrors.IsAlreadyExists(err) {
		return errors.Wrapf(err, "Error create jenkins deployment")
	}

	registryService := getRegistryService(nsName)
	if _, err := l.services.Create(registryService); err != nil && !apierrors.IsAlreadyExists(err) {
		return errors.Wrapf(err, "Error create registry service")
	}
	registryDeployment := getRegistryDeployment(nsName)
	if _, err := l.deployments.Create(registryDeployment); err != nil && !apierrors.IsAlreadyExists(err) {
		return errors.Wrapf(err, "Error create registry deployment")
	}

	minioService := getMinioService(nsName)
	if _, err := l.services.Create(minioService); err != nil && !apierrors.IsAlreadyExists(err) {
		return errors.Wrapf(err, "Error create minio service")
	}
	minioDeployment := getMinioDeployment(nsName)
	if _, err := l.deployments.Create(minioDeployment); err != nil && !apierrors.IsAlreadyExists(err) {
		return errors.Wrapf(err, "Error create minio deployment")
	}

	//docker credential for local registry
	if l.reconcileRegistryCredential(obj, token); err != nil {
		return err
	}

	if err := l.reconcileProxyConfigMap(pname); err != nil {
		return err
	}
	nginxDaemonset := getProxyDeployment(nsName)
	if _, err := l.deployments.Create(nginxDaemonset); err != nil && !apierrors.IsAlreadyExists(err) {
		return errors.Wrapf(err, "Error create nginx proxy")
	}

	return l.reconcileRb(obj)
}

func getPipelineSecret(ns string, token string) *corev1.Secret {
	hashed, err := utils.BCryptHash(token)
	if err != nil {
		logrus.Warningf("warning hash registry token got - %v", err)
	}
	registryToken := fmt.Sprintf("%s:%s", utils.PipelineSecretDefaultUser, hashed)

	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: ns,
			Name:      utils.PipelineSecretName,
		},
		Data: map[string][]byte{
			utils.PipelineSecretTokenKey:         []byte(token),
			utils.PipelineSecretUserKey:          []byte(utils.PipelineSecretDefaultUser),
			utils.PipelineSecretRegistryTokenKey: []byte(registryToken),
		},
	}
}

func getRegistryCrtSecret(ns string) *corev1.Secret {
	//CACrt, CAKey, err := generateRegistryCACertAndKey(ns)
	cn := fmt.Sprintf("%s.%s", utils.RegistryName, ns)
	CACrt, CAKey, err := pki.GenerateCACertAndKey(cn)
	if err != nil {
		logrus.Error(err)
	}
	var keyBlock = &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(CAKey),
	}
	b := &bytes.Buffer{}
	if err := pem.Encode(b, keyBlock); err != nil {
		logrus.Error(err)
	}

	var crtBlock = &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: CACrt.Raw,
	}
	crtb := &bytes.Buffer{}
	if err := pem.Encode(crtb, crtBlock); err != nil {
		logrus.Error(err)
	}

	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: ns,
			Name:      utils.RegistryCrtSecretName,
		},
		Data: map[string][]byte{
			utils.RegistryCrt:   crtb.Bytes(),
			utils.RegistryCAKey: b.Bytes(),
		},
	}
}

func getRegistryCredential(projectID string, token string, hostname string) (*corev1.Secret, error) {
	_, ns := ref.Parse(projectID)
	config := credentialprovider.DockerConfigJson{
		Auths: credentialprovider.DockerConfig{
			hostname: credentialprovider.DockerConfigEntry{
				Username: utils.PipelineSecretDefaultUser,
				Password: token,
			},
		},
	}
	configJson, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      utils.DockerCredentialName,
			Namespace: ns,
			Annotations: map[string]string{
				projectIDFieldLabel: projectID,
			},
		},
		Data: map[string][]byte{
			corev1.DockerConfigJsonKey: configJson,
		},
		Type: corev1.SecretTypeDockerConfigJson,
	}, nil
}

func getServiceAccount(ns string) *corev1.ServiceAccount {
	return &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: ns,
			Name:      utils.JenkinsName,
		},
	}
}

func getRoleBindings(rbNs string, commonName string) *rbacv1.RoleBinding {
	return &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      commonName,
			Namespace: rbNs,
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     roleAdmin,
		},
		Subjects: []rbacv1.Subject{{
			Kind:      rbacv1.ServiceAccountKind,
			Namespace: commonName,
			Name:      utils.JenkinsName,
		}},
	}
}

func getClusterRoleBindings(ns string, roleName string) *rbacv1.ClusterRoleBinding {
	return &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: ns + "-" + roleName,
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     roleName,
		},
		Subjects: []rbacv1.Subject{{
			Kind:      rbacv1.ServiceAccountKind,
			Namespace: ns,
			Name:      utils.JenkinsName,
		}},
	}
}

func getJenkinsService(ns string) *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: ns,
			Name:      utils.JenkinsName,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				utils.LabelKeyApp:     utils.JenkinsName,
				utils.LabelKeyJenkins: utils.JenkinsMaster,
			},
			Ports: []corev1.ServicePort{
				{
					Name: "http",
					Port: utils.JenkinsPort,
				},
				{
					Name: "agent",
					Port: utils.JenkinsJNLPPort,
				},
			},
		},
	}
}

func getJenkinsDeployment(ns string) *appsv1beta2.Deployment {
	replicas := int32(1)
	return &appsv1beta2.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: ns,
			Name:      utils.JenkinsName,
		},
		Spec: appsv1beta2.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{utils.LabelKeyApp: utils.JenkinsName},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						utils.LabelKeyApp:     utils.JenkinsName,
						utils.LabelKeyJenkins: utils.JenkinsMaster,
					},
					Name: utils.JenkinsName,
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: utils.JenkinsName,
					Containers: []corev1.Container{
						{
							Name:  utils.JenkinsName,
							Image: images.Resolve(mv3.ToolsSystemImages.PipelineSystemImages.Jenkins),
							Env: []corev1.EnvVar{
								{
									Name: "ADMIN_PASSWORD",
									ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: &corev1.SecretKeySelector{
											LocalObjectReference: corev1.LocalObjectReference{
												Name: utils.PipelineSecretName,
											},
											Key: utils.PipelineSecretTokenKey,
										}},
								}, {
									Name: "ADMIN_USER",
									ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: &corev1.SecretKeySelector{
											LocalObjectReference: corev1.LocalObjectReference{
												Name: utils.PipelineSecretName,
											},
											Key: utils.PipelineSecretUserKey,
										}},
								}, {
									Name:  "JAVA_OPTS",
									Value: "-Xmx300m -Dhudson.slaves.NodeProvisioner.initialDelay=0 -Dhudson.slaves.NodeProvisioner.MARGIN=50 -Dhudson.slaves.NodeProvisioner.MARGIN0=0.85 -Dhudson.model.LoadStatistics.clock=2000 -Dhudson.slaves.NodeProvisioner.recurrencePeriod=2000",
								}, {
									Name:  "NAMESPACE",
									Value: ns,
								}, {
									Name: "JENKINS_POD_IP",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "status.podIP",
										},
									},
								},
							},
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									ContainerPort: utils.JenkinsPort,
								},
								{
									Name:          "agent",
									ContainerPort: utils.JenkinsJNLPPort,
								},
							},
							ReadinessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/login",
										Port: intstr.FromInt(utils.JenkinsPort),
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func getNetworkPolicy(ns string) *v1.NetworkPolicy {
	return &v1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: ns,
			Name:      utils.NetWorkPolicyName,
		},
		Spec: v1.NetworkPolicySpec{
			PodSelector: metav1.LabelSelector{
				MatchExpressions: []metav1.LabelSelectorRequirement{
					{
						Key:      utils.LabelKeyApp,
						Operator: metav1.LabelSelectorOpIn,
						Values:   []string{utils.JenkinsName, utils.MinioName},
					},
				},
			},
			Ingress: []v1.NetworkPolicyIngressRule{{}},
		},
	}
}

func getRegistryService(ns string) *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: ns,
			Name:      utils.RegistryName,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				utils.LabelKeyApp: utils.RegistryName,
			},
			Ports: []corev1.ServicePort{
				{
					Name: utils.RegistryName,
					Port: utils.RegistryPort,
				},
			},
			Type: corev1.ServiceTypeNodePort,
		},
	}
}

func getRegistryDeployment(ns string) *appsv1beta2.Deployment {
	replicas := int32(1)
	return &appsv1beta2.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: ns,
			Name:      utils.RegistryName,
		},
		Spec: appsv1beta2.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{utils.LabelKeyApp: utils.RegistryName},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{utils.LabelKeyApp: utils.RegistryName},
					Name:   utils.RegistryName,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            utils.RegistryName,
							Image:           images.Resolve(mv3.ToolsSystemImages.PipelineSystemImages.Registry),
							ImagePullPolicy: corev1.PullAlways,
							Ports: []corev1.ContainerPort{
								{
									Name:          utils.RegistryName,
									ContainerPort: utils.RegistryPort,
								},
							},
							Env: []corev1.EnvVar{
								{
									Name:  "REGISTRY_HTTP_ADDR",
									Value: "0.0.0.0:443",
								},
								{
									Name:  "REGISTRY_HTTP_TLS_CERTIFICATE",
									Value: utils.PipelineSecretRegistryCrtPath + utils.RegistryCrt,
								},
								{
									Name:  "REGISTRY_HTTP_TLS_KEY",
									Value: utils.PipelineSecretRegistryCrtPath + utils.RegistryCAKey,
								},
								{
									Name:  "REGISTRY_AUTH",
									Value: "htpasswd",
								},
								{
									Name:  "REGISTRY_AUTH_HTPASSWD_REALM",
									Value: "Registry Realm",
								},
								{
									Name:  "REGISTRY_AUTH_HTPASSWD_PATH",
									Value: utils.PipelineSecretRegistryAuthPath + utils.PipelineSecretRegistryTokenKey,
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      utils.RegistryCrtVolumeName,
									MountPath: utils.PipelineSecretRegistryCrtPath,
									ReadOnly:  true,
								},
								{
									Name:      utils.RegistryAuthVolumeName,
									MountPath: utils.PipelineSecretRegistryAuthPath,
									ReadOnly:  true,
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: utils.RegistryCrtVolumeName,
							VolumeSource: corev1.VolumeSource{
								Secret: &corev1.SecretVolumeSource{
									SecretName: utils.RegistryCrtSecretName,
								},
							},
						},
						{
							Name: utils.RegistryAuthVolumeName,
							VolumeSource: corev1.VolumeSource{
								Secret: &corev1.SecretVolumeSource{
									SecretName: utils.PipelineSecretName,
									Items: []corev1.KeyToPath{
										{
											Key:  utils.PipelineSecretRegistryTokenKey,
											Path: utils.PipelineSecretRegistryTokenKey,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (l *Lifecycle) reconcileProxyConfigMap(projectID string) error {
	exist := true
	cm, err := l.configMapLister.Get(utils.PipelineNamespace, utils.ProxyConfigMapName)
	if apierrors.IsNotFound(err) {
		exist = false
	} else if err != nil {
		return err
	}

	if !exist {
		port := getPort()
		toCreate := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: utils.PipelineNamespace,
				Name:      utils.ProxyConfigMapName,
			},
			Data: map[string]string{},
		}
		toCreate.Data[projectID] = port

		yamlMap := map[string]string{}
		portMap := map[string]string{projectID: port}
		b, err := json.Marshal(portMap)
		if err != nil {
			return err
		}
		yamlMap[utils.RegistryPortMappingKey] = string(b)
		b, err = yaml.Marshal(yamlMap)
		if err != nil {
			return err
		}
		toCreate.Data[utils.RegistryPortMappingFile] = string(b)

		_, err = l.configMaps.Create(toCreate)
		return err
	}

	if _, ok := cm.Data[projectID]; !ok {
		toUpdate := cm.DeepCopy()
		port := getPort()
		toUpdate.Data[projectID] = port
		portMap := map[string]string{}
		yamlMap := map[string]string{}
		curYaml := toUpdate.Data[utils.RegistryPortMappingFile]
		yaml.Unmarshal([]byte(curYaml), &yamlMap)
		if err != nil {
			return err
		}
		err := json.Unmarshal([]byte(yamlMap[utils.RegistryPortMappingKey]), &portMap)
		if err != nil {
			return err
		}
		portMap[projectID] = port
		b, err := json.Marshal(portMap)
		if err != nil {
			return err
		}
		yamlMap[utils.RegistryPortMappingKey] = string(b)
		b, err = yaml.Marshal(yamlMap)
		if err != nil {
			return err
		}
		toUpdate.Data[utils.RegistryPortMappingFile] = string(b)
		_, err = l.configMaps.Update(toUpdate)
		return err
	}

	return nil
}

func getPort() string {
	port := 34000
	//if portMap == nil{
	//	return strconv.Itoa(port)
	//}
	rand.Seed(time.Now().UnixNano())
	port = port + rand.Intn(1000)

	return strconv.Itoa(port)

}

func getProxyDeployment(ns string) *appsv1beta2.Deployment {
	replicas := int32(1)
	return &appsv1beta2.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: utils.PipelineNamespace,
			Name:      utils.RegistryProxyName,
		},
		Spec: appsv1beta2.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{utils.LabelKeyApp: utils.RegistryProxyName},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{utils.LabelKeyApp: utils.RegistryProxyName},
					Name:   utils.RegistryProxyName,
				},
				Spec: corev1.PodSpec{
					HostNetwork: true,
					DNSPolicy:   corev1.DNSClusterFirstWithHostNet,
					Containers: []corev1.Container{
						{
							Name:    utils.RegistryProxyName,
							Image:   "lawr/proxy", //FIXME images.Resolve(mv3.ToolsSystemImages.PipelineSystemImages.Registry),
							Command: []string{"nginx-proxy"},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "mappings",
									MountPath: "/opt",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "mappings",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{Name: utils.ProxyConfigMapName},
									Items: []corev1.KeyToPath{
										{Key: "mappings.yaml",
											Path: "mappings.yaml",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func getMinioService(ns string) *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: ns,
			Name:      utils.MinioName,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				utils.LabelKeyApp: utils.MinioName,
			},
			Ports: []corev1.ServicePort{
				{
					Name: utils.MinioName,
					Port: utils.MinioPort,
				},
			},
		},
	}
}

func getMinioDeployment(ns string) *appsv1beta2.Deployment {
	replicas := int32(1)
	return &appsv1beta2.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: ns,
			Name:      utils.MinioName,
		},
		Spec: appsv1beta2.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{utils.LabelKeyApp: utils.MinioName},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{utils.LabelKeyApp: utils.MinioName},
					Name:   utils.MinioName,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            utils.MinioName,
							Image:           images.Resolve(mv3.ToolsSystemImages.PipelineSystemImages.Minio),
							ImagePullPolicy: corev1.PullAlways,
							Args:            []string{"server", "/data"},
							Env: []corev1.EnvVar{
								{
									Name: "MINIO_SECRET_KEY",
									ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: &corev1.SecretKeySelector{
											LocalObjectReference: corev1.LocalObjectReference{
												Name: utils.PipelineSecretName,
											},
											Key: utils.PipelineSecretTokenKey,
										}},
								}, {
									Name: "MINIO_ACCESS_KEY",
									ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: &corev1.SecretKeySelector{
											LocalObjectReference: corev1.LocalObjectReference{
												Name: utils.PipelineSecretName,
											},
											Key: utils.PipelineSecretUserKey,
										}},
								},
							},
							Ports: []corev1.ContainerPort{
								{
									Name:          utils.MinioName,
									ContainerPort: utils.MinioPort,
								},
							},
						},
					},
				},
			},
		},
	}
}

func generateRegistryCACertAndKey(ns string) (*x509.Certificate, *rsa.PrivateKey, error) {
	cn := fmt.Sprintf("%s.%s", utils.RegistryName, ns)
	rootKey, err := cert.NewPrivateKey()
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to generate private key for CA certificate: %v", err)
	}
	caConfig := cert.Config{
		CommonName: cn,
		AltNames: cert.AltNames{
			DNSNames: []string{cn},
		},
	}
	kubeCACert, err := cert.NewSelfSignedCACert(caConfig, rootKey)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to generate CA certificate: %v", err)
	}

	return kubeCACert, rootKey, nil
}

func (l *Lifecycle) reconcileRegistryCredential(obj *v3.PipelineExecution, token string) error {
	cm, err := l.configMapLister.Get(utils.PipelineNamespace, utils.ProxyConfigMapName)
	if err != nil {
		return err
	}
	_, projectID := ref.Parse(obj.Spec.ProjectName)
	port, ok := cm.Data[projectID]
	if !ok || port == "" {
		return errors.New("Found no port for local registry")
	}
	regHostname := "127.0.0.1:" + port
	dockerCredential, err := getRegistryCredential(obj.Spec.ProjectName, token, regHostname)
	if err != nil {
		return err
	}
	if _, err := l.secrets.Create(dockerCredential); err != nil && !apierrors.IsAlreadyExists(err) {
		return errors.Wrapf(err, "Error create credential for local registry")
	}
	return nil
}
