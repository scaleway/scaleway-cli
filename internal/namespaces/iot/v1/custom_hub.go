package iot

import (
	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/internal/human"
	"github.com/scaleway/scaleway-sdk-go/api/iot/v1"
)

var (
	hubStatusMarshalSpecs = human.EnumMarshalSpecs{
		iot.HubStatusDisabled:  &human.EnumMarshalSpec{Attribute: color.FgRed, Value: "disabled"},
		iot.HubStatusDisabling: &human.EnumMarshalSpec{Attribute: color.FgBlue, Value: "disabling"},
		iot.HubStatusEnabling:  &human.EnumMarshalSpec{Attribute: color.FgBlue, Value: "enabling"},
		iot.HubStatusError:     &human.EnumMarshalSpec{Attribute: color.FgRed, Value: "error"},
		iot.HubStatusReady:     &human.EnumMarshalSpec{Attribute: color.FgGreen, Value: "ready"},
	}
)
