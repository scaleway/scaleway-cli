package vpc

import (
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	humanv2 "github.com/scaleway/scaleway-cli/v2/core/human/v2"
	"github.com/scaleway/scaleway-sdk-go/api/vpc/v2"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("vpc").Groups = []string{"network"}

	cmds.Remove("vpc", "post")
	cmds.MustFind("vpc", "private-network", "get").Override(privateNetworkGetBuilder)
	cmds.MustFind("vpc", "private-network", "delete").Override(privateNetworkDeleteBuilder)
	humanv2.RegisterLayout(reflect.TypeFor[vpc.PrivateNetwork](), privateNetworkLayout)

	cmds.Merge(core.NewCommands(
		vpcACLEditCommand(),
		privateNetworkCockpitCommand(),
	))

	return cmds
}
