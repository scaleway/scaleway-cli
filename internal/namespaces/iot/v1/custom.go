package iot

import (
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	"github.com/scaleway/scaleway-sdk-go/api/iot/v1"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("iot").Groups = []string{"integration"}

	human.RegisterMarshalerFunc(iot.HubStatus(""), human.EnumMarshalFunc(hubStatusMarshalSpecs))
	human.RegisterMarshalerFunc(
		iot.DeviceMessageFiltersRulePolicy(""),
		human.EnumMarshalFunc(deviceMessageFiltersRulePolicyMarshalSpecs),
	)
	human.RegisterMarshalerFunc(
		iot.DeviceStatus(""),
		human.EnumMarshalFunc(deviceStatusMarshalSpecs),
	)
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
