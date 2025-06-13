package lb

import (
	"context"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	"github.com/scaleway/scaleway-sdk-go/api/lb/v1"
)

func lbFrontendMarshalerFunc(i any, opt *human.MarshalOpt) (string, error) {
	type tmp lb.Frontend
	frontend := tmp(i.(lb.Frontend))

	opt.Sections = []*human.MarshalSection{
		{
			FieldName: "LB",
		},
		{
			FieldName: "Backend",
		},
	}

	str, err := human.Marshal(frontend, opt)
	if err != nil {
		return "", err
	}

	return str, nil
}

func frontendGetBuilder(c *core.Command) *core.Command {
	c.Interceptor = interceptFrontend()

	return c
}

func frontendCreateBuilder(c *core.Command) *core.Command {
	c.Interceptor = interceptFrontend()

	return c
}

func frontendUpdateBuilder(c *core.Command) *core.Command {
	c.Interceptor = interceptFrontend()

	return c
}

func frontendDeleteBuilder(c *core.Command) *core.Command {
	c.Interceptor = interceptFrontend()

	return c
}

func interceptFrontend() core.CommandInterceptor {
	return func(ctx context.Context, argsI any, runner core.CommandRunner) (any, error) {
		var getFrontend *lb.Frontend
		var err error

		client := core.ExtractClient(ctx)
		api := lb.NewZonedAPI(client)

		if _, ok := argsI.(*lb.ZonedAPIDeleteFrontendRequest); ok {
			getFrontend, err = api.GetFrontend(&lb.ZonedAPIGetFrontendRequest{
				Zone:       argsI.(*lb.ZonedAPIDeleteFrontendRequest).Zone,
				FrontendID: argsI.(*lb.ZonedAPIDeleteFrontendRequest).FrontendID,
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
			if len(getFrontend.LB.Tags) != 0 && getFrontend.LB.Tags[0] == kapsuleTag {
				return warningKapsuleTaggedMessageView(), nil
			}
		}

		return res, nil
	}
}
