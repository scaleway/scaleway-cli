package iot

import (
	"context"
	"time"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	humanv2 "github.com/scaleway/scaleway-cli/v2/core/human/v2"
	"github.com/scaleway/scaleway-sdk-go/api/iot/v1"
)

const (
	hubActionTimeout = 5 * time.Minute
)

var hubStatusFormat = humanv2.EnumFormats[iot.HubStatus]{
	iot.HubStatusDisabled:  humanv2.EnumFormat{Color: color.FgRed, Value: "disabled"},
	iot.HubStatusDisabling: humanv2.EnumFormat{Color: color.FgBlue, Value: "disabling"},
	iot.HubStatusEnabling:  humanv2.EnumFormat{Color: color.FgBlue, Value: "enabling"},
	iot.HubStatusError:     humanv2.EnumFormat{Color: color.FgRed, Value: "error"},
	iot.HubStatusReady:     humanv2.EnumFormat{Color: color.FgGreen, Value: "ready"},
}

func hubCreateBuilder(c *core.Command) *core.Command {
	c.WaitFunc = func(ctx context.Context, _, respI any) (any, error) {
		api := iot.NewAPI(core.ExtractClient(ctx))

		return api.WaitForHub(&iot.WaitForHubRequest{
			HubID:         respI.(*iot.Hub).ID,
			Region:        respI.(*iot.Hub).Region,
			Timeout:       new(hubActionTimeout),
			RetryInterval: core.DefaultRetryInterval,
		})
	}

	return c
}
