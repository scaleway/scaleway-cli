package k8s

import (
	"context"

	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-cli/internal/human"
	k8s "github.com/scaleway/scaleway-sdk-go/api/k8s/v1"
)

func versionListBuilder(c *core.Command) *core.Command {
	c.AddInterceptors(func(ctx context.Context, argsI interface{}, runner core.CommandRunner) (interface{}, error) {
		originalRes, err := runner(ctx, argsI)
		if err != nil {
			return nil, err
		}

		versionsResponse := originalRes.(*k8s.ListVersionsResponse)
		return versionsResponse.Versions, nil
	})

	return c
}

func versionMarshalerFunc(i interface{}, opt *human.MarshalOpt) (string, error) {
	type tmp k8s.Version
	version := tmp(i.(k8s.Version))

	// Sections
	opt.Sections = []*human.MarshalSection{
		{
			FieldName: "AvailableCnis",
			Title:     "Available CNIs",
		},
		{
			FieldName: "AvailableIngresses",
			Title:     "Available Ingresses",
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
