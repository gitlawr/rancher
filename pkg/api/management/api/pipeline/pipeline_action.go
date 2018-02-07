package pipeline

import (
	"strings"

	"fmt"
	"github.com/rancher/norman/api/access"
	"github.com/rancher/norman/types"
	"github.com/rancher/types/apis/management.cattle.io/v3"
	"github.com/rancher/types/client/management/v3"
	"github.com/rancher/types/config"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"time"
)

type Handler struct {
	Management config.ManagementContext
}

func Formatter(apiContext *types.APIContext, resource *types.RawResource) {
	resource.AddAction(apiContext, "activate")
	resource.AddAction(apiContext, "deactivate")
	resource.AddAction(apiContext, "run")
	resource.AddAction(apiContext, "export")
}

func (h *Handler) ActionHandler(actionName string, action *types.Action, apiContext *types.APIContext) error {
	logrus.Infof("do activity action:%s", actionName)

	switch actionName {
	case "activate":
		return h.activate(apiContext)

	case "deactivate":
		return h.deactivate(apiContext)

	case "run":
		return h.run(apiContext)
	case "export":
		return h.export(apiContext)
	}
	return nil
}

func (h *Handler) activate(apiContext *types.APIContext) error {
	//parts := strings.Split(apiContext.ID, ":")
	//ns := parts[0]
	//id := parts[1]
	//
	//pipelineClient := h.Management.Management.Pipelines(ns)
	//pipeline, err := pipelineClient.Get(id, metav1.GetOptions{})
	//if err != nil {
	//	logrus.Errorf("Error while getting pipeline for %s :%v", apiContext.ID, err)
	//	return err
	//}
	//
	//if pipeline.Spec.Active == false {
	//	pipeline.Spec.Active = true
	//} else {
	//	return errors.New("the pipeline is already activated")
	//}
	//_, err = pipelineClient.Update(pipeline)
	//if err != nil {
	//	logrus.Errorf("Error while updating pipeline:%v", err)
	//	return err
	//}
	//
	//data := map[string]interface{}{}
	//if err := access.ByID(apiContext, apiContext.Version, apiContext.Type, apiContext.ID, &data); err != nil {
	//	return err
	//}
	//
	//apiContext.WriteResponse(http.StatusOK, data)
	return nil
}

func (h *Handler) deactivate(apiContext *types.APIContext) error {
	//parts := strings.Split(apiContext.ID, ":")
	//ns := parts[0]
	//id := parts[1]
	//
	//pipelineClient := h.Management.Management.Pipelines(ns)
	//pipeline, err := pipelineClient.Get(id, metav1.GetOptions{})
	//if err != nil {
	//	logrus.Errorf("Error while getting pipeline for %s :%v", apiContext.ID, err)
	//	return err
	//}
	//
	//if pipeline.Spec.Active == true {
	//	pipeline.Spec.Active = false
	//} else {
	//	return errors.New("the pipeline is already deactivated")
	//}
	//_, err = pipelineClient.Update(pipeline)
	//if err != nil {
	//	logrus.Errorf("Error while updating pipeline:%v", err)
	//	return err
	//}
	//
	//data := map[string]interface{}{}
	//if err := access.ByID(apiContext, apiContext.Version, apiContext.Type, apiContext.ID, &data); err != nil {
	//	return err
	//}
	//
	//apiContext.WriteResponse(http.StatusOK, data)
	return nil
}

func (h *Handler) run(apiContext *types.APIContext) error {
	parts := strings.Split(apiContext.ID, ":")
	ns := parts[0]
	id := parts[1]
	pipelineClient := h.Management.Management.Pipelines(ns)
	pipeline, err := pipelineClient.Get(id, metav1.GetOptions{})
	if err != nil {
		logrus.Errorf("Error while getting pipeline for %s :%v", apiContext.ID, err)
		return err
	}

	historyClient := h.Management.Management.PipelineExecutions(ns)
	history := initHistory(pipeline, v3.TriggerTypeManual)
	history, err = historyClient.Create(history)
	if err != nil {
		return err
	}
	pipeline.Status.NextRun++
	pipeline.Status.LastExecutionId = history.Name
	pipeline.Status.LastStarted = time.Now().String()

	_, err = pipelineClient.Update(pipeline)

	data := map[string]interface{}{}
	if err := access.ByID(apiContext, apiContext.Version, client.PipelineExecutionType, ns+":"+history.Name, &data); err != nil {
		return err
	}

	apiContext.WriteResponse(http.StatusOK, data)
	return err
}

func (h *Handler) export(apiContext *types.APIContext) error {
	return nil
}

func getNextHistoryName(p *v3.Pipeline) string {
	if p == nil {
		return ""
	}
	return fmt.Sprintf("%s-%d", p.Name, p.Status.NextRun)
}

func initHistory(p *v3.Pipeline, triggerType string) *v3.PipelineExecution {
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
