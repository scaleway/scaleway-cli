package iot

import (
	"context"
	"time"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-sdk-go/api/iot/v1"
)

const (
	hubActionTimeout = 5 * time.Minute
)

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
