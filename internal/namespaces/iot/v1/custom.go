package iot

import (
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	humanv2 "github.com/scaleway/scaleway-cli/v2/core/human/v2"
	"github.com/scaleway/scaleway-sdk-go/api/iot/v1"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("iot").Groups = []string{"integration"}

	humanv2.RegisterEnumFormat(hubStatusFormat)
	humanv2.RegisterEnumFormat(deviceMessageFiltersRulePolicyFormat)
	humanv2.RegisterEnumFormat(deviceStatusFormat)
	human.RegisterMarshalerFunc(
		&iot.CreateNetworkResponse{},
		iotNetworkCreateResponseMarshalerFunc,
	)
	human.RegisterMarshalerFunc(
		&iot.CreateDeviceResponse{},
		iotDeviceCreateResponseMarshalerFunc,
	)

	cmds.MustFind("iot", "hub", "create").Override(hubCreateBuilder)

	return cmds
}
