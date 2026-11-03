package secret

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/cockpit/v1"
	secret "github.com/scaleway/scaleway-sdk-go/api/secret/v1beta1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type secretCockpitArgs struct {
	SecretID string
	Region   scw.Region
}

func secretCockpitCommand() *core.Command {
	return &core.Command{
		Namespace: "secret",
		Resource:  "secret",
		Verb:      "cockpit",
		Short:     `Open the Cockpit Grafana dashboard for this secret`,
		ArgsType:  reflect.TypeFor[secretCockpitArgs](),
		ArgSpecs: core.ArgSpecs{
			{Name: "secret-id", Short: "Secret ID", Required: true, Positional: true},
			core.RegionArgSpec(),
		},
		Run: func(ctx context.Context, argsI any) (any, error) {
			args := argsI.(*secretCockpitArgs)

			api := secret.NewAPI(core.ExtractClient(ctx))
			sec, err := api.GetSecret(&secret.GetSecretRequest{
				Region:   args.Region,
				SecretID: args.SecretID,
			})
			if err != nil {
				return nil, err
			}

			return cockpit.OpenDashboard(
				ctx,
				"scw-secret-manager-overview",
				map[string]string{
					"region":      string(sec.Region),
					"secret_name": sec.Name,
				},
			)
		},
	}
}
