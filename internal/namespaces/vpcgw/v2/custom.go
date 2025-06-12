package vpcgw

import (
	"strings"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	"github.com/scaleway/scaleway-sdk-go/api/vpcgw/v2"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()
	for _, cmd := range cmds.GetAll() {
		if cmd.Verb != "" && !strings.HasSuffix(cmd.Verb, "-v2") {
			cmd.Verb = strings.TrimSpace(cmd.Verb) + "-v2"
		}
	}

	cmds.MustFind("vpc-gw").Groups = []string{"network"}

	human.RegisterMarshalerFunc(
		vpcgw.GatewayNetworkStatus(""),
		human.EnumMarshalFunc(gatewayNetworkStatusMarshalSpecs))
	human.RegisterMarshalerFunc(
		vpcgw.GatewayStatus(""),
		human.EnumMarshalFunc(gatewayStatusMarshalSpecs))
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
