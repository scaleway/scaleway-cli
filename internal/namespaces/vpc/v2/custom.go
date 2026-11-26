package vpc

import (
	"github.com/scaleway/scaleway-cli/v2/core"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("vpc").Groups = []string{"network"}

	cmds.Remove("vpc", "post")
	cmds.MustFind("vpc", "private-network", "get").Override(privateNetworkGetBuilder)
	cmds.MustFind("vpc", "private-network", "delete").Override(privateNetworkDeleteBuilder)
	// human.RegisterMarshalerFunc(vpc.PrivateNetwork{}, privateNetworkMarshalerFunc)

	cmds.Merge(core.NewCommands(
		vpcACLEditCommand(),
		privateNetworkCockpitCommand(),
	))

	return cmds
}
