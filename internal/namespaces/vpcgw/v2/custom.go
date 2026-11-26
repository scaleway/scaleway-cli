package vpcgw

import (
	"github.com/scaleway/scaleway-cli/v2/core"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()
	for _, cmd := range cmds.GetAll() {
		if cmd.Verb != "" {
			cmd.Aliases = append(cmd.Aliases, cmd.Verb+"-v2")
		}
	}

	cmds.MustFind("vpc-gw").Groups = []string{"network"}

	// TODO
	// human.RegisterMarshalerFunc(
	// 	vpcgw.GatewayNetworkStatus(""),
	// 	human.EnumMarshalFunc(gatewayNetworkStatusMarshalSpecs))
	// human.RegisterMarshalerFunc(
	// 	vpcgw.GatewayStatus(""),
	// 	human.EnumMarshalFunc(gatewayStatusMarshalSpecs))
	// human.RegisterMarshalerFunc(
	// 	vpcgw.Gateway{},
	// 	gatewayMarshalerFunc)

	cmds.MustFind("vpc-gw", "gateway-type", "list").Override(vpcgwGatewayTypeListBuilder)
	cmds.MustFind("vpc-gw", "gateway", "create").Override(gatewayCreateBuilder)
	cmds.MustFind("vpc-gw", "gateway-network", "create").Override(gatewayNetworkCreateBuilder)
	cmds.MustFind("vpc-gw", "gateway-network", "delete").Override(gatewayNetworkDeleteBuilder)

	cmds.Merge(core.NewCommands(
		vpcgwPATRulesEditCommand(),
	))

	return cmds
}
