package pipeline

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/rancher/norman/types"
	"github.com/rancher/rancher/pkg/pipeline/remote/booter"
	"github.com/rancher/types/apis/management.cattle.io/v3"
	"github.com/rancher/types/config"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

type ClusterPipelineHandler struct {
	Management config.ManagementContext
}

func ClusterPipelineFormatter(apiContext *types.APIContext, resource *types.RawResource) {
	resource.AddAction(apiContext, "deploy")
	resource.AddAction(apiContext, "destroy")
	resource.AddAction(apiContext, "auth")
}

func (h *ClusterPipelineHandler) ActionHandler(actionName string, action *types.Action, apiContext *types.APIContext) error {

	logrus.Infof("get id:%s", apiContext.ID)
	parts := strings.Split(apiContext.ID, ":")
	if len(parts) <= 1 {
		return errors.New("invalid ID")
	}
	ns := parts[0]
	id := parts[1]
	client := h.Management.Management.ClusterPipelines(ns)
	clusterPipeline, err := client.Get(id, metav1.GetOptions{})
	if err != nil {
		return err
	}
	logrus.Infof("do cluster pipeline action:%s", actionName)

	switch actionName {
	case "deploy":
		clusterPipeline.Spec.Deploy = true
	case "destroy":
		clusterPipeline.Spec.Deploy = false
	case "auth":
		return h.auth(apiContext)
	}
	_, err = client.Update(clusterPipeline)
	return err
}

func (h *ClusterPipelineHandler) auth(apiContext *types.APIContext) error {

	parts := strings.Split(apiContext.ID, ":")
	if len(parts) <= 1 {
		return errors.New("invalid ID")
	}
	ns := parts[0]
	id := parts[1]

	requestBody := make(map[string]interface{})
	requestBytes, err := ioutil.ReadAll(apiContext.Request.Body)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(requestBytes, &requestBody); err != nil {
		return err
	}
	var code, remoteType, clientID, clientSecret, redirectURL, scheme, host string

	if requestBody["code"] != nil {
		code = requestBody["code"].(string)
	}
	if requestBody["clientId"] != nil {
		clientID = requestBody["clientId"].(string)
	}
	if requestBody["clientSecret"] != nil {
		clientSecret = requestBody["clientSecret"].(string)
	}
	if requestBody["redirectUrl"] != nil {
		redirectURL = requestBody["redirectUrl"].(string)
	}
	if requestBody["scheme"] != nil {
		scheme = requestBody["scheme"].(string)
	}
	if requestBody["host"] != nil {
		host = requestBody["host"].(string)
	}
	if requestBody["type"] != nil {
		remoteType = requestBody["type"].(string)
	}

	clusterPipeline, err := h.Management.Management.ClusterPipelines(ns).Get(id, metav1.GetOptions{})
	if err != nil {
		return err
	}
	if clientID == "" || clientSecret == "" {

		if clusterPipeline.Spec.GithubConfig == nil {
			return errors.New("github not configured")
		}

		//oauth and add user
		userName := apiContext.Request.Header.Get("Impersonate-User")
		if err := h.auth_add_account(clusterPipeline, remoteType, userName, redirectURL, code); err != nil {
			return err
		}

	} else {
		if remoteType == "github" && clusterPipeline.Spec.GithubConfig == nil {
			clusterPipeline.Spec.GithubConfig = &v3.GithubConfig{
				Scheme:       scheme,
				Host:         host,
				ClientId:     clientID,
				ClientSecret: clientSecret,
			}
		}
		//oauth and add user
		userName := apiContext.Request.Header.Get("Impersonate-User")
		if err := h.auth_add_account(clusterPipeline, remoteType, userName, redirectURL, code); err != nil {
			return err
		}
		//update cluster pipeline config
		if _, err := h.Management.Management.ClusterPipelines(ns).Update(clusterPipeline); err != nil {
			return err
		}
	}
	apiContext.WriteResponse(200, clusterPipeline)
	return nil
}

func (h *ClusterPipelineHandler) auth_add_account(clusterPipeline *v3.ClusterPipeline, remoteType string, userName string, redirectURL string, code string) error {

	if userName == "" {
		return errors.New("unauth")
	}

	remote, err := booter.New(*clusterPipeline, remoteType)
	if err != nil {
		return err
	}
	account, err := remote.Login(redirectURL, code)
	if err != nil {
		return err
	}
	account.Spec.UserName = userName

	if _, err := h.Management.Management.RemoteAccounts("user-" + userName).Create(account); err != nil {
		return err
	}
	return nil
}

func (h *ClusterPipelineHandler) test_auth(apiContext *types.APIContext) error {

	parts := strings.Split(apiContext.ID, ":")
	if len(parts) <= 1 {
		return errors.New("invalid ID")
	}
	ns := parts[0]
	id := parts[1]

	requestBody := make(map[string]interface{})
	requestBytes, err := ioutil.ReadAll(apiContext.Request.Body)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(requestBytes, &requestBody); err != nil {
		return err
	}
	var code, remoteType, clientID, clientSecret, redirectURL, scheme, host string

	if requestBody["code"] != nil {
		code = requestBody["code"].(string)
	}
	if requestBody["clientId"] != nil {
		clientID = requestBody["clientId"].(string)
	}
	if requestBody["clientSecret"] != nil {
		clientSecret = requestBody["clientSecret"].(string)
	}
	if requestBody["redirectUrl"] != nil {
		redirectURL = requestBody["redirectUrl"].(string)
	}
	if requestBody["scheme"] != nil {
		scheme = requestBody["scheme"].(string)
	}
	if requestBody["host"] != nil {
		host = requestBody["host"].(string)
	}
	if requestBody["type"] != nil {
		remoteType = requestBody["type"].(string)
	}

	clusterPipeline, err := h.Management.Management.ClusterPipelines(ns).Get(id, metav1.GetOptions{})
	if err != nil {
		return err
	}
	if clientID == "" || clientSecret == "" {

		if clusterPipeline.Spec.GithubConfig == nil {
			return errors.New("github not configured")
		}

		clientID = clusterPipeline.Spec.GithubConfig.ClientId
		clientSecret = clusterPipeline.Spec.GithubConfig.ClientSecret

		//oauth and add user
		userName := apiContext.Request.Header.Get("Impersonate-User")
		if err := h.test_auth_add_account(clusterPipeline, remoteType, userName, redirectURL, code); err != nil {
			return err
		}

	} else {
		if remoteType == "github" && clusterPipeline.Spec.GithubConfig == nil {
			clusterPipeline.Spec.GithubConfig = &v3.GithubConfig{
				Scheme:       scheme,
				Host:         host,
				ClientId:     clientID,
				ClientSecret: clientSecret,
			}
		}
		//oauth and add user
		userName := apiContext.Request.Header.Get("Impersonate-User")
		if err := h.test_auth_add_account(clusterPipeline, remoteType, userName, redirectURL, code); err != nil {
			return err
		}
		//update cluster pipeline config
		if _, err := h.Management.Management.ClusterPipelines(ns).Update(clusterPipeline); err != nil {
			return err
		}
	}
	apiContext.WriteResponse(200, clusterPipeline)
	return nil
}

func (h *ClusterPipelineHandler) test_auth_add_account(clusterPipeline *v3.ClusterPipeline, remoteType string, userName string, redirectURL string, code string) error {
	/*
		remote, err := booter.New(*clusterPipeline, scmType)
		if err != nil {
			return err
		}
		account, err := remote.Login(redirectURL, code)
		if err != nil {
			return err
		}
	*/
	account := &v3.RemoteAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name: "github-test",
		},
		Spec: v3.RemoteAccountSpec{

			Type: "github",
		},
	}
	if userName == "" {
		return errors.New("unauth")
	}
	account.Spec.UserName = userName
	if _, err := h.Management.Management.RemoteAccounts("user-" + userName).Create(account); err != nil {
		return err
	}
	return nil
}
