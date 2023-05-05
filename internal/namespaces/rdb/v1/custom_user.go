package rdb

import (
	"context"

	"github.com/scaleway/scaleway-cli/v2/internal/core"
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
			&rdb.ListPrivilegesRequest{InstanceID: listUserRequest.InstanceID},
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
				core.ExtractLogger(ctx).Errorf("unsupported permission value %s", privilege.Permission)
			}
		}
		return res, nil
	}

	return c
}
