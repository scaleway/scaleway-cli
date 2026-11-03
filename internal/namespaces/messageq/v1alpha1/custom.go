package messageq

import "github.com/scaleway/scaleway-cli/v2/core"

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("messageq", "deployment", "list").Override(core.ListViewBuilder("ID", "Name", "OrganizationID", "ProjectID", "Status", "Tags", "NodeCount", "NodeType", "Endpoints", "CreatedAt", "UpdatedAt", "Version", "Region"))
	cmds.MustFind("messageq", "node-type", "list").Override(core.ListViewBuilder("StockStatus", "Name", "Description", "Vcpus", "MemoryBytes", "Disabled", "Beta", "InstanceRange", "AvailableVolumeTypes", "Region"))
	cmds.MustFind("messageq", "version", "list").Override(core.ListViewBuilder("Version", "EndOfLife", "Disabled", "Beta", "Region"))

	cmds.MustFind("messageq", "user", "list").Override(core.ListViewBuilder("Username", "DeploymentID", "Region"))

	return cmds
}
