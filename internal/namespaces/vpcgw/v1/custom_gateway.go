package vpcgw

import (
	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/internal/human"
	"github.com/scaleway/scaleway-sdk-go/api/vpcgw/v1"
)

var (
	gatewayStatusMarshalSpecs = human.EnumMarshalSpecs{
		vpcgw.GatewayStatusAllocating:  &human.EnumMarshalSpec{Attribute: color.FgBlue},
		vpcgw.GatewayStatusConfiguring: &human.EnumMarshalSpec{Attribute: color.FgBlue},
		vpcgw.GatewayStatusDeleted:     &human.EnumMarshalSpec{Attribute: color.FgRed},
		vpcgw.GatewayStatusDeleting:    &human.EnumMarshalSpec{Attribute: color.FgBlue},
		vpcgw.GatewayStatusFailed:      &human.EnumMarshalSpec{Attribute: color.FgRed},
		vpcgw.GatewayStatusRunning:     &human.EnumMarshalSpec{Attribute: color.FgGreen},
		vpcgw.GatewayStatusStopped:     &human.EnumMarshalSpec{Attribute: color.FgRed},
		vpcgw.GatewayStatusStopping:    &human.EnumMarshalSpec{Attribute: color.FgBlue},
		vpcgw.GatewayStatusUnknown:     &human.EnumMarshalSpec{Attribute: color.Faint},
	}
)
