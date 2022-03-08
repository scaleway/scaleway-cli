package iot

import (
	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/internal/human"
	"github.com/scaleway/scaleway-sdk-go/api/iot/v1"
)

var (
	deviceMessageFiltersRulePolicyMarshalSpecs = human.EnumMarshalSpecs{
		iot.DeviceMessageFiltersRulePolicyAccept: &human.EnumMarshalSpec{Attribute: color.FgGreen, Value: "accept"},
		iot.DeviceMessageFiltersRulePolicyReject: &human.EnumMarshalSpec{Attribute: color.FgRed, Value: "reject"},
	}

	deviceStatusMarshalSpecs = human.EnumMarshalSpecs{
		iot.DeviceStatusEnabled:  &human.EnumMarshalSpec{Attribute: color.FgGreen, Value: "enabled"},
		iot.DeviceStatusDisabled: &human.EnumMarshalSpec{Attribute: color.FgRed, Value: "disabled"},
		iot.DeviceStatusError:    &human.EnumMarshalSpec{Attribute: color.FgRed, Value: "error"},
	}
)
