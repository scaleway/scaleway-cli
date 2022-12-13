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
	aclMarshalSpecs = human.EnumMarshalSpecs{
		lb.ACLActionTypeAllow: &human.EnumMarshalSpec{Attribute: color.FgGreen, Value: "allow"},
		lb.ACLActionTypeDeny:  &human.EnumMarshalSpec{Attribute: color.FgRed, Value: "deny"},
	}
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
	return func(ctx context.Context, argsI interface{}, runner core.CommandRunner) (interface{}, error) {
		client := core.ExtractClient(ctx)
		api := lb.NewZonedAPI(client)

		res, err := runner(ctx, argsI)
		if err != nil {
			return nil, err
		}

		switch res.(type) {
		case *lb.ACL:
			if len(res.(*lb.ACL).Frontend.LB.Tags) != 0 && res.(*lb.ACL).Frontend.LB.Tags[0] == kapsuleTag {
				ACLResp, err := human.Marshal(res.(*lb.ACL), nil)
				if err != nil {
					return "", err
				}

				return strings.Join([]string{
					ACLResp,
					warningKapsuleTaggedMessageView(),
				}, "\n\n"), nil
			}
		case *core.SuccessResult:
			getACL, err := api.GetACL(&lb.ZonedAPIGetACLRequest{
				Zone:  argsI.(*lb.ZonedAPIDeleteACLRequest).Zone,
				ACLID: argsI.(*lb.ZonedAPIDeleteACLRequest).ACLID,
			})
			if err != nil {
				return nil, err
			}

			if len(getACL.Frontend.LB.Tags) != 0 && getACL.Frontend.LB.Tags[0] == kapsuleTag {
				ACLResp, err := human.Marshal(res.(*core.SuccessResult), nil)
				if err != nil {
					return "", err
				}

				return strings.Join([]string{
					ACLResp,
					warningKapsuleTaggedMessageView(),
				}, "\n\n"), nil
			}
		}

		return res, nil
	}
}
