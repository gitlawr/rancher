package engine

import (
	"github.com/rancher/rancher/pkg/pipeline/engine/jenkins"
	"github.com/rancher/types/apis/management.cattle.io/v3"
	"github.com/rancher/types/config"
)

type PipelineEngine interface {
	RunPipeline(pipeline *v3.Pipeline, triggerType string) error
	RerunHistory(history *v3.PipelineExecution) error
	StopHistory(history *v3.PipelineExecution) error
	GetStepLog(history *v3.PipelineExecution, stageOrdinal int, stepOrdinal int, paras map[string]interface{}) (string, error)
	OnHistoryCompelte(history *v3.PipelineExecution)
}

func New(cluster *config.ClusterContext, url string, user string, token string) (PipelineEngine, error) {
	client, err := jenkins.New(url, user, token)
	if err != nil {
		return nil, err
	}
	engine := &jenkins.JenkinsEngine{
		Client:  client,
		Cluster: cluster,
	}
	return engine, nil
}
