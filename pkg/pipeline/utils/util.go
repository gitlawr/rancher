package utils

import (
	"errors"
	"fmt"
	"github.com/rancher/norman/types"
	corev1 "github.com/rancher/types/apis/core/v1"
	"github.com/rancher/types/apis/management.cattle.io/v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"net/url"
	"time"
)

var CI_ENDPOINT = ""

var PIPELINE_FINISH_LABEL = labels.Set(map[string]string{"pipeline.management.cattle.io/finish": "true"})
var PIPELINE_INPROGRESS_LABEL = labels.Set(map[string]string{"pipeline.management.cattle.io/finish": "false"})

func InitHistory(p *v3.Pipeline, triggerType string) *v3.PipelineExecution {
	history := &v3.PipelineExecution{
		ObjectMeta: metav1.ObjectMeta{
			Name:      getNextHistoryName(p),
			Namespace: p.Namespace,
			Labels:    PIPELINE_INPROGRESS_LABEL,
		},
		Spec: v3.PipelineExecutionSpec{
			ProjectName:  p.Spec.ProjectName,
			PipelineName: p.Name,
			Run:          p.Status.NextRun,
			TriggeredBy:  triggerType,
			Pipeline:     *p,
			//DisplayName: getNextHistoryName(pipeline),
		},
	}
	history.Status.State = v3.StateWaiting
	history.Status.Stages = make([]v3.StageStatus, len(p.Spec.Stages))

	for i := 0; i < len(history.Status.Stages); i++ {
		stage := &history.Status.Stages[i]
		stage.State = v3.StateWaiting
		stepsize := len(p.Spec.Stages[i].Steps)
		stage.Steps = make([]v3.StepStatus, stepsize)
		for j := 0; j < stepsize; j++ {
			step := &stage.Steps[j]
			step.State = v3.StateWaiting
		}
	}
	return history
}

func getNextHistoryName(p *v3.Pipeline) string {
	if p == nil {
		return ""
	}
	return fmt.Sprintf("%s-%d", p.Name, p.Status.NextRun)
}

func IsStageSuccess(stage v3.StageStatus) bool {
	if stage.State == v3.StateSuccess {
		return true
	} else if stage.State == v3.StateFail || stage.State == v3.StateDenied {
		return false
	}
	successSteps := 0
	for _, step := range stage.Steps {
		if step.State == v3.StateSuccess || step.State == v3.StateSkip {
			successSteps++
		}
	}
	return successSteps == len(stage.Steps)
}

func UpdateEndpoint(apiContext *types.APIContext) error {

	reqUrl := apiContext.URLBuilder.Current()
	u, err := url.Parse(reqUrl)
	if err != nil {
		return err
	}
	CI_ENDPOINT = fmt.Sprintf("%s://%s/hooks", u.Scheme, u.Host)
	return nil
}

//FIXME proper way to connect to Jenkins in cluster
func GetJenkinsURL(nodeLister corev1.NodeLister, serviceLister corev1.ServiceLister) (string, error) {
	nodes, err := nodeLister.List("", labels.NewSelector())
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
	service, err := serviceLister.Get("cattle-pipeline", "jenkins")
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

func IsExecutionFinish(execution *v3.PipelineExecution) bool {
	if execution == nil {
		return false
	}
	if execution.Status.State != v3.StateWaiting && execution.Status.State != v3.StateBuilding {
		return true
	}
	return false
}

func RunPipeline(pipelines v3.PipelineInterface, pipelineExecutions v3.PipelineExecutionInterface, pipeline *v3.Pipeline, triggerType string) error {

	//Generate a new pipeline execution
	execution := InitHistory(pipeline, triggerType)
	execution, err := pipelineExecutions.Create(execution)
	if err != nil {
		return err
	}
	pipeline.Status.NextRun++
	pipeline.Status.LastExecutionId = execution.Name
	pipeline.Status.LastStarted = time.Now().Format(time.RFC3339)

	_, err = pipelines.Update(pipeline)
	if err != nil {
		return err
	}
	return nil
}
