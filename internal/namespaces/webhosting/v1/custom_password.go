package webhosting

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/interactive"
	"github.com/scaleway/scaleway-cli/v2/internal/passwordgenerator"
	"github.com/scaleway/scaleway-sdk-go/api/webhosting/v1"
)

func databaseUserCreateBuilder(c *core.Command) *core.Command {
	type requestCustom struct {
		*webhosting.DatabaseAPICreateDatabaseUserRequest
		GeneratePassword bool
	}

	c.ArgSpecs.AddBefore("password", passwordgenerator.ArgSpec(true))
	c.ArgSpecs.GetByName("password").Required = false
	c.ArgsType = reflect.TypeFor[requestCustom]()

	c.Run = func(ctx context.Context, argsI any) (any, error) {
		api := webhosting.NewDatabaseAPI(core.ExtractClient(ctx))
		customRequest := argsI.(*requestCustom)
		req := customRequest.DatabaseAPICreateDatabaseUserRequest

		generated := customRequest.GeneratePassword && req.Password == ""
		var err error
		req.Password, err = passwordgenerator.Resolve(customRequest.GeneratePassword, req.Password)
		if err != nil {
			return nil, err
		}
		if generated {
			_, _ = interactive.Println("Your generated password is", req.Password)
		}

		return api.CreateDatabaseUser(req)
	}

	return c
}

func databaseUserChangePasswordBuilder(c *core.Command) *core.Command {
	type requestCustom struct {
		*webhosting.DatabaseAPIChangeDatabaseUserPasswordRequest
		GeneratePassword bool
	}

	c.ArgSpecs.AddBefore("password", passwordgenerator.ArgSpec(true))
	c.ArgSpecs.GetByName("password").Required = false
	c.ArgsType = reflect.TypeFor[requestCustom]()

	c.Run = func(ctx context.Context, argsI any) (any, error) {
		api := webhosting.NewDatabaseAPI(core.ExtractClient(ctx))
		customRequest := argsI.(*requestCustom)
		req := customRequest.DatabaseAPIChangeDatabaseUserPasswordRequest

		generated := customRequest.GeneratePassword && req.Password == ""
		var err error
		req.Password, err = passwordgenerator.Resolve(customRequest.GeneratePassword, req.Password)
		if err != nil {
			return nil, err
		}
		if generated {
			_, _ = interactive.Println("Your generated password is", req.Password)
		}

		return api.ChangeDatabaseUserPassword(req)
	}

	return c
}

func ftpAccountCreateBuilder(c *core.Command) *core.Command {
	type requestCustom struct {
		*webhosting.FtpAccountAPICreateFtpAccountRequest
		GeneratePassword bool
	}

	c.ArgSpecs.AddBefore("password", passwordgenerator.ArgSpec(true))
	c.ArgSpecs.GetByName("password").Required = false
	c.ArgsType = reflect.TypeFor[requestCustom]()

	c.Run = func(ctx context.Context, argsI any) (any, error) {
		api := webhosting.NewFtpAccountAPI(core.ExtractClient(ctx))
		customRequest := argsI.(*requestCustom)
		req := customRequest.FtpAccountAPICreateFtpAccountRequest

		generated := customRequest.GeneratePassword && req.Password == ""
		var err error
		req.Password, err = passwordgenerator.Resolve(customRequest.GeneratePassword, req.Password)
		if err != nil {
			return nil, err
		}
		if generated {
			_, _ = interactive.Println("Your generated password is", req.Password)
		}

		return api.CreateFtpAccount(req)
	}

	return c
}

func mailAccountCreateBuilder(c *core.Command) *core.Command {
	type requestCustom struct {
		*webhosting.MailAccountAPICreateMailAccountRequest
		GeneratePassword bool
	}

	c.ArgSpecs.AddBefore("password", passwordgenerator.ArgSpec(true))
	c.ArgSpecs.GetByName("password").Required = false
	c.ArgsType = reflect.TypeFor[requestCustom]()

	c.Run = func(ctx context.Context, argsI any) (any, error) {
		api := webhosting.NewMailAccountAPI(core.ExtractClient(ctx))
		customRequest := argsI.(*requestCustom)
		req := customRequest.MailAccountAPICreateMailAccountRequest

		generated := customRequest.GeneratePassword && req.Password == ""
		var err error
		req.Password, err = passwordgenerator.Resolve(customRequest.GeneratePassword, req.Password)
		if err != nil {
			return nil, err
		}
		if generated {
			_, _ = interactive.Println("Your generated password is", req.Password)
		}

		return api.CreateMailAccount(req)
	}

	return c
}

func mailAccountChangePasswordBuilder(c *core.Command) *core.Command {
	type requestCustom struct {
		*webhosting.MailAccountAPIChangeMailAccountPasswordRequest
		GeneratePassword bool
	}

	c.ArgSpecs.AddBefore("password", passwordgenerator.ArgSpec(true))
	c.ArgSpecs.GetByName("password").Required = false
	c.ArgsType = reflect.TypeFor[requestCustom]()

	c.Run = func(ctx context.Context, argsI any) (any, error) {
		api := webhosting.NewMailAccountAPI(core.ExtractClient(ctx))
		customRequest := argsI.(*requestCustom)
		req := customRequest.MailAccountAPIChangeMailAccountPasswordRequest

		generated := customRequest.GeneratePassword && req.Password == ""
		var err error
		req.Password, err = passwordgenerator.Resolve(customRequest.GeneratePassword, req.Password)
		if err != nil {
			return nil, err
		}
		if generated {
			_, _ = interactive.Println("Your generated password is", req.Password)
		}

		return api.ChangeMailAccountPassword(req)
	}

	return c
}
