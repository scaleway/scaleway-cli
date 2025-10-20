// This file was automatically generated. DO NOT EDIT.
// If you have any remark or suggestion do not hesitate to open an issue.

package dedibox

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-sdk-go/api/dedibox/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// always import dependencies
var (
	_ = scw.RegionFrPar
)

func GetGeneratedCommands() *core.Commands {
	return core.NewCommands(
		dediboxRoot(),
		dediboxServer(),
		dediboxService(),
		dediboxOffer(),
		dediboxOs(),
		dediboxBmc(),
		dediboxReverseIP(),
		dediboxRaid(),
		dediboxRescue(),
		dediboxFip(),
		dediboxBilling(),
		dediboxIPv6Block(),
		dediboxRpnInfo(),
		dediboxSan(),
		dediboxRpnV1(),
		dediboxRpnV2(),
		dediboxServerList(),
		dediboxServerGet(),
		dediboxServerListOptions(),
		dediboxServerSubscribeOption(),
		dediboxServerCreate(),
		dediboxServerSubscribeStorage(),
		dediboxServerUpdate(),
		dediboxServerReboot(),
		dediboxServerStart(),
		dediboxServerStop(),
		dediboxServerDelete(),
		dediboxServerListEvents(),
		dediboxServerListDisks(),
		dediboxServiceGet(),
		dediboxServiceDelete(),
		dediboxServiceList(),
		dediboxServerInstall(),
		dediboxServerGetInstall(),
		dediboxServerCancelInstall(),
		dediboxServerGetPartitioning(),
		dediboxBmcStart(),
		dediboxBmcGet(),
		dediboxBmcStop(),
		dediboxOfferList(),
		dediboxOfferGet(),
		dediboxOsList(),
		dediboxOsGet(),
		dediboxReverseIPUpdate(),
		dediboxFipCreate(),
		dediboxFipAttach(),
		dediboxFipDetach(),
		dediboxFipAttachMac(),
		dediboxFipDetachMac(),
		dediboxFipDelete(),
		dediboxFipList(),
		dediboxFipGet(),
		dediboxFipGetQuota(),
		dediboxRaidGet(),
		dediboxRaidUpdate(),
		dediboxRescueStart(),
		dediboxRescueGet(),
		dediboxRescueStop(),
		dediboxBillingListInvoice(),
		dediboxBillingGetInvoice(),
		dediboxBillingDownloadInvoice(),
		dediboxBillingListRefund(),
		dediboxBillingGetRefund(),
		dediboxBillingDownloadRefund(),
		dediboxBillingGetOrderCapacity(),
		dediboxIPv6BlockGetQuota(),
		dediboxIPv6BlockCreate(),
		dediboxIPv6BlockGet(),
		dediboxIPv6BlockUpdate(),
		dediboxIPv6BlockDelete(),
		dediboxIPv6BlockCreateSubnet(),
		dediboxIPv6BlockListSubnet(),
		dediboxRpnInfoList(),
		dediboxRpnInfoGet(),
		dediboxSanList(),
		dediboxSanGet(),
		dediboxSanDelete(),
		dediboxSanCreate(),
		dediboxSanListIPs(),
		dediboxSanAddIP(),
		dediboxSanRemoveIP(),
		dediboxSanListAvailableIPs(),
		dediboxRpnV1List(),
		dediboxRpnV1Get(),
		dediboxRpnV1Create(),
		dediboxRpnV1Delete(),
		dediboxRpnV1Update(),
		dediboxRpnV1ListMembers(),
		dediboxRpnV1Invite(),
		dediboxRpnV1Leave(),
		dediboxRpnV1AddMembers(),
		dediboxRpnV1DeleteMembers(),
		dediboxRpnV1ListCapableServer(),
		dediboxRpnV1ListCapableSanServer(),
		dediboxRpnV1ListInvites(),
		dediboxRpnV1AcceptInvite(),
		dediboxRpnV1RefuseInvite(),
		dediboxRpnV2List(),
		dediboxRpnV2ListMembers(),
		dediboxRpnV2Get(),
		dediboxRpnV2Create(),
		dediboxRpnV2Delete(),
		dediboxRpnV2Update(),
		dediboxRpnV2AddMembers(),
		dediboxRpnV2DeleteMembers(),
		dediboxRpnV2ListCapableResources(),
		dediboxRpnV2ListLogs(),
		dediboxRpnV2UpdateVlanMembers(),
		dediboxRpnV2EnableCompatibility(),
		dediboxRpnV2DisableCompatibility(),
	)
}

func dediboxRoot() *core.Command {
	return &core.Command{
		Short:     `Dedibox Phoenix API`,
		Long:      `Dedibox Phoenix API.`,
		Namespace: "dedibox",
	}
}

func dediboxServer() *core.Command {
	return &core.Command{
		Short:     `Baremetal server commands`,
		Long:      `Baremetal server commands.`,
		Namespace: "dedibox",
		Resource:  "server",
	}
}

func dediboxService() *core.Command {
	return &core.Command{
		Short:     `Service commands`,
		Long:      `Service commands.`,
		Namespace: "dedibox",
		Resource:  "service",
	}
}

func dediboxOffer() *core.Command {
	return &core.Command{
		Short:     `Offer commands`,
		Long:      `Offer commands.`,
		Namespace: "dedibox",
		Resource:  "offer",
	}
}

func dediboxOs() *core.Command {
	return &core.Command{
		Short:     `OS commands`,
		Long:      `OS commands.`,
		Namespace: "dedibox",
		Resource:  "os",
	}
}

func dediboxBmc() *core.Command {
	return &core.Command{
		Short:     `BMC (Baseboard Management Controller) access commands`,
		Long:      `BMC (Baseboard Management Controller) access commands.`,
		Namespace: "dedibox",
		Resource:  "bmc",
	}
}

func dediboxReverseIP() *core.Command {
	return &core.Command{
		Short:     `Reverse-IP commands`,
		Long:      `Reverse-IP commands.`,
		Namespace: "dedibox",
		Resource:  "reverse-ip",
	}
}

func dediboxRaid() *core.Command {
	return &core.Command{
		Short:     `RAID commands`,
		Long:      `RAID commands.`,
		Namespace: "dedibox",
		Resource:  "raid",
	}
}

func dediboxRescue() *core.Command {
	return &core.Command{
		Short:     `Rescue commands`,
		Long:      `Rescue commands.`,
		Namespace: "dedibox",
		Resource:  "rescue",
	}
}

func dediboxFip() *core.Command {
	return &core.Command{
		Short:     `Failover IPs commands`,
		Long:      `Failover IPs commands.`,
		Namespace: "dedibox",
		Resource:  "fip",
	}
}

func dediboxBilling() *core.Command {
	return &core.Command{
		Short:     `Billing commands`,
		Long:      `Billing commands.`,
		Namespace: "dedibox",
		Resource:  "billing",
	}
}

func dediboxIPv6Block() *core.Command {
	return &core.Command{
		Short:     `IPv6 block commands`,
		Long:      `IPv6 block commands.`,
		Namespace: "dedibox",
		Resource:  "ipv6-block",
	}
}

func dediboxRpnInfo() *core.Command {
	return &core.Command{
		Short:     `RPN's information commands`,
		Long:      `RPN's information commands.`,
		Namespace: "dedibox",
		Resource:  "rpn-info",
	}
}

func dediboxSan() *core.Command {
	return &core.Command{
		Short:     `RPN SAN (Storage Area Network) commands`,
		Long:      `RPN SAN (Storage Area Network) commands.`,
		Namespace: "dedibox",
		Resource:  "san",
	}
}

func dediboxRpnV1() *core.Command {
	return &core.Command{
		Short:     `RPN V1 commands`,
		Long:      `RPN V1 commands.`,
		Namespace: "dedibox",
		Resource:  "rpn-v1",
	}
}

func dediboxRpnV2() *core.Command {
	return &core.Command{
		Short:     ``,
		Long:      ``,
		Namespace: "dedibox",
		Resource:  "rpn-v2",
	}
}

func dediboxServerList() *core.Command {
	return &core.Command{
		Short:     `List baremetal servers for project`,
		Long:      `List baremetal servers for project.`,
		Namespace: "dedibox",
		Resource:  "server",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.ListServersRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `Order of the servers`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"created_at_asc",
					"created_at_desc",
				},
			},
			{
				Name:       "project-id",
				Short:      `Filter servers by project ID`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "search",
				Short:      `Filter servers by hostname`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneNlAms1,
				scw.Zone(core.AllLocalities),
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.ListServersRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			if request.Zone == scw.Zone(core.AllLocalities) {
				opts = append(opts, scw.WithZones(api.Zones()...))
				request.Zone = ""
			}
			resp, err := api.ListServers(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.Servers, nil
		},
	}
}

func dediboxServerGet() *core.Command {
	return &core.Command{
		Short:     `Get a specific baremetal server`,
		Long:      `Get the server associated with the given ID.`,
		Namespace: "dedibox",
		Resource:  "server",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.GetServerRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "server-id",
				Short:      `ID of the server`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneNlAms1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.GetServerRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewAPI(client)

			return api.GetServer(request)
		},
	}
}

func dediboxServerListOptions() *core.Command {
	return &core.Command{
		Short:     `List subscribable server options`,
		Long:      `List subscribable options associated to the given server ID.`,
		Namespace: "dedibox",
		Resource:  "server",
		Verb:      "list-options",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.ListSubscribableServerOptionsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "server-id",
				Short:      `Server ID of the subscribable server options`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneNlAms1,
				scw.Zone(core.AllLocalities),
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.ListSubscribableServerOptionsRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			if request.Zone == scw.Zone(core.AllLocalities) {
				opts = append(opts, scw.WithZones(api.Zones()...))
				request.Zone = ""
			}
			resp, err := api.ListSubscribableServerOptions(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.ServerOptions, nil
		},
	}
}

func dediboxServerSubscribeOption() *core.Command {
	return &core.Command{
		Short:     `Subscribe server option`,
		Long:      `Subscribe option for the given server ID.`,
		Namespace: "dedibox",
		Resource:  "server",
		Verb:      "subscribe-option",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.SubscribeServerOptionRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "server-id",
				Short:      `Server ID to subscribe server option`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "option-id",
				Short:      `Option ID to subscribe`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneNlAms1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.SubscribeServerOptionRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewAPI(client)

			return api.SubscribeServerOption(request)
		},
	}
}

func dediboxServerCreate() *core.Command {
	return &core.Command{
		Short:     `Create a baremetal server`,
		Long:      `Create a new baremetal server. The order return you a service ID to follow the provisionning status you could call GetService.`,
		Namespace: "dedibox",
		Resource:  "server",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.CreateServerRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "offer-id",
				Short:      `Offer ID of the new server`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "server-option-ids.{index}",
				Short:      `Server option IDs of the new server`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.ProjectIDArgSpec(),
			{
				Name:       "datacenter-name",
				Short:      `Datacenter name of the new server`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneNlAms1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.CreateServerRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewAPI(client)

			return api.CreateServer(request)
		},
	}
}

func dediboxServerSubscribeStorage() *core.Command {
	return &core.Command{
		Short:     `Subscribe storage server option`,
		Long:      `Subscribe storage option for the given server ID.`,
		Namespace: "dedibox",
		Resource:  "server",
		Verb:      "subscribe-storage",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.SubscribeStorageOptionsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "server-id",
				Short:      `Server ID of the storage options to subscribe`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "options-ids.{index}",
				Short:      `Option IDs of the storage options to subscribe`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneNlAms1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.SubscribeStorageOptionsRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewAPI(client)

			return api.SubscribeStorageOptions(request)
		},
	}
}

func dediboxServerUpdate() *core.Command {
	return &core.Command{
		Short:     `Update a baremetal server`,
		Long:      `Update the server associated with the given ID.`,
		Namespace: "dedibox",
		Resource:  "server",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.UpdateServerRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "server-id",
				Short:      `Server ID to update`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "hostname",
				Short:      `Hostname of the server to update`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "enable-ipv6",
				Short:      `Flag to enable or not the IPv6 of server`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneNlAms1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.UpdateServerRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewAPI(client)

			return api.UpdateServer(request)
		},
	}
}

func dediboxServerReboot() *core.Command {
	return &core.Command{
		Short:     `Reboot a baremetal server`,
		Long:      `Reboot the server associated with the given ID, use boot param to reboot in rescue.`,
		Namespace: "dedibox",
		Resource:  "server",
		Verb:      "reboot",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.RebootServerRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "server-id",
				Short:      `Server ID to reboot`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneNlAms1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.RebootServerRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewAPI(client)
			e = api.RebootServer(request)
			if e != nil {
				return nil, e
			}

			return &core.SuccessResult{
				Resource: "server",
				Verb:     "reboot",
			}, nil
		},
	}
}

func dediboxServerStart() *core.Command {
	return &core.Command{
		Short:     `Start a baremetal server`,
		Long:      `Start the server associated with the given ID.`,
		Namespace: "dedibox",
		Resource:  "server",
		Verb:      "start",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.StartServerRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "server-id",
				Short:      `Server ID to start`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneNlAms1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.StartServerRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewAPI(client)
			e = api.StartServer(request)
			if e != nil {
				return nil, e
			}

			return &core.SuccessResult{
				Resource: "server",
				Verb:     "start",
			}, nil
		},
	}
}

func dediboxServerStop() *core.Command {
	return &core.Command{
		Short:     `Stop a baremetal server`,
		Long:      `Stop the server associated with the given ID.`,
		Namespace: "dedibox",
		Resource:  "server",
		Verb:      "stop",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.StopServerRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "server-id",
				Short:      `Server ID to stop`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneNlAms1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.StopServerRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewAPI(client)
			e = api.StopServer(request)
			if e != nil {
				return nil, e
			}

			return &core.SuccessResult{
				Resource: "server",
				Verb:     "stop",
			}, nil
		},
	}
}

func dediboxServerDelete() *core.Command {
	return &core.Command{
		Short:     `Delete a baremetal server`,
		Long:      `Delete the server associated with the given ID.`,
		Namespace: "dedibox",
		Resource:  "server",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.DeleteServerRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "server-id",
				Short:      `Server ID to delete`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneNlAms1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.DeleteServerRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewAPI(client)
			e = api.DeleteServer(request)
			if e != nil {
				return nil, e
			}

			return &core.SuccessResult{
				Resource: "server",
				Verb:     "delete",
			}, nil
		},
	}
}

func dediboxServerListEvents() *core.Command {
	return &core.Command{
		Short:     `List server events`,
		Long:      `List events associated to the given server ID.`,
		Namespace: "dedibox",
		Resource:  "server",
		Verb:      "list-events",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.ListServerEventsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `Order of the server events`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"created_at_asc",
					"created_at_desc",
				},
			},
			{
				Name:       "server-id",
				Short:      `Server ID of the server events`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneNlAms1,
				scw.Zone(core.AllLocalities),
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.ListServerEventsRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			if request.Zone == scw.Zone(core.AllLocalities) {
				opts = append(opts, scw.WithZones(api.Zones()...))
				request.Zone = ""
			}
			resp, err := api.ListServerEvents(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.Events, nil
		},
	}
}

func dediboxServerListDisks() *core.Command {
	return &core.Command{
		Short:     `List server disks`,
		Long:      `List disks associated to the given server ID.`,
		Namespace: "dedibox",
		Resource:  "server",
		Verb:      "list-disks",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.ListServerDisksRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `Order of the server disks`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"created_at_asc",
					"created_at_desc",
				},
			},
			{
				Name:       "server-id",
				Short:      `Server ID of the server disks`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneNlAms1,
				scw.Zone(core.AllLocalities),
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.ListServerDisksRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			if request.Zone == scw.Zone(core.AllLocalities) {
				opts = append(opts, scw.WithZones(api.Zones()...))
				request.Zone = ""
			}
			resp, err := api.ListServerDisks(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.Disks, nil
		},
	}
}

func dediboxServiceGet() *core.Command {
	return &core.Command{
		Short:     `Get a specific service`,
		Long:      `Get the service associated with the given ID.`,
		Namespace: "dedibox",
		Resource:  "service",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.GetServiceRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "service-id",
				Short:      `ID of the service`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneNlAms1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.GetServiceRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewAPI(client)

			return api.GetService(request)
		},
	}
}

func dediboxServiceDelete() *core.Command {
	return &core.Command{
		Short:     `Delete a specific service`,
		Long:      `Delete the service associated with the given ID.`,
		Namespace: "dedibox",
		Resource:  "service",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.DeleteServiceRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "service-id",
				Short:      `ID of the service`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneNlAms1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.DeleteServiceRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewAPI(client)

			return api.DeleteService(request)
		},
	}
}

func dediboxServiceList() *core.Command {
	return &core.Command{
		Short:     `List services`,
		Long:      `List services.`,
		Namespace: "dedibox",
		Resource:  "service",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.ListServicesRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `Order of the services`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"created_at_asc",
					"created_at_desc",
				},
			},
			{
				Name:       "project-id",
				Short:      `Project ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneNlAms1,
				scw.Zone(core.AllLocalities),
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.ListServicesRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			if request.Zone == scw.Zone(core.AllLocalities) {
				opts = append(opts, scw.WithZones(api.Zones()...))
				request.Zone = ""
			}
			resp, err := api.ListServices(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.Services, nil
		},
	}
}

func dediboxServerInstall() *core.Command {
	return &core.Command{
		Short:     `Install a baremetal server`,
		Long:      `Install an OS on the server associated with the given ID.`,
		Namespace: "dedibox",
		Resource:  "server",
		Verb:      "install",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.InstallServerRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "server-id",
				Short:      `Server ID to install`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "os-id",
				Short:      `OS ID to install on the server`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "hostname",
				Short:      `Hostname of the server`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "user-login",
				Short:      `User to install on the server`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "user-password",
				Short:      `User password to install on the server`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "panel-password",
				Short:      `Panel password to install on the server`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "root-password",
				Short:      `Root password to install on the server`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "partitions.{index}.file-system",
				Short:      `File system of the installation partition`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown",
					"efi",
					"swap",
					"ext4",
					"ext3",
					"ext2",
					"xfs",
					"ntfs",
					"fat32",
					"ufs",
				},
			},
			{
				Name:       "partitions.{index}.mount-point",
				Short:      `Mount point of the installation partition`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "partitions.{index}.raid-level",
				Short:      `RAID level of the installation partition`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"no_raid",
					"raid0",
					"raid1",
					"raid5",
					"raid6",
					"raid10",
				},
			},
			{
				Name:       "partitions.{index}.capacity",
				Short:      `Capacity of the installation partition`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "partitions.{index}.connectors.{index}",
				Short:      `Connectors of the installation partition`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "ssh-key-ids.{index}",
				Short:      `SSH key IDs authorized on the server`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "license-offer-id",
				Short:      `Offer ID of license to install on server`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "ip-id",
				Short:      `IP to link at the license to install on server`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneNlAms1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.InstallServerRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewAPI(client)

			return api.InstallServer(request)
		},
	}
}

func dediboxServerGetInstall() *core.Command {
	return &core.Command{
		Short:     `Get a specific server installation status`,
		Long:      `Get the server installation status associated with the given server ID.`,
		Namespace: "dedibox",
		Resource:  "server",
		Verb:      "get-install",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.GetServerInstallRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "server-id",
				Short:      `Server ID of the server to install`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneNlAms1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.GetServerInstallRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewAPI(client)

			return api.GetServerInstall(request)
		},
	}
}

func dediboxServerCancelInstall() *core.Command {
	return &core.Command{
		Short:     `Cancels the current (running) server installation`,
		Long:      `Cancels the current server installation associated with the given server ID.`,
		Namespace: "dedibox",
		Resource:  "server",
		Verb:      "cancel-install",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.CancelServerInstallRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "server-id",
				Short:      `Server ID of the server to cancel install`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneNlAms1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.CancelServerInstallRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewAPI(client)
			e = api.CancelServerInstall(request)
			if e != nil {
				return nil, e
			}

			return &core.SuccessResult{
				Resource: "server",
				Verb:     "cancel-install",
			}, nil
		},
	}
}

func dediboxServerGetPartitioning() *core.Command {
	return &core.Command{
		Short:     `Get server default partitioning`,
		Long:      `Get the server default partitioning schema associated with the given server ID and OS ID.`,
		Namespace: "dedibox",
		Resource:  "server",
		Verb:      "get-partitioning",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.GetServerDefaultPartitioningRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "server-id",
				Short:      `ID of the server`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "os-id",
				Short:      `OS ID of the default partitioning`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneNlAms1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.GetServerDefaultPartitioningRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewAPI(client)

			return api.GetServerDefaultPartitioning(request)
		},
	}
}

func dediboxBmcStart() *core.Command {
	return &core.Command{
		Short: `Start BMC (Baseboard Management Controller) access for a given baremetal server`,
		Long: `Start BMC (Baseboard Management Controller) access associated with the given ID.
The BMC (Baseboard Management Controller) access is available one hour after the installation of the server.`,
		Namespace: "dedibox",
		Resource:  "bmc",
		Verb:      "start",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.StartBMCAccessRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "server-id",
				Short:      `ID of the server to start the BMC access`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "ip",
				Short:      `The IP authorized to connect to the given server`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneNlAms1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.StartBMCAccessRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewAPI(client)
			e = api.StartBMCAccess(request)
			if e != nil {
				return nil, e
			}

			return &core.SuccessResult{
				Resource: "bmc",
				Verb:     "start",
			}, nil
		},
	}
}

func dediboxBmcGet() *core.Command {
	return &core.Command{
		Short:     `Get BMC (Baseboard Management Controller) access for a given baremetal server`,
		Long:      `Get the BMC (Baseboard Management Controller) access associated with the given ID.`,
		Namespace: "dedibox",
		Resource:  "bmc",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.GetBMCAccessRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "server-id",
				Short:      `ID of the server to get BMC access`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneNlAms1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.GetBMCAccessRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewAPI(client)

			return api.GetBMCAccess(request)
		},
	}
}

func dediboxBmcStop() *core.Command {
	return &core.Command{
		Short:     `Stop BMC (Baseboard Management Controller) access for a given baremetal server`,
		Long:      `Stop BMC (Baseboard Management Controller) access associated with the given ID.`,
		Namespace: "dedibox",
		Resource:  "bmc",
		Verb:      "stop",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.StopBMCAccessRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "server-id",
				Short:      `ID of the server to stop BMC access`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneNlAms1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.StopBMCAccessRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewAPI(client)
			e = api.StopBMCAccess(request)
			if e != nil {
				return nil, e
			}

			return &core.SuccessResult{
				Resource: "bmc",
				Verb:     "stop",
			}, nil
		},
	}
}

func dediboxOfferList() *core.Command {
	return &core.Command{
		Short:     `List offers`,
		Long:      `List all available server offers.`,
		Namespace: "dedibox",
		Resource:  "offer",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.ListOffersRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `Order of the offers`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"created_at_asc",
					"created_at_desc",
					"name_asc",
					"name_desc",
					"price_asc",
					"price_desc",
				},
			},
			{
				Name:       "commercial-range",
				Short:      `Filter on commercial range`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "catalog",
				Short:      `Filter on catalog`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"all",
					"default",
					"beta",
					"reseller",
					"premium",
					"volume",
					"admin",
					"inactive",
				},
			},
			{
				Name:       "project-id",
				Short:      `Project ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "is-failover-ip",
				Short:      `Get the current failover IP offer`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "is-failover-block",
				Short:      `Get the current failover IP block offer`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "sold-in",
				Short:      `Filter offers depending on their datacenter`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "available-only",
				Short:      `Set this filter to true to only return available offers`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "is-rpn-san",
				Short:      `Get the RPN SAN offers`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneNlAms1,
				scw.Zone(core.AllLocalities),
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.ListOffersRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			if request.Zone == scw.Zone(core.AllLocalities) {
				opts = append(opts, scw.WithZones(api.Zones()...))
				request.Zone = ""
			}
			resp, err := api.ListOffers(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.Offers, nil
		},
	}
}

func dediboxOfferGet() *core.Command {
	return &core.Command{
		Short:     `Get offer`,
		Long:      `Return specific offer for the given ID.`,
		Namespace: "dedibox",
		Resource:  "offer",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.GetOfferRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "offer-id",
				Short:      `ID of offer`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "project-id",
				Short:      `Project ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneNlAms1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.GetOfferRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewAPI(client)

			return api.GetOffer(request)
		},
	}
}

func dediboxOsList() *core.Command {
	return &core.Command{
		Short:     `List all available OS that can be install on a baremetal server`,
		Long:      `List all available OS that can be install on a baremetal server.`,
		Namespace: "dedibox",
		Resource:  "os",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.ListOSRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `Order of the OS`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"created_at_asc",
					"created_at_desc",
					"released_at_asc",
					"released_at_desc",
				},
			},
			{
				Name:       "type",
				Short:      `Type of the OS`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_type",
					"server",
					"virtu",
					"panel",
					"desktop",
					"custom",
					"rescue",
				},
			},
			{
				Name:       "server-id",
				Short:      `Filter OS by compatible server ID`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "project-id",
				Short:      `Project ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneNlAms1,
				scw.Zone(core.AllLocalities),
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.ListOSRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			if request.Zone == scw.Zone(core.AllLocalities) {
				opts = append(opts, scw.WithZones(api.Zones()...))
				request.Zone = ""
			}
			resp, err := api.ListOS(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.Os, nil
		},
	}
}

func dediboxOsGet() *core.Command {
	return &core.Command{
		Short:     `Get an OS with a given ID`,
		Long:      `Return specific OS for the given ID.`,
		Namespace: "dedibox",
		Resource:  "os",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.GetOSRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "os-id",
				Short:      `ID of the OS`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "server-id",
				Short:      `ID of the server`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "project-id",
				Short:      `Project ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneNlAms1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.GetOSRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewAPI(client)

			return api.GetOS(request)
		},
	}
}

func dediboxReverseIPUpdate() *core.Command {
	return &core.Command{
		Short:     `Update reverse of ip`,
		Long:      `Update reverse of ip associated with the given ID.`,
		Namespace: "dedibox",
		Resource:  "reverse-ip",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.UpdateReverseRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "ip-id",
				Short:      `ID of the IP`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "reverse",
				Short:      `Reverse to apply on the IP`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneNlAms1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.UpdateReverseRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewAPI(client)

			return api.UpdateReverse(request)
		},
	}
}

func dediboxFipCreate() *core.Command {
	return &core.Command{
		Short:     `Order failover IPs`,
		Long:      `Order X failover IPs.`,
		Namespace: "dedibox",
		Resource:  "fip",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.CreateFailoverIPsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "offer-id",
				Short:      `Failover IP offer ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ProjectIDArgSpec(),
			{
				Name:       "quantity",
				Short:      `Quantity`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneNlAms1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.CreateFailoverIPsRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewAPI(client)

			return api.CreateFailoverIPs(request)
		},
	}
}

func dediboxFipAttach() *core.Command {
	return &core.Command{
		Short:     `Attach failovers on baremetal server`,
		Long:      `Attach failovers on the server associated with the given ID.`,
		Namespace: "dedibox",
		Resource:  "fip",
		Verb:      "attach",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.AttachFailoverIPsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "server-id",
				Short:      `ID of the server`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "fips-ids.{index}",
				Short:      `List of ID of failovers IP to attach`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneNlAms1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.AttachFailoverIPsRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewAPI(client)
			e = api.AttachFailoverIPs(request)
			if e != nil {
				return nil, e
			}

			return &core.SuccessResult{
				Resource: "fip",
				Verb:     "attach",
			}, nil
		},
	}
}

func dediboxFipDetach() *core.Command {
	return &core.Command{
		Short:     `Detach failovers on baremetal server`,
		Long:      `Detach failovers on the server associated with the given ID.`,
		Namespace: "dedibox",
		Resource:  "fip",
		Verb:      "detach",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.DetachFailoverIPsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "fips-ids.{index}",
				Short:      `List of IDs of failovers IP to detach`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneNlAms1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.DetachFailoverIPsRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewAPI(client)
			e = api.DetachFailoverIPs(request)
			if e != nil {
				return nil, e
			}

			return &core.SuccessResult{
				Resource: "fip",
				Verb:     "detach",
			}, nil
		},
	}
}

func dediboxFipAttachMac() *core.Command {
	return &core.Command{
		Short:     `Attach a failover IP to a MAC address`,
		Long:      `Attach a failover IP to a MAC address.`,
		Namespace: "dedibox",
		Resource:  "fip",
		Verb:      "attach-mac",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.AttachFailoverIPToMacAddressRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "ip-id",
				Short:      `ID of the failover IP`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "type",
				Short:      `A mac type`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"mac_type_unknown",
					"vmware",
					"kvm",
					"xen",
				},
			},
			{
				Name:       "mac",
				Short:      `A valid mac address (existing or not)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneNlAms1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.AttachFailoverIPToMacAddressRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewAPI(client)

			return api.AttachFailoverIPToMacAddress(request)
		},
	}
}

func dediboxFipDetachMac() *core.Command {
	return &core.Command{
		Short:     `Detach a failover IP from a MAC address`,
		Long:      `Detach a failover IP from a MAC address.`,
		Namespace: "dedibox",
		Resource:  "fip",
		Verb:      "detach-mac",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.DetachFailoverIPFromMacAddressRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "ip-id",
				Short:      `ID of the failover IP`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneNlAms1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.DetachFailoverIPFromMacAddressRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewAPI(client)

			return api.DetachFailoverIPFromMacAddress(request)
		},
	}
}

func dediboxFipDelete() *core.Command {
	return &core.Command{
		Short:     `Delete a failover server`,
		Long:      `Delete the failover associated with the given ID.`,
		Namespace: "dedibox",
		Resource:  "fip",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.DeleteFailoverIPRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "ip-id",
				Short:      `ID of the failover IP to delete`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneNlAms1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.DeleteFailoverIPRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewAPI(client)
			e = api.DeleteFailoverIP(request)
			if e != nil {
				return nil, e
			}

			return &core.SuccessResult{
				Resource: "fip",
				Verb:     "delete",
			}, nil
		},
	}
}

func dediboxFipList() *core.Command {
	return &core.Command{
		Short:     `List failovers for project`,
		Long:      `List failovers servers for project.`,
		Namespace: "dedibox",
		Resource:  "fip",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.ListFailoverIPsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `Order of the failovers IP`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"ip_asc",
					"ip_desc",
				},
			},
			{
				Name:       "project-id",
				Short:      `Filter failovers IP by project ID`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "search",
				Short:      `Filter failovers IP which matching with this field`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "only-available",
				Short:      `True: return all failovers IP not attached on server`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneNlAms1,
				scw.Zone(core.AllLocalities),
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.ListFailoverIPsRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			if request.Zone == scw.Zone(core.AllLocalities) {
				opts = append(opts, scw.WithZones(api.Zones()...))
				request.Zone = ""
			}
			resp, err := api.ListFailoverIPs(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.FailoverIPs, nil
		},
	}
}

func dediboxFipGet() *core.Command {
	return &core.Command{
		Short:     `Get a specific baremetal server`,
		Long:      `Get the server associated with the given ID.`,
		Namespace: "dedibox",
		Resource:  "fip",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.GetFailoverIPRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "ip-id",
				Short:      `ID of the failover IP`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneNlAms1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.GetFailoverIPRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewAPI(client)

			return api.GetFailoverIP(request)
		},
	}
}

func dediboxFipGetQuota() *core.Command {
	return &core.Command{
		Short:     `Get remaining quota`,
		Long:      `Get remaining quota.`,
		Namespace: "dedibox",
		Resource:  "fip",
		Verb:      "get-quota",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.GetRemainingQuotaRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "project-id",
				Short:      `Project ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneNlAms1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.GetRemainingQuotaRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewAPI(client)

			return api.GetRemainingQuota(request)
		},
	}
}

func dediboxRaidGet() *core.Command {
	return &core.Command{
		Short:     `Get raid`,
		Long:      `Return raid for the given server ID.`,
		Namespace: "dedibox",
		Resource:  "raid",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.GetRaidRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "server-id",
				Short:      `ID of the server`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneNlAms1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.GetRaidRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewAPI(client)

			return api.GetRaid(request)
		},
	}
}

func dediboxRaidUpdate() *core.Command {
	return &core.Command{
		Short:     `Update RAID`,
		Long:      `Update RAID associated with the given server ID.`,
		Namespace: "dedibox",
		Resource:  "raid",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.UpdateRaidRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "server-id",
				Short:      `ID of the server`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "raid-arrays.{index}.raid-level",
				Short:      `The RAID level`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"no_raid",
					"raid0",
					"raid1",
					"raid5",
					"raid6",
					"raid10",
				},
			},
			{
				Name:       "raid-arrays.{index}.disk-ids.{index}",
				Short:      `The list of Disk ID of the updatable RAID`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneNlAms1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.UpdateRaidRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewAPI(client)
			e = api.UpdateRaid(request)
			if e != nil {
				return nil, e
			}

			return &core.SuccessResult{
				Resource: "raid",
				Verb:     "update",
			}, nil
		},
	}
}

func dediboxRescueStart() *core.Command {
	return &core.Command{
		Short:     `Start in rescue baremetal server`,
		Long:      `Start in rescue the server associated with the given ID.`,
		Namespace: "dedibox",
		Resource:  "rescue",
		Verb:      "start",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.StartRescueRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "server-id",
				Short:      `ID of the server to start rescue`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "os-id",
				Short:      `OS ID to use to start rescue`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneNlAms1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.StartRescueRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewAPI(client)

			return api.StartRescue(request)
		},
	}
}

func dediboxRescueGet() *core.Command {
	return &core.Command{
		Short:     `Get rescue information`,
		Long:      `Return rescue information for the given server ID.`,
		Namespace: "dedibox",
		Resource:  "rescue",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.GetRescueRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "server-id",
				Short:      `ID of the server to get rescue`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneNlAms1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.GetRescueRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewAPI(client)

			return api.GetRescue(request)
		},
	}
}

func dediboxRescueStop() *core.Command {
	return &core.Command{
		Short:     `Stop rescue on baremetal server`,
		Long:      `Stop rescue on the server associated with the given ID.`,
		Namespace: "dedibox",
		Resource:  "rescue",
		Verb:      "stop",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.StopRescueRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "server-id",
				Short:      `ID of the server to stop rescue`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneNlAms1,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.StopRescueRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewAPI(client)
			e = api.StopRescue(request)
			if e != nil {
				return nil, e
			}

			return &core.SuccessResult{
				Resource: "rescue",
				Verb:     "stop",
			}, nil
		},
	}
}

func dediboxBillingListInvoice() *core.Command {
	return &core.Command{
		Short:     `List-invoice dedibox resources`,
		Long:      `List-invoice dedibox resources.`,
		Namespace: "dedibox",
		Resource:  "billing",
		Verb:      "list-invoice",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.BillingAPIListInvoicesRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"created_at_asc",
					"created_at_desc",
				},
			},
			{
				Name:       "project-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.BillingAPIListInvoicesRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewBillingAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			resp, err := api.ListInvoices(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.Invoices, nil
		},
	}
}

func dediboxBillingGetInvoice() *core.Command {
	return &core.Command{
		Short:     `Get-invoice dedibox resources`,
		Long:      `Get-invoice dedibox resources.`,
		Namespace: "dedibox",
		Resource:  "billing",
		Verb:      "get-invoice",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.BillingAPIGetInvoiceRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "invoice-id",
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.BillingAPIGetInvoiceRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewBillingAPI(client)

			return api.GetInvoice(request)
		},
	}
}

func dediboxBillingDownloadInvoice() *core.Command {
	return &core.Command{
		Short:     `Download-invoice dedibox resources`,
		Long:      `Download-invoice dedibox resources.`,
		Namespace: "dedibox",
		Resource:  "billing",
		Verb:      "download-invoice",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.BillingAPIDownloadInvoiceRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "invoice-id",
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.BillingAPIDownloadInvoiceRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewBillingAPI(client)

			return api.DownloadInvoice(request)
		},
	}
}

func dediboxBillingListRefund() *core.Command {
	return &core.Command{
		Short:     `List-refund dedibox resources`,
		Long:      `List-refund dedibox resources.`,
		Namespace: "dedibox",
		Resource:  "billing",
		Verb:      "list-refund",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.BillingAPIListRefundsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"created_at_asc",
					"created_at_desc",
				},
			},
			{
				Name:       "project-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.BillingAPIListRefundsRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewBillingAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			resp, err := api.ListRefunds(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.Refunds, nil
		},
	}
}

func dediboxBillingGetRefund() *core.Command {
	return &core.Command{
		Short:     `Get-refund dedibox resources`,
		Long:      `Get-refund dedibox resources.`,
		Namespace: "dedibox",
		Resource:  "billing",
		Verb:      "get-refund",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.BillingAPIGetRefundRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "refund-id",
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.BillingAPIGetRefundRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewBillingAPI(client)

			return api.GetRefund(request)
		},
	}
}

func dediboxBillingDownloadRefund() *core.Command {
	return &core.Command{
		Short:     `Download-refund dedibox resources`,
		Long:      `Download-refund dedibox resources.`,
		Namespace: "dedibox",
		Resource:  "billing",
		Verb:      "download-refund",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.BillingAPIDownloadRefundRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "refund-id",
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.BillingAPIDownloadRefundRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewBillingAPI(client)

			return api.DownloadRefund(request)
		},
	}
}

func dediboxBillingGetOrderCapacity() *core.Command {
	return &core.Command{
		Short:     `Get-order-capacity dedibox resources`,
		Long:      `Get-order-capacity dedibox resources.`,
		Namespace: "dedibox",
		Resource:  "billing",
		Verb:      "get-order-capacity",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.BillingAPICanOrderRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ProjectIDArgSpec(),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.BillingAPICanOrderRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewBillingAPI(client)

			return api.CanOrder(request)
		},
	}
}

func dediboxIPv6BlockGetQuota() *core.Command {
	return &core.Command{
		Short: `Get IPv6 block quota`,
		Long: `Get IPv6 block quota with the given project ID.
/48 one per organization.
/56 link to your number of server.
/64 link to your number of failover IP.`,
		Namespace: "dedibox",
		Resource:  "ipv6-block",
		Verb:      "get-quota",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.IPv6BlockAPIGetIPv6BlockQuotasRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "project-id",
				Short:      `ID of the project`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.IPv6BlockAPIGetIPv6BlockQuotasRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewIPv6BlockAPI(client)

			return api.GetIPv6BlockQuotas(request)
		},
	}
}

func dediboxIPv6BlockCreate() *core.Command {
	return &core.Command{
		Short:     `Create IPv6 block for baremetal server`,
		Long:      `Create IPv6 block associated with the given project ID.`,
		Namespace: "dedibox",
		Resource:  "ipv6-block",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.IPv6BlockAPICreateIPv6BlockRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "project-id",
				Short:      `ID of the project`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.IPv6BlockAPICreateIPv6BlockRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewIPv6BlockAPI(client)

			return api.CreateIPv6Block(request)
		},
	}
}

func dediboxIPv6BlockGet() *core.Command {
	return &core.Command{
		Short:     `Get first IPv6 block`,
		Long:      `Get the first IPv6 block associated with the given project ID.`,
		Namespace: "dedibox",
		Resource:  "ipv6-block",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.IPv6BlockAPIGetIPv6BlockRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "project-id",
				Short:      `ID of the project`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.IPv6BlockAPIGetIPv6BlockRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewIPv6BlockAPI(client)

			return api.GetIPv6Block(request)
		},
	}
}

func dediboxIPv6BlockUpdate() *core.Command {
	return &core.Command{
		Short: `Update IPv6 block`,
		Long: `Update DNS associated to IPv6 block.
If DNS is used, minimum of 2 is necessary and maximum of 5 (no duplicate).`,
		Namespace: "dedibox",
		Resource:  "ipv6-block",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.IPv6BlockAPIUpdateIPv6BlockRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "block-id",
				Short:      `ID of the IPv6 block`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "nameservers.{index}",
				Short:      `DNS to link to the IPv6`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.IPv6BlockAPIUpdateIPv6BlockRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewIPv6BlockAPI(client)

			return api.UpdateIPv6Block(request)
		},
	}
}

func dediboxIPv6BlockDelete() *core.Command {
	return &core.Command{
		Short:     `Delete IPv6 block`,
		Long:      `Delete IPv6 block subnet with the given ID.`,
		Namespace: "dedibox",
		Resource:  "ipv6-block",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.IPv6BlockAPIDeleteIPv6BlockRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "block-id",
				Short:      `ID of the IPv6 block to delete`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.IPv6BlockAPIDeleteIPv6BlockRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewIPv6BlockAPI(client)
			e = api.DeleteIPv6Block(request)
			if e != nil {
				return nil, e
			}

			return &core.SuccessResult{
				Resource: "ipv6-block",
				Verb:     "delete",
			}, nil
		},
	}
}

func dediboxIPv6BlockCreateSubnet() *core.Command {
	return &core.Command{
		Short: `Create IPv6 block subnet`,
		Long: `Create IPv6 block subnet for the given IP ID.
/48 could create subnet in /56 (quota link to your number of server).
/56 could create subnet in /64 (quota link to your number of failover IP).`,
		Namespace: "dedibox",
		Resource:  "ipv6-block",
		Verb:      "create-subnet",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.IPv6BlockAPICreateIPv6BlockSubnetRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "block-id",
				Short:      `ID of the IPv6 block`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "address",
				Short:      `Address of the IPv6`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "cidr",
				Short:      `Classless InterDomain Routing notation of the IPv6`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.IPv6BlockAPICreateIPv6BlockSubnetRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewIPv6BlockAPI(client)

			return api.CreateIPv6BlockSubnet(request)
		},
	}
}

func dediboxIPv6BlockListSubnet() *core.Command {
	return &core.Command{
		Short:     `List available IPv6 block subnets`,
		Long:      `List all available IPv6 block subnets for given IP ID.`,
		Namespace: "dedibox",
		Resource:  "ipv6-block",
		Verb:      "list-subnet",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.IPv6BlockAPIListIPv6BlockSubnetsAvailableRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "block-id",
				Short:      `ID of the IPv6 block`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.IPv6BlockAPIListIPv6BlockSubnetsAvailableRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewIPv6BlockAPI(client)

			return api.ListIPv6BlockSubnetsAvailable(request)
		},
	}
}

func dediboxRpnInfoList() *core.Command {
	return &core.Command{
		Short:     `List dedibox resources`,
		Long:      `List dedibox resources.`,
		Namespace: "dedibox",
		Resource:  "rpn-info",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.RpnAPIListRpnServerCapabilitiesRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `Order of the servers`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"created_at_asc",
					"created_at_desc",
				},
			},
			{
				Name:       "project-id",
				Short:      `Filter servers by project ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.RpnAPIListRpnServerCapabilitiesRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewRpnAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			resp, err := api.ListRpnServerCapabilities(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.Servers, nil
		},
	}
}

func dediboxRpnInfoGet() *core.Command {
	return &core.Command{
		Short:     `Get dedibox resources`,
		Long:      `Get dedibox resources.`,
		Namespace: "dedibox",
		Resource:  "rpn-info",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.RpnAPIGetRpnStatusRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "project-id",
				Short:      `A project ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "rpnv1-group-id",
				Short:      `An RPN v1 group ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "rpnv2-group-id",
				Short:      `An RPN v2 group ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.RpnAPIGetRpnStatusRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewRpnAPI(client)

			return api.GetRpnStatus(request)
		},
	}
}

func dediboxSanList() *core.Command {
	return &core.Command{
		Short:     `List dedibox resources`,
		Long:      `List dedibox resources.`,
		Namespace: "dedibox",
		Resource:  "san",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.RpnSanAPIListRpnSansRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `Order of the RPN SANs`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"created_at_asc",
					"created_at_desc",
				},
			},
			{
				Name:       "project-id",
				Short:      `Filter RPN SANs by project ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.RpnSanAPIListRpnSansRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewRpnSanAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			resp, err := api.ListRpnSans(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.RpnSans, nil
		},
	}
}

func dediboxSanGet() *core.Command {
	return &core.Command{
		Short:     `Get dedibox resources`,
		Long:      `Get dedibox resources.`,
		Namespace: "dedibox",
		Resource:  "san",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.RpnSanAPIGetRpnSanRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "rpn-san-id",
				Short:      `RPN SAN ID`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.RpnSanAPIGetRpnSanRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewRpnSanAPI(client)

			return api.GetRpnSan(request)
		},
	}
}

func dediboxSanDelete() *core.Command {
	return &core.Command{
		Short:     `Delete dedibox resources`,
		Long:      `Delete dedibox resources.`,
		Namespace: "dedibox",
		Resource:  "san",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.RpnSanAPIDeleteRpnSanRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "rpn-san-id",
				Short:      `RPN SAN ID`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.RpnSanAPIDeleteRpnSanRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewRpnSanAPI(client)
			e = api.DeleteRpnSan(request)
			if e != nil {
				return nil, e
			}

			return &core.SuccessResult{
				Resource: "san",
				Verb:     "delete",
			}, nil
		},
	}
}

func dediboxSanCreate() *core.Command {
	return &core.Command{
		Short:     `Create dedibox resources`,
		Long:      `Create dedibox resources.`,
		Namespace: "dedibox",
		Resource:  "san",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.RpnSanAPICreateRpnSanRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "offer-id",
				Short:      `Offer ID`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.ProjectIDArgSpec(),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.RpnSanAPICreateRpnSanRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewRpnSanAPI(client)

			return api.CreateRpnSan(request)
		},
	}
}

func dediboxSanListIPs() *core.Command {
	return &core.Command{
		Short:     `List-ips dedibox resources`,
		Long:      `List-ips dedibox resources.`,
		Namespace: "dedibox",
		Resource:  "san",
		Verb:      "list-ips",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.RpnSanAPIListIPsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "rpn-san-id",
				Short:      `RPN SAN ID`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "type",
				Short:      `Filter by IP type (server | rpnv2_subnet)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown",
					"server_ip",
					"rpnv2_subnet",
				},
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.RpnSanAPIListIPsRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewRpnSanAPI(client)

			return api.ListIPs(request)
		},
	}
}

func dediboxSanAddIP() *core.Command {
	return &core.Command{
		Short:     `Add-ip dedibox resources`,
		Long:      `Add-ip dedibox resources.`,
		Namespace: "dedibox",
		Resource:  "san",
		Verb:      "add-ip",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.RpnSanAPIAddIPRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "rpn-san-id",
				Short:      `RPN SAN ID`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "ip-ids.{index}",
				Short:      `An array of IP ID`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.RpnSanAPIAddIPRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewRpnSanAPI(client)
			e = api.AddIP(request)
			if e != nil {
				return nil, e
			}

			return &core.SuccessResult{
				Resource: "san",
				Verb:     "add-ip",
			}, nil
		},
	}
}

func dediboxSanRemoveIP() *core.Command {
	return &core.Command{
		Short:     `Remove-ip dedibox resources`,
		Long:      `Remove-ip dedibox resources.`,
		Namespace: "dedibox",
		Resource:  "san",
		Verb:      "remove-ip",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.RpnSanAPIRemoveIPRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "rpn-san-id",
				Short:      `RPN SAN ID`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "ip-ids.{index}",
				Short:      `An array of IP ID`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.RpnSanAPIRemoveIPRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewRpnSanAPI(client)
			e = api.RemoveIP(request)
			if e != nil {
				return nil, e
			}

			return &core.SuccessResult{
				Resource: "san",
				Verb:     "remove-ip",
			}, nil
		},
	}
}

func dediboxSanListAvailableIPs() *core.Command {
	return &core.Command{
		Short:     `List-available-ips dedibox resources`,
		Long:      `List-available-ips dedibox resources.`,
		Namespace: "dedibox",
		Resource:  "san",
		Verb:      "list-available-ips",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.RpnSanAPIListAvailableIPsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "rpn-san-id",
				Short:      `RPN SAN ID`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "type",
				Short:      `Filter by IP type (server | rpnv2_subnet)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown",
					"server_ip",
					"rpnv2_subnet",
				},
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.RpnSanAPIListAvailableIPsRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewRpnSanAPI(client)

			return api.ListAvailableIPs(request)
		},
	}
}

func dediboxRpnV1List() *core.Command {
	return &core.Command{
		Short:     `List dedibox resources`,
		Long:      `List dedibox resources.`,
		Namespace: "dedibox",
		Resource:  "rpn-v1",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.RpnV1ApiListRpnGroupsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `Order of the rpn v1 groups`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"created_at_asc",
					"created_at_desc",
				},
			},
			{
				Name:       "project-id",
				Short:      `Filter rpn v1 groups by project ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.RpnV1ApiListRpnGroupsRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewRpnV1API(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			resp, err := api.ListRpnGroups(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.RpnGroups, nil
		},
	}
}

func dediboxRpnV1Get() *core.Command {
	return &core.Command{
		Short:     `Get dedibox resources`,
		Long:      `Get dedibox resources.`,
		Namespace: "dedibox",
		Resource:  "rpn-v1",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.RpnV1ApiGetRpnGroupRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "group-id",
				Short:      `Rpn v1 group ID`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.RpnV1ApiGetRpnGroupRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewRpnV1API(client)

			return api.GetRpnGroup(request)
		},
	}
}

func dediboxRpnV1Create() *core.Command {
	return &core.Command{
		Short:     `Create dedibox resources`,
		Long:      `Create dedibox resources.`,
		Namespace: "dedibox",
		Resource:  "rpn-v1",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.RpnV1ApiCreateRpnGroupRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "name",
				Short:      `Rpn v1 group name`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "server-ids.{index}",
				Short:      `A collection of rpn v1 capable servers`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "san-server-ids.{index}",
				Short:      `A collection of rpn v1 capable rpn sans servers`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ProjectIDArgSpec(),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.RpnV1ApiCreateRpnGroupRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewRpnV1API(client)

			return api.CreateRpnGroup(request)
		},
	}
}

func dediboxRpnV1Delete() *core.Command {
	return &core.Command{
		Short:     `Delete dedibox resources`,
		Long:      `Delete dedibox resources.`,
		Namespace: "dedibox",
		Resource:  "rpn-v1",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.RpnV1ApiDeleteRpnGroupRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "group-id",
				Short:      `Rpn v1 group ID`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.RpnV1ApiDeleteRpnGroupRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewRpnV1API(client)
			e = api.DeleteRpnGroup(request)
			if e != nil {
				return nil, e
			}

			return &core.SuccessResult{
				Resource: "rpn-v1",
				Verb:     "delete",
			}, nil
		},
	}
}

func dediboxRpnV1Update() *core.Command {
	return &core.Command{
		Short:     `Update dedibox resources`,
		Long:      `Update dedibox resources.`,
		Namespace: "dedibox",
		Resource:  "rpn-v1",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.RpnV1ApiUpdateRpnGroupNameRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "group-id",
				Short:      `Rpn v1 group ID`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "name",
				Short:      `New rpn v1 group name`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.RpnV1ApiUpdateRpnGroupNameRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewRpnV1API(client)

			return api.UpdateRpnGroupName(request)
		},
	}
}

func dediboxRpnV1ListMembers() *core.Command {
	return &core.Command{
		Short:     `List-members dedibox resources`,
		Long:      `List-members dedibox resources.`,
		Namespace: "dedibox",
		Resource:  "rpn-v1",
		Verb:      "list-members",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.RpnV1ApiListRpnGroupMembersRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `Order of the rpn v1 group members`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"created_at_asc",
					"created_at_desc",
				},
			},
			{
				Name:       "group-id",
				Short:      `Filter rpn v1 group members by group ID`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "project-id",
				Short:      `A project ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.RpnV1ApiListRpnGroupMembersRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewRpnV1API(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			resp, err := api.ListRpnGroupMembers(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.Members, nil
		},
	}
}

func dediboxRpnV1Invite() *core.Command {
	return &core.Command{
		Short:     `Invite dedibox resources`,
		Long:      `Invite dedibox resources.`,
		Namespace: "dedibox",
		Resource:  "rpn-v1",
		Verb:      "invite",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.RpnV1ApiRpnGroupInviteRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "group-id",
				Short:      `The RPN V1 group ID`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "server-ids.{index}",
				Short:      `A collection of external server IDs`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.ProjectIDArgSpec(),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.RpnV1ApiRpnGroupInviteRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewRpnV1API(client)
			e = api.RpnGroupInvite(request)
			if e != nil {
				return nil, e
			}

			return &core.SuccessResult{
				Resource: "rpn-v1",
				Verb:     "invite",
			}, nil
		},
	}
}

func dediboxRpnV1Leave() *core.Command {
	return &core.Command{
		Short:     `Leave dedibox resources`,
		Long:      `Leave dedibox resources.`,
		Namespace: "dedibox",
		Resource:  "rpn-v1",
		Verb:      "leave",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.RpnV1ApiLeaveRpnGroupRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "group-id",
				Short:      `The RPN V1 group ID`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "member-ids.{index}",
				Short:      `A collection of rpn v1 group members IDs`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.RpnV1ApiLeaveRpnGroupRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewRpnV1API(client)
			e = api.LeaveRpnGroup(request)
			if e != nil {
				return nil, e
			}

			return &core.SuccessResult{
				Resource: "rpn-v1",
				Verb:     "leave",
			}, nil
		},
	}
}

func dediboxRpnV1AddMembers() *core.Command {
	return &core.Command{
		Short:     `Add-members dedibox resources`,
		Long:      `Add-members dedibox resources.`,
		Namespace: "dedibox",
		Resource:  "rpn-v1",
		Verb:      "add-members",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.RpnV1ApiAddRpnGroupMembersRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "group-id",
				Short:      `The rpn v1 group ID`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "server-ids.{index}",
				Short:      `A collection of rpn v1 capable server IDs`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "san-server-ids.{index}",
				Short:      `A collection of rpn v1 capable RPN SAN server IDs`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.RpnV1ApiAddRpnGroupMembersRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewRpnV1API(client)

			return api.AddRpnGroupMembers(request)
		},
	}
}

func dediboxRpnV1DeleteMembers() *core.Command {
	return &core.Command{
		Short:     `Delete-members dedibox resources`,
		Long:      `Delete-members dedibox resources.`,
		Namespace: "dedibox",
		Resource:  "rpn-v1",
		Verb:      "delete-members",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.RpnV1ApiDeleteRpnGroupMembersRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "group-id",
				Short:      `The rpn v1 group ID`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "member-ids.{index}",
				Short:      `A collection of rpn v1 group members IDs`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.RpnV1ApiDeleteRpnGroupMembersRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewRpnV1API(client)

			return api.DeleteRpnGroupMembers(request)
		},
	}
}

func dediboxRpnV1ListCapableServer() *core.Command {
	return &core.Command{
		Short:     `List-capable-server dedibox resources`,
		Long:      `List-capable-server dedibox resources.`,
		Namespace: "dedibox",
		Resource:  "rpn-v1",
		Verb:      "list-capable-server",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.RpnV1ApiListRpnCapableServersRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `Order of the rpn capable resources`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"created_at_asc",
					"created_at_desc",
				},
			},
			{
				Name:       "project-id",
				Short:      `Filter rpn capable resources by project ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.RpnV1ApiListRpnCapableServersRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewRpnV1API(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			resp, err := api.ListRpnCapableServers(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.Servers, nil
		},
	}
}

func dediboxRpnV1ListCapableSanServer() *core.Command {
	return &core.Command{
		Short:     `List-capable-san-server dedibox resources`,
		Long:      `List-capable-san-server dedibox resources.`,
		Namespace: "dedibox",
		Resource:  "rpn-v1",
		Verb:      "list-capable-san-server",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.RpnV1ApiListRpnCapableSanServersRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `Order of the rpn capable resources`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"created_at_asc",
					"created_at_desc",
				},
			},
			{
				Name:       "project-id",
				Short:      `Filter rpn capable resources by project ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.RpnV1ApiListRpnCapableSanServersRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewRpnV1API(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			resp, err := api.ListRpnCapableSanServers(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.SanServers, nil
		},
	}
}

func dediboxRpnV1ListInvites() *core.Command {
	return &core.Command{
		Short:     `List-invites dedibox resources`,
		Long:      `List-invites dedibox resources.`,
		Namespace: "dedibox",
		Resource:  "rpn-v1",
		Verb:      "list-invites",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.RpnV1ApiListRpnInvitesRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `Order of the rpn capable resources`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"created_at_asc",
					"created_at_desc",
				},
			},
			core.ProjectIDArgSpec(),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.RpnV1ApiListRpnInvitesRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewRpnV1API(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			resp, err := api.ListRpnInvites(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.Members, nil
		},
	}
}

func dediboxRpnV1AcceptInvite() *core.Command {
	return &core.Command{
		Short:     `Accept-invite dedibox resources`,
		Long:      `Accept-invite dedibox resources.`,
		Namespace: "dedibox",
		Resource:  "rpn-v1",
		Verb:      "accept-invite",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.RpnV1ApiAcceptRpnInviteRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "member-id",
				Short:      `The member ID`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.RpnV1ApiAcceptRpnInviteRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewRpnV1API(client)
			e = api.AcceptRpnInvite(request)
			if e != nil {
				return nil, e
			}

			return &core.SuccessResult{
				Resource: "rpn-v1",
				Verb:     "accept-invite",
			}, nil
		},
	}
}

func dediboxRpnV1RefuseInvite() *core.Command {
	return &core.Command{
		Short:     `Refuse-invite dedibox resources`,
		Long:      `Refuse-invite dedibox resources.`,
		Namespace: "dedibox",
		Resource:  "rpn-v1",
		Verb:      "refuse-invite",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.RpnV1ApiRefuseRpnInviteRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "member-id",
				Short:      `The member ID`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.RpnV1ApiRefuseRpnInviteRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewRpnV1API(client)
			e = api.RefuseRpnInvite(request)
			if e != nil {
				return nil, e
			}

			return &core.SuccessResult{
				Resource: "rpn-v1",
				Verb:     "refuse-invite",
			}, nil
		},
	}
}

func dediboxRpnV2List() *core.Command {
	return &core.Command{
		Short:     `List dedibox resources`,
		Long:      `List dedibox resources.`,
		Namespace: "dedibox",
		Resource:  "rpn-v2",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.RpnV2ApiListRpnV2GroupsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `Order of the rpn v2 groups`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"created_at_asc",
					"created_at_desc",
				},
			},
			{
				Name:       "project-id",
				Short:      `Filter rpn v2 groups by project ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.RpnV2ApiListRpnV2GroupsRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewRpnV2API(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			resp, err := api.ListRpnV2Groups(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.RpnGroups, nil
		},
	}
}

func dediboxRpnV2ListMembers() *core.Command {
	return &core.Command{
		Short:     `List-members dedibox resources`,
		Long:      `List-members dedibox resources.`,
		Namespace: "dedibox",
		Resource:  "rpn-v2",
		Verb:      "list-members",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.RpnV2ApiListRpnV2MembersRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `Order of the rpn v2 group members`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"created_at_asc",
					"created_at_desc",
				},
			},
			{
				Name:       "group-id",
				Short:      `RPN V2 group ID`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "type",
				Short:      `Filter members by type`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_type",
					"rpnv1_group",
					"server",
				},
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.RpnV2ApiListRpnV2MembersRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewRpnV2API(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			resp, err := api.ListRpnV2Members(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.Members, nil
		},
	}
}

func dediboxRpnV2Get() *core.Command {
	return &core.Command{
		Short:     `Get dedibox resources`,
		Long:      `Get dedibox resources.`,
		Namespace: "dedibox",
		Resource:  "rpn-v2",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.RpnV2ApiGetRpnV2GroupRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "group-id",
				Short:      `RPN V2 group ID`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.RpnV2ApiGetRpnV2GroupRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewRpnV2API(client)

			return api.GetRpnV2Group(request)
		},
	}
}

func dediboxRpnV2Create() *core.Command {
	return &core.Command{
		Short:     `Create dedibox resources`,
		Long:      `Create dedibox resources.`,
		Namespace: "dedibox",
		Resource:  "rpn-v2",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.RpnV2ApiCreateRpnV2GroupRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ProjectIDArgSpec(),
			{
				Name:       "type",
				Short:      `RPN V2 group type (qing / standard)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_type",
					"standard",
					"qinq",
				},
			},
			{
				Name:       "name",
				Short:      `RPN V2 group name`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "servers.{index}",
				Short:      `A collection of server IDs`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.RpnV2ApiCreateRpnV2GroupRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewRpnV2API(client)

			return api.CreateRpnV2Group(request)
		},
	}
}

func dediboxRpnV2Delete() *core.Command {
	return &core.Command{
		Short:     `Delete dedibox resources`,
		Long:      `Delete dedibox resources.`,
		Namespace: "dedibox",
		Resource:  "rpn-v2",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.RpnV2ApiDeleteRpnV2GroupRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "group-id",
				Short:      `RPN V2 group ID`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.RpnV2ApiDeleteRpnV2GroupRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewRpnV2API(client)
			e = api.DeleteRpnV2Group(request)
			if e != nil {
				return nil, e
			}

			return &core.SuccessResult{
				Resource: "rpn-v2",
				Verb:     "delete",
			}, nil
		},
	}
}

func dediboxRpnV2Update() *core.Command {
	return &core.Command{
		Short:     `Update dedibox resources`,
		Long:      `Update dedibox resources.`,
		Namespace: "dedibox",
		Resource:  "rpn-v2",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.RpnV2ApiUpdateRpnV2GroupNameRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "group-id",
				Short:      `RPN V2 group ID`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "name",
				Short:      `RPN V2 group name`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.RpnV2ApiUpdateRpnV2GroupNameRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewRpnV2API(client)

			return api.UpdateRpnV2GroupName(request)
		},
	}
}

func dediboxRpnV2AddMembers() *core.Command {
	return &core.Command{
		Short:     `Add-members dedibox resources`,
		Long:      `Add-members dedibox resources.`,
		Namespace: "dedibox",
		Resource:  "rpn-v2",
		Verb:      "add-members",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.RpnV2ApiAddRpnV2MembersRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "group-id",
				Short:      `RPN V2 group ID`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "servers.{index}",
				Short:      `A collection of server IDs`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.RpnV2ApiAddRpnV2MembersRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewRpnV2API(client)
			e = api.AddRpnV2Members(request)
			if e != nil {
				return nil, e
			}

			return &core.SuccessResult{
				Resource: "rpn-v2",
				Verb:     "add-members",
			}, nil
		},
	}
}

func dediboxRpnV2DeleteMembers() *core.Command {
	return &core.Command{
		Short:     `Delete-members dedibox resources`,
		Long:      `Delete-members dedibox resources.`,
		Namespace: "dedibox",
		Resource:  "rpn-v2",
		Verb:      "delete-members",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.RpnV2ApiDeleteRpnV2MembersRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "group-id",
				Short:      `RPN V2 group ID`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "member-ids.{index}",
				Short:      `A collection of member IDs`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.RpnV2ApiDeleteRpnV2MembersRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewRpnV2API(client)
			e = api.DeleteRpnV2Members(request)
			if e != nil {
				return nil, e
			}

			return &core.SuccessResult{
				Resource: "rpn-v2",
				Verb:     "delete-members",
			}, nil
		},
	}
}

func dediboxRpnV2ListCapableResources() *core.Command {
	return &core.Command{
		Short:     `List-capable-resources dedibox resources`,
		Long:      `List-capable-resources dedibox resources.`,
		Namespace: "dedibox",
		Resource:  "rpn-v2",
		Verb:      "list-capable-resources",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.RpnV2ApiListRpnV2CapableResourcesRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `Order of the rpn v2 capable resources`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"created_at_asc",
					"created_at_desc",
				},
			},
			{
				Name:       "project-id",
				Short:      `Filter rpn v2 capable resources by project ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.RpnV2ApiListRpnV2CapableResourcesRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewRpnV2API(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			resp, err := api.ListRpnV2CapableResources(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.Servers, nil
		},
	}
}

func dediboxRpnV2ListLogs() *core.Command {
	return &core.Command{
		Short:     `List-logs dedibox resources`,
		Long:      `List-logs dedibox resources.`,
		Namespace: "dedibox",
		Resource:  "rpn-v2",
		Verb:      "list-logs",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.RpnV2ApiListRpnV2GroupLogsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `Order of the rpn v2 group logs`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"created_at_asc",
					"created_at_desc",
				},
			},
			{
				Name:       "group-id",
				Short:      `RPN V2 group ID`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.RpnV2ApiListRpnV2GroupLogsRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewRpnV2API(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			resp, err := api.ListRpnV2GroupLogs(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.Logs, nil
		},
	}
}

func dediboxRpnV2UpdateVlanMembers() *core.Command {
	return &core.Command{
		Short:     `Update-vlan-members dedibox resources`,
		Long:      `Update-vlan-members dedibox resources.`,
		Namespace: "dedibox",
		Resource:  "rpn-v2",
		Verb:      "update-vlan-members",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.RpnV2ApiUpdateRpnV2VlanForMembersRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "group-id",
				Short:      `RPN V2 group ID`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "member-ids.{index}",
				Short:      `RPN V2 member IDs`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "vlan",
				Short:      `A vlan`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.RpnV2ApiUpdateRpnV2VlanForMembersRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewRpnV2API(client)
			e = api.UpdateRpnV2VlanForMembers(request)
			if e != nil {
				return nil, e
			}

			return &core.SuccessResult{
				Resource: "rpn-v2",
				Verb:     "update-vlan-members",
			}, nil
		},
	}
}

func dediboxRpnV2EnableCompatibility() *core.Command {
	return &core.Command{
		Short:     `Enable-compatibility dedibox resources`,
		Long:      `Enable-compatibility dedibox resources.`,
		Namespace: "dedibox",
		Resource:  "rpn-v2",
		Verb:      "enable-compatibility",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.RpnV2ApiEnableRpnV2GroupCompatibilityRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "group-id",
				Short:      `RPN V2 group ID`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "rpnv1-group-id",
				Short:      `RPN V1 group ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.RpnV2ApiEnableRpnV2GroupCompatibilityRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewRpnV2API(client)
			e = api.EnableRpnV2GroupCompatibility(request)
			if e != nil {
				return nil, e
			}

			return &core.SuccessResult{
				Resource: "rpn-v2",
				Verb:     "enable-compatibility",
			}, nil
		},
	}
}

func dediboxRpnV2DisableCompatibility() *core.Command {
	return &core.Command{
		Short:     `Disable-compatibility dedibox resources`,
		Long:      `Disable-compatibility dedibox resources.`,
		Namespace: "dedibox",
		Resource:  "rpn-v2",
		Verb:      "disable-compatibility",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(dedibox.RpnV2ApiDisableRpnV2GroupCompatibilityRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "group-id",
				Short:      `RPN V2 group ID`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dedibox.RpnV2ApiDisableRpnV2GroupCompatibilityRequest)

			client := core.ExtractClient(ctx)
			api := dedibox.NewRpnV2API(client)
			e = api.DisableRpnV2GroupCompatibility(request)
			if e != nil {
				return nil, e
			}

			return &core.SuccessResult{
				Resource: "rpn-v2",
				Verb:     "disable-compatibility",
			}, nil
		},
	}
}
