package iam

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	iam "github.com/scaleway/scaleway-sdk-go/api/iam/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

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
	response.UserType = user.Type

	if options.WithPolicies {
		listPolicyRequest := &iam.ListPoliciesRequest{
			UserIDs: []string{*apiKey.UserID},
		}
		// if user is owner, list all policies attached to the organization
		// because the user has no policies attached directly
		if user.Type == iam.UserTypeOwner {
			listPolicyRequest.OrganizationID = apiKey.DefaultProjectID
			listPolicyRequest.UserIDs = []string{}
		}
		policies, err := api.ListPolicies(
			listPolicyRequest,
			scw.WithAllPages(),
			scw.WithContext(ctx),
		)
		if err != nil {
			return response, err
		}
		response.Policies = policies.Policies
	}

	return response, nil
}

type apiKeyResponse struct {
	APIKey   *iam.APIKey
	UserType iam.UserType  `json:"user_type"`
	Policies []*iam.Policy `json:"policies"`
}

func apiKeyMarshalerFunc(i interface{}, opt *human.MarshalOpt) (string, error) {
	type tmp apiKeyResponse
	resp := tmp(i.(apiKeyResponse))

	sections := []*human.MarshalSection{
		{
			FieldName: "UserType",
			Title:     "User Type",
		},
		{
			FieldName: "APIKey",
			Title:     "API Key",
		},
	}

	if len(resp.Policies) > 0 {
		sections = append(sections, &human.MarshalSection{
			FieldName: "Policies",
			Title:     "Policies",
		})
	}

	opt.Sections = sections

	return human.Marshal(resp, opt)
}

func iamAPIKeyGetBuilder(c *core.Command) *core.Command {
	human.RegisterMarshalerFunc(apiKeyResponse{}, apiKeyMarshalerFunc)

	return &core.Command{
		Short:     `Get an API key`,
		Long:      `Retrieve information about an API key, specified by the ` + "`" + `access_key` + "`" + ` parameter. The API key's details, including either the ` + "`" + `user_id` + "`" + ` or ` + "`" + `application_id` + "`" + ` of its bearer are returned in the response. Note that the string value for the ` + "`" + `secret_key` + "`" + ` is nullable, and therefore is not displayed in the response. The ` + "`" + `secret_key` + "`" + ` value is only displayed upon API key creation.`,
		Namespace: "iam",
		Resource:  "api-key",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(iamGetAPIKeyArgs{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "access-key",
				Short:      `Access key to search for`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "with-policies",
				Short:      `Display policies associated with the API key`,
				Default:    core.DefaultValueSetter("false"),
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			arguments := args.(*iamGetAPIKeyArgs)

			client := core.ExtractClient(ctx)
			api := iam.NewAPI(client)

			return getApiKey(ctx, api, arguments.AccessKey, apiKeyOptions{
				WithPolicies: arguments.WithPolicies,
			})
		},
	}
}
