package instance

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
)

// Builders

func ipCreateBuilder(c *core.Command) *core.Command {
	type customCreateIPRequest struct {
		*instance.CreateIPRequest
		ProjectID *string
	}

	renameProjectIDArgSpec(c.ArgSpecs)

	c.ArgsType = reflect.TypeOf(customCreateIPRequest{})

	c.AddInterceptors(func(ctx context.Context, argsI interface{}, runner core.CommandRunner) (i interface{}, err error) {
		args := argsI.(*customCreateIPRequest)

		if args.CreateIPRequest == nil {
			args.CreateIPRequest = &instance.CreateIPRequest{}
		}
		request := args.CreateIPRequest
		request.Project = args.ProjectID

		return runner(ctx, request)
	})

	return c
}

func ipListBuilder(c *core.Command) *core.Command {
	type customListIPsRequest struct {
		*instance.ListIPsRequest
		ProjectID *string
	}

	renameProjectIDArgSpec(c.ArgSpecs)

	c.ArgsType = reflect.TypeOf(customListIPsRequest{})

	c.AddInterceptors(func(ctx context.Context, argsI interface{}, runner core.CommandRunner) (i interface{}, err error) {
		args := argsI.(*customListIPsRequest)

		if args.ListIPsRequest == nil {
			args.ListIPsRequest = &instance.ListIPsRequest{}
		}
		request := args.ListIPsRequest
		request.Project = args.ProjectID

		return runner(ctx, request)
	})
	return c
}
