package mongodb

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/passwordgenerator"
	mongodb "github.com/scaleway/scaleway-sdk-go/api/mongodb/v1"
)

type UserResult struct {
	*mongodb.User
	Password string `json:"password,omitempty"`
}

func userCreateBuilder(c *core.Command) *core.Command {
	type mongodbCreateUserRequestCustom struct {
		*mongodb.CreateUserRequest
		GeneratePassword bool
	}

	c.ArgSpecs.AddBefore("password", passwordgenerator.ArgSpec(true))
	c.ArgSpecs.GetByName("password").Required = false
	c.ArgsType = reflect.TypeFor[mongodbCreateUserRequestCustom]()

	c.Run = func(ctx context.Context, argsI any) (any, error) {
		api := mongodb.NewAPI(core.ExtractClient(ctx))
		customRequest := argsI.(*mongodbCreateUserRequestCustom)
		req := customRequest.CreateUserRequest

		var err error
		req.Password, err = passwordgenerator.Resolve(customRequest.GeneratePassword, req.Password)
		if err != nil {
			return nil, err
		}

		user, err := api.CreateUser(req)
		if err != nil {
			return nil, err
		}

		return UserResult{User: user, Password: req.Password}, nil
	}

	return c
}

func userUpdateBuilder(c *core.Command) *core.Command {
	type mongodbUpdateUserRequestCustom struct {
		*mongodb.UpdateUserRequest
		GeneratePassword bool
	}

	c.ArgSpecs.AddBefore("password", passwordgenerator.ArgSpec(false))
	c.ArgsType = reflect.TypeFor[mongodbUpdateUserRequestCustom]()

	c.Run = func(ctx context.Context, argsI any) (any, error) {
		api := mongodb.NewAPI(core.ExtractClient(ctx))
		customRequest := argsI.(*mongodbUpdateUserRequestCustom)
		req := customRequest.UpdateUserRequest

		password, err := passwordgenerator.ResolvePtr(customRequest.GeneratePassword, req.Password)
		if err != nil {
			return nil, err
		}
		req.Password = password

		user, err := api.UpdateUser(req)
		if err != nil {
			return nil, err
		}

		result := UserResult{User: user}
		if password != nil {
			result.Password = *password
		}

		return result, nil
	}

	return c
}
