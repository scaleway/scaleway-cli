// This file was automatically generated. DO NOT EDIT.
// If you have any remark or suggestion do not hesitate to open an issue.

package baremetal

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-sdk-go/api/baremetal/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// always import dependencies
var (
	_ = scw.RegionFrPar
)

func GetGeneratedCommands() *core.Commands {
	return core.NewCommands(
		baremetalRoot(),
		baremetalServer(),
		baremetalOs(),
		baremetalBmc(),
		baremetalServerList(),
		baremetalServerGet(),
		baremetalServerCreate(),
		baremetalServerUpdate(),
		baremetalServerInstall(),
		baremetalServerDelete(),
		baremetalServerReboot(),
		baremetalServerStart(),
		baremetalServerStop(),
		baremetalBmcStart(),
		baremetalBmcGet(),
		baremetalBmcStop(),
		baremetalOsList(),
		baremetalOsGet(),
	)
}
func baremetalRoot() *core.Command {
	return &core.Command{
		Short:     `Elastic metal API`,
		Long:      ``,
		Namespace: "baremetal",
	}
}

func baremetalServer() *core.Command {
	return &core.Command{
		Short:     `Server management commands`,
		Long:      `A server is a denomination of a type of instances provided by Scaleway`,
		Namespace: "baremetal",
		Resource:  "server",
	}
}

func baremetalOs() *core.Command {
	return &core.Command{
		Short:     `Operating System (OS) management commands`,
		Long:      `An Operating System (OS) is the underlying software installed on your server`,
		Namespace: "baremetal",
		Resource:  "os",
	}
}

func baremetalBmc() *core.Command {
	return &core.Command{
		Short: `Baseboard Management Controller (BMC) management commands`,
		Long: `Baseboard Management Controller (BMC) allows you to remotely access the low-level parameters of your dedicated server.
For instance, your KVM-IP management console could be accessed with it.
You need first to create an option Remote Access. You will find the ID and the price with a call to listOffers (https://developers.scaleway.com/en/products/baremetal/api/#get-78db92). Then you can add the option https://developers.scaleway.com/en/products/baremetal/api/#post-b14abd. Do not forget to delete the Option.
Then you need to create Remote Access https://developers.scaleway.com/en/products/baremetal/api/#post-1af723.
And finally Get Remote Access to get the login/password https://developers.scaleway.com/en/products/baremetal/api/#get-cefc0f.
`,
		Namespace: "baremetal",
		Resource:  "bmc",
	}
}

func baremetalServerList() *core.Command {
	return &core.Command{
		Short:     `List elastic metal servers for organization`,
		Long:      `List elastic metal servers for organization.`,
		Namespace: "baremetal",
		Resource:  "server",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(baremetal.ListServersRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `Order of the servers`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"created_at_asc", "created_at_desc"},
			},
			{
				Name:       "tags.{index}",
				Short:      `Filter by tags`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "status.{index}",
				Short:      `Filter by status`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "name",
				Short:      `Filter by name`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "project-id",
				Short:      `Filter by project ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "option-id",
				Short:      `Filter by option ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "organization-id",
				Short:      `Filter by organization ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*baremetal.ListServersRequest)

			client := core.ExtractClient(ctx)
			api := baremetal.NewAPI(client)
			resp, err := api.ListServers(request, scw.WithAllPages())
			if err != nil {
				return nil, err
			}
			return resp.Servers, nil

		},
		Examples: []*core.Example{
			{
				Short:    "List all servers on your default zone",
				ArgsJSON: `null`,
			},
		},
	}
}

func baremetalServerGet() *core.Command {
	return &core.Command{
		Short:     `Get a specific elastic metal server`,
		Long:      `Get the server associated with the given ID.`,
		Namespace: "baremetal",
		Resource:  "server",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(baremetal.GetServerRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "server-id",
				Short:      `ID of the server`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*baremetal.GetServerRequest)

			client := core.ExtractClient(ctx)
			api := baremetal.NewAPI(client)
			return api.GetServer(request)

		},
		Examples: []*core.Example{
			{
				Short:    "Get a given server",
				ArgsJSON: `{"server_id":"11111111-1111-1111-1111-111111111111"}`,
			},
		},
	}
}

func baremetalServerCreate() *core.Command {
	return &core.Command{
		Short:     `Create an elastic metal server`,
		Long:      `Create a new elastic metal server. Once the server is created, you probably want to install an OS.`,
		Namespace: "baremetal",
		Resource:  "server",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(baremetal.CreateServerRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "offer-id",
				Short:      `Offer ID of the new server`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.ProjectIDArgSpec(),
			{
				Name:       "name",
				Short:      `Name of the server (≠hostname)`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "description",
				Short:      `Description associated to the server, max 255 characters`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `Tags to associate to the server`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "install.os-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "install.hostname",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "install.ssh-key-ids.{index}",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "install.user",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "install.password",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "install.service-user",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "install.service-password",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "option-ids.{index}",
				Short:      `IDs of options to enable on server`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.OrganizationIDArgSpec(),
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*baremetal.CreateServerRequest)

			client := core.ExtractClient(ctx)
			api := baremetal.NewAPI(client)
			return api.CreateServer(request)

		},
		Examples: []*core.Example{
			{
				Short:    "Create instance a default elastic metal instance",
				ArgsJSON: `null`,
			},
		},
	}
}

func baremetalServerUpdate() *core.Command {
	return &core.Command{
		Short:     `Update an elastic metal server`,
		Long:      `Update the server associated with the given ID.`,
		Namespace: "baremetal",
		Resource:  "server",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(baremetal.UpdateServerRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "server-id",
				Short:      `ID of the server to update`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "name",
				Short:      `Name of the server (≠hostname), not updated if null`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "description",
				Short:      `Description associated to the server, max 255 characters, not updated if null`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `Tags associated to the server, not updated if null`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*baremetal.UpdateServerRequest)

			client := core.ExtractClient(ctx)
			api := baremetal.NewAPI(client)
			return api.UpdateServer(request)

		},
	}
}

func baremetalServerInstall() *core.Command {
	return &core.Command{
		Short:     `Install an elastic metal server`,
		Long:      `Install an OS on the server associated with the given ID.`,
		Namespace: "baremetal",
		Resource:  "server",
		Verb:      "install",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(baremetal.InstallServerRequest{}),
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
				Short:      `ID of the OS to install on the server`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "hostname",
				Short:      `Hostname of the server`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "ssh-key-ids.{index}",
				Short:      `SSH key IDs authorized on the server`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "user",
				Short:      `User used for the installation`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "password",
				Short:      `Password used for the installation`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "service-user",
				Short:      `User used for the service to install`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "service-password",
				Short:      `Password used for the service to install`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*baremetal.InstallServerRequest)

			client := core.ExtractClient(ctx)
			api := baremetal.NewAPI(client)
			return api.InstallServer(request)

		},
		Examples: []*core.Example{
			{
				Short:    "Install an OS on a given server with a particular SSH key ID",
				ArgsJSON: `{"os_id":"11111111-1111-1111-1111-111111111111","server_id":"11111111-1111-1111-1111-111111111111","ssh_key_ids":["11111111-1111-1111-1111-111111111111"]}`,
			},
		},
		SeeAlsos: []*core.SeeAlso{
			{
				Command: "scw account ssh-key list",
				Short:   "List all SSH keys",
			},
			{
				Command: "scw baremetal os list",
				Short:   "List OS (useful to get all OS IDs)",
			},
			{
				Command: "scw baremetal server create",
				Short:   "Create an elastic metal server",
			},
		},
	}
}

func baremetalServerDelete() *core.Command {
	return &core.Command{
		Short:     `Delete an elastic metal server`,
		Long:      `Delete the server associated with the given ID.`,
		Namespace: "baremetal",
		Resource:  "server",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(baremetal.DeleteServerRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "server-id",
				Short:      `ID of the server to delete`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*baremetal.DeleteServerRequest)

			client := core.ExtractClient(ctx)
			api := baremetal.NewAPI(client)
			return api.DeleteServer(request)

		},
		Examples: []*core.Example{
			{
				Short:    "Delete an elastic metal server",
				ArgsJSON: `{"server_id":"11111111-1111-1111-1111-111111111111"}`,
			},
		},
	}
}

func baremetalServerReboot() *core.Command {
	return &core.Command{
		Short:     `Reboot an elastic metal server`,
		Long:      `Reboot the server associated with the given ID, use boot param to reboot in rescue.`,
		Namespace: "baremetal",
		Resource:  "server",
		Verb:      "reboot",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(baremetal.RebootServerRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "server-id",
				Short:      `ID of the server to reboot`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "boot-type",
				Short:      `The type of boot`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"unknown_boot_type", "normal", "rescue"},
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*baremetal.RebootServerRequest)

			client := core.ExtractClient(ctx)
			api := baremetal.NewAPI(client)
			return api.RebootServer(request)

		},
		Examples: []*core.Example{
			{
				Short:    "Reboot a server using the same os",
				ArgsJSON: `{"server_id":"11111111-1111-1111-1111-111111111111"}`,
			},
			{
				Short:    "Reboot a server in rescue mode",
				ArgsJSON: `{"boot_type":"rescue","server_id":"11111111-1111-1111-1111-111111111111"}`,
			},
		},
	}
}

func baremetalServerStart() *core.Command {
	return &core.Command{
		Short:     `Start an elastic metal server`,
		Long:      `Start the server associated with the given ID.`,
		Namespace: "baremetal",
		Resource:  "server",
		Verb:      "start",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(baremetal.StartServerRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "server-id",
				Short:      `ID of the server to start`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "boot-type",
				Short:      `The type of boot`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"unknown_boot_type", "normal", "rescue"},
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*baremetal.StartServerRequest)

			client := core.ExtractClient(ctx)
			api := baremetal.NewAPI(client)
			return api.StartServer(request)

		},
		Examples: []*core.Example{
			{
				Short:    "Start an elastic metal server",
				ArgsJSON: `{"server_id":"11111111-1111-1111-1111-111111111111"}`,
			},
			{
				Short:    "Start a server in rescue mode",
				ArgsJSON: `{"boot_type":"rescue","server_id":"11111111-1111-1111-1111-111111111111"}`,
			},
		},
	}
}

func baremetalServerStop() *core.Command {
	return &core.Command{
		Short:     `Stop an elastic metal server`,
		Long:      `Stop the server associated with the given ID.`,
		Namespace: "baremetal",
		Resource:  "server",
		Verb:      "stop",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(baremetal.StopServerRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "server-id",
				Short:      `ID of the server to stop`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*baremetal.StopServerRequest)

			client := core.ExtractClient(ctx)
			api := baremetal.NewAPI(client)
			return api.StopServer(request)

		},
		Examples: []*core.Example{
			{
				Short:    "Stop an elastic metal server",
				ArgsJSON: `{"server_id":"11111111-1111-1111-1111-111111111111"}`,
			},
		},
	}
}

func baremetalBmcStart() *core.Command {
	return &core.Command{
		Short: `Start BMC (Baseboard Management Controller) access for a given elastic metal server`,
		Long: `Start BMC (Baseboard Management Controller) access associated with the given ID.
The BMC (Baseboard Management Controller) access is available one hour after the installation of the server.
You need first to create an option Remote Access. You will find the ID and the price with a call to listOffers (https://developers.scaleway.com/en/products/baremetal/api/#get-78db92). Then you can add the option https://developers.scaleway.com/en/products/baremetal/api/#post-b14abd. Do not forget to delete the Option.
 After start BMC, you need to Get Remote Access to get the login/password https://developers.scaleway.com/en/products/baremetal/api/#get-cefc0f.
`,
		Namespace: "baremetal",
		Resource:  "bmc",
		Verb:      "start",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(baremetal.StartBMCAccessRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "server-id",
				Short:      `ID of the server`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "ip",
				Short:      `The IP authorized to connect to the given server`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*baremetal.StartBMCAccessRequest)

			client := core.ExtractClient(ctx)
			api := baremetal.NewAPI(client)
			return api.StartBMCAccess(request)

		},
	}
}

func baremetalBmcGet() *core.Command {
	return &core.Command{
		Short:     `Get BMC (Baseboard Management Controller) access for a given elastic metal server`,
		Long:      `Get the BMC (Baseboard Management Controller) access associated with the given ID.`,
		Namespace: "baremetal",
		Resource:  "bmc",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(baremetal.GetBMCAccessRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "server-id",
				Short:      `ID of the server`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*baremetal.GetBMCAccessRequest)

			client := core.ExtractClient(ctx)
			api := baremetal.NewAPI(client)
			return api.GetBMCAccess(request)

		},
	}
}

func baremetalBmcStop() *core.Command {
	return &core.Command{
		Short:     `Stop BMC (Baseboard Management Controller) access for a given elastic metal server`,
		Long:      `Stop BMC (Baseboard Management Controller) access associated with the given ID.`,
		Namespace: "baremetal",
		Resource:  "bmc",
		Verb:      "stop",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(baremetal.StopBMCAccessRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "server-id",
				Short:      `ID of the server`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*baremetal.StopBMCAccessRequest)

			client := core.ExtractClient(ctx)
			api := baremetal.NewAPI(client)
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

func baremetalOsList() *core.Command {
	return &core.Command{
		Short:     `List all available OS that can be install on an elastic metal server`,
		Long:      `List all available OS that can be install on an elastic metal server.`,
		Namespace: "baremetal",
		Resource:  "os",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(baremetal.ListOSRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "offer-id",
				Short:      `Filter OS by offer ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*baremetal.ListOSRequest)

			client := core.ExtractClient(ctx)
			api := baremetal.NewAPI(client)
			resp, err := api.ListOS(request, scw.WithAllPages())
			if err != nil {
				return nil, err
			}
			return resp.Os, nil

		},
	}
}

func baremetalOsGet() *core.Command {
	return &core.Command{
		Short:     `Get an OS with a given ID`,
		Long:      `Return specific OS for the given ID.`,
		Namespace: "baremetal",
		Resource:  "os",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(baremetal.GetOSRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "os-id",
				Short:      `ID of the OS`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*baremetal.GetOSRequest)

			client := core.ExtractClient(ctx)
			api := baremetal.NewAPI(client)
			return api.GetOS(request)

		},
		Examples: []*core.Example{
			{
				Short:    "Get a specific OS ID",
				ArgsJSON: `{}`,
			},
		},
	}
}
