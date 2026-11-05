package iot

import (
	"github.com/fatih/color"
	humanv2 "github.com/scaleway/scaleway-cli/v2/core/human/v2"
	"github.com/scaleway/scaleway-sdk-go/api/iot/v1"
)

var (
	deviceMessageFiltersRulePolicyFormat = humanv2.EnumFormats[iot.DeviceMessageFiltersRulePolicy]{
		iot.DeviceMessageFiltersRulePolicyAccept: humanv2.EnumFormat{
			Color: color.FgGreen,
			Value: "accept",
		},
		iot.DeviceMessageFiltersRulePolicyReject: humanv2.EnumFormat{
			Color: color.FgRed,
			Value: "reject",
		},
	}

	deviceStatusFormat = humanv2.EnumFormats[iot.DeviceStatus]{
		iot.DeviceStatusEnabled: humanv2.EnumFormat{
			Color: color.FgGreen,
			Value: "enabled",
		},
		iot.DeviceStatusDisabled: humanv2.EnumFormat{Color: color.FgRed, Value: "disabled"},
		iot.DeviceStatusError:    humanv2.EnumFormat{Color: color.FgRed, Value: "error"},
	}

	iotDeviceCreateLayout = humanv2.Layout{
		Sections: []humanv2.Section{
			{
				Path:   "Device",
				Title:  "Device",
				Footer: "WARNING: The keys below are displayed only once. Make sure to store them securely",
			},
			{
				Path:  "Certificate.Crt",
				Title: "Certificate",
			},
			{
				Path:  "Certificate.Key",
				Title: "Private Key",
			},
		},
	}
)
