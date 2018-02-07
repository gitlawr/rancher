package handler

import "net/http"
import (
	"fmt"
	"github.com/rancher/types/apis/management.cattle.io/v3"
	"github.com/rancher/types/config"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

type WebhookHandler struct {
	Management *config.ManagementContext
}

func (h *WebhookHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	id := req.FormValue("pipelineId")
	logrus.Info("receieve webhook,id:%s", id)
	if err := h.runPipeline(id); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
	}
	rw.WriteHeader(http.StatusOK)

}

func (h *WebhookHandler) runPipeline(id string) error {
	pipeline, err := h.Management.Management.Pipelines("").Get(id, v1.GetOptions{})
	if err != nil {
		return err
	}
	ns := pipeline.Namespace
	pipelineClient := h.Management.Management.Pipelines(ns)

	historyClient := h.Management.Management.PipelineExecutions(ns)
	history := &v3.PipelineExecution{
		ObjectMeta: v1.ObjectMeta{
			Name: getNextHistoryName(pipeline),
		},
		Spec: v3.PipelineExecutionSpec{
			Run:         pipeline.Status.NextRun,
			TriggeredBy: v3.TriggerTypeWebhook,
			Pipeline:    *pipeline,
			//DisplayName: getNextHistoryName(pipeline),
		},
	}
	if _, err := historyClient.Create(history); err != nil {
		return err
	}
	pipeline.Status.NextRun++
	pipeline.Status.LastExecutionId = history.Name
	pipeline.Status.LastStarted = time.Now().String()

	_, err = pipelineClient.Update(pipeline)
	if err != nil {
		return err
	}
	return nil
}

func getNextHistoryName(p *v3.Pipeline) string {
	if p == nil {
		return ""
	}
	return fmt.Sprintf("%s-%d", p.Name, p.Status.NextRun)
}
