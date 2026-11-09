package vpcgw

import (
	"context"
	"time"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	humanv2 "github.com/scaleway/scaleway-cli/v2/core/human/v2"
	"github.com/scaleway/scaleway-sdk-go/api/vpcgw/v2"
)

const (
	gatewayActionTimeout = 60 * time.Minute
)

var gatewayStatusFormat = humanv2.EnumFormats[vpcgw.GatewayStatus]{
	vpcgw.GatewayStatusUnknownStatus: humanv2.EnumFormat{Color: color.Faint},
	vpcgw.GatewayStatusAllocating:    humanv2.EnumFormat{Color: color.FgBlue},
	vpcgw.GatewayStatusConfiguring:   humanv2.EnumFormat{Color: color.FgBlue},
	vpcgw.GatewayStatusDeleting:      humanv2.EnumFormat{Color: color.FgBlue},
	vpcgw.GatewayStatusFailed:        humanv2.EnumFormat{Color: color.FgRed},
	vpcgw.GatewayStatusRunning:       humanv2.EnumFormat{Color: color.FgGreen},
	vpcgw.GatewayStatusStopped:       humanv2.EnumFormat{Color: color.FgRed},
	vpcgw.GatewayStatusLocked:        humanv2.EnumFormat{Color: color.FgRed},
	vpcgw.GatewayStatusStopping:      humanv2.EnumFormat{Color: color.FgBlue},
}

func gatewayCreateBuilder(c *core.Command) *core.Command {
	c.WaitFunc = func(ctx context.Context, _, respI any) (any, error) {
		getResp := respI.(*vpcgw.Gateway)
		api := vpcgw.NewAPI(core.ExtractClient(ctx))

		return api.WaitForGateway(&vpcgw.WaitForGatewayRequest{
			GatewayID:     getResp.ID,
			Zone:          getResp.Zone,
			Timeout:       new(gatewayActionTimeout),
			RetryInterval: core.DefaultRetryInterval,
		})
	}

	return c
}

var gatewayLayout = humanv2.Layout{
	Sections: []humanv2.Section{
		{Path: "IPv4", Title: "IPv4"},
		{Path: "GatewayNetworks", Title: "GatewayNetworks"},
	},
}
