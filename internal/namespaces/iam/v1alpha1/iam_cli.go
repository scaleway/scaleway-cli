// This file was automatically generated. DO NOT EDIT.
// If you have any remark or suggestion do not hesitate to open an issue.

package iam

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/internal/core"
	"github.com/scaleway/scaleway-sdk-go/api/iam/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// always import dependencies
var (
	_ = scw.RegionFrPar
)

func GetGeneratedCommands() *core.Commands {
	return core.NewCommands(
		iamRoot(),
		iamSSHKey(),
		iamGroup(),
		iamAPIKey(),
		iamUser(),
		iamApplication(),
		iamPolicy(),
		iamRule(),
		iamPermissionSet(),
		iamSSHKeyList(),
		iamSSHKeyCreate(),
		iamSSHKeyGet(),
		iamSSHKeyUpdate(),
		iamSSHKeyDelete(),
		iamUserList(),
		iamUserGet(),
		iamUserDelete(),
		iamApplicationList(),
		iamApplicationCreate(),
		iamApplicationGet(),
		iamApplicationUpdate(),
		iamApplicationDelete(),
		iamGroupList(),
		iamGroupCreate(),
		iamGroupGet(),
		iamGroupUpdate(),
		iamGroupAddMember(),
		iamGroupRemoveMember(),
		iamGroupDelete(),
		iamPolicyList(),
		iamPolicyCreate(),
		iamPolicyGet(),
		iamPolicyUpdate(),
		iamPolicyDelete(),
		iamRuleUpdate(),
		iamRuleList(),
		iamPermissionSetList(),
		iamAPIKeyList(),
		iamAPIKeyCreate(),
		iamAPIKeyGet(),
		iamAPIKeyUpdate(),
		iamAPIKeyDelete(),
	)
}
func iamRoot() *core.Command {
	return &core.Command{
		Short:     `IAM API`,
		Long:      ``,
		Namespace: "iam",
	}
}

func iamSSHKey() *core.Command {
	return &core.Command{
		Short:     `SSH keys management commands`,
		Long:      `SSH keys management commands.`,
		Namespace: "iam",
		Resource:  "ssh-key",
	}
}

func iamGroup() *core.Command {
	return &core.Command{
		Short:     `Groups management commands`,
		Long:      `Groups management commands.`,
		Namespace: "iam",
		Resource:  "group",
	}
}

func iamAPIKey() *core.Command {
	return &core.Command{
		Short:     `API keys management commands`,
		Long:      `API keys management commands.`,
		Namespace: "iam",
		Resource:  "api-key",
	}
}

func iamUser() *core.Command {
	return &core.Command{
		Short:     `Users management commands`,
		Long:      `Users management commands.`,
		Namespace: "iam",
		Resource:  "user",
	}
}

func iamApplication() *core.Command {
	return &core.Command{
		Short:     `Applications management commands`,
		Long:      `Applications management commands.`,
		Namespace: "iam",
		Resource:  "application",
	}
}

func iamPolicy() *core.Command {
	return &core.Command{
		Short:     `Policies management commands`,
		Long:      `Policies management commands.`,
		Namespace: "iam",
		Resource:  "policy",
	}
}

func iamRule() *core.Command {
	return &core.Command{
		Short:     `Rules management commands`,
		Long:      `Rules management commands.`,
		Namespace: "iam",
		Resource:  "rule",
	}
}

func iamPermissionSet() *core.Command {
	return &core.Command{
		Short:     `Permission sets management commands`,
		Long:      `Permission sets management commands.`,
		Namespace: "iam",
		Resource:  "permission-set",
	}
}

func iamSSHKeyList() *core.Command {
	return &core.Command{
		Short:     `List SSH keys`,
		Long:      `List SSH keys.`,
		Namespace: "iam",
		Resource:  "ssh-key",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(iam.ListSSHKeysRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `Sort order of SSH keys`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				Default:    core.DefaultValueSetter("created_at_asc"),
				EnumValues: []string{"created_at_asc", "created_at_desc", "updated_at_asc", "updated_at_desc", "name_asc", "name_desc"},
			},
			{
				Name:       "name",
				Short:      `Name of group to find`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "project-id",
				Short:      `Filter by project ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "disabled",
				Short:      `Filter out disabled SSH keys or not`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "organization-id",
				Short:      `Filter by organization ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*iam.ListSSHKeysRequest)

			client := core.ExtractClient(ctx)
			api := iam.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			resp, err := api.ListSSHKeys(request, opts...)
			if err != nil {
				return nil, err
			}
			return resp.SSHKeys, nil

		},
		View: &core.View{Fields: []*core.ViewField{
			{
				FieldName: "ID",
			},
			{
				FieldName: "Name",
			},
			{
				FieldName: "CreatedAt",
			},
			{
				FieldName: "Fingerprint",
			},
			{
				FieldName: "ProjectID",
			},
			{
				FieldName: "Disabled",
			},
		}},
	}
}

func iamSSHKeyCreate() *core.Command {
	return &core.Command{
		Short:     `Create an SSH key`,
		Long:      `Create an SSH key.`,
		Namespace: "iam",
		Resource:  "ssh-key",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(iam.CreateSSHKeyRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "name",
				Short:      `The name of the SSH key. Max length is 1000`,
				Required:   true,
				Deprecated: false,
				Positional: false,
				Default:    core.RandomValueGenerator("key"),
			},
			{
				Name:       "public-key",
				Short:      `SSH public key. Currently ssh-rsa, ssh-dss (DSA), ssh-ed25519 and ecdsa keys with NIST curves are supported. Max length is 65000`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.ProjectIDArgSpec(),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*iam.CreateSSHKeyRequest)

			client := core.ExtractClient(ctx)
			api := iam.NewAPI(client)
			return api.CreateSSHKey(request)

		},
		Examples: []*core.Example{
			{
				Short: "Add a given ssh key",
				Raw:   `scw iam ssh-key create name=foobar public-key="$(cat <path/to/your/public/key>)"`,
			},
		},
		SeeAlsos: []*core.SeeAlso{
			{
				Command: "scw iam ssh-key list",
				Short:   "List all SSH keys",
			},
			{
				Command: "scw iam ssh-key delete",
				Short:   "Delete an SSH key",
			},
		},
	}
}

func iamSSHKeyGet() *core.Command {
	return &core.Command{
		Short:     `Get an SSH key`,
		Long:      `Get an SSH key.`,
		Namespace: "iam",
		Resource:  "ssh-key",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(iam.GetSSHKeyRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "ssh-key-id",
				Short:      `The ID of the SSH key`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*iam.GetSSHKeyRequest)

			client := core.ExtractClient(ctx)
			api := iam.NewAPI(client)
			return api.GetSSHKey(request)

		},
	}
}

func iamSSHKeyUpdate() *core.Command {
	return &core.Command{
		Short:     `Update an SSH key`,
		Long:      `Update an SSH key.`,
		Namespace: "iam",
		Resource:  "ssh-key",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(iam.UpdateSSHKeyRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "ssh-key-id",
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "name",
				Short:      `Name of the SSH key. Max length is 1000`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "disabled",
				Short:      `Enable or disable the SSH key`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*iam.UpdateSSHKeyRequest)

			client := core.ExtractClient(ctx)
			api := iam.NewAPI(client)
			return api.UpdateSSHKey(request)

		},
	}
}

func iamSSHKeyDelete() *core.Command {
	return &core.Command{
		Short:     `Delete an SSH key`,
		Long:      `Delete an SSH key.`,
		Namespace: "iam",
		Resource:  "ssh-key",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(iam.DeleteSSHKeyRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "ssh-key-id",
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*iam.DeleteSSHKeyRequest)

			client := core.ExtractClient(ctx)
			api := iam.NewAPI(client)
			e = api.DeleteSSHKey(request)
			if e != nil {
				return nil, e
			}
			return &core.SuccessResult{
				Resource: "ssh-key",
				Verb:     "delete",
			}, nil
		},
		Examples: []*core.Example{
			{
				Short:    "Delete a given SSH key",
				ArgsJSON: `{"ssh_key_id":"11111111-1111-1111-1111-111111111111"}`,
			},
		},
		SeeAlsos: []*core.SeeAlso{
			{
				Command: "scw iam ssh-key list",
				Short:   "List all SSH keys",
			},
			{
				Command: "scw iam ssh-key create",
				Short:   "Add a SSH key",
			},
		},
	}
}

func iamUserList() *core.Command {
	return &core.Command{
		Short:     `List users of an organization`,
		Long:      `List users of an organization.`,
		Namespace: "iam",
		Resource:  "user",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(iam.ListUsersRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `Criteria for sorting results`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				Default:    core.DefaultValueSetter("created_at_asc"),
				EnumValues: []string{"created_at_asc", "created_at_desc", "updated_at_asc", "updated_at_desc", "email_asc", "email_desc", "last_login_asc", "last_login_desc"},
			},
			{
				Name:       "user-ids.{index}",
				Short:      `Filter out by a list of ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "organization-id",
				Short:      `ID of organization to filter`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*iam.ListUsersRequest)

			client := core.ExtractClient(ctx)
			api := iam.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			resp, err := api.ListUsers(request, opts...)
			if err != nil {
				return nil, err
			}
			return resp.Users, nil

		},
	}
}

func iamUserGet() *core.Command {
	return &core.Command{
		Short:     `Retrieve a user from its ID`,
		Long:      `Retrieve a user from its ID.`,
		Namespace: "iam",
		Resource:  "user",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(iam.GetUserRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "user-id",
				Short:      `ID of user to find`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*iam.GetUserRequest)

			client := core.ExtractClient(ctx)
			api := iam.NewAPI(client)
			return api.GetUser(request)

		},
	}
}

func iamUserDelete() *core.Command {
	return &core.Command{
		Short:     `Delete a guest user from an organization`,
		Long:      `Delete a guest user from an organization.`,
		Namespace: "iam",
		Resource:  "user",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(iam.DeleteUserRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "user-id",
				Short:      `ID of user to delete`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*iam.DeleteUserRequest)

			client := core.ExtractClient(ctx)
			api := iam.NewAPI(client)
			e = api.DeleteUser(request)
			if e != nil {
				return nil, e
			}
			return &core.SuccessResult{
				Resource: "user",
				Verb:     "delete",
			}, nil
		},
	}
}

func iamApplicationList() *core.Command {
	return &core.Command{
		Short:     `List applications of an organization`,
		Long:      `List applications of an organization.`,
		Namespace: "iam",
		Resource:  "application",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(iam.ListApplicationsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `Criteria for sorting results`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				Default:    core.DefaultValueSetter("created_at_asc"),
				EnumValues: []string{"created_at_asc", "created_at_desc", "updated_at_asc", "updated_at_desc", "name_asc", "name_desc"},
			},
			{
				Name:       "name",
				Short:      `Name of application to filter`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "editable",
				Short:      `Filter out editable applications or not`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "application-ids.{index}",
				Short:      `Filter out by a list of ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "organization-id",
				Short:      `ID of organization to filter`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*iam.ListApplicationsRequest)

			client := core.ExtractClient(ctx)
			api := iam.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			resp, err := api.ListApplications(request, opts...)
			if err != nil {
				return nil, err
			}
			return resp.Applications, nil

		},
	}
}

func iamApplicationCreate() *core.Command {
	return &core.Command{
		Short:     `Create a new application`,
		Long:      `Create a new application.`,
		Namespace: "iam",
		Resource:  "application",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(iam.CreateApplicationRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "name",
				Short:      `Name of application to create (max length is 64 chars)`,
				Required:   true,
				Deprecated: false,
				Positional: false,
				Default:    core.RandomValueGenerator("app"),
			},
			{
				Name:       "description",
				Short:      `Description of application (max length is 200 chars)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.OrganizationIDArgSpec(),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*iam.CreateApplicationRequest)

			client := core.ExtractClient(ctx)
			api := iam.NewAPI(client)
			return api.CreateApplication(request)

		},
	}
}

func iamApplicationGet() *core.Command {
	return &core.Command{
		Short:     `Get an existing application`,
		Long:      `Get an existing application.`,
		Namespace: "iam",
		Resource:  "application",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(iam.GetApplicationRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "application-id",
				Short:      `ID of application to find`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*iam.GetApplicationRequest)

			client := core.ExtractClient(ctx)
			api := iam.NewAPI(client)
			return api.GetApplication(request)

		},
	}
}

func iamApplicationUpdate() *core.Command {
	return &core.Command{
		Short:     `Update an existing application`,
		Long:      `Update an existing application.`,
		Namespace: "iam",
		Resource:  "application",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(iam.UpdateApplicationRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "application-id",
				Short:      `ID of application to update`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "name",
				Short:      `New name of application (max length is 64 chars)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "description",
				Short:      `New description of application (max length is 200 chars)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*iam.UpdateApplicationRequest)

			client := core.ExtractClient(ctx)
			api := iam.NewAPI(client)
			return api.UpdateApplication(request)

		},
	}
}

func iamApplicationDelete() *core.Command {
	return &core.Command{
		Short:     `Delete an application`,
		Long:      `Delete an application.`,
		Namespace: "iam",
		Resource:  "application",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(iam.DeleteApplicationRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "application-id",
				Short:      `ID of application to delete`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*iam.DeleteApplicationRequest)

			client := core.ExtractClient(ctx)
			api := iam.NewAPI(client)
			e = api.DeleteApplication(request)
			if e != nil {
				return nil, e
			}
			return &core.SuccessResult{
				Resource: "application",
				Verb:     "delete",
			}, nil
		},
	}
}

func iamGroupList() *core.Command {
	return &core.Command{
		Short:     `List groups`,
		Long:      `List groups.`,
		Namespace: "iam",
		Resource:  "group",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(iam.ListGroupsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `Sort order of groups`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				Default:    core.DefaultValueSetter("created_at_asc"),
				EnumValues: []string{"created_at_asc", "created_at_desc", "updated_at_asc", "updated_at_desc", "name_asc", "name_desc"},
			},
			{
				Name:       "name",
				Short:      `Name of group to find`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "application-ids.{index}",
				Short:      `Filter out by a list of application ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "user-ids.{index}",
				Short:      `Filter out by a list of user ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "group-ids.{index}",
				Short:      `Filter out by a list of group ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "organization-id",
				Short:      `Filter by organization ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*iam.ListGroupsRequest)

			client := core.ExtractClient(ctx)
			api := iam.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			resp, err := api.ListGroups(request, opts...)
			if err != nil {
				return nil, err
			}
			return resp.Groups, nil

		},
		View: &core.View{Fields: []*core.ViewField{
			{
				FieldName: "ID",
			},
			{
				FieldName: "Name",
			},
			{
				FieldName: "UserIDs",
			},
			{
				FieldName: "ApplicationIDs",
			},
		}},
	}
}

func iamGroupCreate() *core.Command {
	return &core.Command{
		Short:     `Create a new group`,
		Long:      `Create a new group.`,
		Namespace: "iam",
		Resource:  "group",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(iam.CreateGroupRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "name",
				Short:      `Name of the group to create (max length is 64 chars). MUST be unique inside an organization`,
				Required:   true,
				Deprecated: false,
				Positional: false,
				Default:    core.RandomValueGenerator("grp"),
			},
			{
				Name:       "description",
				Short:      `Description of the group to create (max length is 200 chars)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.OrganizationIDArgSpec(),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*iam.CreateGroupRequest)

			client := core.ExtractClient(ctx)
			api := iam.NewAPI(client)
			return api.CreateGroup(request)

		},
		Examples: []*core.Example{
			{
				Short: "Create a group",
				Raw:   `scw iam group create name=foobar`,
			},
		},
		SeeAlsos: []*core.SeeAlso{
			{
				Command: "scw iam group add-member",
				Short:   "Add a group member",
			},
			{
				Command: "scw iam group delete",
				Short:   "Delete a group",
			},
			{
				Command: "scw iam policy create",
				Short:   "Create a policy for a group",
			},
		},
	}
}

func iamGroupGet() *core.Command {
	return &core.Command{
		Short:     `Get a group`,
		Long:      `Get a group.`,
		Namespace: "iam",
		Resource:  "group",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(iam.GetGroupRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "group-id",
				Short:      `ID of group`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*iam.GetGroupRequest)

			client := core.ExtractClient(ctx)
			api := iam.NewAPI(client)
			return api.GetGroup(request)

		},
	}
}

func iamGroupUpdate() *core.Command {
	return &core.Command{
		Short:     `Update a group`,
		Long:      `Update a group.`,
		Namespace: "iam",
		Resource:  "group",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(iam.UpdateGroupRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "group-id",
				Short:      `ID of group to update`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "name",
				Short:      `New name for the group (max length is 64 chars). MUST be unique inside an organization`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "description",
				Short:      `New description for the group (max length is 200 chars)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*iam.UpdateGroupRequest)

			client := core.ExtractClient(ctx)
			api := iam.NewAPI(client)
			return api.UpdateGroup(request)

		},
	}
}

func iamGroupAddMember() *core.Command {
	return &core.Command{
		Short:     `Add a user of an application to a group`,
		Long:      `Add a user of an application to a group.`,
		Namespace: "iam",
		Resource:  "group",
		Verb:      "add-member",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(iam.AddGroupMemberRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "group-id",
				Short:      `ID of group`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "user-id",
				Short:      `ID of the user to add`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "application-id",
				Short:      `ID of the application to add`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*iam.AddGroupMemberRequest)

			client := core.ExtractClient(ctx)
			api := iam.NewAPI(client)
			return api.AddGroupMember(request)

		},
	}
}

func iamGroupRemoveMember() *core.Command {
	return &core.Command{
		Short:     `Remove a user or an application from a group`,
		Long:      `Remove a user or an application from a group.`,
		Namespace: "iam",
		Resource:  "group",
		Verb:      "remove-member",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(iam.RemoveGroupMemberRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "group-id",
				Short:      `ID of group`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "user-id",
				Short:      `ID of the user to remove`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "application-id",
				Short:      `ID of the application to remove`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*iam.RemoveGroupMemberRequest)

			client := core.ExtractClient(ctx)
			api := iam.NewAPI(client)
			return api.RemoveGroupMember(request)

		},
		SeeAlsos: []*core.SeeAlso{
			{
				Command: "scw iam group remove-member",
				Short:   "Remove a group member",
			},
			{
				Command: "scw iam group create",
				Short:   "Create a group",
			},
		},
	}
}

func iamGroupDelete() *core.Command {
	return &core.Command{
		Short:     `Delete a group`,
		Long:      `Delete a group.`,
		Namespace: "iam",
		Resource:  "group",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(iam.DeleteGroupRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "group-id",
				Short:      `ID of group to delete`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*iam.DeleteGroupRequest)

			client := core.ExtractClient(ctx)
			api := iam.NewAPI(client)
			e = api.DeleteGroup(request)
			if e != nil {
				return nil, e
			}
			return &core.SuccessResult{
				Resource: "group",
				Verb:     "delete",
			}, nil
		},
		Examples: []*core.Example{
			{
				Short:    "Delete a given group",
				ArgsJSON: `{"group_id":"11111111-1111-1111-1111-111111111111"}`,
			},
		},
		SeeAlsos: []*core.SeeAlso{
			{
				Command: "scw iam group list",
				Short:   "List all groups",
			},
			{
				Command: "scw iam group delete",
				Short:   "Delete a group",
			},
		},
	}
}

func iamPolicyList() *core.Command {
	return &core.Command{
		Short:     `List policies of an organization`,
		Long:      `List policies of an organization.`,
		Namespace: "iam",
		Resource:  "policy",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(iam.ListPoliciesRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `Criteria for sorting results`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				Default:    core.DefaultValueSetter("created_at_asc"),
				EnumValues: []string{"policy_name_asc", "policy_name_desc", "created_at_asc", "created_at_desc"},
			},
			{
				Name:       "editable",
				Short:      `Filter out editable policies or not`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "user-ids.{index}",
				Short:      `Filter out by a list of user ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "group-ids.{index}",
				Short:      `Filter out by a list of group ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "application-ids.{index}",
				Short:      `Filter out by a list of application ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "no-principal",
				Short:      `True when the policy do not belong to any principal`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "policy-name",
				Short:      `Name of policy to fetch`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "organization-id",
				Short:      `ID of organization to filter`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*iam.ListPoliciesRequest)

			client := core.ExtractClient(ctx)
			api := iam.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			resp, err := api.ListPolicies(request, opts...)
			if err != nil {
				return nil, err
			}
			return resp.Policies, nil

		},
	}
}

func iamPolicyCreate() *core.Command {
	return &core.Command{
		Short:     `Create a new policy`,
		Long:      `Create a new policy.`,
		Namespace: "iam",
		Resource:  "policy",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(iam.CreatePolicyRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "name",
				Short:      `Name of policy to create (max length is 64 chars)`,
				Required:   true,
				Deprecated: false,
				Positional: false,
				Default:    core.RandomValueGenerator("pol"),
			},
			{
				Name:       "description",
				Short:      `Description of policy to create (max length is 200 chars)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "rules.{index}.permission-set-names.{index}",
				Short:      `Names of permission sets bound to the rule`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "rules.{index}.project-ids.{index}",
				Short:      `List of project IDs scoped to the rule`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "rules.{index}.organization-id",
				Short:      `ID of organization scoped to the rule`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "user-id",
				Short:      `ID of user, owner of the policy`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "group-id",
				Short:      `ID of group, owner of the policy`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "application-id",
				Short:      `ID of application, owner of the policy`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "no-principal",
				Short:      `True when the policy do not belong to any principal`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.OrganizationIDArgSpec(),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*iam.CreatePolicyRequest)

			client := core.ExtractClient(ctx)
			api := iam.NewAPI(client)
			return api.CreatePolicy(request)

		},
		Examples: []*core.Example{
			{
				Short: "Add a policy for a group that gives InstanceFullAccess on all projects",
				Raw:   `scw iam policy create group-id=11111111-1111-1111-1111-111111111111 rules.0.organization-id=11111111-1111-1111-1111-111111111111 rules.0.permission-set-names.0=InstanceFullAccess`,
			},
		},
	}
}

func iamPolicyGet() *core.Command {
	return &core.Command{
		Short:     `Get an existing policy`,
		Long:      `Get an existing policy.`,
		Namespace: "iam",
		Resource:  "policy",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(iam.GetPolicyRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "policy-id",
				Short:      `Id of policy to search`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*iam.GetPolicyRequest)

			client := core.ExtractClient(ctx)
			api := iam.NewAPI(client)
			return api.GetPolicy(request)

		},
	}
}

func iamPolicyUpdate() *core.Command {
	return &core.Command{
		Short:     `Update an existing policy`,
		Long:      `Update an existing policy.`,
		Namespace: "iam",
		Resource:  "policy",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(iam.UpdatePolicyRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "policy-id",
				Short:      `Id of policy to update`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "name",
				Short:      `New name of policy (max length is 64 chars)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "description",
				Short:      `New description of policy (max length is 200 chars)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "user-id",
				Short:      `New ID of user, owner of the policy`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "group-id",
				Short:      `New ID of group, owner of the policy`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "application-id",
				Short:      `New ID of application, owner of the policy`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "no-principal",
				Short:      `True when the policy do not belong to any principal`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*iam.UpdatePolicyRequest)

			client := core.ExtractClient(ctx)
			api := iam.NewAPI(client)
			return api.UpdatePolicy(request)

		},
	}
}

func iamPolicyDelete() *core.Command {
	return &core.Command{
		Short:     `Delete a policy`,
		Long:      `Delete a policy.`,
		Namespace: "iam",
		Resource:  "policy",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(iam.DeletePolicyRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "policy-id",
				Short:      `Id of policy to delete`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*iam.DeletePolicyRequest)

			client := core.ExtractClient(ctx)
			api := iam.NewAPI(client)
			e = api.DeletePolicy(request)
			if e != nil {
				return nil, e
			}
			return &core.SuccessResult{
				Resource: "policy",
				Verb:     "delete",
			}, nil
		},
	}
}

func iamRuleUpdate() *core.Command {
	return &core.Command{
		Short:     `Set rules of an existing policy`,
		Long:      `Set rules of an existing policy.`,
		Namespace: "iam",
		Resource:  "rule",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(iam.SetRulesRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "policy-id",
				Short:      `Id of policy to update`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "rules.{index}.permission-set-names.{index}",
				Short:      `Names of permission sets bound to the rule`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "rules.{index}.project-ids.{index}",
				Short:      `List of project IDs scoped to the rule`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "rules.{index}.organization-id",
				Short:      `ID of organization scoped to the rule`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*iam.SetRulesRequest)

			client := core.ExtractClient(ctx)
			api := iam.NewAPI(client)
			return api.SetRules(request)

		},
	}
}

func iamRuleList() *core.Command {
	return &core.Command{
		Short:     `List rules of an existing policy`,
		Long:      `List rules of an existing policy.`,
		Namespace: "iam",
		Resource:  "rule",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(iam.ListRulesRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "policy-id",
				Short:      `Id of policy to search`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*iam.ListRulesRequest)

			client := core.ExtractClient(ctx)
			api := iam.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			resp, err := api.ListRules(request, opts...)
			if err != nil {
				return nil, err
			}
			return resp.Rules, nil

		},
	}
}

func iamPermissionSetList() *core.Command {
	return &core.Command{
		Short:     `List permission sets`,
		Long:      `List permission sets.`,
		Namespace: "iam",
		Resource:  "permission-set",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(iam.ListPermissionSetsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `Criteria for sorting results`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				Default:    core.DefaultValueSetter("created_at_asc"),
				EnumValues: []string{"name_asc", "name_desc", "created_at_asc", "created_at_desc"},
			},
			core.OrganizationIDArgSpec(),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*iam.ListPermissionSetsRequest)

			client := core.ExtractClient(ctx)
			api := iam.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			resp, err := api.ListPermissionSets(request, opts...)
			if err != nil {
				return nil, err
			}
			return resp.PermissionSets, nil

		},
	}
}

func iamAPIKeyList() *core.Command {
	return &core.Command{
		Short:     `List API keys`,
		Long:      `List API keys.`,
		Namespace: "iam",
		Resource:  "api-key",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(iam.ListAPIKeysRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `Criteria for sorting results`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				Default:    core.DefaultValueSetter("created_at_asc"),
				EnumValues: []string{"created_at_asc", "created_at_desc", "updated_at_asc", "updated_at_desc", "expires_at_asc", "expires_at_desc", "access_key_asc", "access_key_desc"},
			},
			{
				Name:       "application-id",
				Short:      `ID of an application bearer`,
				Required:   false,
				Deprecated: true,
				Positional: false,
			},
			{
				Name:       "user-id",
				Short:      `ID of a user bearer`,
				Required:   false,
				Deprecated: true,
				Positional: false,
			},
			{
				Name:       "editable",
				Short:      `Filter out editable API keys or not`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "expirable",
				Short:      `Filter out expirable API keys or not`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "access-key",
				Short:      `Filter out by access key`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "description",
				Short:      `Filter out by description`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "bearer-id",
				Short:      `Filter out by bearer ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "bearer-type",
				Short:      `Filter out by type of bearer`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"unknown_bearer_type", "user", "application"},
			},
			{
				Name:       "organization-id",
				Short:      `ID of organization`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*iam.ListAPIKeysRequest)

			client := core.ExtractClient(ctx)
			api := iam.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			resp, err := api.ListAPIKeys(request, opts...)
			if err != nil {
				return nil, err
			}
			return resp.APIKeys, nil

		},
		View: &core.View{Fields: []*core.ViewField{
			{
				FieldName: "AccessKey",
			},
			{
				FieldName: "SecretKey",
			},
			{
				FieldName: "CreatedAt",
			},
			{
				FieldName: "ExpiresAt",
			},
			{
				FieldName: "DefaultProjectID",
			},
		}},
	}
}

func iamAPIKeyCreate() *core.Command {
	return &core.Command{
		Short:     `Create an API key`,
		Long:      `Create an API key.`,
		Namespace: "iam",
		Resource:  "api-key",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(iam.CreateAPIKeyRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "application-id",
				Short:      `ID of application principal`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "user-id",
				Short:      `ID of user principal`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "expires-at",
				Short:      `Expiration date of the API key`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "default-project-id",
				Short:      `The default project ID to use with object storage`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "description",
				Short:      `The description of the API key (max length is 200 chars)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*iam.CreateAPIKeyRequest)

			client := core.ExtractClient(ctx)
			api := iam.NewAPI(client)
			return api.CreateAPIKey(request)

		},
		SeeAlsos: []*core.SeeAlso{
			{
				Command: "scw iam api-key list",
				Short:   "List all API keys",
			},
			{
				Command: "scw iam api-key delete",
				Short:   "Delete an API key",
			},
		},
	}
}

func iamAPIKeyGet() *core.Command {
	return &core.Command{
		Short:     `Get an API key`,
		Long:      `Get an API key.`,
		Namespace: "iam",
		Resource:  "api-key",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(iam.GetAPIKeyRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "access-key",
				Short:      `Access key to search for`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*iam.GetAPIKeyRequest)

			client := core.ExtractClient(ctx)
			api := iam.NewAPI(client)
			return api.GetAPIKey(request)

		},
	}
}

func iamAPIKeyUpdate() *core.Command {
	return &core.Command{
		Short:     `Update an API key`,
		Long:      `Update an API key.`,
		Namespace: "iam",
		Resource:  "api-key",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(iam.UpdateAPIKeyRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "access-key",
				Short:      `Access key to update`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "default-project-id",
				Short:      `The new default project ID to set`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "description",
				Short:      `The new description to update`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*iam.UpdateAPIKeyRequest)

			client := core.ExtractClient(ctx)
			api := iam.NewAPI(client)
			return api.UpdateAPIKey(request)

		},
	}
}

func iamAPIKeyDelete() *core.Command {
	return &core.Command{
		Short:     `Delete an API key`,
		Long:      `Delete an API key.`,
		Namespace: "iam",
		Resource:  "api-key",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(iam.DeleteAPIKeyRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "access-key",
				Short:      `Access key to delete`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*iam.DeleteAPIKeyRequest)

			client := core.ExtractClient(ctx)
			api := iam.NewAPI(client)
			e = api.DeleteAPIKey(request)
			if e != nil {
				return nil, e
			}
			return &core.SuccessResult{
				Resource: "api-key",
				Verb:     "delete",
			}, nil
		},
		Examples: []*core.Example{
			{
				Short:    "Delete a given API key",
				ArgsJSON: `{"access_key":"SCW00000000000"}`,
			},
		},
		SeeAlsos: []*core.SeeAlso{
			{
				Command: "scw iam api-key list",
				Short:   "List all API keys",
			},
			{
				Command: "scw iam api-key create",
				Short:   "Create an API key",
			},
		},
	}
}
