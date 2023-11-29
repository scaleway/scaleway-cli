// This file was automatically generated. DO NOT EDIT.
// If you have any remark or suggestion do not hesitate to open an issue.

package vpc

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/internal/core"
	"github.com/scaleway/scaleway-sdk-go/api/vpc/v2"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// always import dependencies
var (
	_ = scw.RegionFrPar
)

func GetGeneratedCommands() *core.Commands {
	return core.NewCommands(
		vpcRoot(),
		vpcVpc(),
		vpcPrivateNetwork(),
		vpcSubnet(),
		vpcVpcList(),
		vpcVpcCreate(),
		vpcVpcGet(),
		vpcVpcUpdate(),
		vpcVpcDelete(),
		vpcPrivateNetworkList(),
		vpcPrivateNetworkCreate(),
		vpcPrivateNetworkGet(),
		vpcPrivateNetworkUpdate(),
		vpcPrivateNetworkDelete(),
		vpcPrivateNetworkPost(),
	)
}
func vpcRoot() *core.Command {
	return &core.Command{
		Short:     `VPC API`,
		Long:      `VPC API.`,
		Namespace: "vpc",
	}
}

func vpcVpc() *core.Command {
	return &core.Command{
		Short: `VPC management command`,
		Long: `A Virtual Private Cloud (VPC) allows you to group your regional 
Private Networks together. Note that a Private Network can be a 
part of only one VPC.`,
		Namespace: "vpc",
		Resource:  "vpc",
	}
}

func vpcPrivateNetwork() *core.Command {
	return &core.Command{
		Short: `Private network management command`,
		Long: `A Private Network allows you to interconnect your Scaleway resources
in an isolated and private network. Network reachability is limited
to resources that are on the same Private Network. Note that a
resource can be a part of multiple private networks.`,
		Namespace: "vpc",
		Resource:  "private-network",
	}
}

func vpcSubnet() *core.Command {
	return &core.Command{
		Short:     `Subnet management command`,
		Long:      `CIDR Subnet.`,
		Namespace: "vpc",
		Resource:  "subnet",
	}
}

func vpcVpcList() *core.Command {
	return &core.Command{
		Short:     `List VPCs`,
		Long:      `List existing VPCs in the specified region.`,
		Namespace: "vpc",
		Resource:  "vpc",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(vpc.ListVPCsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `Sort order of the returned VPCs`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"created_at_asc", "created_at_desc", "name_asc", "name_desc"},
			},
			{
				Name:       "name",
				Short:      `Name to filter for. Only VPCs with names containing this string will be returned`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `Tags to filter for. Only VPCs with one more more matching tags will be returned`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "project-id",
				Short:      `Project ID to filter for. Only VPCs belonging to this Project will be returned`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "is-default",
				Short:      `Defines whether to filter only for VPCs which are the default one for their Project`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "organization-id",
				Short:      `Organization ID to filter for. Only VPCs belonging to this Organization will be returned`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw, scw.Region(core.AllLocalities)),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*vpc.ListVPCsRequest)

			client := core.ExtractClient(ctx)
			api := vpc.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			if request.Region == scw.Region(core.AllLocalities) {
				opts = append(opts, scw.WithRegions(api.Regions()...))
				request.Region = ""
			}
			resp, err := api.ListVPCs(request, opts...)
			if err != nil {
				return nil, err
			}
			return resp.Vpcs, nil

		},
	}
}

func vpcVpcCreate() *core.Command {
	return &core.Command{
		Short:     `Create a VPC`,
		Long:      `Create a new VPC in the specified region.`,
		Namespace: "vpc",
		Resource:  "vpc",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(vpc.CreateVPCRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "name",
				Short:      `Name for the VPC`,
				Required:   true,
				Deprecated: false,
				Positional: false,
				Default:    core.RandomValueGenerator("vpc"),
			},
			core.ProjectIDArgSpec(),
			{
				Name:       "tags.{index}",
				Short:      `Tags for the VPC`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*vpc.CreateVPCRequest)

			client := core.ExtractClient(ctx)
			api := vpc.NewAPI(client)
			return api.CreateVPC(request)

		},
	}
}

func vpcVpcGet() *core.Command {
	return &core.Command{
		Short:     `Get a VPC`,
		Long:      `Retrieve details of an existing VPC, specified by its VPC ID.`,
		Namespace: "vpc",
		Resource:  "vpc",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(vpc.GetVPCRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "vpc-id",
				Short:      `VPC ID`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*vpc.GetVPCRequest)

			client := core.ExtractClient(ctx)
			api := vpc.NewAPI(client)
			return api.GetVPC(request)

		},
	}
}

func vpcVpcUpdate() *core.Command {
	return &core.Command{
		Short:     `Update VPC`,
		Long:      `Update parameters including name and tags of the specified VPC.`,
		Namespace: "vpc",
		Resource:  "vpc",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(vpc.UpdateVPCRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "vpc-id",
				Short:      `VPC ID`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "name",
				Short:      `Name for the VPC`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `Tags for the VPC`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*vpc.UpdateVPCRequest)

			client := core.ExtractClient(ctx)
			api := vpc.NewAPI(client)
			return api.UpdateVPC(request)

		},
	}
}

func vpcVpcDelete() *core.Command {
	return &core.Command{
		Short:     `Delete a VPC`,
		Long:      `Delete a VPC specified by its VPC ID.`,
		Namespace: "vpc",
		Resource:  "vpc",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(vpc.DeleteVPCRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "vpc-id",
				Short:      `VPC ID`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*vpc.DeleteVPCRequest)

			client := core.ExtractClient(ctx)
			api := vpc.NewAPI(client)
			e = api.DeleteVPC(request)
			if e != nil {
				return nil, e
			}
			return &core.SuccessResult{
				Resource: "vpc",
				Verb:     "delete",
			}, nil
		},
	}
}

func vpcPrivateNetworkList() *core.Command {
	return &core.Command{
		Short:     `List Private Networks`,
		Long:      `List existing Private Networks in the specified region. By default, the Private Networks returned in the list are ordered by creation date in ascending order, though this can be modified via the order_by field.`,
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
				Name:       "vpc-id",
				Short:      `VPC ID to filter for. Only Private Networks belonging to this VPC will be returned`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "dhcp-enabled",
				Short:      `DHCP status to filter for. When true, only Private Networks with managed DHCP enabled will be returned`,
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
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw, scw.Region(core.AllLocalities)),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*vpc.ListPrivateNetworksRequest)

			client := core.ExtractClient(ctx)
			api := vpc.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			if request.Region == scw.Region(core.AllLocalities) {
				opts = append(opts, scw.WithRegions(api.Regions()...))
				request.Region = ""
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
		Long:      `Create a new Private Network. Once created, you can attach Scaleway resources which are in the same region.`,
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
			{
				Name:       "vpc-id",
				Short:      `VPC in which to create the Private Network`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
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
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
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
				Short:      `Name for the Private Network`,
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
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
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
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
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

func vpcPrivateNetworkPost() *core.Command {
	return &core.Command{
		Short:     `Migrate Private Networks from zoned to regional`,
		Long:      `Transform multiple existing zoned Private Networks (scoped to a single Availability Zone) into regional Private Networks, scoped to an entire region. You can transform one or many Private Networks (specified by their Private Network IDs) within a single Scaleway Organization or Project, with the same call.`,
		Namespace: "vpc",
		Resource:  "private-network",
		Verb:      "post",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(vpc.MigrateZonalPrivateNetworksRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ProjectIDArgSpec(),
			{
				Name:       "private-network-ids.{index}",
				Short:      `IDs of the Private Networks to migrate`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.OrganizationIDArgSpec(),
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*vpc.MigrateZonalPrivateNetworksRequest)

			client := core.ExtractClient(ctx)
			api := vpc.NewAPI(client)
			e = api.MigrateZonalPrivateNetworks(request)
			if e != nil {
				return nil, e
			}
			return &core.SuccessResult{
				Resource: "private-network",
				Verb:     "post",
			}, nil
		},
	}
}
