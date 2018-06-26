package pipeline

import (
	"github.com/rancher/norman/httperror"
	"github.com/rancher/norman/types"
	"github.com/rancher/rancher/pkg/clustermanager"
	"github.com/rancher/rancher/pkg/pipeline/utils"
	"github.com/rancher/rancher/pkg/ref"
	"github.com/rancher/types/apis/project.cattle.io/v3"
	"time"
)

type ExecutionHandler struct {
	ClusterManager *clustermanager.Manager

	PipelineExecutionLister v3.PipelineExecutionLister
	PipelineExecutions      v3.PipelineExecutionInterface
}

func (h *ExecutionHandler) ExecutionFormatter(apiContext *types.APIContext, resource *types.RawResource) {
	if e, ok := resource.Values["executionState"].(string); ok && e != utils.StateBuilding && e != utils.StateWaiting {
		resource.AddAction(apiContext, "rerun")
	}
	if e, ok := resource.Values["executionState"].(string); ok && (e == utils.StateBuilding || e == utils.StateWaiting) {
		resource.AddAction(apiContext, "stop")
	}
	resource.Links["log"] = apiContext.URLBuilder.Link("log", resource)
}

func (h *ExecutionHandler) LinkHandler(apiContext *types.APIContext, next types.RequestHandler) error {
	if apiContext.Link == "log" {
		return h.handleLog(apiContext)
	}

	return httperror.NewAPIError(httperror.NotFound, "Link not found")

}

func (h *ExecutionHandler) ActionHandler(actionName string, action *types.Action, apiContext *types.APIContext) error {

	switch actionName {
	case "rerun":
		return h.rerun(apiContext)
	case "stop":
		return h.stop(apiContext)
	}

	return httperror.NewAPIError(httperror.InvalidAction, "unsupported action")
}

func (h *ExecutionHandler) rerun(apiContext *types.APIContext) error {
	ns, name := ref.Parse(apiContext.ID)
	execution, err := h.PipelineExecutionLister.Get(ns, name)
	if err != nil {
		return err
	}

	//reset execution
	toUpdate := execution.DeepCopy()
	toUpdate.Status.ExecutionState = utils.StateWaiting
	toUpdate.Status.Started = time.Now().Format(time.RFC3339)
	toUpdate.Status.Ended = ""
	toUpdate.Status.Conditions = nil
	toUpdate.Labels[utils.PipelineFinishLabel] = "false"
	for i := 0; i < len(toUpdate.Status.Stages); i++ {
		stage := &toUpdate.Status.Stages[i]
		stage.State = utils.StateWaiting
		stage.Started = ""
		stage.Ended = ""
		for j := 0; j < len(stage.Steps); j++ {
			step := &stage.Steps[j]
			step.State = utils.StateWaiting
			step.Started = ""
			step.Ended = ""
		}
	}
	if _, err := h.PipelineExecutions.Update(toUpdate); err != nil {
		return err
	}

	return nil
}

func (h *ExecutionHandler) stop(apiContext *types.APIContext) error {
	ns, name := ref.Parse(apiContext.ID)
	execution, err := h.PipelineExecutionLister.Get(ns, name)
	if err != nil {
		return err
	}

	toUpdate := execution.DeepCopy()
	toUpdate.Status.ExecutionState = utils.StateAborted
	toUpdate.Status.Ended = time.Now().Format(time.RFC3339)
	toUpdate.Labels[utils.PipelineFinishLabel] = "true"
	if _, err := h.PipelineExecutions.Update(toUpdate); err != nil {
		return err
	}
	return nil
}
