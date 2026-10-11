package datawarehouse

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/connect"
	datawarehouse "github.com/scaleway/scaleway-sdk-go/api/datawarehouse/v1beta1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func userUpdateBuilder(c *core.Command) *core.Command {
	type datawarehouseUpdateUserRequestCustom struct {
		*datawarehouse.UpdateUserRequest
		CliConnect bool
	}

	c.ArgSpecs.AddBefore("password", connect.CliConnectArgSpec())
	c.ArgsType = reflect.TypeFor[datawarehouseUpdateUserRequestCustom]()

	c.Run = func(ctx context.Context, args any) (any, error) {
		customRequest := args.(*datawarehouseUpdateUserRequestCustom)
		request := customRequest.UpdateUserRequest

		if customRequest.CliConnect && request.Password == nil {
			password, err := connect.PromptPassword(ctx)
			if err != nil {
				return nil, err
			}
			request.Password = &password
		}

		api := datawarehouse.NewAPI(core.ExtractClient(ctx))
		user, err := api.UpdateUser(request, scw.WithContext(ctx))
		if err != nil {
			return nil, err
		}

		if customRequest.CliConnect {
			if err := connect.PersistCredentials(
				ctx,
				request.Region,
				"datawarehouse",
				string(request.Region),
				request.DeploymentID,
				request.Name,
				*request.Password,
				"deployment user",
			); err != nil {
				return nil, err
			}
		}

		return user, nil
	}

	return c
}
