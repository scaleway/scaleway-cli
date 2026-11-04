package vpcgw

import (
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	humanv2 "github.com/scaleway/scaleway-cli/v2/core/human/v2"
	"github.com/scaleway/scaleway-sdk-go/api/vpcgw/v2"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()
	for _, cmd := range cmds.GetAll() {
		if cmd.Verb != "" {
			cmd.Aliases = append(cmd.Aliases, cmd.Verb+"-v2")
		}
	}

	cmds.MustFind("vpc-gw").Groups = []string{"network"}

	humanv2.RegisterEnumFormat(gatewayNetworkStatusFormat)
	humanv2.RegisterEnumFormat(gatewayStatusFormat)
	human.RegisterMarshalerFunc(
		vpcgw.Gateway{},
		gatewayMarshalerFunc)

	cmds.MustFind("vpc-gw", "gateway-type", "list").Override(vpcgwGatewayTypeListBuilder)
	cmds.MustFind("vpc-gw", "gateway", "create").Override(gatewayCreateBuilder)
	cmds.MustFind("vpc-gw", "gateway-network", "create").Override(gatewayNetworkCreateBuilder)
	cmds.MustFind("vpc-gw", "gateway-network", "delete").Override(gatewayNetworkDeleteBuilder)

	cmds.Merge(core.NewCommands(
		vpcgwPATRulesEditCommand(),
	))

	return cmds
}
