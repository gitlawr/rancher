package jenkins

import (
	"encoding/xml"
	"fmt"

	"github.com/rancher/types/apis/management.cattle.io/v3"
	"github.com/rancher/types/config"
	"github.com/sirupsen/logrus"
)

type JenkinsEngine struct {
	Client  *Client
	Cluster *config.ClusterContext
}

func (j JenkinsEngine) RunPipeline(pipeline *v3.Pipeline, triggerType string) error {

	jobName := JENKINS_JOB_PREFIX + pipeline.Name

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

	jobName := JENKINS_JOB_PREFIX + pipeline.Name
	bconf, _ := xml.MarshalIndent(jobconf, "  ", "    ")
	if err := j.Client.CreateJob(jobName, bconf); err != nil {
		return err
	}
	return nil
}

func (j JenkinsEngine) UpdatePipelineJob(pipeline *v3.Pipeline) error {
	logrus.Debug("update jenkins job for pipeline")
	jobconf := ConvertPipelineToJenkinsPipeline(pipeline)

	jobName := JENKINS_JOB_PREFIX + pipeline.Name
	bconf, _ := xml.MarshalIndent(jobconf, "  ", "    ")
	if err := j.Client.UpdateJob(jobName, bconf); err != nil {
		return err
	}
	return nil
}

func (j JenkinsEngine) RerunHistory(history *v3.PipelineExecution) error {
	return j.RunPipeline(&history.Spec.Pipeline, TriggerTypeManual)
}

func (j JenkinsEngine) StopHistory(history *v3.PipelineExecution) error {
	jobName := JENKINS_JOB_PREFIX + history.Spec.Pipeline.Name
	buildNumber := history.Spec.Run
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
