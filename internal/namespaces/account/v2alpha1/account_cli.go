// This file was automatically generated. DO NOT EDIT.
// If you have any remark or suggestion do not hesitate to open an issue.

package account

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-sdk-go/api/account/v2alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// always import dependencies
var (
	_ = scw.RegionFrPar
)

func GetGeneratedCommands() *core.Commands {
	return core.NewCommands(
		accountRoot(),
		accountSSHKey(),
		accountSSHKeyList(),
		accountSSHKeyAdd(),
		accountSSHKeyGet(),
		accountSSHKeyUpdate(),
		accountSSHKeyRemove(),
	)
}
func accountRoot() *core.Command {
	return &core.Command{
		Short:     `Account API`,
		Long:      ``,
		Namespace: "account",
	}
}

func accountSSHKey() *core.Command {
	return &core.Command{
		Short:     `SSH keys management commands`,
		Long:      `SSH keys management commands.`,
		Namespace: "account",
		Resource:  "ssh-key",
	}
}

func accountSSHKeyList() *core.Command {
	return &core.Command{
		Short:     `List all SSH keys of your project`,
		Long:      `List all SSH keys of your project.`,
		Namespace: "account",
		Resource:  "ssh-key",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(account.ListSSHKeysRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"created_at_asc", "created_at_desc", "updated_at_asc", "updated_at_desc", "name_asc", "name_desc"},
			},
			{
				Name:       "name",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "project-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "organization-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*account.ListSSHKeysRequest)

			client := core.ExtractClient(ctx)
			api := account.NewAPI(client)
			resp, err := api.ListSSHKeys(request, scw.WithAllPages())
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
				FieldName: "UpdatedAt",
			},
			{
				FieldName: "ProjectID",
			},
			{
				FieldName: "CreationInfo.Address",
			},
			{
				FieldName: "CreationInfo.CountryCode",
			},
			{
				FieldName: "CreationInfo.UserAgent",
			},
			{
				FieldName: "OrganizationID",
			},
		}},
	}
}

func accountSSHKeyAdd() *core.Command {
	return &core.Command{
		Short:     `Add an SSH key to your project`,
		Long:      `Add an SSH key to your project.`,
		Namespace: "account",
		Resource:  "ssh-key",
		Verb:      "add",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(account.CreateSSHKeyRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "name",
				Short:      `The name of the SSH key`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "public-key",
				Short:      `SSH public key. Currently ssh-rsa, ssh-dss (DSA), ssh-ed25519 and ecdsa keys with NIST curves are supported`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.ProjectIDArgSpec(),
			core.OrganizationIDArgSpec(),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*account.CreateSSHKeyRequest)

			client := core.ExtractClient(ctx)
			api := account.NewAPI(client)
			return api.CreateSSHKey(request)

		},
		Examples: []*core.Example{
			{
				Short: "Add a given ssh key",
				Raw:   `scw account ssh-key add name=foobar public-key="$(cat <path/to/your/public/key>)"`,
			},
		},
		SeeAlsos: []*core.SeeAlso{
			{
				Command: "scw account ssh-key list",
				Short:   "List all SSH keys",
			},
			{
				Command: "scw account ssh-key remove",
				Short:   "Remove an SSH key",
			},
		},
	}
}

func accountSSHKeyGet() *core.Command {
	return &core.Command{
		Short:     `Get an SSH key from your project`,
		Long:      `Get an SSH key from your project.`,
		Namespace: "account",
		Resource:  "ssh-key",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(account.GetSSHKeyRequest{}),
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
			request := args.(*account.GetSSHKeyRequest)

			client := core.ExtractClient(ctx)
			api := account.NewAPI(client)
			return api.GetSSHKey(request)

		},
	}
}

func accountSSHKeyUpdate() *core.Command {
	return &core.Command{
		Short:     `Update an SSH key on your project`,
		Long:      `Update an SSH key on your project.`,
		Namespace: "account",
		Resource:  "ssh-key",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(account.UpdateSSHKeyRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "ssh-key-id",
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "name",
				Short:      `Name of the SSH key`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*account.UpdateSSHKeyRequest)

			client := core.ExtractClient(ctx)
			api := account.NewAPI(client)
			return api.UpdateSSHKey(request)

		},
	}
}

func accountSSHKeyRemove() *core.Command {
	return &core.Command{
		Short:     `Remove an SSH key from your project`,
		Long:      `Remove an SSH key from your project.`,
		Namespace: "account",
		Resource:  "ssh-key",
		Verb:      "remove",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(account.DeleteSSHKeyRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "ssh-key-id",
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*account.DeleteSSHKeyRequest)

			client := core.ExtractClient(ctx)
			api := account.NewAPI(client)
			e = api.DeleteSSHKey(request)
			if e != nil {
				return nil, e
			}
			return &core.SuccessResult{
				Resource: "ssh-key",
				Verb:     "remove",
			}, nil
		},
		Examples: []*core.Example{
			{
				Short:    "Remove a given SSH key",
				ArgsJSON: `{"ssh_key_id":"11111111-1111-1111-1111-111111111111"}`,
			},
		},
	}
}
