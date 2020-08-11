// This file was automatically generated. DO NOT EDIT.
// If you have any remark or suggestion do not hesitate to open an issue.

package iot

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-sdk-go/api/iot/v1beta1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// always import dependencies
var (
	_ = scw.RegionFrPar
)

func GetGeneratedCommands() *core.Commands {
	return core.NewCommands(
		iotRoot(),
		iotHub(),
		iotDevice(),
		iotRoute(),
		iotNetwork(),
		iotHubList(),
		iotHubCreate(),
		iotHubGet(),
		iotHubUpdate(),
		iotHubEnable(),
		iotHubDisable(),
		iotHubDelete(),
		iotHubGetMetrics(),
		iotDeviceList(),
		iotDeviceCreate(),
		iotDeviceGet(),
		iotDeviceUpdate(),
		iotDeviceEnable(),
		iotDeviceDisable(),
		iotDeviceDelete(),
		iotDeviceGetMetrics(),
		iotNetworkList(),
		iotNetworkCreate(),
		iotNetworkGet(),
		iotNetworkDelete(),
	)
}
func iotRoot() *core.Command {
	return &core.Command{
		Short:     `This API allows you to manage IoT hubs and devices`,
		Long:      ``,
		Namespace: "iot",
	}
}

func iotHub() *core.Command {
	return &core.Command{
		Short:     `IoT Hub commands`,
		Long:      `IoT Hub commands.`,
		Namespace: "iot",
		Resource:  "hub",
	}
}

func iotDevice() *core.Command {
	return &core.Command{
		Short:     `IoT Device commands`,
		Long:      `IoT Device commands.`,
		Namespace: "iot",
		Resource:  "device",
	}
}

func iotRoute() *core.Command {
	return &core.Command{
		Short:     `IoT Route commands`,
		Long:      `IoT Route commands.`,
		Namespace: "iot",
		Resource:  "route",
	}
}

func iotNetwork() *core.Command {
	return &core.Command{
		Short:     `IoT Network commands`,
		Long:      `IoT Network commands.`,
		Namespace: "iot",
		Resource:  "network",
	}
}

func iotHubList() *core.Command {
	return &core.Command{
		Short:     `List hubs`,
		Long:      `List hubs.`,
		Namespace: "iot",
		Resource:  "hub",
		Verb:      "list",
		ArgsType:  reflect.TypeOf(iot.ListHubsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `Ordering of requested hub`,
				Required:   false,
				Positional: false,
				EnumValues: []string{"name_asc", "name_desc", "status_asc", "status_desc", "product_plan_asc", "product_plan_desc", "created_at_asc", "created_at_desc", "updated_at_asc", "updated_at_desc", "enabled_asc", "enabled_desc"},
			},
			{
				Name:       "name",
				Short:      `Filter on the name`,
				Required:   false,
				Positional: false,
			},
			{
				Name:       "organization-id",
				Short:      `Filter on the organization`,
				Required:   false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*iot.ListHubsRequest)

			client := core.ExtractClient(ctx)
			api := iot.NewAPI(client)
			resp, err := api.ListHubs(request, scw.WithAllPages())
			if err != nil {
				return nil, err
			}
			return resp.Hubs, nil

		},
	}
}

func iotHubCreate() *core.Command {
	return &core.Command{
		Short:     `Create a hub`,
		Long:      `Create a hub.`,
		Namespace: "iot",
		Resource:  "hub",
		Verb:      "create",
		ArgsType:  reflect.TypeOf(iot.CreateHubRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "name",
				Short:      `Hub name (up to 255 characters)`,
				Required:   true,
				Positional: false,
			},
			{
				Name:       "product-plan",
				Short:      `Hub feature set`,
				Required:   true,
				Positional: false,
				EnumValues: []string{"plan_unknown", "plan_shared", "plan_dedicated", "plan_ha"},
			},
			{
				Name:       "disable-events",
				Short:      `Disable Hub events (default false)`,
				Required:   false,
				Positional: false,
			},
			{
				Name:       "events-topic-prefix",
				Short:      `Hub events topic prefix (default '$SCW/events')`,
				Required:   false,
				Positional: false,
			},
			core.OrganizationIDArgSpec(),
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*iot.CreateHubRequest)

			client := core.ExtractClient(ctx)
			api := iot.NewAPI(client)
			return api.CreateHub(request)

		},
	}
}

func iotHubGet() *core.Command {
	return &core.Command{
		Short:     `Get a hub`,
		Long:      `Get a hub.`,
		Namespace: "iot",
		Resource:  "hub",
		Verb:      "get",
		ArgsType:  reflect.TypeOf(iot.GetHubRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "hub-id",
				Short:      `Hub ID`,
				Required:   true,
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*iot.GetHubRequest)

			client := core.ExtractClient(ctx)
			api := iot.NewAPI(client)
			return api.GetHub(request)

		},
	}
}

func iotHubUpdate() *core.Command {
	return &core.Command{
		Short:     `Update a hub`,
		Long:      `Update a hub.`,
		Namespace: "iot",
		Resource:  "hub",
		Verb:      "update",
		ArgsType:  reflect.TypeOf(iot.UpdateHubRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "hub-id",
				Short:      `Hub ID`,
				Required:   true,
				Positional: true,
			},
			{
				Name:       "name",
				Short:      `Hub name (up to 255 characters)`,
				Required:   false,
				Positional: false,
			},
			{
				Name:       "disable-events",
				Short:      `Disable events`,
				Required:   false,
				Positional: false,
			},
			{
				Name:       "events-topic-prefix",
				Short:      `Hub events topic prefix`,
				Required:   false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*iot.UpdateHubRequest)

			client := core.ExtractClient(ctx)
			api := iot.NewAPI(client)
			return api.UpdateHub(request)

		},
	}
}

func iotHubEnable() *core.Command {
	return &core.Command{
		Short:     `Enable a hub`,
		Long:      `Enable a hub.`,
		Namespace: "iot",
		Resource:  "hub",
		Verb:      "enable",
		ArgsType:  reflect.TypeOf(iot.EnableHubRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "hub-id",
				Short:      `Hub ID`,
				Required:   true,
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*iot.EnableHubRequest)

			client := core.ExtractClient(ctx)
			api := iot.NewAPI(client)
			return api.EnableHub(request)

		},
	}
}

func iotHubDisable() *core.Command {
	return &core.Command{
		Short:     `Disable a hub`,
		Long:      `Disable a hub.`,
		Namespace: "iot",
		Resource:  "hub",
		Verb:      "disable",
		ArgsType:  reflect.TypeOf(iot.DisableHubRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "hub-id",
				Short:      `Hub ID`,
				Required:   true,
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*iot.DisableHubRequest)

			client := core.ExtractClient(ctx)
			api := iot.NewAPI(client)
			return api.DisableHub(request)

		},
	}
}

func iotHubDelete() *core.Command {
	return &core.Command{
		Short:     `Delete a hub`,
		Long:      `Delete a hub.`,
		Namespace: "iot",
		Resource:  "hub",
		Verb:      "delete",
		ArgsType:  reflect.TypeOf(iot.DeleteHubRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "hub-id",
				Short:      `Hub ID`,
				Required:   true,
				Positional: true,
			},
			{
				Name:       "delete-devices",
				Short:      `Force deletion of devices added to this hub instead of rejecting operation`,
				Required:   false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*iot.DeleteHubRequest)

			client := core.ExtractClient(ctx)
			api := iot.NewAPI(client)
			e = api.DeleteHub(request)
			if e != nil {
				return nil, e
			}
			return &core.SuccessResult{
				Resource: "hub",
				Verb:     "delete",
			}, nil
		},
	}
}

func iotHubGetMetrics() *core.Command {
	return &core.Command{
		Short:     `Get a hub's metrics`,
		Long:      `Get a hub's metrics.`,
		Namespace: "iot",
		Resource:  "hub",
		Verb:      "get-metrics",
		ArgsType:  reflect.TypeOf(iot.GetHubMetricsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "hub-id",
				Short:      `Hub ID`,
				Required:   true,
				Positional: true,
			},
			{
				Name:       "period",
				Short:      `Period over which the metrics span`,
				Required:   true,
				Positional: false,
				EnumValues: []string{"hour", "day", "week", "month", "year"},
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*iot.GetHubMetricsRequest)

			client := core.ExtractClient(ctx)
			api := iot.NewAPI(client)
			return api.GetHubMetrics(request)

		},
	}
}

func iotDeviceList() *core.Command {
	return &core.Command{
		Short:     `List devices`,
		Long:      `List devices.`,
		Namespace: "iot",
		Resource:  "device",
		Verb:      "list",
		ArgsType:  reflect.TypeOf(iot.ListDevicesRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `Ordering of requested devices`,
				Required:   false,
				Positional: false,
				EnumValues: []string{"name_asc", "name_desc", "status_asc", "status_desc", "hub_id_asc", "hub_id_desc", "created_at_asc", "created_at_desc", "updated_at_asc", "updated_at_desc", "enabled_asc", "enabled_desc", "allow_insecure_asc", "allow_insecure_desc", "last_seen_at_asc", "last_seen_at_desc"},
			},
			{
				Name:       "name",
				Short:      `Filter on the name`,
				Required:   false,
				Positional: false,
			},
			{
				Name:       "hub-id",
				Short:      `Filter on the hub`,
				Required:   false,
				Positional: false,
			},
			{
				Name:       "enabled",
				Short:      `Filter on the enabled flag`,
				Required:   false,
				Positional: false,
			},
			{
				Name:       "allow-insecure",
				Short:      `Filter on the allow_insecure flag`,
				Required:   false,
				Positional: false,
			},
			{
				Name:       "is-connected",
				Short:      `Filter on the is_connected state`,
				Required:   false,
				Positional: false,
			},
			{
				Name:       "organization-id",
				Short:      `Filter on the organization`,
				Required:   false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*iot.ListDevicesRequest)

			client := core.ExtractClient(ctx)
			api := iot.NewAPI(client)
			resp, err := api.ListDevices(request, scw.WithAllPages())
			if err != nil {
				return nil, err
			}
			return resp.Devices, nil

		},
	}
}

func iotDeviceCreate() *core.Command {
	return &core.Command{
		Short:     `Add a device`,
		Long:      `Add a device.`,
		Namespace: "iot",
		Resource:  "device",
		Verb:      "create",
		ArgsType:  reflect.TypeOf(iot.CreateDeviceRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "name",
				Short:      `Device name`,
				Required:   true,
				Positional: false,
			},
			{
				Name:       "hub-id",
				Short:      `ID of the device's hub`,
				Required:   true,
				Positional: false,
			},
			{
				Name:       "allow-insecure",
				Short:      `Allow plain and server-authenticated SSL connections in addition to mutually-authenticated ones`,
				Required:   false,
				Positional: false,
			},
			{
				Name:       "message-filters.publish.policy",
				Required:   false,
				Positional: false,
				EnumValues: []string{"unknown", "accept", "reject"},
			},
			{
				Name:       "message-filters.publish.topics.{index}",
				Required:   false,
				Positional: false,
			},
			{
				Name:       "message-filters.subscribe.policy",
				Required:   false,
				Positional: false,
				EnumValues: []string{"unknown", "accept", "reject"},
			},
			{
				Name:       "message-filters.subscribe.topics.{index}",
				Required:   false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*iot.CreateDeviceRequest)

			client := core.ExtractClient(ctx)
			api := iot.NewAPI(client)
			return api.CreateDevice(request)

		},
	}
}

func iotDeviceGet() *core.Command {
	return &core.Command{
		Short:     `Get a device`,
		Long:      `Get a device.`,
		Namespace: "iot",
		Resource:  "device",
		Verb:      "get",
		ArgsType:  reflect.TypeOf(iot.GetDeviceRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "device-id",
				Short:      `Device ID`,
				Required:   true,
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*iot.GetDeviceRequest)

			client := core.ExtractClient(ctx)
			api := iot.NewAPI(client)
			return api.GetDevice(request)

		},
	}
}

func iotDeviceUpdate() *core.Command {
	return &core.Command{
		Short:     `Update a device`,
		Long:      `Update a device.`,
		Namespace: "iot",
		Resource:  "device",
		Verb:      "update",
		ArgsType:  reflect.TypeOf(iot.UpdateDeviceRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "device-id",
				Short:      `Device ID`,
				Required:   true,
				Positional: true,
			},
			{
				Name:       "name",
				Short:      `Device name`,
				Required:   false,
				Positional: false,
			},
			{
				Name:       "allow-insecure",
				Short:      `Allow plain and server-authenticated SSL connections in addition to mutually-authenticated ones`,
				Required:   false,
				Positional: false,
			},
			{
				Name:       "message-filters.publish.policy",
				Required:   false,
				Positional: false,
				EnumValues: []string{"unknown", "accept", "reject"},
			},
			{
				Name:       "message-filters.publish.topics.{index}",
				Required:   false,
				Positional: false,
			},
			{
				Name:       "message-filters.subscribe.policy",
				Required:   false,
				Positional: false,
				EnumValues: []string{"unknown", "accept", "reject"},
			},
			{
				Name:       "message-filters.subscribe.topics.{index}",
				Required:   false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*iot.UpdateDeviceRequest)

			client := core.ExtractClient(ctx)
			api := iot.NewAPI(client)
			return api.UpdateDevice(request)

		},
	}
}

func iotDeviceEnable() *core.Command {
	return &core.Command{
		Short:     `Enable a device`,
		Long:      `Enable a device.`,
		Namespace: "iot",
		Resource:  "device",
		Verb:      "enable",
		ArgsType:  reflect.TypeOf(iot.EnableDeviceRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "device-id",
				Short:      `Device ID`,
				Required:   true,
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*iot.EnableDeviceRequest)

			client := core.ExtractClient(ctx)
			api := iot.NewAPI(client)
			return api.EnableDevice(request)

		},
	}
}

func iotDeviceDisable() *core.Command {
	return &core.Command{
		Short:     `Disable a device`,
		Long:      `Disable a device.`,
		Namespace: "iot",
		Resource:  "device",
		Verb:      "disable",
		ArgsType:  reflect.TypeOf(iot.DisableDeviceRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "device-id",
				Short:      `Device ID`,
				Required:   true,
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*iot.DisableDeviceRequest)

			client := core.ExtractClient(ctx)
			api := iot.NewAPI(client)
			return api.DisableDevice(request)

		},
	}
}

func iotDeviceDelete() *core.Command {
	return &core.Command{
		Short:     `Remove a device`,
		Long:      `Remove a device.`,
		Namespace: "iot",
		Resource:  "device",
		Verb:      "delete",
		ArgsType:  reflect.TypeOf(iot.DeleteDeviceRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "device-id",
				Short:      `Device ID`,
				Required:   true,
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*iot.DeleteDeviceRequest)

			client := core.ExtractClient(ctx)
			api := iot.NewAPI(client)
			e = api.DeleteDevice(request)
			if e != nil {
				return nil, e
			}
			return &core.SuccessResult{
				Resource: "device",
				Verb:     "delete",
			}, nil
		},
	}
}

func iotDeviceGetMetrics() *core.Command {
	return &core.Command{
		Short:     `Get a device's metrics`,
		Long:      `Get a device's metrics.`,
		Namespace: "iot",
		Resource:  "device",
		Verb:      "get-metrics",
		ArgsType:  reflect.TypeOf(iot.GetDeviceMetricsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "device-id",
				Short:      `Device ID`,
				Required:   true,
				Positional: true,
			},
			{
				Name:       "period",
				Short:      `Period over which the metrics span`,
				Required:   true,
				Positional: false,
				EnumValues: []string{"hour", "day", "week", "month", "year"},
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*iot.GetDeviceMetricsRequest)

			client := core.ExtractClient(ctx)
			api := iot.NewAPI(client)
			return api.GetDeviceMetrics(request)

		},
	}
}

func iotNetworkList() *core.Command {
	return &core.Command{
		Short:     `List the Networks`,
		Long:      `List the Networks.`,
		Namespace: "iot",
		Resource:  "network",
		Verb:      "list",
		ArgsType:  reflect.TypeOf(iot.ListNetworksRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `Ordering of requested routes`,
				Required:   false,
				Positional: false,
				EnumValues: []string{"name_asc", "name_desc", "type_asc", "type_desc", "created_at_asc", "created_at_desc"},
			},
			{
				Name:       "name",
				Short:      `Filter on Network name`,
				Required:   false,
				Positional: false,
			},
			{
				Name:       "hub-id",
				Short:      `Filter on the hub`,
				Required:   false,
				Positional: false,
			},
			{
				Name:       "topic-prefix",
				Short:      `Filter on the topic prefix`,
				Required:   false,
				Positional: false,
			},
			{
				Name:       "organization-id",
				Short:      `Filter on the organization`,
				Required:   false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*iot.ListNetworksRequest)

			client := core.ExtractClient(ctx)
			api := iot.NewAPI(client)
			resp, err := api.ListNetworks(request, scw.WithAllPages())
			if err != nil {
				return nil, err
			}
			return resp.Networks, nil

		},
	}
}

func iotNetworkCreate() *core.Command {
	return &core.Command{
		Short:     `Create a new Network`,
		Long:      `Create a new Network.`,
		Namespace: "iot",
		Resource:  "network",
		Verb:      "create",
		ArgsType:  reflect.TypeOf(iot.CreateNetworkRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "name",
				Short:      `Network name`,
				Required:   true,
				Positional: false,
			},
			{
				Name:       "type",
				Short:      `Type of network to connect with`,
				Required:   true,
				Positional: false,
				EnumValues: []string{"unknown", "sigfox", "rest"},
			},
			{
				Name:       "hub-id",
				Short:      `Hub ID to connect the Network to`,
				Required:   true,
				Positional: false,
			},
			{
				Name:       "topic-prefix",
				Short:      `Topic prefix for the Network`,
				Required:   true,
				Positional: false,
			},
			core.OrganizationIDArgSpec(),
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*iot.CreateNetworkRequest)

			client := core.ExtractClient(ctx)
			api := iot.NewAPI(client)
			return api.CreateNetwork(request)

		},
	}
}

func iotNetworkGet() *core.Command {
	return &core.Command{
		Short:     `Retrieve a specific Network`,
		Long:      `Retrieve a specific Network.`,
		Namespace: "iot",
		Resource:  "network",
		Verb:      "get",
		ArgsType:  reflect.TypeOf(iot.GetNetworkRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "network-id",
				Short:      `Network ID`,
				Required:   true,
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*iot.GetNetworkRequest)

			client := core.ExtractClient(ctx)
			api := iot.NewAPI(client)
			return api.GetNetwork(request)

		},
	}
}

func iotNetworkDelete() *core.Command {
	return &core.Command{
		Short:     `Delete a Network`,
		Long:      `Delete a Network.`,
		Namespace: "iot",
		Resource:  "network",
		Verb:      "delete",
		ArgsType:  reflect.TypeOf(iot.DeleteNetworkRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "network-id",
				Short:      `Network ID`,
				Required:   true,
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*iot.DeleteNetworkRequest)

			client := core.ExtractClient(ctx)
			api := iot.NewAPI(client)
			e = api.DeleteNetwork(request)
			if e != nil {
				return nil, e
			}
			return &core.SuccessResult{
				Resource: "network",
				Verb:     "delete",
			}, nil
		},
	}
}
