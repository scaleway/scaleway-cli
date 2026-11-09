package k8s

import (
	"context"
	"fmt"

	"github.com/scaleway/scaleway-cli/v2/core"
	humanv2 "github.com/scaleway/scaleway-cli/v2/core/human/v2"
	"github.com/scaleway/scaleway-cli/v2/internal/semver"
	k8s "github.com/scaleway/scaleway-sdk-go/api/k8s/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func versionListBuilder(c *core.Command) *core.Command {
	c.AddInterceptors(
		func(ctx context.Context, argsI any, runner core.CommandRunner) (any, error) {
			originalRes, err := runner(ctx, argsI)
			if err != nil {
				return nil, err
			}

			versionsResponse := originalRes.(*k8s.ListVersionsResponse)

			return versionsResponse.Versions, nil
		},
	)

	return c
}

var versionLayout = humanv2.Layout{
	Sections: []humanv2.Section{
		{Path: "AvailableKubeletArgs", Title: "Available Kubelet Arguments"},
		{Path: "AvailableCnis", Title: "Available CNIs"},
		{Path: "AvailableFeatureGates", Title: "Available Feature Gates"},
		{Path: "AvailableAdmissionPlugins", Title: "Available Admission Plugins"},
	},
}

func getLatestK8SVersion(scwClient *scw.Client) (string, error) {
	api := k8s.NewAPI(scwClient)
	versions, err := api.ListVersions(&k8s.ListVersionsRequest{})
	if err != nil {
		return "", fmt.Errorf("could not get latest K8S version: %w", err)
	}

	names := make([]string, 0, len(versions.Versions))
	for _, v := range versions.Versions {
		names = append(names, v.Name)
	}

	return semver.LatestVersion(names)
}
