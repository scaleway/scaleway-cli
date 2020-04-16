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
		Short:     `This API allows to manage your scaleway account`,
		Long:      ``,
		Namespace: "account",
	}
}

func accountSSHKey() *core.Command {
	return &core.Command{
		Short:     `Manage your Scaleway SSH keys`,
		Long:      `Manage your Scaleway SSH keys.`,
		Namespace: "account",
		Resource:  "ssh-key",
	}
}

func accountSSHKeyList() *core.Command {
	return &core.Command{
		Short:     `List all SSH keys`,
		Long:      `List all SSH keys.`,
		Namespace: "account",
		Resource:  "ssh-key",
		Verb:      "list",
		ArgsType:  reflect.TypeOf(account.ListSSHKeysRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Required:   false,
				Positional: false,
				EnumValues: []string{"created_at_asc", "created_at_desc", "updated_at_asc", "updated_at_desc", "name_asc", "name_desc"},
			},
			{
				Name:       "name",
				Required:   false,
				Positional: false,
			},
			{
				Name:       "organization-id",
				Required:   false,
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
				FieldName: "id",
			},
			{
				FieldName: "name",
			},
			{
				FieldName: "created_at",
			},
			{
				FieldName: "updated_at",
			},
			{
				FieldName: "organization_id",
			},
			{
				FieldName: "creation_info.address",
			},
			{
				FieldName: "creation_info.country_code",
			},
			{
				FieldName: "creation_info.user_agent",
			},
		}},
	}
}

func accountSSHKeyAdd() *core.Command {
	return &core.Command{
		Short:     `Add a SSH key to your Scaleway account`,
		Long:      `Add a SSH key to your Scaleway account.`,
		Namespace: "account",
		Resource:  "ssh-key",
		Verb:      "add",
		ArgsType:  reflect.TypeOf(account.CreateSSHKeyRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "name",
				Short:      `The name of the SSH key`,
				Required:   false,
				Positional: false,
			},
			{
				Name:       "public-key",
				Short:      `SSH public key. Currently ssh-rsa, ssh-dss (DSA), ssh-ed25519 and ecdsa keys with NIST curves are supported`,
				Required:   true,
				Positional: false,
			},
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
				Raw:   `scw account ssh-key add name=foobar public_key="$(cat <path/to/your/public/key>)"`,
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
		Short:     `Get SSH key details`,
		Long:      `Get SSH key details.`,
		Namespace: "account",
		Resource:  "ssh-key",
		Verb:      "get",
		ArgsType:  reflect.TypeOf(account.GetSSHKeyRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "ssh-key-id",
				Short:      `The ID of the SSH key`,
				Required:   true,
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
		Short:     `Update an SSH key`,
		Long:      `Update an SSH key.`,
		Namespace: "account",
		Resource:  "ssh-key",
		Verb:      "update",
		ArgsType:  reflect.TypeOf(account.UpdateSSHKeyRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "ssh-key-id",
				Required:   true,
				Positional: true,
			},
			{
				Name:       "name",
				Required:   false,
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
		Short:     `Remove a SSH key from your Scaleway account`,
		Long:      `Remove a SSH key from your Scaleway account.`,
		Namespace: "account",
		Resource:  "ssh-key",
		Verb:      "remove",
		ArgsType:  reflect.TypeOf(account.DeleteSSHKeyRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "ssh-key-id",
				Required:   true,
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
				Short:   "Remove a given SSH key",
				Request: `{"ssh_key_id":"11111111-1111-1111-1111-111111111111"}`,
			},
		},
	}
}
