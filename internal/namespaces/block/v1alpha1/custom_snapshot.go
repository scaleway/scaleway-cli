package block

import (
	"context"
	"reflect"
	"time"

	"github.com/scaleway/scaleway-cli/v2/core"
	block "github.com/scaleway/scaleway-sdk-go/api/block/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

const (
	snapshotActionTimeout = 5 * time.Minute
)

type snapshotWaitRequest struct {
	Zone       scw.Zone
	SnapshotID string
	Timeout    time.Duration

	TerminalStatus *block.SnapshotStatus
}

func snapshotWaitCommand() *core.Command {
	snapshotsStatuses := block.SnapshotStatus("").Values()
	snapshotsStatusStrings := make([]string, len(snapshotsStatuses))
	for k, v := range snapshotsStatuses {
		snapshotsStatusStrings[k] = v.String()
	}

	return &core.Command{
		Short:     `Wait for snapshot to reach a stable state`,
		Long:      `Wait for snapshot to reach a stable state. This is similar to using --wait flag on other action commands, but without requiring a new action on the snapshot.`,
		Namespace: "block",
		Resource:  "snapshot",
		Verb:      "wait",
		Groups:    []string{"workflow"},
		ArgsType:  reflect.TypeOf(snapshotWaitRequest{}),
		Run: func(ctx context.Context, argsI any) (i any, err error) {
			args := argsI.(*snapshotWaitRequest)

			return block.NewAPI(core.ExtractClient(ctx)).
				WaitForSnapshot(&block.WaitForSnapshotRequest{
					Zone:          args.Zone,
					SnapshotID:    args.SnapshotID,
					Timeout:       scw.TimeDurationPtr(args.Timeout),
					RetryInterval: core.DefaultRetryInterval,

					TerminalStatus: args.TerminalStatus,
				})
		},
		ArgSpecs: core.ArgSpecs{
			core.WaitTimeoutArgSpec(snapshotActionTimeout),
			{
				Name:       "snapshot-id",
				Short:      `ID of the snapshot affected by the action.`,
				Required:   true,
				Positional: true,
			},
			{
				Name:       "terminal-status",
				Short:      `Expected terminal status, will wait until this status is reached.`,
				EnumValues: snapshotsStatusStrings,
			},
			core.ZoneArgSpec((*instance.API)(nil).Zones()...),
		},
		Examples: []*core.Example{
			{
				Short:    "Wait for a snapshot to be available",
				ArgsJSON: `{"snapshot_id": "11111111-1111-1111-1111-111111111111", "terminal_status": "available"}`,
			},
		},
	}
}

func blockSnapshotCreateBuilder(c *core.Command) *core.Command {
	c.WaitFunc = func(ctx context.Context, _, respI any) (any, error) {
		resp := respI.(*block.Snapshot)

		return block.NewAPI(core.ExtractClient(ctx)).WaitForSnapshot(&block.WaitForSnapshotRequest{
			SnapshotID:    resp.ID,
			Zone:          resp.Zone,
			RetryInterval: core.DefaultRetryInterval,
		})
	}

	return c
}
