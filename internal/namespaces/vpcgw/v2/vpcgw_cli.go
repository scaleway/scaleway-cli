// This file was automatically generated. DO NOT EDIT.
// If you have any remark or suggestion do not hesitate to open an issue.

package vpcgw

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-sdk-go/api/vpcgw/v2"
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
		vpcGwPatRule(),
		vpcGwIP(),
		vpcGwGatewayType(),
		vpcGwGatewayList(),
		vpcGwGatewayGet(),
		vpcGwGatewayCreate(),
		vpcGwGatewayUpdate(),
		vpcGwGatewayDelete(),
		vpcGwGatewayUpgrade(),
		vpcGwGatewayNetworkList(),
		vpcGwGatewayNetworkGet(),
		vpcGwGatewayNetworkCreate(),
		vpcGwGatewayNetworkUpdate(),
		vpcGwGatewayNetworkDelete(),
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
		Long:      `Public Gateways are building blocks for your infrastructure on Scaleway's shared public cloud. They provide a set of managed network services and features for Scaleway's Private Networks such NAT and PAT rules.`,
		Namespace: "vpc-gw",
		Resource:  "gateway",
	}
}

func vpcGwGatewayNetwork() *core.Command {
	return &core.Command{
		Short:     `Gateway Networks management`,
		Long:      `A Gateway Network represents the connection of a Private Network to a Public Gateway.`,
		Namespace: "vpc-gw",
		Resource:  "gateway-network",
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
		// Deprecated:    false,
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
				Name:       "types.{index}",
				Short:      `Filter for gateways of these types`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "status.{index}",
				Short:      `Filter for gateways with these status. Use ` + "`" + `unknown` + "`" + ` to include all statuses`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_status",
					"stopped",
					"allocating",
					"configuring",
					"running",
					"stopping",
					"failed",
					"deleting",
					"locked",
				},
			},
			{
				Name:       "private-network-ids.{index}",
				Short:      `Filter for gateways attached to these Private Networks`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "include-legacy",
				Short:      `Include also legacy gateways`,
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
	}
}

func vpcGwGatewayGet() *core.Command {
	return &core.Command{
		Short:     `Get a Public Gateway`,
		Long:      `Get details of a Public Gateway, specified by its gateway ID. The response object contains full details of the gateway, including its **name**, **type**, **status** and more.`,
		Namespace: "vpc-gw",
		Resource:  "gateway",
		Verb:      "get",
		// Deprecated:    false,
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
		// Deprecated:    false,
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
		// Deprecated:    false,
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
		// Deprecated:    false,
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
				Name:       "delete-ip",
				Short:      `Defines whether the PGW's IP should be deleted`,
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

			return api.DeleteGateway(request)
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
		// Deprecated:    false,
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

func vpcGwGatewayNetworkList() *core.Command {
	return &core.Command{
		Short:     `List Public Gateway connections to Private Networks`,
		Long:      `List the connections between Public Gateways and Private Networks (a connection = a GatewayNetwork). You can choose to filter by ` + "`" + `gateway-id` + "`" + ` to list all Private Networks attached to the specified Public Gateway, or by ` + "`" + `private_network_id` + "`" + ` to list all Public Gateways attached to the specified Private Network. Other query parameters are also available. The result is an array of GatewayNetwork objects, each giving details of the connection between a given Public Gateway and a given Private Network.`,
		Namespace: "vpc-gw",
		Resource:  "gateway-network",
		Verb:      "list",
		// Deprecated:    false,
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
				Name:       "status.{index}",
				Short:      `Filter for GatewayNetworks with these status. Use ` + "`" + `unknown` + "`" + ` to include all statuses`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_status",
					"created",
					"attaching",
					"configuring",
					"ready",
					"detaching",
				},
			},
			{
				Name:       "gateway-ids.{index}",
				Short:      `Filter for GatewayNetworks connected to these gateways`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "private-network-ids.{index}",
				Short:      `Filter for GatewayNetworks connected to these Private Networks`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "masquerade-enabled",
				Short:      `Filter for GatewayNetworks with this ` + "`" + `enable_masquerade` + "`" + ` setting`,
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
	}
}

func vpcGwGatewayNetworkGet() *core.Command {
	return &core.Command{
		Short:     `Get a Public Gateway connection to a Private Network`,
		Long:      `Get details of a given connection between a Public Gateway and a Private Network (this connection = a GatewayNetwork), specified by its ` + "`" + `gateway_network_id` + "`" + `. The response object contains details of the connection including the IDs of the Public Gateway and Private Network, the dates the connection was created/updated and its configuration settings.`,
		Namespace: "vpc-gw",
		Resource:  "gateway-network",
		Verb:      "get",
		// Deprecated:    false,
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
		Long:      `Attach a specific Public Gateway to a specific Private Network (create a GatewayNetwork). You can configure parameters for the connection including whether to enable masquerade (dynamic NAT), and more.`,
		Namespace: "vpc-gw",
		Resource:  "gateway-network",
		Verb:      "create",
		// Deprecated:    false,
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
				Name:       "push-default-route",
				Short:      `Enabling the default route also enables masquerading`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "ipam-ip-id",
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
		Long:      `Update the configuration parameters of a connection between a given Public Gateway and Private Network (the connection = a GatewayNetwork). Updatable parameters include whether to enable traffic masquerade (dynamic NAT).`,
		Namespace: "vpc-gw",
		Resource:  "gateway-network",
		Verb:      "update",
		// Deprecated:    false,
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
				Name:       "push-default-route",
				Short:      `Enabling the default route also enables masquerading`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "ipam-ip-id",
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
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(vpcgw.DeleteGatewayNetworkRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "gateway-network-id",
				Short:      `ID of the GatewayNetwork to delete`,
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
			request := args.(*vpcgw.DeleteGatewayNetworkRequest)

			client := core.ExtractClient(ctx)
			api := vpcgw.NewAPI(client)

			return api.DeleteGatewayNetwork(request)
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
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(vpcgw.ListPatRulesRequest{}),
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
				Name:       "gateway-ids.{index}",
				Short:      `Filter for PAT rules on these gateways`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "private-ips.{index}",
				Short:      `Filter for PAT rules targeting these private ips`,
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
					"unknown_protocol",
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
			request := args.(*vpcgw.ListPatRulesRequest)

			client := core.ExtractClient(ctx)
			api := vpcgw.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			if request.Zone == scw.Zone(core.AllLocalities) {
				opts = append(opts, scw.WithZones(api.Zones()...))
				request.Zone = ""
			}
			resp, err := api.ListPatRules(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.PatRules, nil
		},
	}
}

func vpcGwPatRuleGet() *core.Command {
	return &core.Command{
		Short:     `Get a PAT rule`,
		Long:      `Get a PAT rule, specified by its PAT rule ID. The response object gives full details of the PAT rule, including the Public Gateway it belongs to and the configuration settings in terms of public / private ports, private IP and protocol.`,
		Namespace: "vpc-gw",
		Resource:  "pat-rule",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(vpcgw.GetPatRuleRequest{}),
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
			request := args.(*vpcgw.GetPatRuleRequest)

			client := core.ExtractClient(ctx)
			api := vpcgw.NewAPI(client)

			return api.GetPatRule(request)
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
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(vpcgw.CreatePatRuleRequest{}),
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
					"unknown_protocol",
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
			request := args.(*vpcgw.CreatePatRuleRequest)

			client := core.ExtractClient(ctx)
			api := vpcgw.NewAPI(client)

			return api.CreatePatRule(request)
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
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(vpcgw.UpdatePatRuleRequest{}),
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
					"unknown_protocol",
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
			request := args.(*vpcgw.UpdatePatRuleRequest)

			client := core.ExtractClient(ctx)
			api := vpcgw.NewAPI(client)

			return api.UpdatePatRule(request)
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
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(vpcgw.SetPatRulesRequest{}),
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
					"unknown_protocol",
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
			request := args.(*vpcgw.SetPatRulesRequest)

			client := core.ExtractClient(ctx)
			api := vpcgw.NewAPI(client)

			return api.SetPatRules(request)
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
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(vpcgw.DeletePatRuleRequest{}),
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
			request := args.(*vpcgw.DeletePatRuleRequest)

			client := core.ExtractClient(ctx)
			api := vpcgw.NewAPI(client)
			e = api.DeletePatRule(request)
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
		// Deprecated:    false,
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
		// Deprecated:    false,
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
					"address_asc",
					"address_desc",
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
	}
}

func vpcGwIPGet() *core.Command {
	return &core.Command{
		Short:     `Get an IP`,
		Long:      `Get details of a Public Gateway flexible IP address, identified by its IP ID. The response object contains information including which (if any) Public Gateway using this IP address, the reverse and various other metadata.`,
		Namespace: "vpc-gw",
		Resource:  "ip",
		Verb:      "get",
		// Deprecated:    false,
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
		// Deprecated:    false,
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
		// Deprecated:    false,
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
		// Deprecated:    false,
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
		// Deprecated:    false,
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
