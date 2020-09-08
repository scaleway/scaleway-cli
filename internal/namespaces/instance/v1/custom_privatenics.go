package instance

import (
	"context"

	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
)

func privateNicListBuilder(c *core.Command) *core.Command {
	c.AddInterceptors(func(ctx context.Context, argsI interface{}, runner core.CommandRunner) (i interface{}, err error) {
		listPrivateNicResp, err := runner(ctx, argsI)
		if err != nil {
			return listPrivateNicResp, err
		}
		l := listPrivateNicResp.(*instance.ListPrivateNICsResponse)
		privateNic := l.PrivateNics

		return privateNic, nil
	})

	return c
}
