package ipam

import "github.com/scaleway/scaleway-cli/v2/core"

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("ipam").Groups = []string{"network"}

	cmds.MustFind("ipam", "ip", "list").Override(core.ListViewBuilder("ID", "Address", "ProjectID", "IsIPv6", "CreatedAt", "UpdatedAt", "Tags", "Reverses", "Region", "Zone"))

	return cmds
}
