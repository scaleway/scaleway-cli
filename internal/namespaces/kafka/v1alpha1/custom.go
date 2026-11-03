package kafka

import "github.com/scaleway/scaleway-cli/v2/core"

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("kafka", "cluster", "list").Override(core.ListViewBuilder("ID", "Name", "ProjectID", "OrganizationID", "Status", "Version", "Tags", "Settings", "NodeAmount", "NodeType", "Endpoints", "CreatedAt", "UpdatedAt", "Region"))
	cmds.MustFind("kafka", "node-type", "list").Override(core.ListViewBuilder("Name", "StockStatus", "Description", "Vcpus", "MemoryBytes", "AvailableVolumeTypes", "Disabled", "Beta", "ClusterRange"))
	cmds.MustFind("kafka", "version", "list").Override(core.ListViewBuilder("Version", "EndOfLifeAt", "AvailableSettings"))

	cmds.MustFind("kafka", "users", "list").Override(core.ListViewBuilder("Username"))

	return cmds
}
