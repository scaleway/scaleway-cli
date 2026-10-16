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

	c.SeeAlsos = connectSeeAlsos

	c.Run = func(ctx context.Context, args any) (any, error) {
		req := args.(*datawarehouseUpdateUserRequestCustom)

		return connect.RunUserUpdate(
			ctx,
			req.CliConnect,
			connect.Datawarehouse,
			req.Region,
			req.Region.String(),
			req.DeploymentID,
			req.Name,
			&req.Password,
			func(ctx context.Context) (any, error) {
				return datawarehouse.NewAPI(core.ExtractClient(ctx)).
					UpdateUser(req.UpdateUserRequest, scw.WithContext(ctx))
			},
		)
	}

	return c
}
