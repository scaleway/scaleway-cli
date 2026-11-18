package datawarehouse

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/interactive"
	"github.com/scaleway/scaleway-cli/v2/internal/passwordgenerator"
	datawarehouse "github.com/scaleway/scaleway-sdk-go/api/datawarehouse/v1beta1"
)

func deploymentCreateBuilder(c *core.Command) *core.Command {
	type createDeploymentRequestCustom struct {
		*datawarehouse.CreateDeploymentRequest
		GeneratePassword bool
	}

	c.ArgSpecs.AddBefore("password", passwordgenerator.ArgSpec(true))
	c.ArgSpecs.GetByName("password").Required = false
	c.ArgsType = reflect.TypeFor[createDeploymentRequestCustom]()

	c.Run = func(ctx context.Context, argsI any) (any, error) {
		api := datawarehouse.NewAPI(core.ExtractClient(ctx))
		customRequest := argsI.(*createDeploymentRequestCustom)
		req := customRequest.CreateDeploymentRequest

		generated := customRequest.GeneratePassword && req.Password == ""
		var err error
		req.Password, err = passwordgenerator.Resolve(customRequest.GeneratePassword, req.Password)
		if err != nil {
			return nil, err
		}
		if generated {
			_, _ = interactive.Println("Your generated password is", req.Password)
		}

		return api.CreateDeployment(req)
	}

	return c
}

func userCreateBuilder(c *core.Command) *core.Command {
	type createUserRequestCustom struct {
		*datawarehouse.CreateUserRequest
		GeneratePassword bool
	}

	c.ArgSpecs.AddBefore("password", passwordgenerator.ArgSpec(true))
	c.ArgSpecs.GetByName("password").Required = false
	c.ArgsType = reflect.TypeFor[createUserRequestCustom]()

	c.Run = func(ctx context.Context, argsI any) (any, error) {
		api := datawarehouse.NewAPI(core.ExtractClient(ctx))
		customRequest := argsI.(*createUserRequestCustom)
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
	type updateUserRequestCustom struct {
		*datawarehouse.UpdateUserRequest
		GeneratePassword bool
	}

	c.ArgSpecs.AddBefore("password", passwordgenerator.ArgSpec(false))
	c.ArgsType = reflect.TypeFor[updateUserRequestCustom]()

	c.Run = func(ctx context.Context, argsI any) (any, error) {
		api := datawarehouse.NewAPI(core.ExtractClient(ctx))
		customRequest := argsI.(*updateUserRequestCustom)
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
