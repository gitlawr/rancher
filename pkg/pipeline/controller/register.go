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
	pipelineLifecycle := &PipelineLifecycle{
		cluster: cluster,
	}
	pipelineClient.AddLifecycle("pipeline-controller", pipelineLifecycle)

	pipelineHistoryClient := cluster.Management.Management.PipelineExecutions("")
	pipelineHistoryLifecycle := &PipelineHistoryLifecycle{
		cluster: cluster,
	}
	pipelineHistoryClient.AddLifecycle("pipeline-history-controller", pipelineHistoryLifecycle)

}
