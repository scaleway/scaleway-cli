package instance

import (
	"context"
	"errors"
	"fmt"
	"sort"

	"github.com/scaleway/scaleway-cli/v2/core"
	instance "github.com/scaleway/scaleway-sdk-go/api/instance/v2alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func userDataSetBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("content").Required = true
	c.ArgSpecs.GetByName("content").CanLoadFile = true
	c.ArgSpecs.GetByName("server-id").Positional = true

	return c
}

func userDataGetBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("server-id").Positional = true

	c.AddInterceptors(
		func(ctx context.Context, argsI any, runner core.CommandRunner) (any, error) {
			req := argsI.(*instance.GetUserDataRequest)
			res, err := runner(ctx, argsI)
			if err != nil {
				if _, ok := errors.AsType[*scw.ResourceNotFoundError](err); ok {
					return nil, fmt.Errorf("'%s' key does not exist in 'user_data'", req.Key)
				}

				return nil, err
			}

			return res, nil
		},
	)

	return c
}

func userDataListBuilder(c *core.Command) *core.Command {
	type userDataRow struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}

	c.ArgSpecs.GetByName("server-id").Positional = true

	c.AddInterceptors(
		func(ctx context.Context, argsI any, _ core.CommandRunner) (any, error) {
			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			args := argsI.(*instance.ListUserDataKeysRequest)

			res, err := api.ListUserDataKeys(&instance.ListUserDataKeysRequest{
				Zone:     args.Zone,
				ServerID: args.ServerID,
			}, scw.WithContext(ctx), scw.WithAllPages())
			if err != nil {
				return nil, fmt.Errorf("failed to list user_data keys: %w", err)
			}

			var udRows []userDataRow
			for _, key := range res.Keys {
				ud, err := api.GetUserData(&instance.GetUserDataRequest{
					Zone:     args.Zone,
					ServerID: args.ServerID,
					Key:      key,
				}, scw.WithContext(ctx))
				if err != nil {
					return nil, fmt.Errorf("failed to get user_data for key '%s': %w", key, err)
				}

				udRows = append(udRows, userDataRow{
					Key:   ud.Key,
					Value: string(ud.Content),
				})
			}

			sort.Slice(udRows, func(i, j int) bool {
				return udRows[i].Key < udRows[j].Key
			})

			return udRows, nil
		},
	)

	return c
}

func userDataDeleteBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("server-id").Positional = true

	c.AddInterceptors(
		func(ctx context.Context, argsI any, runner core.CommandRunner) (any, error) {
			req := argsI.(*instance.DeleteUserDataRequest)
			res, err := runner(ctx, argsI)
			if err != nil {
				return res, err
			}

			return &core.SuccessResult{
				Resource: fmt.Sprintf("user-data key '%s'", req.Key),
				Verb:     "delete",
			}, nil
		},
	)

	return c
}

func userDataSetCloudInitBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("content").Required = true
	c.ArgSpecs.GetByName("content").CanLoadFile = true
	c.ArgSpecs.GetByName("server-id").Positional = true

	return c
}

func userDataGetCloudInitBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("server-id").Positional = true

	return c
}
