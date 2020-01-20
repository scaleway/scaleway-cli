package instance

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
)

//
// Commands
//

func userDataCommand() *core.Command {
	return &core.Command{
		Namespace: "instance",
		Resource:  "user-data",
	}
}

func userDataListCommand() *core.Command {
	return &core.Command{
		Short:     `List user data`,
		Long:      `List user data for the given server.`,
		Namespace: "instance",
		Resource:  "user-data",
		Verb:      "list",
		ArgsType:  reflect.TypeOf(instance.ListServerUserDataRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ZoneArgSpec(),
			{
				Name:     "server-id",
				Short:    `ID of a server`,
				Required: true,
			},
		},
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
			return instance.NewAPI(core.ExtractClient(ctx)).ListServerUserData(argsI.(*instance.ListServerUserDataRequest))
		},
	}
}

func userDataDeleteCommand() *core.Command {
	return &core.Command{
		Short:     `Delete user data by key`,
		Long:      `Delete user data key for the given server.`,
		Namespace: "instance",
		Resource:  "user-data",
		Verb:      "delete",
		ArgsType:  reflect.TypeOf(instance.DeleteServerUserDataRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ZoneArgSpec(),
			{
				Name:     "server-id",
				Short:    `ID of a server`,
				Required: true,
			},
			{
				Name:     "key",
				Short:    `Key of the user data`,
				Required: true,
			},
		},
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
			err := instance.NewAPI(core.ExtractClient(ctx)).DeleteServerUserData(argsI.(*instance.DeleteServerUserDataRequest))
			if err != nil {
				return nil, err
			}
			return &core.SuccessResult{}, nil
		},
	}
}

func userDataGetCommand() *core.Command {
	return &core.Command{
		Short:     `Get user data key`,
		Long:      `Get user data key for the given server.`,
		Namespace: "instance",
		Resource:  "user-data",
		Verb:      "get",
		ArgsType:  reflect.TypeOf(instance.GetServerUserDataRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ZoneArgSpec(),
			{
				Name:     "server-id",
				Short:    `ID of a server`,
				Required: true,
			},
			{
				Name:     "key",
				Short:    `Key of the user data`,
				Required: true,
			},
		},
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
			return instance.NewAPI(core.ExtractClient(ctx)).GetServerUserData(argsI.(*instance.GetServerUserDataRequest))
		},
	}
}

func userDataSetCommand() *core.Command {
	return &core.Command{
		Short:     `Set a user data`,
		Long:      `Set a user data for the given server.`,
		Namespace: "instance",
		Resource:  "user-data",
		Verb:      "set",
		ArgsType:  reflect.TypeOf(instance.SetServerUserDataRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ZoneArgSpec(),
			{
				Name:     "server-id",
				Short:    `ID of a server`,
				Required: true,
			},
			{
				Name:     "key",
				Short:    `Key of the user data`,
				Required: true,
			},
			{
				Name:     "content",
				Short:    `Content of the user data`,
				Required: true,
			},
		},
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
			err := instance.NewAPI(core.ExtractClient(ctx)).SetServerUserData(argsI.(*instance.SetServerUserDataRequest))
			if err != nil {
				return nil, err
			}
			return &core.SuccessResult{}, nil
		},
	}
}
