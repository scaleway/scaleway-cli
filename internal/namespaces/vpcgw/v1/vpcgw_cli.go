// This file was automatically generated. DO NOT EDIT.
// If you have any remark or suggestion do not hesitate to open an issue.

package vpcgw

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-sdk-go/api/vpcgw/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// always import dependencies
var (
	_ = scw.RegionFrPar
)

func GetGeneratedCommands() *core.Commands {
	return core.NewCommands(
		vpcGwRoot(),
		vpcGwGateway(),
		vpcGwGatewayNetwork(),
		vpcGwDHCP(),
		vpcGwDHCPEntry(),
		vpcGwPatRule(),
		vpcGwIP(),
		vpcGwGatewayType(),
		vpcGwGatewayList(),
		vpcGwGatewayGet(),
		vpcGwGatewayCreate(),
		vpcGwGatewayUpdate(),
		vpcGwGatewayDelete(),
		vpcGwGatewayUpgrade(),
		vpcGwGatewayEnableIPMobility(),
		vpcGwGatewayNetworkList(),
		vpcGwGatewayNetworkGet(),
		vpcGwGatewayNetworkCreate(),
		vpcGwGatewayNetworkUpdate(),
		vpcGwGatewayNetworkDelete(),
		vpcGwDHCPList(),
		vpcGwDHCPGet(),
		vpcGwDHCPCreate(),
		vpcGwDHCPUpdate(),
		vpcGwDHCPDelete(),
		vpcGwDHCPEntryList(),
		vpcGwDHCPEntryGet(),
		vpcGwDHCPEntryCreate(),
		vpcGwDHCPEntryUpdate(),
		vpcGwDHCPEntrySet(),
		vpcGwDHCPEntryDelete(),
		vpcGwPatRuleList(),
		vpcGwPatRuleGet(),
		vpcGwPatRuleCreate(),
		vpcGwPatRuleUpdate(),
		vpcGwPatRuleSet(),
		vpcGwPatRuleDelete(),
		vpcGwGatewayTypeList(),
		vpcGwIPList(),
		vpcGwIPGet(),
		vpcGwIPCreate(),
		vpcGwIPUpdate(),
		vpcGwIPDelete(),
		vpcGwGatewayRefreshSSHKeys(),
		vpcGwGatewayMigrateToV2(),
	)
}

func vpcGwRoot() *core.Command {
	return &core.Command{
		Short:     `This API allows you to manage your Public Gateways`,
		Long:      `This API allows you to manage your Public Gateways.`,
		Namespace: "vpc-gw",
	}
}

func vpcGwGateway() *core.Command {
	return &core.Command{
		Short:     `Public Gateway management`,
		Long:      `Public Gateways are building blocks for your infrastructure on Scaleway's shared public cloud. They provide a set of managed network services and features for Scaleway's Private Networks such as DHCP, NAT and routing.`,
		Namespace: "vpc-gw",
		Resource:  "gateway",
	}
}

func vpcGwGatewayNetwork() *core.Command {
	return &core.Command{
		Short:     `Gateway Networks management`,
		Long:      `A Gateway Network represents the connection of a Private Network to a Public Gateway. It holds configuration options relative to this specific connection, such as the DHCP configuration.`,
		Namespace: "vpc-gw",
		Resource:  "gateway-network",
	}
}

func vpcGwDHCP() *core.Command {
	return &core.Command{
		Short:     `DHCP configuration management`,
		Long:      `These objects define a DHCP configuration, i.e. how IP addresses should be assigned to devices on a Private Network attached to a Public Gateway. Definable parameters include the subnet for the DHCP server, the validity  period for DHCP entries, whether to use dynamic pooling, and more. A DHCP configuration object has a DHCP ID, which can then be used as part of a  call to create or update a Gateway Network. This lets you attach an existing DHCP configuration to a Public Gateway attached to a Private Network. Similarly, you can use a DHCP ID as a query parameter to list Gateway Networks which use this DHCP configuration object.`,
		Namespace: "vpc-gw",
		Resource:  "dhcp",
	}
}

func vpcGwDHCPEntry() *core.Command {
	return &core.Command{
		Short:     `DHCP entries management`,
		Long:      `DHCP entries belong to a specified Gateway Network (Public Gateway / Private Network connection). A DHCP entry can hold either a dynamic DHCP lease (an IP address dynamically assigned by the Public Gateway to a device) or a static, user-created DHCP reservation.`,
		Namespace: "vpc-gw",
		Resource:  "dhcp-entry",
	}
}

func vpcGwPatRule() *core.Command {
	return &core.Command{
		Short:     `PAT rules management`,
		Long:      `PAT (Port Address Translation) rules, aka static NAT rules, belong to a specified Public Gateway.  They define the forwarding of a public port to a specific device on a Private Network, enabling enables ingress traffic from the public Internet  to reach the correct device in the Private Network.`,
		Namespace: "vpc-gw",
		Resource:  "pat-rule",
	}
}

func vpcGwIP() *core.Command {
	return &core.Command{
		Short:     `IP address management`,
		Long:      `Public, flexible IP addresses for Public Gateways, allowing the gateway to reach the public internet, as well as forward (masquerade) traffic from member devices of attached Private Networks.`,
		Namespace: "vpc-gw",
		Resource:  "ip",
	}
}

func vpcGwGatewayType() *core.Command {
	return &core.Command{
		Short:     `Gateway types information`,
		Long:      `Public Gateways come in various shapes, sizes and prices, which are  described by gateway types. They represent the different commercial  offer types for Public Gateways available at Scaleway.`,
		Namespace: "vpc-gw",
		Resource:  "gateway-type",
	}
}

func vpcGwGatewayList() *core.Command {
	return &core.Command{
		Short:     `List Public Gateways`,
		Long:      `List Public Gateways in a given Scaleway Organization or Project. By default, results are displayed in ascending order of creation date.`,
		Namespace: "vpc-gw",
		Resource:  "gateway",
		Verb:      "list",
		// Deprecated:    true,
		ArgsType: reflect.TypeOf(vpcgw.ListGatewaysRequest{}),
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
				Short:      `Include only gateways in this Project`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "name",
				Short:      `Filter for gateways which have this search term in their name`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `Filter for gateways with these tags`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "type",
				Short:      `Filter for gateways of this type`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "status",
				Short:      `Filter for gateways with this current status. Use ` + "`" + `unknown` + "`" + ` to include all statuses`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown",
					"stopped",
					"allocating",
					"configuring",
					"running",
					"stopping",
					"failed",
					"deleting",
					"deleted",
					"locked",
				},
			},
			{
				Name:       "private-network-id",
				Short:      `Filter for gateways attached to this Private nNetwork`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "organization-id",
				Short:      `Include only gateways in this Organization`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
				scw.Zone(core.AllLocalities),
			),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*vpcgw.ListGatewaysRequest)

			client := core.ExtractClient(ctx)
			api := vpcgw.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			if request.Zone == scw.Zone(core.AllLocalities) {
				opts = append(opts, scw.WithZones(api.Zones()...))
				request.Zone = ""
			}
			resp, err := api.ListGateways(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.Gateways, nil
		},
		View: &core.View{Fields: []*core.ViewField{
			{
				FieldName: "ID",
			},
			{
				FieldName: "Name",
			},
			{
				FieldName: "Status",
			},
			{
				FieldName: "Tags",
			},
			{
				FieldName: "GatewayNetworks",
			},
			{
				FieldName: "UpstreamDNSServers",
			},
			{
				FieldName: "UpdatedAt",
			},
			{
				FieldName: "CreatedAt",
			},
			{
				FieldName: "Zone",
			},
			{
				FieldName: "ProjectID",
			},
			{
				FieldName: "OrganizationID",
			},
		}},
	}
}

func vpcGwGatewayGet() *core.Command {
	return &core.Command{
		Short:     `Get a Public Gateway`,
		Long:      `Get details of a Public Gateway, specified by its gateway ID. The response object contains full details of the gateway, including its **name**, **type**, **status** and more.`,
		Namespace: "vpc-gw",
		Resource:  "gateway",
		Verb:      "get",
		// Deprecated:    true,
		ArgsType: reflect.TypeOf(vpcgw.GetGatewayRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "gateway-id",
				Short:      `ID of the gateway to fetch`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
			),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*vpcgw.GetGatewayRequest)

			client := core.ExtractClient(ctx)
			api := vpcgw.NewAPI(client)

			return api.GetGateway(request)
		},
	}
}

func vpcGwGatewayCreate() *core.Command {
	return &core.Command{
		Short:     `Create a Public Gateway`,
		Long:      `Create a new Public Gateway in the specified Scaleway Project, defining its **name**, **type** and other configuration details such as whether to enable SSH bastion.`,
		Namespace: "vpc-gw",
		Resource:  "gateway",
		Verb:      "create",
		// Deprecated:    true,
		ArgsType: reflect.TypeOf(vpcgw.CreateGatewayRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ProjectIDArgSpec(),
			{
				Name:       "name",
				Short:      `Name for the gateway`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				Default:    core.RandomValueGenerator("gw"),
			},
			{
				Name:       "tags.{index}",
				Short:      `Tags for the gateway`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "type",
				Short:      `Gateway type (commercial offer type)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				Default:    core.DefaultValueSetter("VPC-GW-S"),
			},
			{
				Name:       "upstream-dns-servers.{index}",
				Short:      `Array of DNS server IP addresses to override the gateway's default recursive DNS servers`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "ip-id",
				Short:      `Existing IP address to attach to the gateway`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "enable-smtp",
				Short:      `Defines whether SMTP traffic should be allowed pass through the gateway`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "enable-bastion",
				Short:      `Defines whether SSH bastion should be enabled the gateway`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "bastion-port",
				Short:      `Port of the SSH bastion`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
			),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*vpcgw.CreateGatewayRequest)

			client := core.ExtractClient(ctx)
			api := vpcgw.NewAPI(client)

			return api.CreateGateway(request)
		},
	}
}

func vpcGwGatewayUpdate() *core.Command {
	return &core.Command{
		Short:     `Update a Public Gateway`,
		Long:      `Update the parameters of an existing Public Gateway, for example, its **name**, **tags**, **SSH bastion configuration**, and **DNS servers**.`,
		Namespace: "vpc-gw",
		Resource:  "gateway",
		Verb:      "update",
		// Deprecated:    true,
		ArgsType: reflect.TypeOf(vpcgw.UpdateGatewayRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "gateway-id",
				Short:      `ID of the gateway to update`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "name",
				Short:      `Name for the gateway`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `Tags for the gateway`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "upstream-dns-servers.{index}",
				Short:      `Array of DNS server IP addresses to override the gateway's default recursive DNS servers`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "enable-bastion",
				Short:      `Defines whether SSH bastion should be enabled the gateway`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "bastion-port",
				Short:      `Port of the SSH bastion`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "enable-smtp",
				Short:      `Defines whether SMTP traffic should be allowed to pass through the gateway`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
			),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*vpcgw.UpdateGatewayRequest)

			client := core.ExtractClient(ctx)
			api := vpcgw.NewAPI(client)

			return api.UpdateGateway(request)
		},
	}
}

func vpcGwGatewayDelete() *core.Command {
	return &core.Command{
		Short:     `Delete a Public Gateway`,
		Long:      `Delete an existing Public Gateway, specified by its gateway ID. This action is irreversible.`,
		Namespace: "vpc-gw",
		Resource:  "gateway",
		Verb:      "delete",
		// Deprecated:    true,
		ArgsType: reflect.TypeOf(vpcgw.DeleteGatewayRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "gateway-id",
				Short:      `ID of the gateway to delete`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "cleanup-dhcp",
				Short:      `Defines whether to clean up attached DHCP configurations (if any, and if not attached to another Gateway Network)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
			),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*vpcgw.DeleteGatewayRequest)

			client := core.ExtractClient(ctx)
			api := vpcgw.NewAPI(client)
			e = api.DeleteGateway(request)
			if e != nil {
				return nil, e
			}

			return &core.SuccessResult{
				Resource: "gateway",
				Verb:     "delete",
			}, nil
		},
	}
}

func vpcGwGatewayUpgrade() *core.Command {
	return &core.Command{
		Short:     `Upgrade a Public Gateway to the latest version and/or to a different commercial offer type`,
		Long:      `Upgrade a given Public Gateway to the newest software version or to a different commercial offer type. This applies the latest bugfixes and features to your Public Gateway. Note that gateway service will be interrupted during the update.`,
		Namespace: "vpc-gw",
		Resource:  "gateway",
		Verb:      "upgrade",
		// Deprecated:    true,
		ArgsType: reflect.TypeOf(vpcgw.UpgradeGatewayRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "gateway-id",
				Short:      `ID of the gateway to upgrade`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "type",
				Short:      `Gateway type (commercial offer)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
			),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*vpcgw.UpgradeGatewayRequest)

			client := core.ExtractClient(ctx)
			api := vpcgw.NewAPI(client)

			return api.UpgradeGateway(request)
		},
	}
}

func vpcGwGatewayEnableIPMobility() *core.Command {
	return &core.Command{
		Short:     `Upgrade a Public Gateway to IP mobility`,
		Long:      `Upgrade a Public Gateway to IP mobility (move from NAT IP to routed IP). This is idempotent: repeated calls after the first will return no error but have no effect.`,
		Namespace: "vpc-gw",
		Resource:  "gateway",
		Verb:      "enable-ip-mobility",
		// Deprecated:    true,
		ArgsType: reflect.TypeOf(vpcgw.EnableIPMobilityRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "gateway-id",
				Short:      `ID of the gateway to upgrade to IP mobility`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
			),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*vpcgw.EnableIPMobilityRequest)

			client := core.ExtractClient(ctx)
			api := vpcgw.NewAPI(client)
			e = api.EnableIPMobility(request)
			if e != nil {
				return nil, e
			}

			return &core.SuccessResult{
				Resource: "gateway",
				Verb:     "enable-ip-mobility",
			}, nil
		},
	}
}

func vpcGwGatewayNetworkList() *core.Command {
	return &core.Command{
		Short:     `List Public Gateway connections to Private Networks`,
		Long:      `List the connections between Public Gateways and Private Networks (a connection = a GatewayNetwork). You can choose to filter by ` + "`" + `gateway-id` + "`" + ` to list all Private Networks attached to the specified Public Gateway, or by ` + "`" + `private_network_id` + "`" + ` to list all Public Gateways attached to the specified Private Network. Other query parameters are also available. The result is an array of GatewayNetwork objects, each giving details of the connection between a given Public Gateway and a given Private Network.`,
		Namespace: "vpc-gw",
		Resource:  "gateway-network",
		Verb:      "list",
		// Deprecated:    true,
		ArgsType: reflect.TypeOf(vpcgw.ListGatewayNetworksRequest{}),
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
					"status_asc",
					"status_desc",
				},
			},
			{
				Name:       "gateway-id",
				Short:      `Filter for GatewayNetworks connected to this gateway`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "private-network-id",
				Short:      `Filter for GatewayNetworks connected to this Private Network`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "enable-masquerade",
				Short:      `Filter for GatewayNetworks with this ` + "`" + `enable_masquerade` + "`" + ` setting`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "dhcp-id",
				Short:      `Filter for GatewayNetworks using this DHCP configuration`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "status",
				Short:      `Filter for GatewayNetworks with this current status this status. Use ` + "`" + `unknown` + "`" + ` to include all statuses`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown",
					"created",
					"attaching",
					"configuring",
					"ready",
					"detaching",
					"deleted",
				},
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
				scw.Zone(core.AllLocalities),
			),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*vpcgw.ListGatewayNetworksRequest)

			client := core.ExtractClient(ctx)
			api := vpcgw.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			if request.Zone == scw.Zone(core.AllLocalities) {
				opts = append(opts, scw.WithZones(api.Zones()...))
				request.Zone = ""
			}
			resp, err := api.ListGatewayNetworks(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.GatewayNetworks, nil
		},
		View: &core.View{Fields: []*core.ViewField{
			{
				FieldName: "ID",
			},
			{
				FieldName: "GatewayID",
			},
			{
				FieldName: "PrivateNetworkID",
			},
			{
				FieldName: "Status",
			},
			{
				FieldName: "Address",
			},
			{
				FieldName: "MacAddress",
			},
			{
				FieldName: "EnableDHCP",
			},
			{
				FieldName: "EnableMasquerade",
			},
			{
				FieldName: "CreatedAt",
			},
			{
				FieldName: "UpdatedAt",
			},
			{
				FieldName: "Zone",
			},
		}},
	}
}

func vpcGwGatewayNetworkGet() *core.Command {
	return &core.Command{
		Short:     `Get a Public Gateway connection to a Private Network`,
		Long:      `Get details of a given connection between a Public Gateway and a Private Network (this connection = a GatewayNetwork), specified by its ` + "`" + `gateway_network_id` + "`" + `. The response object contains details of the connection including the IDs of the Public Gateway and Private Network, the dates the connection was created/updated and its configuration settings.`,
		Namespace: "vpc-gw",
		Resource:  "gateway-network",
		Verb:      "get",
		// Deprecated:    true,
		ArgsType: reflect.TypeOf(vpcgw.GetGatewayNetworkRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "gateway-network-id",
				Short:      `ID of the GatewayNetwork to fetch`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
			),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*vpcgw.GetGatewayNetworkRequest)

			client := core.ExtractClient(ctx)
			api := vpcgw.NewAPI(client)

			return api.GetGatewayNetwork(request)
		},
	}
}

func vpcGwGatewayNetworkCreate() *core.Command {
	return &core.Command{
		Short:     `Attach a Public Gateway to a Private Network`,
		Long:      `Attach a specific Public Gateway to a specific Private Network (create a GatewayNetwork). You can configure parameters for the connection including DHCP settings, whether to enable masquerade (dynamic NAT), and more.`,
		Namespace: "vpc-gw",
		Resource:  "gateway-network",
		Verb:      "create",
		// Deprecated:    true,
		ArgsType: reflect.TypeOf(vpcgw.CreateGatewayNetworkRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "gateway-id",
				Short:      `Public Gateway to connect`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "private-network-id",
				Short:      `Private Network to connect`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "enable-masquerade",
				Short:      `Defines whether to enable masquerade (dynamic NAT) on the GatewayNetwork.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "enable-dhcp",
				Short:      `Defines whether to enable DHCP on this Private Network.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "dhcp-id",
				Short:      `ID of an existing DHCP configuration object to use for this GatewayNetwork`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "address",
				Short:      `Static IP address in CIDR format to use without DHCP`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "ipam-config.push-default-route",
				Short:      `Enabling the default route also enables masquerading`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "ipam-config.ipam-ip-id",
				Short:      `Use this IPAM-booked IP ID as the Gateway's IP in this Private Network`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
			),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*vpcgw.CreateGatewayNetworkRequest)

			client := core.ExtractClient(ctx)
			api := vpcgw.NewAPI(client)

			return api.CreateGatewayNetwork(request)
		},
	}
}

func vpcGwGatewayNetworkUpdate() *core.Command {
	return &core.Command{
		Short:     `Update a Public Gateway's connection to a Private Network`,
		Long:      `Update the configuration parameters of a connection between a given Public Gateway and Private Network (the connection = a GatewayNetwork). Updatable parameters include DHCP settings and whether to enable traffic masquerade (dynamic NAT).`,
		Namespace: "vpc-gw",
		Resource:  "gateway-network",
		Verb:      "update",
		// Deprecated:    true,
		ArgsType: reflect.TypeOf(vpcgw.UpdateGatewayNetworkRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "gateway-network-id",
				Short:      `ID of the GatewayNetwork to update`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "enable-masquerade",
				Short:      `Defines whether to enable masquerade (dynamic NAT) on the GatewayNetwork.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "enable-dhcp",
				Short:      `Defines whether to enable DHCP on this Private Network.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "dhcp-id",
				Short:      `ID of the new DHCP configuration object to use with this GatewayNetwork`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "address",
				Short:      `New static IP address`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "ipam-config.push-default-route",
				Short:      `Enabling the default route also enables masquerading`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "ipam-config.ipam-ip-id",
				Short:      `Use this IPAM-booked IP ID as the Gateway's IP in this Private Network`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
			),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*vpcgw.UpdateGatewayNetworkRequest)

			client := core.ExtractClient(ctx)
			api := vpcgw.NewAPI(client)

			return api.UpdateGatewayNetwork(request)
		},
	}
}

func vpcGwGatewayNetworkDelete() *core.Command {
	return &core.Command{
		Short:     `Detach a Public Gateway from a Private Network`,
		Long:      `Detach a given Public Gateway from a given Private Network, i.e. delete a GatewayNetwork specified by a gateway_network_id.`,
		Namespace: "vpc-gw",
		Resource:  "gateway-network",
		Verb:      "delete",
		// Deprecated:    true,
		ArgsType: reflect.TypeOf(vpcgw.DeleteGatewayNetworkRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "gateway-network-id",
				Short:      `ID of the GatewayNetwork to delete`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "cleanup-dhcp",
				Short:      `Defines whether to clean up attached DHCP configurations (if any, and if not attached to another Gateway Network)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
			),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*vpcgw.DeleteGatewayNetworkRequest)

			client := core.ExtractClient(ctx)
			api := vpcgw.NewAPI(client)
			e = api.DeleteGatewayNetwork(request)
			if e != nil {
				return nil, e
			}

			return &core.SuccessResult{
				Resource: "gateway-network",
				Verb:     "delete",
			}, nil
		},
	}
}

func vpcGwDHCPList() *core.Command {
	return &core.Command{
		Short:     `List DHCP configurations`,
		Long:      `List DHCP configurations, optionally filtering by Organization, Project, Public Gateway IP address or more. The response is an array of DHCP configuration objects, each identified by a DHCP ID and containing configuration settings for the assignment of IP addresses to devices on a Private Network attached to a Public Gateway. Note that the response does not contain the IDs of any Private Network / Public Gateway the configuration is attached to. Use the ` + "`" + `List Public Gateway connections to Private Networks` + "`" + ` method for that purpose, filtering on DHCP ID.`,
		Namespace: "vpc-gw",
		Resource:  "dhcp",
		Verb:      "list",
		// Deprecated:    true,
		ArgsType: reflect.TypeOf(vpcgw.ListDHCPsRequest{}),
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
					"subnet_asc",
					"subnet_desc",
				},
			},
			{
				Name:       "project-id",
				Short:      `Include only DHCP configuration objects in this Project`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "address",
				Short:      `Filter for DHCP configuration objects with this DHCP server IP address (the gateway's address in the Private Network)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "has-address",
				Short:      `Filter for DHCP configuration objects with subnets containing this IP address`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "organization-id",
				Short:      `Include only DHCP configuration objects in this Organization`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
				scw.Zone(core.AllLocalities),
			),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*vpcgw.ListDHCPsRequest)

			client := core.ExtractClient(ctx)
			api := vpcgw.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			if request.Zone == scw.Zone(core.AllLocalities) {
				opts = append(opts, scw.WithZones(api.Zones()...))
				request.Zone = ""
			}
			resp, err := api.ListDHCPs(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.Dhcps, nil
		},
		View: &core.View{Fields: []*core.ViewField{
			{
				FieldName: "ID",
			},
			{
				FieldName: "Subnet",
			},
			{
				FieldName: "Address",
			},
			{
				FieldName: "EnableDynamic",
			},
			{
				FieldName: "PoolLow",
			},
			{
				FieldName: "PoolHigh",
			},
			{
				FieldName: "PushDefaultRoute",
			},
			{
				FieldName: "PushDNSServer",
			},
			{
				FieldName: "DNSLocalName",
			},
			{
				FieldName: "DNSServersOverride",
			},
			{
				FieldName: "DNSSearch",
			},
			{
				FieldName: "ValidLifetime",
			},
			{
				FieldName: "RenewTimer",
			},
			{
				FieldName: "RebindTimer",
			},
			{
				FieldName: "UpdatedAt",
			},
			{
				FieldName: "CreatedAt",
			},
			{
				FieldName: "Zone",
			},
			{
				FieldName: "ProjectID",
			},
			{
				FieldName: "OrganizationID",
			},
		}},
	}
}

func vpcGwDHCPGet() *core.Command {
	return &core.Command{
		Short:     `Get a DHCP configuration`,
		Long:      `Get a DHCP configuration object, identified by its DHCP ID. The response object contains configuration settings for the assignment of IP addresses to devices on a Private Network attached to a Public Gateway. Note that the response does not contain the IDs of any Private Network / Public Gateway the configuration is attached to. Use the ` + "`" + `List Public Gateway connections to Private Networks` + "`" + ` method for that purpose, filtering on DHCP ID.`,
		Namespace: "vpc-gw",
		Resource:  "dhcp",
		Verb:      "get",
		// Deprecated:    true,
		ArgsType: reflect.TypeOf(vpcgw.GetDHCPRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "dhcp-id",
				Short:      `ID of the DHCP configuration to fetch`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
			),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*vpcgw.GetDHCPRequest)

			client := core.ExtractClient(ctx)
			api := vpcgw.NewAPI(client)

			return api.GetDHCP(request)
		},
	}
}

func vpcGwDHCPCreate() *core.Command {
	return &core.Command{
		Short:     `Create a DHCP configuration`,
		Long:      `Create a new DHCP configuration object, containing settings for the assignment of IP addresses to devices on a Private Network attached to a Public Gateway. The response object includes the ID of the DHCP configuration object. You can use this ID as part of a call to ` + "`" + `Create a Public Gateway connection to a Private Network` + "`" + ` or ` + "`" + `Update a Public Gateway connection to a Private Network` + "`" + ` to directly apply this DHCP configuration.`,
		Namespace: "vpc-gw",
		Resource:  "dhcp",
		Verb:      "create",
		// Deprecated:    true,
		ArgsType: reflect.TypeOf(vpcgw.CreateDHCPRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ProjectIDArgSpec(),
			{
				Name:       "subnet",
				Short:      `Subnet for the DHCP server`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "address",
				Short:      `IP address of the DHCP server. This will be the gateway's address in the Private Network. Defaults to the first address of the subnet`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "pool-low",
				Short:      `Low IP (inclusive) of the dynamic address pool. Must be in the config's subnet. Defaults to the second address of the subnet`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "pool-high",
				Short:      `High IP (inclusive) of the dynamic address pool. Must be in the config's subnet. Defaults to the last address of the subnet`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "enable-dynamic",
				Short:      `Defines whether to enable dynamic pooling of IPs. When false, only pre-existing DHCP reservations will be handed out. Defaults to true`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "valid-lifetime",
				Short:      `How long DHCP entries will be valid for. Defaults to 1h (3600s)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "renew-timer",
				Short:      `After how long a renew will be attempted. Must be 30s lower than ` + "`" + `rebind_timer` + "`" + `. Defaults to 50m (3000s)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "rebind-timer",
				Short:      `After how long a DHCP client will query for a new lease if previous renews fail. Must be 30s lower than ` + "`" + `valid_lifetime` + "`" + `. Defaults to 51m (3060s)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "push-default-route",
				Short:      `Defines whether the gateway should push a default route to DHCP clients or only hand out IPs. Defaults to true`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "push-dns-server",
				Short:      `Defines whether the gateway should push custom DNS servers to clients. This allows for Instance hostname -> IP resolution. Defaults to true`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "dns-servers-override.{index}",
				Short:      `Array of DNS server IP addresses used to override the DNS server list pushed to DHCP clients, instead of the gateway itself`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "dns-search.{index}",
				Short:      `Array of search paths in addition to the pushed DNS configuration`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "dns-local-name",
				Short:      `TLD given to hostnames in the Private Network. Allowed characters are ` + "`" + `a-z0-9-.` + "`" + `. Defaults to the slugified Private Network name if created along a GatewayNetwork, or else to ` + "`" + `priv` + "`" + ``,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
			),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*vpcgw.CreateDHCPRequest)

			client := core.ExtractClient(ctx)
			api := vpcgw.NewAPI(client)

			return api.CreateDHCP(request)
		},
	}
}

func vpcGwDHCPUpdate() *core.Command {
	return &core.Command{
		Short:     `Update a DHCP configuration`,
		Long:      `Update a DHCP configuration object, identified by its DHCP ID.`,
		Namespace: "vpc-gw",
		Resource:  "dhcp",
		Verb:      "update",
		// Deprecated:    true,
		ArgsType: reflect.TypeOf(vpcgw.UpdateDHCPRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "dhcp-id",
				Short:      `DHCP configuration to update`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "subnet",
				Short:      `Subnet for the DHCP server`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "address",
				Short:      `IP address of the DHCP server. This will be the Public Gateway's address in the Private Network. It must be part of config's subnet`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "pool-low",
				Short:      `Low IP (inclusive) of the dynamic address pool. Must be in the config's subnet`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "pool-high",
				Short:      `High IP (inclusive) of the dynamic address pool. Must be in the config's subnet`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "enable-dynamic",
				Short:      `Defines whether to enable dynamic pooling of IPs. When false, only pre-existing DHCP reservations will be handed out. Defaults to true`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "valid-lifetime",
				Short:      `How long DHCP entries will be valid for`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "renew-timer",
				Short:      `After how long a renew will be attempted. Must be 30s lower than ` + "`" + `rebind_timer` + "`" + ``,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "rebind-timer",
				Short:      `After how long a DHCP client will query for a new lease if previous renews fail. Must be 30s lower than ` + "`" + `valid_lifetime` + "`" + ``,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "push-default-route",
				Short:      `Defines whether the gateway should push a default route to DHCP clients, or only hand out IPs`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "push-dns-server",
				Short:      `Defines whether the gateway should push custom DNS servers to clients. This allows for instance hostname -> IP resolution`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "dns-servers-override.{index}",
				Short:      `Array of DNS server IP addresses used to override the DNS server list pushed to DHCP clients, instead of the gateway itself`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "dns-search.{index}",
				Short:      `Array of search paths in addition to the pushed DNS configuration`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "dns-local-name",
				Short:      `TLD given to hostnames in the Private Networks. If an instance with hostname ` + "`" + `foo` + "`" + ` gets a lease, and this is set to ` + "`" + `bar` + "`" + `, ` + "`" + `foo.bar` + "`" + ` will resolve. Allowed characters are ` + "`" + `a-z0-9-.` + "`" + ``,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
			),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*vpcgw.UpdateDHCPRequest)

			client := core.ExtractClient(ctx)
			api := vpcgw.NewAPI(client)

			return api.UpdateDHCP(request)
		},
	}
}

func vpcGwDHCPDelete() *core.Command {
	return &core.Command{
		Short:     `Delete a DHCP configuration`,
		Long:      `Delete a DHCP configuration object, identified by its DHCP ID. Note that you cannot delete a DHCP configuration object that is currently being used by a Gateway Network.`,
		Namespace: "vpc-gw",
		Resource:  "dhcp",
		Verb:      "delete",
		// Deprecated:    true,
		ArgsType: reflect.TypeOf(vpcgw.DeleteDHCPRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "dhcp-id",
				Short:      `DHCP configuration ID to delete`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
			),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*vpcgw.DeleteDHCPRequest)

			client := core.ExtractClient(ctx)
			api := vpcgw.NewAPI(client)
			e = api.DeleteDHCP(request)
			if e != nil {
				return nil, e
			}

			return &core.SuccessResult{
				Resource: "dhcp",
				Verb:     "delete",
			}, nil
		},
	}
}

func vpcGwDHCPEntryList() *core.Command {
	return &core.Command{
		Short:     `List DHCP entries`,
		Long:      `List DHCP entries, whether dynamically assigned and/or statically reserved. DHCP entries can be filtered by the Gateway Network they are on, their MAC address, IP address, type or hostname.`,
		Namespace: "vpc-gw",
		Resource:  "dhcp-entry",
		Verb:      "list",
		// Deprecated:    true,
		ArgsType: reflect.TypeOf(vpcgw.ListDHCPEntriesRequest{}),
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
					"ip_address_asc",
					"ip_address_desc",
					"hostname_asc",
					"hostname_desc",
				},
			},
			{
				Name:       "gateway-network-id",
				Short:      `Filter for entries on this GatewayNetwork`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "mac-address",
				Short:      `Filter for entries with this MAC address`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "ip-address",
				Short:      `Filter for entries with this IP address`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "hostname",
				Short:      `Filter for entries with this hostname substring`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "type",
				Short:      `Filter for entries of this type`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown",
					"reservation",
					"lease",
				},
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
				scw.Zone(core.AllLocalities),
			),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*vpcgw.ListDHCPEntriesRequest)

			client := core.ExtractClient(ctx)
			api := vpcgw.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			if request.Zone == scw.Zone(core.AllLocalities) {
				opts = append(opts, scw.WithZones(api.Zones()...))
				request.Zone = ""
			}
			resp, err := api.ListDHCPEntries(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.DHCPEntries, nil
		},
		View: &core.View{Fields: []*core.ViewField{
			{
				FieldName: "ID",
			},
			{
				FieldName: "GatewayNetworkID",
			},
			{
				FieldName: "IPAddress",
			},
			{
				FieldName: "MacAddress",
			},
			{
				FieldName: "Hostname",
			},
			{
				FieldName: "Type",
			},
			{
				FieldName: "UpdatedAt",
			},
			{
				FieldName: "CreatedAt",
			},
			{
				FieldName: "Zone",
			},
		}},
	}
}

func vpcGwDHCPEntryGet() *core.Command {
	return &core.Command{
		Short:     `Get a DHCP entry`,
		Long:      `Get a DHCP entry, specified by its DHCP entry ID.`,
		Namespace: "vpc-gw",
		Resource:  "dhcp-entry",
		Verb:      "get",
		// Deprecated:    true,
		ArgsType: reflect.TypeOf(vpcgw.GetDHCPEntryRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "dhcp-entry-id",
				Short:      `ID of the DHCP entry to fetch`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
			),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*vpcgw.GetDHCPEntryRequest)

			client := core.ExtractClient(ctx)
			api := vpcgw.NewAPI(client)

			return api.GetDHCPEntry(request)
		},
	}
}

func vpcGwDHCPEntryCreate() *core.Command {
	return &core.Command{
		Short:     `Create a DHCP entry`,
		Long:      `Create a static DHCP reservation, specifying the Gateway Network for the reservation, the MAC address of the target device and the IP address to assign this device. The response is a DHCP entry object, confirming the ID and configuration details of the static DHCP reservation.`,
		Namespace: "vpc-gw",
		Resource:  "dhcp-entry",
		Verb:      "create",
		// Deprecated:    true,
		ArgsType: reflect.TypeOf(vpcgw.CreateDHCPEntryRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "gateway-network-id",
				Short:      `GatewayNetwork on which to create a DHCP reservation`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "mac-address",
				Short:      `MAC address to give a static entry to`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "ip-address",
				Short:      `IP address to give to the device`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
			),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*vpcgw.CreateDHCPEntryRequest)

			client := core.ExtractClient(ctx)
			api := vpcgw.NewAPI(client)

			return api.CreateDHCPEntry(request)
		},
	}
}

func vpcGwDHCPEntryUpdate() *core.Command {
	return &core.Command{
		Short:     `Update a DHCP entry`,
		Long:      `Update the IP address for a DHCP entry, specified by its DHCP entry ID. You can update an existing DHCP entry of any type (` + "`" + `reservation` + "`" + ` (static), ` + "`" + `lease` + "`" + ` (dynamic) or ` + "`" + `unknown` + "`" + `), but in manually updating the IP address the entry will necessarily be of type ` + "`" + `reservation` + "`" + ` after the update.`,
		Namespace: "vpc-gw",
		Resource:  "dhcp-entry",
		Verb:      "update",
		// Deprecated:    true,
		ArgsType: reflect.TypeOf(vpcgw.UpdateDHCPEntryRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "dhcp-entry-id",
				Short:      `ID of the DHCP entry to update`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "ip-address",
				Short:      `New IP address to give to the device`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
			),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*vpcgw.UpdateDHCPEntryRequest)

			client := core.ExtractClient(ctx)
			api := vpcgw.NewAPI(client)

			return api.UpdateDHCPEntry(request)
		},
	}
}

func vpcGwDHCPEntrySet() *core.Command {
	return &core.Command{
		Short:     `Set all DHCP reservations on a Gateway Network`,
		Long:      `Set the list of DHCP reservations attached to a Gateway Network. Reservations are identified by their MAC address, and will sync the current DHCP entry list to the given list, creating, updating or deleting DHCP entries accordingly.`,
		Namespace: "vpc-gw",
		Resource:  "dhcp-entry",
		Verb:      "set",
		// Deprecated:    true,
		ArgsType: reflect.TypeOf(vpcgw.SetDHCPEntriesRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "gateway-network-id",
				Short:      `ID of the Gateway Network on which to set DHCP reservation list`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "dhcp-entries.{index}.mac-address",
				Short:      `MAC address to give a static entry to`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "dhcp-entries.{index}.ip-address",
				Short:      `IP address to give to the device`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
			),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*vpcgw.SetDHCPEntriesRequest)

			client := core.ExtractClient(ctx)
			api := vpcgw.NewAPI(client)

			return api.SetDHCPEntries(request)
		},
	}
}

func vpcGwDHCPEntryDelete() *core.Command {
	return &core.Command{
		Short:     `Delete a DHCP entry`,
		Long:      `Delete a static DHCP reservation, identified by its DHCP entry ID. Note that you cannot delete DHCP entries of type ` + "`" + `lease` + "`" + `, these are deleted automatically when their time-to-live expires.`,
		Namespace: "vpc-gw",
		Resource:  "dhcp-entry",
		Verb:      "delete",
		// Deprecated:    true,
		ArgsType: reflect.TypeOf(vpcgw.DeleteDHCPEntryRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "dhcp-entry-id",
				Short:      `ID of the DHCP entry to delete`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
			),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*vpcgw.DeleteDHCPEntryRequest)

			client := core.ExtractClient(ctx)
			api := vpcgw.NewAPI(client)
			e = api.DeleteDHCPEntry(request)
			if e != nil {
				return nil, e
			}

			return &core.SuccessResult{
				Resource: "dhcp-entry",
				Verb:     "delete",
			}, nil
		},
	}
}

func vpcGwPatRuleList() *core.Command {
	return &core.Command{
		Short:     `List PAT rules`,
		Long:      `List PAT rules. You can filter by gateway ID to list all PAT rules for a particular gateway, or filter for PAT rules targeting a specific IP address or using a specific protocol.`,
		Namespace: "vpc-gw",
		Resource:  "pat-rule",
		Verb:      "list",
		// Deprecated:    true,
		ArgsType: reflect.TypeOf(vpcgw.ListPATRulesRequest{}),
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
					"public_port_asc",
					"public_port_desc",
				},
			},
			{
				Name:       "gateway-id",
				Short:      `Filter for PAT rules on this Gateway`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "private-ip",
				Short:      `Filter for PAT rules targeting this private ip`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "protocol",
				Short:      `Filter for PAT rules with this protocol`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown",
					"both",
					"tcp",
					"udp",
				},
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
				scw.Zone(core.AllLocalities),
			),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*vpcgw.ListPATRulesRequest)

			client := core.ExtractClient(ctx)
			api := vpcgw.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			if request.Zone == scw.Zone(core.AllLocalities) {
				opts = append(opts, scw.WithZones(api.Zones()...))
				request.Zone = ""
			}
			resp, err := api.ListPATRules(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.PatRules, nil
		},
		View: &core.View{Fields: []*core.ViewField{
			{
				FieldName: "ID",
			},
			{
				FieldName: "GatewayID",
			},
			{
				FieldName: "PublicPort",
			},
			{
				FieldName: "PrivateIP",
			},
			{
				FieldName: "PrivatePort",
			},
			{
				FieldName: "Protocol",
			},
			{
				FieldName: "UpdatedAt",
			},
			{
				FieldName: "CreatedAt",
			},
			{
				FieldName: "Zone",
			},
		}},
	}
}

func vpcGwPatRuleGet() *core.Command {
	return &core.Command{
		Short:     `Get a PAT rule`,
		Long:      `Get a PAT rule, specified by its PAT rule ID. The response object gives full details of the PAT rule, including the Public Gateway it belongs to and the configuration settings in terms of public / private ports, private IP and protocol.`,
		Namespace: "vpc-gw",
		Resource:  "pat-rule",
		Verb:      "get",
		// Deprecated:    true,
		ArgsType: reflect.TypeOf(vpcgw.GetPATRuleRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "pat-rule-id",
				Short:      `ID of the PAT rule to get`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
			),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*vpcgw.GetPATRuleRequest)

			client := core.ExtractClient(ctx)
			api := vpcgw.NewAPI(client)

			return api.GetPATRule(request)
		},
	}
}

func vpcGwPatRuleCreate() *core.Command {
	return &core.Command{
		Short:     `Create a PAT rule`,
		Long:      `Create a new PAT rule on a specified Public Gateway, defining the protocol to use, public port to listen on, and private port / IP address to map to.`,
		Namespace: "vpc-gw",
		Resource:  "pat-rule",
		Verb:      "create",
		// Deprecated:    true,
		ArgsType: reflect.TypeOf(vpcgw.CreatePATRuleRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "gateway-id",
				Short:      `ID of the Gateway on which to create the rule`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "public-port",
				Short:      `Public port to listen on`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "private-ip",
				Short:      `Private IP to forward data to`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "private-port",
				Short:      `Private port to translate to`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "protocol",
				Short:      `Protocol the rule should apply to`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown",
					"both",
					"tcp",
					"udp",
				},
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
			),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*vpcgw.CreatePATRuleRequest)

			client := core.ExtractClient(ctx)
			api := vpcgw.NewAPI(client)

			return api.CreatePATRule(request)
		},
	}
}

func vpcGwPatRuleUpdate() *core.Command {
	return &core.Command{
		Short:     `Update a PAT rule`,
		Long:      `Update a PAT rule, specified by its PAT rule ID. Configuration settings including private/public port, private IP address and protocol can all be updated.`,
		Namespace: "vpc-gw",
		Resource:  "pat-rule",
		Verb:      "update",
		// Deprecated:    true,
		ArgsType: reflect.TypeOf(vpcgw.UpdatePATRuleRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "pat-rule-id",
				Short:      `ID of the PAT rule to update`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "public-port",
				Short:      `Public port to listen on`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "private-ip",
				Short:      `Private IP to forward data to`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "private-port",
				Short:      `Private port to translate to`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "protocol",
				Short:      `Protocol the rule should apply to`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown",
					"both",
					"tcp",
					"udp",
				},
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
			),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*vpcgw.UpdatePATRuleRequest)

			client := core.ExtractClient(ctx)
			api := vpcgw.NewAPI(client)

			return api.UpdatePATRule(request)
		},
	}
}

func vpcGwPatRuleSet() *core.Command {
	return &core.Command{
		Short:     `Set all PAT rules`,
		Long:      `Set a definitive list of PAT rules attached to a Public Gateway. Each rule is identified by its public port and protocol. This will sync the current PAT rule list on the gateway with the new list, creating, updating or deleting PAT rules accordingly.`,
		Namespace: "vpc-gw",
		Resource:  "pat-rule",
		Verb:      "set",
		// Deprecated:    true,
		ArgsType: reflect.TypeOf(vpcgw.SetPATRulesRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "gateway-id",
				Short:      `ID of the gateway on which to set the PAT rules`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "pat-rules.{index}.public-port",
				Short:      `Public port to listen on`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "pat-rules.{index}.private-ip",
				Short:      `Private IP to forward data to`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "pat-rules.{index}.private-port",
				Short:      `Private port to translate to`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "pat-rules.{index}.protocol",
				Short:      `Protocol the rule should apply to`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown",
					"both",
					"tcp",
					"udp",
				},
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
			),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*vpcgw.SetPATRulesRequest)

			client := core.ExtractClient(ctx)
			api := vpcgw.NewAPI(client)

			return api.SetPATRules(request)
		},
	}
}

func vpcGwPatRuleDelete() *core.Command {
	return &core.Command{
		Short:     `Delete a PAT rule`,
		Long:      `Delete a PAT rule, identified by its PAT rule ID. This action is irreversible.`,
		Namespace: "vpc-gw",
		Resource:  "pat-rule",
		Verb:      "delete",
		// Deprecated:    true,
		ArgsType: reflect.TypeOf(vpcgw.DeletePATRuleRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "pat-rule-id",
				Short:      `ID of the PAT rule to delete`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
			),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*vpcgw.DeletePATRuleRequest)

			client := core.ExtractClient(ctx)
			api := vpcgw.NewAPI(client)
			e = api.DeletePATRule(request)
			if e != nil {
				return nil, e
			}

			return &core.SuccessResult{
				Resource: "pat-rule",
				Verb:     "delete",
			}, nil
		},
	}
}

func vpcGwGatewayTypeList() *core.Command {
	return &core.Command{
		Short:     `List Public Gateway types`,
		Long:      `List the different Public Gateway commercial offer types available at Scaleway. The response is an array of objects describing the name and technical details of each available gateway type.`,
		Namespace: "vpc-gw",
		Resource:  "gateway-type",
		Verb:      "list",
		// Deprecated:    true,
		ArgsType: reflect.TypeOf(vpcgw.ListGatewayTypesRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
			),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*vpcgw.ListGatewayTypesRequest)

			client := core.ExtractClient(ctx)
			api := vpcgw.NewAPI(client)

			return api.ListGatewayTypes(request)
		},
	}
}

func vpcGwIPList() *core.Command {
	return &core.Command{
		Short:     `List IPs`,
		Long:      `List Public Gateway flexible IP addresses. A number of filter options are available for limiting results in the response.`,
		Namespace: "vpc-gw",
		Resource:  "ip",
		Verb:      "list",
		// Deprecated:    true,
		ArgsType: reflect.TypeOf(vpcgw.ListIPsRequest{}),
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
					"ip_asc",
					"ip_desc",
					"reverse_asc",
					"reverse_desc",
				},
			},
			{
				Name:       "project-id",
				Short:      `Filter for IP addresses in this Project`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `Filter for IP addresses with these tags`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "reverse",
				Short:      `Filter for IP addresses that have a reverse containing this string`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "is-free",
				Short:      `Filter based on whether the IP is attached to a gateway or not`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "organization-id",
				Short:      `Filter for IP addresses in this Organization`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
				scw.Zone(core.AllLocalities),
			),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*vpcgw.ListIPsRequest)

			client := core.ExtractClient(ctx)
			api := vpcgw.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			if request.Zone == scw.Zone(core.AllLocalities) {
				opts = append(opts, scw.WithZones(api.Zones()...))
				request.Zone = ""
			}
			resp, err := api.ListIPs(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.IPs, nil
		},
		View: &core.View{Fields: []*core.ViewField{
			{
				FieldName: "ID",
			},
			{
				FieldName: "Address",
			},
			{
				FieldName: "Reverse",
			},
			{
				FieldName: "Tags",
			},
			{
				FieldName: "UpdatedAt",
			},
			{
				FieldName: "CreatedAt",
			},
			{
				FieldName: "Zone",
			},
			{
				FieldName: "ProjectID",
			},
			{
				FieldName: "OrganizationID",
			},
		}},
	}
}

func vpcGwIPGet() *core.Command {
	return &core.Command{
		Short:     `Get an IP`,
		Long:      `Get details of a Public Gateway flexible IP address, identified by its IP ID. The response object contains information including which (if any) Public Gateway using this IP address, the reverse and various other metadata.`,
		Namespace: "vpc-gw",
		Resource:  "ip",
		Verb:      "get",
		// Deprecated:    true,
		ArgsType: reflect.TypeOf(vpcgw.GetIPRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "ip-id",
				Short:      `ID of the IP address to get`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
			),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*vpcgw.GetIPRequest)

			client := core.ExtractClient(ctx)
			api := vpcgw.NewAPI(client)

			return api.GetIP(request)
		},
	}
}

func vpcGwIPCreate() *core.Command {
	return &core.Command{
		Short:     `Reserve an IP`,
		Long:      `Create (reserve) a new flexible IP address that can be used for a Public Gateway in a specified Scaleway Project.`,
		Namespace: "vpc-gw",
		Resource:  "ip",
		Verb:      "create",
		// Deprecated:    true,
		ArgsType: reflect.TypeOf(vpcgw.CreateIPRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ProjectIDArgSpec(),
			{
				Name:       "tags.{index}",
				Short:      `Tags to give to the IP address`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
			),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*vpcgw.CreateIPRequest)

			client := core.ExtractClient(ctx)
			api := vpcgw.NewAPI(client)

			return api.CreateIP(request)
		},
	}
}

func vpcGwIPUpdate() *core.Command {
	return &core.Command{
		Short:     `Update an IP`,
		Long:      `Update details of an existing flexible IP address, including its tags, reverse and the Public Gateway it is assigned to.`,
		Namespace: "vpc-gw",
		Resource:  "ip",
		Verb:      "update",
		// Deprecated:    true,
		ArgsType: reflect.TypeOf(vpcgw.UpdateIPRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "ip-id",
				Short:      `ID of the IP address to update`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "tags.{index}",
				Short:      `Tags to give to the IP address`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "reverse",
				Short:      `Reverse to set on the address. Empty string to unset`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "gateway-id",
				Short:      `Gateway to attach the IP address to. Empty string to detach`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
			),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*vpcgw.UpdateIPRequest)

			client := core.ExtractClient(ctx)
			api := vpcgw.NewAPI(client)

			return api.UpdateIP(request)
		},
	}
}

func vpcGwIPDelete() *core.Command {
	return &core.Command{
		Short:     `Delete an IP`,
		Long:      `Delete a flexible IP address from your account. This action is irreversible.`,
		Namespace: "vpc-gw",
		Resource:  "ip",
		Verb:      "delete",
		// Deprecated:    true,
		ArgsType: reflect.TypeOf(vpcgw.DeleteIPRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "ip-id",
				Short:      `ID of the IP address to delete`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
			),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*vpcgw.DeleteIPRequest)

			client := core.ExtractClient(ctx)
			api := vpcgw.NewAPI(client)
			e = api.DeleteIP(request)
			if e != nil {
				return nil, e
			}

			return &core.SuccessResult{
				Resource: "ip",
				Verb:     "delete",
			}, nil
		},
	}
}

func vpcGwGatewayRefreshSSHKeys() *core.Command {
	return &core.Command{
		Short:     `Refresh a Public Gateway's SSH keys`,
		Long:      `Refresh the SSH keys of a given Public Gateway, specified by its gateway ID. This adds any new SSH keys in the gateway's Scaleway Project to the gateway itself.`,
		Namespace: "vpc-gw",
		Resource:  "gateway",
		Verb:      "refresh-ssh-keys",
		// Deprecated:    true,
		ArgsType: reflect.TypeOf(vpcgw.RefreshSSHKeysRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "gateway-id",
				Short:      `ID of the gateway to refresh SSH keys on`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
			),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*vpcgw.RefreshSSHKeysRequest)

			client := core.ExtractClient(ctx)
			api := vpcgw.NewAPI(client)

			return api.RefreshSSHKeys(request)
		},
	}
}

func vpcGwGatewayMigrateToV2() *core.Command {
	return &core.Command{
		Short:     `Put a Public Gateway in IPAM mode`,
		Long:      `Put a Public Gateway in IPAM mode, so that it can be used with the Public Gateways API v2. This call is idempotent.`,
		Namespace: "vpc-gw",
		Resource:  "gateway",
		Verb:      "migrate-to-v2",
		// Deprecated:    true,
		ArgsType: reflect.TypeOf(vpcgw.MigrateToV2Request{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "gateway-id",
				Short:      `ID of the gateway to put into IPAM mode`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZoneNlAms3,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
				scw.ZonePlWaw3,
			),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*vpcgw.MigrateToV2Request)

			client := core.ExtractClient(ctx)
			api := vpcgw.NewAPI(client)
			e = api.MigrateToV2(request)
			if e != nil {
				return nil, e
			}

			return &core.SuccessResult{
				Resource: "gateway",
				Verb:     "migrate-to-v2",
			}, nil
		},
	}
}
