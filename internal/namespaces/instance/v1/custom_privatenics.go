package instance

import (
	"context"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-cli/internal/human"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
)

var (
	privateNICStateMarshalSpecs = human.EnumMarshalSpecs{
		instance.PrivateNICStateAvailable:    &human.EnumMarshalSpec{Attribute: color.FgGreen},
		instance.PrivateNICStateSyncing:      &human.EnumMarshalSpec{Attribute: color.FgBlue},
		instance.PrivateNICStateSyncingError: &human.EnumMarshalSpec{Attribute: color.FgRed},
	}
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
