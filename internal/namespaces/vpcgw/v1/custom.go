package vpcgw

import (
	"github.com/scaleway/scaleway-cli/v2/internal/core"
	"github.com/scaleway/scaleway-cli/v2/internal/human"
	"github.com/scaleway/scaleway-sdk-go/api/vpcgw/v1"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	human.RegisterMarshalerFunc(vpcgw.GatewayNetworkStatus(""), human.EnumMarshalFunc(gatewayNetworkStatusMarshalSpecs))
	human.RegisterMarshalerFunc(vpcgw.GatewayStatus(""), human.EnumMarshalFunc(gatewayStatusMarshalSpecs))
	human.RegisterMarshalerFunc(vpcgw.Gateway{}, gatewayNetworkMarshallerFunc)

	cmds.MustFind("vpc-gw", "gateway-type", "list").Override(vpcgwGatewayTypeListBuilder)
	cmds.MustFind("vpc-gw", "gateway", "create").Override(gatewayCreateBuilder)
	cmds.MustFind("vpc-gw", "gateway-network", "create").Override(gatewayNetworkCreateBuilder)
	cmds.MustFind("vpc-gw", "gateway-network", "delete").Override(gatewayNetworkDeleteBuilder)

	return cmds
}
