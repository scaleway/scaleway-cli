package iam

import (
	"context"
	"fmt"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	iam "github.com/scaleway/scaleway-sdk-go/api/iam/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type iamRuleCreateCommandRequest struct {
	PolicyID string
	iam.RuleSpecs
}

func iamRuleCreateCommand() *core.Command {
	return &core.Command{
		Namespace: "iam",
		Resource:  "rule",
		Verb:      "create",
		Short:     "Create a rule for a specific IAM policy",
		ArgsType:  reflect.TypeOf(iamRuleCreateCommandRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "policy-id",
				Short:      "Id of policy to update",
				Positional: true,
			},
			{
				Name:       "permission-set-names.{index}",
				Short:      `Names of permission sets bound to the rule`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "project-ids.{index}",
				Short:      `List of Project IDs the rule is scoped to`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "organization-id",
				Short:      `ID of Organization the rule is scoped to`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Examples: nil,
		SeeAlsos: nil,
		Run: func(ctx context.Context, argsI any) (any, error) {
			args := argsI.(*iamRuleCreateCommandRequest)
			api := iam.NewAPI(core.ExtractClient(ctx))

			resp, err := api.ListRules(&iam.ListRulesRequest{
				PolicyID: args.PolicyID,
			}, scw.WithContext(ctx), scw.WithAllPages())
			if err != nil {
				return nil, err
			}

			rulesSpecs := make([]*iam.RuleSpecs, 0, len(resp.Rules)+1)
			for _, rule := range resp.Rules {
				rulesSpecs = append(rulesSpecs, &iam.RuleSpecs{
					PermissionSetNames: rule.PermissionSetNames,
					ProjectIDs:         rule.ProjectIDs,
					OrganizationID:     rule.OrganizationID,
				})
			}

			rulesSpecs = append(rulesSpecs, &args.RuleSpecs)

			policy, err := api.SetRules(&iam.SetRulesRequest{
				PolicyID: args.PolicyID,
				Rules:    rulesSpecs,
			}, scw.WithContext(ctx))
			if err != nil {
				return nil, err
			}

			return policy.Rules, err
		},
		Groups: []string{"utility"},
	}
}

type iamRuleDeleteCommandRequest struct {
	PolicyID string
	RuleID   string
}

func iamRuleDeleteCommand() *core.Command {
	return &core.Command{
		Namespace: "iam",
		Resource:  "rule",
		Verb:      "delete",
		Short:     "Delete a rule for a specific IAM policy",
		ArgsType:  reflect.TypeOf(iamRuleDeleteCommandRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "policy-id",
				Short:      "Id of policy to update",
				Positional: true,
			},
			{
				Name:  "rule-id",
				Short: "Id of rule to delete",
			},
		},
		Examples: nil,
		SeeAlsos: nil,
		Run: func(ctx context.Context, argsI any) (any, error) {
			args := argsI.(*iamRuleDeleteCommandRequest)
			api := iam.NewAPI(core.ExtractClient(ctx))

			resp, err := api.ListRules(&iam.ListRulesRequest{
				PolicyID: args.PolicyID,
			}, scw.WithContext(ctx), scw.WithAllPages())
			if err != nil {
				return nil, err
			}

			found := false
			rulesSpecs := make([]*iam.RuleSpecs, 0, len(resp.Rules))
			for _, rule := range resp.Rules {
				if rule.ID == args.RuleID {
					found = true
				} else {
					rulesSpecs = append(rulesSpecs, &iam.RuleSpecs{
						PermissionSetNames: rule.PermissionSetNames,
						ProjectIDs:         rule.ProjectIDs,
						OrganizationID:     rule.OrganizationID,
					})
				}
			}

			if !found {
				return nil, fmt.Errorf("rule %s not found for given policy", args.RuleID)
			}

			policy, err := api.SetRules(&iam.SetRulesRequest{
				PolicyID: args.PolicyID,
				Rules:    rulesSpecs,
			}, scw.WithContext(ctx))
			if err != nil {
				return nil, err
			}

			return policy.Rules, err
		},
		Groups: []string{"utility"},
	}
}
