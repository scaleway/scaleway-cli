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
		ArgsType:  reflect.TypeOf(vpc.ListPrivateNetworksRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `The sort order of the returned private networks`,
				Required:   false,
				Positional: false,
				EnumValues: []string{"created_at_asc", "created_at_desc", "name_asc", "name_desc"},
			},
			{
				Name:       "name",
				Short:      `Filter private networks with names containing this string`,
				Required:   false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `Filter private networks with one or more matching tags`,
				Required:   false,
				Positional: false,
			},
			{
				Name:       "project-id",
				Short:      `The project ID on which to filter the returned private networks`,
				Required:   false,
				Positional: false,
			},
			{
				Name:       "organization-id",
				Short:      `The organization ID on which to filter the returned private networks`,
				Required:   false,
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
		ArgsType:  reflect.TypeOf(vpc.CreatePrivateNetworkRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "name",
				Short:      `The name of the private network`,
				Required:   true,
				Positional: false,
				Default:    core.RandomValueGenerator("pn"),
			},
			{
				Name:       "project-id",
				Short:      `The project ID of the private network`,
				Required:   true,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `The private networks tags`,
				Required:   false,
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

func vpcPrivateNetworkUpdate() *core.Command {
	return &core.Command{
		Short:     `Update private network`,
		Long:      `Update private network.`,
		Namespace: "vpc",
		Resource:  "private-network",
		Verb:      "update",
		ArgsType:  reflect.TypeOf(vpc.UpdatePrivateNetworkRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "private-network-id",
				Short:      `The private network ID`,
				Required:   true,
				Positional: false,
			},
			{
				Name:       "name",
				Short:      `The name of the private network`,
				Required:   false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `The private networks tags`,
				Required:   false,
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
		ArgsType:  reflect.TypeOf(vpc.DeletePrivateNetworkRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "private-network-id",
				Short:      `The private network ID`,
				Required:   true,
				Positional: false,
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
