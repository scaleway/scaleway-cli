package instance

import (
	"context"
	"fmt"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/internal/core"
	"github.com/scaleway/scaleway-cli/v2/internal/human"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
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

func privateNicCreateBuilder(c *core.Command) *core.Command {
	c.Interceptor = func(ctx context.Context, argsI interface{}, runner core.CommandRunner) (i interface{}, err error) {
		res, err := runner(ctx, argsI)
		if err == nil {
			return res, nil
		}

		resErr, isResErr := err.(*scw.InvalidArgumentsError)
		if !isResErr {
			return nil, err
		}

		if resErr.Details[0].ArgumentName == "private_network_id" && resErr.Details[0].HelpMessage == "required key not provided" {
			return "", &core.CliError{
				Err: fmt.Errorf("missing required argument 'private-network-id'"),
			}
		}

		return nil, resErr
	}

	return c
}
