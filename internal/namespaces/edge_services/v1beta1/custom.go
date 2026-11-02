package edge_services

import "github.com/scaleway/scaleway-cli/v2/core"

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("edge-services").Groups = []string{"network"}

	cmds.MustFind("edge-services", "plan", "list").Override(planListBuilder)

	cmds.Merge(core.NewCommands(
		edgeServicesRouteRulesEditCommand(),
	))

	return cmds
}
