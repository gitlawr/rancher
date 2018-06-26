package pipeline

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/rancher/norman/api/access"
	"github.com/rancher/norman/httperror"
	"github.com/rancher/norman/types"
	"github.com/rancher/rancher/pkg/pipeline/providers"
	"github.com/rancher/rancher/pkg/pipeline/remote"
	"github.com/rancher/rancher/pkg/pipeline/remote/model"
	"github.com/rancher/rancher/pkg/pipeline/utils"
	"github.com/rancher/rancher/pkg/ref"
	"github.com/rancher/types/apis/project.cattle.io/v3"
	"github.com/rancher/types/client/project/v3"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/labels"
	"net/http"
	"time"
)

type Handler struct {
	Pipelines                  v3.PipelineInterface
	PipelineLister             v3.PipelineLister
	PipelineExecutions         v3.PipelineExecutionInterface
	SourceCodeCredentialLister v3.SourceCodeCredentialLister
}

func Formatter(apiContext *types.APIContext, resource *types.RawResource) {
	resource.AddAction(apiContext, "run")
	resource.AddAction(apiContext, "reload")
	resource.AddAction(apiContext, "pushconfig")
	resource.Links["export"] = apiContext.URLBuilder.Link("export", resource)
	resource.Links["configs"] = apiContext.URLBuilder.Link("configs", resource)
	resource.Links["yaml"] = apiContext.URLBuilder.Link("yaml", resource)
	resource.Links["branches"] = apiContext.URLBuilder.Link("branches", resource)
}

func (h *Handler) LinkHandler(apiContext *types.APIContext, next types.RequestHandler) error {

	ns, name := ref.Parse(apiContext.ID)
	pipeline, err := h.PipelineLister.Get(ns, name)
	if err != nil {
		return err
	}

	if apiContext.Link == "export" {
		pipelineConfig, err := h.getPipelineConfig(apiContext)
		if err != nil {
			return err
		}
		content, err := utils.PipelineConfigToYaml(pipelineConfig)
		if err != nil {
			return err
		}
		fileName := fmt.Sprintf("pipeline-%s.yaml", pipeline.Spec.DisplayName)
		apiContext.Response.Header().Add("Content-Disposition", "attachment; filename="+fileName)
		http.ServeContent(apiContext.Response, apiContext.Request, fileName, time.Now(), bytes.NewReader(content))
		return nil
	} else if apiContext.Link == "yaml" {
		pipelineConfig, err := h.getPipelineConfig(apiContext)
		if err != nil {
			return err
		}
		content, err := utils.PipelineConfigToYaml(pipelineConfig)
		if err != nil {
			return err
		}
		_, err = apiContext.Response.Write([]byte(content))
		return err
	} else if apiContext.Link == "configs" {
		return h.getPipelineConfigs(apiContext)
	} else if apiContext.Link == "branches" {
		return h.getValidBranches(apiContext)
	}

	return httperror.NewAPIError(httperror.NotFound, "Link not found")
}

func (h *Handler) getPipelineConfig(apiContext *types.APIContext) (*v3.PipelineConfig, error) {
	ns, name := ref.Parse(apiContext.ID)
	pipeline, err := h.PipelineLister.Get(ns, name)
	if err != nil {
		return nil, err
	}
	branch := apiContext.Request.URL.Query().Get("branch")
	if branch == "" {
		return nil, httperror.NewAPIError(httperror.InvalidOption, "Branch is not specified")
	}

	return providers.GetPipelineConfigByBranch(h.SourceCodeCredentialLister, pipeline, branch)
}

func (h *Handler) ActionHandler(actionName string, action *types.Action, apiContext *types.APIContext) error {

	switch actionName {
	case "reload":
		return h.reload(apiContext)
	case "run":
		return h.run(apiContext)
	case "pushconfig":
		return h.pushConfig(apiContext)
	}
	return httperror.NewAPIError(httperror.InvalidAction, "unsupported action")
}

func (h *Handler) run(apiContext *types.APIContext) error {

	ns, name := ref.Parse(apiContext.ID)
	pipeline, err := h.PipelineLister.Get(ns, name)
	if err != nil {
		return err
	}
	runPipelineInput := v3.RunPipelineInput{}
	requestBytes, err := ioutil.ReadAll(apiContext.Request.Body)
	if err != nil {
		return err
	}
	if string(requestBytes) != "" {
		if err := json.Unmarshal(requestBytes, &runPipelineInput); err != nil {
			return err
		}
	}

	branch := runPipelineInput.Branch
	if branch == "" {
		return httperror.NewAPIError(httperror.InvalidBodyContent, "Error branch is not specified for the pipeline to run")
	}

	userName := apiContext.Request.Header.Get("Impersonate-User")
	pipelineConfig, err := providers.GetPipelineConfigByBranch(h.SourceCodeCredentialLister, pipeline, branch)
	if err != nil {
		return err
	}

	if pipelineConfig == nil {
		return fmt.Errorf("find no pipeline config to run in the branch")
	}

	info, err := h.getBuildInfoByBranch(pipeline, branch)
	if err != nil {
		return err
	}
	info.TriggerType = utils.TriggerTypeUser
	info.TriggerUserName = userName
	execution, err := utils.GenerateExecution(h.PipelineExecutions, pipeline, pipelineConfig, info)
	if err != nil {
		return err
	}

	if execution == nil {
		return errors.New("condition is not match, no build is triggered")
	}

	data := map[string]interface{}{}
	if err := access.ByID(apiContext, apiContext.Version, client.PipelineExecutionType, ns+":"+execution.Name, &data); err != nil {
		return err
	}

	apiContext.WriteResponse(http.StatusOK, data)
	return err
}

func (h *Handler) pushConfig(apiContext *types.APIContext) error {
	ns, name := ref.Parse(apiContext.ID)
	pipeline, err := h.PipelineLister.Get(ns, name)
	if err != nil {
		return err
	}

	pushConfigInput := v3.PushPipelineConfigInput{}
	requestBytes, err := ioutil.ReadAll(apiContext.Request.Body)
	if err != nil {
		return err
	}
	if string(requestBytes) != "" {
		if err := json.Unmarshal(requestBytes, &pushConfigInput); err != nil {
			return err
		}
	}

	//use current user's auth to do the push
	userName := apiContext.Request.Header.Get("Impersonate-User")
	creds, err := h.SourceCodeCredentialLister.List(userName, labels.Everything())
	if err != nil {
		return err
	}
	accessToken := ""
	sourceCodeType := model.GithubType
	for _, cred := range creds {
		if cred.Spec.ProjectName == pipeline.Spec.ProjectName && !cred.Status.Logout {
			accessToken = cred.Spec.AccessToken
			sourceCodeType = cred.Spec.SourceCodeType
		}
	}

	_, projID := ref.Parse(pipeline.Spec.ProjectName)
	scpConfig, err := providers.GetSourceCodeProviderConfig(sourceCodeType, projID)
	if err != nil {
		return err
	}
	remote, err := remote.New(scpConfig)
	if err != nil {
		return err
	}

	for branch, config := range pushConfigInput.Configs {
		content, err := utils.PipelineConfigToYaml(&config)
		if err != nil {
			return err
		}
		if err := remote.SetPipelineFileInRepo(pipeline.Spec.RepositoryURL, branch, accessToken, content); err != nil {
			if apierr, ok := err.(*httperror.APIError); ok && apierr.Code.Status == http.StatusNotFound {
				//github returns 404 for unauth request to prevent leakage of private repos
				return httperror.NewAPIError(httperror.Unauthorized, "current git account is unauthorized for the action")
			}
			return err
		}
	}

	data := map[string]interface{}{}
	if err := access.ByID(apiContext, apiContext.Version, apiContext.Type, apiContext.ID, &data); err != nil {
		return err
	}

	apiContext.WriteResponse(http.StatusOK, data)
	return nil
}

func (h *Handler) reload(apiContext *types.APIContext) error {
	ns, name := ref.Parse(apiContext.ID)
	pipeline, err := h.PipelineLister.Get(ns, name)
	if err != nil {
		return err
	}

	reloadPipelineInput := v3.ReloadPipelineInput{}
	requestBytes, err := ioutil.ReadAll(apiContext.Request.Body)
	if err != nil {
		return err
	}
	if string(requestBytes) != "" {
		if err := json.Unmarshal(requestBytes, &reloadPipelineInput); err != nil {
			return err
		}
	}

	branch := reloadPipelineInput.Branch
	if branch != "" {
		delete(pipeline.Spec.UnSyncConfigs, branch)
	} else {
		pipeline.Spec.UnSyncConfigs = nil
	}
	if _, err = h.Pipelines.Update(pipeline); err != nil {
		return err
	}
	data := map[string]interface{}{}
	if err := access.ByID(apiContext, apiContext.Version, apiContext.Type, apiContext.ID, &data); err != nil {
		return err
	}

	apiContext.WriteResponse(http.StatusOK, data)
	return nil
}

func (h *Handler) getBuildInfoByBranch(pipeline *v3.Pipeline, branch string) (*model.BuildInfo, error) {

	credentialName := pipeline.Spec.SourceCodeCredentialName
	repoURL := pipeline.Spec.RepositoryURL
	accessToken := ""
	sourceCodeType := model.GithubType

	if credentialName != "" {
		ns, name := ref.Parse(credentialName)
		credential, err := h.SourceCodeCredentialLister.Get(ns, name)
		if err != nil {
			return nil, err
		}
		sourceCodeType = credential.Spec.SourceCodeType
		accessToken = credential.Spec.AccessToken
	}
	_, projID := ref.Parse(pipeline.Spec.ProjectName)
	scpConfig, err := providers.GetSourceCodeProviderConfig(sourceCodeType, projID)
	if err != nil {
		return nil, err
	}
	remote, err := remote.New(scpConfig)
	if err != nil {
		return nil, err
	}
	info, err := remote.GetHeadInfo(repoURL, branch, accessToken)
	if err != nil {
		return nil, err
	}
	return info, nil

}

func (h *Handler) getValidBranches(apiContext *types.APIContext) error {
	ns, name := ref.Parse(apiContext.ID)
	pipeline, err := h.PipelineLister.Get(ns, name)
	if err != nil {
		return err
	}

	accessKey := ""
	sourceCodeType := model.GithubType
	if pipeline.Spec.SourceCodeCredentialName != "" {
		ns, name = ref.Parse(pipeline.Spec.SourceCodeCredentialName)
		cred, err := h.SourceCodeCredentialLister.Get(ns, name)
		if err != nil {
			return err
		}
		accessKey = cred.Spec.AccessToken
		sourceCodeType = cred.Spec.SourceCodeType
	}

	_, projID := ref.Parse(pipeline.Spec.ProjectName)
	scpConfig, err := providers.GetSourceCodeProviderConfig(sourceCodeType, projID)
	if err != nil {
		return err
	}
	remote, err := remote.New(scpConfig)
	if err != nil {
		return err
	}

	validBranches := map[string]bool{}

	branches, err := remote.GetBranches(pipeline.Spec.RepositoryURL, accessKey)
	if err != nil {
		return err
	}
	for _, b := range branches {
		content, err := remote.GetPipelineFileInRepo(pipeline.Spec.RepositoryURL, b, accessKey)
		if err != nil {
			return err
		}
		if content != nil {
			validBranches[b] = true
		}
	}
	for b := range pipeline.Spec.UnSyncConfigs {
		validBranches[b] = true
	}

	result := []string{}
	for b := range validBranches {
		result = append(result, b)
	}

	bytes, err := json.Marshal(result)
	if err != nil {
		return err
	}
	apiContext.Response.Write(bytes)
	return nil
}

func (h *Handler) getPipelineConfigs(apiContext *types.APIContext) error {
	// get configs in all branches
	ns, name := ref.Parse(apiContext.ID)
	pipeline, err := h.PipelineLister.Get(ns, name)
	if err != nil {
		return err
	}

	accessToken := ""
	sourceCodeType := model.GithubType
	if pipeline.Spec.SourceCodeCredentialName != "" {
		ns, name = ref.Parse(pipeline.Spec.SourceCodeCredentialName)
		cred, err := h.SourceCodeCredentialLister.Get(ns, name)
		if err != nil {
			return err
		}
		sourceCodeType = cred.Spec.SourceCodeType
		accessToken = cred.Spec.AccessToken
	}

	_, projID := ref.Parse(pipeline.Spec.ProjectName)
	scpConfig, err := providers.GetSourceCodeProviderConfig(sourceCodeType, projID)
	if err != nil {
		return err
	}
	remote, err := remote.New(scpConfig)
	if err != nil {
		return err
	}

	m := map[string]interface{}{}

	branch := apiContext.Request.URL.Query().Get("branch")
	if branch != "" {
		content, err := remote.GetPipelineFileInRepo(pipeline.Spec.RepositoryURL, branch, accessToken)
		if err != nil {
			return err
		}
		if content != nil {
			spec, err := utils.PipelineConfigFromYaml(content)
			if err != nil {
				return errors.Wrapf(err, "Error fetching pipeline config in Branch '%s'", branch)
			}
			m[branch] = spec
		}
	} else {
		branches, err := remote.GetBranches(pipeline.Spec.RepositoryURL, accessToken)
		if err != nil {
			return err
		}
		for _, b := range branches {
			content, err := remote.GetPipelineFileInRepo(pipeline.Spec.RepositoryURL, b, accessToken)
			if err != nil {
				return err
			}
			if content != nil {
				spec, err := utils.PipelineConfigFromYaml(content)
				if err != nil {
					return errors.Wrapf(err, "Error fetching pipeline config in Branch '%s'", b)
				}
				m[b] = spec
			} else {
				m[b] = nil
			}
		}
	}

	bytes, err := json.Marshal(m)
	if err != nil {
		return err
	}
	apiContext.Response.Write(bytes)
	return nil
}
