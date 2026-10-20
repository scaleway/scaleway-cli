package redis

import "github.com/scaleway/scaleway-cli/v2/core"

func addCockpitURLs(cmds *core.Commands) {
	const dashboard = "scw-redis-overview"

	cmds.MustFind("redis").CockpitDashboard = dashboard
	cmds.MustFind("redis", "cluster", "get").CockpitDashboard = dashboard
	cmds.MustFind("redis", "cluster", "get").CockpitResourceVar = "resource_name"
}
