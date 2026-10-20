package mongodb

import "github.com/scaleway/scaleway-cli/v2/core"

func addCockpitURLs(cmds *core.Commands) {
	const dashboard = "scw-mongodb-overview"

	cmds.MustFind("mongodb").CockpitDashboard = dashboard
	cmds.MustFind("mongodb", "instance", "get").CockpitDashboard = dashboard
	cmds.MustFind("mongodb", "instance", "get").CockpitResourceVar = "cl_name"
}
