package iot

import (
	"reflect"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	humanv2 "github.com/scaleway/scaleway-cli/v2/core/human/v2"
	"github.com/scaleway/scaleway-sdk-go/api/iot/v1"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("iot").Groups = []string{"integration"}

	humanv2.RegisterEnumFormat(humanv2.EnumFormats[iot.DeviceMessageFiltersRulePolicy]{
		iot.DeviceMessageFiltersRulePolicyAccept: humanv2.EnumFormat{Color: color.FgGreen, Value: "accept"},
		iot.DeviceMessageFiltersRulePolicyReject: humanv2.EnumFormat{Color: color.FgRed, Value: "reject"},
	})
	humanv2.RegisterEnumFormat(humanv2.EnumFormats[iot.DeviceStatus]{
		iot.DeviceStatusEnabled:  humanv2.EnumFormat{Color: color.FgGreen, Value: "enabled"},
		iot.DeviceStatusDisabled: humanv2.EnumFormat{Color: color.FgRed, Value: "disabled"},
		iot.DeviceStatusError:    humanv2.EnumFormat{Color: color.FgRed, Value: "error"},
	})
	humanv2.RegisterEnumFormat(humanv2.EnumFormats[iot.HubStatus]{
		iot.HubStatusDisabled:  humanv2.EnumFormat{Color: color.FgRed, Value: "disabled"},
		iot.HubStatusDisabling: humanv2.EnumFormat{Color: color.FgBlue, Value: "disabling"},
		iot.HubStatusEnabling:  humanv2.EnumFormat{Color: color.FgBlue, Value: "enabling"},
		iot.HubStatusError:     humanv2.EnumFormat{Color: color.FgRed, Value: "error"},
		iot.HubStatusReady:     humanv2.EnumFormat{Color: color.FgGreen, Value: "ready"},
	})

	humanv2.RegisterLayout(reflect.TypeFor[*iot.CreateDeviceResponse](), humanv2.Layout{
		Sections: []humanv2.Section{
			{Path: "Device", Main: true, Footer: "WARNING: The keys below are displayed only once. Make sure to store them securely"},
			{Path: "Certificate.Crt", Title: "Certificate"},
			{Path: "Certificate.Key", Title: "Private Key"},
		},
	})
	humanv2.RegisterLayout(reflect.TypeFor[*iot.CreateNetworkResponse](), humanv2.Layout{
		Sections: []humanv2.Section{
			{Path: "Network", Main: true, Footer: "WARNING: The secret below is displayed only once, we do not keep it. Make sure to store it securely"},
			{Path: "Secret"},
		},
	})

	cmds.MustFind("iot", "hub", "create").Override(hubCreateBuilder)

	return cmds
}
