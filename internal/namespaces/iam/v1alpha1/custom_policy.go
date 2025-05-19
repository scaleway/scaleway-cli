package iam

import (
	"context"
	"fmt"

	"github.com/scaleway/scaleway-cli/v2/core"
	iam "github.com/scaleway/scaleway-sdk-go/api/iam/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func iamPolicyCreateBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("rules.{index}.permission-set-names.{index}").AutoCompleteFunc = func(ctx context.Context, _ string, _ any) core.AutocompleteSuggestions {
		client := core.ExtractClient(ctx)
		api := iam.NewAPI(client)
		// TODO: store result in a CLI cache
		resp, err := api.ListPermissionSets(&iam.ListPermissionSetsRequest{
			PageSize: scw.Uint32Ptr(100),
		}, scw.WithAllPages())
		if err != nil {
			return nil
		}
		suggestions := core.AutocompleteSuggestions{}
		for _, ps := range resp.PermissionSets {
			suggestions = append(suggestions, ps.Name)
		}

		return suggestions
	}

	return c
}

type PolicyGetInterceptorResponse struct {
	*iam.Policy
	Rules []*iam.Rule
}

func iamPolicyGetBuilder(c *core.Command) *core.Command {
	c.View = &core.View{
		Title: "Policy",
		Sections: []*core.ViewSection{
			{
				Title:     "Rules",
				FieldName: "Rules",
			},
		},
	}
	c.AddInterceptors(
		func(ctx context.Context, argsI interface{}, runner core.CommandRunner) (interface{}, error) {
			args := argsI.(*iam.GetPolicyRequest)
			api := iam.NewAPI(core.ExtractClient(ctx))

			respI, err := runner(ctx, argsI)
			if err != nil {
				return respI, err
			}
			resp := &PolicyGetInterceptorResponse{
				Policy: respI.(*iam.Policy),
			}

			rules, err := api.ListRules(&iam.ListRulesRequest{
				PolicyID: args.PolicyID,
			}, scw.WithContext(ctx), scw.WithAllPages())
			if err != nil {
				return nil, fmt.Errorf("failed to list rules for given policy: %w", err)
			}
			resp.Rules = rules.Rules

			return resp, nil
		},
	)

	return c
}
