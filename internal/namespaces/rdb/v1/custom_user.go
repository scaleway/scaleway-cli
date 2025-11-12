package rdb

import (
	"context"
	"errors"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/interactive"
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
		Password string `json:"password"`
	}

	c.ArgSpecs.AddBefore("password", &core.ArgSpec{
		Name:       "generate-password",
		Short:      `Will generate a 21 character-length password that contains a mix of upper/lower case letters, numbers and special symbols`,
		Required:   false,
		Deprecated: false,
		Positional: false,
		Default:    core.DefaultValueSetter("true"),
	})
	c.ArgsType = reflect.TypeOf(rdbCreateUserRequestCustom{})

	c.Run = func(ctx context.Context, argsI any) (any, error) {
		api := rdb.NewAPI(core.ExtractClient(ctx))
		customRequest := argsI.(*rdbCreateUserRequestCustom)
		req := customRequest.CreateUserRequest

		if req.Name != "" {
			name := req.Name
			users, err := api.ListUsers(&rdb.ListUsersRequest{
				Region:     req.Region,
				InstanceID: req.InstanceID,
				Name:       &name,
			}, scw.WithAllPages())
			if err == nil && users.TotalCount > 0 {
				return rdbCreateUserResponseCustom{
					User:     users.Users[0],
					Password: "",
				}, nil
			}
		}

		if customRequest.GeneratePassword && customRequest.Password == "" {
			password, err := passwordgenerator.GeneratePassword(21, 1, 1, 1, 1)
			if err != nil {
				return nil, err
			}
			req.Password = password
			// password is returned in the command output; avoid logging it to stdout
		}

		user, err := api.CreateUser(req)
		if err != nil {
			return nil, err
		}

		result := rdbCreateUserResponseCustom{
			User:     user,
			Password: req.Password,
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
		Password string `json:"password"`
	}

	c.ArgSpecs.AddBefore("password", &core.ArgSpec{
		Name:       "generate-password",
		Short:      `Will generate a 21 character-length password that contains a mix of upper/lower case letters, numbers and special symbols`,
		Required:   false,
		Deprecated: false,
		Positional: false,
		Default:    core.DefaultValueSetter("true"),
	})
	c.ArgsType = reflect.TypeOf(rdbUpdateUserRequestCustom{})

	c.Run = func(ctx context.Context, argsI any) (any, error) {
		client := core.ExtractClient(ctx)
		api := rdb.NewAPI(client)

		customRequest := argsI.(*rdbUpdateUserRequestCustom)
		updateUserRequest := customRequest.UpdateUserRequest

		var err error

		if customRequest.GeneratePassword || customRequest.Password != nil {
			switch {
			case customRequest.GeneratePassword && customRequest.Password == nil:
				updateUserRequest.Password = new(string)
				pwd, err := passwordgenerator.GeneratePassword(21, 1, 1, 1, 1)
				if err != nil {
					return nil, err
				}
				*updateUserRequest.Password = pwd

				_, err = interactive.Println("Your generated password is", pwd)
				if err != nil {
					return nil, err
				}

			case !customRequest.GeneratePassword && customRequest.Password == nil:
				return nil, errors.New(
					"you must provide a password when generate-password is set to false",
				)

			default:
				updateUserRequest.Password = customRequest.Password
			}
		}

		user, err := api.UpdateUser(updateUserRequest)
		if err != nil {
			return nil, err
		}

		respPwd := ""
		if updateUserRequest.Password != nil {
			respPwd = *updateUserRequest.Password
		}

		result := rdbUpdateUserResponseCustom{
			User:     user,
			Password: respPwd,
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
		ArgsType:  reflect.TypeOf(rdbGetURLArgs{}),
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
