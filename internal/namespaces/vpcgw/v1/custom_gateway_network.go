package vpcgw

import (
	"context"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/internal/core"
	"github.com/scaleway/scaleway-cli/v2/internal/human"
	"github.com/scaleway/scaleway-sdk-go/api/vpcgw/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
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

func gatewayNetworkCreateBuilder(c *core.Command) *core.Command {
	c.WaitFunc = func(ctx context.Context, _, respI interface{}) (interface{}, error) {
		getResp := respI.(*vpcgw.GatewayNetwork)
		api := vpcgw.NewAPI(core.ExtractClient(ctx))
		return api.WaitForGatewayNetwork(&vpcgw.WaitForGatewayNetworkRequest{
			GatewayNetworkID: getResp.ID,
			Zone:             getResp.Zone,
			Timeout:          scw.TimeDurationPtr(gatewayActionTimeout),
			RetryInterval:    core.DefaultRetryInterval,
		})
	}
	return c
}

func gatewayNetworkMarshallerFunc(i interface{}, opt *human.MarshalOpt) (string, error) {
	type tmp vpcgw.Gateway
	vpcgtwNetwork := tmp(i.(vpcgw.Gateway))
	opt.Sections = []*human.MarshalSection{
		{
			FieldName: "IP",
			Title:     "IP",
		},
		{
			FieldName: "GatewayNetworks",
			Title:     "GatewayNetworks",
		},
	}
	str, err := human.Marshal(vpcgtwNetwork, opt)
	if err != nil {
		return "", err
	}
	return str, nil
}
