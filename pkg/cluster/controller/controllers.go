package controller

import (
	"context"

	"github.com/rancher/rancher/pkg/cluster/controller/authz"
	"github.com/rancher/rancher/pkg/cluster/controller/eventssyncer"
	"github.com/rancher/rancher/pkg/cluster/controller/healthsyncer"
	"github.com/rancher/rancher/pkg/cluster/controller/nodesyncer"
	"github.com/rancher/rancher/pkg/cluster/controller/secret"
	helmController "github.com/rancher/rancher/pkg/helm/controller"
	pipelineController "github.com/rancher/rancher/pkg/pipeline/controller"
	workloadController "github.com/rancher/rancher/pkg/workload/controller"
	"github.com/rancher/types/config"
)

func Register(ctx context.Context, cluster *config.ClusterContext) error {
	nodesyncer.Register(cluster)
	healthsyncer.Register(ctx, cluster)
	authz.Register(cluster)
	eventssyncer.Register(cluster)
	secret.Register(cluster)
	helmController.Register(cluster)
	pipelineController.Register(cluster)

	workloadContext := cluster.WorkloadContext()
	return workloadController.Register(ctx, workloadContext)
}
