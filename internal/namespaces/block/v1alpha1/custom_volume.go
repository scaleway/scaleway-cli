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
	volumeActionTimeout = 5 * time.Minute
)

type volumeWaitRequest struct {
	Zone     scw.Zone
	VolumeID string
	Timeout  time.Duration

	TerminalStatus *block.VolumeStatus
}

func volumeWaitCommand() *core.Command {
	volumeStatuses := block.VolumeStatus("").Values()
	volumeStatusStrings := make([]string, len(volumeStatuses))
	for k, v := range volumeStatuses {
		volumeStatusStrings[k] = v.String()
	}

	return &core.Command{
		Short:     `Wait for volume to reach a stable state`,
		Long:      `Wait for volume to reach a stable state. This is similar to using --wait flag on other action commands, but without requiring a new action on the volume.`,
		Namespace: "block",
		Resource:  "volume",
		Verb:      "wait",
		Groups:    []string{"workflow"},
		ArgsType:  reflect.TypeOf(volumeWaitRequest{}),
		Run: func(ctx context.Context, argsI any) (i any, err error) {
			args := argsI.(*volumeWaitRequest)

			return block.NewAPI(core.ExtractClient(ctx)).WaitForVolume(&block.WaitForVolumeRequest{
				Zone:          args.Zone,
				VolumeID:      args.VolumeID,
				Timeout:       scw.TimeDurationPtr(args.Timeout),
				RetryInterval: core.DefaultRetryInterval,

				TerminalStatus: args.TerminalStatus,
			})
		},
		ArgSpecs: core.ArgSpecs{
			core.WaitTimeoutArgSpec(volumeActionTimeout),
			{
				Name:       "volume-id",
				Short:      `ID of the volume affected by the action.`,
				Required:   true,
				Positional: true,
			},
			{
				Name:       "terminal-status",
				Short:      `Expected terminal status, will wait until this status is reached.`,
				EnumValues: volumeStatusStrings,
			},
			core.ZoneArgSpec((*instance.API)(nil).Zones()...),
		},
		Examples: []*core.Example{
			{
				Short:    "Wait for a volume to be available",
				ArgsJSON: `{"volume_id": "11111111-1111-1111-1111-111111111111", "terminal_status": "available"}`,
			},
		},
	}
}

func blockVolumeCreateBuilder(c *core.Command) *core.Command {
	c.WaitFunc = func(ctx context.Context, _, respI any) (any, error) {
		resp := respI.(*block.Volume)

		return block.NewAPI(core.ExtractClient(ctx)).WaitForVolume(&block.WaitForVolumeRequest{
			VolumeID: resp.ID,
			Zone:     resp.Zone,
		})
	}

	return c
}
