package vpc

import (
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	"github.com/scaleway/scaleway-sdk-go/api/vpc/v2"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("vpc").Groups = []string{"network"}

	cmds.Remove("vpc", "post")
	cmds.MustFind("vpc", "private-network", "get").Override(privateNetworkGetBuilder)
	human.RegisterMarshalerFunc(vpc.PrivateNetwork{}, privateNetworkMarshalerFunc)

	cmds.Merge(core.NewCommands(
		vpcACLEditCommand(),
	))

	return cmds
}
