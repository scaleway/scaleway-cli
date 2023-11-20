package rdb

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/internal/core"
	"github.com/scaleway/scaleway-cli/v2/internal/editor"
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
		aclResult := make([]string, 0, len(aclDeleteResponse.Rules))

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

var rdbACLEditYamlExample = `rules:
- description: your description
  ip: 0.0.0.0/0
`

type rdbACLEditArgs struct {
	Region     scw.Region
	InstanceID string
	Mode       editor.MarshalMode
}

func aclEditCommand() *core.Command {
	return &core.Command{
		Short:     "Edit a database instance's ACL",
		Long:      editor.LongDescription,
		Namespace: "rdb",
		Resource:  "acl",
		Verb:      "edit",
		ArgsType:  reflect.TypeOf(rdbACLEditArgs{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "instance-id",
				Short:      "ID of the database instance ",
				Required:   true,
				Positional: true,
			},
			editor.MarshalModeArgSpec(),
			core.RegionArgSpec(),
		},
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
			args := argsI.(*rdbACLEditArgs)

			client := core.ExtractClient(ctx)
			api := rdb.NewAPI(client)

			setRequest := &rdb.SetInstanceACLRulesRequest{
				Region:     args.Region,
				InstanceID: args.InstanceID,
			}

			rules, err := api.ListInstanceACLRules(&rdb.ListInstanceACLRulesRequest{
				Region:     args.Region,
				InstanceID: args.InstanceID,
			}, scw.WithAllPages(), scw.WithContext(ctx))
			if err != nil {
				return nil, fmt.Errorf("failed to list ACL rules: %w", err)
			}

			editedSetRequest, err := editor.UpdateResourceEditor(rules, setRequest, &editor.Config{
				PutRequest:  true,
				MarshalMode: args.Mode,
				Template:    rdbACLEditYamlExample,
			})
			if err != nil {
				return nil, err
			}

			setRequest = editedSetRequest.(*rdb.SetInstanceACLRulesRequest)

			resp, err := api.SetInstanceACLRules(setRequest, scw.WithContext(ctx))
			if err != nil {
				return nil, fmt.Errorf("failed to set ACL rules: %w", err)
			}

			return resp.Rules, nil
		},
	}
}
