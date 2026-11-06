package file

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	billing "github.com/scaleway/scaleway-cli/v2/internal/namespaces/billing/v2beta1"
	file "github.com/scaleway/scaleway-sdk-go/api/file/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type filesystemChargeArgs struct {
	FilesystemID string
	Region       scw.Region
}

func filesystemChargeCommand() *core.Command {
	return &core.Command{
		Namespace: "file",
		Resource:  "filesystem",
		Verb:      "charge",
		Short:     `List all the charges of this filesystem since its creation`,
		ArgsType:  reflect.TypeFor[filesystemChargeArgs](),
		ArgSpecs: core.ArgSpecs{
			{Name: "filesystem-id", Short: "Filesystem ID", Required: true, Positional: true},
			core.RegionArgSpec(),
		},
		View: billing.ChargesView(),
		Run: func(ctx context.Context, argsI any) (any, error) {
			args := argsI.(*filesystemChargeArgs)

			api := file.NewAPI(core.ExtractClient(ctx))
			filesystem, err := api.GetFileSystem(&file.GetFileSystemRequest{
				Region:       args.Region,
				FilesystemID: args.FilesystemID,
			})
			if err != nil {
				return nil, err
			}

			return billing.ListChargesForResource(ctx, args.FilesystemID, filesystem.CreatedAt)
		},
	}
}
