package rdb

import (
	"context"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-cli/internal/human"
	"github.com/scaleway/scaleway-sdk-go/api/rdb/v1"
)

var (
	aclRuleActionMarshalSpecs = human.EnumMarshalSpecs{
		rdb.ACLRuleActionAllow: &human.EnumMarshalSpec{Attribute: color.FgGreen, Value: "allow"},
		rdb.ACLRuleActionDeny:  &human.EnumMarshalSpec{Attribute: color.FgRed, Value: "deny"},
	}
)

func aclAddBuilder(c *core.Command) *core.Command {
	c.Interceptor = func(ctx context.Context, argsI interface{}, runner core.CommandRunner) (interface{}, error) {
		aclAddResponseI, err := runner(ctx, argsI)
		if err != nil {
			return nil, err
		}
		aclAddResponse := aclAddResponseI.(*rdb.AddInstanceACLRulesResponse)
		return aclAddResponse.Rules, nil
	}

	return c
}

func aclDeleteBuilder(c *core.Command) *core.Command {
	c.Interceptor = func(ctx context.Context, argsI interface{}, runner core.CommandRunner) (interface{}, error) {
		aclDeleteResponseI, err := runner(ctx, argsI)
		if err != nil {
			return nil, err
		}
		aclDeleteResponse := aclDeleteResponseI.(*rdb.DeleteInstanceACLRulesResponse)
		return rdb.ListInstanceACLRulesResponse{
			Rules:      aclDeleteResponse.Rules,
			TotalCount: uint32(len(aclDeleteResponse.Rules)),
		}, nil
	}

	return c
}
