package key_manager

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	billing "github.com/scaleway/scaleway-cli/v2/internal/namespaces/billing/v2beta1"
	keymanager "github.com/scaleway/scaleway-sdk-go/api/key_manager/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type keyChargeArgs struct {
	KeyID  string
	Region scw.Region
}

func keyChargeCommand() *core.Command {
	return &core.Command{
		Namespace: "keymanager",
		Resource:  "key",
		Verb:      "charge",
		Short:     `List all the charges of this key since its creation`,
		ArgsType:  reflect.TypeFor[keyChargeArgs](),
		ArgSpecs: core.ArgSpecs{
			{Name: "key-id", Short: "Key ID", Required: true, Positional: true},
			core.RegionArgSpec(),
		},
		View: billing.ChargesView(),
		Run: func(ctx context.Context, argsI any) (any, error) {
			args := argsI.(*keyChargeArgs)

			api := keymanager.NewAPI(core.ExtractClient(ctx))
			key, err := api.GetKey(&keymanager.GetKeyRequest{
				Region: args.Region,
				KeyID:  args.KeyID,
			})
			if err != nil {
				return nil, err
			}

			return billing.ListChargesForResource(ctx, args.KeyID, key.CreatedAt)
		},
	}
}
