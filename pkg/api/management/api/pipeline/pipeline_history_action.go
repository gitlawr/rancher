package pipeline

import (
	"github.com/pkg/errors"
	"github.com/rancher/norman/types"
	"github.com/rancher/types/apis/management.cattle.io/v3"
	"github.com/rancher/types/config"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strconv"
	"strings"
)

type HistoryHandler struct {
	Management config.ManagementContext
}

func HistoryFormatter(apiContext *types.APIContext, resource *types.RawResource) {
	resource.AddAction(apiContext, "rerun")
	resource.AddAction(apiContext, "stop")
}

func (h *HistoryHandler) ActionHandler(actionName string, action *types.Action, apiContext *types.APIContext) error {
	logrus.Infof("do activity action:%s", actionName)

	switch actionName {
	case "rerun":
		return h.rerun(apiContext)

	case "deactivate":
		return h.stop(apiContext)
	case "notify":
		return h.notify(apiContext)
	}
	return nil
}

func (h *HistoryHandler) rerun(apiContext *types.APIContext) error {
	return nil
}

func (h *HistoryHandler) stop(apiContext *types.APIContext) error {
	return nil
}

func (h *HistoryHandler) notify(apiContext *types.APIContext) error {
	stepName := apiContext.Request.FormValue("stepName")
	state := apiContext.Request.FormValue("state")
	//TODO token check
	//token := apiContext.Request.FormValue("token")
	parts := strings.Split(stepName, "-")
	if len(parts) < 3 {
		return errors.New("invalid stepName")
	}
	stageOrdinal, err := strconv.Atoi(parts[1])
	if err != nil {
		return err
	}
	stepOrdinal, err := strconv.Atoi(parts[2])
	if err != nil {
		return err
	}
	parts = strings.Split(apiContext.ID, ":")
	ns := parts[0]
	id := parts[1]
	pipelineHistoryClient := h.Management.Management.PipelineExecutions(ns)
	pipelineHistory, err := pipelineHistoryClient.Get(id, metav1.GetOptions{})
	if err != nil {
		return err
	}
	if len(pipelineHistory.Status.Stages) < stageOrdinal ||
		len(pipelineHistory.Status.Stages[stageOrdinal].Steps) < stepOrdinal {
		return errors.New("invalid status")
	}
	if state == "start" {
		pipelineHistory.Status.Stages[stageOrdinal].Steps[stepOrdinal].State = v3.StateBuilding
	} else if state == "success" {
		pipelineHistory.Status.Stages[stageOrdinal].Steps[stepOrdinal].State = v3.StateSuccess
	} else if state == "fail" {
		pipelineHistory.Status.Stages[stageOrdinal].Steps[stepOrdinal].State = v3.StateFail
	} else {
		return errors.New("unknown state")
	}
	if _, err := pipelineHistoryClient.Update(pipelineHistory); err != nil {
		return err
	}

	return nil
}
