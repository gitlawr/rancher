package controller

import (
	"context"
	"errors"
	"fmt"
	"github.com/rancher/rancher/pkg/pipeline/engine"
	"github.com/rancher/types/apis/management.cattle.io/v3"
	"github.com/rancher/types/config"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"time"
)

const (
	syncInterval = 20 * time.Second
)

type ExecutionStateSyncer struct {
	pipelineExecutionLister v3.PipelineExecutionLister
	pipelineExecutions      v3.PipelineExecutionInterface
	cluster                 *config.ClusterContext
}

func Registerxxx(ctx context.Context, cluster *config.ClusterContext) {
	pipelineExecutions := cluster.Management.Management.PipelineExecutions("")
	pipelineExecutionLister := pipelineExecutions.Controller().Lister()
	s := &ExecutionStateSyncer{
		pipelineExecutionLister: pipelineExecutionLister,
		pipelineExecutions:      pipelineExecutions,
	}
	//pipelineExecutions.Controller().AddHandler(s.GetName(), s.sync)
	fmt.Println(s)
	//go s.syncState(ctx, syncInterval)
}

func (s *ExecutionStateSyncer) Sync(key string, obj *v3.Pipeline) error {

	return nil
}

func (s *ExecutionStateSyncer) syncState(ctx context.Context, syncInterval time.Duration) {
	executions, err := s.pipelineExecutionLister.List("", labels.NewSelector())
	if err != nil {
		logrus.Errorf("Error listing PipelineExecutions - %v", err)
	}
	url, err := s.getJenkinsURL()
	if err != nil {
		logrus.Errorf("Error get Jenkins url - %v", err)
	}
	pipelineEngine, err := engine.New(s.cluster, url)
	if err != nil {
		logrus.Errorf("Error get Jenkins engine - %v", err)
	}
	for _, e := range executions {
		if e.Status.State == v3.StateWaiting || e.Status.State == v3.StateBuilding {
			updated, err := pipelineEngine.SyncExecution(e)
			if err != nil {
				logrus.Errorf("Error sync pipeline execution - %v", err)
				e.Status.State = v3.StateFail
				if _, err := s.pipelineExecutions.Update(e); err != nil {
					logrus.Errorf("Error update pipeline execution - %v", err)
				}
			} else if updated {
				if _, err := s.pipelineExecutions.Update(e); err != nil {
					logrus.Errorf("Error update pipeline execution - %v", err)
				}
			}
		}
	}
}

//FIXME proper way to connect to Jenkins in cluster
func (l *ExecutionStateSyncer) getJenkinsURL() (string, error) {

	nodes, err := l.cluster.Core.Nodes("").List(metav1.ListOptions{})
	if err != nil {
		return "", err
	}
	if len(nodes.Items) < 1 {
		return "", errors.New("no available nodes")
	}
	if len(nodes.Items[0].Status.Addresses) < 1 {
		return "", errors.New("no available address")
	}
	host := nodes.Items[0].Status.Addresses[0].Address

	svcport := 0
	service, err := l.cluster.Core.Services("cattle-pipeline").Get("jenkins", metav1.GetOptions{})
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

func (s *ExecutionStateSyncer) GetName() string {
	return "pipelineexecution-statesync-controller"
}
