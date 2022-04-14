// This file was automatically generated. DO NOT EDIT.
// If you have any remark or suggestion do not hesitate to open an issue.

package vpcgw

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/internal/core"
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
	)
}
func vpcGwRoot() *core.Command {
	return &core.Command{
		Short:     `VPC Public Gateway API`,
		Long:      ``,
		Namespace: "vpc-gw",
	}
}

func vpcGwGateway() *core.Command {
	return &core.Command{
		Short: `VPC Public Gateway management`,
		Long: `The VPC Public Gateway is a building block for your infrastructure on Scaleway's shared public cloud. It provides a set of managed network services and features for Scaleway's Private Networks such as DHCP, NAT and routing.
`,
		Namespace: "vpc-gw",
		Resource:  "gateway",
	}
}

func vpcGwGatewayNetwork() *core.Command {
	return &core.Command{
		Short: `Gateway Networks management`,
		Long: `A Gateway Network represents the connection of a Private Network to a VPC Public Gateway. It holds configuration options relative to this specific connection, such as the DHCP configuration.
`,
		Namespace: "vpc-gw",
		Resource:  "gateway-network",
	}
}

func vpcGwDHCP() *core.Command {
	return &core.Command{
		Short: `DHCP configuration management`,
		Long: `DHCP configuration allows you to set parameters for assignment of IP addresses to devices on a Private Network attached to a VPC Public Gateway (subnet, lease time etc).
`,
		Namespace: "vpc-gw",
		Resource:  "dhcp",
	}
}

func vpcGwDHCPEntry() *core.Command {
	return &core.Command{
		Short: `DHCP entries management`,
		Long: `DHCP entries hold both dynamic DHCP leases (IP addresses dynamically assigned by the gateway to instances) and static user-created DHCP reservations.
`,
		Namespace: "vpc-gw",
		Resource:  "dhcp-entry",
	}
}

func vpcGwPatRule() *core.Command {
	return &core.Command{
		Short: `PAT rules management`,
		Long: `PAT (Port Address Translation) rules are global to a gateway. They define the forwarding of a public port to a specific instance on a Private Network.
`,
		Namespace: "vpc-gw",
		Resource:  "pat-rule",
	}
}

func vpcGwIP() *core.Command {
	return &core.Command{
		Short: `IP address management`,
		Long: `A VPC Public Gateway has a public IP address, allowing it to reach the public internet, as well as forward (masquerade) traffic from member instances of attached Private Networks.
`,
		Namespace: "vpc-gw",
		Resource:  "ip",
	}
}

func vpcGwGatewayType() *core.Command {
	return &core.Command{
		Short: ``,
		Long: `Gateways come in multiple shapes and size, which are described by the various gateway types.
`,
		Namespace: "vpc-gw",
		Resource:  "gateway-type",
	}
}

func vpcGwGatewayList() *core.Command {
	return &core.Command{
		Short:     `List VPC Public Gateways`,
		Long:      `List VPC Public Gateways.`,
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
				EnumValues: []string{"created_at_asc", "created_at_desc", "name_asc", "name_desc", "type_asc", "type_desc", "status_asc", "status_desc"},
			},
			{
				Name:       "project-id",
				Short:      `Include only gateways in this project`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "name",
				Short:      `Filter gateways including this name`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `Filter gateways with these tags`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "type",
				Short:      `Filter gateways of this type`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "status",
				Short:      `Filter gateways in this status (unknown for any)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"unknown", "stopped", "allocating", "configuring", "running", "stopping", "failed", "deleting", "deleted", "locked"},
			},
			{
				Name:       "private-network-id",
				Short:      `Filter gateways attached to this private network`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "organization-id",
				Short:      `Include only gateways in this organization`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZonePlWaw1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*vpcgw.ListGatewaysRequest)

			client := core.ExtractClient(ctx)
			api := vpcgw.NewAPI(client)
			resp, err := api.ListGateways(request, scw.WithAllPages())
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
		Short:     `Get a VPC Public Gateway`,
		Long:      `Get a VPC Public Gateway.`,
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
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZonePlWaw1),
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
		Short:     `Create a VPC Public Gateway`,
		Long:      `Create a VPC Public Gateway.`,
		Namespace: "vpc-gw",
		Resource:  "gateway",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(vpcgw.CreateGatewayRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ProjectIDArgSpec(),
			{
				Name:       "name",
				Short:      `Name of the gateway`,
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
				Short:      `Gateway type`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				Default:    core.DefaultValueSetter("VPC-GW-S"),
			},
			{
				Name:       "upstream-dns-servers.{index}",
				Short:      `Override the gateway's default recursive DNS servers, if DNS features are enabled`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "ip-id",
				Short:      `Attach an existing IP to the gateway`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZonePlWaw1),
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
		Short:     `Update a VPC Public Gateway`,
		Long:      `Update a VPC Public Gateway.`,
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
				Short:      `Name fo the gateway`,
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
				Short:      `Override the gateway's default recursive DNS servers, if DNS features are enabled`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "enable-bastion",
				Short:      `Enable SSH bastion on the gateway`,
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
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZonePlWaw1),
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
		Short:     `Delete a VPC Public Gateway`,
		Long:      `Delete a VPC Public Gateway.`,
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
				Name:       "cleanup-dhcp",
				Short:      `Whether to cleanup attached DHCP configurations`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZonePlWaw1),
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
		Short:     `Upgrade a VPC Public Gateway to the latest version`,
		Long:      `Upgrade a VPC Public Gateway to the latest version.`,
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
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZonePlWaw1),
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
		Short:     `List gateway connections to Private Networks`,
		Long:      `List gateway connections to Private Networks.`,
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
				EnumValues: []string{"created_at_asc", "created_at_desc", "status_asc", "status_desc"},
			},
			{
				Name:       "gateway-id",
				Short:      `Filter by gateway`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "private-network-id",
				Short:      `Filter by private network`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "enable-masquerade",
				Short:      `Filter by masquerade enablement`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "dhcpid",
				Short:      `Filter by DHCP configuration`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "status",
				Short:      `Filter GatewayNetworks by this status (unknown for any)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"unknown", "created", "attaching", "configuring", "ready", "detaching", "deleted"},
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZonePlWaw1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*vpcgw.ListGatewayNetworksRequest)

			client := core.ExtractClient(ctx)
			api := vpcgw.NewAPI(client)
			resp, err := api.ListGatewayNetworks(request, scw.WithAllPages())
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
		Short:     `Get a gateway connection to a Private Network`,
		Long:      `Get a gateway connection to a Private Network.`,
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
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZonePlWaw1),
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
		Short:     `Attach a gateway to a Private Network`,
		Long:      `Attach a gateway to a Private Network.`,
		Namespace: "vpc-gw",
		Resource:  "gateway-network",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(vpcgw.CreateGatewayNetworkRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "gateway-id",
				Short:      `Gateway to connect`,
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
				Short:      `Whether to enable masquerade on this network`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "dhcpid",
				Short:      `Existing configuration`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "address",
				Short:      `Static IP address in CIDR format to to use without DHCP`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "enable-dhcp",
				Short:      `Whether to enable DHCP on this Private Network`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZonePlWaw1),
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
		Short:     `Update a gateway connection to a Private Network`,
		Long:      `Update a gateway connection to a Private Network.`,
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
				Short:      `New masquerade enablement`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "dhcpid",
				Short:      `New DHCP configuration`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "enable-dhcp",
				Short:      `Whether to enable DHCP on the connected Private Network`,
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
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZonePlWaw1),
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
		Short:     `Detach a gateway from a Private Network`,
		Long:      `Detach a gateway from a Private Network.`,
		Namespace: "vpc-gw",
		Resource:  "gateway-network",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(vpcgw.DeleteGatewayNetworkRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "gateway-network-id",
				Short:      `GatewayNetwork to delete`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "cleanup-dhcp",
				Short:      `Whether to cleanup the attached DHCP configuration`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZonePlWaw1),
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
		Long:      `List DHCP configurations.`,
		Namespace: "vpc-gw",
		Resource:  "dhcp",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(vpcgw.ListDHCPsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `Order in which to return results`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"created_at_asc", "created_at_desc", "subnet_asc", "subnet_desc"},
			},
			{
				Name:       "project-id",
				Short:      `Include only DHCPs in this project`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "address",
				Short:      `Filter on gateway address`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "has-address",
				Short:      `Filter on subnets containing address`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "organization-id",
				Short:      `Include only DHCPs in this organization`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZonePlWaw1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*vpcgw.ListDHCPsRequest)

			client := core.ExtractClient(ctx)
			api := vpcgw.NewAPI(client)
			resp, err := api.ListDHCPs(request, scw.WithAllPages())
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
		Long:      `Get a DHCP configuration.`,
		Namespace: "vpc-gw",
		Resource:  "dhcp",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(vpcgw.GetDHCPRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "dhcpid",
				Short:      `ID of the DHCP config to fetch`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZonePlWaw1),
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
		Long:      `Create a DHCP configuration.`,
		Namespace: "vpc-gw",
		Resource:  "dhcp",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(vpcgw.CreateDHCPRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ProjectIDArgSpec(),
			{
				Name:       "subnet",
				Short:      `Subnet for the DHCP server`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "address",
				Short:      `Address of the DHCP server. This will be the gateway's address in the private network. Defaults to the first address of the subnet`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "pool-low",
				Short:      `Low IP (included) of the dynamic address pool`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "pool-high",
				Short:      `High IP (included) of the dynamic address pool`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "enable-dynamic",
				Short:      `Whether to enable dynamic pooling of IPs`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "valid-lifetime.seconds",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "valid-lifetime.nanos",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "renew-timer.seconds",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "renew-timer.nanos",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "rebind-timer.seconds",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "rebind-timer.nanos",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "push-default-route",
				Short:      `Whether the gateway should push a default route to DHCP clients or only hand out IPs. Defaults to true`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "push-dns-server",
				Short:      `Whether the gateway should push custom DNS servers to clients`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "dns-servers-override.{index}",
				Short:      `Override the DNS server list pushed to DHCP clients, instead of the gateway itself`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "dns-search.{index}",
				Short:      `Additional DNS search paths`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "dns-local-name",
				Short:      `TLD given to hosts in the Private Network`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZonePlWaw1),
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
		Long:      `Update a DHCP configuration.`,
		Namespace: "vpc-gw",
		Resource:  "dhcp",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(vpcgw.UpdateDHCPRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "dhcpid",
				Short:      `DHCP config to update`,
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
				Short:      `Address of the DHCP server. This will be the gateway's address in the private network`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "pool-low",
				Short:      `Low IP (included) of the dynamic address pool`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "pool-high",
				Short:      `High IP (included) of the dynamic address pool`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "enable-dynamic",
				Short:      `Whether to enable dynamic pooling of IPs`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "valid-lifetime.seconds",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "valid-lifetime.nanos",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "renew-timer.seconds",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "renew-timer.nanos",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "rebind-timer.seconds",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "rebind-timer.nanos",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "push-default-route",
				Short:      `Whether the gateway should push a default route to DHCP clients or only hand out IPs`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "push-dns-server",
				Short:      `Whether the gateway should push custom DNS servers to clients`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "dns-servers-override.{index}",
				Short:      `Override the DNS server list pushed to DHCP clients, instead of the gateway itself`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "dns-search.{index}",
				Short:      `Additional DNS search paths`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "dns-local-name",
				Short:      `TLD given to hosts in the Private Network`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZonePlWaw1),
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
		Long:      `Delete a DHCP configuration.`,
		Namespace: "vpc-gw",
		Resource:  "dhcp",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(vpcgw.DeleteDHCPRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "dhcpid",
				Short:      `DHCP config id to delete`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZonePlWaw1),
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
		Long:      `List DHCP entries.`,
		Namespace: "vpc-gw",
		Resource:  "dhcp-entry",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(vpcgw.ListDHCPEntriesRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `Order in which to return results`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"created_at_asc", "created_at_desc", "ip_address_asc", "ip_address_desc", "hostname_asc", "hostname_desc"},
			},
			{
				Name:       "gateway-network-id",
				Short:      `Filter entries based on the gateway network they are on`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "mac-address",
				Short:      `Filter entries on their MAC address`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "ip-address",
				Short:      `Filter entries on their IP address`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "hostname",
				Short:      `Filter entries on their hostname substring`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "type",
				Short:      `Filter entries on their type`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"unknown", "reservation", "lease"},
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZonePlWaw1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*vpcgw.ListDHCPEntriesRequest)

			client := core.ExtractClient(ctx)
			api := vpcgw.NewAPI(client)
			resp, err := api.ListDHCPEntries(request, scw.WithAllPages())
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
		Short:     `Get DHCP entries`,
		Long:      `Get DHCP entries.`,
		Namespace: "vpc-gw",
		Resource:  "dhcp-entry",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(vpcgw.GetDHCPEntryRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "dhcp-entry-id",
				Short:      `ID of the DHCP entry to fetch`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZonePlWaw1),
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
		Short:     `Create a static DHCP reservation`,
		Long:      `Create a static DHCP reservation.`,
		Namespace: "vpc-gw",
		Resource:  "dhcp-entry",
		Verb:      "create",
		// Deprecated:    false,
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
				Short:      `IP address to give to the machine`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZonePlWaw1),
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
		Long:      `Update a DHCP entry.`,
		Namespace: "vpc-gw",
		Resource:  "dhcp-entry",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(vpcgw.UpdateDHCPEntryRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "dhcp-entry-id",
				Short:      `DHCP entry ID to update`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "ip-address",
				Short:      `New IP address to give to the machine`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZonePlWaw1),
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
		Short: `Set all DHCP reservations on a Gateway Network`,
		Long: `Set the list of DHCP reservations attached to a Gateway Network. Reservations are identified by their MAC address, and will sync the current DHCP entry list to the given list, creating, updating or deleting DHCP entries.
`,
		Namespace: "vpc-gw",
		Resource:  "dhcp-entry",
		Verb:      "set",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(vpcgw.SetDHCPEntriesRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "gateway-network-id",
				Short:      `Gateway Network on which to set DHCP reservation list`,
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
				Short:      `IP address to give to the machine`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZonePlWaw1),
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
		Short:     `Delete a DHCP reservation`,
		Long:      `Delete a DHCP reservation.`,
		Namespace: "vpc-gw",
		Resource:  "dhcp-entry",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(vpcgw.DeleteDHCPEntryRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "dhcp-entry-id",
				Short:      `DHCP entry ID to delete`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZonePlWaw1),
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
		Long:      `List PAT rules.`,
		Namespace: "vpc-gw",
		Resource:  "pat-rule",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(vpcgw.ListPATRulesRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `Order in which to return results`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"created_at_asc", "created_at_desc", "public_port_asc", "public_port_desc"},
			},
			{
				Name:       "gateway-id",
				Short:      `Fetch rules for this gateway`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "private-ip",
				Short:      `Fetch rules targeting this private ip`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "protocol",
				Short:      `Fetch rules for this protocol`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"unknown", "both", "tcp", "udp"},
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZonePlWaw1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*vpcgw.ListPATRulesRequest)

			client := core.ExtractClient(ctx)
			api := vpcgw.NewAPI(client)
			resp, err := api.ListPATRules(request, scw.WithAllPages())
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
		Long:      `Get a PAT rule.`,
		Namespace: "vpc-gw",
		Resource:  "pat-rule",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(vpcgw.GetPATRuleRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "pat-rule-id",
				Short:      `PAT rule to get`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZonePlWaw1),
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
		Long:      `Create a PAT rule.`,
		Namespace: "vpc-gw",
		Resource:  "pat-rule",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(vpcgw.CreatePATRuleRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "gateway-id",
				Short:      `Gateway on which to attach the rule to`,
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
				EnumValues: []string{"unknown", "both", "tcp", "udp"},
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZonePlWaw1),
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
		Long:      `Update a PAT rule.`,
		Namespace: "vpc-gw",
		Resource:  "pat-rule",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(vpcgw.UpdatePATRuleRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "pat-rule-id",
				Short:      `PAT rule to update`,
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
				EnumValues: []string{"unknown", "both", "tcp", "udp"},
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZonePlWaw1),
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
		Short: `Set all PAT rules on a Gateway`,
		Long: `Set the list of PAT rules attached to a Gateway. Rules are identified by their public port and protocol. This will sync the current PAT rule list with the givent list, creating, updating or deleting PAT rules.
`,
		Namespace: "vpc-gw",
		Resource:  "pat-rule",
		Verb:      "set",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(vpcgw.SetPATRulesRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "gateway-id",
				Short:      `Gateway on which to set the PAT rules`,
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
				EnumValues: []string{"unknown", "both", "tcp", "udp"},
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZonePlWaw1),
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
		Long:      `Delete a PAT rule.`,
		Namespace: "vpc-gw",
		Resource:  "pat-rule",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(vpcgw.DeletePATRuleRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "pat-rule-id",
				Short:      `PAT rule to delete`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZonePlWaw1),
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
		Short:     `List VPC Public Gateway types`,
		Long:      `List VPC Public Gateway types.`,
		Namespace: "vpc-gw",
		Resource:  "gateway-type",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(vpcgw.ListGatewayTypesRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZonePlWaw1),
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
		Long:      `List IPs.`,
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
				EnumValues: []string{"created_at_asc", "created_at_desc", "ip_asc", "ip_desc", "reverse_asc", "reverse_desc"},
			},
			{
				Name:       "project-id",
				Short:      `Include only IPs in this project`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `Filter IPs with these tags`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "reverse",
				Short:      `Filter by reverse containing this string`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "is-free",
				Short:      `Filter whether the IP is attached to a gateway or not`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "organization-id",
				Short:      `Include only IPs in this organization`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZonePlWaw1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*vpcgw.ListIPsRequest)

			client := core.ExtractClient(ctx)
			api := vpcgw.NewAPI(client)
			resp, err := api.ListIPs(request, scw.WithAllPages())
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
		Long:      `Get an IP.`,
		Namespace: "vpc-gw",
		Resource:  "ip",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(vpcgw.GetIPRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "ip-id",
				Short:      `ID of the IP to get`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZonePlWaw1),
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
		Long:      `Reserve an IP.`,
		Namespace: "vpc-gw",
		Resource:  "ip",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(vpcgw.CreateIPRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ProjectIDArgSpec(),
			{
				Name:       "tags.{index}",
				Short:      `Tags to give to the IP`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZonePlWaw1),
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
		Long:      `Update an IP.`,
		Namespace: "vpc-gw",
		Resource:  "ip",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(vpcgw.UpdateIPRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "ip-id",
				Short:      `ID of the IP to update`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "tags.{index}",
				Short:      `Tags to give to the IP`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "reverse",
				Short:      `Reverse to set on the IP. Empty string to unset`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "gateway-id",
				Short:      `Gateway to attach the IP to. Empty string to detach`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZonePlWaw1),
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
		Long:      `Delete an IP.`,
		Namespace: "vpc-gw",
		Resource:  "ip",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(vpcgw.DeleteIPRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "ip-id",
				Short:      `ID of the IP to delete`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.ZonePlWaw1),
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
