package mongodb

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/connect"
	mongodb "github.com/scaleway/scaleway-sdk-go/api/mongodb/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func userUpdateBuilder(c *core.Command) *core.Command {
	type mongodbUpdateUserRequestCustom struct {
		*mongodb.UpdateUserRequest
		CliConnect bool
	}

	c.ArgSpecs.AddBefore("password", connect.CliConnectArgSpec())
	c.ArgsType = reflect.TypeFor[mongodbUpdateUserRequestCustom]()

	c.Run = func(ctx context.Context, args any) (any, error) {
		customRequest := args.(*mongodbUpdateUserRequestCustom)
		request := customRequest.UpdateUserRequest

		if customRequest.CliConnect && request.Password == nil {
			password, err := connect.PromptPassword(ctx)
			if err != nil {
				return nil, err
			}
			request.Password = &password
		}

		api := mongodb.NewAPI(core.ExtractClient(ctx))
		user, err := api.UpdateUser(request, scw.WithContext(ctx))
		if err != nil {
			return nil, err
		}

		if customRequest.CliConnect {
			if err := connect.PersistCredentials(
				ctx,
				request.Region,
				"mongodb",
				string(request.Region),
				request.InstanceID,
				request.Name,
				*request.Password,
				"instance user",
			); err != nil {
				return nil, err
			}
		}

		return user, nil
	}

	return c
}
