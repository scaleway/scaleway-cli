// This file was automatically generated. DO NOT EDIT.
// If you have any remark or suggestion do not hesitate to open an issue.

package vpc

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-sdk-go/api/vpc/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// always import dependencies
var (
	_ = scw.RegionFrPar
)

func GetGeneratedCommands() *core.Commands {
	return core.NewCommands(
		vpcRoot(),
		vpcPrivateNetwork(),
		vpcPrivateNetworkList(),
		vpcPrivateNetworkCreate(),
		vpcPrivateNetworkGet(),
		vpcPrivateNetworkUpdate(),
		vpcPrivateNetworkDelete(),
	)
}
func vpcRoot() *core.Command {
	return &core.Command{
		Short:     `VPC API`,
		Long:      ``,
		Namespace: "vpc",
	}
}

func vpcPrivateNetwork() *core.Command {
	return &core.Command{
		Short: `Private network management command`,
		Long: `A private network allows interconnecting your instances in an
isolated and private network. The network reachability is limited
to the instances that are on the same private network.  Network
Interface Controllers (NICs) are available on the instance and can
be freely managed (adding IP addresses, shutdown interface...)

Note that an instance can be a part of multiple private networks.
`,
		Namespace: "vpc",
		Resource:  "private-network",
	}
}

func vpcPrivateNetworkList() *core.Command {
	return &core.Command{
		Short:     `List private networks`,
		Long:      `List private networks.`,
		Namespace: "vpc",
		Resource:  "private-network",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(vpc.ListPrivateNetworksRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `The sort order of the returned private networks`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"created_at_asc", "created_at_desc", "name_asc", "name_desc"},
			},
			{
				Name:       "name",
				Short:      `Filter private networks with names containing this string`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `Filter private networks with one or more matching tags`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "project-id",
				Short:      `The project ID on which to filter the returned private networks`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "organization-id",
				Short:      `The organization ID on which to filter the returned private networks`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneNlAms1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*vpc.ListPrivateNetworksRequest)

			client := core.ExtractClient(ctx)
			api := vpc.NewAPI(client)
			resp, err := api.ListPrivateNetworks(request, scw.WithAllPages())
			if err != nil {
				return nil, err
			}
			return resp.PrivateNetworks, nil

		},
	}
}

func vpcPrivateNetworkCreate() *core.Command {
	return &core.Command{
		Short:     `Create a private network`,
		Long:      `Create a private network.`,
		Namespace: "vpc",
		Resource:  "private-network",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(vpc.CreatePrivateNetworkRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "name",
				Short:      `The name of the private network`,
				Required:   true,
				Deprecated: false,
				Positional: false,
				Default:    core.RandomValueGenerator("pn"),
			},
			core.ProjectIDArgSpec(),
			{
				Name:       "tags.{index}",
				Short:      `The private networks tags`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneNlAms1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*vpc.CreatePrivateNetworkRequest)

			client := core.ExtractClient(ctx)
			api := vpc.NewAPI(client)
			return api.CreatePrivateNetwork(request)

		},
	}
}

func vpcPrivateNetworkGet() *core.Command {
	return &core.Command{
		Short:     `Get a private network`,
		Long:      `Get a private network.`,
		Namespace: "vpc",
		Resource:  "private-network",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(vpc.GetPrivateNetworkRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "private-network-id",
				Short:      `The private network id`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneNlAms1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*vpc.GetPrivateNetworkRequest)

			client := core.ExtractClient(ctx)
			api := vpc.NewAPI(client)
			return api.GetPrivateNetwork(request)

		},
	}
}

func vpcPrivateNetworkUpdate() *core.Command {
	return &core.Command{
		Short:     `Update private network`,
		Long:      `Update private network.`,
		Namespace: "vpc",
		Resource:  "private-network",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(vpc.UpdatePrivateNetworkRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "private-network-id",
				Short:      `The private network ID`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "name",
				Short:      `The name of the private network`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `The private networks tags`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneNlAms1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*vpc.UpdatePrivateNetworkRequest)

			client := core.ExtractClient(ctx)
			api := vpc.NewAPI(client)
			return api.UpdatePrivateNetwork(request)

		},
	}
}

func vpcPrivateNetworkDelete() *core.Command {
	return &core.Command{
		Short:     `Delete a private network`,
		Long:      `Delete a private network.`,
		Namespace: "vpc",
		Resource:  "private-network",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(vpc.DeletePrivateNetworkRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "private-network-id",
				Short:      `The private network ID`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneNlAms1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*vpc.DeletePrivateNetworkRequest)

			client := core.ExtractClient(ctx)
			api := vpc.NewAPI(client)
			e = api.DeletePrivateNetwork(request)
			if e != nil {
				return nil, e
			}
			return &core.SuccessResult{
				Resource: "private-network",
				Verb:     "delete",
			}, nil
		},
	}
}
