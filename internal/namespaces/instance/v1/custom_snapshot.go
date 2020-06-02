package instance

import (
	"context"
	"reflect"
	"time"

	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
)

const (
	snapshotActionTimeout = 60 * time.Minute
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

	c.WaitFunc = func(ctx context.Context, argsI, respI interface{}) (interface{}, error) {
		api := instance.NewAPI(core.ExtractClient(ctx))
		return api.WaitForSnapshot(&instance.WaitForSnapshotRequest{
			SnapshotID:    respI.(*instance.CreateSnapshotResponse).Snapshot.ID,
			Zone:          argsI.(*customCreateSnapshotRequest).Zone,
			Timeout:       serverActionTimeout,
			RetryInterval: core.DefaultRetryInterval,
		})
	}

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

func snapshotWaitCommand() *core.Command {
	return &core.Command{
		Short:     `Wait for snapshot to reach a stable state`,
		Long:      `Wait for snapshot to reach a stable state. This is similar to using --wait flag on other action commands, but without requiring a new action on the snapshot.`,
		Namespace: "instance",
		Resource:  "snapshot",
		Verb:      "wait",
		ArgsType:  reflect.TypeOf(instance.WaitForSnapshotRequest{}),
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, err error) {
			api := instance.NewAPI(core.ExtractClient(ctx))
			return api.WaitForSnapshot(&instance.WaitForSnapshotRequest{
				Zone:          argsI.(*instance.WaitForSnapshotRequest).Zone,
				SnapshotID:    argsI.(*instance.WaitForSnapshotRequest).SnapshotID,
				Timeout:       snapshotActionTimeout,
				RetryInterval: core.DefaultRetryInterval,
			})

		},
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "snapshot-id",
				Short:      `ID of the snapshot.`,
				Required:   true,
				Positional: true,
			},
			core.ZoneArgSpec(),
		},
		Examples: []*core.Example{
			{
				Short:   "Wait for a snapshot to reach a stable state",
				Request: `{"snapshot_id": "11111111-1111-1111-1111-111111111111"}`,
			},
		},
	}
}
