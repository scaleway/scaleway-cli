package iam

import (
	"context"
	"fmt"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	iam "github.com/scaleway/scaleway-sdk-go/api/iam/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type apiKeyResponse struct {
	APIKey   *iam.APIKey
	UserType string              `json:"user_type"`
	Policies map[string][]string `json:"policies"`
}
type iamGetAPIKeyArgs struct {
	AccessKey    string
	WithPolicies bool
}

type apiKeyOptions struct {
	WithPolicies bool
}

func WithPolicies(withPolicies bool) apiKeyOptions {
	return apiKeyOptions{
		WithPolicies: withPolicies,
	}
}

func getApiKey(
	ctx context.Context,
	api *iam.API,
	accessKey string,
	options apiKeyOptions,
) (apiKeyResponse, error) {
	var response apiKeyResponse
	apiKey, err := api.GetAPIKey(&iam.GetAPIKeyRequest{
		AccessKey: accessKey,
	}, scw.WithContext(ctx))
	if err != nil {
		return response, err
	}

	user, err := api.GetUser(&iam.GetUserRequest{
		UserID: *apiKey.UserID,
	}, scw.WithContext(ctx))
	if err != nil {
		return response, err
	}

	response.APIKey = apiKey
	response.UserType = string(user.Type)

	if user.Type == iam.UserTypeOwner {
		response.UserType = fmt.Sprintf(
			"%s (owner has all permissions over the organization)",
			user.Type,
		)

		return response, nil
	}

	if options.WithPolicies {
		listPolicyRequest := &iam.ListPoliciesRequest{
			UserIDs: []string{*apiKey.UserID},
		}
		policies, err := api.ListPolicies(
			listPolicyRequest,
			scw.WithAllPages(),
			scw.WithContext(ctx),
		)
		if err != nil {
			return response, err
		}
		// Build a map of policies -> [rules...]
		policyMap := map[string][]string{}
		for _, policy := range policies.Policies {
			rules, err := api.ListRules(
				&iam.ListRulesRequest{
					PolicyID: policy.ID,
				},
				scw.WithContext(ctx),
			)
			if err != nil {
				return response, err
			}

			for _, rules := range rules.Rules {
				policyMap[fmt.Sprintf("%s (%s)", policy.Name, policy.ID)] = append(
					policyMap[fmt.Sprintf("%s (%s)", policy.Name, policy.ID)],
					*rules.PermissionSetNames...)
			}
		}
		response.Policies = policyMap
	}

	return response, nil
}

func apiKeyMarshalerFunc(i interface{}, opt *human.MarshalOpt) (string, error) {
	type tmp apiKeyResponse
	resp := tmp(i.(apiKeyResponse))

	opt.Sections = []*human.MarshalSection{
		{
			FieldName: "UserType",
		},
		{
			FieldName: "APIKey",
		},
		{
			FieldName:   "Policies",
			HideIfEmpty: true,
		},
	}

	return human.Marshal(resp, opt)
}

var iamApiKeyCustomBuilder = struct {
	argSpecs core.ArgSpecs
	argType  reflect.Type
	run      func(ctx context.Context, args interface{}) (i interface{}, e error)
}{
	argSpecs: core.ArgSpecs{
		{
			Name:       "access-key",
			Short:      `Access key to search for`,
			Required:   true,
			Deprecated: false,
			Positional: true,
		},
		{
			Name:       "with-policies",
			Short:      `Display the set of policies associated with the API key`,
			Default:    core.DefaultValueSetter("true"),
			Required:   false,
			Deprecated: false,
			Positional: false,
		},
	},
	argType: reflect.TypeOf(iamGetAPIKeyArgs{}),
	run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
		arguments := args.(*iamGetAPIKeyArgs)

		client := core.ExtractClient(ctx)
		api := iam.NewAPI(client)

		return getApiKey(ctx, api, arguments.AccessKey, apiKeyOptions{
			WithPolicies: arguments.WithPolicies,
		})
	},
}
