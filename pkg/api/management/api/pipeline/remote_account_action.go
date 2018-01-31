package pipeline

import (
	"github.com/rancher/norman/types"
	"github.com/rancher/types/config"
	"github.com/sirupsen/logrus"
)

type RemoteAccountHandler struct {
	Management config.ManagementContext
}

func RemoteAccountCollectionFormatter(apiContext *types.APIContext, collection *types.GenericCollection) {
	collection.AddAction(apiContext, "auth")
}

func RemoteAccountFormatter(apiContext *types.APIContext, resource *types.RawResource) {
	resource.AddAction(apiContext, "refreshrepos")
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
	/*
		//TODO test cluster name
		clusterName := "local"
		clusterPipeline, err := h.Management.Management.ClusterPipelines("").Get(clusterName, metav1.GetOptions{})
		if err != nil {
			return err
		}
		h.Management.Management.RemoteAccounts("").Get("")
		booter.New()
	*/
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
