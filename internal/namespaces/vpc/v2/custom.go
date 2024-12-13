package vpc

import (
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/human"
	"github.com/scaleway/scaleway-sdk-go/api/vpc/v2"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.Remove("vpc", "post")
	cmds.MustFind("vpc", "private-network", "get").Override(privateNetworkGetBuilder)
	human.RegisterMarshalerFunc(vpc.PrivateNetwork{}, privateNetworkMarshalerFunc)

	return cmds
}
