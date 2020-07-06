package rdb

import (
	"context"
	"time"

	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-sdk-go/api/rdb/v1"
)

func engineListBuilder(c *core.Command) *core.Command {
	type customList struct {
		Name       string
		EngineType string
		EndOfLife  time.Time
	}

	c.View = &core.View{
		Fields: []*core.ViewField{
			{
				Label:     "Name",
				FieldName: "Name",
			},
			{
				Label:     "Engine Type",
				FieldName: "EngineType",
			},
			{
				Label:     "End of Life",
				FieldName: "EndOfLife",
			},
		},
	}

	c.AddInterceptors(func(ctx context.Context, argsI interface{}, runner core.CommandRunner) (i interface{}, err error) {
		listEngineResp, err := runner(ctx, argsI)
		if err != nil {
			return listEngineResp, err
		}
		engineList := listEngineResp.([]*rdb.DatabaseEngine)
		var res []customList
		for _, engine := range engineList {
			for _, version := range engine.Versions {
				res = append(res, customList{
					Name:       version.Name,
					EngineType: engine.Name,
					EndOfLife:  version.EndOfLife,
				})
			}
		}

		return res, nil
	})

	return c
}
