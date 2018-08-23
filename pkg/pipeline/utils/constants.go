package utils

const (
	PipelineNamespace              = "cattle-pipeline"
	PipelineNamespaceSuffix        = "-pipeline"
	JenkinsName                    = "jenkins"
	PipelineName                   = "pipeline"
	PipelineSecretName             = "pipeline-secret"
	PipelineSecretUserKey          = "admin-user"
	PipelineSecretTokenKey         = "admin-token"
	PipelineSecretRegistryTokenKey = "registry-token"
	PipelineSecretRegistryAuthPath = "/opt/auth/"
	PipelineSecretRegistryCrtPath  = "/opt/crt/"
	PipelineSecretDefaultUser      = "admin"
	PipelineSecretDefaultToken     = "admin123"
	RegistryName                   = "docker-registry"
	RegistryProxyName              = "registry-proxy"
	RegistryCrt                    = "domain.crt"
	RegistryCAKey                  = "ca.crt"
	RegistryCrtSecretName          = "registry-crt"
	RegistryCrtVolumeName          = "crt"
	RegistryAuthVolumeName         = "auth"
	RegistryPortMappingKey         = "mappings"
	RegistryPortMappingFile        = "mappings.yaml"
	PublishSecretUserKey           = "username"
	PublishSecretPwKey             = "password"
	DockerCredentialName           = "pipeline-docker-registry"
	ProxyConfigMapName             = "proxy-mappings"
	MinioName                      = "minio"
	MinioBucketLocation            = "local"
	MinioLogBucket                 = "pipeline-logs"
	NetWorkPolicyName              = "pipeline-np"
	LabelKeyApp                    = "app"
	LabelKeyJenkins                = "jenkins"
	JenkinsMaster                  = "master"
	LabelKeyExecution              = "execution"
	DefaultRegistry                = "index.docker.io"
	LocalRegistry                  = "docker-registry"
	DefaultTag                     = "latest"
	JenkinsPort                    = 8080
	JenkinsJNLPPort                = 50000
	RegistryPort                   = 443
	MinioPort                      = 9000

	WebhookEventPush        = "push"
	WebhookEventPullRequest = "pull_request"
	WebhookEventTag         = "tag"

	TriggerTypeUser    = "user"
	TriggerTypeWebhook = "webhook"

	StateWaiting  = "Waiting"
	StateBuilding = "Building"
	StateSuccess  = "Success"
	StateFailed   = "Failed"
	StateSkipped  = "Skipped"
	StateAborted  = "Aborted"
	StateQueueing = "Queueing"
	StatePending  = "Pending"
	StateDenied   = "Denied"

	PipelineFinishLabel    = "pipeline.project.cattle.io/finish"
	LocalRegistryPortLabel = "pipeline.project.cattle.io/localRegistryPort"

	PipelineFileYml  = ".rancher-pipeline.yml"
	PipelineFileYaml = ".rancher-pipeline.yaml"

	EnvGitRepoName       = "CICD_GIT_REPO_NAME"
	EnvGitURL            = "CICD_GIT_URL"
	EnvGitCommit         = "CICD_GIT_COMMIT"
	EnvGitRef            = "CICD_GIT_REF"
	EnvGitBranch         = "CICD_GIT_BRANCH"
	EnvGitTag            = "CICD_GIT_TAG"
	EnvTriggerType       = "CICD_TRIGGER_TYPE"
	EnvEvent             = "CICD_EVENT"
	EnvExecutionID       = "CICD_EXECUTION_ID"
	EnvExecutionSequence = "CICD_EXECUTION_SEQUENCE"
	EnvPipelineID        = "CICD_PIPELINE_ID"
	EnvProjectID         = "CICD_PROJECT_ID"
	EnvClusterID         = "CICD_CLUSTER_ID"
	EnvRegistry          = "CICD_REGISTRY"
	EnvImageRepo         = "CICD_IMAGE_REPO"
	EnvLocalRegistry     = "CICD_LOCAL_REGISTRY"

	SettingExecutorQuota        = "executor-quota"
	SettingExecutorQuotaDefault = "2"

	DefaultTimeout = 60
)
