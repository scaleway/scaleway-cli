package lb

import (
	"context"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-sdk-go/api/lb/v1"
)

func ACLGetBuilder(c *core.Command) *core.Command {
	c.Interceptor = interceptACL()

	return c
}

func ACLCreateBuilder(c *core.Command) *core.Command {
	c.Interceptor = interceptACL()

	return c
}

func ACLUpdateBuilder(c *core.Command) *core.Command {
	c.Interceptor = interceptACL()

	return c
}

func ACLDeleteBuilder(c *core.Command) *core.Command {
	c.Interceptor = interceptACL()

	return c
}

func interceptACL() core.CommandInterceptor {
	return func(ctx context.Context, argsI any, runner core.CommandRunner) (any, error) {
		var getACL *lb.ACL
		var err error

		client := core.ExtractClient(ctx)
		api := lb.NewZonedAPI(client)

		if _, ok := argsI.(*lb.ZonedAPIDeleteACLRequest); ok {
			getACL, err = api.GetACL(&lb.ZonedAPIGetACLRequest{
				Zone:  argsI.(*lb.ZonedAPIDeleteACLRequest).Zone,
				ACLID: argsI.(*lb.ZonedAPIDeleteACLRequest).ACLID,
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
			if getACL != nil && getACL.Frontend != nil && getACL.Frontend.LB != nil &&
				len(getACL.Frontend.LB.Tags) != 0 &&
				getACL.Frontend.LB.Tags[0] == kapsuleTag {
				return warningKapsuleTaggedMessageView(), nil
			}
		}

		return res, nil
	}
}
