package rdb

import (
	"context"
	"fmt"
	"reflect"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	"github.com/scaleway/scaleway-cli/v2/internal/editor"
	"github.com/scaleway/scaleway-sdk-go/api/rdb/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

var aclRuleActionMarshalSpecs = human.EnumMarshalSpecs{
	rdb.ACLRuleActionAllow: &human.EnumMarshalSpec{Attribute: color.FgGreen, Value: "allow"},
	rdb.ACLRuleActionDeny:  &human.EnumMarshalSpec{Attribute: color.FgRed, Value: "deny"},
}

type rdbACLCustomArgs struct {
	Region     scw.Region
	InstanceID string
	ACLRuleIPs []scw.IPNet
}

type rdbACLAddCustomArgs struct {
	Region       scw.Region
	InstanceID   string
	ACLRuleIPs   []scw.IPNet
	Description  string
	Descriptions []string
}

type CustomACLResult struct {
	Rules   []*rdb.ACLRule
	Success core.SuccessResult
}

func rdbACLCustomResultMarshalerFunc(i any, opt *human.MarshalOpt) (string, error) {
	result := i.(CustomACLResult)
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
	c.ArgsType = reflect.TypeOf(rdbACLAddCustomArgs{})
	c.ArgSpecs = core.ArgSpecs{
		{
			Name:       "acl-rule-ips",
			Short:      "IP addresses defined in the ACL rules of the Database Instance",
			Required:   false,
			Positional: true,
		},
		{
			Name:       "instance-id",
			Short:      "ID of the Database Instance",
			Required:   true,
			Positional: false,
		},
		{
			Name:       "description",
			Short:      "Description of the ACL rule. If multiple IPs are provided, this description will be applied to all rules unless specific descriptions are provided.",
			Required:   false,
			Positional: false,
		},
		{
			Name:       "descriptions",
			Short:      "Descriptions of the ACL rules",
			Required:   false,
			Positional: false,
		},
		core.RegionArgSpec(),
	}
	c.AcceptMultiplePositionalArgs = true

	c.Interceptor = func(ctx context.Context, argsI any, runner core.CommandRunner) (any, error) {
		respI, err := runner(ctx, argsI)
		if err != nil {
			return nil, err
		}

		return respI.(*CustomACLResult), nil
	}

	c.Run = func(ctx context.Context, argsI any) (i any, e error) {
		args := argsI.(*rdbACLAddCustomArgs)
		client := core.ExtractClient(ctx)
		api := rdb.NewAPI(client)

		// Build rules with general and specific descriptions
		rules := make([]*rdb.ACLRuleRequest, 0, len(args.ACLRuleIPs))
		for i, ip := range args.ACLRuleIPs {
			description := args.Description
			if description == "" {
				description = "Allow " + ip.String()
			}
			if i < len(args.Descriptions) && args.Descriptions[i] != "" {
				description = args.Descriptions[i]
			}
			rules = append(rules, &rdb.ACLRuleRequest{
				IP:          ip,
				Description: description,
			})
		}

		rule, err := api.AddInstanceACLRules(&rdb.AddInstanceACLRulesRequest{
			Region:     args.Region,
			InstanceID: args.InstanceID,
			Rules:      rules,
		}, scw.WithContext(ctx))
		if err != nil {
			return nil, fmt.Errorf("failed to add ACL rule: %w", err)
		}

		// Create success message
		var message string
		if len(args.ACLRuleIPs) == 1 {
			message = fmt.Sprintf("ACL rule %s successfully added", args.ACLRuleIPs[0].String())
		} else {
			message = fmt.Sprintf("%d ACL rules successfully added", len(args.ACLRuleIPs))
		}

		return &CustomACLResult{
			Rules: rule.Rules,
			Success: core.SuccessResult{
				Message: message,
			},
		}, nil
	}

	c.WaitFunc = func(ctx context.Context, argsI, respI any) (any, error) {
		args := argsI.(*rdbACLAddCustomArgs)
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

		return respI.(*CustomACLResult), nil
	}

	return c
}

func aclDeleteBuilder(c *core.Command) *core.Command {
	c.ArgsType = reflect.TypeOf(rdbACLCustomArgs{})
	c.ArgSpecs = core.ArgSpecs{
		{
			Name:       "acl-rule-ips",
			Short:      "IP addresses defined in the ACL rules of the Database Instance",
			Required:   false,
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
	c.AcceptMultiplePositionalArgs = true

	c.Interceptor = func(ctx context.Context, argsI any, runner core.CommandRunner) (any, error) {
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

		resp := respI.(*CustomACLResult)
		resp.Rules = rules.Rules

		return resp, nil
	}

	c.Run = func(ctx context.Context, argsI any) (i any, e error) {
		args := argsI.(*rdbACLCustomArgs)
		client := core.ExtractClient(ctx)
		api := rdb.NewAPI(client)

		// The API returns 200 OK even if the rule was not set in the first place, so we have to check if the rule was present
		// before deleting it to warn them if nothing was done
		rules, err := api.ListInstanceACLRules(&rdb.ListInstanceACLRulesRequest{
			Region:     args.Region,
			InstanceID: args.InstanceID,
		}, scw.WithContext(ctx), scw.WithAllPages())
		if err != nil {
			return nil, fmt.Errorf("failed to list ACL rules: %w", err)
		}

		// Check which rules were actually set
		existingIPs := make(map[string]bool)
		for _, rule := range rules.Rules {
			existingIPs[rule.IP.String()] = true
		}

		// Convert IPs to strings for deletion
		ipStrings := make([]string, len(args.ACLRuleIPs))
		for i, ip := range args.ACLRuleIPs {
			ipStrings[i] = ip.String()
		}

		_, err = api.DeleteInstanceACLRules(&rdb.DeleteInstanceACLRulesRequest{
			Region:     args.Region,
			InstanceID: args.InstanceID,
			ACLRuleIPs: ipStrings,
		}, scw.WithContext(ctx))
		if err != nil {
			return nil, fmt.Errorf("failed to remove ACL rules: %w", err)
		}

		// Count how many rules were actually deleted
		deletedCount := 0
		for _, ip := range args.ACLRuleIPs {
			if existingIPs[ip.String()] {
				deletedCount++
			}
		}

		var message string
		if len(args.ACLRuleIPs) == 1 {
			if deletedCount > 0 {
				message = fmt.Sprintf(
					"ACL rule %s successfully deleted",
					args.ACLRuleIPs[0].String(),
				)
			} else {
				message = fmt.Sprintf("ACL rule %s was not set", args.ACLRuleIPs[0].String())
			}
		} else {
			message = fmt.Sprintf("%d ACL rules successfully deleted", deletedCount)
			if deletedCount < len(args.ACLRuleIPs) {
				message += fmt.Sprintf(" (%d were not set)", len(args.ACLRuleIPs)-deletedCount)
			}
		}

		return &CustomACLResult{
			Success: core.SuccessResult{
				Message: message,
			},
		}, nil
	}

	c.WaitFunc = func(ctx context.Context, argsI, respI any) (any, error) {
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

		return respI.(*CustomACLResult), nil
	}

	return c
}

type rdbACLSetCustomArgs struct {
	Region       scw.Region
	InstanceID   string
	ACLRuleIPs   []scw.IPNet
	Descriptions []string
}

func aclSetBuilder(c *core.Command) *core.Command {
	c.ArgsType = reflect.TypeOf(rdbACLSetCustomArgs{})
	c.ArgSpecs = core.ArgSpecs{
		{
			Name:       "acl-rule-ips",
			Short:      "IP addresses defined in the ACL rules of the Database Instance",
			Required:   false,
			Positional: true,
		},
		{
			Name:       "instance-id",
			Short:      "ID of the Database Instance",
			Required:   true,
			Positional: false,
		},
		{
			Name:       "descriptions",
			Short:      "Descriptions of the ACL rules",
			Required:   false,
			Positional: false,
		},
		core.RegionArgSpec(),
	}
	c.AcceptMultiplePositionalArgs = true

	c.Run = func(ctx context.Context, argsI any) (i any, e error) {
		args := argsI.(*rdbACLSetCustomArgs)
		client := core.ExtractClient(ctx)
		api := rdb.NewAPI(client)

		aclRules := []*rdb.ACLRuleRequest(nil)
		for _, ip := range args.ACLRuleIPs {
			aclRules = append(aclRules, &rdb.ACLRuleRequest{
				IP:          ip,
				Description: "Allow " + ip.String(),
			})
		}

		for i, desc := range args.Descriptions {
			aclRules[i].Description = desc
		}

		rule, err := api.SetInstanceACLRules(&rdb.SetInstanceACLRulesRequest{
			Region:     args.Region,
			InstanceID: args.InstanceID,
			Rules:      aclRules,
		}, scw.WithContext(ctx))
		if err != nil {
			return nil, fmt.Errorf("failed to set ACL rule: %w", err)
		}

		return &CustomACLResult{
			Rules: rule.Rules,
			Success: core.SuccessResult{
				Message: "ACL rules successfully set",
			},
		}, nil
	}

	c.WaitFunc = func(ctx context.Context, argsI, respI any) (any, error) {
		args := argsI.(*rdbACLSetCustomArgs)
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

		return respI.(*CustomACLResult), nil
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
		Run: func(ctx context.Context, argsI any) (i any, e error) {
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
