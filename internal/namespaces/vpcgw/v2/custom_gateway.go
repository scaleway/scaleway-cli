package vpcgw

import (
	"context"
	"time"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	"github.com/scaleway/scaleway-sdk-go/api/vpcgw/v2"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

const (
	gatewayActionTimeout = 60 * time.Minute
)

var gatewayStatusMarshalSpecs = human.EnumMarshalSpecs{
	vpcgw.GatewayStatusUnknownStatus: &human.EnumMarshalSpec{Attribute: color.Faint},
	vpcgw.GatewayStatusAllocating:    &human.EnumMarshalSpec{Attribute: color.FgBlue},
	vpcgw.GatewayStatusConfiguring:   &human.EnumMarshalSpec{Attribute: color.FgBlue},
	vpcgw.GatewayStatusDeleting:      &human.EnumMarshalSpec{Attribute: color.FgBlue},
	vpcgw.GatewayStatusFailed:        &human.EnumMarshalSpec{Attribute: color.FgRed},
	vpcgw.GatewayStatusRunning:       &human.EnumMarshalSpec{Attribute: color.FgGreen},
	vpcgw.GatewayStatusStopped:       &human.EnumMarshalSpec{Attribute: color.FgRed},
	vpcgw.GatewayStatusLocked:        &human.EnumMarshalSpec{Attribute: color.FgRed},
	vpcgw.GatewayStatusStopping:      &human.EnumMarshalSpec{Attribute: color.FgBlue},
}

func gatewayCreateBuilder(c *core.Command) *core.Command {
	c.WaitFunc = func(ctx context.Context, _, respI interface{}) (interface{}, error) {
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

func gatewayMarshalerFunc(i interface{}, opt *human.MarshalOpt) (string, error) {
	type tmp vpcgw.Gateway
	vpcgtw := tmp(i.(vpcgw.Gateway))
	opt.Sections = []*human.MarshalSection{
		{
			FieldName: "IPv4",
			Title:     "IPv4",
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
