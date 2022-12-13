package lb

import (
	"context"
	"strings"

	"github.com/scaleway/scaleway-cli/v2/internal/core"
	"github.com/scaleway/scaleway-cli/v2/internal/human"
	"github.com/scaleway/scaleway-sdk-go/api/lb/v1"
)

func frontendGetBuilder(c *core.Command) *core.Command {
	c.View = &core.View{
		Sections: []*core.ViewSection{
			{
				FieldName: "LB",
			},
			{
				FieldName: "Backend",
			},
		},
	}
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
	return func(ctx context.Context, argsI interface{}, runner core.CommandRunner) (interface{}, error) {
		client := core.ExtractClient(ctx)
		api := lb.NewZonedAPI(client)

		res, err := runner(ctx, argsI)
		if err != nil {
			return nil, err
		}

		switch res.(type) {
		case *lb.Frontend:
			frontendResp, err := human.Marshal(res.(*lb.Frontend), nil)
			if err != nil {
				return "", err
			}
			if len(res.(*lb.Frontend).LB.Tags) != 0 && res.(*lb.Frontend).LB.Tags[0] == kapsuleTag {
				return strings.Join([]string{
					frontendResp,
					warningKapsuleTaggedMessageView(),
				}, "\n\n"), nil
			}
		case *core.SuccessResult:
			getFrontend, err := api.GetFrontend(&lb.ZonedAPIGetFrontendRequest{
				Zone:       argsI.(*lb.ZonedAPIDeleteFrontendRequest).Zone,
				FrontendID: argsI.(*lb.ZonedAPIDeleteFrontendRequest).FrontendID,
			})
			if err != nil {
				return nil, err
			}

			if len(getFrontend.LB.Tags) != 0 && getFrontend.LB.Tags[0] == kapsuleTag {
				frontendResp, err := human.Marshal(res.(*core.SuccessResult), nil)
				if err != nil {
					return "", err
				}

				return strings.Join([]string{
					frontendResp,
					warningKapsuleTaggedMessageView(),
				}, "\n\n"), nil
			}
		}

		return res, nil
	}
}
