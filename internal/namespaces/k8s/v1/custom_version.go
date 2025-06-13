package k8s

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-version"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
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

func versionMarshalerFunc(i any, opt *human.MarshalOpt) (string, error) {
	type tmp k8s.Version
	version := tmp(i.(k8s.Version))

	// Sections
	opt.Sections = []*human.MarshalSection{
		{
			FieldName: "AvailableKubeletArgs",
			Title:     "Available Kubelet Arguments",
		},
		{
			FieldName: "AvailableCnis",
			Title:     "Available CNIs",
		},
		{
			FieldName: "AvailableContainerRuntimes",
			Title:     "Available Container Runtimes",
		},
		{
			FieldName: "AvailableFeatureGates",
			Title:     "Available Feature Gates",
		},
		{
			FieldName: "AvailableAdmissionPlugins",
			Title:     "Available Admission Plugins",
		},
	}

	str, err := human.Marshal(version, opt)
	if err != nil {
		return "", err
	}

	return str, nil
}

func getLatestK8SVersion(scwClient *scw.Client) (string, error) {
	api := k8s.NewAPI(scwClient)
	versions, err := api.ListVersions(&k8s.ListVersionsRequest{})
	if err != nil {
		return "", fmt.Errorf("could not get latest K8S version: %s", err)
	}

	latestVersion, _ := version.NewVersion("0.0.0")
	for _, v := range versions.Versions {
		newVersion, _ := version.NewVersion(v.Name)
		if newVersion.GreaterThan(latestVersion) {
			latestVersion = newVersion
		}
	}

	return latestVersion.String(), nil
}
