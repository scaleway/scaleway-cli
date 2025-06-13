package lb

import (
	"context"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	"github.com/scaleway/scaleway-sdk-go/api/lb/v1"
)

var aclMarshalSpecs = human.EnumMarshalSpecs{
	lb.ACLActionTypeAllow: &human.EnumMarshalSpec{Attribute: color.FgGreen, Value: "allow"},
	lb.ACLActionTypeDeny:  &human.EnumMarshalSpec{Attribute: color.FgRed, Value: "deny"},
}

func lbACLMarshalerFunc(i any, opt *human.MarshalOpt) (string, error) {
	type tmp lb.ACL
	acl := tmp(i.(lb.ACL))

	opt.Sections = []*human.MarshalSection{
		{
			FieldName: "Frontend",
		},
	}

	str, err := human.Marshal(acl, opt)
	if err != nil {
		return "", err
	}

	return str, nil
}

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

		if _, ok := argsI.(*lb.ZonedAPIDeleteCertificateRequest); ok {
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
			if len(getACL.Frontend.LB.Tags) != 0 && getACL.Frontend.LB.Tags[0] == kapsuleTag {
				return warningKapsuleTaggedMessageView(), nil
			}
		}

		return res, nil
	}
}
