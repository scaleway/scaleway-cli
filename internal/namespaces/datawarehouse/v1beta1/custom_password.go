package datawarehouse

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/passwordgenerator"
	datawarehouse "github.com/scaleway/scaleway-sdk-go/api/datawarehouse/v1beta1"
)

type CreateDeploymentResult struct {
	*datawarehouse.Deployment
	Password string `json:"password,omitempty"`
}

type UserResult struct {
	*datawarehouse.User
	Password string `json:"password,omitempty"`
}

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

		var err error
		req.Password, err = passwordgenerator.Resolve(customRequest.GeneratePassword, req.Password)
		if err != nil {
			return nil, err
		}

		deployment, err := api.CreateDeployment(req)
		if err != nil {
			return nil, err
		}

		return CreateDeploymentResult{Deployment: deployment, Password: req.Password}, nil
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
