// This file was automatically generated. DO NOT EDIT.
// If you have any remark or suggestion do not hesitate to open an issue.

package vpc

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/internal/core"
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
		Long:      `VPC API.`,
		Namespace: "vpc",
	}
}

func vpcPrivateNetwork() *core.Command {
	return &core.Command{
		Short: `Private network management command`,
		Long: `A Private Network allows you to interconnect your Scaleway resources in an
isolated and private network. Network reachability is limited
to resources that are on the same Private Network. Note that a resource can
be part of multiple Private Networks.
`,
		Namespace: "vpc",
		Resource:  "private-network",
	}
}

func vpcPrivateNetworkList() *core.Command {
	return &core.Command{
		Short:     `List Private Networks`,
		Long:      `List existing Private Networks in a specified Availability Zone. By default, the Private Networks returned in the list are ordered by creation date in ascending order, though this can be modified via the order_by field.`,
		Namespace: "vpc",
		Resource:  "private-network",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(vpc.ListPrivateNetworksRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `Sort order of the returned Private Networks`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"created_at_asc", "created_at_desc", "name_asc", "name_desc"},
			},
			{
				Name:       "name",
				Short:      `Name to filter for. Only Private Networks with names containing this string will be returned`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `Tags to filter for. Only Private Networks with one or more matching tags will be returned`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "project-id",
				Short:      `Project ID to filter for. Only Private Networks belonging to this Project will be returned`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "private-network-ids.{index}",
				Short:      `Private Network IDs to filter for. Only Private Networks with one of these IDs will be returned`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "include-regional",
				Short:      `Defines whether to include regional Private Networks in the response`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "organization-id",
				Short:      `Organization ID to filter for. Only Private Networks belonging to this Organization will be returned`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneFrPar3, scw.ZoneNlAms1, scw.ZoneNlAms2, scw.ZoneNlAms3, scw.ZonePlWaw1, scw.ZonePlWaw2, scw.Zone(core.AllLocalities)),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*vpc.ListPrivateNetworksRequest)

			client := core.ExtractClient(ctx)
			api := vpc.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			if request.Zone == scw.Zone(core.AllLocalities) {
				opts = append(opts, scw.WithZones(api.Zones()...))
				request.Zone = ""
			}
			resp, err := api.ListPrivateNetworks(request, opts...)
			if err != nil {
				return nil, err
			}
			return resp.PrivateNetworks, nil

		},
	}
}

func vpcPrivateNetworkCreate() *core.Command {
	return &core.Command{
		Short:     `Create a Private Network`,
		Long:      `Create a new Private Network. Once created, you can attach Scaleway resources in the same Availability Zone.`,
		Namespace: "vpc",
		Resource:  "private-network",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(vpc.CreatePrivateNetworkRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "name",
				Short:      `Name for the Private Network`,
				Required:   true,
				Deprecated: false,
				Positional: false,
				Default:    core.RandomValueGenerator("pn"),
			},
			core.ProjectIDArgSpec(),
			{
				Name:       "tags.{index}",
				Short:      `Tags for the Private Network`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "subnets.{index}",
				Short:      `Private Network subnets CIDR`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneFrPar3, scw.ZoneNlAms1, scw.ZoneNlAms2, scw.ZoneNlAms3, scw.ZonePlWaw1, scw.ZonePlWaw2),
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
		Short:     `Get a Private Network`,
		Long:      `Retrieve information about an existing Private Network, specified by its Private Network ID. Its full details are returned in the response object.`,
		Namespace: "vpc",
		Resource:  "private-network",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(vpc.GetPrivateNetworkRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "private-network-id",
				Short:      `Private Network ID`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneFrPar3, scw.ZoneNlAms1, scw.ZoneNlAms2, scw.ZoneNlAms3, scw.ZonePlWaw1, scw.ZonePlWaw2),
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
		Short:     `Update Private Network`,
		Long:      `Update parameters (such as name or tags) of an existing Private Network, specified by its Private Network ID.`,
		Namespace: "vpc",
		Resource:  "private-network",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(vpc.UpdatePrivateNetworkRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "private-network-id",
				Short:      `Private Network ID`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "name",
				Short:      `Name of the private network`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `Tags for the Private Network`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "subnets.{index}",
				Short:      `Private Network subnets CIDR (deprecated)`,
				Required:   false,
				Deprecated: true,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneFrPar3, scw.ZoneNlAms1, scw.ZoneNlAms2, scw.ZoneNlAms3, scw.ZonePlWaw1, scw.ZonePlWaw2),
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
		Short:     `Delete a Private Network`,
		Long:      `Delete an existing Private Network. Note that you must first detach all resources from the network, in order to delete it.`,
		Namespace: "vpc",
		Resource:  "private-network",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(vpc.DeletePrivateNetworkRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "private-network-id",
				Short:      `Private Network ID`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneFrPar3, scw.ZoneNlAms1, scw.ZoneNlAms2, scw.ZoneNlAms3, scw.ZonePlWaw1, scw.ZonePlWaw2),
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
