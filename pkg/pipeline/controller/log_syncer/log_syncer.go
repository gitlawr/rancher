package log_syncer

import (
	"context"
	"errors"
	"fmt"
	"github.com/rancher/rancher/pkg/cluster/utils"
	"github.com/rancher/rancher/pkg/pipeline/engine"
	utils2 "github.com/rancher/rancher/pkg/pipeline/utils"
	"github.com/rancher/types/apis/core/v1"
	"github.com/rancher/types/apis/management.cattle.io/v3"
	"github.com/rancher/types/config"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/labels"
	"time"
)

const (
	syncInterval = 20 * time.Second
)

type ExecutionLogSyncer struct {
	pipelineExecutionLister    v3.PipelineExecutionLister
	pipelineExecutionLogLister v3.PipelineExecutionLogLister
	pipelineExecutionLogs      v3.PipelineExecutionLogInterface
	nodeLister                 v1.NodeLister
	serviceLister              v1.ServiceLister
	cluster                    *config.ClusterContext
}

func Register(ctx context.Context, cluster *config.ClusterContext) {
	pipelineExecutionLister := cluster.Management.Management.PipelineExecutions("").Controller().Lister()
	pipelineExecutionLogs := cluster.Management.Management.PipelineExecutionLogs("")
	pipelineExecutionLogLister := pipelineExecutionLogs.Controller().Lister()

	nodeLister := cluster.Core.Nodes("").Controller().Lister()
	serviceLister := cluster.Core.Services("").Controller().Lister()
	s := &ExecutionLogSyncer{
		pipelineExecutionLister:    pipelineExecutionLister,
		pipelineExecutionLogLister: pipelineExecutionLogLister,
		pipelineExecutionLogs:      pipelineExecutionLogs,
		nodeLister:                 nodeLister,
		serviceLister:              serviceLister,
	}
	go s.syncState(ctx, syncInterval)
}

func (s *ExecutionLogSyncer) syncState(ctx context.Context, syncInterval time.Duration) {
	for range utils.TickerContext(ctx, syncInterval) {
		logrus.Debugf("Start heartbeat")
		s.syncLogs()
		logrus.Debugf("Heartbeat complete")
	}

}

func (s *ExecutionLogSyncer) syncLogs() {
	Logs, err := s.pipelineExecutionLogLister.List("", utils2.PIPELINE_INPROGRESS_LABEL.AsSelector())
	if err != nil {
		logrus.Errorf("Error listing PipelineExecutionLogs - %v", err)
	}
	url, err := s.getJenkinsURL()
	if err != nil {
		logrus.Errorf("Error get Jenkins url - %v", err)
	}
	pipelineEngine, err := engine.New(s.cluster, url)
	if err != nil {
		logrus.Errorf("Error get Jenkins engine - %v", err)
	}
	for _, e := range Logs {
		execution, err := s.pipelineExecutionLister.Get(e.Namespace, e.Spec.PipelineExecutionName)
		if err != nil {
			logrus.Errorf("Error get pipeline execution - %v", err)
		}
		logText, err := pipelineEngine.GetStepLog(execution, e.Spec.Stage, e.Spec.Step, nil)
		if err != nil {
			logrus.Errorf("Error get pipeline execution log - %v", err)
			e.Spec.Message += fmt.Sprintf("\nError get pipeline execution log - %v", err)
			e.Labels = utils2.PIPELINE_FINISH_LABEL
			if _, err := s.pipelineExecutionLogs.Update(e); err != nil {
				logrus.Errorf("Error update pipeline execution log - %v", err)
			}
			continue
		}
		//TODO trim message
		e.Spec.Message = logText
		stepState := execution.Status.Stages[e.Spec.Stage].Steps[e.Spec.Step].State
		if stepState != v3.StateWaiting && stepState != v3.StateBuilding {
			e.Labels = utils2.PIPELINE_FINISH_LABEL
		}
		if _, err := s.pipelineExecutionLogs.Update(e); err != nil {
			logrus.Errorf("Error update pipeline execution log - %v", err)
		}
	}
}

//FIXME proper way to connect to Jenkins in cluster
func (l *ExecutionLogSyncer) getJenkinsURL() (string, error) {

	nodes, err := l.nodeLister.List("", labels.NewSelector())
	if err != nil {
		return "", err
	}
	if len(nodes) < 1 {
		return "", errors.New("no available nodes")
	}
	if len(nodes[0].Status.Addresses) < 1 {
		return "", errors.New("no available address")
	}
	host := nodes[0].Status.Addresses[0].Address

	svcport := 0
	service, err := l.serviceLister.Get("cattle-pipeline", "jenkins")
	if err != nil {
		return "", err
	}

	ports := service.Spec.Ports
	for _, port := range ports {
		if port.NodePort != 0 && port.Name == "http" {
			svcport = int(port.NodePort)
			break
		}
	}
	return fmt.Sprintf("http://%s:%d", host, svcport), nil
}
