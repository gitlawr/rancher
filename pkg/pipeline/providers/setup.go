package providers

import (
	"github.com/rancher/norman/types"
	"github.com/rancher/rancher/pkg/pipeline/providers/github"
	"github.com/rancher/rancher/pkg/pipeline/providers/gitlab"
	"github.com/rancher/rancher/pkg/pipeline/remote/model"
	"github.com/rancher/types/apis/project.cattle.io/v3/schema"
	"github.com/rancher/types/client/project/v3"
	"github.com/rancher/types/config"
)

var (
	providers       = make(map[string]SourceCodeProvider)
	providersByType = make(map[string]SourceCodeProvider)
)

var sourceCodeProviderConfigTypes = []string{
	client.GithubPipelineConfigType,
	client.GitlabPipelineConfigType,
}

func SetupSourceCodeProviderConfig(management *config.ScaledContext, schemas *types.Schemas) {
	configure(management)

	providerBaseSchema := schemas.Schema(&schema.Version, client.SourceCodeProviderType)
	setSourceCodeProviderStore(providerBaseSchema, management)

	for _, scpSubtype := range sourceCodeProviderConfigTypes {
		providersByType[scpSubtype].CustomizeSchemas(schemas)
	}
}

func configure(management *config.ScaledContext) {

	ghProvider := &github.GhProvider{
		SourceCodeProviderConfigs:  management.Project.SourceCodeProviderConfigs(""),
		SourceCodeCredentialLister: management.Project.SourceCodeCredentials("").Controller().Lister(),
		SourceCodeCredentials:      management.Project.SourceCodeCredentials(""),
		SourceCodeRepositories:     management.Project.SourceCodeRepositories(""),
		SourceCodeRepositoryLister: management.Project.SourceCodeRepositories("").Controller().Lister(),
		Pipelines:                  management.Project.Pipelines(""),
		PipelineLister:             management.Project.Pipelines("").Controller().Lister(),
		PipelineExecutions:         management.Project.PipelineExecutions(""),
		PipelineExecutionLister:    management.Project.PipelineExecutions("").Controller().Lister(),

		AuthConfigs: management.Management.AuthConfigs(""),
	}
	providers[model.GithubType] = ghProvider
	providersByType[client.GithubPipelineConfigType] = ghProvider

	glProvider := &gitlab.GlProvider{
		SourceCodeProviderConfigs:  management.Project.SourceCodeProviderConfigs(""),
		SourceCodeCredentialLister: management.Project.SourceCodeCredentials("").Controller().Lister(),
		SourceCodeCredentials:      management.Project.SourceCodeCredentials(""),
		SourceCodeRepositories:     management.Project.SourceCodeRepositories(""),
		SourceCodeRepositoryLister: management.Project.SourceCodeRepositories("").Controller().Lister(),
		Pipelines:                  management.Project.Pipelines(""),
		PipelineLister:             management.Project.Pipelines("").Controller().Lister(),
		PipelineExecutions:         management.Project.PipelineExecutions(""),
		PipelineExecutionLister:    management.Project.PipelineExecutions("").Controller().Lister(),

		AuthConfigs: management.Management.AuthConfigs(""),
	}
	providers[model.GitlabType] = glProvider
	providersByType[client.GitlabPipelineConfigType] = glProvider

}
