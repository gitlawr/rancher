package controller

import (
	"github.com/rancher/types/apis/management.cattle.io/v3"
	"github.com/rancher/types/config"
)

type PipelineLifecycle struct {
	cluster *config.ClusterContext
}

func (l *PipelineLifecycle) Create(obj *v3.Pipeline) (*v3.Pipeline, error) {
	return obj, nil
}

func (l *PipelineLifecycle) Updated(obj *v3.Pipeline) (*v3.Pipeline, error) {
	return obj, nil
}

func (l *PipelineLifecycle) Remove(obj *v3.Pipeline) (*v3.Pipeline, error) {
	return obj, nil
}
