package rdb

import (
	"context"
	"fmt"
	"reflect"

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

type rdbACLCustomArgs struct {
	Region     scw.Region
	InstanceID string
	ACLRuleIPs scw.IPNet
}

type rdbACLCustomResult struct {
	Rules   []*rdb.ACLRule
	Success core.SuccessResult
}

func rdbACLCustomResultMarshalerFunc(i interface{}, opt *human.MarshalOpt) (string, error) {
	result := i.(rdbACLCustomResult)
	messageStr, err := result.Success.MarshalHuman()
	if err != nil {
		return "", err
	}
	aclStr, err := human.Marshal(result.Rules, opt)
	if err != nil {
		return "", err
	}
	return messageStr + "\n" + aclStr, nil
}

func aclAddBuilder(c *core.Command) *core.Command {
	c.ArgsType = reflect.TypeOf(rdbACLCustomArgs{})
	c.ArgSpecs = core.ArgSpecs{
		{
			Name:       "acl-rule-ips",
			Short:      "IP addresses defined in the ACL rules of the Database Instance",
			Required:   true,
			Positional: true,
		},
		{
			Name:       "instance-id",
			Short:      "ID of the Database Instance",
			Required:   true,
			Positional: false,
		},
		core.RegionArgSpec(),
	}

	c.Interceptor = func(ctx context.Context, argsI interface{}, runner core.CommandRunner) (interface{}, error) {
		respI, err := runner(ctx, argsI)
		if err != nil {
			return nil, err
		}
		return respI.(*rdbACLCustomResult), nil
	}

	c.Run = func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
		args := argsI.(*rdbACLCustomArgs)
		client := core.ExtractClient(ctx)
		api := rdb.NewAPI(client)

		rule, err := api.AddInstanceACLRules(&rdb.AddInstanceACLRulesRequest{
			Region:     args.Region,
			InstanceID: args.InstanceID,
			Rules: []*rdb.ACLRuleRequest{
				{
					IP:          args.ACLRuleIPs,
					Description: fmt.Sprintf("Allow %s", args.ACLRuleIPs.String()),
				},
			},
		}, scw.WithContext(ctx))
		if err != nil {
			return nil, fmt.Errorf("failed to add ACL rule: %w", err)
		}

		return &rdbACLCustomResult{
			Rules: rule.Rules,
			Success: core.SuccessResult{
				Message: fmt.Sprintf("ACL rule %s successfully added", args.ACLRuleIPs.String()),
			},
		}, nil
	}

	c.WaitFunc = func(ctx context.Context, argsI, respI interface{}) (interface{}, error) {
		args := argsI.(*rdbACLCustomArgs)
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

		return respI.(*rdbACLCustomResult), nil
	}

	return c
}

func aclDeleteBuilder(c *core.Command) *core.Command {
	c.ArgsType = reflect.TypeOf(rdbACLCustomArgs{})
	c.ArgSpecs = core.ArgSpecs{
		{
			Name:       "acl-rule-ips",
			Short:      "IP addresses defined in the ACL rules of the Database Instance",
			Required:   true,
			Positional: true,
		},
		{
			Name:       "instance-id",
			Short:      "ID of the Database Instance",
			Required:   true,
			Positional: false,
		},
		core.RegionArgSpec(),
	}

	c.Interceptor = func(ctx context.Context, argsI interface{}, runner core.CommandRunner) (interface{}, error) {
		respI, err := runner(ctx, argsI)
		if err != nil {
			return nil, err
		}
		api := rdb.NewAPI(core.ExtractClient(ctx))
		args := argsI.(*rdbACLCustomArgs)
		rules, err := api.ListInstanceACLRules(&rdb.ListInstanceACLRulesRequest{
			Region:     args.Region,
			InstanceID: args.InstanceID,
		}, scw.WithContext(ctx), scw.WithAllPages())
		if err != nil {
			return nil, fmt.Errorf("failed to list ACL rules: %w", err)
		}

		resp := respI.(*rdbACLCustomResult)
		resp.Rules = rules.Rules

		return resp, nil
	}

	c.Run = func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
		args := argsI.(*rdbACLCustomArgs)
		client := core.ExtractClient(ctx)
		api := rdb.NewAPI(client)

		// The API returns 200 OK even if the rule was not set in the first place, so we have to check if the rule was present
		// before deleting it to warn them if nothing was done
		ruleWasSet := false
		rules, err := api.ListInstanceACLRules(&rdb.ListInstanceACLRulesRequest{
			Region:     args.Region,
			InstanceID: args.InstanceID,
		}, scw.WithContext(ctx), scw.WithAllPages())
		if err != nil {
			return nil, fmt.Errorf("failed to list ACL rules: %w", err)
		}
		for _, rule := range rules.Rules {
			if rule.IP.String() == args.ACLRuleIPs.String() {
				ruleWasSet = true
			}
		}

		_, err = api.DeleteInstanceACLRules(&rdb.DeleteInstanceACLRulesRequest{
			Region:     args.Region,
			InstanceID: args.InstanceID,
			ACLRuleIPs: []string{args.ACLRuleIPs.String()},
		}, scw.WithContext(ctx))
		if err != nil {
			return nil, fmt.Errorf("failed to remove ACL rule: %w", err)
		}

		message := ""
		if ruleWasSet {
			message = fmt.Sprintf("ACL rule %s successfully deleted", args.ACLRuleIPs.String())
		} else {
			message = fmt.Sprintf("ACL rule %s was not set", args.ACLRuleIPs.String())
		}

		return &rdbACLCustomResult{
			Success: core.SuccessResult{
				Message: message,
			},
		}, nil
	}

	c.WaitFunc = func(ctx context.Context, argsI, respI interface{}) (interface{}, error) {
		args := argsI.(*rdbACLCustomArgs)
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

		return respI.(*rdbACLCustomResult), nil
	}

	return c
}

func aclSetBuilder(c *core.Command) *core.Command {
	c.Run = func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
		args := argsI.(*rdb.SetInstanceACLRulesRequest)
		client := core.ExtractClient(ctx)
		api := rdb.NewAPI(client)

		rule, err := api.SetInstanceACLRules(args, scw.WithContext(ctx))
		if err != nil {
			return nil, fmt.Errorf("failed to set ACL rule: %w", err)
		}

		return &rdbACLCustomResult{
			Rules: rule.Rules,
			Success: core.SuccessResult{
				Message: "ACL rules successfully set",
			},
		}, nil
	}

	c.WaitFunc = func(ctx context.Context, argsI, respI interface{}) (interface{}, error) {
		args := argsI.(*rdb.SetInstanceACLRulesRequest)
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

		return respI.(*rdbACLCustomResult), nil
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
				Short:      "ID of the Database Instance",
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
