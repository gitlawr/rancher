package pipeline

import (
	"errors"
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
	parts := strings.Split(apiContext.ID, ":")
	ns := parts[0]
	id := parts[1]

	pipelineClient := h.Management.Management.Pipelines(ns)
	pipeline, err := pipelineClient.Get(id, metav1.GetOptions{})
	if err != nil {
		logrus.Errorf("Error while getting pipeline for %s :%v", apiContext.ID, err)
		return err
	}

	if pipeline.Spec.Active == false {
		pipeline.Spec.Active = true
	} else {
		return errors.New("the pipeline is already activated")
	}
	_, err = pipelineClient.Update(pipeline)
	if err != nil {
		logrus.Errorf("Error while updating pipeline:%v", err)
		return err
	}

	data := map[string]interface{}{}
	if err := access.ByID(apiContext, apiContext.Version, apiContext.Type, apiContext.ID, &data); err != nil {
		return err
	}

	apiContext.WriteResponse(http.StatusOK, data)
	return nil
}

func (h *Handler) deactivate(apiContext *types.APIContext) error {
	parts := strings.Split(apiContext.ID, ":")
	ns := parts[0]
	id := parts[1]

	pipelineClient := h.Management.Management.Pipelines(ns)
	pipeline, err := pipelineClient.Get(id, metav1.GetOptions{})
	if err != nil {
		logrus.Errorf("Error while getting pipeline for %s :%v", apiContext.ID, err)
		return err
	}

	if pipeline.Spec.Active == true {
		pipeline.Spec.Active = false
	} else {
		return errors.New("the pipeline is already deactivated")
	}
	_, err = pipelineClient.Update(pipeline)
	if err != nil {
		logrus.Errorf("Error while updating pipeline:%v", err)
		return err
	}

	data := map[string]interface{}{}
	if err := access.ByID(apiContext, apiContext.Version, apiContext.Type, apiContext.ID, &data); err != nil {
		return err
	}

	apiContext.WriteResponse(http.StatusOK, data)
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

	historyClient := h.Management.Management.PipelineHistories(ns)
	history := initHistory(pipeline, v3.TriggerTypeManual)
	history, err = historyClient.Create(history)
	if err != nil {
		return err
	}
	pipeline.Status.NextRunNumber++
	pipeline.Status.LastRunId = history.Name
	pipeline.Status.LastRunTime = time.Now().UnixNano() / int64(time.Millisecond)

	_, err = pipelineClient.Update(pipeline)

	data := map[string]interface{}{}
	if err := access.ByID(apiContext, apiContext.Version, client.PipelineHistoryType, ns+":"+history.Name, &data); err != nil {
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
	return fmt.Sprintf("%s-%d", p.Name, p.Status.NextRunNumber)
}

func initHistory(p *v3.Pipeline, triggerType string) *v3.PipelineHistory {
	history := &v3.PipelineHistory{
		ObjectMeta: metav1.ObjectMeta{
			Name: getNextHistoryName(p),
		},
		Spec: v3.PipelineHistorySpec{
			RunNumber:   p.Status.NextRunNumber,
			TriggerType: triggerType,
			Pipeline:    *p,
			//DisplayName: getNextHistoryName(pipeline),
		},
	}
	history.Status.State = v3.StateWaiting
	history.Status.StageStatus = make([]v3.StageStatus, len(p.Spec.Stages))

	for i, stage := range history.Status.StageStatus {
		stage.State = v3.StateWaiting
		stepsize := len(p.Spec.Stages[i].Steps)
		stage.StepStatus = make([]v3.StepStatus, stepsize)
		for _, step := range stage.StepStatus {
			step.State = v3.StateWaiting
		}
	}
	return history
}
