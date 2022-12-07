package lb

import (
	"context"
	"strings"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/internal/core"
	"github.com/scaleway/scaleway-cli/v2/internal/human"
	"github.com/scaleway/scaleway-sdk-go/api/lb/v1"
)

var (
	backendServerStatsHealthCheckStatusMarshalSpecs = human.EnumMarshalSpecs{
		lb.BackendServerStatsHealthCheckStatusPassed:   &human.EnumMarshalSpec{Attribute: color.FgGreen, Value: "passed"},
		lb.BackendServerStatsHealthCheckStatusFailed:   &human.EnumMarshalSpec{Attribute: color.FgRed, Value: "failed"},
		lb.BackendServerStatsHealthCheckStatusUnknown:  &human.EnumMarshalSpec{Attribute: color.Faint, Value: "unknown"},
		lb.BackendServerStatsHealthCheckStatusNeutral:  &human.EnumMarshalSpec{Attribute: color.Faint, Value: "neutral"},
		lb.BackendServerStatsHealthCheckStatusCondpass: &human.EnumMarshalSpec{Attribute: color.FgBlue, Value: "condition passed"},
	}
	backendServerStatsServerStateMarshalSpecs = human.EnumMarshalSpecs{
		lb.BackendServerStatsServerStateStopped:  &human.EnumMarshalSpec{Attribute: color.FgRed, Value: "stopped"},
		lb.BackendServerStatsServerStateStarting: &human.EnumMarshalSpec{Attribute: color.FgBlue, Value: "starting"},
		lb.BackendServerStatsServerStateRunning:  &human.EnumMarshalSpec{Attribute: color.FgGreen, Value: "running"},
		lb.BackendServerStatsServerStateStopping: &human.EnumMarshalSpec{Attribute: color.FgBlue, Value: "stopping"},
	}
)

func lbBackendMarshalerFunc(i interface{}, opt *human.MarshalOpt) (string, error) {
	type tmp lb.Backend
	backend := tmp(i.(lb.Backend))

	opt.Sections = []*human.MarshalSection{
		{
			FieldName: "HealthCheck",
		},
		{
			FieldName: "Pool",
		},
		{
			FieldName: "LB",
		},
	}

	str, err := human.Marshal(backend, opt)
	if err != nil {
		return "", err
	}

	return str, nil
}

func backendGetBuilder(c *core.Command) *core.Command {
	c.Interceptor = interceptBackend()
	return c
}

func backendCreateBuilder(c *core.Command) *core.Command {
	c.Interceptor = interceptBackend()
	return c
}

func backendUpdateBuilder(c *core.Command) *core.Command {
	c.Interceptor = interceptBackend()
	return c
}

func backendDeleteBuilder(c *core.Command) *core.Command {
	c.Interceptor = interceptBackend()
	return c
}

func backendAddServersBuilder(c *core.Command) *core.Command {
	c.Interceptor = interceptBackend()
	return c
}

func backendRemoveServersBuilder(c *core.Command) *core.Command {
	c.Interceptor = interceptBackend()
	return c
}

func backendSetServersBuilder(c *core.Command) *core.Command {
	c.Interceptor = interceptBackend()
	return c
}

func backendUpdateHealthcheckBuilder(c *core.Command) *core.Command {
	c.Interceptor = func(ctx context.Context, argsI interface{}, runner core.CommandRunner) (interface{}, error) {
		res, err := runner(ctx, argsI)
		if err != nil {
			return nil, err
		}

		backendResp, err := human.Marshal(res.(*lb.HealthCheck), nil)
		if err != nil {
			return "", err
		}

		client := core.ExtractClient(ctx)
		api := lb.NewZonedAPI(client)

		getBackend, err := api.GetBackend(&lb.ZonedAPIGetBackendRequest{
			Zone:      argsI.(*lb.ZonedAPIUpdateHealthCheckRequest).Zone,
			BackendID: argsI.(*lb.ZonedAPIUpdateHealthCheckRequest).BackendID,
		})
		if err != nil {
			return nil, err
		}

		if len(getBackend.LB.Tags) != 0 && getBackend.LB.Tags[0] == kapsuleTag {
			return strings.Join([]string{
				backendResp,
				warningKapsuleTaggedMessageView(),
			}, "\n\n"), nil
		}

		return res, nil
	}
	return c
}

func interceptBackend() core.CommandInterceptor {
	return func(ctx context.Context, argsI interface{}, runner core.CommandRunner) (interface{}, error) {
		res, err := runner(ctx, argsI)
		if err != nil {
			return nil, err
		}

		backendResp, err := human.Marshal(res.(*lb.Backend), nil)
		if err != nil {
			return "", err
		}

		if len(res.(*lb.Backend).LB.Tags[0]) != 0 && res.(*lb.Backend).LB.Tags[0] == kapsuleTag {
			return strings.Join([]string{
				backendResp,
				warningKapsuleTaggedMessageView(),
			}, "\n\n"), nil
		}

		return res, nil
	}
}
