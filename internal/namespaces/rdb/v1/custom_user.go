package rdb

import (
	"context"
	"fmt"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/passwordgenerator"
	"github.com/scaleway/scaleway-sdk-go/api/rdb/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// userListBuilder creates a table visualization of user's permission across different database in a given RDB instance
func userListBuilder(c *core.Command) *core.Command {
	c.Interceptor = func(ctx context.Context, argsI interface{}, runner core.CommandRunner) (interface{}, error) {
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

	c.Run = func(ctx context.Context, argsI interface{}) (interface{}, error) {
		client := core.ExtractClient(ctx)
		api := rdb.NewAPI(client)

		customRequest := argsI.(*rdbCreateUserRequestCustom)
		createUserRequest := customRequest.CreateUserRequest

		var err error
		if customRequest.GeneratePassword && customRequest.Password == "" {
			createUserRequest.Password, err = passwordgenerator.GeneratePassword(21, 1, 1, 1, 1)
			if err != nil {
				return nil, err
			}
			fmt.Printf("Your generated password is %s \n", createUserRequest.Password)
			fmt.Printf("\n")
		}

		user, err := api.CreateUser(createUserRequest)
		if err != nil {
			return nil, err
		}

		result := rdbCreateUserResponseCustom{
			User:     user,
			Password: createUserRequest.Password,
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

	c.Run = func(ctx context.Context, argsI interface{}) (interface{}, error) {
		client := core.ExtractClient(ctx)
		api := rdb.NewAPI(client)

		customRequest := argsI.(*rdbUpdateUserRequestCustom)

		updateUserRequest := customRequest.UpdateUserRequest

		var err error
		if customRequest.GeneratePassword && customRequest.Password == nil {
			updateUserRequest.Password = new(string)
			*updateUserRequest.Password, err = passwordgenerator.GeneratePassword(21, 1, 1, 1, 1)
			if err != nil {
				return nil, err
			}
			fmt.Printf("Your generated password is %v \n", *updateUserRequest.Password)
			fmt.Printf("\n")
		}

		user, err := api.UpdateUser(updateUserRequest)
		if err != nil {
			return nil, err
		}

		result := rdbUpdateUserResponseCustom{
			User:     user,
			Password: *updateUserRequest.Password,
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
