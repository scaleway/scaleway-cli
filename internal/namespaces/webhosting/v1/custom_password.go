package webhosting

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/passwordgenerator"
	"github.com/scaleway/scaleway-sdk-go/api/webhosting/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func databaseUserCreateBuilder(c *core.Command) *core.Command {
	type requestCustom struct {
		*webhosting.DatabaseAPICreateDatabaseUserRequest
		GeneratePassword bool
	}
	type responseCustom struct {
		*webhosting.DatabaseUser
		Password string `json:"password,omitempty"`
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

		user, err := api.CreateDatabaseUser(req, scw.WithContext(ctx))
		if err != nil {
			return nil, err
		}

		result := responseCustom{DatabaseUser: user}
		if generated {
			result.Password = req.Password
		}

		return result, nil
	}

	return c
}

func databaseUserChangePasswordBuilder(c *core.Command) *core.Command {
	type requestCustom struct {
		*webhosting.DatabaseAPIChangeDatabaseUserPasswordRequest
		GeneratePassword bool
	}
	type responseCustom struct {
		*webhosting.DatabaseUser
		Password string `json:"password,omitempty"`
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

		user, err := api.ChangeDatabaseUserPassword(req, scw.WithContext(ctx))
		if err != nil {
			return nil, err
		}

		result := responseCustom{DatabaseUser: user}
		if generated {
			result.Password = req.Password
		}

		return result, nil
	}

	return c
}

func ftpAccountCreateBuilder(c *core.Command) *core.Command {
	type requestCustom struct {
		*webhosting.FtpAccountAPICreateFtpAccountRequest
		GeneratePassword bool
	}
	type responseCustom struct {
		*webhosting.FtpAccount
		Password string `json:"password,omitempty"`
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

		account, err := api.CreateFtpAccount(req, scw.WithContext(ctx))
		if err != nil {
			return nil, err
		}

		result := responseCustom{FtpAccount: account}
		if generated {
			result.Password = req.Password
		}

		return result, nil
	}

	return c
}

func mailAccountCreateBuilder(c *core.Command) *core.Command {
	type requestCustom struct {
		*webhosting.MailAccountAPICreateMailAccountRequest
		GeneratePassword bool
	}
	type responseCustom struct {
		*webhosting.MailAccount
		Password string `json:"password,omitempty"`
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

		account, err := api.CreateMailAccount(req, scw.WithContext(ctx))
		if err != nil {
			return nil, err
		}

		result := responseCustom{MailAccount: account}
		if generated {
			result.Password = req.Password
		}

		return result, nil
	}

	return c
}

func mailAccountChangePasswordBuilder(c *core.Command) *core.Command {
	type requestCustom struct {
		*webhosting.MailAccountAPIChangeMailAccountPasswordRequest
		GeneratePassword bool
	}
	type responseCustom struct {
		*webhosting.MailAccount
		Password string `json:"password,omitempty"`
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

		account, err := api.ChangeMailAccountPassword(req, scw.WithContext(ctx))
		if err != nil {
			return nil, err
		}

		result := responseCustom{MailAccount: account}
		if generated {
			result.Password = req.Password
		}

		return result, nil
	}

	return c
}
