package jenkins

import (
	"bytes"
	"fmt"
	"github.com/rancher/types/apis/management.cattle.io/v3"
)

func ConvertPipelineToJenkinsPipeline(pipeline *v3.Pipeline) PipelineJob {
	pipelineJob := PipelineJob{
		Plugin: WORKFLOW_JOB_PLUGIN,
		Definition: Definition{
			Class:   FLOW_DEFINITION_CLASS,
			Plugin:  FLOW_DEFINITION_PLUGIN,
			Sandbox: true,
			//TODO test script
			Script: convertPipeline(pipeline),
		},
	}

	return pipelineJob
}

func convertStep(step *v3.Step, stageOrdinal int, stepOrdinal int) string {
	stepContent := ""
	stepName := fmt.Sprintf("step_%d_%d", stageOrdinal, stepOrdinal)

	if step.SourceCodeConfig != nil {
		stepContent = fmt.Sprintf("git '%s'", step.SourceCodeConfig.Url)
	} else if step.RunScriptConfig == nil {
		stepContent = fmt.Sprintf("sh \"\"\"\n%s\n\"\"\"", step.RunScriptConfig.ShellScript)
	} else if step.PublishImageConfig == nil {
		stepContent = fmt.Sprintf(`sh """"\necho dopublishimage\n"""`)
	} else {
		return ""
	}
	return fmt.Sprintf(stepBlock, stepName, stepName, stepName, stepContent)
}

func convertStage(stage *v3.Stage, stageOrdinal int) string {
	var buffer bytes.Buffer
	for i, step := range stage.Steps {
		buffer.WriteString(convertStep(&step, stageOrdinal, i+1))
		if i != len(stage.Steps)-1 {
			buffer.WriteString(",")
		}
	}

	return fmt.Sprintf(stageBlock, stage.Name, buffer.String())
}

func convertPipeline(pipeline *v3.Pipeline) string {
	var containerbuffer bytes.Buffer
	var pipelinebuffer bytes.Buffer
	for j, stage := range pipeline.Spec.Stages {
		pipelinebuffer.WriteString(convertStage(&stage, j+1))
		pipelinebuffer.WriteString("\n")
		for k, step := range stage.Steps {
			stepName := fmt.Sprintf("step_%d_%d", j+1, k+1)
			image := ""
			if step.SourceCodeConfig != nil {
				image = "alpine/git"
			} else if step.RunScriptConfig != nil {
				image = step.RunScriptConfig.Image
			} else if step.PublishImageConfig != nil {
				image = "docker"
			} else {
				return ""
			}
			containerDef := fmt.Sprintf(containerBlock, stepName, image)
			containerbuffer.WriteString(containerDef)

		}
	}

	return fmt.Sprintf(pipelineBlock, containerbuffer.String(), pipelinebuffer.String())
}

const stageBlock = `stage('%s'){
parallel %s
}
`

const stepBlock = `'%s': {
  stage('%s'){
    container(name: '%s') {
      %s
    }
  }
}
`

const pipelineBlock = `def label = "buildpod.${env.JOB_NAME}.${env.BUILD_NUMBER}".replace('-', '_').replace('/', '_')
podTemplate(label: label, containers: [
%s
containerTemplate(name: 'jnlp', image: 'jenkinsci/jnlp-slave:alpine', envVars: [
envVar(key: 'JENKINS_URL', value: 'http://jenkins:8080')], args: '${computer.jnlpmac} ${computer.name}', ttyEnabled: false)]) {
node(label) {
timestamps {
%s
}
}
}`

const containerBlock = `containerTemplate(name: '%s', image: '%s', ttyEnabled: true, command: 'cat'),`
