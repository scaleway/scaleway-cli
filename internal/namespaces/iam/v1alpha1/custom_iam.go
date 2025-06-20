package iam

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	iam "github.com/scaleway/scaleway-sdk-go/api/iam/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type apiKeyResponse struct {
	APIKey     *iam.APIKey
	EntityType string              `json:"entity_type"`
	Policies   map[string][]string `json:"policies"`
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

type userEntity struct {
	UserID string
}

type applicationEntity struct {
	ApplicationID string
}

type entity interface {
	entityType(ctx context.Context, api *iam.API) (string, error)
	getPolicies(ctx context.Context, api *iam.API) ([]*iam.Policy, error)
}

func (u userEntity) entityType(ctx context.Context, api *iam.API) (string, error) {
	user, err := api.GetUser(&iam.GetUserRequest{
		UserID: u.UserID,
	}, scw.WithContext(ctx))
	if err != nil {
		return "", err
	}

	return string(user.Type), nil
}

func (a applicationEntity) entityType(ctx context.Context, api *iam.API) (string, error) {
	return "application", nil
}

func buildEntity(apiKey *iam.APIKey) (entity, error) {
	if apiKey == nil {
		return nil, errors.New("invalid API key")
	}
	if apiKey.UserID != nil {
		return userEntity{UserID: *apiKey.UserID}, nil
	}
	if apiKey.ApplicationID != nil {
		return applicationEntity{ApplicationID: *apiKey.ApplicationID}, nil
	}

	return nil, errors.New("invalid API key")
}

func (u userEntity) getPolicies(ctx context.Context, api *iam.API) ([]*iam.Policy, error) {
	policies, err := api.ListPolicies(&iam.ListPoliciesRequest{
		UserIDs: []string{u.UserID},
	}, scw.WithContext(ctx), scw.WithAllPages())
	if err != nil {
		return nil, err
	}
	if policies == nil {
		return nil, errors.New("no policies found")
	}

	return policies.Policies, nil
}

func (a applicationEntity) getPolicies(ctx context.Context, api *iam.API) ([]*iam.Policy, error) {
	policies, err := api.ListPolicies(&iam.ListPoliciesRequest{
		ApplicationIDs: []string{a.ApplicationID},
	}, scw.WithContext(ctx), scw.WithAllPages())
	if err != nil {
		return nil, err
	}
	if policies == nil {
		return nil, errors.New("no policies found")
	}

	return policies.Policies, nil
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

	entity, err := buildEntity(apiKey)
	if err != nil {
		return response, err
	}

	entityType, err := entity.entityType(ctx, api)
	if err != nil {
		return response, err
	}

	response.APIKey = apiKey
	response.EntityType = entityType

	if entityType == string(iam.UserTypeOwner) {
		response.EntityType = entityType + " (owner has all permissions over the organization)"

		return response, nil
	}

	if options.WithPolicies {
		policies, err := entity.getPolicies(ctx, api)
		if err != nil {
			return response, err
		}

		// Build a map of policies -> [rules...]
		policyMap := map[string][]string{}
		for _, policy := range policies {
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

func apiKeyMarshalerFunc(i any, opt *human.MarshalOpt) (string, error) {
	type tmp apiKeyResponse
	resp := tmp(i.(apiKeyResponse))

	opt.Sections = []*human.MarshalSection{
		{
			FieldName: "EntityType",
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
	run      func(ctx context.Context, args any) (i any, e error)
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
	run: func(ctx context.Context, args any) (i any, e error) {
		arguments := args.(*iamGetAPIKeyArgs)

		client := core.ExtractClient(ctx)
		api := iam.NewAPI(client)

		return getApiKey(ctx, api, arguments.AccessKey, apiKeyOptions{
			WithPolicies: arguments.WithPolicies,
		})
	},
}
