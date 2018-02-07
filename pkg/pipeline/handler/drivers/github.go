package drivers

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/pkg/errors"
	"github.com/rancher/norman/types"
	"github.com/rancher/types/apis/management.cattle.io/v3"
	"github.com/rancher/types/config"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

const GITHUB_WEBHOOK_HEADER = "X-GitHub-Event"

type GithubDriver struct {
	Management *config.ManagementContext
}

func (g GithubDriver) Execute(apiContext types.APIContext) (int, error) {
	var signature string
	if signature = apiContext.Request.Header.Get("X-Hub-Signature"); len(signature) == 0 {
		logrus.Errorf("receive github webhook,no signature")
		return 422, errors.New("github webhook missing signature")
	}
	event := apiContext.Request.Header.Get(GITHUB_WEBHOOK_HEADER)
	if event == "ping" {
		return 200, nil
	} else if event != "push" && event != "pullrequest" && event != "tag" {
		//TODO check event types
		return 422, fmt.Errorf("not trigger for event:%s", event)
	}

	pipelineId := apiContext.Request.FormValue("pipelineId")
	parts := strings.Split(pipelineId, ":")
	if len(parts) < 0 {
		return 433, errors.New("pipeline id not valid")
	}
	ns := parts[0]
	id := parts[1]
	pipeline, err := g.Management.Management.Pipelines("").GetNamespaced(ns, id, metav1.GetOptions{})
	if err != nil {
		return 500, err
	}

	//////
	body, err := ioutil.ReadAll(apiContext.Request.Body)
	if err != nil {
		logrus.Errorf("receive github webhook, got error:%v", err)
		return 422, err
	}
	if match := VerifyGithubWebhookSignature([]byte(pipeline.Status.Token), signature, body); !match {
		logrus.Errorf("receive github webhook, invalid signature")
		return 422, errors.New("github webhook invalid signature")
	}
	//check branch
	payload := &github.WebHookPayload{}
	if err := json.Unmarshal(body, payload); err != nil {
		logrus.Error("fail to parse github webhook payload")
		return 422, err
	}
	//TODO
	if *payload.Ref != "refs/heads/"+pipeline.Spec.Stages[0].Steps[0].SourceCodeStepConfig.Branch {
		logrus.Warningf("branch not match:%v,%v", *payload.Ref, pipeline.Spec.Stages[0].Steps[0].SourceCodeStepConfig.Branch)
		return 500, errors.New("branch not match")
	}

	return 0, nil
}

func VerifyGithubWebhookSignature(secret []byte, signature string, body []byte) bool {

	const signaturePrefix = "sha1="
	const signatureLength = 45 // len(SignaturePrefix) + len(hex(sha1))

	if len(signature) != signatureLength || !strings.HasPrefix(signature, signaturePrefix) {
		return false
	}

	actual := make([]byte, 20)
	hex.Decode(actual, []byte(signature[5:]))
	computed := hmac.New(sha1.New, secret)
	computed.Write(body)

	return hmac.Equal([]byte(computed.Sum(nil)), actual)
}

func VerifyGitRefs(pipeline v3.Pipeline, refs string) {
	payload := &github.WebHookPayload{}
	payload.GetRef()
	if refs == "refs/heads/"+pipeline.Spec.Stages[0].Steps[0].SourceCodeStepConfig.Branch {

	}
}
