package monitor

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	promapi "github.com/prometheus/client_golang/api"
	promapiv1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/rancher/rancher/pkg/ref"
	"github.com/rancher/rancher/pkg/settings"
	"golang.org/x/sync/errgroup"

	"github.com/rancher/norman/httperror"
	"github.com/rancher/norman/parse"
	"github.com/rancher/norman/types"
	"github.com/rancher/norman/types/convert"
	"github.com/rancher/rancher/pkg/clustermanager"
	v3 "github.com/rancher/types/apis/management.cattle.io/v3"
	pv3 "github.com/rancher/types/apis/project.cattle.io/v3"
	"github.com/rancher/types/config"
	"github.com/rancher/types/config/dialer"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

func NewGlobalGraphHandler(dialerFactory dialer.Factory, clustermanager *clustermanager.Manager) *GlobalGraphHandler {
	return &GlobalGraphHandler{
		dialerFactory:  dialerFactory,
		clustermanager: clustermanager,
		projectLister:  clustermanager.ScaledContext.Management.Projects(metav1.NamespaceAll).Controller().Lister(),
		appLister:      clustermanager.ScaledContext.Project.Apps(metav1.NamespaceAll).Controller().Lister(),
		graphLister:    clustermanager.ScaledContext.Management.GlobalMonitorGraphs("").Controller().Lister(),
	}
}

type GlobalGraphHandler struct {
	dialerFactory  dialer.Factory
	clustermanager *clustermanager.Manager
	projectLister  v3.ProjectLister
	appLister      pv3.AppLister
	graphLister    v3.GlobalMonitorGraphLister
}

func (h *GlobalGraphHandler) QuerySeriesAction(actionName string, action *types.Action, apiContext *types.APIContext) error {
	var queryGraphInput v3.QueryGraphInput
	actionInput, err := parse.ReadBody(apiContext.Request)
	if err != nil {
		return err
	}

	if err = convert.ToObj(actionInput, &queryGraphInput); err != nil {
		return err
	}

	return h.QueryGlobalSeriesAction(actionName, action, apiContext, queryGraphInput)
}

func (h *GlobalGraphHandler) QueryGlobalSeriesAction(actionName string, action *types.Action, apiContext *types.APIContext, queryGraphInput v3.QueryGraphInput) error {
	adminClusterID := settings.AdminClusterID.Get()
	if adminClusterID == "" {
		return errors.New("admin cluster is not set")
	}

	if queryGraphInput.Filters == nil {
		queryGraphInput.Filters = map[string]string{"clusterId": adminClusterID}
	} else {
		queryGraphInput.Filters["clusterId"] = adminClusterID
	}

	inputParser := newClusterGraphInputParser(queryGraphInput)
	if err := inputParser.parse(); err != nil {
		return err
	}

	reqContext, cancel := context.WithTimeout(context.Background(), prometheusReqTimeout)
	defer cancel()

	adminClusterName := settings.AdminClusterID.Get()
	dialer, err := h.dialerFactory.ClusterDialer(adminClusterName)
	if err != nil {
		return fmt.Errorf("get dailer failed, %v", err)
	}
	prometheusQuery, err := newGlobalPrometheusQuery(reqContext, dialer)
	if err != nil {
		return err
	}

	mgmtClient := h.clustermanager.ScaledContext.Management
	nodeLister := mgmtClient.Nodes(metav1.NamespaceAll).Controller().Lister()

	nodeMap, err := getNodeName2InternalIPMap(nodeLister, "")
	if err != nil {
		return err
	}

	graphs, err := h.graphLister.List("cattle-global-data", labels.NewSelector())
	if err != nil {
		return err
	}

	logrus.Infof("get graphs: %v", graphs)
	var queries []*PrometheusQuery
	for _, g := range graphs {
		var set labels.Set
		if inputParser.Input.IsDetails && g.Spec.DetailsMetricsSelector != nil {
			set = labels.Set(g.Spec.DetailsMetricsSelector)
		} else {
			set = labels.Set(g.Spec.MetricsSelector)
		}
		metrics, err := mgmtClient.MonitorMetrics("cattle-global-data").List(metav1.ListOptions{LabelSelector: set.AsSelector().String()})
		if err != nil {
			return fmt.Errorf("list metrics failed, %v", err)
		}
		logrus.Infof("get monitorMetrics: %v", metrics)
		for _, v := range metrics.Items {
			id := getPrometheusQueryID(g.Name, g.Spec.ResourceType, v.Name)
			queries = append(queries, InitPromQuery(id, inputParser.Start, inputParser.End, inputParser.Step, v.Spec.Expression, v.Spec.LegendFormat, isInstanceGraph(g.Spec.GraphType)))
		}
	}

	logrus.Infof("get queries: %v", queries)
	seriesSlice, err := prometheusQuery.Do(queries)
	if err != nil {
		logrus.Error(err)
		logrus.WithError(err).Warn("query series failed")
		return httperror.NewAPIError(httperror.ServerError, "Failed to obtain metrics. The metrics service may not be available."+err.Error())
	}

	if seriesSlice == nil {
		apiContext.WriteResponse(http.StatusNoContent, nil)
		return nil
	}
	logrus.Infof("seriesSlice: %v", seriesSlice)

	collection := v3.QueryClusterGraphOutput{Type: "collection"}
	for k, v := range seriesSlice {
		graphName, resourceType, _ := parseID(k)
		series := convertInstance(v, nodeMap, resourceType)
		queryGraph := v3.QueryClusterGraph{
			GraphName: graphName,
			Series:    series,
		}
		collection.Data = append(collection.Data, queryGraph)
	}

	res, err := json.Marshal(collection)
	if err != nil {
		return fmt.Errorf("marshal query series result failed, %v", err)
	}

	apiContext.Response.Write(res)
	return nil
}

func newGlobalPrometheusQuery(ctx context.Context, dial dialer.Dialer) (*Queries, error) {
	url := "http://global-monitoring-thanos.cattle-global-data"
	cfg := promapi.Config{
		Address: url,
		RoundTripper: &http.Transport{
			Dial:                  dial,
			ExpectContinueTimeout: 1 * time.Second,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
		},
	}

	client, err := promapi.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("create prometheus client failed: %v", err)
	}
	q := &Queries{
		ctx: ctx,
		api: promapiv1.NewAPI(client),
	}
	q.eg, q.ctx = errgroup.WithContext(q.ctx)
	return q, nil
}

func globalGraph2Metrics(userContext *config.UserContext, mgmtClient v3.Interface, clusterName, resourceType, refGraphName string, metricSelector, detailsMetricSelector map[string]string, metricParams map[string]string, isDetails bool) ([]*metricWrap, error) {
	projectName, _ := ref.Parse(refGraphName)
	nodeLister := mgmtClient.Nodes(metav1.NamespaceAll).Controller().Lister()
	newMetricParams, err := parseMetricParams(userContext, nodeLister, resourceType, clusterName, projectName, metricParams)
	if err != nil {
		return nil, err
	}

	var excuteMetrics []*metricWrap
	var set labels.Set
	if isDetails && detailsMetricSelector != nil {
		set = labels.Set(detailsMetricSelector)
	} else {
		set = labels.Set(metricSelector)
	}
	metrics, err := mgmtClient.MonitorMetrics("cattle-global-data").List(metav1.ListOptions{LabelSelector: set.AsSelector().String()})
	if err != nil {
		return nil, fmt.Errorf("list metrics failed, %v", err)
	}

	for _, v := range metrics.Items {
		executeExpression := replaceParams(newMetricParams, v.Spec.Expression)
		logrus.Infof("wrap metrics: name/%v, expr/%v, refGraphName/%v, resourceType/%v", v.Name, executeExpression, refGraphName, resourceType)
		excuteMetrics = append(excuteMetrics, &metricWrap{
			MonitorMetric:              *v.DeepCopy(),
			ExecuteExpression:          executeExpression,
			ReferenceGraphName:         refGraphName,
			ReferenceGraphResourceType: resourceType,
		})
	}
	return excuteMetrics, nil
}
