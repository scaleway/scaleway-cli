package k8s

import (
	"context"
	"fmt"
	"reflect"

	"github.com/scaleway/scaleway-cli/internal/core"
	k8s "github.com/scaleway/scaleway-sdk-go/api/k8s/v1beta4"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

//
// Marshalers
//

func versionListBuilder(c *core.Command) *core.Command {
	originalRun := c.Run

	c.Run = func(ctx context.Context, argsI interface{}) (interface{}, error) {
		originalRes, err := originalRun(ctx, argsI)
		if err != nil {
			return nil, err
		}

		versionsResponse := originalRes.(*k8s.ListVersionsResponse)
		return versionsResponse.Versions, nil
	}

	return c
}

type k8sVersionGetRequest struct {
	Version string
	Region  scw.Region
}

func k8sVersionGetCommand() *core.Command {
	return &core.Command{
		Short:     `Get a Kubernetes version`,
		Long:      `Get all the details about a specific Kubernetes version.`,
		Namespace: "k8s",
		Verb:      "get",
		Resource:  "version",
		ArgsType:  reflect.TypeOf(k8sVersionGetRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.RegionArgSpec(),
			{
				Name:       "version",
				Short:      "Version from which to get details",
				Required:   true,
				Positional: true,
			},
		},
		Run: k8sVersionGetRun,
	}
}

func k8sVersionGetRun(ctx context.Context, argsI interface{}) (i interface{}, e error) {
	args := argsI.(*k8sVersionGetRequest)

	client := core.ExtractClient(ctx)
	apiK8s := k8s.NewAPI(client)

	versions, err := apiK8s.ListVersions(&k8s.ListVersionsRequest{
		Region: args.Region,
	})

	if err != nil {
		return nil, err
	}

	for _, version := range versions.Versions {
		if version.Name == args.Version {
			return version, nil
		}
	}
	return nil, fmt.Errorf("version %s not found", args.Version)
}
