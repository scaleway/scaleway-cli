package k8s

import (
	"context"

	"github.com/scaleway/scaleway-cli/internal/core"
	k8s "github.com/scaleway/scaleway-sdk-go/api/k8s/v1beta4"
)

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
