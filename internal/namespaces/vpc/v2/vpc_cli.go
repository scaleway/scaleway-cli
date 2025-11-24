// This file was automatically generated. DO NOT EDIT.
// If you have any remark or suggestion do not hesitate to open an issue.

package vpc

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
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
		vpcRoute(),
		vpcRule(),
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
		vpcPrivateNetworkEnableDHCP(),
		vpcRouteEnableRouting(),
		vpcRouteCreate(),
		vpcRouteGet(),
		vpcRouteUpdate(),
		vpcRouteDelete(),
		vpcRuleGet(),
		vpcRuleSet(),
		vpcRouteList(),
	)
}

func vpcRoot() *core.Command {
	return &core.Command{
		Short:     `This API allows you to manage your Virtual Private Clouds (VPCs) and Private Networks`,
		Long:      `This API allows you to manage your Virtual Private Clouds (VPCs) and Private Networks.`,
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

func vpcRoute() *core.Command {
	return &core.Command{
		Short:     `Route management command`,
		Long:      `Custom routes.`,
		Namespace: "vpc",
		Resource:  "route",
	}
}

func vpcRule() *core.Command {
	return &core.Command{
		Short:     `Rule management command`,
		Long:      `ACL Rules.`,
		Namespace: "vpc",
		Resource:  "rule",
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
				EnumValues: []string{
					"created_at_asc",
					"created_at_desc",
					"name_asc",
					"name_desc",
				},
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
				Short:      `Tags to filter for. Only VPCs with one or more matching tags will be returned`,
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
				Name:       "routing-enabled",
				Short:      `Defines whether to filter only for VPCs which route traffic between their Private Networks`,
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
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
				scw.Region(core.AllLocalities),
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
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
			{
				Name:       "enable-routing",
				Short:      `Enable routing between Private Networks in the VPC`,
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
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
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
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
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
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
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
				EnumValues: []string{
					"created_at_asc",
					"created_at_desc",
					"name_asc",
					"name_desc",
				},
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
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
				scw.Region(core.AllLocalities),
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
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
			{
				Name:       "default-route-propagation-enabled",
				Short:      `Defines whether default v4 and v6 routes are propagated for this Private Network`,
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
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
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
			{
				Name:       "default-route-propagation-enabled",
				Short:      `Defines whether default v4 and v6 routes are propagated for this Private Network`,
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
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
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

func vpcPrivateNetworkEnableDHCP() *core.Command {
	return &core.Command{
		Short:     `Enable DHCP on a Private Network`,
		Long:      `Enable DHCP managed on an existing Private Network. Note that you will not be able to deactivate it afterwards.`,
		Namespace: "vpc",
		Resource:  "private-network",
		Verb:      "enable-dhcp",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(vpc.EnableDHCPRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "private-network-id",
				Short:      `Private Network ID`,
				Required:   true,
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
			request := args.(*vpc.EnableDHCPRequest)

			client := core.ExtractClient(ctx)
			api := vpc.NewAPI(client)

			return api.EnableDHCP(request)
		},
	}
}

func vpcRouteEnableRouting() *core.Command {
	return &core.Command{
		Short:     `Enable routing on a VPC`,
		Long:      `Enable routing on an existing VPC. Note that you will not be able to deactivate it afterwards.`,
		Namespace: "vpc",
		Resource:  "route",
		Verb:      "enable-routing",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(vpc.EnableRoutingRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "vpc-id",
				Short:      `VPC ID`,
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
			request := args.(*vpc.EnableRoutingRequest)

			client := core.ExtractClient(ctx)
			api := vpc.NewAPI(client)

			return api.EnableRouting(request)
		},
	}
}

func vpcRouteCreate() *core.Command {
	return &core.Command{
		Short:     `Create a Route`,
		Long:      `Create a new custom Route.`,
		Namespace: "vpc",
		Resource:  "route",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(vpc.CreateRouteRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "description",
				Short:      `Route description`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `Tags of the Route`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "vpc-id",
				Short:      `VPC the Route belongs to`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "destination",
				Short:      `Destination of the Route`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "nexthop-resource-id",
				Short:      `ID of the nexthop resource`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "nexthop-private-network-id",
				Short:      `ID of the nexthop private network`,
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
			request := args.(*vpc.CreateRouteRequest)

			client := core.ExtractClient(ctx)
			api := vpc.NewAPI(client)

			return api.CreateRoute(request)
		},
	}
}

func vpcRouteGet() *core.Command {
	return &core.Command{
		Short:     `Get a Route`,
		Long:      `Retrieve details of an existing Route, specified by its Route ID.`,
		Namespace: "vpc",
		Resource:  "route",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(vpc.GetRouteRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "route-id",
				Short:      `Route ID`,
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
			request := args.(*vpc.GetRouteRequest)

			client := core.ExtractClient(ctx)
			api := vpc.NewAPI(client)

			return api.GetRoute(request)
		},
	}
}

func vpcRouteUpdate() *core.Command {
	return &core.Command{
		Short:     `Update Route`,
		Long:      `Update parameters of the specified Route.`,
		Namespace: "vpc",
		Resource:  "route",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(vpc.UpdateRouteRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "route-id",
				Short:      `Route ID`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "description",
				Short:      `Route description`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `Tags of the Route`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "destination",
				Short:      `Destination of the Route`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "nexthop-resource-id",
				Short:      `ID of the nexthop resource`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "nexthop-private-network-id",
				Short:      `ID of the nexthop private network`,
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
			request := args.(*vpc.UpdateRouteRequest)

			client := core.ExtractClient(ctx)
			api := vpc.NewAPI(client)

			return api.UpdateRoute(request)
		},
	}
}

func vpcRouteDelete() *core.Command {
	return &core.Command{
		Short:     `Delete a Route`,
		Long:      `Delete a Route specified by its Route ID.`,
		Namespace: "vpc",
		Resource:  "route",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(vpc.DeleteRouteRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "route-id",
				Short:      `Route ID`,
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
			request := args.(*vpc.DeleteRouteRequest)

			client := core.ExtractClient(ctx)
			api := vpc.NewAPI(client)
			e = api.DeleteRoute(request)
			if e != nil {
				return nil, e
			}

			return &core.SuccessResult{
				Resource: "route",
				Verb:     "delete",
			}, nil
		},
	}
}

func vpcRuleGet() *core.Command {
	return &core.Command{
		Short:     `Get ACL Rules for VPC`,
		Long:      `Retrieve a list of ACL rules for a VPC, specified by its VPC ID.`,
		Namespace: "vpc",
		Resource:  "rule",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(vpc.GetACLRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "vpc-id",
				Short:      `ID of the Network ACL's VPC`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "is-ipv6",
				Short:      `Defines whether this set of ACL rules is for IPv6 (false = IPv4). Each Network ACL can have rules for only one IP type.`,
				Required:   true,
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
			request := args.(*vpc.GetACLRequest)

			client := core.ExtractClient(ctx)
			api := vpc.NewAPI(client)

			return api.GetACL(request)
		},
	}
}

func vpcRuleSet() *core.Command {
	return &core.Command{
		Short:     `Set VPC ACL rules`,
		Long:      `Set the list of ACL rules and the default routing policy for a VPC.`,
		Namespace: "vpc",
		Resource:  "rule",
		Verb:      "set",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(vpc.SetACLRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "vpc-id",
				Short:      `ID of the Network ACL's VPC`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "rules.{index}.protocol",
				Short:      `Protocol to which this rule applies`,
				Required:   true,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"ANY",
					"TCP",
					"UDP",
					"ICMP",
				},
			},
			{
				Name:       "rules.{index}.source",
				Short:      `Source IP range to which this rule applies (CIDR notation with subnet mask)`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "rules.{index}.src-port-low",
				Short:      `Starting port of the source port range to which this rule applies (inclusive)`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "rules.{index}.src-port-high",
				Short:      `Ending port of the source port range to which this rule applies (inclusive)`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "rules.{index}.destination",
				Short:      `Destination IP range to which this rule applies (CIDR notation with subnet mask)`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "rules.{index}.dst-port-low",
				Short:      `Starting port of the destination port range to which this rule applies (inclusive)`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "rules.{index}.dst-port-high",
				Short:      `Ending port of the destination port range to which this rule applies (inclusive)`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "rules.{index}.action",
				Short:      `Policy to apply to the packet`,
				Required:   true,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_action",
					"accept",
					"drop",
				},
			},
			{
				Name:       "rules.{index}.description",
				Short:      `Rule description`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "is-ipv6",
				Short:      `Defines whether this set of ACL rules is for IPv6 (false = IPv4). Each Network ACL can have rules for only one IP type.`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "default-policy",
				Short:      `Action to take for packets which do not match any rules`,
				Required:   true,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_action",
					"accept",
					"drop",
				},
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*vpc.SetACLRequest)

			client := core.ExtractClient(ctx)
			api := vpc.NewAPI(client)

			return api.SetACL(request)
		},
	}
}

func vpcRouteList() *core.Command {
	return &core.Command{
		Short:     `Return routes with associated next hop data`,
		Long:      `Return routes with associated next hop data.`,
		Namespace: "vpc",
		Resource:  "route",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(vpc.RoutesWithNexthopAPIListRoutesWithNexthopRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `Sort order of the returned routes`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"created_at_asc",
					"created_at_desc",
					"destination_asc",
					"destination_desc",
					"prefix_len_asc",
					"prefix_len_desc",
				},
			},
			{
				Name:       "vpc-id",
				Short:      `VPC to filter for. Only routes within this VPC will be returned`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "nexthop-resource-id",
				Short:      `Next hop resource ID to filter for. Only routes with a matching next hop resource ID will be returned`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "nexthop-private-network-id",
				Short:      `Next hop private network ID to filter for. Only routes with a matching next hop private network ID will be returned`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "nexthop-resource-type",
				Short:      `Next hop resource type to filter for. Only Routes with a matching next hop resource type will be returned`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_type",
					"vpc_gateway_network",
					"instance_private_nic",
					"baremetal_private_nic",
					"apple_silicon_private_nic",
				},
			},
			{
				Name:       "contains",
				Short:      `Only routes whose destination is contained in this subnet will be returned`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `Tags to filter for, only routes with one or more matching tags will be returned`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "is-ipv6",
				Short:      `Only routes with an IPv6 destination will be returned`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.Region(core.AllLocalities)),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*vpc.RoutesWithNexthopAPIListRoutesWithNexthopRequest)

			client := core.ExtractClient(ctx)
			api := vpc.NewRoutesWithNexthopAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			resp, err := api.ListRoutesWithNexthop(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.Routes, nil
		},
	}
}
