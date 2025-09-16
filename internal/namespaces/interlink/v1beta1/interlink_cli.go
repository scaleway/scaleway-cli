// This file was automatically generated. DO NOT EDIT.
// If you have any remark or suggestion do not hesitate to open an issue.

package interlink

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	interlink "github.com/scaleway/scaleway-sdk-go/api/interlink/v1beta1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// always import dependencies
var (
	_ = scw.RegionFrPar
)

func GetGeneratedCommands() *core.Commands {
	return core.NewCommands(
		interlinkRoot(),
		interlinkPartner(),
		interlinkPop(),
		interlinkConnection(),
		interlinkLink(),
		interlinkRoutingPolicy(),
		interlinkPartnerList(),
		interlinkPartnerGet(),
		interlinkPopList(),
		interlinkPopGet(),
		interlinkLinkList(),
		interlinkLinkGet(),
		interlinkLinkCreate(),
		interlinkLinkUpdate(),
		interlinkLinkDelete(),
		interlinkLinkAttachVpc(),
		interlinkLinkDetachVpc(),
		interlinkLinkAttachPolicy(),
		interlinkLinkDetachPolicy(),
		interlinkLinkEnablePropagation(),
		interlinkLinkDisablePropagation(),
		interlinkRoutingPolicyList(),
		interlinkRoutingPolicyGet(),
		interlinkRoutingPolicyCreate(),
		interlinkRoutingPolicyUpdate(),
		interlinkRoutingPolicyDelete(),
	)
}

func interlinkRoot() *core.Command {
	return &core.Command{
		Short:     `This API allows you to manage your InterLink services`,
		Long:      `This API allows you to manage your InterLink services.`,
		Namespace: "interlink",
	}
}

func interlinkPartner() *core.Command {
	return &core.Command{
		Short:     `Partner commands`,
		Long:      `Partner commands.`,
		Namespace: "interlink",
		Resource:  "partner",
	}
}

func interlinkPop() *core.Command {
	return &core.Command{
		Short:     `PoP commands`,
		Long:      `PoP commands.`,
		Namespace: "interlink",
		Resource:  "pop",
	}
}

func interlinkConnection() *core.Command {
	return &core.Command{
		Short:     `Connection commands`,
		Long:      `Connection commands.`,
		Namespace: "interlink",
		Resource:  "connection",
	}
}

func interlinkLink() *core.Command {
	return &core.Command{
		Short:     `Link commands`,
		Long:      `Link commands.`,
		Namespace: "interlink",
		Resource:  "link",
	}
}

func interlinkRoutingPolicy() *core.Command {
	return &core.Command{
		Short:     `Routing policy commands`,
		Long:      `Routing policy commands.`,
		Namespace: "interlink",
		Resource:  "routing-policy",
	}
}

func interlinkPartnerList() *core.Command {
	return &core.Command{
		Short:     `List available partners`,
		Long:      `List all available partners. By default, the partners returned in the list are ordered by name in ascending order, though this can be modified via the ` + "`" + `order_by` + "`" + ` field.`,
		Namespace: "interlink",
		Resource:  "partner",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(interlink.ListPartnersRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `Order in which to return results`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"name_asc",
					"name_desc",
				},
			},
			{
				Name:       "pop-ids.{index}",
				Short:      `Filter for partners present (offering a connection) in one of these PoPs`,
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
			request := args.(*interlink.ListPartnersRequest)

			client := core.ExtractClient(ctx)
			api := interlink.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			if request.Region == scw.Region(core.AllLocalities) {
				opts = append(opts, scw.WithRegions(api.Regions()...))
				request.Region = ""
			}
			resp, err := api.ListPartners(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.Partners, nil
		},
	}
}

func interlinkPartnerGet() *core.Command {
	return &core.Command{
		Short:     `Get a partner`,
		Long:      `Get a partner for the given partner IP. The response object includes information such as the partner's name, email address and portal URL.`,
		Namespace: "interlink",
		Resource:  "partner",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(interlink.GetPartnerRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "partner-id",
				Short:      `ID of partner to get`,
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
			request := args.(*interlink.GetPartnerRequest)

			client := core.ExtractClient(ctx)
			api := interlink.NewAPI(client)

			return api.GetPartner(request)
		},
	}
}

func interlinkPopList() *core.Command {
	return &core.Command{
		Short:     `List PoPs`,
		Long:      `List all available PoPs (locations) for a given region. By default, the results are returned in ascending alphabetical order by name.`,
		Namespace: "interlink",
		Resource:  "pop",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(interlink.ListPopsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `Order in which to return results`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"name_asc",
					"name_desc",
				},
			},
			{
				Name:       "name",
				Short:      `PoP name to filter for`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "hosting-provider-name",
				Short:      `Hosting provider name to filter for`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "partner-id",
				Short:      `Filter for PoPs hosting an available shared connection from this partner`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "link-bandwidth-mbps",
				Short:      `Filter for PoPs with a shared connection allowing this bandwidth size. Note that we cannot guarantee that PoPs returned will have available capacity.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "dedicated-available",
				Short:      `Filter for PoPs with a dedicated connection available for self-hosted links.`,
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
			request := args.(*interlink.ListPopsRequest)

			client := core.ExtractClient(ctx)
			api := interlink.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			if request.Region == scw.Region(core.AllLocalities) {
				opts = append(opts, scw.WithRegions(api.Regions()...))
				request.Region = ""
			}
			resp, err := api.ListPops(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.Pops, nil
		},
	}
}

func interlinkPopGet() *core.Command {
	return &core.Command{
		Short:     `Get a PoP`,
		Long:      `Get a PoP for the given PoP ID. The response object includes the PoP's name and information about its physical location.`,
		Namespace: "interlink",
		Resource:  "pop",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(interlink.GetPopRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "pop-id",
				Short:      `ID of PoP to get`,
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
			request := args.(*interlink.GetPopRequest)

			client := core.ExtractClient(ctx)
			api := interlink.NewAPI(client)

			return api.GetPop(request)
		},
	}
}

func interlinkLinkList() *core.Command {
	return &core.Command{
		Short:     `List links`,
		Long:      `List all your links (InterLink connections). A number of filters are available, including Project ID, name, tags and status.`,
		Namespace: "interlink",
		Resource:  "link",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(interlink.ListLinksRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `Order in which to return results`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"created_at_asc",
					"created_at_desc",
					"name_asc",
					"name_desc",
					"status_asc",
					"status_desc",
				},
			},
			{
				Name:       "project-id",
				Short:      `Project ID to filter for`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "name",
				Short:      `Link name to filter for`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `Tags to filter for`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "status",
				Short:      `Link status to filter for`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_link_status",
					"configuring",
					"failed",
					"requested",
					"refused",
					"expired",
					"provisioning",
					"active",
					"limited_connectivity",
					"all_down",
					"deprovisioning",
					"deleted",
					"locked",
					"ready",
				},
			},
			{
				Name:       "bgp-v4-status",
				Short:      `BGP IPv4 status to filter for`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_bgp_status",
					"up",
					"down",
					"disabled",
				},
			},
			{
				Name:       "bgp-v6-status",
				Short:      `BGP IPv6 status to filter for`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_bgp_status",
					"up",
					"down",
					"disabled",
				},
			},
			{
				Name:       "pop-id",
				Short:      `Filter for links attached to this PoP (via connections)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "bandwidth-mbps",
				Short:      `Filter for link bandwidth (in Mbps)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "partner-id",
				Short:      `Filter for links hosted by this partner`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "vpc-id",
				Short:      `Filter for links attached to this VPC`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "routing-policy-id",
				Short:      `Filter for links using this routing policy`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "pairing-key",
				Short:      `Filter for the link with this pairing_key`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "kind",
				Short:      `Filter for hosted or self-hosted links`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"hosted",
					"self_hosted",
				},
			},
			{
				Name:       "connection-id",
				Short:      `Filter for links self-hosted on this connection`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "organization-id",
				Short:      `Organization ID to filter for`,
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
			request := args.(*interlink.ListLinksRequest)

			client := core.ExtractClient(ctx)
			api := interlink.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			if request.Region == scw.Region(core.AllLocalities) {
				opts = append(opts, scw.WithRegions(api.Regions()...))
				request.Region = ""
			}
			resp, err := api.ListLinks(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.Links, nil
		},
	}
}

func interlinkLinkGet() *core.Command {
	return &core.Command{
		Short:     `Get a link`,
		Long:      `Get a link (InterLink session / logical InterLink resource) for the given link ID. The response object includes information about the link's various configuration details.`,
		Namespace: "interlink",
		Resource:  "link",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(interlink.GetLinkRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "link-id",
				Short:      `ID of the link to get`,
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
			request := args.(*interlink.GetLinkRequest)

			client := core.ExtractClient(ctx)
			api := interlink.NewAPI(client)

			return api.GetLink(request)
		},
	}
}

func interlinkLinkCreate() *core.Command {
	return &core.Command{
		Short:     `Create a link`,
		Long:      `Create a link (InterLink session / logical InterLink resource) in a given PoP, specifying its various configuration details. Links can either be hosted (facilitated by partners' shared physical connections) or self-hosted (for users who have purchased a dedicated physical connection).`,
		Namespace: "interlink",
		Resource:  "link",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(interlink.CreateLinkRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ProjectIDArgSpec(),
			{
				Name:       "name",
				Short:      `Name of the link`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `List of tags to apply to the link`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "pop-id",
				Short:      `PoP (location) where the link will be created`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "bandwidth-mbps",
				Short:      `Desired bandwidth for the link. Must be compatible with available link bandwidths and remaining bandwidth capacity of the connection`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "connection-id",
				Short:      `If set, creates a self-hosted link using this dedicated physical connection. As the customer, specify the ID of the physical connection you already have for this link.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "partner-id",
				Short:      `If set, creates a hosted link on a partner's connection. Specify the ID of the chosen partner, who already has a shared connection with available bandwidth.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "peer-asn",
				Short:      `For self-hosted links we need the peer AS Number to establish BGP session. If not given, a default one will be assigned.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "vlan",
				Short:      `For self-hosted links only, it is possible to choose the VLAN ID. If the VLAN is not available (ie already taken or out of range), an error is returned.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "routing-policy-v4-id",
				Short:      `If set, attaches this routing policy containing IPv4 prefixes to the Link. Hence, a BGP IPv4 session will be created.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "routing-policy-v6-id",
				Short:      `If set, attaches this routing policy containing IPv6 prefixes to the Link. Hence, a BGP IPv6 session will be created.`,
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
			request := args.(*interlink.CreateLinkRequest)

			client := core.ExtractClient(ctx)
			api := interlink.NewAPI(client)

			return api.CreateLink(request)
		},
	}
}

func interlinkLinkUpdate() *core.Command {
	return &core.Command{
		Short:     `Update a link`,
		Long:      `Update an existing link, specified by its link ID. Only its name and tags can be updated.`,
		Namespace: "interlink",
		Resource:  "link",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(interlink.UpdateLinkRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "link-id",
				Short:      `ID of the link to update`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "name",
				Short:      `Name of the link`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `List of tags to apply to the link`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "peer-asn",
				Short:      `For self-hosted links, AS Number to establish BGP session.`,
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
			request := args.(*interlink.UpdateLinkRequest)

			client := core.ExtractClient(ctx)
			api := interlink.NewAPI(client)

			return api.UpdateLink(request)
		},
	}
}

func interlinkLinkDelete() *core.Command {
	return &core.Command{
		Short:     `Delete a link`,
		Long:      `Delete an existing link, specified by its link ID. Note that as well as deleting the link here on the Scaleway side, it is also necessary to request deletion from the partner on their side. Only when this action has been carried out on both sides will the resource be completely deleted.`,
		Namespace: "interlink",
		Resource:  "link",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(interlink.DeleteLinkRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "link-id",
				Short:      `ID of the link to delete`,
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
			request := args.(*interlink.DeleteLinkRequest)

			client := core.ExtractClient(ctx)
			api := interlink.NewAPI(client)

			return api.DeleteLink(request)
		},
	}
}

func interlinkLinkAttachVpc() *core.Command {
	return &core.Command{
		Short:     `Attach a VPC`,
		Long:      `Attach a VPC to an existing link. This facilitates communication between the resources in your Scaleway VPC, and your on-premises infrastructure.`,
		Namespace: "interlink",
		Resource:  "link",
		Verb:      "attach_vpc",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(interlink.AttachVpcRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "link-id",
				Short:      `ID of the link to attach VPC to`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "vpc-id",
				Short:      `ID of the VPC to attach`,
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
			request := args.(*interlink.AttachVpcRequest)

			client := core.ExtractClient(ctx)
			api := interlink.NewAPI(client)

			return api.AttachVpc(request)
		},
	}
}

func interlinkLinkDetachVpc() *core.Command {
	return &core.Command{
		Short:     `Detach a VPC`,
		Long:      `Detach a VPC from an existing link.`,
		Namespace: "interlink",
		Resource:  "link",
		Verb:      "detach_vpc",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(interlink.DetachVpcRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "link-id",
				Short:      `ID of the link to detach the VPC from`,
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
			request := args.(*interlink.DetachVpcRequest)

			client := core.ExtractClient(ctx)
			api := interlink.NewAPI(client)

			return api.DetachVpc(request)
		},
	}
}

func interlinkLinkAttachPolicy() *core.Command {
	return &core.Command{
		Short:     `Attach a routing policy`,
		Long:      `Attach a routing policy to an existing link. As all routes across the link are blocked by default, you must attach a routing policy to set IP prefix filters for allowed routes, facilitating traffic flow.`,
		Namespace: "interlink",
		Resource:  "link",
		Verb:      "attach_policy",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(interlink.AttachRoutingPolicyRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "link-id",
				Short:      `ID of the link to attach a routing policy to`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "routing-policy-id",
				Short:      `ID of the routing policy to be attached`,
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
			request := args.(*interlink.AttachRoutingPolicyRequest)

			client := core.ExtractClient(ctx)
			api := interlink.NewAPI(client)

			return api.AttachRoutingPolicy(request)
		},
	}
}

func interlinkLinkDetachPolicy() *core.Command {
	return &core.Command{
		Short:     `Detach a routing policy`,
		Long:      `Detach a routing policy from an existing link. Without a routing policy, all routes across the link are blocked by default.`,
		Namespace: "interlink",
		Resource:  "link",
		Verb:      "detach_policy",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(interlink.DetachRoutingPolicyRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "link-id",
				Short:      `ID of the link to detach a routing policy from`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "routing-policy-id",
				Short:      `ID of the routing policy to be detached`,
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
			request := args.(*interlink.DetachRoutingPolicyRequest)

			client := core.ExtractClient(ctx)
			api := interlink.NewAPI(client)

			return api.DetachRoutingPolicy(request)
		},
	}
}

func interlinkLinkEnablePropagation() *core.Command {
	return &core.Command{
		Short:     `Enable route propagation`,
		Long:      `Enable all allowed prefixes (defined in a routing policy) to be announced in the BGP session. This allows traffic to flow between the attached VPC and the on-premises infrastructure along the announced routes. Note that by default, even when route propagation is enabled, all routes are blocked. It is essential to attach a routing policy to define the ranges of routes to announce.`,
		Namespace: "interlink",
		Resource:  "link",
		Verb:      "enable_propagation",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(interlink.EnableRoutePropagationRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "link-id",
				Short:      `ID of the link on which to enable route propagation`,
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
			request := args.(*interlink.EnableRoutePropagationRequest)

			client := core.ExtractClient(ctx)
			api := interlink.NewAPI(client)

			return api.EnableRoutePropagation(request)
		},
	}
}

func interlinkLinkDisablePropagation() *core.Command {
	return &core.Command{
		Short:     `Disable route propagation`,
		Long:      `Prevent any prefixes from being announced in the BGP session. Traffic will not be able to flow over the InterLink until route propagation is re-enabled.`,
		Namespace: "interlink",
		Resource:  "link",
		Verb:      "disable_propagation",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(interlink.DisableRoutePropagationRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "link-id",
				Short:      `ID of the link on which to disable route propagation`,
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
			request := args.(*interlink.DisableRoutePropagationRequest)

			client := core.ExtractClient(ctx)
			api := interlink.NewAPI(client)

			return api.DisableRoutePropagation(request)
		},
	}
}

func interlinkRoutingPolicyList() *core.Command {
	return &core.Command{
		Short:     `List routing policies`,
		Long:      `List all routing policies in a given region. A routing policy can be attached to one or multiple links (InterLink connections).`,
		Namespace: "interlink",
		Resource:  "routing-policy",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(interlink.ListRoutingPoliciesRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `Order in which to return results`,
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
				Name:       "project-id",
				Short:      `Project ID to filter for`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "name",
				Short:      `Routing policy name to filter for`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `Tags to filter for`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "ipv6",
				Short:      `Filter for the routing policies based on IP prefixes version`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "organization-id",
				Short:      `Organization ID to filter for`,
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
			request := args.(*interlink.ListRoutingPoliciesRequest)

			client := core.ExtractClient(ctx)
			api := interlink.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			if request.Region == scw.Region(core.AllLocalities) {
				opts = append(opts, scw.WithRegions(api.Regions()...))
				request.Region = ""
			}
			resp, err := api.ListRoutingPolicies(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.RoutingPolicies, nil
		},
	}
}

func interlinkRoutingPolicyGet() *core.Command {
	return &core.Command{
		Short:     `Get routing policy`,
		Long:      `Get a routing policy for the given routing policy ID. The response object gives information including the policy's name, tags and prefix filters.`,
		Namespace: "interlink",
		Resource:  "routing-policy",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(interlink.GetRoutingPolicyRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "routing-policy-id",
				Short:      `ID of the routing policy to get`,
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
			request := args.(*interlink.GetRoutingPolicyRequest)

			client := core.ExtractClient(ctx)
			api := interlink.NewAPI(client)

			return api.GetRoutingPolicy(request)
		},
	}
}

func interlinkRoutingPolicyCreate() *core.Command {
	return &core.Command{
		Short:     `Create a routing policy`,
		Long:      `Create a routing policy. Routing policies allow you to set IP prefix filters to define the incoming route announcements to accept from the peer, and the outgoing routes to announce to the peer.`,
		Namespace: "interlink",
		Resource:  "routing-policy",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(interlink.CreateRoutingPolicyRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ProjectIDArgSpec(),
			{
				Name:       "name",
				Short:      `Name of the routing policy`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `List of tags to apply to the routing policy`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "prefix-filter-in.{index}",
				Short:      `IP prefixes to accept from the peer (ranges of route announcements to accept)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "prefix-filter-out.{index}",
				Short:      `IP prefix filters to advertise to the peer (ranges of routes to advertise)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "is-ipv6",
				Short:      `IP prefixes version of the routing policy`,
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
			request := args.(*interlink.CreateRoutingPolicyRequest)

			client := core.ExtractClient(ctx)
			api := interlink.NewAPI(client)

			return api.CreateRoutingPolicy(request)
		},
	}
}

func interlinkRoutingPolicyUpdate() *core.Command {
	return &core.Command{
		Short:     `Update a routing policy`,
		Long:      `Update an existing routing policy, specified by its routing policy ID. Its name, tags and incoming/outgoing prefix filters can be updated.`,
		Namespace: "interlink",
		Resource:  "routing-policy",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(interlink.UpdateRoutingPolicyRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "routing-policy-id",
				Short:      `ID of the routing policy to update`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "name",
				Short:      `Name of the routing policy`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `List of tags to apply to the routing policy`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "prefix-filter-in.{index}",
				Short:      `IP prefixes to accept from the peer (ranges of route announcements to accept)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "prefix-filter-out.{index}",
				Short:      `IP prefix filters for routes to advertise to the peer (ranges of routes to advertise)`,
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
			request := args.(*interlink.UpdateRoutingPolicyRequest)

			client := core.ExtractClient(ctx)
			api := interlink.NewAPI(client)

			return api.UpdateRoutingPolicy(request)
		},
	}
}

func interlinkRoutingPolicyDelete() *core.Command {
	return &core.Command{
		Short:     `Delete a routing policy`,
		Long:      `Delete an existing routing policy, specified by its routing policy ID.`,
		Namespace: "interlink",
		Resource:  "routing-policy",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(interlink.DeleteRoutingPolicyRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "routing-policy-id",
				Short:      `ID of the routing policy to delete`,
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
			request := args.(*interlink.DeleteRoutingPolicyRequest)

			client := core.ExtractClient(ctx)
			api := interlink.NewAPI(client)
			e = api.DeleteRoutingPolicy(request)
			if e != nil {
				return nil, e
			}

			return &core.SuccessResult{
				Resource: "routing-policy",
				Verb:     "delete",
			}, nil
		},
	}
}
