package dedibox

import (
	"context"
	"github.com/scaleway/scaleway-cli/v2/internal/core"
	"github.com/scaleway/scaleway-sdk-go/api/dedibox/v1"
)

func serviceCreateBuilder(c *core.Command) *core.Command {
	c.WaitFunc = func(ctx context.Context, argsI, respI interface{}) (interface{}, error) {
		api := dedibox.NewAPI(core.ExtractClient(ctx))

	}
}
