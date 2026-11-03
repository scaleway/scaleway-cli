package datawarehouse

import (
	"github.com/scaleway/scaleway-cli/v2/core"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("datawarehouse").Groups = []string{"database"}

	cmds.MustFind("datawarehouse", "deployment", "list").Override(core.ListViewBuilder("ID", "Name", "OrganizationID", "ProjectID", "Status", "Tags", "CreatedAt", "UpdatedAt", "Version", "ReplicaCount", "ShardCount", "CPUMin", "CPUMax", "Endpoints", "RAMPerCPU", "MoveFactor", "Region"))
	cmds.MustFind("datawarehouse", "preset", "list").Override(core.ListViewBuilder("Name", "Category", "CPUMin", "CPUMax", "RAMPerCPU", "ReplicaCount", "ShardCount"))
	cmds.MustFind("datawarehouse", "version", "list").Override(core.ListViewBuilder("Version", "EndOfLifeAt"))

	cmds.MustFind("datawarehouse", "database", "list").Override(core.ListViewBuilder("Name", "Size", "DeploymentID", "Region"))
	cmds.MustFind("datawarehouse", "user", "list").Override(core.ListViewBuilder("Name", "IsAdmin", "DeploymentID", "Region"))

	return cmds
}
