package rdb

import (
	"context"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	"github.com/scaleway/scaleway-sdk-go/api/rdb/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

var logStatusMarshalSpecs = human.EnumMarshalSpecs{
	rdb.InstanceLogStatusUnknown: &human.EnumMarshalSpec{Attribute: color.Faint, Value: "unknown"},
	rdb.InstanceLogStatusReady:   &human.EnumMarshalSpec{Attribute: color.FgGreen, Value: "ready"},
	rdb.InstanceLogStatusCreating: &human.EnumMarshalSpec{
		Attribute: color.FgBlue,
		Value:     "creating",
	},
	rdb.InstanceLogStatusError: &human.EnumMarshalSpec{Attribute: color.FgRed, Value: "error"},
}

func logPrepareBuilder(c *core.Command) *core.Command {
	c.WaitFunc = func(ctx context.Context, _, respI any) (any, error) {
		getResp := respI.(*rdb.PrepareInstanceLogsResponse)
		api := rdb.NewAPI(core.ExtractClient(ctx))
		readyLogs := make([]*rdb.InstanceLog, len(getResp.InstanceLogs))
		for i := range getResp.InstanceLogs {
			logs, err := api.WaitForInstanceLog(&rdb.WaitForInstanceLogRequest{
				InstanceLogID: getResp.InstanceLogs[i].ID,
				Region:        getResp.InstanceLogs[i].Region,
				Timeout:       scw.TimeDurationPtr(instanceActionTimeout),
				RetryInterval: core.DefaultRetryInterval,
			})
			if err != nil {
				return nil, err
			}
			readyLogs[i] = logs
		}
		respI.(*rdb.PrepareInstanceLogsResponse).InstanceLogs = readyLogs

		return respI.(*rdb.PrepareInstanceLogsResponse), nil
	}

	return c
}
