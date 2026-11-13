package container

import (
	"context"
	"reflect"
	"time"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/cockpitlogs"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type containerLogsRequest struct {
	ContainerID string
	Region      scw.Region
	TimeSpan    time.Duration
	EntryCount  *int
}

func (r *containerLogsRequest) resource() cockpitlogs.Resource {
	return cockpitlogs.Resource{
		Type:   "serverless_container",
		ID:     r.ContainerID,
		Region: r.Region,
	}
}

func containerLogs() *core.Command {
	return &core.Command{
		Short:     `Show container logs`,
		Long:      `Display the logs of a container from the last 2 hours`,
		Namespace: "container",
		Resource:  "container",
		Verb:      "logs",
		Groups:    []string{"utility"},
		ArgsType:  reflect.TypeFor[containerLogsRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "container-id",
				Short:      "ID of the container which logs are to be displayed",
				Positional: true,
			},
			{
				Name:    "time-span",
				Short:   "Time range for which to retrieve container logs in duration format, defaults to 2h",
				Default: core.DefaultValueSetter("2h"),
			},
			{
				Name:  "entry-count",
				Short: "Maximum number of log entries to be displayed",
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
			),
		},
		Run: func(ctx context.Context, argsI any) (any, error) {
			args := argsI.(*containerLogsRequest)

			return cockpitlogs.Recent(ctx, args.resource(), args.TimeSpan, args.EntryCount)
		},
		FollowFunc: func(ctx context.Context, argsI any, emit func(any) error) error {
			args := argsI.(*containerLogsRequest)

			return cockpitlogs.Follow(ctx, args.resource(), args.TimeSpan, emit)
		},
	}
}
