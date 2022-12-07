package lb

import (
	"context"
	"reflect"
	"strings"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/internal/core"
	"github.com/scaleway/scaleway-cli/v2/internal/human"
	"github.com/scaleway/scaleway-sdk-go/api/lb/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

var (
	lbStatusMarshalSpecs = human.EnumMarshalSpecs{
		lb.LBStatusError:     &human.EnumMarshalSpec{Attribute: color.FgRed, Value: "error"},
		lb.LBStatusLocked:    &human.EnumMarshalSpec{Attribute: color.FgRed, Value: "locked"},
		lb.LBStatusMigrating: &human.EnumMarshalSpec{Attribute: color.FgBlue, Value: "migrating"},
		lb.LBStatusPending:   &human.EnumMarshalSpec{Attribute: color.FgBlue, Value: "pending"},
		lb.LBStatusReady:     &human.EnumMarshalSpec{Attribute: color.FgGreen, Value: "ready"},
		lb.LBStatusStopped:   &human.EnumMarshalSpec{Attribute: color.Faint, Value: "stopped"},
		lb.LBStatusUnknown:   &human.EnumMarshalSpec{Attribute: color.Faint, Value: "unknown"},
	}
)

func lbMarshalerFunc(i interface{}, opt *human.MarshalOpt) (string, error) {
	type tmp lb.LB
	loadbalancer := tmp(i.(lb.LB))

	opt.Sections = []*human.MarshalSection{
		{
			FieldName: "IP",
			Title:     "IPs",
		},
		{
			FieldName: "Instances",
			Title:     "Backends",
		},
	}

	str, err := human.Marshal(loadbalancer, opt)
	if err != nil {
		return "", err
	}

	return str, nil
}

func lbWaitCommand() *core.Command {
	return &core.Command{
		Short:     `Wait for a load balancer to reach a stable state`,
		Long:      `Wait for a load balancer to reach a stable state. This is similar to using --wait flag.`,
		Namespace: "lb",
		Resource:  "lb",
		Verb:      "wait",
		ArgsType:  reflect.TypeOf(lb.ZonedAPIWaitForLBRequest{}),
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, err error) {
			api := lb.NewZonedAPI(core.ExtractClient(ctx))
			args := argsI.(*lb.ZonedAPIWaitForLBRequest)
			return api.WaitForLb(&lb.ZonedAPIWaitForLBRequest{
				LBID:          args.LBID,
				Zone:          args.Zone,
				RetryInterval: core.DefaultRetryInterval,
			})
		},
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "lb-id",
				Short:      `ID of the load balancer you want to wait for.`,
				Required:   true,
				Positional: true,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZonePlWaw1, scw.ZoneNlAms1),
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
	c.ArgSpecs.GetByName("type").ValidateFunc = func(argSpec *core.ArgSpec, value interface{}) error {
		// Allow all lb types
		return nil
	}

	c.WaitFunc = func(ctx context.Context, argsI, respI interface{}) (interface{}, error) {
		api := lb.NewZonedAPI(core.ExtractClient(ctx))
		return api.WaitForLb(&lb.ZonedAPIWaitForLBRequest{
			LBID:          respI.(*lb.LB).ID,
			Zone:          respI.(*lb.LB).Zone,
			RetryInterval: core.DefaultRetryInterval,
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
	c.ArgSpecs.GetByName("type").ValidateFunc = func(argSpec *core.ArgSpec, value interface{}) error {
		// Allow all lb types
		return nil
	}
	c.Interceptor = interceptLB()
	return c
}

func lbGetBuilder(c *core.Command) *core.Command {
	c.Interceptor = interceptLB()
	return c
}

func lbUpdateBuilder(c *core.Command) *core.Command {
	c.Interceptor = interceptLB()
	return c
}

func lbDeleteBuilder(c *core.Command) *core.Command {
	c.Interceptor = interceptLB()
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

func interceptLB() core.CommandInterceptor {
	return func(ctx context.Context, argsI interface{}, runner core.CommandRunner) (interface{}, error) {
		res, err := runner(ctx, argsI)
		if err != nil {
			return nil, err
		}

		lbResp, err := human.Marshal(res.(*lb.LB), nil)
		if err != nil {
			return "", err
		}

		if len(res.(*lb.LB).Tags) != 0 && res.(*lb.LB).Tags[0] == kapsuleTag {
			return strings.Join([]string{
				lbResp,
				warningKapsuleTaggedMessageView(),
			}, "\n\n"), nil
		}

		return res, nil
	}
}
