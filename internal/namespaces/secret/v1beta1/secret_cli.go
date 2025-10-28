// This file was automatically generated. DO NOT EDIT.
// If you have any remark or suggestion do not hesitate to open an issue.

package secret

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	secret "github.com/scaleway/scaleway-sdk-go/api/secret/v1beta1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// always import dependencies
var (
	_ = scw.RegionFrPar
)

func GetGeneratedCommands() *core.Commands {
	return core.NewCommands(
		secretRoot(),
		secretSecret(),
		secretVersion(),
		secretSecretCreate(),
		secretSecretGet(),
		secretSecretUpdate(),
		secretSecretDelete(),
		secretSecretList(),
		secretSecretProtect(),
		secretSecretUnprotect(),
		secretSecretAddOwner(),
		secretVersionCreate(),
		secretVersionGet(),
		secretVersionUpdate(),
		secretVersionDelete(),
		secretVersionList(),
		secretVersionAccess(),
		secretVersionAccessByPath(),
		secretVersionEnable(),
		secretVersionDisable(),
	)
}

func secretRoot() *core.Command {
	return &core.Command{
		Short:     `Secret Manager API`,
		Long:      `This API allows you to manage your Secret Manager services, for storing, accessing and sharing sensitive data such as passwords, API keys and certificates.`,
		Namespace: "secret",
	}
}

func secretSecret() *core.Command {
	return &core.Command{
		Short:     `Secret management commands`,
		Long:      `Secrets are logical containers made up of zero or more immutable versions, that contain sensitive data.`,
		Namespace: "secret",
		Resource:  "secret",
	}
}

func secretVersion() *core.Command {
	return &core.Command{
		Short:     `Secret Version management commands`,
		Long:      `Versions store the sensitive data contained in your secrets (API keys, passwords, or certificates).`,
		Namespace: "secret",
		Resource:  "version",
	}
}

func secretSecretCreate() *core.Command {
	return &core.Command{
		Short:     `Create a secret`,
		Long:      `Create a secret in a given region specified by the ` + "`" + `region` + "`" + ` parameter.`,
		Namespace: "secret",
		Resource:  "secret",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(secret.CreateSecretRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ProjectIDArgSpec(),
			{
				Name:       "name",
				Short:      `Name of the secret`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `List of the secret's tags`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "description",
				Short:      `Description of the secret`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "type",
				Short:      `Type of the secret`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_type",
					"opaque",
					"certificate",
					"key_value",
					"basic_credentials",
					"database_credentials",
					"ssh_key",
				},
			},
			{
				Name:       "path",
				Short:      `Path of the secret`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "ephemeral-policy.time-to-live",
				Short:      `Time frame, from one second and up to one year, during which the secret's versions are valid.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "ephemeral-policy.expires-once-accessed",
				Short:      `Returns ` + "`" + `true` + "`" + ` if the version expires after a single user access.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "ephemeral-policy.action",
				Short:      `Action to perform when the version of a secret expires`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_action",
					"delete",
					"disable",
				},
			},
			{
				Name:       "protected",
				Short:      `Returns ` + "`" + `true` + "`" + ` if secret protection is applied to a given secret`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "key-id",
				Short:      `ID of the Scaleway Key Manager key`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*secret.CreateSecretRequest)

			client := core.ExtractClient(ctx)
			api := secret.NewAPI(client)

			return api.CreateSecret(request)
		},
		Examples: []*core.Example{
			{
				Short: "Create a given secret",
				Raw:   `scw secret secret create name=foobar description="$(cat <path/to/your/secret>)"`,
			},
		},
	}
}

func secretSecretGet() *core.Command {
	return &core.Command{
		Short:     `Get metadata using the secret's ID`,
		Long:      `Retrieve the metadata of a secret specified by the ` + "`" + `region` + "`" + ` and ` + "`" + `secret_id` + "`" + ` parameters.`,
		Namespace: "secret",
		Resource:  "secret",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(secret.GetSecretRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "secret-id",
				Short:      `ID of the secret`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*secret.GetSecretRequest)

			client := core.ExtractClient(ctx)
			api := secret.NewAPI(client)

			return api.GetSecret(request)
		},
	}
}

func secretSecretUpdate() *core.Command {
	return &core.Command{
		Short:     `Update metadata of a secret`,
		Long:      `Edit a secret's metadata such as name, tag(s), description and ephemeral policy. The secret to update is specified by the ` + "`" + `secret_id` + "`" + ` and ` + "`" + `region` + "`" + ` parameters.`,
		Namespace: "secret",
		Resource:  "secret",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(secret.UpdateSecretRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "secret-id",
				Short:      `ID of the secret`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "name",
				Short:      `Secret's updated name (optional)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `Secret's updated list of tags (optional)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "description",
				Short:      `Description of the secret`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "path",
				Short:      `Path of the folder`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "ephemeral-policy.time-to-live",
				Short:      `Time frame, from one second and up to one year, during which the secret's versions are valid.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "ephemeral-policy.expires-once-accessed",
				Short:      `Returns ` + "`" + `true` + "`" + ` if the version expires after a single user access.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "ephemeral-policy.action",
				Short:      `Action to perform when the version of a secret expires`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_action",
					"delete",
					"disable",
				},
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*secret.UpdateSecretRequest)

			client := core.ExtractClient(ctx)
			api := secret.NewAPI(client)

			return api.UpdateSecret(request)
		},
	}
}

func secretSecretDelete() *core.Command {
	return &core.Command{
		Short:     `Delete a secret`,
		Long:      `Delete a given secret specified by the ` + "`" + `region` + "`" + ` and ` + "`" + `secret_id` + "`" + ` parameters.`,
		Namespace: "secret",
		Resource:  "secret",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(secret.DeleteSecretRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "secret-id",
				Short:      `ID of the secret`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*secret.DeleteSecretRequest)

			client := core.ExtractClient(ctx)
			api := secret.NewAPI(client)
			e = api.DeleteSecret(request)
			if e != nil {
				return nil, e
			}

			return &core.SuccessResult{
				Resource: "secret",
				Verb:     "delete",
			}, nil
		},
		Examples: []*core.Example{
			{
				Short:    "Delete a given secret",
				ArgsJSON: `{"secret_id":"11111111-1111-1111-1111-111111111111"}`,
			},
		},
	}
}

func secretSecretList() *core.Command {
	return &core.Command{
		Short:     `List secrets`,
		Long:      `Retrieve the list of secrets created within an Organization and/or Project. You must specify either the ` + "`" + `organization_id` + "`" + ` or the ` + "`" + `project_id` + "`" + ` and the ` + "`" + `region` + "`" + `.`,
		Namespace: "secret",
		Resource:  "secret",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(secret.ListSecretsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "project-id",
				Short:      `Filter by Project ID (optional)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "order-by",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"name_asc",
					"name_desc",
					"created_at_asc",
					"created_at_desc",
					"updated_at_asc",
					"updated_at_desc",
				},
			},
			{
				Name:       "tags.{index}",
				Short:      `List of tags to filter on (optional)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "name",
				Short:      `Filter by secret name (optional)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "path",
				Short:      `Filter by exact path (optional)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "ephemeral",
				Short:      `Filter by ephemeral / not ephemeral (optional)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "type",
				Short:      `Filter by secret type (optional)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_type",
					"opaque",
					"certificate",
					"key_value",
					"basic_credentials",
					"database_credentials",
					"ssh_key",
				},
			},
			{
				Name:       "scheduled-for-deletion",
				Short:      `Filter by whether the secret was scheduled for deletion / not scheduled for deletion. By default, it will display only not scheduled for deletion secrets.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "organization-id",
				Short:      `Filter by Organization ID (optional)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
				scw.Region(core.AllLocalities),
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*secret.ListSecretsRequest)

			client := core.ExtractClient(ctx)
			api := secret.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			if request.Region == scw.Region(core.AllLocalities) {
				opts = append(opts, scw.WithRegions(api.Regions()...))
				request.Region = ""
			}
			resp, err := api.ListSecrets(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.Secrets, nil
		},
	}
}

func secretSecretProtect() *core.Command {
	return &core.Command{
		Short:     `Enable secret protection`,
		Long:      `Enable secret protection for a given secret specified by the ` + "`" + `secret_id` + "`" + ` parameter. Enabling secret protection means that your secret can be read and modified, but it cannot be deleted.`,
		Namespace: "secret",
		Resource:  "secret",
		Verb:      "protect",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(secret.ProtectSecretRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "secret-id",
				Short:      `ID of the secret to enable secret protection for`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*secret.ProtectSecretRequest)

			client := core.ExtractClient(ctx)
			api := secret.NewAPI(client)

			return api.ProtectSecret(request)
		},
		Examples: []*core.Example{
			{
				Short:    "Enable secret protection",
				ArgsJSON: `{"secret_id":"11111111-1111-1111-1111-111111111111"}`,
			},
		},
	}
}

func secretSecretUnprotect() *core.Command {
	return &core.Command{
		Short:     `Disable secret protection`,
		Long:      `Disable secret protection for a given secret specified by the ` + "`" + `secret_id` + "`" + ` parameter. Disabling secret protection means that your secret can be read, modified and deleted.`,
		Namespace: "secret",
		Resource:  "secret",
		Verb:      "unprotect",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(secret.UnprotectSecretRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "secret-id",
				Short:      `ID of the secret to disable secret protection for`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*secret.UnprotectSecretRequest)

			client := core.ExtractClient(ctx)
			api := secret.NewAPI(client)

			return api.UnprotectSecret(request)
		},
		Examples: []*core.Example{
			{
				Short:    "Disable secret protection",
				ArgsJSON: `{"secret_id":"11111111-1111-1111-1111-111111111111"}`,
			},
		},
	}
}

func secretSecretAddOwner() *core.Command {
	return &core.Command{
		Short:     `Allow a product to use the secret`,
		Long:      `Allow a product to use the secret.`,
		Namespace: "secret",
		Resource:  "secret",
		Verb:      "add-owner",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(secret.AddSecretOwnerRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "secret-id",
				Short:      `ID of the secret`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "product",
				Short:      `ID of the product to add`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_product",
					"edge_services",
					"s2s_vpn",
				},
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*secret.AddSecretOwnerRequest)

			client := core.ExtractClient(ctx)
			api := secret.NewAPI(client)
			e = api.AddSecretOwner(request)
			if e != nil {
				return nil, e
			}

			return &core.SuccessResult{
				Resource: "secret",
				Verb:     "add-owner",
			}, nil
		},
	}
}

func secretVersionCreate() *core.Command {
	return &core.Command{
		Short:     `Create a version`,
		Long:      `Create a version of a given secret specified by the ` + "`" + `region` + "`" + ` and ` + "`" + `secret_id` + "`" + ` parameters.`,
		Namespace: "secret",
		Resource:  "version",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(secret.CreateSecretVersionRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "secret-id",
				Short:      `ID of the secret`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "data",
				Short:      `The base64-encoded secret payload of the version`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "description",
				Short:      `Description of the version`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "disable-previous",
				Short:      `Disable the previous secret version`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "data-crc32",
				Short:      `(Optional.) The CRC32 checksum of the data as a base-10 integer`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*secret.CreateSecretVersionRequest)

			client := core.ExtractClient(ctx)
			api := secret.NewAPI(client)

			return api.CreateSecretVersion(request)
		},
	}
}

func secretVersionGet() *core.Command {
	return &core.Command{
		Short:     `Get metadata of a secret's version using the secret's ID`,
		Long:      `Retrieve the metadata of a secret's given version specified by the ` + "`" + `region` + "`" + `, ` + "`" + `secret_id` + "`" + ` and ` + "`" + `revision` + "`" + ` parameters.`,
		Namespace: "secret",
		Resource:  "version",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(secret.GetSecretVersionRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "secret-id",
				Short:      `ID of the secret`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "revision",
				Short:      `Version number`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*secret.GetSecretVersionRequest)

			client := core.ExtractClient(ctx)
			api := secret.NewAPI(client)

			return api.GetSecretVersion(request)
		},
	}
}

func secretVersionUpdate() *core.Command {
	return &core.Command{
		Short:     `Update metadata of a version`,
		Long:      `Edit the metadata of a secret's given version, specified by the ` + "`" + `region` + "`" + `, ` + "`" + `secret_id` + "`" + ` and ` + "`" + `revision` + "`" + ` parameters.`,
		Namespace: "secret",
		Resource:  "version",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(secret.UpdateSecretVersionRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "secret-id",
				Short:      `ID of the secret`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "revision",
				Short:      `Version number`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "description",
				Short:      `Description of the version`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "ephemeral-properties.expires-at",
				Short:      `The version's expiration date`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "ephemeral-properties.expires-once-accessed",
				Short:      `Returns ` + "`" + `true` + "`" + ` if the version expires after a single user access.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "ephemeral-properties.action",
				Short:      `Action to perform when the version of a secret expires`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_action",
					"delete",
					"disable",
				},
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*secret.UpdateSecretVersionRequest)

			client := core.ExtractClient(ctx)
			api := secret.NewAPI(client)

			return api.UpdateSecretVersion(request)
		},
	}
}

func secretVersionDelete() *core.Command {
	return &core.Command{
		Short:     `Delete a version`,
		Long:      `Delete a secret's version and the sensitive data contained in it. Deleting a version is permanent and cannot be undone.`,
		Namespace: "secret",
		Resource:  "version",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(secret.DeleteSecretVersionRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "secret-id",
				Short:      `ID of the secret`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "revision",
				Short:      `Version number`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*secret.DeleteSecretVersionRequest)

			client := core.ExtractClient(ctx)
			api := secret.NewAPI(client)
			e = api.DeleteSecretVersion(request)
			if e != nil {
				return nil, e
			}

			return &core.SuccessResult{
				Resource: "version",
				Verb:     "delete",
			}, nil
		},
		Examples: []*core.Example{
			{
				Short:    "Delete a given Secret Version",
				ArgsJSON: `{"revision":"1","secret_id":"11111111-1111-1111-1111-111111111111"}`,
			},
		},
	}
}

func secretVersionList() *core.Command {
	return &core.Command{
		Short:     `List versions of a secret using the secret's ID`,
		Long:      `Retrieve the list of a given secret's versions specified by the ` + "`" + `secret_id` + "`" + ` and ` + "`" + `region` + "`" + ` parameters.`,
		Namespace: "secret",
		Resource:  "version",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(secret.ListSecretVersionsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "secret-id",
				Short:      `ID of the secret`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "status.{index}",
				Short:      `Filter results by status`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_status",
					"enabled",
					"disabled",
					"deleted",
					"scheduled_for_deletion",
				},
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
				scw.Region(core.AllLocalities),
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*secret.ListSecretVersionsRequest)

			client := core.ExtractClient(ctx)
			api := secret.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			if request.Region == scw.Region(core.AllLocalities) {
				opts = append(opts, scw.WithRegions(api.Regions()...))
				request.Region = ""
			}
			resp, err := api.ListSecretVersions(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.Versions, nil
		},
	}
}

func secretVersionAccess() *core.Command {
	return &core.Command{
		Short:     `Access a secret's version using the secret's ID`,
		Long:      `Access sensitive data in a secret's version specified by the ` + "`" + `region` + "`" + `, ` + "`" + `secret_id` + "`" + ` and ` + "`" + `revision` + "`" + ` parameters.`,
		Namespace: "secret",
		Resource:  "version",
		Verb:      "access",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(secret.AccessSecretVersionRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "secret-id",
				Short:      `ID of the secret`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "revision",
				Short:      `Version number`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*secret.AccessSecretVersionRequest)

			client := core.ExtractClient(ctx)
			api := secret.NewAPI(client)

			return api.AccessSecretVersion(request)
		},
	}
}

func secretVersionAccessByPath() *core.Command {
	return &core.Command{
		Short:     `Access a secret's version using the secret's name and path`,
		Long:      `Access sensitive data in a secret's version specified by the ` + "`" + `region` + "`" + `, ` + "`" + `secret_name` + "`" + `, ` + "`" + `secret_path` + "`" + ` and ` + "`" + `revision` + "`" + ` parameters.`,
		Namespace: "secret",
		Resource:  "version",
		Verb:      "access-by-path",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(secret.AccessSecretVersionByPathRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "secret-path",
				Short:      `Secret's path`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "secret-name",
				Short:      `Secret's name`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "revision",
				Short:      `Version number`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.ProjectIDArgSpec(),
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*secret.AccessSecretVersionByPathRequest)

			client := core.ExtractClient(ctx)
			api := secret.NewAPI(client)

			return api.AccessSecretVersionByPath(request)
		},
	}
}

func secretVersionEnable() *core.Command {
	return &core.Command{
		Short:     `Enable a version`,
		Long:      `Make a specific version accessible. You must specify the ` + "`" + `region` + "`" + `, ` + "`" + `secret_id` + "`" + ` and ` + "`" + `revision` + "`" + ` parameters.`,
		Namespace: "secret",
		Resource:  "version",
		Verb:      "enable",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(secret.EnableSecretVersionRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "secret-id",
				Short:      `ID of the secret`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "revision",
				Short:      `Version number`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*secret.EnableSecretVersionRequest)

			client := core.ExtractClient(ctx)
			api := secret.NewAPI(client)

			return api.EnableSecretVersion(request)
		},
	}
}

func secretVersionDisable() *core.Command {
	return &core.Command{
		Short:     `Disable a version`,
		Long:      `Make a specific version inaccessible. You must specify the ` + "`" + `region` + "`" + `, ` + "`" + `secret_id` + "`" + ` and ` + "`" + `revision` + "`" + ` parameters.`,
		Namespace: "secret",
		Resource:  "version",
		Verb:      "disable",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(secret.DisableSecretVersionRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "secret-id",
				Short:      `ID of the secret`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "revision",
				Short:      `Version number`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*secret.DisableSecretVersionRequest)

			client := core.ExtractClient(ctx)
			api := secret.NewAPI(client)

			return api.DisableSecretVersion(request)
		},
	}
}
