package pipeline

import (
	"github.com/pkg/errors"
	"github.com/rancher/norman/api/access"
	"github.com/rancher/norman/types"
	"github.com/rancher/rancher/pkg/pipeline/remote/booter"
	"github.com/rancher/types/apis/management.cattle.io/v3"
	"github.com/rancher/types/client/management/v3"
	"github.com/rancher/types/config"
	"github.com/sirupsen/logrus"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"strings"
)

type RemoteAccountHandler struct {
	Management config.ManagementContext
}

func RemoteAccountCollectionFormatter(apiContext *types.APIContext, collection *types.GenericCollection) {
	collection.AddAction(apiContext, "auth")
}

func RemoteAccountFormatter(apiContext *types.APIContext, resource *types.RawResource) {
	resource.AddAction(apiContext, "refreshrepos")
	resource.Links["repos"] = apiContext.URLBuilder.Link("repos", resource)
}

func (h RemoteAccountHandler) LinkHandler(apiContext *types.APIContext) error {
	logrus.Infof("get id:%s", apiContext.ID)

	_, err := h.Management.Management.GitRepoCaches("").Get(apiContext.ID, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		return h.refreshrepos(apiContext)
	} else if err != nil {
		return err
	}
	data := map[string]interface{}{}
	if err := access.ByID(apiContext, apiContext.Version, "gitrepocache", apiContext.ID, &data); err != nil {
		return err
	}

	apiContext.WriteResponse(http.StatusOK, data)
	return nil
}
func (h *RemoteAccountHandler) ActionHandler(actionName string, action *types.Action, apiContext *types.APIContext) error {
	logrus.Debugf("do remote account action:%s", actionName)

	switch actionName {
	case "refreshrepos":
		return h.refreshrepos(apiContext)
	}
	return nil
}

func (h *RemoteAccountHandler) refreshrepos(apiContext *types.APIContext) error {
	logrus.Infof("get id:%s", apiContext.ID)
	account, err := h.Management.Management.RemoteAccounts("").Get(apiContext.ID, metav1.GetOptions{})
	if err != nil {
		return err
	}
	parts := strings.SplitN(apiContext.ID, "-", 2)
	if len(parts) < 2 {
		return errors.New("invalid remote account")
	}
	remoteType := parts[0]

	mockConfig := v3.ClusterPipeline{
		Spec: v3.ClusterPipelineSpec{
			GithubConfig: &v3.GithubConfig{},
		},
	}
	remote, err := booter.New(mockConfig, remoteType)
	if err != nil {
		return err
	}
	repos, err := remote.Repos(account)
	if err != nil {
		return err
	}
	repocache := &v3.GitRepoCache{
		ObjectMeta: metav1.ObjectMeta{
			Name: account.Name,
		},
		Spec: v3.GitRepoCacheSpec{
			RemoteType:   remoteType,
			Repositories: repos,
		},
	}
	_, err = h.Management.Management.GitRepoCaches("").Get(apiContext.ID, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		repocache, err = h.Management.Management.GitRepoCaches("").Create(repocache)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	repocache, err = h.Management.Management.GitRepoCaches("").Update(repocache)
	if err != nil {
		return err
	}
	data := map[string]interface{}{}
	if err := access.ByID(apiContext, apiContext.Version, client.GitRepoCacheType, apiContext.ID, &data); err != nil {
		return err
	}

	apiContext.WriteResponse(http.StatusOK, data)
	return nil
}

func getRedirectURL(apiContext *types.APIContext) string {
	return "https://example.com/redirect"
}

/*
func (h *RemoteAccountHandler) oauth2(apiContext *types.APIContext) error {

	logrus.Infof("get id:%s", apiContext.ID)
	//TODO test cluster name
	clusterName := "local"

	clusterPipeline := v3.ClusterPipeline{}
	requestBytes, err := ioutil.ReadAll(apiContext.Request.Body)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(requestBytes, &clusterPipeline); err != nil {
		return err
	}
	var code, scmType, clientId, clientSecret, redirectURL, scheme, host string

	if clusterPipeline.Spec.GithubConfig != nil {
		scmType = "github"
	}

	//TODO test cluster name
	if apiContext.Request.FormValue("clusterName") != "" {
		clusterName = apiContext.Request.FormValue("clusterName")
	}

	if apiContext.Request.FormValue("code") != "" {
		code = apiContext.Request.FormValue("code")
	}
	if apiContext.Request.FormValue("redirectUrl") != "" {
		code = apiContext.Request.FormValue("redirectUrl")
	}

	if clientId == "" || clientSecret == "" || redirectURL == "" {
		clusterPipeline, err := h.Management.Management.ClusterPipelines("").Get(clusterName, metav1.GetOptions{})
		if err != nil {
			return err
		}
		if clusterPipeline.Spec.GithubConfig == nil {
			return errors.New("github not configured")
		}

		clientId = clusterPipeline.Spec.GithubConfig.ClientId
		clientSecret = clusterPipeline.Spec.GithubConfig.ClientSecret

		remote, err := booter.New(*clusterPipeline, scmType)
		if err != nil {
			return err
		}
		account, err := remote.Login(redirectURL, code)
		if err != nil {
			return err
		}
		if _, err := h.Management.Management.RemoteAccounts("").Create(account); err != nil {
			return err
		}

	} else {
		clusterPipeline := &v3.ClusterPipeline{
			ObjectMeta: metav1.ObjectMeta{
				Name: clusterName,
			},
			Spec: v3.ClusterPipelineSpec{
				ClusterName: clusterName,
				GithubConfig: &v3.GithubConfig{
					Scheme:       scheme,
					Host:         host,
					ClientId:     clientId,
					ClientSecret: clientSecret,
				},
			},
		}

		remote, err := booter.New(*clusterPipeline, scmType)
		if err != nil {
			return err
		}
		if _, err := remote.Login(redirectURL, code); err != nil {
			return err
		}
		if _, err := h.Management.Management.ClusterPipelines("").Create(clusterPipeline); err != nil {
			return err
		}
	}
	apiContext.WriteResponse(200, clusterPipeline)
	return nil
}
*/
