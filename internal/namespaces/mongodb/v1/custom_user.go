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

	c.SeeAlsos = connectSeeAlsos

	c.Run = func(ctx context.Context, args any) (any, error) {
		req := args.(*mongodbUpdateUserRequestCustom)

		return connect.RunUserUpdate(
			ctx,
			req.CliConnect,
			"mongodb",
			"instance user",
			req.Region,
			req.InstanceID,
			req.Name,
			&req.Password,
			func(ctx context.Context) (any, error) {
				return mongodb.NewAPI(core.ExtractClient(ctx)).
					UpdateUser(req.UpdateUserRequest, scw.WithContext(ctx))
			},
		)
	}

	return c
}
