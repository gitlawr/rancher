package jenkins

import (
	"encoding/xml"
	"fmt"

	"github.com/pkg/errors"
	"github.com/rancher/rancher/pkg/pipeline/utils"
	"github.com/rancher/types/apis/management.cattle.io/v3"
	"github.com/rancher/types/config"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

type JenkinsEngine struct {
	Client  *Client
	Cluster *config.ClusterContext
}

func (j JenkinsEngine) RunPipeline(pipeline *v3.Pipeline, triggerType string) error {

	jobName := getJobName(pipeline)

	if _, err := j.Client.GetJobInfo(jobName); err == ErrJenkinsJobNotFound {
		if err = j.CreatePipelineJob(pipeline); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	if err := j.UpdatePipelineJob(pipeline); err != nil {
		return err
	}

	if _, err := j.Client.BuildJob(jobName, map[string]string{}); err != nil {
		logrus.Errorf("run %s error:%v", jobName, err)
		return err
	}
	return nil
}

func (j JenkinsEngine) CreatePipelineJob(pipeline *v3.Pipeline) error {
	logrus.Debug("create jenkins job for pipeline")
	jobconf := ConvertPipelineToJenkinsPipeline(pipeline)

	jobName := getJobName(pipeline)
	bconf, _ := xml.MarshalIndent(jobconf, "  ", "    ")
	if err := j.Client.CreateJob(jobName, bconf); err != nil {
		return err
	}
	return nil
}

func (j JenkinsEngine) UpdatePipelineJob(pipeline *v3.Pipeline) error {
	logrus.Debug("update jenkins job for pipeline")
	jobconf := ConvertPipelineToJenkinsPipeline(pipeline)

	jobName := getJobName(pipeline)
	bconf, _ := xml.MarshalIndent(jobconf, "  ", "    ")
	if err := j.Client.UpdateJob(jobName, bconf); err != nil {
		return err
	}
	return nil
}

func (j JenkinsEngine) RerunHistory(history *v3.PipelineExecution) error {
	return j.RunPipeline(&history.Spec.Pipeline, TriggerTypeManual)
}

func (j JenkinsEngine) StopHistory(execution *v3.PipelineExecution) error {
	jobName := getJobName(&execution.Spec.Pipeline)
	buildNumber := execution.Spec.Run
	info, err := j.Client.GetJobInfo(jobName)
	if err == ErrJenkinsJobNotFound {
		return nil
	} else if err != nil {
		return err
	}
	if info.InQueue {
		//delete in queue
		//TODO filter build number
		queueItem, ok := info.QueueItem.(map[string]interface{})
		if !ok {
			return fmt.Errorf("type assertion fail for queueitem")
		}
		queueId, ok := queueItem["id"].(float64)
		if !ok {
			return fmt.Errorf("type assertion fail for queueId")
		}
		if err := j.Client.CancelQueueItem(int(queueId)); err != nil {
			return fmt.Errorf("cancel queueitem error:%v", err)
		}
	} else {
		buildInfo, err := j.Client.GetBuildInfo(jobName, buildNumber)
		if err != nil {
			return err
		}
		if buildInfo.Building {
			if err := j.Client.StopJob(jobName, buildNumber); err != nil {
				return err
			}
		}
	}
	return nil
}

func (j JenkinsEngine) SyncExecution(execution *v3.PipelineExecution) (bool, error) {
	jobName := getJobName(&execution.Spec.Pipeline)
	info, err := j.Client.GetWFBuildInfo(jobName)
	if err != nil {
		return false, err
	}
	updated := false
	for _, jenkinsStage := range info.Stages {
		//handle those in step-1-1 format
		parts := strings.Split(jenkinsStage.Name, "-")
		if len(parts) == 3 {
			stage, err := strconv.Atoi(parts[1])
			if err != nil {
				return false, err
			}
			step, err := strconv.Atoi(parts[2])
			if err != nil {
				return false, err
			}
			if len(execution.Status.Stages) <= stage || len(execution.Status.Stages[stage].Steps) <= step {
				return false, errors.New("error sync execution - index out of range")
			}
			status := jenkinsStage.Status
			if status == "SUCCESS" && execution.Status.Stages[stage].Steps[step].State != v3.StateSuccess {
				updated = true
				successStep(execution, stage, step, jenkinsStage)
			} else if status == "FAILED" && execution.Status.Stages[stage].Steps[step].State != v3.StateFail {
				updated = true
				failStep(execution, stage, step, jenkinsStage)
			} else if status == "IN_PROGRESS" && execution.Status.Stages[stage].Steps[step].State != v3.StateBuilding {
				updated = true
				buildingStep(execution, stage, step, jenkinsStage)
			}
		}
	}

	return updated, nil
}

func successStep(execution *v3.PipelineExecution, stage int, step int, jenkinsStage JenkinsStage) {

	startTime := time.Unix(jenkinsStage.StartTimeMillis/1000, 0).Format(time.RFC3339)
	endTime := time.Unix((jenkinsStage.StartTimeMillis+jenkinsStage.DurationMillis)/1000, 0).Format(time.RFC3339)
	execution.Status.Stages[stage].Steps[step].State = v3.StateSuccess
	if execution.Status.Stages[stage].Steps[step].Started == "" {
		execution.Status.Stages[stage].Steps[step].Started = startTime
	}
	execution.Status.Stages[stage].Steps[step].Ended = endTime
	if execution.Status.Stages[stage].Started == "" {
		execution.Status.Stages[stage].Started = startTime
	}
	if execution.Status.Started == "" {
		execution.Status.Started = startTime
	}
	if utils.IsStageSuccess(execution.Status.Stages[stage]) {
		execution.Status.Stages[stage].State = v3.StateSuccess
		execution.Status.Stages[stage].Ended = endTime
		if stage == len(execution.Status.Stages)-1 {
			execution.Status.State = v3.StateSuccess
			execution.Status.Ended = endTime
		}
	}
}

func failStep(execution *v3.PipelineExecution, stage int, step int, jenkinsStage JenkinsStage) {

	startTime := time.Unix(jenkinsStage.StartTimeMillis/1000, 0).Format(time.RFC3339)
	endTime := time.Unix((jenkinsStage.StartTimeMillis+jenkinsStage.DurationMillis)/1000, 0).Format(time.RFC3339)
	execution.Status.Stages[stage].Steps[step].State = v3.StateFail
	execution.Status.Stages[stage].State = v3.StateFail
	execution.Status.State = v3.StateFail
	if execution.Status.Stages[stage].Steps[step].Started == "" {
		execution.Status.Stages[stage].Steps[step].Started = startTime
	}
	execution.Status.Stages[stage].Steps[step].Ended = endTime
	if execution.Status.Stages[stage].Started == "" {
		execution.Status.Stages[stage].Started = startTime
	}
	if execution.Status.Stages[stage].Ended == "" {
		execution.Status.Stages[stage].Ended = endTime
	}
	if execution.Status.Started == "" {
		execution.Status.Started = startTime
	}
	if execution.Status.Ended == "" {
		execution.Status.Ended = endTime
	}
}

func buildingStep(execution *v3.PipelineExecution, stage int, step int, jenkinsStage JenkinsStage) {
	startTime := time.Unix(jenkinsStage.StartTimeMillis/1000, 0).Format(time.RFC3339)
	execution.Status.Stages[stage].Steps[step].State = v3.StateBuilding
	if execution.Status.Stages[stage].Steps[step].Started == "" {
		execution.Status.Stages[stage].Steps[step].Started = startTime
	}
	if execution.Status.Stages[stage].Started == "" {
		execution.Status.Stages[stage].Started = startTime
	}
	if execution.Status.Started == "" {
		execution.Status.Started = startTime
	}
}

//OnActivityCompelte helps clean up
func (j JenkinsEngine) OnHistoryCompelte(history *v3.PipelineExecution) {
	//TODO
	return
	/*
		//clean related container by label
		command := fmt.Sprintf("docker ps --filter label=activityid=%s -q | xargs docker rm -f", activity.Id)
		cleanServiceScript := fmt.Sprintf(ScriptSkel, activity.NodeName, strings.Replace(command, "\"", "\\\"", -1))
		logrus.Debugf("cleanservicescript is: %v", cleanServiceScript)
		res, err := ExecScript(cleanServiceScript)
		logrus.Debugf("clean services result:%v,%v", res, err)
		if err != nil {
			logrus.Errorf("error cleanning up on worker node: %v, got result '%s'", err, res)
		}
		logrus.Infof("activity '%s' complete", activity.Id)
		//clean workspace
		if !activity.Pipeline.KeepWorkspace {
			command = "rm -rf ${System.getenv('JENKINS_HOME')}/workspace/" + activity.Id
			cleanWorkspaceScript := fmt.Sprintf(ScriptSkel, activity.NodeName, strings.Replace(command, "\"", "\\\"", -1))
			res, err = ExecScript(cleanWorkspaceScript)
			if err != nil {
				logrus.Errorf("error cleanning up on worker node: %v, got result '%s'", err, res)
			}
			logrus.Debugf("clean workspace result:%v,%v", res, err)
		}
	*/
}

func (j JenkinsEngine) GetStepLog(history *v3.PipelineExecution, stageOrdinal int, stepOrdinal int, paras map[string]interface{}) (string, error) {
	//TODO
	return "", nil
	/*
		if stageOrdinal < 0 || stageOrdinal >= len(activity.ActivityStages) || stepOrdinal < 0 || stepOrdinal >= len(activity.ActivityStages[stageOrdinal].ActivitySteps) {
			return "", errors.New("ordinal out of range")
		}
		jobName := getJobName(activity, stageOrdinal, stepOrdinal)
		var logText *string
		if val, ok := paras["prevLog"]; ok {
			logText = val.(*string)
		}
		startLine := len(strings.Split(*logText, "\n"))

		rawOutput, err := GetBuildRawOutput(jobName, startLine)
		if err != nil {
			return "", err
		}
		token := "\\n\\w{14}\\s{2}\\[.*?\\].*?\\.sh"
		*logText = *logText + rawOutput
		outputs := regexp.MustCompile(token).Split(*logText, -1)
		if len(outputs) > 1 && stageOrdinal == 0 && stepOrdinal == 0 {
			// SCM
			return trimFirstLine(outputs[1]), nil
		}
		if len(outputs) < 3 {
			//no printed log
			return "", nil
		}
		//hide set +x
		return trimFirstLine(outputs[2]), nil
	*/
}

func getJobName(pipeline *v3.Pipeline) string {
	return fmt.Sprintf("%s%s-%d", JENKINS_JOB_PREFIX, pipeline.Name, pipeline.Status.NextRun)
}
