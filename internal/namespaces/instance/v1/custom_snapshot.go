package instance

import (
	"context"
	"reflect"
	"time"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

const (
	snapshotActionTimeout = 60 * time.Minute
)

// Builders

func snapshotCreateBuilder(c *core.Command) *core.Command {
	type customCreateSnapshotRequest struct {
		*instance.CreateSnapshotRequest
		OrganizationID *string
		ProjectID      *string
		Unified        bool
	}

	renameOrganizationIDArgSpec(c.ArgSpecs)
	renameProjectIDArgSpec(c.ArgSpecs)
	c.ArgSpecs.DeleteByName("volume-type")
	c.ArgSpecs.AddBefore("tags.{index}", &core.ArgSpec{
		Name:  "unified",
		Short: "Whether a snapshot is unified or not.",
	})

	c.ArgsType = reflect.TypeOf(customCreateSnapshotRequest{})

	c.AddInterceptors(
		func(ctx context.Context, argsI interface{}, runner core.CommandRunner) (i interface{}, err error) {
			args := argsI.(*customCreateSnapshotRequest)

			if args.CreateSnapshotRequest == nil {
				args.CreateSnapshotRequest = &instance.CreateSnapshotRequest{}
			}

			request := args.CreateSnapshotRequest
			request.Organization = args.OrganizationID
			request.Project = args.ProjectID

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			if args.Unified {
				request.VolumeType = instance.SnapshotVolumeTypeUnified
			} else if request.VolumeID != nil {
				// If the snapshot is not unified, we need to set the snapshot volume type to the same type as the volume we target.
				// Done only when creating snapshot from volume
				volume, err := api.GetVolume(&instance.GetVolumeRequest{
					VolumeID: *args.VolumeID,
					Zone:     args.Zone,
				})
				if err != nil {
					return nil, err
				}

				request.VolumeType = instance.SnapshotVolumeType(volume.Volume.VolumeType)
			}

			return runner(ctx, request)
		},
	)

	c.WaitFunc = func(ctx context.Context, argsI, respI interface{}) (interface{}, error) {
		api := instance.NewAPI(core.ExtractClient(ctx))

		return api.WaitForSnapshot(&instance.WaitForSnapshotRequest{
			SnapshotID:    respI.(*instance.CreateSnapshotResponse).Snapshot.ID,
			Zone:          argsI.(*customCreateSnapshotRequest).Zone,
			Timeout:       scw.TimeDurationPtr(serverActionTimeout),
			RetryInterval: core.DefaultRetryInterval,
		})
	}

	return c
}

func snapshotListBuilder(c *core.Command) *core.Command {
	type customListSnapshotsRequest struct {
		*instance.ListSnapshotsRequest
		OrganizationID *string
		ProjectID      *string
	}

	renameOrganizationIDArgSpec(c.ArgSpecs)
	renameProjectIDArgSpec(c.ArgSpecs)

	c.ArgsType = reflect.TypeOf(customListSnapshotsRequest{})

	c.AddInterceptors(
		func(ctx context.Context, argsI interface{}, runner core.CommandRunner) (i interface{}, err error) {
			args := argsI.(*customListSnapshotsRequest)

			if args.ListSnapshotsRequest == nil {
				args.ListSnapshotsRequest = &instance.ListSnapshotsRequest{}
			}

			request := args.ListSnapshotsRequest
			request.Organization = args.OrganizationID
			request.Project = args.ProjectID

			return runner(ctx, request)
		},
	)

	return c
}

func snapshotWaitCommand() *core.Command {
	return &core.Command{
		Short:     `Wait for snapshot to reach a stable state`,
		Long:      `Wait for snapshot to reach a stable state. This is similar to using --wait flag on other action commands, but without requiring a new action on the snapshot.`,
		Namespace: "instance",
		Resource:  "snapshot",
		Verb:      "wait",
		Groups:    []string{"workflow"},
		ArgsType:  reflect.TypeOf(instance.WaitForSnapshotRequest{}),
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, err error) {
			api := instance.NewAPI(core.ExtractClient(ctx))

			return api.WaitForSnapshot(&instance.WaitForSnapshotRequest{
				Zone:          argsI.(*instance.WaitForSnapshotRequest).Zone,
				SnapshotID:    argsI.(*instance.WaitForSnapshotRequest).SnapshotID,
				Timeout:       argsI.(*instance.WaitForSnapshotRequest).Timeout,
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
			core.ZoneArgSpec((*instance.API)(nil).Zones()...),
			core.WaitTimeoutArgSpec(snapshotActionTimeout),
		},
		Examples: []*core.Example{
			{
				Short:    "Wait for a snapshot to reach a stable state",
				ArgsJSON: `{"snapshot_id": "11111111-1111-1111-1111-111111111111"}`,
			},
		},
	}
}

func snapshotUpdateBuilder(c *core.Command) *core.Command {
	c.WaitFunc = func(ctx context.Context, _, respI interface{}) (interface{}, error) {
		snapshot := respI.(*instance.UpdateSnapshotResponse).Snapshot
		api := instance.NewAPI(core.ExtractClient(ctx))

		return api.WaitForSnapshot(&instance.WaitForSnapshotRequest{
			SnapshotID:    snapshot.ID,
			Zone:          snapshot.Zone,
			Timeout:       scw.TimeDurationPtr(snapshotActionTimeout),
			RetryInterval: core.DefaultRetryInterval,
		})
	}

	return c
}

func snapshotPlanMigrationCommand() *core.Command {
	cmd := instanceVolumePlanMigration()
	cmd.Resource = "snapshot"

	cmd.ArgSpecs.DeleteByName("volume-id")
	cmd.ArgSpecs.GetByName("snapshot-id").Positional = true

	return cmd
}

func snapshotApplyMigrationCommand() *core.Command {
	cmd := instanceVolumeApplyMigration()
	cmd.Resource = "snapshot"

	cmd.ArgSpecs.DeleteByName("volume-id")
	cmd.ArgSpecs.GetByName("snapshot-id").Positional = true

	return cmd
}
