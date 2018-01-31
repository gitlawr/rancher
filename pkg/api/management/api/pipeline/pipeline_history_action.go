package pipeline

import (
	"github.com/rancher/norman/types"
	"github.com/rancher/types/config"
	"github.com/sirupsen/logrus"
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

	}
	return nil
}

func (h *HistoryHandler) rerun(apiContext *types.APIContext) error {
	return nil
}

func (h *HistoryHandler) stop(apiContext *types.APIContext) error {
	return nil
}
