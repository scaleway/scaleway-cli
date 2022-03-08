// This file was automatically generated. DO NOT EDIT.
// If you have any remark or suggestion do not hesitate to open an issue.

package applesilicon

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-sdk-go/api/applesilicon/v1alpha1"
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
	)
}
func appleSiliconRoot() *core.Command {
	return &core.Command{
		Short: `Apple silicon API`,
		Long: `Scaleway Apple silicon M1 as-a-Service is built using the latest generation of Apple Mac mini hardware (fifth generation).

These dedicated Mac mini M1s are designed for developing, building, testing, and signing applications for Apple devices, including iPhones, iPads, Mac computers and much more.

Get set to explore, learn and build on a dedicated Mac mini M1 with more performance and speed than you ever thought possible.

**Apple silicon as a Service comes with a minimum allocation period of 24 hours**.

Mac mini and macOS are trademarks of Apple Inc., registered in the U.S. and other countries and regions.
IOS is a trademark or registered trademark of Cisco in the U.S. and other countries and is used by Apple under license.
Scaleway is not affiliated with Apple Inc.
`,
		Namespace: "apple-silicon",
	}
}

func appleSiliconServer() *core.Command {
	return &core.Command{
		Short:     `Apple silicon management commands`,
		Long:      `Apple silicon management commands`,
		Namespace: "apple-silicon",
		Resource:  "server",
	}
}

func appleSiliconOs() *core.Command {
	return &core.Command{
		Short:     `OS management commands`,
		Long:      `OS management commands`,
		Namespace: "apple-silicon",
		Resource:  "os",
	}
}

func appleSiliconServerType() *core.Command {
	return &core.Command{
		Short:     `Server-Types management commands`,
		Long:      `Server-Types management commands`,
		Namespace: "apple-silicon",
		Resource:  "server-type",
	}
}

func appleSiliconServerTypeList() *core.Command {
	return &core.Command{
		Short:     `List server types`,
		Long:      `List all server types technical details.`,
		Namespace: "apple-silicon",
		Resource:  "server-type",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(applesilicon.ListServerTypesRequest{}),
		ArgSpecs: core.ArgSpecs{},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
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
		Long:      `Get a server technical details.`,
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
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
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
		Long:      `Create a server.`,
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
				Default:    core.DefaultValueSetter("M1-M"),
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
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
		Long:      `List all servers.`,
		Namespace: "apple-silicon",
		Resource:  "server",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(applesilicon.ListServersRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `The sort order of the returned servers`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"created_at_asc", "created_at_desc"},
			},
			{
				Name:       "project-id",
				Short:      `List only servers of this project ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "organization-id",
				Short:      `List only servers of this organization ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*applesilicon.ListServersRequest)

			client := core.ExtractClient(ctx)
			api := applesilicon.NewAPI(client)
			resp, err := api.ListServers(request, scw.WithAllPages())
			if err != nil {
				return nil, err
			}
			return resp.Servers, nil

		},
	}
}

func appleSiliconOsList() *core.Command {
	return &core.Command{
		Short:     `List all Operating System (OS)`,
		Long:      `List all Operating System (OS).`,
		Namespace: "apple-silicon",
		Resource:  "os",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(applesilicon.ListOSRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "server-type",
				Short:      `List of compatible server type`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "name",
				Short:      `Filter os by name (for eg. "11.1" will return "11.1.2" and "11.1" but not "12")`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*applesilicon.ListOSRequest)

			client := core.ExtractClient(ctx)
			api := applesilicon.NewAPI(client)
			resp, err := api.ListOS(request, scw.WithAllPages())
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
		Long:      `Get an Operating System (OS).`,
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
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
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
		Long:      `Get a server.`,
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
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
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
		Long:      `Update a server.`,
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
				Positional: false,
			},
			{
				Name:       "name",
				Short:      `Updated name for your server`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
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
		Long:      `Delete a server.`,
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
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
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
		Long:      `Reboot a server.`,
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
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
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
		Long:      `Reinstall a server.`,
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
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*applesilicon.ReinstallServerRequest)

			client := core.ExtractClient(ctx)
			api := applesilicon.NewAPI(client)
			return api.ReinstallServer(request)

		},
	}
}
