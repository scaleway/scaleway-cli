package applesilicon

import (
	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/internal/human"
	applesilicon "github.com/scaleway/scaleway-sdk-go/api/applesilicon/v1alpha1"
)

var (
	serverStatusMarshalSpecs = human.EnumMarshalSpecs{
		applesilicon.ServerStatusError:     &human.EnumMarshalSpec{Attribute: color.FgRed, Value: "error"},
		applesilicon.ServerStatusReady:     &human.EnumMarshalSpec{Attribute: color.FgGreen, Value: "ready"},
		applesilicon.ServerStatusRebooting: &human.EnumMarshalSpec{Attribute: color.FgBlue, Value: "rebooting"},
		applesilicon.ServerStatusStarting:  &human.EnumMarshalSpec{Attribute: color.FgBlue, Value: "starting"},
		applesilicon.ServerStatusUpdating:  &human.EnumMarshalSpec{Attribute: color.FgBlue, Value: "updating"},
	}
)
