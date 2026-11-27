package datawarehouse

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/passwordgenerator"
	datawarehouse "github.com/scaleway/scaleway-sdk-go/api/datawarehouse/v1beta1"
	"github.com/scaleway/scaleway-sdk-go/scw"
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

		generated := customRequest.GeneratePassword && req.Password == ""
		var err error
		req.Password, err = passwordgenerator.Resolve(customRequest.GeneratePassword, req.Password)
		if err != nil {
			return nil, err
		}

		deployment, err := api.CreateDeployment(req, scw.WithContext(ctx))
		if err != nil {
			return nil, err
		}

		result := CreateDeploymentResult{Deployment: deployment}
		if generated {
			result.Password = req.Password
		}

		return result, nil
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

		user, err := api.CreateUser(req, scw.WithContext(ctx))
		if err != nil {
			return nil, err
		}

		result := UserResult{User: user}
		if generated {
			result.Password = req.Password
		}

		return result, nil
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

		user, err := api.UpdateUser(req, scw.WithContext(ctx))
		if err != nil {
			return nil, err
		}

		result := UserResult{User: user}
		if generated {
			result.Password = *password
		}

		return result, nil
	}

	return c
}
