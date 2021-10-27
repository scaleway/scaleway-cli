package vpcgw

import (
	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/internal/human"
	"github.com/scaleway/scaleway-sdk-go/api/vpcgw/v1"
)

var (
	gatewayNetworkStatusMarshalSpecs = human.EnumMarshalSpecs{
		vpcgw.GatewayNetworkStatusAttaching:   &human.EnumMarshalSpec{Attribute: color.FgBlue},
		vpcgw.GatewayNetworkStatusConfiguring: &human.EnumMarshalSpec{Attribute: color.FgBlue},
		vpcgw.GatewayNetworkStatusCreated:     &human.EnumMarshalSpec{Attribute: color.FgGreen},
		vpcgw.GatewayNetworkStatusDeleted:     &human.EnumMarshalSpec{Attribute: color.FgRed},
		vpcgw.GatewayNetworkStatusDetaching:   &human.EnumMarshalSpec{Attribute: color.FgBlue},
		vpcgw.GatewayNetworkStatusReady:       &human.EnumMarshalSpec{Attribute: color.FgGreen},
		vpcgw.GatewayNetworkStatusUnknown:     &human.EnumMarshalSpec{Attribute: color.Faint},
	}
)
