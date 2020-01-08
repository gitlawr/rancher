package globalmonitoring

import (
	"context"

	"github.com/rancher/types/config"
)

// Register initializes the controllers and registers
func Register(ctx context.Context, management *config.ManagementContext) {
	ah := appHandler{
		appLister:     management.Project.Apps("").Controller().Lister(),
		appClient:     management.Project.Apps(""),
		projectLister: management.Management.Projects("").Controller().Lister(),
		clusterLister: management.Management.Clusters("").Controller().Lister(),
		secretClient:  management.Core.Secrets(""),
		secretLister:  management.Core.Secrets("").Controller().Lister(),
		userLister:    management.Management.Users("").Controller().Lister(),
		userManager:   management.UserManager,
	}

	ch := clusterHandler{
		appLister:     management.Project.Apps("").Controller().Lister(),
		appClient:     management.Project.Apps(""),
		projectLister: management.Management.Projects("").Controller().Lister(),
		deployer: deployer{
			graphClient:  management.Management.GlobalMonitorGraphs(""),
			graphLister:  management.Management.GlobalMonitorGraphs("").Controller().Lister(),
			metricClient: management.Management.MonitorMetrics(""),
			metricLister: management.Management.MonitorMetrics("").Controller().Lister(),
		},
	}

	management.Project.Apps("").AddHandler(ctx, "globalmonitoring-app-handler", ah.sync)
	management.Management.Clusters("").AddHandler(ctx, "globalmonitoring-cluster-handler", ch.sync)

}
