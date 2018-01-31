package controller

import (
	"github.com/rancher/types/config"
)

func Register(cluster *config.ClusterContext) {
	clusterPipelineClient := cluster.Management.Management.ClusterPipelines("")
	clusterPipelineLifecycle := &ClusterPipelineLifecycle{
		cluster: cluster,
	}
	clusterPipelineClient.AddLifecycle("cluster-pipeline-controller", clusterPipelineLifecycle)
	clusterPipelineClient.AddClusterScopedHandler("cluster-pipeline-maintainer", cluster.ClusterName, clusterPipelineLifecycle.Sync)

	pipelineClient := cluster.Management.Management.Pipelines("")
	pipelineLifecycle := &PipelineLifecycle{}
	pipelineClient.AddLifecycle("pipeline-controller", pipelineLifecycle)

}
