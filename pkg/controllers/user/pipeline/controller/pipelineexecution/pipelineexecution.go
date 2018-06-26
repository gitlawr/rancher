package pipelineexecution

import (
	"context"
	"github.com/pkg/errors"
	"github.com/rancher/rancher/pkg/pipeline/engine"
	"github.com/rancher/rancher/pkg/pipeline/utils"
	"github.com/rancher/rancher/pkg/ref"
	"github.com/rancher/types/apis/apps/v1beta2"
	"github.com/rancher/types/apis/core/v1"
	networkv1 "github.com/rancher/types/apis/networking.k8s.io/v1"
	"github.com/rancher/types/apis/project.cattle.io/v3"
	rbacv1 "github.com/rancher/types/apis/rbac.authorization.k8s.io/v1"
	"github.com/rancher/types/config"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"strconv"
	"strings"
)

const (
	projectIDLabel   = "field.cattle.io/projectId"
	roleCreateNs     = "create-ns"
	roleGetNs        = "get-ns"
	roleEditNsSuffix = "-namespaces-edit"
	roleAdmin        = "admin"
)

//Lifecycle is responsible for initializing logs for pipeline execution
//and calling the run for the execution.
type Lifecycle struct {
	namespaceLister v1.NamespaceLister
	namespaces      v1.NamespaceInterface
	secrets         v1.SecretInterface
	services        v1.ServiceInterface
	serviceAccounts v1.ServiceAccountInterface
	networkPolicies networkv1.NetworkPolicyInterface

	clusterRoleBindings rbacv1.ClusterRoleBindingInterface
	roleBindings        rbacv1.RoleBindingInterface
	deployments         v1beta2.DeploymentInterface

	pipelineLister             v3.PipelineLister
	pipelines                  v3.PipelineInterface
	pipelineExecutionLister    v3.PipelineExecutionLister
	pipelineExecutions         v3.PipelineExecutionInterface
	pipelineSettingLister      v3.PipelineSettingLister
	pipelineEngine             engine.PipelineEngine
	sourceCodeCredentialLister v3.SourceCodeCredentialLister
}

func Register(ctx context.Context, cluster *config.UserContext) {
	clusterName := cluster.ClusterName

	namespaces := cluster.Core.Namespaces("")
	namespaceLister := cluster.Core.Namespaces("").Controller().Lister()
	secrets := cluster.Core.Secrets("")
	services := cluster.Core.Services("")
	serviceAccounts := cluster.Core.ServiceAccounts("")
	networkPolicies := cluster.Networking.NetworkPolicies("")
	clusterRoleBindings := cluster.RBAC.ClusterRoleBindings("")
	roleBindings := cluster.RBAC.RoleBindings("")
	deployments := cluster.Apps.Deployments("")

	pipelines := cluster.Management.Project.Pipelines("")
	pipelineLister := pipelines.Controller().Lister()
	pipelineExecutions := cluster.Management.Project.PipelineExecutions("")
	pipelineExecutionLister := pipelineExecutions.Controller().Lister()
	pipelineSettingLister := cluster.Management.Project.PipelineSettings("").Controller().Lister()
	sourceCodeCredentialLister := cluster.Management.Project.SourceCodeCredentials("").Controller().Lister()

	pipelineEngine := engine.New(cluster)
	pipelineExecutionLifecycle := &Lifecycle{
		namespaces:          namespaces,
		namespaceLister:     namespaceLister,
		secrets:             secrets,
		services:            services,
		serviceAccounts:     serviceAccounts,
		networkPolicies:     networkPolicies,
		clusterRoleBindings: clusterRoleBindings,
		roleBindings:        roleBindings,
		deployments:         deployments,

		pipelineLister:             pipelineLister,
		pipelines:                  pipelines,
		pipelineExecutionLister:    pipelineExecutionLister,
		pipelineExecutions:         pipelineExecutions,
		pipelineSettingLister:      pipelineSettingLister,
		pipelineEngine:             pipelineEngine,
		sourceCodeCredentialLister: sourceCodeCredentialLister,
	}
	stateSyncer := &ExecutionStateSyncer{
		clusterName: clusterName,

		pipelineLister:          pipelineLister,
		pipelines:               pipelines,
		pipelineExecutionLister: pipelineExecutionLister,
		pipelineExecutions:      pipelineExecutions,
		pipelineEngine:          pipelineEngine,
	}

	pipelineExecutions.AddClusterScopedLifecycle(pipelineExecutionLifecycle.GetName(), cluster.ClusterName, pipelineExecutionLifecycle)

	go stateSyncer.sync(ctx, syncStateInterval)

}

func (l *Lifecycle) Create(obj *v3.PipelineExecution) (*v3.PipelineExecution, error) {
	return l.Sync(obj)
}

func (l *Lifecycle) Updated(obj *v3.PipelineExecution) (*v3.PipelineExecution, error) {
	return l.Sync(obj)
}

func (l *Lifecycle) Sync(obj *v3.PipelineExecution) (*v3.PipelineExecution, error) {

	if obj == nil {
		return obj, nil
	}

	//doIfAbort
	if obj.Status.ExecutionState == utils.StateAborted {
		if err := l.doStop(obj); err != nil {
			return obj, err
		}
	}

	//doIfFinish
	if obj.Labels != nil && obj.Labels[utils.PipelineFinishLabel] == "true" {
		//start a queueing execution if there is any
		if err := l.startQueueingExecution(obj); err != nil {
			return obj, err
		}
		return obj, nil
	}

	//doIfRunning
	if v3.PipelineExecutionConditionInitialized.GetStatus(obj) != "" {
		return obj, nil
	}

	//doIfExceedQuota
	exceed, err := l.exceedQuota(obj)
	if err != nil {
		return obj, err
	}
	if exceed {
		obj.Status.ExecutionState = utils.StateQueueing
		obj.Labels[utils.PipelineFinishLabel] = ""

		if err := l.newExecutionUpdateLastRunState(obj); err != nil {
			return obj, err
		}

		return obj, nil
	}

	//doIfOnCreation
	if err := l.newExecutionUpdateLastRunState(obj); err != nil {
		return obj, err
	}
	v3.PipelineExecutionConditionInitialized.CreateUnknownIfNotExists(obj)
	obj.Labels[utils.PipelineFinishLabel] = "false"

	if err := l.deploy(obj); err != nil {
		obj.Labels[utils.PipelineFinishLabel] = "true"
		obj.Status.ExecutionState = utils.StateFailed
		v3.PipelineExecutionConditionInitialized.False(obj)
		v3.PipelineExecutionConditionInitialized.ReasonAndMessageFromError(obj, err)
	}

	return obj, nil
}

func (l *Lifecycle) newExecutionUpdateLastRunState(obj *v3.PipelineExecution) error {

	ns, name := ref.Parse(obj.Spec.PipelineName)
	pipeline, err := l.pipelineLister.Get(ns, name)
	if err != nil {
		return err
	}
	if obj.Spec.Run == pipeline.Status.NextRun {
		pipeline.Status.NextRun++
		pipeline.Status.LastExecutionID = ref.Ref(obj)
		pipeline.Status.LastStarted = obj.Status.Started
	}
	if pipeline.Status.LastExecutionID == ref.Ref(obj) {
		pipeline.Status.LastRunState = obj.Status.ExecutionState
	}
	_, err = l.pipelines.Update(pipeline)
	return err
}
func (l *Lifecycle) startQueueingExecution(obj *v3.PipelineExecution) error {
	_, projectID := ref.Parse(obj.Spec.ProjectName)
	set := labels.Set(map[string]string{utils.PipelineFinishLabel: ""})
	queueingExecutions, err := l.pipelineExecutionLister.List(projectID, set.AsSelector())
	if err != nil {
		return err
	}
	if len(queueingExecutions) == 0 {
		return nil
	}
	oldestTime := queueingExecutions[0].CreationTimestamp
	toRunExecution := queueingExecutions[0]
	for _, e := range queueingExecutions {
		if e.CreationTimestamp.Before(&oldestTime) {
			oldestTime = e.CreationTimestamp
			toRunExecution = e
		}
	}
	toRunExecution = toRunExecution.DeepCopy()
	toRunExecution.Status.ExecutionState = utils.StateWaiting
	_, err = l.pipelineExecutions.Update(toRunExecution)
	return err
}

func (l *Lifecycle) exceedQuota(obj *v3.PipelineExecution) (bool, error) {
	_, projectID := ref.Parse(obj.Spec.ProjectName)
	quotaSetting, err := l.pipelineSettingLister.Get(projectID, utils.SettingExecutorQuota)
	if apierrors.IsNotFound(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	quotaStr := quotaSetting.Default
	if quotaSetting.Value != "" {
		quotaStr = quotaSetting.Value
	}
	quota, err := strconv.Atoi(quotaStr)
	if err != nil || quota <= 0 {
		return false, nil
	}
	set := labels.Set(map[string]string{utils.PipelineFinishLabel: "false"})
	runningExecutions, err := l.pipelineExecutionLister.List(projectID, set.AsSelector())
	if err != nil {
		return false, err
	}
	if len(runningExecutions) >= quota {
		return true, nil
	}
	return false, nil
}

func (l *Lifecycle) doStop(obj *v3.PipelineExecution) error {

	if v3.PipelineExecutionConditionInitialized.IsTrue(obj) &&
		v3.PipelineExecutionConditionBuilt.IsUnknown(obj) {
		if err := l.pipelineEngine.StopExecution(obj); err != nil {
			return err
		}
		if _, err := l.pipelineEngine.SyncExecution(obj); err != nil {
			return err
		}
	}
	v3.PipelineExecutionConditionBuilt.False(obj)
	v3.PipelineExecutionConditionBuilt.Message(obj, "aborted by user")
	for i := range obj.Status.Stages {
		if obj.Status.Stages[i].State == utils.StateBuilding {
			obj.Status.Stages[i].State = utils.StateAborted
		}
		for j := range obj.Status.Stages[i].Steps {
			if obj.Status.Stages[i].Steps[j].State == utils.StateBuilding {
				obj.Status.Stages[i].Steps[j].State = utils.StateAborted
			}
		}
	}

	//check and update lastrunstate of the pipeline when necessary
	ns, name := ref.Parse(obj.Spec.PipelineName)
	p, err := l.pipelineLister.Get(ns, name)
	if err != nil && !apierrors.IsNotFound(err) {
		return err
	}

	if p != nil && p.Status.LastExecutionID == obj.Namespace+":"+obj.Name &&
		p.Status.LastRunState != obj.Status.ExecutionState {
		p.Status.LastRunState = obj.Status.ExecutionState
		if _, err := l.pipelines.Update(p); err != nil {
			return err
		}
	}

	return nil
}

func (l *Lifecycle) Remove(obj *v3.PipelineExecution) (*v3.PipelineExecution, error) {
	return obj, nil
}

func (l *Lifecycle) GetName() string {
	return "pipeline-execution-controller"
}

//reconcileRb grant access to pipeline service account inside project namespaces
func (l *Lifecycle) reconcileRb(obj *v3.PipelineExecution) error {
	commonName := utils.GetPipelineCommonName(obj)
	_, projectID := ref.Parse(obj.Spec.ProjectName)

	namespaces, err := l.namespaceLister.List("", labels.NewSelector())
	if err != nil {
		return errors.Wrapf(err, "Error list cluster namespaces")
	}
	var namespacesInProject []*corev1.Namespace
	for _, namespace := range namespaces {
		parts := strings.Split(namespace.Annotations[projectIDLabel], ":")
		if len(parts) == 2 && parts[1] == projectID {
			namespacesInProject = append(namespacesInProject, namespace)
		} else {
			if err := l.roleBindings.DeleteNamespaced(namespace.Name, commonName, &metav1.DeleteOptions{}); err != nil && !apierrors.IsNotFound(err) {
				return err
			}
		}
	}
	for _, namespace := range namespacesInProject {
		rb := getRoleBindings(namespace.Name, commonName)
		if _, err := l.roleBindings.Create(rb); err != nil && !apierrors.IsAlreadyExists(err) {
			return errors.Wrapf(err, "Error create role binding")
		}
	}

	clusterRbs := []string{roleCreateNs, projectID + roleEditNsSuffix}
	for _, crbName := range clusterRbs {
		crb := getClusterRoleBindings(commonName, crbName)
		if _, err := l.clusterRoleBindings.Create(crb); err != nil && !apierrors.IsAlreadyExists(err) {
			return errors.Wrapf(err, "Error create cluster role binding")
		}
	}

	return nil
}
