package instance

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/cockpit/v1"
	instance "github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type serverCockpitArgs struct {
	ServerID string
	Zone     scw.Zone
}

func serverCockpitCommand() *core.Command {
	return &core.Command{
		Namespace: "instance",
		Resource:  "server",
		Verb:      "cockpit",
		Short:     `Open the Cockpit Grafana dashboard for this server`,
		ArgsType:  reflect.TypeFor[serverCockpitArgs](),
		ArgSpecs: core.ArgSpecs{
			{Name: "server-id", Short: "Server ID", Required: true, Positional: true},
			core.ZoneArgSpec(),
		},
		Run: func(ctx context.Context, argsI any) (any, error) {
			args := argsI.(*serverCockpitArgs)

			api := instance.NewAPI(core.ExtractClient(ctx))
			server, err := api.GetServer(&instance.GetServerRequest{
				Zone:     args.Zone,
				ServerID: args.ServerID,
			})
			if err != nil {
				return nil, err
			}

			region, _ := server.Server.Zone.Region()

			return cockpit.OpenDashboard(ctx, "scw-instance-overview", map[string]string{
				"region":            string(region),
				"availability_zone": string(server.Server.Zone),
				"server_name":       server.Server.Name,
			})
		},
	}
}
