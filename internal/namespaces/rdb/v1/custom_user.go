package rdb

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/passwordgenerator"
	"github.com/scaleway/scaleway-sdk-go/api/rdb/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// userListBuilder creates a table visualization of user's permission across different database in a given RDB instance
func userListBuilder(c *core.Command) *core.Command {
	c.Interceptor = func(ctx context.Context, argsI any, runner core.CommandRunner) (any, error) {
		type customUser struct {
			Name      string   `json:"name"`
			IsAdmin   bool     `json:"is_admin"`
			ReadOnly  []string `json:"readonly"`
			ReadWrite []string `json:"readwrite"`
			All       []string `json:"all"`
			Custom    []string `json:"custom"`
			None      []string `json:"none"`
		}

		resI, err := runner(ctx, argsI)
		if err != nil {
			return nil, err
		}

		// We index user by their name and use customUser as the type holding the different privileges across databases
		index := make(map[string]*customUser)
		res := []*customUser(nil)
		listUserRequest := argsI.(*rdb.ListUsersRequest)
		listUserResponse := resI.([]*rdb.User)
		for _, user := range listUserResponse {
			user := &customUser{
				Name:      user.Name,
				IsAdmin:   user.IsAdmin,
				ReadOnly:  []string{},
				ReadWrite: []string{},
				All:       []string{},
				Custom:    []string{},
				None:      []string{},
			}
			res = append(res, user)
			index[user.Name] = user
		}

		api := rdb.NewAPI(core.ExtractClient(ctx))
		listPrivileges, err := api.ListPrivileges(
			&rdb.ListPrivilegesRequest{
				InstanceID: listUserRequest.InstanceID,
				Region:     listUserRequest.Region,
			},
			scw.WithAllPages(),
		)
		if err != nil {
			return resI, err
		}

		for _, privilege := range listPrivileges.Privileges {
			user, userExists := index[privilege.UserName]
			if !userExists {
				continue
			}

			switch privilege.Permission {
			case rdb.PermissionAll:
				user.All = append(user.All, privilege.DatabaseName)
			case rdb.PermissionReadonly:
				user.ReadOnly = append(user.ReadOnly, privilege.DatabaseName)
			case rdb.PermissionCustom:
				user.Custom = append(user.Custom, privilege.DatabaseName)
			case rdb.PermissionNone:
				user.None = append(user.None, privilege.DatabaseName)
			case rdb.PermissionReadwrite:
				user.ReadWrite = append(user.ReadWrite, privilege.DatabaseName)
			default:
				core.ExtractLogger(ctx).
					Errorf("unsupported permission value %s", privilege.Permission)
			}
		}

		return res, nil
	}

	return c
}

func userCreateBuilder(c *core.Command) *core.Command {
	type rdbCreateUserRequestCustom struct {
		*rdb.CreateUserRequest
		GeneratePassword bool
	}

	type rdbCreateUserResponseCustom struct {
		*rdb.User
		Password string `json:"password,omitempty"`
	}

	c.ArgSpecs.AddBefore("password", passwordgenerator.ArgSpec(true))
	c.ArgsType = reflect.TypeFor[rdbCreateUserRequestCustom]()

	c.Run = func(ctx context.Context, argsI any) (any, error) {
		api := rdb.NewAPI(core.ExtractClient(ctx))
		customRequest := argsI.(*rdbCreateUserRequestCustom)
		req := customRequest.CreateUserRequest

		if req.Name != "" {
			users, err := api.ListUsers(&rdb.ListUsersRequest{
				Region:     req.Region,
				InstanceID: req.InstanceID,
				Name:       new(req.Name),
			}, scw.WithAllPages(), scw.WithContext(ctx))
			if err == nil && users.TotalCount > 0 {
				return rdbCreateUserResponseCustom{User: users.Users[0]}, nil
			}
		}

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

		result := rdbCreateUserResponseCustom{User: user}
		if generated {
			result.Password = req.Password
		}

		return result, nil
	}

	return c
}

func userUpdateBuilder(c *core.Command) *core.Command {
	type rdbUpdateUserRequestCustom struct {
		*rdb.UpdateUserRequest
		GeneratePassword bool
	}

	type rdbUpdateUserResponseCustom struct {
		*rdb.User
		Password string `json:"password,omitempty"`
	}

	c.ArgSpecs.AddBefore("password", passwordgenerator.ArgSpec(true))
	c.ArgsType = reflect.TypeFor[rdbUpdateUserRequestCustom]()

	c.Run = func(ctx context.Context, argsI any) (any, error) {
		api := rdb.NewAPI(core.ExtractClient(ctx))
		customRequest := argsI.(*rdbUpdateUserRequestCustom)
		updateUserRequest := customRequest.UpdateUserRequest

		generated := customRequest.GeneratePassword && customRequest.Password == nil
		if generated {
			pwd, err := passwordgenerator.Generate()
			if err != nil {
				return nil, err
			}
			updateUserRequest.Password = &pwd
		}

		user, err := api.UpdateUser(updateUserRequest, scw.WithContext(ctx))
		if err != nil {
			return nil, err
		}

		result := rdbUpdateUserResponseCustom{User: user}
		if generated {
			result.Password = *updateUserRequest.Password
		}

		return result, nil
	}

	return c
}

func userGetURLCommand() *core.Command {
	return &core.Command{
		Namespace: "rdb",
		Resource:  "user",
		Verb:      "get-url",
		Short:     "Gets the URL to connect to the Database",
		Long:      "Provides the URL to connect to a Database on an Instance as the given user",
		ArgsType:  reflect.TypeFor[rdbGetURLArgs](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "instance-id",
				Short:      `ID of the Database Instance`,
				Required:   true,
				Positional: true,
			},
			{
				Name:       "user",
				Short:      `User of the Database`,
				Required:   false,
				Positional: false,
			},
			{
				Name:       "db",
				Short:      `Name of the Database to connect to`,
				Required:   false,
				Positional: false,
			},
		},
		Run: generateURL,
	}
}
