package rdb

import (
	"context"
	"reflect"
	"time"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-cli/internal/human"
	"github.com/scaleway/scaleway-sdk-go/api/rdb/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

const (
	logsTimeout = 20 * time.Minute
)

var (
	logStatusMarshalSpecs = human.EnumMarshalSpecs{
		rdb.InstanceLogStatusUnknown:  &human.EnumMarshalSpec{Attribute: color.Faint, Value: "unknown"},
		rdb.InstanceLogStatusReady:    &human.EnumMarshalSpec{Attribute: color.FgGreen, Value: "ready"},
		rdb.InstanceLogStatusCreating: &human.EnumMarshalSpec{Attribute: color.FgBlue, Value: "creating"},
		rdb.InstanceLogStatusError:    &human.EnumMarshalSpec{Attribute: color.FgRed, Value: "error"},
	}
)

type logWaitRequest struct {
	InstanceLogID string
	Region        scw.Region
}

func logWaitCommand() *core.Command {
	timeout := logsTimeout
	return &core.Command{
		Short:     `Wait for an instance logs to reach a stable state`,
		Long:      `Wait for an instance logs to reach a stable state. This is similar to using --wait flag.`,
		Namespace: "rdb",
		Resource:  "log",
		Verb:      "wait",
		ArgsType:  reflect.TypeOf(logWaitRequest{}),
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, err error) {
			api := rdb.NewAPI(core.ExtractClient(ctx))
			return api.WaitForInstanceLog(&rdb.WaitForInstanceLogRequest{
				InstanceLogID: argsI.(*logWaitRequest).InstanceLogID,
				Region:        argsI.(*logWaitRequest).Region,
				Timeout:       &timeout,
				RetryInterval: core.DefaultRetryInterval,
			})
		},
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "instance-log-id",
				Short:      `ID of the instance logs you want to wait for.`,
				Required:   true,
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms),
		},
		Examples: []*core.Example{
			{
				Short:    "Wait for an instance logs to reach a stable state",
				ArgsJSON: `{"instance_log_id": "11111111-1111-1111-1111-111111111111"}`,
			},
		},
	}
}
