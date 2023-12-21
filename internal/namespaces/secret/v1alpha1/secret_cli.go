// This file was automatically generated. DO NOT EDIT.
// If you have any remark or suggestion do not hesitate to open an issue.

package secret

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/internal/core"
	"github.com/scaleway/scaleway-sdk-go/api/secret/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// always import dependencies
var (
	_ = scw.RegionFrPar
)

func GetGeneratedCommands() *core.Commands {
	return core.NewCommands(
		secretRoot(),
		secretFolder(),
		secretSecret(),
		secretVersion(),
		secretTag(),
		secretSecretCreate(),
		secretFolderCreate(),
		secretSecretGet(),
		secretSecretUpdate(),
		secretSecretList(),
		secretFolderList(),
		secretSecretDelete(),
		secretFolderDelete(),
		secretSecretProtect(),
		secretSecretUnprotect(),
		secretSecretAddOwner(),
		secretVersionCreate(),
		secretVersionGeneratePassword(),
		secretVersionGet(),
		secretVersionUpdate(),
		secretVersionList(),
		secretVersionEnable(),
		secretVersionDisable(),
		secretVersionAccess(),
		secretVersionDelete(),
		secretTagList(),
	)
}
func secretRoot() *core.Command {
	return &core.Command{
		Short:     `Secret Manager API`,
		Long:      `This API allows you to conveniently store, access and share sensitive data.`,
		Namespace: "secret",
	}
}

func secretFolder() *core.Command {
	return &core.Command{
		Short:     `Folder management commands`,
		Long:      `Location of the secret in the directory structure.`,
		Namespace: "secret",
		Resource:  "folder",
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

func secretTag() *core.Command {
	return &core.Command{
		Short:     `Tag management commands`,
		Long:      `Tag management commands.`,
		Namespace: "secret",
		Resource:  "tag",
	}
}

func secretSecretCreate() *core.Command {
	return &core.Command{
		Short:     `Create a secret`,
		Long:      `You must specify the ` + "`" + `region` + "`" + ` to create a secret.`,
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
				EnumValues: []string{"unknown_secret_type", "opaque", "certificate", "key_value"},
			},
			{
				Name:       "path",
				Short:      `Path of the secret`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "ephemeral-policy-template.time-to-live",
				Short:      `Time frame, from one second and up to one year, during which the secret's versions are valid.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "ephemeral-policy-template.expires-once-accessed",
				Short:      `Returns ` + "`" + `true` + "`" + ` if the version expires after a single user access.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "ephemeral-policy-template.action",
				Short:      `Action to perform when the version of a secret expires`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"unknown_action", "delete", "disable"},
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*secret.CreateSecretRequest)

			client := core.ExtractClient(ctx)
			api := secret.NewAPI(client)
			return api.CreateSecret(request)

		},
		Examples: []*core.Example{
			{
				Short: "Add a given secret",
				Raw:   `scw secret secret create name=foobar description="$(cat <path/to/your/secret>)"`,
			},
		},
	}
}

func secretFolderCreate() *core.Command {
	return &core.Command{
		Short:     `Create folder`,
		Long:      `Create folder.`,
		Namespace: "secret",
		Resource:  "folder",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(secret.CreateFolderRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ProjectIDArgSpec(),
			{
				Name:       "name",
				Short:      `Name of the folder`,
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
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*secret.CreateFolderRequest)

			client := core.ExtractClient(ctx)
			api := secret.NewAPI(client)
			return api.CreateFolder(request)

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
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
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
				Positional: false,
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
				Name:       "ephemeral-policy-template.time-to-live",
				Short:      `Time frame, from one second and up to one year, during which the secret's versions are valid.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "ephemeral-policy-template.expires-once-accessed",
				Short:      `Returns ` + "`" + `true` + "`" + ` if the version expires after a single user access.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "ephemeral-policy-template.action",
				Short:      `Action to perform when the version of a secret expires`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"unknown_action", "delete", "disable"},
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*secret.UpdateSecretRequest)

			client := core.ExtractClient(ctx)
			api := secret.NewAPI(client)
			return api.UpdateSecret(request)

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
				EnumValues: []string{"name_asc", "name_desc", "created_at_asc", "created_at_desc", "updated_at_asc", "updated_at_desc"},
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
				Name:       "is-managed",
				Short:      `Filter by managed / not managed (optional)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "path",
				Short:      `Filter by path (optional)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "is-ephemeral",
				Short:      `Filter by ephemeral / not ephemeral (optional)`,
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
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw, scw.Region(core.AllLocalities)),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
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

func secretFolderList() *core.Command {
	return &core.Command{
		Short:     `List folders`,
		Long:      `Retrieve the list of folders created within a Project.`,
		Namespace: "secret",
		Resource:  "folder",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(secret.ListFoldersRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "project-id",
				Short:      `Filter by Project ID (optional)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "path",
				Short:      `Filter by path (optional)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "order-by",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"created_at_asc", "created_at_desc", "name_asc", "name_desc"},
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw, scw.Region(core.AllLocalities)),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*secret.ListFoldersRequest)

			client := core.ExtractClient(ctx)
			api := secret.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			if request.Region == scw.Region(core.AllLocalities) {
				opts = append(opts, scw.WithRegions(api.Regions()...))
				request.Region = ""
			}
			resp, err := api.ListFolders(request, opts...)
			if err != nil {
				return nil, err
			}
			return resp.Folders, nil

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
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
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

func secretFolderDelete() *core.Command {
	return &core.Command{
		Short:     `Delete a given folder specified by the ` + "`" + `region` + "`" + ` and ` + "`" + `folder_id` + "`" + ` parameters`,
		Long:      `Delete a given folder specified by the ` + "`" + `region` + "`" + ` and ` + "`" + `folder_id` + "`" + ` parameters.`,
		Namespace: "secret",
		Resource:  "folder",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(secret.DeleteFolderRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "folder-id",
				Short:      `ID of the folder`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*secret.DeleteFolderRequest)

			client := core.ExtractClient(ctx)
			api := secret.NewAPI(client)
			e = api.DeleteFolder(request)
			if e != nil {
				return nil, e
			}
			return &core.SuccessResult{
				Resource: "folder",
				Verb:     "delete",
			}, nil
		},
	}
}

func secretSecretProtect() *core.Command {
	return &core.Command{
		Short:     `Protect a secret`,
		Long:      `Protect a given secret specified by the ` + "`" + `secret_id` + "`" + ` parameter. A protected secret can be read and modified but cannot be deleted.`,
		Namespace: "secret",
		Resource:  "secret",
		Verb:      "protect",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(secret.ProtectSecretRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "secret-id",
				Short:      `ID of the secret to protect`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*secret.ProtectSecretRequest)

			client := core.ExtractClient(ctx)
			api := secret.NewAPI(client)
			return api.ProtectSecret(request)

		},
		Examples: []*core.Example{
			{
				Short:    "Protect a secret",
				ArgsJSON: `{"secret_id":"11111111-1111-1111-1111-111111111111"}`,
			},
		},
	}
}

func secretSecretUnprotect() *core.Command {
	return &core.Command{
		Short:     `Unprotect a secret`,
		Long:      `Unprotect a given secret specified by the ` + "`" + `secret_id` + "`" + ` parameter. An unprotected secret can be read, modified and deleted.`,
		Namespace: "secret",
		Resource:  "secret",
		Verb:      "unprotect",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(secret.UnprotectSecretRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "secret-id",
				Short:      `ID of the secret to unprotect`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*secret.UnprotectSecretRequest)

			client := core.ExtractClient(ctx)
			api := secret.NewAPI(client)
			return api.UnprotectSecret(request)

		},
		Examples: []*core.Example{
			{
				Short:    "Unprotect a secret",
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
				Name:       "product-name",
				Short:      `(Deprecated: use ` + "`" + `product` + "`" + ` field) Name of the product to add`,
				Required:   false,
				Deprecated: true,
				Positional: false,
			},
			{
				Name:       "product",
				Short:      `ID of the product to add`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"unknown", "edge_services"},
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
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
				Positional: false,
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
				Name:       "password-generation.length",
				Short:      `Length of the password to generate (between 1 and 1024)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "password-generation.no-lowercase-letters",
				Short:      `Do not include lower case letters by default in the alphabet`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "password-generation.no-uppercase-letters",
				Short:      `Do not include upper case letters by default in the alphabet`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "password-generation.no-digits",
				Short:      `Do not include digits by default in the alphabet`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "password-generation.additional-chars",
				Short:      `Additional ascii characters to be included in the alphabet`,
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
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*secret.CreateSecretVersionRequest)

			client := core.ExtractClient(ctx)
			api := secret.NewAPI(client)
			return api.CreateSecretVersion(request)

		},
	}
}

func secretVersionGeneratePassword() *core.Command {
	return &core.Command{
		Short:     `Generate a password in a new version`,
		Long:      `Generate a password for the given secret specified by the ` + "`" + `region` + "`" + ` and ` + "`" + `secret_id` + "`" + ` parameters. This will also create a new version of the secret that will store the password.`,
		Namespace: "secret",
		Resource:  "version",
		Verb:      "generate-password",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(secret.GeneratePasswordRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "secret-id",
				Short:      `ID of the secret`,
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
				Name:       "disable-previous",
				Short:      `(Optional.) Disable the previous secret version`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "length",
				Short:      `Length of the password to generate (between 1 and 1024 characters)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "no-lowercase-letters",
				Short:      `(Optional.) Exclude lower case letters by default in the password character set`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "no-uppercase-letters",
				Short:      `(Optional.) Exclude upper case letters by default in the password character set`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "no-digits",
				Short:      `(Optional.) Exclude digits by default in the password character set`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "additional-chars",
				Short:      `(Optional.) Additional ASCII characters to be included in the password character set`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*secret.GeneratePasswordRequest)

			client := core.ExtractClient(ctx)
			api := secret.NewAPI(client)
			return api.GeneratePassword(request)

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
				Positional: false,
			},
			{
				Name:       "revision",
				Short:      `Version number`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
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
				Positional: false,
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
				Name:       "ephemeral-status.expires-at",
				Short:      `The version's expiration date`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "ephemeral-status.expires-once-accessed",
				Short:      `Returns ` + "`" + `true` + "`" + ` if the version expires after a single user access.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "ephemeral-status.action",
				Short:      `Action to perform when the version of a secret expires`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"unknown_action", "delete", "disable"},
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*secret.UpdateSecretVersionRequest)

			client := core.ExtractClient(ctx)
			api := secret.NewAPI(client)
			return api.UpdateSecretVersion(request)

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
				Positional: false,
			},
			{
				Name:       "status.{index}",
				Short:      `Filter results by status`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"unknown", "enabled", "disabled", "destroyed"},
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw, scw.Region(core.AllLocalities)),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
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
				Positional: false,
			},
			{
				Name:       "revision",
				Short:      `Version number`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
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
				Positional: false,
			},
			{
				Name:       "revision",
				Short:      `Version number`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*secret.DisableSecretVersionRequest)

			client := core.ExtractClient(ctx)
			api := secret.NewAPI(client)
			return api.DisableSecretVersion(request)

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
				Positional: false,
			},
			{
				Name:       "revision",
				Short:      `Version number`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*secret.AccessSecretVersionRequest)

			client := core.ExtractClient(ctx)
			api := secret.NewAPI(client)
			return api.AccessSecretVersion(request)

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
		ArgsType: reflect.TypeOf(secret.DestroySecretVersionRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "secret-id",
				Short:      `ID of the secret`,
				Required:   true,
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
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*secret.DestroySecretVersionRequest)

			client := core.ExtractClient(ctx)
			api := secret.NewAPI(client)
			return api.DestroySecretVersion(request)

		},
		Examples: []*core.Example{
			{
				Short:    "Delete a given Secret Version",
				ArgsJSON: `{"revision":"1","secret_id":"11111111-1111-1111-1111-111111111111"}`,
			},
		},
	}
}

func secretTagList() *core.Command {
	return &core.Command{
		Short:     `List tags`,
		Long:      `List all tags associated with secrets within a given Project.`,
		Namespace: "secret",
		Resource:  "tag",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(secret.ListTagsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "project-id",
				Short:      `ID of the Project to target`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw, scw.Region(core.AllLocalities)),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*secret.ListTagsRequest)

			client := core.ExtractClient(ctx)
			api := secret.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			if request.Region == scw.Region(core.AllLocalities) {
				opts = append(opts, scw.WithRegions(api.Regions()...))
				request.Region = ""
			}
			resp, err := api.ListTags(request, opts...)
			if err != nil {
				return nil, err
			}
			return resp.Tags, nil

		},
	}
}
