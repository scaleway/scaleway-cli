package vpc

import (
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/human"
	"github.com/scaleway/scaleway-sdk-go/api/vpc/v2"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.Remove("vpc", "post")
	cmds.RemoveResource("vpc", "route")
	cmds.MustFind("vpc", "private-network", "get").Override(privateNetworkGetBuilder)
	human.RegisterMarshalerFunc(vpc.PrivateNetwork{}, privateNetworkMarshalerFunc)

	return cmds
}
