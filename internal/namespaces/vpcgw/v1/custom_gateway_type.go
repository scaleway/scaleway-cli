package vpcgw

import (
	"context"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-sdk-go/api/vpcgw/v1"
)

func vpcgwGatewayTypeListBuilder(c *core.Command) *core.Command {
	c.AddInterceptors(
		func(ctx context.Context, argsI any, runner core.CommandRunner) (any, error) {
			res, err := runner(ctx, argsI)
			if err != nil {
				return nil, err
			}

			typesResponse := res.(*vpcgw.ListGatewayTypesResponse)

			return typesResponse.Types, nil
		},
	)

	return c
}
