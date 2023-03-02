// This file was automatically generated. DO NOT EDIT.
// If you have any remark or suggestion do not hesitate to open an issue.

package baremetal

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/internal/core"
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
		baremetalOffer(),
		baremetalOs(),
		baremetalBmc(),
		baremetalOptions(),
		baremetalSettings(),
		baremetalPrivateNetwork(),
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
		baremetalOptionsAdd(),
		baremetalOptionsDelete(),
		baremetalOfferList(),
		baremetalOfferGet(),
		baremetalOptionsGet(),
		baremetalOptionsList(),
		baremetalSettingsList(),
		baremetalSettingsUpdate(),
		baremetalOsList(),
		baremetalOsGet(),
		baremetalPrivateNetworkAdd(),
		baremetalPrivateNetworkSet(),
		baremetalPrivateNetworkList(),
		baremetalPrivateNetworkDelete(),
	)
}
func baremetalRoot() *core.Command {
	return &core.Command{
		Short:     `Elastic metal API`,
		Long:      `This API allows to manage your Bare metal server.`,
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

func baremetalOffer() *core.Command {
	return &core.Command{
		Short: `Server offer management commands`,
		Long: `Server offers will answer with all different elastic metal server ranges available in a given zone.
Each of them will contain all the features of the server (cpus, memories, disks) with their associated pricing.
`,
		Namespace: "baremetal",
		Resource:  "offer",
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

func baremetalOptions() *core.Command {
	return &core.Command{
		Short: `Server options management commands`,
		Long: `A Server has additional options that let you personalize it to better fit your needs.
`,
		Namespace: "baremetal",
		Resource:  "options",
	}
}

func baremetalSettings() *core.Command {
	return &core.Command{
		Short: `Settings management commands`,
		Long: `Allows to configure the general settings for your elastic metal server.     
`,
		Namespace: "baremetal",
		Resource:  "settings",
	}
}

func baremetalPrivateNetwork() *core.Command {
	return &core.Command{
		Short: `Private network management command`,
		Long: `A private network allows interconnecting your resources
(servers, instances, ...) in an isolated and private
network. The network reachability is limited to the
resources that are on the same private network.  A VLAN
interface is available on the server and can be freely
managed (adding IP addresses, shutdown interface...).

Note that a resource can be a part of multiple private networks.
`,
		Namespace: "baremetal",
		Resource:  "private-network",
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
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.Zone(core.AllLocalities)),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*baremetal.ListServersRequest)

			client := core.ExtractClient(ctx)
			api := baremetal.NewAPI(client)
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
				Short:      `ID of the OS to install on the server`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "install.hostname",
				Short:      `Hostname of the server`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "install.ssh-key-ids.{index}",
				Short:      `SSH key IDs authorized on the server`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "install.user",
				Short:      `User used for the installation`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "install.password",
				Short:      `Password used for the installation`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "install.service-user",
				Short:      `User used for the service to install`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "install.service-password",
				Short:      `Password used for the service to install`,
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
 After start BMC, you need to Get Remote Access to get the login/password https://developers.scaleway.com/en/products/baremetal/api/#get-cefc0f.`,
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

func baremetalOptionsAdd() *core.Command {
	return &core.Command{
		Short:     `Add server option`,
		Long:      `Add an option to a specific server.`,
		Namespace: "baremetal",
		Resource:  "options",
		Verb:      "add",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(baremetal.AddOptionServerRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "server-id",
				Short:      `ID of the server`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "option-id",
				Short:      `ID of the option to add`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "expires-at",
				Short:      `Auto expire the option after this date`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*baremetal.AddOptionServerRequest)

			client := core.ExtractClient(ctx)
			api := baremetal.NewAPI(client)
			return api.AddOptionServer(request)

		},
		Examples: []*core.Example{
			{
				Short:    "Add a given option to a server",
				ArgsJSON: `{"option_id":"11111111-1111-1111-1111-111111111111","server_id":"11111111-1111-1111-1111-111111111111"}`,
			},
		},
	}
}

func baremetalOptionsDelete() *core.Command {
	return &core.Command{
		Short:     `Delete server option`,
		Long:      `Delete an option from a specific server.`,
		Namespace: "baremetal",
		Resource:  "options",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(baremetal.DeleteOptionServerRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "server-id",
				Short:      `ID of the server`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "option-id",
				Short:      `ID of the option to delete`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*baremetal.DeleteOptionServerRequest)

			client := core.ExtractClient(ctx)
			api := baremetal.NewAPI(client)
			return api.DeleteOptionServer(request)

		},
		Examples: []*core.Example{
			{
				Short:    "Delete a given option from a server",
				ArgsJSON: `{"option_id":"11111111-1111-1111-1111-111111111111","server_id":"11111111-1111-1111-1111-111111111111"}`,
			},
		},
	}
}

func baremetalOfferList() *core.Command {
	return &core.Command{
		Short:     `List offers`,
		Long:      `List all available server offers.`,
		Namespace: "baremetal",
		Resource:  "offer",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(baremetal.ListOffersRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "subscription-period",
				Short:      `Period of subscription to filter offers`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"unknown_subscription_period", "hourly", "monthly"},
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.Zone(core.AllLocalities)),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*baremetal.ListOffersRequest)

			client := core.ExtractClient(ctx)
			api := baremetal.NewAPI(client)
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
		Examples: []*core.Example{
			{
				Short:    "List all server offers in the default zone",
				ArgsJSON: `null`,
			},
			{
				Short:    "List all server offers in fr-par-1 zone",
				ArgsJSON: `{"zone":"fr-par-1"}`,
			},
		},
	}
}

func baremetalOfferGet() *core.Command {
	return &core.Command{
		Short:     `Get offer`,
		Long:      `Return specific offer for the given ID.`,
		Namespace: "baremetal",
		Resource:  "offer",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(baremetal.GetOfferRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "offer-id",
				Short:      `ID of the researched Offer`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*baremetal.GetOfferRequest)

			client := core.ExtractClient(ctx)
			api := baremetal.NewAPI(client)
			return api.GetOffer(request)

		},
		Examples: []*core.Example{
			{
				Short:    "Get a server offer with the given ID",
				ArgsJSON: `{"offer_id":"11111111-1111-1111-1111-111111111111","zone":"fr-par-1"}`,
			},
		},
	}
}

func baremetalOptionsGet() *core.Command {
	return &core.Command{
		Short:     `Get option`,
		Long:      `Return specific option for the given ID.`,
		Namespace: "baremetal",
		Resource:  "options",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(baremetal.GetOptionRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "option-id",
				Short:      `ID of the option`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*baremetal.GetOptionRequest)

			client := core.ExtractClient(ctx)
			api := baremetal.NewAPI(client)
			return api.GetOption(request)

		},
		Examples: []*core.Example{
			{
				Short:    "Get a server option with the given ID",
				ArgsJSON: `{"option_id":"11111111-1111-1111-1111-111111111111","zone":"fr-par-1"}`,
			},
		},
	}
}

func baremetalOptionsList() *core.Command {
	return &core.Command{
		Short:     `List options`,
		Long:      `List all options matching with filters.`,
		Namespace: "baremetal",
		Resource:  "options",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(baremetal.ListOptionsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "offer-id",
				Short:      `Filter options by offer_id`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "name",
				Short:      `Filter options by name`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.Zone(core.AllLocalities)),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*baremetal.ListOptionsRequest)

			client := core.ExtractClient(ctx)
			api := baremetal.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			if request.Zone == scw.Zone(core.AllLocalities) {
				opts = append(opts, scw.WithZones(api.Zones()...))
				request.Zone = ""
			}
			resp, err := api.ListOptions(request, opts...)
			if err != nil {
				return nil, err
			}
			return resp.Options, nil

		},
		Examples: []*core.Example{
			{
				Short:    "List all server options in the default zone",
				ArgsJSON: `null`,
			},
			{
				Short:    "List all server options in fr-par-1 zone",
				ArgsJSON: `{"zone":"fr-par-1"}`,
			},
		},
	}
}

func baremetalSettingsList() *core.Command {
	return &core.Command{
		Short:     `List all settings`,
		Long:      `Return all settings for a project ID.`,
		Namespace: "baremetal",
		Resource:  "settings",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(baremetal.ListSettingsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `Order the response`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"created_at_asc", "created_at_desc"},
			},
			{
				Name:       "project-id",
				Short:      `ID of the project`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.Zone(core.AllLocalities)),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*baremetal.ListSettingsRequest)

			client := core.ExtractClient(ctx)
			api := baremetal.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			if request.Zone == scw.Zone(core.AllLocalities) {
				opts = append(opts, scw.WithZones(api.Zones()...))
				request.Zone = ""
			}
			resp, err := api.ListSettings(request, opts...)
			if err != nil {
				return nil, err
			}
			return resp.Settings, nil

		},
	}
}

func baremetalSettingsUpdate() *core.Command {
	return &core.Command{
		Short:     `Update setting`,
		Long:      `Update a setting for a project ID (enable or disable).`,
		Namespace: "baremetal",
		Resource:  "settings",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(baremetal.UpdateSettingRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "setting-id",
				Short:      `ID of the setting`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "enabled",
				Short:      `Enable/Disable the setting`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*baremetal.UpdateSettingRequest)

			client := core.ExtractClient(ctx)
			api := baremetal.NewAPI(client)
			return api.UpdateSetting(request)

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
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1, scw.Zone(core.AllLocalities)),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*baremetal.ListOSRequest)

			client := core.ExtractClient(ctx)
			api := baremetal.NewAPI(client)
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

func baremetalPrivateNetworkAdd() *core.Command {
	return &core.Command{
		Short:     `Add a server to a private network`,
		Long:      `Add a server to a private network.`,
		Namespace: "baremetal",
		Resource:  "private-network",
		Verb:      "add",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(baremetal.PrivateNetworkAPIAddServerPrivateNetworkRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "server-id",
				Short:      `The ID of the server`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "private-network-id",
				Short:      `The ID of the private network`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar2),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*baremetal.PrivateNetworkAPIAddServerPrivateNetworkRequest)

			client := core.ExtractClient(ctx)
			api := baremetal.NewPrivateNetworkAPI(client)
			return api.AddServerPrivateNetwork(request)

		},
	}
}

func baremetalPrivateNetworkSet() *core.Command {
	return &core.Command{
		Short:     `Set multiple private networks on a server`,
		Long:      `Set multiple private networks on a server.`,
		Namespace: "baremetal",
		Resource:  "private-network",
		Verb:      "set",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(baremetal.PrivateNetworkAPISetServerPrivateNetworksRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "server-id",
				Short:      `The ID of the server`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "private-network-ids.{index}",
				Short:      `The IDs of the private networks`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar2),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*baremetal.PrivateNetworkAPISetServerPrivateNetworksRequest)

			client := core.ExtractClient(ctx)
			api := baremetal.NewPrivateNetworkAPI(client)
			return api.SetServerPrivateNetworks(request)

		},
	}
}

func baremetalPrivateNetworkList() *core.Command {
	return &core.Command{
		Short:     `List the private networks of a server`,
		Long:      `List the private networks of a server.`,
		Namespace: "baremetal",
		Resource:  "private-network",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(baremetal.PrivateNetworkAPIListServerPrivateNetworksRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `The sort order for the returned private networks`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"created_at_asc", "created_at_desc", "updated_at_asc", "updated_at_desc"},
			},
			{
				Name:       "server-id",
				Short:      `Filter private networks by server ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "private-network-id",
				Short:      `Filter private networks by private network ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "project-id",
				Short:      `Filter private networks by project ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "organization-id",
				Short:      `Filter private networks by organization ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar2, scw.Zone(core.AllLocalities)),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*baremetal.PrivateNetworkAPIListServerPrivateNetworksRequest)

			client := core.ExtractClient(ctx)
			api := baremetal.NewPrivateNetworkAPI(client)
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

func baremetalPrivateNetworkDelete() *core.Command {
	return &core.Command{
		Short:     `Delete a private network`,
		Long:      `Delete a private network.`,
		Namespace: "baremetal",
		Resource:  "private-network",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(baremetal.PrivateNetworkAPIDeleteServerPrivateNetworkRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "server-id",
				Short:      `The ID of the server`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "private-network-id",
				Short:      `The ID of the private network`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar2),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*baremetal.PrivateNetworkAPIDeleteServerPrivateNetworkRequest)

			client := core.ExtractClient(ctx)
			api := baremetal.NewPrivateNetworkAPI(client)
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
