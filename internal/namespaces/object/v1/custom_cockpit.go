package object

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/cockpit/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type bucketCockpitArgs struct {
	Name   string
	Region scw.Region
}

func bucketCockpitCommand() *core.Command {
	return &core.Command{
		Namespace: "object",
		Resource:  "bucket",
		Verb:      "cockpit",
		Short:     `Open the Cockpit Grafana dashboard for this bucket`,
		ArgsType:  reflect.TypeFor[bucketCockpitArgs](),
		ArgSpecs: core.ArgSpecs{
			{Name: "name", Short: "Bucket name", Required: true, Positional: true},
			core.RegionArgSpec(),
		},
		Run: func(ctx context.Context, argsI any) (any, error) {
			args := argsI.(*bucketCockpitArgs)

			// The bucket name is the identifying variable, no get is needed.
			return cockpit.OpenDashboard(
				ctx,
				"scw-object-storage-overview",
				map[string]string{
					"region":     string(args.Region),
					"bucketName": args.Name,
				},
			)
		},
	}
}
