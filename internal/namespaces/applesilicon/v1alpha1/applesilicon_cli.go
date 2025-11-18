// This file was automatically generated. DO NOT EDIT.
// If you have any remark or suggestion do not hesitate to open an issue.

package applesilicon

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	applesilicon "github.com/scaleway/scaleway-sdk-go/api/applesilicon/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// always import dependencies
var (
	_ = scw.RegionFrPar
)

func GetGeneratedCommands() *core.Commands {
	return core.NewCommands(
		appleSiliconRoot(),
		appleSiliconServer(),
		appleSiliconOs(),
		appleSiliconServerType(),
		appleSiliconPrivateNetwork(),
		appleSiliconServerTypeList(),
		appleSiliconServerTypeGet(),
		appleSiliconServerCreate(),
		appleSiliconServerList(),
		appleSiliconOsList(),
		appleSiliconOsGet(),
		appleSiliconServerGet(),
		appleSiliconServerUpdate(),
		appleSiliconServerDelete(),
		appleSiliconServerReboot(),
		appleSiliconServerReinstall(),
		appleSiliconPrivateNetworkAdd(),
		appleSiliconPrivateNetworkSet(),
		appleSiliconPrivateNetworkList(),
		appleSiliconPrivateNetworkDelete(),
	)
}

func appleSiliconRoot() *core.Command {
	return &core.Command{
		Short:     `Apple silicon API`,
		Long:      `Apple silicon API.`,
		Namespace: "apple-silicon",
	}
}

func appleSiliconServer() *core.Command {
	return &core.Command{
		Short:     `Apple silicon management commands`,
		Long:      `Apple silicon management commands.`,
		Namespace: "apple-silicon",
		Resource:  "server",
	}
}

func appleSiliconOs() *core.Command {
	return &core.Command{
		Short:     `OS management commands`,
		Long:      `OS management commands.`,
		Namespace: "apple-silicon",
		Resource:  "os",
	}
}

func appleSiliconServerType() *core.Command {
	return &core.Command{
		Short:     `Server-Types management commands`,
		Long:      `Server-Types management commands.`,
		Namespace: "apple-silicon",
		Resource:  "server-type",
	}
}

func appleSiliconPrivateNetwork() *core.Command {
	return &core.Command{
		Short: `Private network management command`,
		Long: `A Private Network allows you to interconnect your resources
in an isolated and private
network. Network reachability is limited to the
resources that are in the same VPC.

Note that a resource can be a part of multiple Private Networks.`,
		Namespace: "apple-silicon",
		Resource:  "private-network",
	}
}

func appleSiliconServerTypeList() *core.Command {
	return &core.Command{
		Short:     `List server types`,
		Long:      `List all technical details about Apple silicon server types available in the specified zone. Since there is only one Availability Zone for Apple silicon servers, the targeted value is ` + "`" + `fr-par-3` + "`" + `.`,
		Namespace: "apple-silicon",
		Resource:  "server-type",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(applesilicon.ListServerTypesRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ZoneArgSpec(scw.ZoneFrPar3),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*applesilicon.ListServerTypesRequest)

			client := core.ExtractClient(ctx)
			api := applesilicon.NewAPI(client)

			return api.ListServerTypes(request)
		},
	}
}

func appleSiliconServerTypeGet() *core.Command {
	return &core.Command{
		Short:     `Get a server type`,
		Long:      `Get technical details (CPU, disk size etc.) of a server type.`,
		Namespace: "apple-silicon",
		Resource:  "server-type",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(applesilicon.GetServerTypeRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "server-type",
				Short:      `Server type identifier`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(scw.ZoneFrPar3),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*applesilicon.GetServerTypeRequest)

			client := core.ExtractClient(ctx)
			api := applesilicon.NewAPI(client)

			return api.GetServerType(request)
		},
	}
}

func appleSiliconServerCreate() *core.Command {
	return &core.Command{
		Short:     `Create a server`,
		Long:      `Create a new server in the targeted zone, specifying its configuration including name and type.`,
		Namespace: "apple-silicon",
		Resource:  "server",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(applesilicon.CreateServerRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "name",
				Short:      `Create a server with this given name`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				Default:    core.RandomValueGenerator("as"),
			},
			core.ProjectIDArgSpec(),
			{
				Name:       "type",
				Short:      `Create a server of the given type`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "os-id",
				Short:      `Create a server with the given os_id`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "enable-vpc",
				Short:      `Activate the Private Network feature for this server`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "commitment-type",
				Short:      `Activate commitment for this server`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"duration_24h",
					"renewed_monthly",
					"none",
				},
			},
			{
				Name:       "public-bandwidth-bps",
				Short:      `Public bandwidth to configure for this server. Expressed in bits per second.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "runner-configuration.name",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "runner-configuration.url",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "runner-configuration.token",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "runner-configuration.provider",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_provider",
					"github",
					"gitlab",
				},
			},
			{
				Name:       "applied-runner-configurations.runner-configuration-ids.{index}",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar3),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*applesilicon.CreateServerRequest)

			client := core.ExtractClient(ctx)
			api := applesilicon.NewAPI(client)

			return api.CreateServer(request)
		},
	}
}

func appleSiliconServerList() *core.Command {
	return &core.Command{
		Short:     `List all servers`,
		Long:      `List all servers in the specified zone. By default, returned servers in the list are ordered by creation date in ascending order, though this can be modified via the ` + "`" + `order_by` + "`" + ` field.`,
		Namespace: "apple-silicon",
		Resource:  "server",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(applesilicon.ListServersRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `Sort order of the returned servers`,
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
				Short:      `Only list servers of this project ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "organization-id",
				Short:      `Only list servers of this Organization ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar3,
				scw.Zone(core.AllLocalities),
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*applesilicon.ListServersRequest)

			client := core.ExtractClient(ctx)
			api := applesilicon.NewAPI(client)
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

func appleSiliconOsList() *core.Command {
	return &core.Command{
		Short:     `List all Operating Systems (OS)`,
		Long:      `List all Operating Systems (OS). The response will include the total number of OS as well as their associated IDs, names and labels.`,
		Namespace: "apple-silicon",
		Resource:  "os",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(applesilicon.ListOSRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "server-type",
				Short:      `List of compatible server types`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "name",
				Short:      `Filter OS by name (note that "11.1" will return "11.1.2" and "11.1" but not "12"))`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar3,
				scw.Zone(core.AllLocalities),
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*applesilicon.ListOSRequest)

			client := core.ExtractClient(ctx)
			api := applesilicon.NewAPI(client)
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

func appleSiliconOsGet() *core.Command {
	return &core.Command{
		Short:     `Get an Operating System (OS)`,
		Long:      `Get an Operating System (OS).  The response will include the OS's unique ID as well as its name and label.`,
		Namespace: "apple-silicon",
		Resource:  "os",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(applesilicon.GetOSRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "os-id",
				Short:      `UUID of the OS you want to get`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(scw.ZoneFrPar3),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*applesilicon.GetOSRequest)

			client := core.ExtractClient(ctx)
			api := applesilicon.NewAPI(client)

			return api.GetOS(request)
		},
	}
}

func appleSiliconServerGet() *core.Command {
	return &core.Command{
		Short:     `Get a server`,
		Long:      `Retrieve information about an existing Apple silicon server, specified by its server ID. Its full details, including name, status and IP address, are returned in the response object.`,
		Namespace: "apple-silicon",
		Resource:  "server",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(applesilicon.GetServerRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "server-id",
				Short:      `UUID of the server you want to get`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(scw.ZoneFrPar3),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*applesilicon.GetServerRequest)

			client := core.ExtractClient(ctx)
			api := applesilicon.NewAPI(client)

			return api.GetServer(request)
		},
	}
}

func appleSiliconServerUpdate() *core.Command {
	return &core.Command{
		Short:     `Update a server`,
		Long:      `Update the parameters of an existing Apple silicon server, specified by its server ID.`,
		Namespace: "apple-silicon",
		Resource:  "server",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(applesilicon.UpdateServerRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "server-id",
				Short:      `UUID of the server you want to update`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "name",
				Short:      `Updated name for your server`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "schedule-deletion",
				Short:      `Specify whether the server should be flagged for automatic deletion`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "enable-vpc",
				Short:      `Activate or deactivate Private Network support for this server`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "commitment-type.commitment-type",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"duration_24h",
					"renewed_monthly",
					"none",
				},
			},
			{
				Name:       "public-bandwidth-bps",
				Short:      `Public bandwidth configured for this server`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "applied-runner-configurations.runner-configuration-ids.{index}",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar3),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*applesilicon.UpdateServerRequest)

			client := core.ExtractClient(ctx)
			api := applesilicon.NewAPI(client)

			return api.UpdateServer(request)
		},
	}
}

func appleSiliconServerDelete() *core.Command {
	return &core.Command{
		Short:     `Delete a server`,
		Long:      `Delete an existing Apple silicon server, specified by its server ID. Deleting a server is permanent, and cannot be undone. Note that the minimum allocation period for Apple silicon-as-a-service is 24 hours, meaning you cannot delete your server prior to that.`,
		Namespace: "apple-silicon",
		Resource:  "server",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(applesilicon.DeleteServerRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "server-id",
				Short:      `UUID of the server you want to delete`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(scw.ZoneFrPar3),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*applesilicon.DeleteServerRequest)

			client := core.ExtractClient(ctx)
			api := applesilicon.NewAPI(client)
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

func appleSiliconServerReboot() *core.Command {
	return &core.Command{
		Short:     `Reboot a server`,
		Long:      `Reboot an existing Apple silicon server, specified by its server ID.`,
		Namespace: "apple-silicon",
		Resource:  "server",
		Verb:      "reboot",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(applesilicon.RebootServerRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "server-id",
				Short:      `UUID of the server you want to reboot`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(scw.ZoneFrPar3),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*applesilicon.RebootServerRequest)

			client := core.ExtractClient(ctx)
			api := applesilicon.NewAPI(client)

			return api.RebootServer(request)
		},
	}
}

func appleSiliconServerReinstall() *core.Command {
	return &core.Command{
		Short:     `Reinstall a server`,
		Long:      `Reinstall an existing Apple silicon server (specified by its server ID) from a new image (OS). All the data on the disk is deleted and all configuration is reset to the default configuration values of the image (OS).`,
		Namespace: "apple-silicon",
		Resource:  "server",
		Verb:      "reinstall",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(applesilicon.ReinstallServerRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "server-id",
				Short:      `UUID of the server you want to reinstall`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "os-id",
				Short:      `Reinstall the server with the OS corresponding to the os_id`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar3),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*applesilicon.ReinstallServerRequest)

			client := core.ExtractClient(ctx)
			api := applesilicon.NewAPI(client)

			return api.ReinstallServer(request)
		},
	}
}

func appleSiliconPrivateNetworkAdd() *core.Command {
	return &core.Command{
		Short:     `Add a server to a Private Network`,
		Long:      `Add an Apple silicon server to a Private Network.`,
		Namespace: "apple-silicon",
		Resource:  "private-network",
		Verb:      "add",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(applesilicon.PrivateNetworkAPIAddServerPrivateNetworkRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "server-id",
				Short:      `ID of the server`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "private-network-id",
				Short:      `ID of the Private Network`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "ipam-ip-ids.{index}",
				Short:      `IPAM IDs of IPs to attach to the server`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar3,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*applesilicon.PrivateNetworkAPIAddServerPrivateNetworkRequest)

			client := core.ExtractClient(ctx)
			api := applesilicon.NewPrivateNetworkAPI(client)

			return api.AddServerPrivateNetwork(request)
		},
	}
}

func appleSiliconPrivateNetworkSet() *core.Command {
	return &core.Command{
		Short:     `Set multiple Private Networks on a server`,
		Long:      `Configure multiple Private Networks on an Apple silicon server.`,
		Namespace: "apple-silicon",
		Resource:  "private-network",
		Verb:      "set",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(applesilicon.PrivateNetworkAPISetServerPrivateNetworksRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "server-id",
				Short:      `ID of the server`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "per-private-network-ipam-ip-ids.{key}",
				Short:      `Object where the keys are the IDs of Private Networks and the values are arrays of IPAM IDs representing the IPs to assign to this Apple silicon server on the Private Network. If the array supplied for a Private Network is empty, the next available IP from the Private Network's CIDR block will automatically be used for attachment.`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar3,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*applesilicon.PrivateNetworkAPISetServerPrivateNetworksRequest)

			client := core.ExtractClient(ctx)
			api := applesilicon.NewPrivateNetworkAPI(client)

			return api.SetServerPrivateNetworks(request)
		},
	}
}

func appleSiliconPrivateNetworkList() *core.Command {
	return &core.Command{
		Short:     `List the Private Networks of a server`,
		Long:      `List the Private Networks of an Apple silicon server.`,
		Namespace: "apple-silicon",
		Resource:  "private-network",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(applesilicon.PrivateNetworkAPIListServerPrivateNetworksRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `Sort order for the returned Private Networks`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"created_at_asc",
					"created_at_desc",
					"updated_at_asc",
					"updated_at_desc",
				},
			},
			{
				Name:       "server-id",
				Short:      `Filter Private Networks by server ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "private-network-id",
				Short:      `Filter Private Networks by Private Network ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "project-id",
				Short:      `Filter Private Networks by Project ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "ipam-ip-ids.{index}",
				Short:      `Filter Private Networks by IPAM IP IDs`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "organization-id",
				Short:      `Filter Private Networks by Organization ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar3,
				scw.Zone(core.AllLocalities),
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*applesilicon.PrivateNetworkAPIListServerPrivateNetworksRequest)

			client := core.ExtractClient(ctx)
			api := applesilicon.NewPrivateNetworkAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			if request.Zone == scw.Zone(core.AllLocalities) {
				opts = append(opts, scw.WithZones(api.Zones()...))
				request.Zone = ""
			}
			resp, err := api.ListServerPrivateNetworks(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.ServerPrivateNetworks, nil
		},
	}
}

func appleSiliconPrivateNetworkDelete() *core.Command {
	return &core.Command{
		Short:     `Delete a Private Network`,
		Long:      `Delete a Private Network.`,
		Namespace: "apple-silicon",
		Resource:  "private-network",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(applesilicon.PrivateNetworkAPIDeleteServerPrivateNetworkRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "server-id",
				Short:      `ID of the server`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "private-network-id",
				Short:      `ID of the Private Network`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar3,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*applesilicon.PrivateNetworkAPIDeleteServerPrivateNetworkRequest)

			client := core.ExtractClient(ctx)
			api := applesilicon.NewPrivateNetworkAPI(client)
			e = api.DeleteServerPrivateNetwork(request)
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
