// This file was automatically generated. DO NOT EDIT.
// If you have any remark or suggestion do not hesitate to open an issue.

package key_manager

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	key_manager "github.com/scaleway/scaleway-sdk-go/api/key_manager/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// always import dependencies
var (
	_ = scw.RegionFrPar
)

func GetGeneratedCommands() *core.Commands {
	return core.NewCommands(
		keymanagerRoot(),
		keymanagerKey(),
		keymanagerKeyCreate(),
		keymanagerKeyGet(),
		keymanagerKeyUpdate(),
		keymanagerKeyDelete(),
		keymanagerKeyRotate(),
		keymanagerKeyProtect(),
		keymanagerKeyUnprotect(),
		keymanagerKeyEnable(),
		keymanagerKeyDisable(),
		keymanagerKeyList(),
		keymanagerKeyGenerateDataKey(),
		keymanagerKeyEncrypt(),
		keymanagerKeyDecrypt(),
	)
}
func keymanagerRoot() *core.Command {
	return &core.Command{
		Short:     `Key Manager API`,
		Long:      `This API allows you to conveniently store and use cryptographic keys.`,
		Namespace: "keymanager",
	}
}

func keymanagerKey() *core.Command {
	return &core.Command{
		Short:     `Key management commands`,
		Long:      `Keys are logical containers which store cryptographic keys.`,
		Namespace: "keymanager",
		Resource:  "key",
	}
}

func keymanagerKeyCreate() *core.Command {
	return &core.Command{
		Short:     `Create a key`,
		Long:      `Create a key in a given region specified by the ` + "`" + `region` + "`" + ` parameter. Keys only support symmetric encryption. You can use keys to encrypt or decrypt arbitrary payloads, or to generate data encryption keys that can be used without being stored in Key Manager.`,
		Namespace: "keymanager",
		Resource:  "key",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(key_manager.CreateKeyRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ProjectIDArgSpec(),
			{
				Name:       "name",
				Short:      `(Optional) Name of the key`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "usage.symmetric-encryption",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"unknown_symmetric_encryption", "aes_256_gcm"},
			},
			{
				Name:       "description",
				Short:      `(Optional) Description of the key`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `(Optional) List of the key's tags`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "rotation-policy.rotation-period",
				Short:      `Rotation period`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "rotation-policy.next-rotation-at",
				Short:      `Key next rotation date`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "unprotected",
				Short:      `(Optional) Defines whether key protection is applied to a key. Protected keys can be used but not deleted`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*key_manager.CreateKeyRequest)

			client := core.ExtractClient(ctx)
			api := key_manager.NewAPI(client)
			return api.CreateKey(request)

		},
	}
}

func keymanagerKeyGet() *core.Command {
	return &core.Command{
		Short:     `Get key metadata`,
		Long:      `Retrieve the metadata of a key specified by the ` + "`" + `region` + "`" + ` and ` + "`" + `key_id` + "`" + ` parameters.`,
		Namespace: "keymanager",
		Resource:  "key",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(key_manager.GetKeyRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "key-id",
				Short:      `ID of the key to target`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*key_manager.GetKeyRequest)

			client := core.ExtractClient(ctx)
			api := key_manager.NewAPI(client)
			return api.GetKey(request)

		},
	}
}

func keymanagerKeyUpdate() *core.Command {
	return &core.Command{
		Short:     `Update a key`,
		Long:      `Update a key's metadata (name, description and tags), specified by the ` + "`" + `key_id` + "`" + ` and ` + "`" + `region` + "`" + ` parameters.`,
		Namespace: "keymanager",
		Resource:  "key",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(key_manager.UpdateKeyRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "key-id",
				Short:      `ID of the key to update`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "name",
				Short:      `(Optional) Updated name of the key`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "description",
				Short:      `(Optional) Updated description of the key`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `(Optional) Updated list of the key's tags`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "rotation-policy.rotation-period",
				Short:      `Rotation period`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "rotation-policy.next-rotation-at",
				Short:      `Key next rotation date`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*key_manager.UpdateKeyRequest)

			client := core.ExtractClient(ctx)
			api := key_manager.NewAPI(client)
			return api.UpdateKey(request)

		},
	}
}

func keymanagerKeyDelete() *core.Command {
	return &core.Command{
		Short:     `Delete a key`,
		Long:      `Delete an existing key specified by the ` + "`" + `region` + "`" + ` and ` + "`" + `key_id` + "`" + ` parameters. Deleting a key is permanent and cannot be undone. All data encrypted using this key, including data encryption keys, will become unusable.`,
		Namespace: "keymanager",
		Resource:  "key",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(key_manager.DeleteKeyRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "key-id",
				Short:      `ID of the key to delete`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*key_manager.DeleteKeyRequest)

			client := core.ExtractClient(ctx)
			api := key_manager.NewAPI(client)
			e = api.DeleteKey(request)
			if e != nil {
				return nil, e
			}
			return &core.SuccessResult{
				Resource: "key",
				Verb:     "delete",
			}, nil
		},
	}
}

func keymanagerKeyRotate() *core.Command {
	return &core.Command{
		Short:     `Rotate a key`,
		Long:      `Generate a new version of an existing key with randomly generated key material. Rotated keys can still be used to decrypt previously encrypted data. The key's new material will be used for subsequent encryption operations and data key generation.`,
		Namespace: "keymanager",
		Resource:  "key",
		Verb:      "rotate",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(key_manager.RotateKeyRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "key-id",
				Short:      `ID of the key to rotate`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*key_manager.RotateKeyRequest)

			client := core.ExtractClient(ctx)
			api := key_manager.NewAPI(client)
			return api.RotateKey(request)

		},
	}
}

func keymanagerKeyProtect() *core.Command {
	return &core.Command{
		Short:     `Apply key protection`,
		Long:      `Apply key protection to a given key specified by the ` + "`" + `key_id` + "`" + ` parameter. Applying key protection means that your key can be used and modified, but it cannot be deleted.`,
		Namespace: "keymanager",
		Resource:  "key",
		Verb:      "protect",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(key_manager.ProtectKeyRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "key-id",
				Short:      `ID of the key to apply key protection to`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*key_manager.ProtectKeyRequest)

			client := core.ExtractClient(ctx)
			api := key_manager.NewAPI(client)
			return api.ProtectKey(request)

		},
	}
}

func keymanagerKeyUnprotect() *core.Command {
	return &core.Command{
		Short:     `Remove key protection`,
		Long:      `Remove key protection from a given key specified by the ` + "`" + `key_id` + "`" + ` parameter. Removing key protection means that your key can be deleted anytime.`,
		Namespace: "keymanager",
		Resource:  "key",
		Verb:      "unprotect",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(key_manager.UnprotectKeyRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "key-id",
				Short:      `ID of the key to remove key protection from`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*key_manager.UnprotectKeyRequest)

			client := core.ExtractClient(ctx)
			api := key_manager.NewAPI(client)
			return api.UnprotectKey(request)

		},
	}
}

func keymanagerKeyEnable() *core.Command {
	return &core.Command{
		Short:     `Enable key`,
		Long:      `Enable a given key to be used for cryptographic operations. Enabling a key allows you to make a disabled key usable again. You must specify the ` + "`" + `region` + "`" + ` and ` + "`" + `key_id` + "`" + ` parameters.`,
		Namespace: "keymanager",
		Resource:  "key",
		Verb:      "enable",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(key_manager.EnableKeyRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "key-id",
				Short:      `ID of the key to enable`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*key_manager.EnableKeyRequest)

			client := core.ExtractClient(ctx)
			api := key_manager.NewAPI(client)
			return api.EnableKey(request)

		},
	}
}

func keymanagerKeyDisable() *core.Command {
	return &core.Command{
		Short:     `Disable key`,
		Long:      `Disable a given key to be used for cryptographic operations. Disabling a key renders it unusable. You must specify the ` + "`" + `region` + "`" + ` and ` + "`" + `key_id` + "`" + ` parameters.`,
		Namespace: "keymanager",
		Resource:  "key",
		Verb:      "disable",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(key_manager.DisableKeyRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "key-id",
				Short:      `ID of the key to disable`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*key_manager.DisableKeyRequest)

			client := core.ExtractClient(ctx)
			api := key_manager.NewAPI(client)
			return api.DisableKey(request)

		},
	}
}

func keymanagerKeyList() *core.Command {
	return &core.Command{
		Short:     `List keys`,
		Long:      `Retrieve the list of keys created within all Projects of an Organization or in a given Project. You must specify the ` + "`" + `region` + "`" + `, and either the ` + "`" + `organization_id` + "`" + ` or the ` + "`" + `project_id` + "`" + `.`,
		Namespace: "keymanager",
		Resource:  "key",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(key_manager.ListKeysRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "project-id",
				Short:      `(Optional) Filter by Project ID`,
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
				Short:      `(Optional) List of tags to filter on`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "name",
				Short:      `(Optional) Filter by key name`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "organization-id",
				Short:      `(Optional) Filter by Organization ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw, scw.Region(core.AllLocalities)),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*key_manager.ListKeysRequest)

			client := core.ExtractClient(ctx)
			api := key_manager.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			if request.Region == scw.Region(core.AllLocalities) {
				opts = append(opts, scw.WithRegions(api.Regions()...))
				request.Region = ""
			}
			resp, err := api.ListKeys(request, opts...)
			if err != nil {
				return nil, err
			}
			return resp.Keys, nil

		},
	}
}

func keymanagerKeyGenerateDataKey() *core.Command {
	return &core.Command{
		Short: `Generate a data encryption key`,
		Long: `Generate a new data encryption key to use for cryptographic operations outside of Key Manager. Note that Key Manager does not store your data encryption key. The data encryption key is encrypted and must be decrypted using the key you have created in Key Manager. The data encryption key's plaintext is returned in the response object, for immediate usage.

Always store the data encryption key's ciphertext, rather than its plaintext, which must not be stored. To retrieve your key's plaintext, call the Decrypt endpoint with your key's ID and ciphertext.`,
		Namespace: "keymanager",
		Resource:  "key",
		Verb:      "generate-data-key",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(key_manager.GenerateDataKeyRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "key-id",
				Short:      `ID of the key`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "algorithm",
				Short:      `Symmetric encryption algorithm of the data encryption key`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"unknown_symmetric_encryption", "aes_256_gcm"},
			},
			{
				Name:       "without-plaintext",
				Short:      `(Optional) Defines whether to return the data encryption key's plaintext in the response object`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*key_manager.GenerateDataKeyRequest)

			client := core.ExtractClient(ctx)
			api := key_manager.NewAPI(client)
			return api.GenerateDataKey(request)

		},
	}
}

func keymanagerKeyEncrypt() *core.Command {
	return &core.Command{
		Short:     `Encrypt data`,
		Long:      `Encrypt data using an existing key, specified by the ` + "`" + `key_id` + "`" + ` parameter. Only keys with a usage set to **symmetric_encryption** are supported by this method. The maximum payload size that can be encrypted is 64KB of plaintext.`,
		Namespace: "keymanager",
		Resource:  "key",
		Verb:      "encrypt",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(key_manager.EncryptRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "key-id",
				Short:      `ID of the key to encrypt`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "plaintext",
				Short:      `Base64 Plaintext data to encrypt`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "associated-data",
				Short:      `(Optional) Additional authenticated data`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*key_manager.EncryptRequest)

			client := core.ExtractClient(ctx)
			api := key_manager.NewAPI(client)
			return api.Encrypt(request)

		},
	}
}

func keymanagerKeyDecrypt() *core.Command {
	return &core.Command{
		Short:     `Decrypt data`,
		Long:      `Decrypt data using an existing key, specified by the ` + "`" + `key_id` + "`" + ` parameter. The maximum payload size that can be decrypted is the result of the encryption of 64KB of data (around 131KB).`,
		Namespace: "keymanager",
		Resource:  "key",
		Verb:      "decrypt",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(key_manager.DecryptRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "key-id",
				Short:      `ID of the key to decrypt`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "ciphertext",
				Short:      `Base64 Ciphertext data to decrypt`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "associated-data",
				Short:      `(Optional) Additional authenticated data`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*key_manager.DecryptRequest)

			client := core.ExtractClient(ctx)
			api := key_manager.NewAPI(client)
			return api.Decrypt(request)

		},
	}
}
