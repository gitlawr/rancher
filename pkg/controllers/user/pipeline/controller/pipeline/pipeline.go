package pipeline

import (
	"context"
	"errors"
	"fmt"
	"github.com/rancher/rancher/pkg/controllers/user/pipeline/remote"
	"github.com/rancher/rancher/pkg/controllers/user/pipeline/utils"
	"github.com/rancher/rancher/pkg/ref"
	"github.com/rancher/types/apis/management.cattle.io/v3"
	"github.com/rancher/types/config"
	"github.com/satori/uuid"
	"github.com/sirupsen/logrus"
	"time"
)

//Lifecycle is responsible for watching pipelines and handling webhook management
//in source code repository. It also helps to maintain labels on pipelines.
type Lifecycle struct {
	pipelines                  v3.PipelineInterface
	pipelineLister             v3.PipelineLister
	sourceCodeCredentialLister v3.SourceCodeCredentialLister
}

func Register(ctx context.Context, cluster *config.UserContext) {
	clusterName := cluster.ClusterName
	clusterPipelineLister := cluster.Management.Management.ClusterPipelines("").Controller().Lister()
	pipelines := cluster.Management.Management.Pipelines("")
	pipelineLister := pipelines.Controller().Lister()
	pipelineExecutions := cluster.Management.Management.PipelineExecutions("")
	sourceCodeCredentialLister := cluster.Management.Management.SourceCodeCredentials("").Controller().Lister()

	pipelineLifecycle := &Lifecycle{
		pipelines:                  pipelines,
		pipelineLister:             pipelineLister,
		sourceCodeCredentialLister: sourceCodeCredentialLister,
	}
	s := &CronSyncer{
		clusterName:           clusterName,
		clusterPipelineLister: clusterPipelineLister,
		pipelineLister:        pipelineLister,
		pipelines:             pipelines,
		pipelineExecutions:    pipelineExecutions,
	}

	pipelines.AddClusterScopedLifecycle("pipeline-controller", cluster.ClusterName, pipelineLifecycle)
	go s.sync(ctx, syncInterval)
}

func (l *Lifecycle) Create(obj *v3.Pipeline) (*v3.Pipeline, error) {

	return l.sync(obj)
}

func (l *Lifecycle) Updated(obj *v3.Pipeline) (*v3.Pipeline, error) {

	return l.sync(obj)
}

func (l *Lifecycle) Remove(obj *v3.Pipeline) (*v3.Pipeline, error) {

	if obj.Status.WebHookID != "" {
		if err := l.deleteHook(obj); err != nil {
			//merely log error to avoid deletion block
			logrus.Warnf("fail to delete previous set webhook for pipeline '%s' - %v", obj.Spec.DisplayName, err)
		}
	}
	return obj, nil
}

func (l *Lifecycle) sync(obj *v3.Pipeline) (*v3.Pipeline, error) {

	if obj.Status.Token == "" {
		//random token for webhook validation
		obj.Status.Token = uuid.NewV4().String()
	}

	//handle sourceCodeCredential info
	if !utils.ValidSourceCodeConfig(obj.Spec) {
		return obj, fmt.Errorf("error invalid definition of pipeline '%s'", obj.Spec.DisplayName)
	}

	sourceCodeCredentialID := obj.Spec.Stages[0].Steps[0].SourceCodeConfig.SourceCodeCredentialName
	ns, name := ref.Parse(sourceCodeCredentialID)
	if obj.Status.SourceCodeCredential == nil ||
		obj.Status.SourceCodeCredential.Namespace != ns ||
		obj.Status.SourceCodeCredential.Name != name {
		updatedCred, err := l.sourceCodeCredentialLister.Get(ns, name)
		if err != nil {
			return obj, err
		}
		updatedCred.Spec.AccessToken = ""
		obj.Status.SourceCodeCredential = updatedCred
	}

	//TODO test
	//if obj.Status.PipelineState == "inactive" {
	//	v3.PipelineConditionActive.False(obj)
	//	v3.PipelineConditionActive.Message(obj, "Pipeline is deactivated")
	//} else {
	//	v3.PipelineConditionActive.True(obj)
	//	v3.PipelineConditionActive.Message(obj, "Pipeline is activated")
	//}

	//handle cron
	if obj.Spec.TriggerCronExpression == "" {
		obj.Labels = map[string]string{utils.PipelineCronLabel: "false"}
		if obj.Status.NextStart != "" {
			obj.Status.NextStart = ""
		}
	} else {
		obj.Labels = map[string]string{utils.PipelineCronLabel: "false"}
		nextStart, err := getNextStartTime(obj.Spec.TriggerCronExpression, obj.Spec.TriggerCronTimezone, time.Now())
		if err != nil {
			return obj, err
		}
		obj.Status.NextStart = nextStart
	}

	//handle webhook
	if obj.Status.WebHookID != "" && !hasWebhookTrigger(obj) {
		if err := l.deleteHook(obj); err != nil {
			logrus.Warnf("fail to delete previous set webhook for pipeline '%s' - %v", obj.Spec.DisplayName, err)
		}
		obj.Status.WebHookID = ""
	} else if hasWebhookTrigger(obj) && obj.Status.WebHookID == "" {
		id, err := l.createHook(obj)
		if err != nil {
			return obj, err
		}
		obj.Status.WebHookID = id
	}

	return obj, nil
}

func (l *Lifecycle) createHook(obj *v3.Pipeline) (string, error) {
	if utils.ValidSourceCodeConfig(obj.Spec) {
		return "", errors.New("invalid pipeline, missing sourcecode step")
	}
	credentialID := obj.Spec.Stages[0].Steps[0].SourceCodeConfig.SourceCodeCredentialName

	ns, name := ref.Parse(credentialID)
	credential, err := l.sourceCodeCredentialLister.Get(ns, name)
	if err != nil {
		return "", err
	}
	accessToken := credential.Spec.AccessToken
	kind := credential.Spec.SourceCodeType
	mockConfig := v3.ClusterPipeline{
		Spec: v3.ClusterPipelineSpec{
			GithubConfig: &v3.GithubClusterConfig{},
		},
	}
	remote, err := remote.New(mockConfig, kind)
	if err != nil {
		return "", err
	}

	id, err := remote.CreateHook(obj, accessToken)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (l *Lifecycle) deleteHook(obj *v3.Pipeline) error {
	if utils.ValidSourceCodeConfig(obj.Spec) {
		return errors.New("invalid pipeline, missing sourcecode step")
	}
	credentialID := obj.Spec.Stages[0].Steps[0].SourceCodeConfig.SourceCodeCredentialName

	ns, name := ref.Parse(credentialID)
	credential, err := l.sourceCodeCredentialLister.Get(ns, name)
	if err != nil {
		return err
	}
	accessToken := credential.Spec.AccessToken
	kind := credential.Spec.SourceCodeType
	mockConfig := v3.ClusterPipeline{
		Spec: v3.ClusterPipelineSpec{
			GithubConfig: &v3.GithubClusterConfig{},
		},
	}
	remote, err := remote.New(mockConfig, kind)
	if err != nil {
		return err
	}

	return remote.DeleteHook(obj, accessToken)
}

func hasWebhookTrigger(obj *v3.Pipeline) bool {
	if obj != nil && (obj.Spec.TriggerWebhookPr || obj.Spec.TriggerWebhookPush) {
		return true
	}
	return false
}
