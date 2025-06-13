package iot

import (
	"context"
	"time"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	"github.com/scaleway/scaleway-sdk-go/api/iot/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

const (
	hubActionTimeout = 5 * time.Minute
)

var hubStatusMarshalSpecs = human.EnumMarshalSpecs{
	iot.HubStatusDisabled:  &human.EnumMarshalSpec{Attribute: color.FgRed, Value: "disabled"},
	iot.HubStatusDisabling: &human.EnumMarshalSpec{Attribute: color.FgBlue, Value: "disabling"},
	iot.HubStatusEnabling:  &human.EnumMarshalSpec{Attribute: color.FgBlue, Value: "enabling"},
	iot.HubStatusError:     &human.EnumMarshalSpec{Attribute: color.FgRed, Value: "error"},
	iot.HubStatusReady:     &human.EnumMarshalSpec{Attribute: color.FgGreen, Value: "ready"},
}

func hubCreateBuilder(c *core.Command) *core.Command {
	c.WaitFunc = func(ctx context.Context, _, respI any) (any, error) {
		api := iot.NewAPI(core.ExtractClient(ctx))

		return api.WaitForHub(&iot.WaitForHubRequest{
			HubID:         respI.(*iot.Hub).ID,
			Region:        respI.(*iot.Hub).Region,
			Timeout:       scw.TimeDurationPtr(hubActionTimeout),
			RetryInterval: core.DefaultRetryInterval,
		})
	}

	return c
}
