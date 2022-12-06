package rdb

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/internal/core"
	"github.com/scaleway/scaleway-cli/v2/internal/human"
	"github.com/scaleway/scaleway-sdk-go/api/rdb/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

var (
	aclRuleActionMarshalSpecs = human.EnumMarshalSpecs{
		rdb.ACLRuleActionAllow: &human.EnumMarshalSpec{Attribute: color.FgGreen, Value: "allow"},
		rdb.ACLRuleActionDeny:  &human.EnumMarshalSpec{Attribute: color.FgRed, Value: "deny"},
	}
)

func aclAddBuilder(c *core.Command) *core.Command {
	type customAddACLRequest struct {
		*rdb.AddInstanceACLRulesRequest
		Rules []*rdb.ACLRuleRequest
	}

	c.ArgSpecs.GetByName("rules.{index}.ip").Name = "rules.{index}.ip"
	c.ArgSpecs.GetByName("rules.{index}.description").Name = "rules.{index}.description"

	c.ArgsType = reflect.TypeOf(customAddACLRequest{})

	c.Interceptor = func(ctx context.Context, argsI interface{}, runner core.CommandRunner) (interface{}, error) {
		args := argsI.(*customAddACLRequest)

		request := args.AddInstanceACLRulesRequest
		request.Rules = args.Rules

		aclAddResponseI, err := runner(ctx, request)
		if err != nil {
			return nil, err
		}

		aclAddResponse := aclAddResponseI.(*rdb.AddInstanceACLRulesResponse)
		return aclAddResponse.Rules, nil
	}

	c.WaitFunc = func(ctx context.Context, argsI, respI interface{}) (interface{}, error) {
		args := argsI.(*customAddACLRequest)

		api := rdb.NewAPI(core.ExtractClient(ctx))
		_, err := api.WaitForInstance(&rdb.WaitForInstanceRequest{
			InstanceID:    args.InstanceID,
			Region:        args.Region,
			Timeout:       scw.TimeDurationPtr(instanceActionTimeout),
			RetryInterval: core.DefaultRetryInterval,
		})
		if err != nil {
			return nil, err
		}

		return respI.([]*rdb.ACLRule), nil
	}

	return c
}

func aclDeleteBuilder(c *core.Command) *core.Command {
	type deleteRules struct {
		IP string `json:"ip"`
	}

	type customDeleteACLRequest struct {
		*rdb.DeleteInstanceACLRulesRequest
		Rules []deleteRules
	}

	c.ArgSpecs.GetByName("acl-rule-ips.{index}").Name = "rules.{index}.ip"

	c.ArgsType = reflect.TypeOf(customDeleteACLRequest{})

	c.Interceptor = func(ctx context.Context, argsI interface{}, runner core.CommandRunner) (interface{}, error) {
		var aclResult []string

		args := argsI.(*customDeleteACLRequest)

		request := args.DeleteInstanceACLRulesRequest
		for _, ip := range args.Rules {
			request.ACLRuleIPs = append(request.ACLRuleIPs, ip.IP)
		}

		aclDeleteResponseI, err := runner(ctx, request)
		if err != nil {
			return nil, err
		}

		aclDeleteResponse := aclDeleteResponseI.(*rdb.DeleteInstanceACLRulesResponse)
		for i := 0; i < len(aclDeleteResponse.Rules); i++ {
			aclResult = append(aclResult, aclDeleteResponse.Rules[i].IP.String())
		}

		return &core.SuccessResult{
			Message: fmt.Sprintf("ACL rule(s) %s successfully deleted", strings.Trim(fmt.Sprint(aclResult), "[]")),
		}, nil
	}

	c.WaitFunc = func(ctx context.Context, argsI, respI interface{}) (interface{}, error) {
		args := argsI.(*customDeleteACLRequest)

		api := rdb.NewAPI(core.ExtractClient(ctx))
		_, err := api.WaitForInstance(&rdb.WaitForInstanceRequest{
			InstanceID:    args.InstanceID,
			Region:        args.Region,
			Timeout:       scw.TimeDurationPtr(instanceActionTimeout),
			RetryInterval: core.DefaultRetryInterval,
		})
		if err != nil {
			return nil, err
		}

		return respI.(*core.SuccessResult), nil
	}

	return c
}
