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

type PipelineHistory struct {
	types.Namespaced

	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PipelineHistorySpec   `json:"spec"`
	Status PipelineHistoryStatus `json:"status"`
}

type PipelineLog struct {
	types.Namespaced

	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec PipelineLogSpec `json:"spec"`
}

type RemoteAccount struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RemoteAccountSpec   `json:"spec"`
	Status RemoteAccountStatus `json:"status"`
}

type GitRepoCache struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GitRepoCacheSpec   `json:"spec"`
	Status GitRepoCacheStatus `json:"status"`
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
	Scheme       string `json:"githubConfig,omitempty"`
	Host         string `json:"host,omitempty"`
	ClientId     string `json:"clientId,omitempty"`
	ClientSecret string `json:"clientSecret,omitempty"`
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
	NextRunNumber int    `json:"nextRunNumber" yaml:"nextRunNumber,omitempty" norman:"default=1,min=1"`
	LastRunId     string `json:"lastRunId,omitempty" yaml:"lastRunId,omitempty"`
	LastRunState  string `json:"lastRunState,omitempty" yaml:"lastRunState,omitempty"`
	LastRunTime   int64  `json:"lastRunTime,omitempty" yaml:"lastRunTime,omitempty"`
	NextRunTime   int64  `json:"nextRunTime,omitempty" yaml:"nextRunTime,omitempty"`
	WebHookId     string `json:"webhookId,omitempty" yaml:"webhookId,omitempty"`
	WebHookToken  string `json:"webhookToken,omitempty" yaml:"webhookToken,omitempty"`
}

type PipelineSpec struct {
	ProjectName string `json:"projectName" yaml:"projectName" norman:"required,type=reference[project]"`

	Active bool `json:"active,omitempty" yaml:"active,omitempty"`
	//DisplayName string   `json:"displayName,omitempty" yaml:"displayName,omitempty" norman:"required"`
	Triggers Triggers `json:"triggers,omitempty" yaml:"triggers,omitempty"`
	Stages   []Stage  `json:"stages,omitempty" yaml:"stages,omitempty" norman:"required"`
}

type Triggers struct {
	WebhookTrigger *WebhookTrigger `json:"webhookTrigger,omitempty" yaml:"webhookTrigger,omitempty"`
	CronTrigger    *CronTrigger    `json:"cronTrigger,omitempty" yaml:"cronTrigger,omitempty"`
}

type WebhookTrigger struct {
	Active bool `json:"active,omitempty" yaml:"active,omitempty"`
}

type CronTrigger struct {
	Active   bool   `json:"active,omitempty" yaml:"active,omitempty"`
	Timezone string `json:"timezone,omitempty" yaml:"timezone,omitempty"`
	Spec     string `json:"spec,omitempty" yaml:"spec,omitempty"`
}

type Stage struct {
	Name  string `json:"name,omitempty" yaml:"name,omitempty" norman:"required"`
	Steps []Step `json:"steps,omitempty" yaml:"steps,omitempty"`
}

type Step struct {
	Type string `json:"type,omitempty" yaml:"type,omitempty" norman:"required,options=sourcecode|runscript|buildimage|pushimage,default=runscript"`

	SourceCodeStepConfig   *SourceCodeStepConfig   `json:"sourceCodeStepConfig,omitempty" yaml:"sourceCodeStepConfig,omitempty"`
	RunScriptStepConfig    *RunScriptStepConfig    `json:"runScriptStepConfig,omitempty" yaml:"runScriptStepConfig,omitempty"`
	PublishImageStepConfig *PublishImageStepConfig `json:"publishImageStepConfig,omitempty" yaml:"publishImageStepConfig,omitempty"`
	//Step timeout in minutes
	Timeout int `json:"timeout,omitempty" yaml:"timeout,omitempty"`
}

type SourceCodeStepConfig struct {
	Repository        string `json:"repository,omitempty" yaml:"repository,omitempty" `
	Branch            string `json:"branch,omitempty" yaml:"branch,omitempty" `
	RemoteAccountName string `json:"remoteAccountName,omitempty" yaml:"remoteAccountName,omitempty" "`
}

type RunScriptStepConfig struct {
	Image       string   `json:"image,omitempty" yaml:"image,omitempty" norman:"required"`
	ShellScript string   `json:"shellScript,omitempty" yaml:"shellScript,omitempty"`
	Entrypoint  string   `json:"entrypoint,omitempty" yaml:"enrtypoint,omitempty"`
	Args        string   `json:"args,omitempty" yaml:"args,omitempty"`
	Env         []string `json:"env,omitempty" yaml:"env,omitempty"`
}

type PublishImageStepConfig struct {
	DockerfilePath string `json:"dockerFilePath,omittempty" yaml:"dockerFilePath,omitempty" norman:"required,default=./Dockerfile"`
	BuildContext   string `json:"buildContext,omitempty" yaml:"buildContext,omitempty" norman:"required,default=."`
	ImageTag       string `json:"imageTag,omitempty" yaml:"imageTag,omitempty" norman:"required,default=${CICD_GIT_REPOSITORY_NAME}:${CICD_GIT_BRANCH}"`
}

type PipelineHistorySpec struct {
	ProjectName string `json:"projectName" norman:"required,type=reference[project]"`

	//DisplayName string   `json:"displayName,omitempty" norman:"required"`
	RunNumber   int      `json:"runNumber,omitempty" norman:"required,min=1"`
	TriggerType string   `json:"triggerType,omitempty" norman:"required,options=manual|cron|webhook"`
	Pipeline    Pipeline `json:"pipeline,omitempty" norman:"required"`
}

type PipelineHistoryStatus struct {
	CommitInfo  string            `json:"commitInfo,omitempty"`
	EnvVars     map[string]string `json:"envVars,omitempty"`
	State       string            `json:"state,omitempty"`
	StartTime   int64             `json:"startTime,omitempty"`
	EndTime     int64             `json:"endTime,omitempty"`
	StageStatus []StageStatus     `json:"stageStatus,omitempty"`
}

type StageStatus struct {
	State      string       `json:"state,omitempty"`
	StartTime  int64        `json:"startTime,omitempty"`
	EndTime    int64        `json:"endTime,omitempty"`
	StepStatus []StepStatus `json:"stepStatus,omitempty"`
}

type StepStatus struct {
	State     string `json:"state,omitempty"`
	StartTime int64  `json:"startTime,omitempty"`
	EndTime   int64  `json:"endTime,omitempty"`
}

type RemoteAccountSpec struct {
	//DisplayName string `json:"displayName,omitempty" norman:"required"`
	//RemoteType        string `json:"remoteType,omitempty" norman:"required,options=github"`
	UserID      string `json:"userId" norman:"required,type=reference[user]"`
	AvatarURL   string `json:"avatarUrl,omitempty"`
	HTMLURL     string `json:"htmlUrl,omitempty"`
	Login       string `json:"login,omitempty"`
	AccountName string `json:"accountName,omitempty"`
	AccessToken string `json:"accessToken,omitempty"`
}

type RemoteAccountStatus struct {
}

type GitRepoCacheSpec struct {
	RemoteType        string          `json:"remoteType,omitempty" norman:"required,options=github"`
	UserID            string          `json:"userId" norman:"required,type=reference[user]"`
	RemoteAccountName string          `json:"remoteAccountName,omitempty" norman:"required,type=reference[remoteaccount]`
	Repositories      []GitRepository `json:"repositories,omitempty"`
}

type GitRepository struct {
	Name        string   `json:"name,omitempty"`
	CloneURL    string   `json:"cloneUrl,omitempty"`
	Permissions RepoPerm `json:"permissions,omitempty"`
	Language    string   `json:"language,omitempty"`
}

type GitRepoCacheStatus struct {
}

type RepoPerm struct {
	Pull  bool `json:"pull,omitempty"`
	Push  bool `json:"push,omitempty"`
	Admin bool `json:"admin,omitempty"`
}

type PipelineLogSpec struct {
	ProjectName string `json:"projectName" yaml:"projectName" norman:"required,type=reference[project]"`

	PipelineHistoryName string `json:"pipelineHistoryName,omitempty" norman:"type=reference[pipelinehistory]`
	StageOrdinal        int    `json:"stageOrdinal,omitempty" norman:"min=1"`
	StepOrdinal         int    `json:"stepOrdinal,omitempty" norman:"min=1"`
	Message             string `json:"message,omitempty"`
}
