package lb

import (
	"context"
	"errors"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	"github.com/scaleway/scaleway-sdk-go/api/lb/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

var (
	defaultLBTimeout     = 10 * time.Minute
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

func lbMarshalerFunc(i any, opt *human.MarshalOpt) (string, error) {
	type tmp lb.LB
	loadbalancer := tmp(i.(lb.LB))

	opt.Sections = []*human.MarshalSection{
		{
			FieldName: "IP",
			Title:     "IPs",
		},
		{
			FieldName: "Instances",
			Title:     "LB Instances",
		},
	}

	if len(loadbalancer.Tags) != 0 && loadbalancer.Tags[0] == kapsuleTag {
		lbResp, err := human.Marshal(loadbalancer, opt)
		if err != nil {
			return "", err
		}

		return strings.Join([]string{
			lbResp,
			warningKapsuleTaggedMessageView(),
		}, "\n\n"), nil
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
		Groups:    []string{"workflow"},
		ArgsType:  reflect.TypeOf(lb.ZonedAPIWaitForLBRequest{}),
		Run: func(ctx context.Context, argsI any) (i any, err error) {
			api := lb.NewZonedAPI(core.ExtractClient(ctx))
			args := argsI.(*lb.ZonedAPIWaitForLBRequest)

			return api.WaitForLb(&lb.ZonedAPIWaitForLBRequest{
				LBID:          args.LBID,
				Zone:          args.Zone,
				RetryInterval: core.DefaultRetryInterval,
				Timeout:       args.Timeout,
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
			core.WaitTimeoutArgSpec(defaultLBTimeout),
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
	c.ArgSpecs.GetByName("type").ValidateFunc = func(_ *core.ArgSpec, _ any) error {
		// Allow all lb types
		return nil
	}

	c.WaitFunc = func(ctx context.Context, _, respI any) (any, error) {
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
	c.ArgSpecs.GetByName("type").ValidateFunc = func(_ *core.ArgSpec, _ any) error {
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
	type lbUpdateRequestCustom struct {
		*lb.ZonedAPIUpdateLBRequest
		AssignFlexibleIPv6 bool   `json:"assign_flexible_ipv6"`
		IPID               string `json:"ip_id"`
	}

	c.ArgsType = reflect.TypeOf(lbUpdateRequestCustom{})

	c.ArgSpecs.AddBefore("tags.{index}", &core.ArgSpec{
		Name:       "assign-flexible-ipv6",
		Short:      "Automatically assign a flexible public IPv6 to the Load Balancer",
		OneOfGroup: "ip",
	})
	c.ArgSpecs.AddBefore("tags.{index}", &core.ArgSpec{
		Name:       "ip-id",
		Short:      "The IP ID to attach to the Load Balancer",
		OneOfGroup: "ip",
	})

	c.Run = func(ctx context.Context, argsI any) (any, error) {
		request := argsI.(*lbUpdateRequestCustom)
		client := core.ExtractClient(ctx)
		lbAPI := lb.NewZonedAPI(client)

		waitRequest := &lb.ZonedAPIWaitForLBRequest{
			LBID:          request.LBID,
			Zone:          request.Zone,
			Timeout:       scw.TimeDurationPtr(defaultLBTimeout),
			RetryInterval: core.DefaultRetryInterval,
		}
		res, err := lbAPI.WaitForLb(waitRequest, scw.WithContext(ctx))
		if err != nil {
			return nil, err
		}

		if request.IPID != "" {
			_, err = lbAPI.UpdateIP(&lb.ZonedAPIUpdateIPRequest{
				Zone: request.Zone,
				IPID: request.IPID,
				LBID: &request.LBID,
			}, scw.WithContext(ctx))
			if err != nil {
				return nil, err
			}
		}

		if request.AssignFlexibleIPv6 {
			ip, err := lbAPI.CreateIP(&lb.ZonedAPICreateIPRequest{
				Zone:      res.Zone,
				ProjectID: &res.ProjectID,
				IsIPv6:    true,
			}, scw.WithContext(ctx))
			if err != nil {
				return nil, err
			}

			_, err = lbAPI.UpdateIP(&lb.ZonedAPIUpdateIPRequest{
				Zone: ip.Zone,
				IPID: ip.ID,
				LBID: &res.ID,
			}, scw.WithContext(ctx))
			if err != nil {
				return nil, err
			}
		}

		_, err = lbAPI.WaitForLb(waitRequest, scw.WithContext(ctx))
		if err != nil {
			return nil, err
		}

		result, err := lbAPI.UpdateLB(request.ZonedAPIUpdateLBRequest)
		if err != nil {
			return nil, err
		}

		return result, err
	}

	c.Interceptor = interceptLB()

	return c
}

func lbDeleteBuilder(c *core.Command) *core.Command {
	c.WaitFunc = func(ctx context.Context, argsI any, _ any) (any, error) {
		api := lb.NewZonedAPI(core.ExtractClient(ctx))
		waitForLb, err := api.WaitForLb(&lb.ZonedAPIWaitForLBRequest{
			LBID:          argsI.(*lb.ZonedAPIDeleteLBRequest).LBID,
			Zone:          argsI.(*lb.ZonedAPIDeleteLBRequest).Zone,
			RetryInterval: core.DefaultRetryInterval,
		})
		if err != nil {
			notFoundError := &scw.ResourceNotFoundError{}
			responseError := &scw.ResponseError{}
			if errors.As(err, &responseError) && responseError.StatusCode == http.StatusNotFound ||
				errors.As(err, &notFoundError) {
				return &core.SuccessResult{
					Resource: "lb",
					Verb:     "delete",
				}, nil
			}

			return nil, err
		}

		return waitForLb, nil
	}
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
	return func(ctx context.Context, argsI any, runner core.CommandRunner) (any, error) {
		var getLB *lb.LB
		var err error

		client := core.ExtractClient(ctx)
		api := lb.NewZonedAPI(client)

		if _, ok := argsI.(*lb.ZonedAPIDeleteLBRequest); ok {
			getLB, err = api.GetLB(&lb.ZonedAPIGetLBRequest{
				Zone: argsI.(*lb.ZonedAPIDeleteLBRequest).Zone,
				LBID: argsI.(*lb.ZonedAPIDeleteLBRequest).LBID,
			})
			if err != nil {
				return nil, err
			}
		}

		res, err := runner(ctx, argsI)
		if err != nil {
			return nil, err
		}

		if _, ok := res.(*core.SuccessResult); ok {
			if len(getLB.Tags) != 0 && getLB.Tags[0] == kapsuleTag {
				return warningKapsuleTaggedMessageView(), nil
			}
		}

		return res, nil
	}
}
