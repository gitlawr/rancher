package v3

import (
	"github.com/rancher/norman/condition"
	"github.com/rancher/norman/types"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	StepTypeSourceCode = "sourceCode"
	StepTypeRunScript  = "runScript"
	StepTypeBuildImage = "buildImage"
	TriggerTypeCron    = "cron"
	TriggerTypeManual  = "manual"
	TriggerTypeWebhook = "webhook"

	StateWaiting  = "Waiting"
	StateBuilding = "Building"
	StateSuccess  = "Success"
	StateFail     = "Fail"
	StateSkip     = "Skipped"
	StateAbort    = "Abort"
	StatePending  = "Pending"
	StateDenied   = "Denied"
)

type ClusterPipeline struct {
	types.Namespaced

	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ClusterPipelineSpec   `json:"spec"`
	Status ClusterPipelineStatus `json:"status"`
}

type Pipeline struct {
	types.Namespaced

	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PipelineSpec   `json:"spec"`
	Status PipelineStatus `json:"status" yaml:"-"`
}

type PipelineExecution struct {
	types.Namespaced

	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PipelineExecutionSpec   `json:"spec"`
	Status PipelineExecutionStatus `json:"status"`
}

type PipelineExecutionLog struct {
	types.Namespaced

	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec PipelineExecutionLogSpec `json:"spec"`
}

type SourceCodeCredential struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SourceCodeCredentialSpec   `json:"spec"`
	Status SourceCodeCredentialStatus `json:"status"`
}

type SourceCodeRepository struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SourceCodeRepositorySpec   `json:"spec"`
	Status SourceCodeRepositoryStatus `json:"status"`
}

type ClusterPipelineSpec struct {
	ClusterName  string        `json:"clusterName" norman:"type=reference[cluster]"`
	Deploy       bool          `json:"deploy"`
	GithubConfig *GithubConfig `json:"githubConfig,omitempty"`
}

type ClusterPipelineStatus struct {
	Conditions []PipelineCondition `json:"conditions,omitempty"`
}

type GithubConfig struct {
	TLS          bool   `json:"tls,omitempty"`
	Host         string `json:"host,omitempty"`
	ClientId     string `json:"clientId,omitempty"`
	ClientSecret string `json:"clientSecret,omitempty"`
	RedirectUrl  string `json:"redirectUrl,omitempty"`
}

var (
	ClusterPipelineConditionInitialized condition.Cond = "Initialized"
	ClusterPipelineConditionProvisioned condition.Cond = "Provisioned"
)

type PipelineCondition struct {
	// Type of cluster condition.
	Type condition.Cond `json:"type"`
	// Status of the condition, one of True, False, Unknown.
	Status v1.ConditionStatus `json:"status"`
	// The last time this condition was updated.
	LastUpdateTime string `json:"lastUpdateTime,omitempty"`
	// Last time the condition transitioned from one status to another.
	LastTransitionTime string `json:"lastTransitionTime,omitempty"`
	// The reason for the condition's last transition.
	Reason string `json:"reason,omitempty"`
	// Human-readable message indicating details about last transition
	Message string `json:"message,omitempty"`
}

type PipelineStatus struct {
	NextRun         int    `json:"nextRun" yaml:"nextRun,omitempty" norman:"default=1,min=1"`
	LastExecutionId string `json:"lastExecutionId,omitempty" yaml:"lastExecutionId,omitempty"`
	LastRunState    string `json:"lastRunState,omitempty" yaml:"lastRunState,omitempty"`
	LastStarted     string `json:"lastStarted,omitempty" yaml:"lastStarted,omitempty"`
	NextStart       string `json:"nextStart,omitempty" yaml:"nextStart,omitempty"`
	WebHookId       string `json:"webhookId,omitempty" yaml:"webhookId,omitempty"`
	Token           string `json:"token,omitempty" yaml:"token,omitempty"`
}

type PipelineSpec struct {
	ProjectName string `json:"projectName" yaml:"projectName" norman:"required,type=reference[project]"`

	//DisplayName string   `json:"displayName,omitempty" yaml:"displayName,omitempty" norman:"required"`
	TriggerWebhook        bool   `json:"triggerWebhook,omitempty" yaml:"triggerWebhook,omitempty"`
	TriggerCronTimezone   string `json:"triggerCronTimezone,omitempty" yaml:"triggerCronTimezone,omitempty"`
	TriggerCronExpression string `json:"triggerCronExpression,omitempty" yaml:"triggerCronExpression,omitempty"`

	Stages []Stage `json:"stages,omitempty" yaml:"stages,omitempty" norman:"required"`
}

type Stage struct {
	Name  string `json:"name,omitempty" yaml:"name,omitempty" norman:"required"`
	Steps []Step `json:"steps,omitempty" yaml:"steps,omitempty"`
}

type Step struct {
	SourceCodeConfig   *SourceCodeConfig   `json:"sourceCodeConfig,omitempty" yaml:"sourceCodeConfig,omitempty"`
	RunScriptConfig    *RunScriptConfig    `json:"runScriptConfig,omitempty" yaml:"runScriptConfig,omitempty"`
	PublishImageConfig *PublishImageConfig `json:"publishImageConfig,omitempty" yaml:"publishImageConfig,omitempty"`
	//Step timeout in minutes
	Timeout int `json:"timeout,omitempty" yaml:"timeout,omitempty"`
}

type SourceCodeConfig struct {
	Url                      string `json:"url,omitempty" yaml:"url,omitempty" `
	Branch                   string `json:"branch,omitempty" yaml:"branch,omitempty" `
	SourceCodeCredentialName string `json:"sourceCodeCredentialName,omitempty" yaml:"sourceCodeCredentialName,omitempty" norman:"type=reference[sourceCodeCredential]"`
}

type RunScriptConfig struct {
	Image       string   `json:"image,omitempty" yaml:"image,omitempty" norman:"required"`
	IsShell     bool     `json:"isShell,omitempty" yaml:"isShell,omitempty"`
	ShellScript string   `json:"shellScript,omitempty" yaml:"shellScript,omitempty"`
	Entrypoint  string   `json:"entrypoint,omitempty" yaml:"enrtypoint,omitempty"`
	Command     string   `json:"command,omitempty" yaml:"command,omitempty"`
	Env         []string `json:"env,omitempty" yaml:"env,omitempty"`
}

type PublishImageConfig struct {
	DockerfilePath string `json:"dockerfilePath,omittempty" yaml:"dockerfilePath,omitempty" norman:"required,default=./Dockerfile"`
	BuildContext   string `json:"buildContext,omitempty" yaml:"buildContext,omitempty" norman:"required,default=."`
	Tag            string `json:"tag,omitempty" yaml:"tag,omitempty" norman:"required,default=${CICD_GIT_REPOSITORY_NAME}:${CICD_GIT_BRANCH}"`
}

type PipelineExecutionSpec struct {
	ProjectName     string   `json:"projectName" norman:"required,type=reference[project]"`
	PipelineName    string   `json:"pipelineName" norman:"type=reference[pipeline]"`
	Run             int      `json:"run,omitempty" norman:"required,min=1"`
	TriggeredBy     string   `json:"triggeredBy,omitempty" norman:"required,options=user|cron|webhook"`
	TriggerUserName string   `json:"triggerUserName,omitempty" norman:"type=reference[user]"`
	Pipeline        Pipeline `json:"pipeline,omitempty" norman:"required"`
}

type PipelineExecutionStatus struct {
	Commit  string        `json:"commit,omitempty"`
	State   string        `json:"state,omitempty"`
	Started string        `json:"started,omitempty"`
	Ended   string        `json:"ended,omitempty"`
	Stages  []StageStatus `json:"stages,omitempty"`
}

type StageStatus struct {
	State   string       `json:"state,omitempty"`
	Started string       `json:"started,omitempty"`
	Ended   string       `json:"ended,omitempty"`
	Steps   []StepStatus `json:"steps,omitempty"`
}

type StepStatus struct {
	State   string `json:"state,omitempty"`
	Started string `json:"started,omitempty"`
	Ended   string `json:"ended,omitempty"`
}

type SourceCodeCredentialSpec struct {
	ClusterName    string `json:"clusterName" norman:"required,type=reference[cluster]"`
	SourceCodeType string `json:"sourceCodeType,omitempty" norman:"required,options=github"`
	UserName       string `json:"userName" norman:"required,type=reference[user]"`
	AvatarURL      string `json:"avatarUrl,omitempty"`
	HTMLURL        string `json:"htmlUrl,omitempty"`
	LoginName      string `json:"loginName,omitempty"`
	DisplayName    string `json:"displayName,omitempty"`
	AccessToken    string `json:"accessToken,omitempty"`
}

type SourceCodeCredentialStatus struct {
}

type SourceCodeRepositorySpec struct {
	ClusterName              string   `json:"clusterName" norman:"required,type=reference[cluster]"`
	SourceCodeType           string   `json:"sourceCodeType,omitempty" norman:"required,options=github"`
	UserName                 string   `json:"userName" norman:"required,type=reference[user]"`
	SourceCodeCredentialName string   `json:"sourceCodeCredentialName,omitempty" norman:"required,type=reference[sourceCodeCredential]`
	Url                      string   `json:"url,omitempty"`
	Permissions              RepoPerm `json:"permissions,omitempty"`
	Language                 string   `json:"language,omitempty"`
}

type SourceCodeRepositoryStatus struct {
}

type RepoPerm struct {
	Pull  bool `json:"pull,omitempty"`
	Push  bool `json:"push,omitempty"`
	Admin bool `json:"admin,omitempty"`
}

type PipelineExecutionLogSpec struct {
	ProjectName string `json:"projectName" yaml:"projectName" norman:"required,type=reference[project]"`

	PipelineExecutionName string `json:"pipelineExecutionName,omitempty" norman:"type=reference[pipelineexecution]`
	Stage                 int    `json:"stage,omitempty" norman:"min=1"`
	Step                  int    `json:"step,omitempty" norman:"min=1"`
	Message               string `json:"message,omitempty"`
}

type AuthAppInput struct {
	Type         string `json:"type,omitempty" norman:"type=string,required"`
	RedirectUrl  string `json:"redirectUrl,omitempty" norman:"type=string"`
	TLS          bool   `json:"tls,omitempty"`
	Host         string `json:"host,omitempty"`
	ClientId     string `json:"clientId,omitempty" norman:"type=string,required"`
	ClientSecret string `json:"clientSecret,omitempty" norman:"type=string,required"`
	Code         string `json:"code,omitempty" norman:"type=string,required"`
}

type AuthUserInput struct {
	Type        string `json:"type,omitempty" norman:"type=string,required"`
	RedirectUrl string `json:"redirectUrl,omitempty" norman:"type=string"`
	Code        string `json:"code,omitempty" norman:"type=string,required"`
}
