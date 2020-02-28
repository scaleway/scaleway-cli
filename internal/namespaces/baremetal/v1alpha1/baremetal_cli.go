// This file was automatically generated. DO NOT EDIT.
// If you have any remark or suggestion do not hesitate to open an issue.

package baremetal

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-sdk-go/api/baremetal/v1alpha1"
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
		baremetalIP(),
		baremetalBmc(),
		baremetalOffer(),
		baremetalServerList(),
	)
}
func baremetalRoot() *core.Command {
	return &core.Command{
		Short:     `Baremetal API`,
		Long:      ``,
		Namespace: "baremetal",
	}
}

func baremetalServer() *core.Command {
	return &core.Command{
		Short:     `A server is a denomination of a type of instances provided by Scaleway`,
		Long:      `A server is a denomination of a type of instances provided by Scaleway.`,
		Namespace: "baremetal",
		Resource:  "server",
	}
}

func baremetalOs() *core.Command {
	return &core.Command{
		Short:     `An Operating System (OS) is the underlying software installed on your server`,
		Long:      `An Operating System (OS) is the underlying software installed on your server.`,
		Namespace: "baremetal",
		Resource:  "os",
	}
}

func baremetalIP() *core.Command {
	return &core.Command{
		Short: `The IPs failovers could be attach to any server in the same zone.
A server could be linked with multiple failovers.
`,
		Long: `The IPs failovers could be attach to any server in the same zone.
A server could be linked with multiple failovers.
.`,
		Namespace: "baremetal",
		Resource:  "ip",
	}
}

func baremetalBmc() *core.Command {
	return &core.Command{
		Short:     `Baseboard Management Controller (BMC) offers a low-level access to your baremetal instance`,
		Long:      `Baseboard Management Controller (BMC) offers a low-level access to your baremetal instance.`,
		Namespace: "baremetal",
		Resource:  "bmc",
	}
}

func baremetalOffer() *core.Command {
	return &core.Command{
		Short:     `Commercial offers`,
		Long:      `Commercial offers.`,
		Namespace: "baremetal",
		Resource:  "offer",
	}
}

func baremetalServerList() *core.Command {
	return &core.Command{
		Short:     `List servers`,
		Long:      `List all created servers.`,
		Namespace: "baremetal",
		Resource:  "server",
		Verb:      "list",
		ArgsType:  reflect.TypeOf(baremetal.ListServersRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `Order of the servers`,
				Required:   false,
				EnumValues: []string{"created_at_asc", "created_at_desc"},
			},
			{
				Name:     "tags.{index}",
				Short:    `Filter servers by tags`,
				Required: false,
			},
			{
				Name:     "status.{index}",
				Short:    `Filter servers by status`,
				Required: false,
			},
			{
				Name:     "name",
				Short:    `Filter servers by name`,
				Required: false,
			},
			{
				Name:     "organization-id",
				Short:    `Filter servers by organization ID`,
				Required: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar2),
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
				Short:   "List all servers on your default zone",
				Request: `null`,
			},
		},
	}
}
