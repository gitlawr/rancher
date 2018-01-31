package pipeline

import (
	"errors"
	"strings"

	"github.com/rancher/norman/types"
	"github.com/rancher/types/apis/management.cattle.io/v3"
	"github.com/rancher/types/config"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strconv"
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
	history := &v3.PipelineHistory{
		Spec: v3.PipelineHistorySpec{
			RunNumber:   pipeline.Status.NextRunNumber,
			TriggerType: v3.TriggerTypeManual,
			Pipeline:    *pipeline,
			DisplayName: pipeline.Name + "-" + strconv.Itoa(pipeline.Status.NextRunNumber),
		},
	}
	if _, err := historyClient.Create(history); err != nil {
		return err
	}
	pipeline.Status.NextRunNumber++
	pipeline.Status.LastRunId = history.Name
	pipeline.Status.LastRunTime = time.Now().UnixNano() / int64(time.Millisecond)

	_, err = pipelineClient.Update(pipeline)

	return err
}

func (h *Handler) export(apiContext *types.APIContext) error {
	return nil
}
