package lb

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-sdk-go/api/lb/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func lbWaitCommand() *core.Command {
	return &core.Command{
		Short:     `Wait for a load balancer to reach a stable state`,
		Long:      `Wait for a load balancer to reach a stable state. This is similar to using --wait flag.`,
		Namespace: "lb",
		Resource:  "lb",
		Verb:      "wait",
		ArgsType:  reflect.TypeOf(lb.WaitForLBRequest{}),
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, err error) {
			api := lb.NewAPI(core.ExtractClient(ctx))
			args := argsI.(*lb.WaitForLBRequest)
			return api.WaitForLb(&lb.WaitForLBRequest{
				LBID:   args.LBID,
				Region: args.Region,
			})
		},
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "lb-id",
				Short:      `ID of the load balancer you want to wait for.`,
				Required:   true,
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms),
		},
		Examples: []*core.Example{
			{
				Short:    "Wait for a load balancer to reach a stable state",
				ArgsJSON: `{"lb_id": "11111111-1111-1111-1111-111111111111"}`,
			},
		},
	}
}

func lbCreateBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("type").EnumValues = typesList
	c.ArgSpecs.GetByName("type").Default = core.DefaultValueSetter("LB-S")

	c.WaitFunc = func(ctx context.Context, argsI, respI interface{}) (interface{}, error) {
		api := lb.NewAPI(core.ExtractClient(ctx))
		return api.WaitForLb(&lb.WaitForLBRequest{
			LBID:   respI.(*lb.LB).ID,
			Region: respI.(*lb.LB).Region,
		})
	}
	return c
}

var typesList = []string{
	"LB-S",
	"LB-GP-M",
	"LB-GP-L",
}

func lbMigrateBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("type").EnumValues = typesList

	return c
}

func lbGetBuilder(c *core.Command) *core.Command {
	c.View = &core.View{
		Sections: []*core.ViewSection{

			{
				FieldName: "IP",
				Title:     "IPs",
			},
			{
				FieldName: "Instances",
			},
		},
	}

	return c
}

func lbGetStatsBuilder(c *core.Command) *core.Command {
	c.View = &core.View{
		Sections: []*core.ViewSection{

			{
				FieldName: "BackendServersStats",
				Title:     "Backends Statistics",
			},
		},
	}

	return c
}
