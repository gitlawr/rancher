package client

const (
	ClusterSpecType                                     = "clusterSpec"
	ClusterSpecFieldDefaultClusterRoleForProjectMembers = "defaultClusterRoleForProjectMembers"
	ClusterSpecFieldDefaultPodSecurityPolicyTemplateID  = "defaultPodSecurityPolicyTemplateId"
	ClusterSpecFieldDescription                         = "description"
	ClusterSpecFieldDesiredAgentImage                   = "desiredAgentImage"
	ClusterSpecFieldDisplayName                         = "displayName"
	ClusterSpecFieldDockerRootDir                       = "dockerRootDir"
	ClusterSpecFieldEnableNetworkPolicy                 = "enableNetworkPolicy"
	ClusterSpecFieldGenericEngineConfig                 = "genericEngineConfig"
	ClusterSpecFieldImportedConfig                      = "importedConfig"
	ClusterSpecFieldInternal                            = "internal"
	ClusterSpecFieldRancherKubernetesEngineConfig       = "rancherKubernetesEngineConfig"
)

type ClusterSpec struct {
	DefaultClusterRoleForProjectMembers string                         `json:"defaultClusterRoleForProjectMembers,omitempty" yaml:"defaultClusterRoleForProjectMembers,omitempty"`
	DefaultPodSecurityPolicyTemplateID  string                         `json:"defaultPodSecurityPolicyTemplateId,omitempty" yaml:"defaultPodSecurityPolicyTemplateId,omitempty"`
	Description                         string                         `json:"description,omitempty" yaml:"description,omitempty"`
	DesiredAgentImage                   string                         `json:"desiredAgentImage,omitempty" yaml:"desiredAgentImage,omitempty"`
	DisplayName                         string                         `json:"displayName,omitempty" yaml:"displayName,omitempty"`
	DockerRootDir                       string                         `json:"dockerRootDir,omitempty" yaml:"dockerRootDir,omitempty"`
	EnableNetworkPolicy                 *bool                          `json:"enableNetworkPolicy,omitempty" yaml:"enableNetworkPolicy,omitempty"`
	GenericEngineConfig                 string                         `json:"genericEngineConfig,omitempty" yaml:"genericEngineConfig,omitempty"`
	ImportedConfig                      *ImportedConfig                `json:"importedConfig,omitempty" yaml:"importedConfig,omitempty"`
	Internal                            bool                           `json:"internal,omitempty" yaml:"internal,omitempty"`
	RancherKubernetesEngineConfig       *RancherKubernetesEngineConfig `json:"rancherKubernetesEngineConfig,omitempty" yaml:"rancherKubernetesEngineConfig,omitempty"`
}
