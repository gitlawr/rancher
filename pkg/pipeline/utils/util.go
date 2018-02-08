package utils

import (
	"fmt"
	"github.com/rancher/types/apis/management.cattle.io/v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func InitHistory(p *v3.Pipeline, triggerType string) *v3.PipelineExecution {
	history := &v3.PipelineExecution{
		ObjectMeta: metav1.ObjectMeta{
			Name: getNextHistoryName(p),
		},
		Spec: v3.PipelineExecutionSpec{
			Run:         p.Status.NextRun,
			TriggeredBy: triggerType,
			Pipeline:    *p,
			//DisplayName: getNextHistoryName(pipeline),
		},
	}
	history.Status.State = v3.StateWaiting
	history.Status.Stages = make([]v3.StageStatus, len(p.Spec.Stages))

	for i, stage := range history.Status.Stages {
		stage.State = v3.StateWaiting
		stepsize := len(p.Spec.Stages[i].Steps)
		stage.Steps = make([]v3.StepStatus, stepsize)
		for _, step := range stage.Steps {
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
