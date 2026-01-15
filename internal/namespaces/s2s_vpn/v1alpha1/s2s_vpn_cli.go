// This file was automatically generated. DO NOT EDIT.
// If you have any remark or suggestion do not hesitate to open an issue.

package s2s_vpn

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	s2s_vpn "github.com/scaleway/scaleway-sdk-go/api/s2s_vpn/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// always import dependencies
var (
	_ = scw.RegionFrPar
)

func GetGeneratedCommands() *core.Commands {
	return core.NewCommands(
		s2sVpnRoot(),
		s2sVpnVpnGateway(),
		s2sVpnVpnGatewayType(),
		s2sVpnConnection(),
		s2sVpnCustomerGateway(),
		s2sVpnRoutingPolicy(),
		s2sVpnVpnGatewayTypeList(),
		s2sVpnVpnGatewayList(),
		s2sVpnVpnGatewayGet(),
		s2sVpnVpnGatewayCreate(),
		s2sVpnVpnGatewayUpdate(),
		s2sVpnVpnGatewayDelete(),
		s2sVpnConnectionList(),
		s2sVpnConnectionGet(),
		s2sVpnConnectionCreate(),
		s2sVpnConnectionUpdate(),
		s2sVpnConnectionDelete(),
		s2sVpnConnectionRenewPsk(),
		s2sVpnConnectionSetRoutingPolicy(),
		s2sVpnConnectionDetachRoutingPolicy(),
		s2sVpnConnectionEnableRoutePropagation(),
		s2sVpnConnectionDisableRoutePropagation(),
		s2sVpnCustomerGatewayList(),
		s2sVpnCustomerGatewayGet(),
		s2sVpnCustomerGatewayCreate(),
		s2sVpnCustomerGatewayUpdate(),
		s2sVpnCustomerGatewayDelete(),
		s2sVpnRoutingPolicyList(),
		s2sVpnRoutingPolicyGet(),
		s2sVpnRoutingPolicyCreate(),
		s2sVpnRoutingPolicyUpdate(),
		s2sVpnRoutingPolicyDelete(),
	)
}

func s2sVpnRoot() *core.Command {
	return &core.Command{
		Short:     `This API allows you to manage your Site-to-Site VPN`,
		Long:      `This API allows you to manage your Site-to-Site VPN.`,
		Namespace: "s2s-vpn",
	}
}

func s2sVpnVpnGateway() *core.Command {
	return &core.Command{
		Short:     `A VPN gateway is an IPsec peer managed by Scaleway. It can support multiple connections to customer gateways.`,
		Long:      `A VPN gateway is an IPsec peer managed by Scaleway. It can support multiple connections to customer gateways.`,
		Namespace: "s2s-vpn",
		Resource:  "vpn-gateway",
	}
}

func s2sVpnVpnGatewayType() *core.Command {
	return &core.Command{
		Short:     `VPN gateways come in various shapes, sizes and prices, which are  described by VPN gateway types. They represent the different commercial  offer types for VPN gateways available at Scaleway.`,
		Long:      `VPN gateways come in various shapes, sizes and prices, which are  described by VPN gateway types. They represent the different commercial  offer types for VPN gateways available at Scaleway.`,
		Namespace: "s2s-vpn",
		Resource:  "vpn-gateway-type",
	}
}

func s2sVpnConnection() *core.Command {
	return &core.Command{
		Short:     `A connection represents the IPsec tunnel between VPN gateway and customer gateway.`,
		Long:      `A connection represents the IPsec tunnel between VPN gateway and customer gateway.`,
		Namespace: "s2s-vpn",
		Resource:  "connection",
	}
}

func s2sVpnCustomerGateway() *core.Command {
	return &core.Command{
		Short:     `A customer gateway represents a Scaleway client's device that communicates with a VPN gateway.`,
		Long:      `A customer gateway represents a Scaleway client's device that communicates with a VPN gateway.`,
		Namespace: "s2s-vpn",
		Resource:  "customer-gateway",
	}
}

func s2sVpnRoutingPolicy() *core.Command {
	return &core.Command{
		Short:     `By default, all routes across the Site-to-Site VPN (between VPN gateway and customer gateway) are blocked. Routing policies allow you to set filters to define the IP prefixes to allow.`,
		Long:      `By default, all routes across the Site-to-Site VPN (between VPN gateway and customer gateway) are blocked. Routing policies allow you to set filters to define the IP prefixes to allow.`,
		Namespace: "s2s-vpn",
		Resource:  "routing-policy",
	}
}

func s2sVpnVpnGatewayTypeList() *core.Command {
	return &core.Command{
		Short:     `List VPN gateway types`,
		Long:      `List the different VPN gateway commercial offer types available at Scaleway. The response is an array of objects describing the name and technical details of each available VPN gateway type.`,
		Namespace: "s2s-vpn",
		Resource:  "vpn-gateway-type",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(s2s_vpn.ListVpnGatewayTypesRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
				scw.Region(core.AllLocalities),
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*s2s_vpn.ListVpnGatewayTypesRequest)

			client := core.ExtractClient(ctx)
			api := s2s_vpn.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			if request.Region == scw.Region(core.AllLocalities) {
				opts = append(opts, scw.WithRegions(api.Regions()...))
				request.Region = ""
			}
			resp, err := api.ListVpnGatewayTypes(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.GatewayTypes, nil
		},
	}
}

func s2sVpnVpnGatewayList() *core.Command {
	return &core.Command{
		Short:     `List VPN gateways`,
		Long:      `List all your VPN gateways. A number of filters are available, including Project ID, name, tags and status.`,
		Namespace: "s2s-vpn",
		Resource:  "vpn-gateway",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(s2s_vpn.ListVpnGatewaysRequest{}),
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
					"type_asc",
					"type_desc",
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
				Short:      `VPN gateway name to filter for`,
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
				Name:       "statuses.{index}",
				Short:      `VPN gateway statuses to filter for`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_status",
					"configuring",
					"failed",
					"provisioning",
					"active",
					"deprovisioning",
					"locked",
				},
			},
			{
				Name:       "gateway-types.{index}",
				Short:      `Filter for VPN gateways of these types`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "private-network-ids.{index}",
				Short:      `Filter for VPN gateways attached to these private networks`,
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
			request := args.(*s2s_vpn.ListVpnGatewaysRequest)

			client := core.ExtractClient(ctx)
			api := s2s_vpn.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			if request.Region == scw.Region(core.AllLocalities) {
				opts = append(opts, scw.WithRegions(api.Regions()...))
				request.Region = ""
			}
			resp, err := api.ListVpnGateways(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.Gateways, nil
		},
	}
}

func s2sVpnVpnGatewayGet() *core.Command {
	return &core.Command{
		Short:     `Get a VPN gateway`,
		Long:      `Get a VPN gateway for the given VPN gateway ID.`,
		Namespace: "s2s-vpn",
		Resource:  "vpn-gateway",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(s2s_vpn.GetVpnGatewayRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "gateway-id",
				Short:      `ID of the requested VPN gateway`,
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
			request := args.(*s2s_vpn.GetVpnGatewayRequest)

			client := core.ExtractClient(ctx)
			api := s2s_vpn.NewAPI(client)

			return api.GetVpnGateway(request)
		},
	}
}

func s2sVpnVpnGatewayCreate() *core.Command {
	return &core.Command{
		Short:     `Create VPN gateway`,
		Long:      `Create VPN gateway.`,
		Namespace: "s2s-vpn",
		Resource:  "vpn-gateway",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(s2s_vpn.CreateVpnGatewayRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ProjectIDArgSpec(),
			{
				Name:       "name",
				Short:      `Name of the VPN gateway`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `List of tags to apply to the VPN gateway`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "gateway-type",
				Short:      `VPN gateway type (commercial offer type)`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "public-config.ipam-ipv4-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "public-config.ipam-ipv6-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "private-network-id",
				Short:      `ID of the Private Network to attach to the VPN gateway`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "ipam-private-ipv4-id",
				Short:      `ID of the IPAM private IPv4 address to attach to the VPN gateway`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "ipam-private-ipv6-id",
				Short:      `ID of the IPAM private IPv6 address to attach to the VPN gateway`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "zone",
				Short:      `Availability Zone where the VPN gateway should be provisioned. If no zone is specified, the VPN gateway will be automatically placed.`,
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
			request := args.(*s2s_vpn.CreateVpnGatewayRequest)

			client := core.ExtractClient(ctx)
			api := s2s_vpn.NewAPI(client)

			return api.CreateVpnGateway(request)
		},
	}
}

func s2sVpnVpnGatewayUpdate() *core.Command {
	return &core.Command{
		Short:     `Update a VPN gateway`,
		Long:      `Update an existing VPN gateway, specified by its VPN gateway ID. Only its name and tags can be updated.`,
		Namespace: "s2s-vpn",
		Resource:  "vpn-gateway",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(s2s_vpn.UpdateVpnGatewayRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "gateway-id",
				Short:      `ID of the VPN gateway to update`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "name",
				Short:      `Name of the VPN gateway`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `List of tags to apply to the VPN Gateway`,
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
			request := args.(*s2s_vpn.UpdateVpnGatewayRequest)

			client := core.ExtractClient(ctx)
			api := s2s_vpn.NewAPI(client)

			return api.UpdateVpnGateway(request)
		},
	}
}

func s2sVpnVpnGatewayDelete() *core.Command {
	return &core.Command{
		Short:     `Delete a VPN gateway`,
		Long:      `Delete an existing VPN gateway, specified by its VPN gateway ID.`,
		Namespace: "s2s-vpn",
		Resource:  "vpn-gateway",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(s2s_vpn.DeleteVpnGatewayRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "gateway-id",
				Short:      `ID of the VPN gateway to delete`,
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
			request := args.(*s2s_vpn.DeleteVpnGatewayRequest)

			client := core.ExtractClient(ctx)
			api := s2s_vpn.NewAPI(client)

			return api.DeleteVpnGateway(request)
		},
	}
}

func s2sVpnConnectionList() *core.Command {
	return &core.Command{
		Short:     `List connections`,
		Long:      `List all your connections. A number of filters are available, including Project ID, name, tags and status.`,
		Namespace: "s2s-vpn",
		Resource:  "connection",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(s2s_vpn.ListConnectionsRequest{}),
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
				Short:      `Connection name to filter for`,
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
				Name:       "statuses.{index}",
				Short:      `Connection statuses to filter for`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_status",
					"active",
					"limited_connectivity",
					"down",
					"locked",
				},
			},
			{
				Name:       "is-ipv6",
				Short:      `Filter connections with IP version of IPSec tunnel`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "routing-policy-ids.{index}",
				Short:      `Filter for connections using these routing policies`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "route-propagation-enabled",
				Short:      `Filter for connections with route propagation enabled`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "vpn-gateway-ids.{index}",
				Short:      `Filter for connections attached to these VPN gateways`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "customer-gateway-ids.{index}",
				Short:      `Filter for connections attached to these customer gateways`,
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
			request := args.(*s2s_vpn.ListConnectionsRequest)

			client := core.ExtractClient(ctx)
			api := s2s_vpn.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			if request.Region == scw.Region(core.AllLocalities) {
				opts = append(opts, scw.WithRegions(api.Regions()...))
				request.Region = ""
			}
			resp, err := api.ListConnections(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.Connections, nil
		},
	}
}

func s2sVpnConnectionGet() *core.Command {
	return &core.Command{
		Short:     `Get a connection`,
		Long:      `Get a connection for the given connection ID. The response object includes information about the connection's various configuration details.`,
		Namespace: "s2s-vpn",
		Resource:  "connection",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(s2s_vpn.GetConnectionRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "connection-id",
				Short:      `ID of the requested connection`,
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
			request := args.(*s2s_vpn.GetConnectionRequest)

			client := core.ExtractClient(ctx)
			api := s2s_vpn.NewAPI(client)

			return api.GetConnection(request)
		},
	}
}

func s2sVpnConnectionCreate() *core.Command {
	return &core.Command{
		Short:     `Create a connection`,
		Long:      `Create a connection.`,
		Namespace: "s2s-vpn",
		Resource:  "connection",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(s2s_vpn.CreateConnectionRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ProjectIDArgSpec(),
			{
				Name:       "name",
				Short:      `Name of the connection`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `List of tags to apply to the connection`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "is-ipv6",
				Short:      `Defines IP version of the IPSec Tunnel`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "initiation-policy",
				Short:      `Who initiates the IPsec tunnel`,
				Required:   true,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_initiation_policy",
					"vpn_gateway",
					"customer_gateway",
				},
			},
			{
				Name:       "ikev2-ciphers.{index}.encryption",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_encryption",
					"aes128",
					"aes192",
					"aes256",
					"aes128gcm",
					"aes192gcm",
					"aes256gcm",
					"aes128ccm",
					"aes256ccm",
					"chacha20poly1305",
				},
			},
			{
				Name:       "ikev2-ciphers.{index}.integrity",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_integrity",
					"sha256",
					"sha384",
					"sha512",
				},
			},
			{
				Name:       "ikev2-ciphers.{index}.dh-group",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_dhgroup",
					"modp2048",
					"modp3072",
					"modp4096",
					"ecp256",
					"ecp384",
					"ecp521",
					"curve25519",
				},
			},
			{
				Name:       "esp-ciphers.{index}.encryption",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_encryption",
					"aes128",
					"aes192",
					"aes256",
					"aes128gcm",
					"aes192gcm",
					"aes256gcm",
					"aes128ccm",
					"aes256ccm",
					"chacha20poly1305",
				},
			},
			{
				Name:       "esp-ciphers.{index}.integrity",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_integrity",
					"sha256",
					"sha384",
					"sha512",
				},
			},
			{
				Name:       "esp-ciphers.{index}.dh-group",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_dhgroup",
					"modp2048",
					"modp3072",
					"modp4096",
					"ecp256",
					"ecp384",
					"ecp521",
					"curve25519",
				},
			},
			{
				Name:       "enable-route-propagation",
				Short:      `Defines whether route propagation is enabled or not.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "vpn-gateway-id",
				Short:      `ID of the VPN gateway to attach to the connection`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "customer-gateway-id",
				Short:      `ID of the customer gateway to attach to the connection`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "bgp-config-ipv4.routing-policy-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "bgp-config-ipv4.private-ip",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "bgp-config-ipv4.peer-private-ip",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "bgp-config-ipv6.routing-policy-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "bgp-config-ipv6.private-ip",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "bgp-config-ipv6.peer-private-ip",
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
			request := args.(*s2s_vpn.CreateConnectionRequest)

			client := core.ExtractClient(ctx)
			api := s2s_vpn.NewAPI(client)

			return api.CreateConnection(request)
		},
	}
}

func s2sVpnConnectionUpdate() *core.Command {
	return &core.Command{
		Short:     `Update a connection`,
		Long:      `Update an existing connection, specified by its connection ID.`,
		Namespace: "s2s-vpn",
		Resource:  "connection",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(s2s_vpn.UpdateConnectionRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "connection-id",
				Short:      `ID of the connection to update`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "name",
				Short:      `Name of the connection`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `List of tags to apply to the connection`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "initiation-policy",
				Short:      `Who initiates the IPsec tunnel`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_initiation_policy",
					"vpn_gateway",
					"customer_gateway",
				},
			},
			{
				Name:       "ikev2-ciphers.{index}.encryption",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_encryption",
					"aes128",
					"aes192",
					"aes256",
					"aes128gcm",
					"aes192gcm",
					"aes256gcm",
					"aes128ccm",
					"aes256ccm",
					"chacha20poly1305",
				},
			},
			{
				Name:       "ikev2-ciphers.{index}.integrity",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_integrity",
					"sha256",
					"sha384",
					"sha512",
				},
			},
			{
				Name:       "ikev2-ciphers.{index}.dh-group",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_dhgroup",
					"modp2048",
					"modp3072",
					"modp4096",
					"ecp256",
					"ecp384",
					"ecp521",
					"curve25519",
				},
			},
			{
				Name:       "esp-ciphers.{index}.encryption",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_encryption",
					"aes128",
					"aes192",
					"aes256",
					"aes128gcm",
					"aes192gcm",
					"aes256gcm",
					"aes128ccm",
					"aes256ccm",
					"chacha20poly1305",
				},
			},
			{
				Name:       "esp-ciphers.{index}.integrity",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_integrity",
					"sha256",
					"sha384",
					"sha512",
				},
			},
			{
				Name:       "esp-ciphers.{index}.dh-group",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_dhgroup",
					"modp2048",
					"modp3072",
					"modp4096",
					"ecp256",
					"ecp384",
					"ecp521",
					"curve25519",
				},
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*s2s_vpn.UpdateConnectionRequest)

			client := core.ExtractClient(ctx)
			api := s2s_vpn.NewAPI(client)

			return api.UpdateConnection(request)
		},
	}
}

func s2sVpnConnectionDelete() *core.Command {
	return &core.Command{
		Short:     `Delete a connection`,
		Long:      `Delete an existing connection, specified by its connection ID.`,
		Namespace: "s2s-vpn",
		Resource:  "connection",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(s2s_vpn.DeleteConnectionRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "connection-id",
				Short:      `ID of the connection to delete`,
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
			request := args.(*s2s_vpn.DeleteConnectionRequest)

			client := core.ExtractClient(ctx)
			api := s2s_vpn.NewAPI(client)
			e = api.DeleteConnection(request)
			if e != nil {
				return nil, e
			}

			return &core.SuccessResult{
				Resource: "connection",
				Verb:     "delete",
			}, nil
		},
	}
}

func s2sVpnConnectionRenewPsk() *core.Command {
	return &core.Command{
		Short:     `Renew pre-shared key`,
		Long:      `Renew pre-shared key for a given connection.`,
		Namespace: "s2s-vpn",
		Resource:  "connection",
		Verb:      "renew-psk",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(s2s_vpn.RenewConnectionPskRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "connection-id",
				Short:      `ID of the connection to renew the PSK`,
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
			request := args.(*s2s_vpn.RenewConnectionPskRequest)

			client := core.ExtractClient(ctx)
			api := s2s_vpn.NewAPI(client)

			return api.RenewConnectionPsk(request)
		},
	}
}

func s2sVpnConnectionSetRoutingPolicy() *core.Command {
	return &core.Command{
		Short:     `Set a new routing policy`,
		Long:      `Set a new routing policy on a connection, overriding the existing one if present, specified by its connection ID.`,
		Namespace: "s2s-vpn",
		Resource:  "connection",
		Verb:      "set-routing-policy",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(s2s_vpn.SetRoutingPolicyRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "connection-id",
				Short:      `ID of the connection whose routing policy is being updated`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "routing-policy-v4",
				Short:      `ID of the routing policy to set for the BGP IPv4 session`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "routing-policy-v6",
				Short:      `ID of the routing policy to set for the BGP IPv6 session`,
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
			request := args.(*s2s_vpn.SetRoutingPolicyRequest)

			client := core.ExtractClient(ctx)
			api := s2s_vpn.NewAPI(client)

			return api.SetRoutingPolicy(request)
		},
	}
}

func s2sVpnConnectionDetachRoutingPolicy() *core.Command {
	return &core.Command{
		Short:     `Detach a routing policy`,
		Long:      `Detach an existing routing policy from a connection, specified by its connection ID.`,
		Namespace: "s2s-vpn",
		Resource:  "connection",
		Verb:      "detach-routing-policy",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(s2s_vpn.DetachRoutingPolicyRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "connection-id",
				Short:      `ID of the connection from which routing policy is being detached`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "routing-policy-v4",
				Short:      `ID of the routing policy to detach from the BGP IPv4 session`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "routing-policy-v6",
				Short:      `ID of the routing policy to detach from the BGP IPv6 session`,
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
			request := args.(*s2s_vpn.DetachRoutingPolicyRequest)

			client := core.ExtractClient(ctx)
			api := s2s_vpn.NewAPI(client)

			return api.DetachRoutingPolicy(request)
		},
	}
}

func s2sVpnConnectionEnableRoutePropagation() *core.Command {
	return &core.Command{
		Short:     `Enable route propagation`,
		Long:      `Enable all allowed prefixes (defined in a routing policy) to be announced in the BGP session. This allows traffic to flow between the attached VPC and the on-premises infrastructure along the announced routes. Note that by default, even when route propagation is enabled, all routes are blocked. It is essential to attach a routing policy to define the ranges of routes to announce.`,
		Namespace: "s2s-vpn",
		Resource:  "connection",
		Verb:      "enable-route-propagation",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(s2s_vpn.EnableRoutePropagationRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "connection-id",
				Short:      `ID of the connection on which to enable route propagation`,
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
			request := args.(*s2s_vpn.EnableRoutePropagationRequest)

			client := core.ExtractClient(ctx)
			api := s2s_vpn.NewAPI(client)

			return api.EnableRoutePropagation(request)
		},
	}
}

func s2sVpnConnectionDisableRoutePropagation() *core.Command {
	return &core.Command{
		Short:     `Disable route propagation`,
		Long:      `Prevent any prefixes from being announced in the BGP session. Traffic will not be able to flow over the VPN Gateway until route propagation is re-enabled.`,
		Namespace: "s2s-vpn",
		Resource:  "connection",
		Verb:      "disable-route-propagation",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(s2s_vpn.DisableRoutePropagationRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "connection-id",
				Short:      `ID of the connection on which to disable route propagation`,
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
			request := args.(*s2s_vpn.DisableRoutePropagationRequest)

			client := core.ExtractClient(ctx)
			api := s2s_vpn.NewAPI(client)

			return api.DisableRoutePropagation(request)
		},
	}
}

func s2sVpnCustomerGatewayList() *core.Command {
	return &core.Command{
		Short:     `List customer gateways`,
		Long:      `List all your customer gateways. A number of filters are available, including Project ID, name, and tags.`,
		Namespace: "s2s-vpn",
		Resource:  "customer-gateway",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(s2s_vpn.ListCustomerGatewaysRequest{}),
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
				Short:      `Customer gateway name to filter for`,
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
			request := args.(*s2s_vpn.ListCustomerGatewaysRequest)

			client := core.ExtractClient(ctx)
			api := s2s_vpn.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			if request.Region == scw.Region(core.AllLocalities) {
				opts = append(opts, scw.WithRegions(api.Regions()...))
				request.Region = ""
			}
			resp, err := api.ListCustomerGateways(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.Gateways, nil
		},
	}
}

func s2sVpnCustomerGatewayGet() *core.Command {
	return &core.Command{
		Short:     `Get a customer gateway`,
		Long:      `Get a customer gateway for the given customer gateway ID.`,
		Namespace: "s2s-vpn",
		Resource:  "customer-gateway",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(s2s_vpn.GetCustomerGatewayRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "gateway-id",
				Short:      `ID of the requested customer gateway`,
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
			request := args.(*s2s_vpn.GetCustomerGatewayRequest)

			client := core.ExtractClient(ctx)
			api := s2s_vpn.NewAPI(client)

			return api.GetCustomerGateway(request)
		},
	}
}

func s2sVpnCustomerGatewayCreate() *core.Command {
	return &core.Command{
		Short:     `Create a customer gateway`,
		Long:      `Create a customer gateway.`,
		Namespace: "s2s-vpn",
		Resource:  "customer-gateway",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(s2s_vpn.CreateCustomerGatewayRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ProjectIDArgSpec(),
			{
				Name:       "name",
				Short:      `Name of the customer gateway`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `List of tags to apply to the customer gateway`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "ipv4-public",
				Short:      `Public IPv4 address of the customer gateway`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "ipv6-public",
				Short:      `Public IPv6 address of the customer gateway`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "asn",
				Short:      `AS Number of the customer gateway`,
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
			request := args.(*s2s_vpn.CreateCustomerGatewayRequest)

			client := core.ExtractClient(ctx)
			api := s2s_vpn.NewAPI(client)

			return api.CreateCustomerGateway(request)
		},
	}
}

func s2sVpnCustomerGatewayUpdate() *core.Command {
	return &core.Command{
		Short:     `Update a customer gateway`,
		Long:      `Update an existing customer gateway, specified by its customer gateway ID. You can update its name, tags, public IPv4 & IPv6 address and AS Number.`,
		Namespace: "s2s-vpn",
		Resource:  "customer-gateway",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(s2s_vpn.UpdateCustomerGatewayRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "gateway-id",
				Short:      `ID of the customer gateway to update`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "name",
				Short:      `Name of the customer gateway`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `List of tags to apply to the customer gateway`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "ipv4-public",
				Short:      `Public IPv4 address of the customer gateway`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "ipv6-public",
				Short:      `Public IPv6 address of the customer gateway`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "asn",
				Short:      `AS Number of the customer gateway`,
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
			request := args.(*s2s_vpn.UpdateCustomerGatewayRequest)

			client := core.ExtractClient(ctx)
			api := s2s_vpn.NewAPI(client)

			return api.UpdateCustomerGateway(request)
		},
	}
}

func s2sVpnCustomerGatewayDelete() *core.Command {
	return &core.Command{
		Short:     `Delete a customer gateway`,
		Long:      `Delete an existing customer gateway, specified by its customer gateway ID.`,
		Namespace: "s2s-vpn",
		Resource:  "customer-gateway",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(s2s_vpn.DeleteCustomerGatewayRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "gateway-id",
				Short:      `ID of the customer gateway to delete`,
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
			request := args.(*s2s_vpn.DeleteCustomerGatewayRequest)

			client := core.ExtractClient(ctx)
			api := s2s_vpn.NewAPI(client)
			e = api.DeleteCustomerGateway(request)
			if e != nil {
				return nil, e
			}

			return &core.SuccessResult{
				Resource: "customer-gateway",
				Verb:     "delete",
			}, nil
		},
	}
}

func s2sVpnRoutingPolicyList() *core.Command {
	return &core.Command{
		Short:     `List routing policies`,
		Long:      `List all routing policies in a given region. A routing policy can be attached to one or multiple connections (S2S VPN connections).`,
		Namespace: "s2s-vpn",
		Resource:  "routing-policy",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(s2s_vpn.ListRoutingPoliciesRequest{}),
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
			request := args.(*s2s_vpn.ListRoutingPoliciesRequest)

			client := core.ExtractClient(ctx)
			api := s2s_vpn.NewAPI(client)
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

func s2sVpnRoutingPolicyGet() *core.Command {
	return &core.Command{
		Short:     `Get routing policy`,
		Long:      `Get a routing policy for the given routing policy ID. The response object gives information including the policy's name, tags and prefix filters.`,
		Namespace: "s2s-vpn",
		Resource:  "routing-policy",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(s2s_vpn.GetRoutingPolicyRequest{}),
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
			request := args.(*s2s_vpn.GetRoutingPolicyRequest)

			client := core.ExtractClient(ctx)
			api := s2s_vpn.NewAPI(client)

			return api.GetRoutingPolicy(request)
		},
	}
}

func s2sVpnRoutingPolicyCreate() *core.Command {
	return &core.Command{
		Short:     `Create a routing policy`,
		Long:      `Create a routing policy. Routing policies allow you to set IP prefix filters to define the incoming route announcements to accept from the customer gateway, and the outgoing routes to announce to the customer gateway.`,
		Namespace: "s2s-vpn",
		Resource:  "routing-policy",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(s2s_vpn.CreateRoutingPolicyRequest{}),
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
				Name:       "is-ipv6",
				Short:      `IP prefixes version of the routing policy`,
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
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*s2s_vpn.CreateRoutingPolicyRequest)

			client := core.ExtractClient(ctx)
			api := s2s_vpn.NewAPI(client)

			return api.CreateRoutingPolicy(request)
		},
	}
}

func s2sVpnRoutingPolicyUpdate() *core.Command {
	return &core.Command{
		Short:     `Update a routing policy`,
		Long:      `Update an existing routing policy, specified by its routing policy ID. Its name, tags and incoming/outgoing prefix filters can be updated.`,
		Namespace: "s2s-vpn",
		Resource:  "routing-policy",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(s2s_vpn.UpdateRoutingPolicyRequest{}),
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
			request := args.(*s2s_vpn.UpdateRoutingPolicyRequest)

			client := core.ExtractClient(ctx)
			api := s2s_vpn.NewAPI(client)

			return api.UpdateRoutingPolicy(request)
		},
	}
}

func s2sVpnRoutingPolicyDelete() *core.Command {
	return &core.Command{
		Short:     `Delete a routing policy`,
		Long:      `Delete an existing routing policy, specified by its routing policy ID.`,
		Namespace: "s2s-vpn",
		Resource:  "routing-policy",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(s2s_vpn.DeleteRoutingPolicyRequest{}),
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
			request := args.(*s2s_vpn.DeleteRoutingPolicyRequest)

			client := core.ExtractClient(ctx)
			api := s2s_vpn.NewAPI(client)
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
