package pipeline

import (
	"context"
	"github.com/gorilla/websocket"
	"github.com/rancher/norman/types"
	"github.com/rancher/rancher/pkg/pipeline/engine"
	"github.com/rancher/rancher/pkg/ref"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	logSyncInterval = 2 * time.Second
	writeWait       = time.Second
	Threshold       = 100000
)

var upgrader = websocket.Upgrader{
	HandshakeTimeout: 5 * time.Second,
	CheckOrigin:      func(r *http.Request) bool { return true },
	Error:            onError,
}

func onError(rw http.ResponseWriter, _ *http.Request, code int, err error) {
	rw.WriteHeader(code)
	rw.Write([]byte(err.Error()))
}

func (h *ExecutionHandler) handleLog(apiContext *types.APIContext) error {
	stageInput := apiContext.Request.URL.Query().Get("stage")
	stepInput := apiContext.Request.URL.Query().Get("step")
	stage, err := strconv.Atoi(stageInput)
	if err != nil {
		return err
	}
	step, err := strconv.Atoi(stepInput)
	if err != nil {
		return err
	}
	ns, name := ref.Parse(apiContext.ID)
	execution, err := h.PipelineExecutionLister.Get(ns, name)
	if err != nil {
		return err
	}
	clusterName, _ := ref.Parse(execution.Spec.ProjectName)
	userContext, err := h.ClusterManager.UserContext(clusterName)
	if err != nil {
		return err
	}

	pipelineEngine := engine.New(userContext)

	c, err := upgrader.Upgrade(apiContext.Response, apiContext.Request, nil)
	if err != nil {
		return err
	}
	defer c.Close()

	cancelCtx, cancel := context.WithCancel(apiContext.Request.Context())
	readerGroup, ctx := errgroup.WithContext(cancelCtx)
	apiContext.Request = apiContext.Request.WithContext(ctx)

	go func() {
		for {
			if _, _, err := c.NextReader(); err != nil {
				cancel()
				c.Close()
				break
			}
		}
	}()

	go func() {
		readerGroup.Wait()
	}()

	syncT := time.NewTicker(logSyncInterval)
	defer syncT.Stop()

	prevLog := ""
	for {
		select {
		case <-syncT.C:
			execution, err = h.PipelineExecutionLister.Get(ns, name)
			if err != nil {
				logrus.Debugf("error in execution get: %v", err)
				if prevLog == "" {
					writeData(c, []byte("Log is unavailable."))
				}
				c.WriteControl(websocket.CloseMessage, []byte{}, time.Now().Add(writeWait))
				return nil
			}
			log, err := pipelineEngine.GetStepLog(execution, stage, step)
			if err != nil {
				logrus.Debug(err)
				if prevLog == "" {
					writeData(c, []byte("Log is unavailable."))
				}
				c.WriteControl(websocket.CloseMessage, []byte{}, time.Now().Add(writeWait))
				return nil
			}
			newLog := getNewLog(prevLog, log)
			prevLog = log
			if newLog != "" {
				if err := writeData(c, []byte(newLog)); err != nil {
					logrus.Debugf("error in writeData: %v", err)
					return nil
				}
			}
			if execution.Status.Stages[stage].Steps[step].Ended != "" {
				c.WriteControl(websocket.CloseMessage, []byte{}, time.Now().Add(writeWait))
				return nil
			}
		}
	}
}

func writeData(c *websocket.Conn, buf []byte) error {
	messageWriter, err := c.NextWriter(websocket.TextMessage)
	if err != nil {
		return err
	}

	if _, err := messageWriter.Write(buf); err != nil {
		return err
	}

	return messageWriter.Close()
}

func getNewLog(prevLog string, currLog string) string {
	if len(prevLog) < Threshold {
		return strings.TrimPrefix(currLog, prevLog)
	}
	//long logs are trimmed so we use previous log tail to do comparison
	prevLogTail := prevLog[len(prevLog)-Threshold:]
	idx := strings.Index(currLog, prevLogTail)
	if idx >= 0 {
		return currLog[idx+Threshold:]
	}
	return currLog
}
