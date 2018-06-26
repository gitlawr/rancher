package gitlab

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/rancher/norman/store/subtype"
	"github.com/rancher/norman/types"
	"github.com/rancher/rancher/pkg/pipeline/remote/model"
	mv3 "github.com/rancher/types/apis/management.cattle.io/v3"
	"github.com/rancher/types/apis/project.cattle.io/v3"
	"github.com/rancher/types/apis/project.cattle.io/v3/schema"
	"github.com/rancher/types/client/project/v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type GlProvider struct {
	SourceCodeProviderConfigs  v3.SourceCodeProviderConfigInterface
	SourceCodeCredentials      v3.SourceCodeCredentialInterface
	SourceCodeCredentialLister v3.SourceCodeCredentialLister
	SourceCodeRepositories     v3.SourceCodeRepositoryInterface
	SourceCodeRepositoryLister v3.SourceCodeRepositoryLister
	Pipelines                  v3.PipelineInterface
	PipelineLister             v3.PipelineLister
	PipelineExecutions         v3.PipelineExecutionInterface
	PipelineExecutionLister    v3.PipelineExecutionLister

	AuthConfigs mv3.AuthConfigInterface
}

func (g *GlProvider) CustomizeSchemas(schemas *types.Schemas) {

	scpConfigBaseSchema := schemas.Schema(&schema.Version, client.SourceCodeProviderConfigType)
	configSchema := schemas.Schema(&schema.Version, client.GitlabPipelineConfigType)
	configSchema.ActionHandler = g.ActionHandler
	configSchema.Formatter = g.Formatter
	configSchema.Store = subtype.NewSubTypeStore(client.GitlabPipelineConfigType, scpConfigBaseSchema.Store)

	providerBaseSchema := schemas.Schema(&schema.Version, client.SourceCodeProviderType)
	providerSchema := schemas.Schema(&schema.Version, client.GitlabProviderType)
	providerSchema.Formatter = g.providerFormatter
	providerSchema.ActionHandler = g.providerActionHandler
	providerSchema.Store = subtype.NewSubTypeStore(client.GitlabProviderType, providerBaseSchema.Store)
}

func (g *GlProvider) GetName() string {
	return model.GitlabType
}

func (g *GlProvider) TransformToSourceCodeProvider(config map[string]interface{}) map[string]interface{} {
	p := transformToSourceCodeProvider(config)
	return p
}

func transformToSourceCodeProvider(config map[string]interface{}) map[string]interface{} {
	result := map[string]interface{}{}

	if m, ok := config["metadata"].(map[string]interface{}); ok {
		result["id"] = fmt.Sprintf("%v:%v", m["namespace"], m["name"])
	}
	if t, ok := config["type"].(string); ok && t != "" {
		result["type"] = "gitlabProvider"
	}
	if t, ok := config["projectName"].(string); ok && t != "" {
		result["projectId"] = t
	}
	result["redirectUrl"] = formGitlabRedirectURLFromMap(config)

	return result
}

func (g *GlProvider) GetProviderConfig(projectID string) (interface{}, error) {
	scpConfigObj, err := g.SourceCodeProviderConfigs.ObjectClient().UnstructuredClient().GetNamespaced(projectID, model.GitlabType, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve GitlabConfig, error: %v", err)
	}

	u, ok := scpConfigObj.(runtime.Unstructured)
	if !ok {
		return nil, fmt.Errorf("failed to retrieve GitlabConfig, cannot read k8s Unstructured data")
	}
	storedGitlabPipelineConfigMap := u.UnstructuredContent()

	storedGitlabPipelineConfig := &v3.GitlabPipelineConfig{}
	mapstructure.Decode(storedGitlabPipelineConfigMap, storedGitlabPipelineConfig)

	metadataMap, ok := storedGitlabPipelineConfigMap["metadata"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("failed to retrieve GitlabConfig metadata, cannot read k8s Unstructured data")
	}

	typemeta := &metav1.ObjectMeta{}
	mapstructure.Decode(metadataMap, typemeta)
	storedGitlabPipelineConfig.ObjectMeta = *typemeta
	storedGitlabPipelineConfig.APIVersion = "project.cattle.io/v3"
	storedGitlabPipelineConfig.Kind = v3.SourceCodeProviderConfigGroupVersionKind.Kind
	return storedGitlabPipelineConfig, nil
}
