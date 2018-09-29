package drivers

import (
	"encoding/json"
	"fmt"
	"github.com/rancher/rancher/pkg/pipeline/providers"
	"github.com/rancher/rancher/pkg/pipeline/remote/bitbucket"
	"github.com/rancher/rancher/pkg/pipeline/remote/model"
	"github.com/rancher/rancher/pkg/pipeline/utils"
	"github.com/rancher/rancher/pkg/ref"
	"github.com/rancher/types/apis/project.cattle.io/v3"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

const (
	BitbucketCloudWebhookHeader  = "X-Hook-UUID"
	bitbucketCloudEventHeader    = "X-Event-Key"
	bitbucketCloudPushEvent      = "repo:push"
	bitbucketCloudPrCreatedEvent = "pullrequest:created"
	bitbucketCloudPrUpdatedEvent = "pullrequest:updated"
	bitbucketCloudStateOpen      = "OPEN"
)

type BitbucketCloudDriver struct {
	PipelineLister             v3.PipelineLister
	PipelineExecutions         v3.PipelineExecutionInterface
	SourceCodeCredentials      v3.SourceCodeCredentialInterface
	SourceCodeCredentialLister v3.SourceCodeCredentialLister
}

func (b BitbucketCloudDriver) Execute(req *http.Request) (int, error) {
	event := req.Header.Get(bitbucketCloudEventHeader)
	if event != bitbucketCloudPushEvent && event != bitbucketCloudPrCreatedEvent && event != bitbucketCloudPrUpdatedEvent {
		return http.StatusUnprocessableEntity, fmt.Errorf("not trigger for event:%s", event)
	}

	pipelineID := req.URL.Query().Get("pipelineId")
	ns, name := ref.Parse(pipelineID)
	pipeline, err := b.PipelineLister.Get(ns, name)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return http.StatusUnprocessableEntity, err
	}
	logrus.Infof("get body:%s", string(body))

	info := &model.BuildInfo{}
	if event == bitbucketCloudPushEvent {
		info, err = parseBitbucketPushPayload(body)
		if err != nil {
			return http.StatusUnprocessableEntity, err
		}
	} else if event == bitbucketCloudPrCreatedEvent || event == bitbucketCloudPrUpdatedEvent {
		info, err = parseBitbucketPullRequestPayload(body)
		if err != nil {
			return http.StatusUnprocessableEntity, err
		}
	}

	if (info.Event == utils.WebhookEventPush && !pipeline.Spec.TriggerWebhookPush) ||
		(info.Event == utils.WebhookEventTag && !pipeline.Spec.TriggerWebhookTag) ||
		(info.Event == utils.WebhookEventPullRequest && !pipeline.Spec.TriggerWebhookPr) {
		return http.StatusUnavailableForLegalReasons, fmt.Errorf("trigger for event '%s' is disabled", info.Event)
	}

	pipelineConfig, err := providers.GetPipelineConfigByBranch(b.SourceCodeCredentials, b.SourceCodeCredentialLister, pipeline, info.Branch)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if pipelineConfig == nil {
		//no pipeline config to run
		return http.StatusOK, nil
	}

	if !utils.Match(pipelineConfig.Branch, info.Branch) {
		return http.StatusUnavailableForLegalReasons, fmt.Errorf("skipped branch '%s'", info.Branch)
	}

	if _, err := utils.GenerateExecution(b.PipelineExecutions, pipeline, pipelineConfig, info); err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func parseBitbucketPushPayload(raw []byte) (*model.BuildInfo, error) {
	info := &model.BuildInfo{}
	payload := bitbucket.PushEventPayload{}
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil, err
	}
	info.TriggerType = utils.TriggerTypeWebhook

	if len(payload.Push.Changes) > 0 {
		change := payload.Push.Changes[0]
		info.Commit = change.New.Target.Hash
		info.Branch = change.New.Name
		info.Message = change.New.Target.Message
		info.Author = change.New.Target.Author.User.UserName
		info.AvatarURL = change.New.Target.Author.User.Links.Avatar.Href
		info.HTMLLink = change.New.Target.Links.Html.Href

		switch change.New.Type {
		case "tag", "annotated_tag", "bookmark":
			info.Event = utils.WebhookEventTag
			info.Ref = RefsTagPrefix + change.New.Name
		default:
			info.Event = utils.WebhookEventPush
			info.Ref = RefsBranchPrefix + change.New.Name
		}
	}
	return info, nil
}

func parseBitbucketPullRequestPayload(raw []byte) (*model.BuildInfo, error) {
	info := &model.BuildInfo{}
	payload := bitbucket.PullRequestEventPayload{}
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil, err
	}

	if payload.PullRequest.State != bitbucketCloudStateOpen {
		return nil, fmt.Errorf("no trigger for closed pull requests")
	}

	info.TriggerType = utils.TriggerTypeWebhook
	info.Event = utils.WebhookEventPullRequest
	info.RepositoryURL = fmt.Sprintf("https://bitbucket.org/%s", payload.PullRequest.Source.Repository.FullName)
	info.Branch = payload.PullRequest.Destination.Branch.Name
	info.Ref = RefsBranchPrefix + payload.PullRequest.Source.Branch.Name
	info.HTMLLink = payload.PullRequest.Links.Html.Href
	info.Title = payload.PullRequest.Title
	info.Message = payload.PullRequest.Title
	info.Commit = payload.PullRequest.Source.Commit.Hash
	info.Author = payload.PullRequest.Author.UserName
	info.AvatarURL = payload.PullRequest.Author.Links.Avatar.Href
	return info, nil
}
