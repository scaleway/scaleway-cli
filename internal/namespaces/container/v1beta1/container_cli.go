// This file was automatically generated. DO NOT EDIT.
// If you have any remark or suggestion do not hesitate to open an issue.

package container

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	container "github.com/scaleway/scaleway-sdk-go/api/container/v1beta1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// always import dependencies
var (
	_ = scw.RegionFrPar
)

func GetGeneratedCommands() *core.Commands {
	return core.NewCommands(
		containerRoot(),
		containerToken(),
		containerTokenCreate(),
		containerTokenGet(),
		containerTokenList(),
		containerTokenDelete(),
	)
}

func containerRoot() *core.Command {
	return &core.Command{
		Short:     `This API allows you to manage your Serverless Containers`,
		Long:      `This API allows you to manage your Serverless Containers.`,
		Namespace: "container",
	}
}

func containerToken() *core.Command {
	return &core.Command{
		Short:     `Token management commands`,
		Long:      `Token management commands.`,
		Namespace: "container",
		Resource:  "token",
	}
}

func containerTokenCreate() *core.Command {
	return &core.Command{
		Short:     `Create a new revocable token`,
		Long:      `Deprecated in favor of IAM authentication.`,
		Namespace: "container",
		Resource:  "token",
		Verb:      "create",
		// Deprecated:    true,
		ArgsType: reflect.TypeOf(container.CreateTokenRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "container-id",
				Short:      `UUID of the container to create the token for`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "namespace-id",
				Short:      `UUID of the namespace to create the token for`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "description",
				Short:      `Description of the token`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "expires-at",
				Short:      `Expiry date of the token`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*container.CreateTokenRequest)

			client := core.ExtractClient(ctx)
			api := container.NewAPI(client)

			return api.CreateToken(request)
		},
	}
}

func containerTokenGet() *core.Command {
	return &core.Command{
		Short:     `Get a token`,
		Long:      `Get a token with a specified ID.`,
		Namespace: "container",
		Resource:  "token",
		Verb:      "get",
		// Deprecated:    true,
		ArgsType: reflect.TypeOf(container.GetTokenRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "token-id",
				Short:      `UUID of the token to get`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*container.GetTokenRequest)

			client := core.ExtractClient(ctx)
			api := container.NewAPI(client)

			return api.GetToken(request)
		},
	}
}

func containerTokenList() *core.Command {
	return &core.Command{
		Short:     `List all tokens`,
		Long:      `List all tokens belonging to a specified Organization or Project.`,
		Namespace: "container",
		Resource:  "token",
		Verb:      "list",
		// Deprecated:    true,
		ArgsType: reflect.TypeOf(container.ListTokensRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `Order of the tokens`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"created_at_asc",
					"created_at_desc",
				},
			},
			{
				Name:       "container-id",
				Short:      `UUID of the container the token belongs to`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "namespace-id",
				Short:      `UUID of the namespace the token belongs to`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
				scw.Region(core.AllLocalities),
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*container.ListTokensRequest)

			client := core.ExtractClient(ctx)
			api := container.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			if request.Region == scw.Region(core.AllLocalities) {
				opts = append(opts, scw.WithRegions(api.Regions()...))
				request.Region = ""
			}
			resp, err := api.ListTokens(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.Tokens, nil
		},
	}
}

func containerTokenDelete() *core.Command {
	return &core.Command{
		Short:     `Delete a token`,
		Long:      `Delete a token with a specified ID.`,
		Namespace: "container",
		Resource:  "token",
		Verb:      "delete",
		// Deprecated:    true,
		ArgsType: reflect.TypeOf(container.DeleteTokenRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "token-id",
				Short:      `UUID of the token to delete`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*container.DeleteTokenRequest)

			client := core.ExtractClient(ctx)
			api := container.NewAPI(client)

			return api.DeleteToken(request)
		},
	}
}
