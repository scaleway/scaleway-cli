package instance

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v2alpha1"
)

const (
	snapshotActionTimeout = 60 * time.Minute
	snapshotCreateTimeout = 10 * time.Minute
)

func snapshotCreateBuilder(c *core.Command) *core.Command {
	type customCreateSnapshotRequest struct {
		*instance.VolumeAPICreateSnapshotRequest
		VolumeID string
	}

	c.ArgSpecs.GetByName("name").Default = core.RandomValueGenerator("snap")

	c.ArgSpecs.AddBefore("public", &core.ArgSpec{
		Name:       "volume-id",
		Short:      "ID of the base volume, please use base-volume-id instead",
		Required:   false,
		Deprecated: true,
	})

	c.ArgsType = reflect.TypeFor[customCreateSnapshotRequest]()

	c.Interceptor = func(ctx context.Context, argsI any, runner core.CommandRunner) (i any, err error) {
		args := argsI.(*customCreateSnapshotRequest)

		if args.VolumeAPICreateSnapshotRequest == nil {
			args.VolumeAPICreateSnapshotRequest = &instance.VolumeAPICreateSnapshotRequest{}
		}

		// If no 'base-volume-id' arg was passed, try deprecated 'volume-id' arg instead or return an error.
		request := args.VolumeAPICreateSnapshotRequest
		if request.BaseVolumeID == "" {
			request.BaseVolumeID = args.VolumeID
			if request.BaseVolumeID == "" {
				return nil, errors.New("invalid arguments: 'base-volume-id' is required")
			}
		}

		return runner(ctx, request)
	}

	c.WaitFunc = func(ctx context.Context, argsI, respI any) (any, error) {
		api := instance.NewVolumeAPI(core.ExtractClient(ctx))

		return api.WaitForSnapshot(&instance.WaitForSnapshotRequest{
			SnapshotID:    respI.(*instance.Snapshot).ID,
			Zone:          argsI.(*customCreateSnapshotRequest).Zone,
			Timeout:       new(snapshotCreateTimeout),
			RetryInterval: core.DefaultRetryInterval,
		})
	}

	return c
}

func snapshotGetBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("snapshot-id").Positional = true

	return c
}

func snapshotDeleteBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("snapshot-id").Positional = true

	return c
}

func snapshotListBuilder(c *core.Command) *core.Command {
	c.Interceptor = func(ctx context.Context, argsI any, runner core.CommandRunner) (i any, err error) {
		rawResp, err := runner(ctx, argsI)
		if err != nil {
			return rawResp, err
		}

		snapList, ok := rawResp.(*instance.ListSnapshotsResponse)
		if !ok {
			return rawResp, fmt.Errorf(
				"expected response of type instance.ListSnapshotsResponse, got %T",
				rawResp,
			)
		}

		return snapList.Snapshots, nil
	}

	c.View = &core.View{
		Fields: []*core.ViewField{
			{
				Label:     "ID",
				FieldName: "ID",
			},
			{
				Label:     "NAME",
				FieldName: "Name",
			},
			{
				Label:     "SIZE",
				FieldName: "Size",
			},
			{
				Label:     "BASE VOLUME ID",
				FieldName: "BaseVolumeID",
			},
			{
				Label:     "VOLUME TYPE",
				FieldName: "VolumeType",
			},
			{
				Label:     "PUBLIC",
				FieldName: "Public",
			},
			{
				Label:     "STATUS",
				FieldName: "Status",
			},
			{
				Label:     "TAGS",
				FieldName: "Tags",
			},
			{
				Label:     "ZONE",
				FieldName: "Zone",
			},
			{
				Label:     "PROJECT ID",
				FieldName: "ProjectID",
			},
			{
				Label:     "CREATED AT",
				FieldName: "CreatedAt",
			},
			{
				Label:     "UPDATED AT",
				FieldName: "UpdatedAt",
			},
			{
				Label:     "SRN",
				FieldName: "Srn",
			},
		},
	}

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
		ArgsType:  reflect.TypeFor[instance.WaitForSnapshotRequest](),
		Run: func(ctx context.Context, argsI any) (i any, err error) {
			api := instance.NewVolumeAPI(core.ExtractClient(ctx))

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
	c.ArgSpecs.GetByName("snapshot-id").Positional = true

	c.WaitFunc = func(ctx context.Context, _, respI any) (any, error) {
		snapshot := respI.(*instance.Snapshot)
		api := instance.NewVolumeAPI(core.ExtractClient(ctx))

		return api.WaitForSnapshot(&instance.WaitForSnapshotRequest{
			SnapshotID:    snapshot.ID,
			Zone:          snapshot.Zone,
			Timeout:       new(snapshotActionTimeout),
			RetryInterval: core.DefaultRetryInterval,
		})
	}

	return c
}

func snapshotExportToObjectStorageBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("snapshot-id").Positional = true

	c.WaitFunc = func(ctx context.Context, _, respI any) (any, error) {
		snapshot := respI.(*instance.Snapshot)
		api := instance.NewVolumeAPI(core.ExtractClient(ctx))

		return api.WaitForSnapshot(&instance.WaitForSnapshotRequest{
			SnapshotID:    snapshot.ID,
			Zone:          snapshot.Zone,
			Timeout:       new(snapshotActionTimeout),
			RetryInterval: core.DefaultRetryInterval,
		})
	}

	return c
}

func snapshotImportToObjectStorageBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("volume-type").Default = core.DefaultValueSetter(
		instance.SnapshotVolumeTypeLSSD.String(),
	)

	return c
}
