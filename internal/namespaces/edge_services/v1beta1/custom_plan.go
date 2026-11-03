package edge_services

import (
	"context"

	"github.com/scaleway/scaleway-cli/v2/core"
	edge_services "github.com/scaleway/scaleway-sdk-go/api/edge_services/v1beta1"
)

// The generated command returns the whole ListPlansResponse, which renders as a
// flat key/value dump. Return the plans slice so the default table rendering
// applies, as in the other list commands of this namespace.
func planListBuilder(c *core.Command) *core.Command {
	c.Run = func(ctx context.Context, args any) (i any, e error) {
		client := core.ExtractClient(ctx)
		api := edge_services.NewAPI(client)

		resp, err := api.ListPlans()
		if err != nil {
			return nil, err
		}

		return resp.Plans, nil
	}

	return c
}
