package baremetal

import "github.com/scaleway/scaleway-cli/v2/core"

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("baremetal", "private-network", "list").
		Override(core.ListViewBuilder("ID", "ProjectID", "ServerID", "PrivateNetworkID", "Vlan", "Status", "CreatedAt", "UpdatedAt", "IpamIPIDs"))

	return cmds
}
