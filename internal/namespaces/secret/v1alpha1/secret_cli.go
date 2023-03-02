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
		Long:      `This API allows you to conveniently store, access and share sensitive data.`,
		Namespace: "secret",
	}
}

func secretSecret() *core.Command {
	return &core.Command{
		Short: `Secret management commands`,
		Long: `Logical container made up of zero or more immutable versions, that hold the sensitive data.
`,
		Namespace: "secret",
		Resource:  "secret",
	}
}

func secretVersion() *core.Command {
	return &core.Command{
		Short: `Secret Version management commands`,
		Long: `Immutable version of a secret.
`,
		Namespace: "secret",
		Resource:  "version",
	}
}

func secretSecretCreate() *core.Command {
	return &core.Command{
		Short:     `Create a Secret containing no versions`,
		Long:      `Create a Secret containing no versions.`,
		Namespace: "secret",
		Resource:  "secret",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(secret.CreateSecretRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ProjectIDArgSpec(),
			{
				Name:       "name",
				Short:      `Name of the Secret`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `List of tags associated to this Secret`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "description",
				Short:      `Description of the Secret`,
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
		Short:     `Get metadata of a Secret`,
		Long:      `Get metadata of a Secret.`,
		Namespace: "secret",
		Resource:  "secret",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(secret.GetSecretRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "secret-id",
				Short:      `ID of the Secret`,
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
		Short:     `Update metadata of a Secret`,
		Long:      `Update metadata of a Secret.`,
		Namespace: "secret",
		Resource:  "secret",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(secret.UpdateSecretRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "secret-id",
				Short:      `ID of the Secret`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "name",
				Short:      `New name of the Secret (optional)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `New list of tags associated to this Secret (optional)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "description",
				Short:      `Description of the Secret`,
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
		Short:     `List Secrets`,
		Long:      `List Secrets.`,
		Namespace: "secret",
		Resource:  "secret",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(secret.ListSecretsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "project-id",
				Short:      `ID of a project to filter on (optional)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "name",
				Short:      `Secret name to filter on (optional)`,
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
				Short:      `ID of an organization to filter on (optional)`,
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
		}},
	}
}

func secretSecretDelete() *core.Command {
	return &core.Command{
		Short:     `Delete a secret`,
		Long:      `Delete a secret.`,
		Namespace: "secret",
		Resource:  "secret",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(secret.DeleteSecretRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "secret-id",
				Short:      `ID of the Secret`,
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
		Short:     `Create a SecretVersion`,
		Long:      `Create a SecretVersion.`,
		Namespace: "secret",
		Resource:  "version",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(secret.CreateSecretVersionRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "secret-id",
				Short:      `ID of the Secret`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "data",
				Short:      `The base64-encoded secret payload of the SecretVersion`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "description",
				Short:      `Description of the SecretVersion`,
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
		Short:     `Get metadata of a SecretVersion`,
		Long:      `Get metadata of a SecretVersion.`,
		Namespace: "secret",
		Resource:  "version",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(secret.GetSecretVersionRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "secret-id",
				Short:      `ID of the Secret`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "revision",
				Short:      `Revision of the SecretVersion (may be a number or "latest")`,
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
		Short:     `Update metadata of a SecretVersion`,
		Long:      `Update metadata of a SecretVersion.`,
		Namespace: "secret",
		Resource:  "version",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(secret.UpdateSecretVersionRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "secret-id",
				Short:      `ID of the Secret`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "revision",
				Short:      `Revision of the SecretVersion (may be a number or "latest")`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "description",
				Short:      `Description of the SecretVersion`,
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
		Short:     `List versions of a secret, not returning any sensitive data`,
		Long:      `List versions of a secret, not returning any sensitive data.`,
		Namespace: "secret",
		Resource:  "version",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(secret.ListSecretVersionsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "secret-id",
				Short:      `ID of the Secret`,
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
		View: &core.View{Fields: []*core.ViewField{
			{
				FieldName: "SecretID",
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
		}},
	}
}

func secretVersionDelete() *core.Command {
	return &core.Command{
		Short:     `Destroy a SecretVersion, permanently destroying the sensitive data`,
		Long:      `Destroy a SecretVersion, permanently destroying the sensitive data.`,
		Namespace: "secret",
		Resource:  "version",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(secret.DestroySecretVersionRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "secret-id",
				Short:      `ID of the Secret`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "revision",
				Short:      `Revision of the SecretVersion (may be a number or "latest")`,
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
		Short:     `Enable a SecretVersion`,
		Long:      `Enable a SecretVersion.`,
		Namespace: "secret",
		Resource:  "version",
		Verb:      "enable",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(secret.EnableSecretVersionRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "secret-id",
				Short:      `ID of the Secret`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "revision",
				Short:      `Revision of the SecretVersion (may be a number or "latest")`,
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
		Short:     `Disable a SecretVersion`,
		Long:      `Disable a SecretVersion.`,
		Namespace: "secret",
		Resource:  "version",
		Verb:      "disable",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(secret.DisableSecretVersionRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "secret-id",
				Short:      `ID of the Secret`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "revision",
				Short:      `Revision of the SecretVersion (may be a number or "latest")`,
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
		Short:     `Access a SecretVersion, returning the sensitive data`,
		Long:      `Access a SecretVersion, returning the sensitive data.`,
		Namespace: "secret",
		Resource:  "version",
		Verb:      "access",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(secret.AccessSecretVersionRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "secret-id",
				Short:      `ID of the Secret`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "revision",
				Short:      `Revision of the SecretVersion (may be a number or "latest")`,
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
