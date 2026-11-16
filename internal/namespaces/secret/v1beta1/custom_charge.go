package secret

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	billing "github.com/scaleway/scaleway-cli/v2/internal/namespaces/billing/v2beta1"
	secret "github.com/scaleway/scaleway-sdk-go/api/secret/v1beta1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type secretChargeArgs struct {
	SecretID string
	Region   scw.Region
}

func secretChargeCommand() *core.Command {
	return &core.Command{
		Namespace: "secret",
		Resource:  "secret",
		Verb:      "charge",
		Short:     `List all the charges of this secret since its creation`,
		ArgsType:  reflect.TypeFor[secretChargeArgs](),
		ArgSpecs: core.ArgSpecs{
			{Name: "secret-id", Short: "Secret ID", Required: true, Positional: true},
			core.RegionArgSpec(),
		},
		View: billing.ChargesView(),
		Run: func(ctx context.Context, argsI any) (any, error) {
			args := argsI.(*secretChargeArgs)

			api := secret.NewAPI(core.ExtractClient(ctx))
			sec, err := api.GetSecret(&secret.GetSecretRequest{
				Region:   args.Region,
				SecretID: args.SecretID,
			})
			if err != nil {
				return nil, err
			}

			return billing.ListChargesForResource(ctx, args.SecretID, sec.CreatedAt)
		},
	}
}
