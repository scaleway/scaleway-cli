package mongodb

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/interactive"
	"github.com/scaleway/scaleway-cli/v2/internal/passwordgenerator"
	mongodb "github.com/scaleway/scaleway-sdk-go/api/mongodb/v1"
)

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

		generated := customRequest.GeneratePassword && req.Password == ""
		var err error
		req.Password, err = passwordgenerator.Resolve(customRequest.GeneratePassword, req.Password)
		if err != nil {
			return nil, err
		}
		if generated {
			_, _ = interactive.Println("Your generated password is", req.Password)
		}

		return api.CreateUser(req)
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

		generated := customRequest.GeneratePassword && (req.Password == nil || *req.Password == "")
		password, err := passwordgenerator.ResolvePtr(customRequest.GeneratePassword, req.Password)
		if err != nil {
			return nil, err
		}
		req.Password = password
		if generated && password != nil {
			_, _ = interactive.Println("Your generated password is", *password)
		}

		return api.UpdateUser(req)
	}

	return c
}
