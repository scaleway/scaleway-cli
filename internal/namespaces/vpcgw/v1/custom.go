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

	return cmds
}
