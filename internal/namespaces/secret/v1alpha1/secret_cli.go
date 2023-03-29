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
		secretSecret(),
		secretVersion(),
		secretSecretCreate(),
		secretSecretGet(),
		secretSecretUpdate(),
		secretSecretList(),
		secretSecretDelete(),
		secretVersionCreate(),
		secretVersionGet(),
		secretVersionUpdate(),
		secretVersionList(),
		secretVersionDelete(),
		secretVersionEnable(),
		secretVersionDisable(),
		secretVersionAccess(),
	)
}
func secretRoot() *core.Command {
	return &core.Command{
		Short:     `This API allows you to conveniently store, access and share sensitive data`,
		Long:      `Secret Manager API documentation.`,
		Namespace: "secret",
	}
}

func secretSecret() *core.Command {
	return &core.Command{
		Short: `Secret management commands`,
		Long: `Secrets are logical containers made up of zero or more immutable versions, that contain sensitive data.
`,
		Namespace: "secret",
		Resource:  "secret",
	}
}

func secretVersion() *core.Command {
	return &core.Command{
		Short: `Secret Version management commands`,
		Long: `Versions store the sensitive data contained in your secrets (API keys, passwords, or certificates).
`,
		Namespace: "secret",
		Resource:  "version",
	}
}

func secretSecretCreate() *core.Command {
	return &core.Command{
		Short:     `Create a secret`,
		Long:      `You must sepcify the ` + "`" + `region` + "`" + ` to create a secret.`,
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
			core.RegionArgSpec(scw.RegionFrPar),
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

func secretSecretGet() *core.Command {
	return &core.Command{
		Short:     `Get metadata using the secret's name`,
		Long:      `Retrieve the metadata of a secret specified by the ` + "`" + `region` + "`" + ` and the ` + "`" + `secret_name` + "`" + ` parameters.`,
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
			core.RegionArgSpec(scw.RegionFrPar),
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
		Long:      `Edit a secret's metadata such as name, tag(s) and description. The secret to update is specified by the ` + "`" + `secret_id` + "`" + ` and ` + "`" + `region` + "`" + ` parameters.`,
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
			core.RegionArgSpec(scw.RegionFrPar),
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
				Name:       "name",
				Short:      `Filter by secret name (optional)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `List of tags to filter on (optional)`,
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
				Name:       "organization-id",
				Short:      `Filter by Organization ID (optional)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.Region(core.AllLocalities)),
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
			core.RegionArgSpec(scw.RegionFrPar),
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
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
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
				Positional: false,
			},
			{
				Name:       "revision",
				Short:      `Version number. The first version of the secret is numbered 1, and all subsequent revisions augment by 1. Value can be a number or "latest"`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
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
				Short:      `Version number. The first version of the secret is numbered 1, and all subsequent revisions augment by 1. Value can be a number or "latest"`,
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
			core.RegionArgSpec(scw.RegionFrPar),
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
			core.RegionArgSpec(scw.RegionFrPar, scw.Region(core.AllLocalities)),
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
				Short:      `Version number. The first version of the secret is numbered 1, and all subsequent revisions augment by 1. Value can be a number or "latest"`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
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
				Short:      `Version number. The first version of the secret is numbered 1, and all subsequent revisions augment by 1. Value can be a number or "latest"`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
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
				Short:      `Version number. The first version of the secret is numbered 1, and all subsequent revisions augment by 1. Value can be a number or "latest"`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
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
				Short:      `Version number. The first version of the secret is numbered 1, and all subsequent revisions augment by 1. Value can be a number or "latest"`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*secret.AccessSecretVersionRequest)

			client := core.ExtractClient(ctx)
			api := secret.NewAPI(client)
			return api.AccessSecretVersion(request)

		},
	}
}
