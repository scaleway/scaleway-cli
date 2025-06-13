package vpcgw

import (
	"context"
	"time"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	"github.com/scaleway/scaleway-sdk-go/api/vpcgw/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

const (
	gatewayActionTimeout = 60 * time.Minute
)

var gatewayStatusMarshalSpecs = human.EnumMarshalSpecs{
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

func gatewayCreateBuilder(c *core.Command) *core.Command {
	c.WaitFunc = func(ctx context.Context, _, respI any) (any, error) {
		getResp := respI.(*vpcgw.Gateway)
		api := vpcgw.NewAPI(core.ExtractClient(ctx))

		return api.WaitForGateway(&vpcgw.WaitForGatewayRequest{
			GatewayID:     getResp.ID,
			Zone:          getResp.Zone,
			Timeout:       scw.TimeDurationPtr(gatewayActionTimeout),
			RetryInterval: core.DefaultRetryInterval,
		})
	}

	return c
}

func gatewayMarshalerFunc(i any, opt *human.MarshalOpt) (string, error) {
	type tmp vpcgw.Gateway
	vpcgtw := tmp(i.(vpcgw.Gateway))
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
	str, err := human.Marshal(vpcgtw, opt)
	if err != nil {
		return "", err
	}

	return str, nil
}
