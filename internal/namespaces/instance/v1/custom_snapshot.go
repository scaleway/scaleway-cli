package instance

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
)

// Builders

func snapshotCreateBuilder(c *core.Command) *core.Command {
	type customCreateSnapshotRequest struct {
		*instance.CreateSnapshotRequest
		OrganizationID string
	}

	renameOrganizationIDArgSpec(c.ArgSpecs)

	c.ArgsType = reflect.TypeOf(customCreateSnapshotRequest{})

	c.AddInterceptors(func(ctx context.Context, argsI interface{}, runner core.CommandRunner) (i interface{}, err error) {
		args := argsI.(*customCreateSnapshotRequest)

		if args.CreateSnapshotRequest == nil {
			args.CreateSnapshotRequest = &instance.CreateSnapshotRequest{}
		}

		request := args.CreateSnapshotRequest
		request.Organization = args.OrganizationID

		return runner(ctx, request)
	})
	return c
}

func snapshotListBuilder(c *core.Command) *core.Command {
	type customListSnapshotsRequest struct {
		*instance.ListSnapshotsRequest
		OrganizationID *string
	}

	renameOrganizationIDArgSpec(c.ArgSpecs)

	c.ArgsType = reflect.TypeOf(customListSnapshotsRequest{})

	c.AddInterceptors(func(ctx context.Context, argsI interface{}, runner core.CommandRunner) (i interface{}, err error) {
		args := argsI.(*customListSnapshotsRequest)

		if args.ListSnapshotsRequest == nil {
			args.ListSnapshotsRequest = &instance.ListSnapshotsRequest{}
		}

		request := args.ListSnapshotsRequest
		request.Organization = args.OrganizationID

		return runner(ctx, request)
	})
	return c
}
