package bitbucket

import (
	"encoding/json"
	"github.com/rancher/rancher/pkg/pipeline/remote/model"
	"github.com/rancher/rancher/pkg/settings"
	"github.com/rancher/types/apis/project.cattle.io/v3"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

var c = client{
	ClientID:     "w9RY6xSrpLBYvm5Gga",
	ClientSecret: "CqVB3tcjV8jg99B6pktDKSRSkkr5uEsP",
	RedirectURL:  "http://127.0.0.1:3000/Callback",
}

var accessToken = "qfH0ZPmAWuZr_grAvKv_epXbTOd1NOqxy6i1-CrIpcjidHG1LU5I5KbbYT2abJTvpQ0McOd6DbbUj6SW9JQ="

func Test_User(t *testing.T) {
	code := "LwEVQgNEpQH4ZHfa6t"
	cred, err := c.Login(code)
	if err != nil {
		logrus.Error(err)
	}
	b, _ := json.Marshal(cred)
	logrus.Infof(string(b))
}

func Test_Repos(t *testing.T) {
	cred := &v3.SourceCodeCredential{
		Spec: v3.SourceCodeCredentialSpec{ProjectName: "p-abc", SourceCodeType: model.BitbucketCloudType, UserName: "lawrlee", DisplayName: "lawrence li", AvatarURL: "", HTMLURL: "", LoginName: "lawrlee", GitLoginName: cloneUserName, AccessToken: accessToken},
	}
	repos, err := c.Repos(cred)
	if err != nil {
		logrus.Error(err)
	}
	b, _ := json.Marshal(repos)
	logrus.Infof(string(b))
}

func Test_Hook(t *testing.T) {
	p := &v3.Pipeline{
		ObjectMeta: v1.ObjectMeta{
			Name: "pipeline-abc",
		},
		Spec: v3.PipelineSpec{
			ProjectName:   "p-abc",
			RepositoryURL: "https://lawrlee@bitbucket.org/lawrlee/test.git",
		},
	}
	settings.ServerURL.Set("http://192.168.99.1")
	id, err := c.CreateHook(p, accessToken)
	if err != nil {
		logrus.Error(err)
	}
	logrus.Infof("get id:%s", id)

}

func Test_DelHook(t *testing.T) {
	p := &v3.Pipeline{
		ObjectMeta: v1.ObjectMeta{
			Name: "pipeline-abc",
		},
		Spec: v3.PipelineSpec{
			ProjectName:   "p-abc",
			RepositoryURL: "https://lawrlee@bitbucket.org/lawrlee/test.git",
		},
	}
	settings.ServerURL.Set("http://54.78.6.13")
	err := c.DeleteHook(p, accessToken)
	if err != nil {
		logrus.Error(err)
	}
}

func Test_GetFile(t *testing.T) {
	b, err := c.GetPipelineFileInRepo("https://lawrlee@bitbucket.org/lawrlee/test.git", "master", accessToken)
	if err != nil {
		logrus.Error(err)
	}
	logrus.Info(string(b))
}

func Test_SetFile(t *testing.T) {
	content := "my pipelinefile 2"
	err := c.SetPipelineFileInRepo("https://lawrlee@bitbucket.org/lawrlee/test.git", "master", accessToken, []byte(content))
	if err != nil {
		logrus.Error(err)
	}
}

func Test_Branch(t *testing.T) {
	bs, err := c.GetBranches("https://lawrlee@bitbucket.org/lawrlee/test.git", accessToken)
	if err != nil {
		logrus.Error(err)
	}
	logrus.Error(bs)
}

func Test_Head(t *testing.T) {
	info, err := c.GetHeadInfo("https://lawrlee@bitbucket.org/lawrlee/test.git", "master", accessToken)
	if err != nil {
		logrus.Error(err)
	}
	b, _ := json.Marshal(info)
	logrus.Error(string(b))
}
