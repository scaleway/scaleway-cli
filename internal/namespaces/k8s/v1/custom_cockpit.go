package k8s

import "github.com/scaleway/scaleway-cli/v2/core"

func addCockpitURLs(cmds *core.Commands) {
	const dashboard = "scw-kubernetes-cluster-overview"

	cmds.MustFind("k8s").CockpitDashboard = dashboard
	cmds.MustFind("k8s", "cluster", "get").CockpitDashboard = dashboard
	cmds.MustFind("k8s", "cluster", "get").CockpitResourceVar = "cluster_name"
}
