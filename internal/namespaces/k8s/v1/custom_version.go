package k8s

import (
	"context"

	"github.com/scaleway/scaleway-cli/internal/core"
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
