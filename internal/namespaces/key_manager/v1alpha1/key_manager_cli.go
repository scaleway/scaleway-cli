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
		keymanagerKeyImportKeyMaterial(),
		keymanagerKeyDeleteKeyMaterial(),
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
		Long:      `Create a key in a given region specified by the ` + "`" + `region` + "`" + ` parameter. Keys only support symmetric encryption. You can use keys to encrypt or decrypt arbitrary payloads, or to generate data encryption keys. **Data encryption keys are not stored in Key Manager**.`,
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
				Short:      `Algorithm used to encrypt and decrypt arbitrary payloads.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_symmetric_encryption",
					"aes_256_gcm",
				},
			},
			{
				Name:       "usage.asymmetric-encryption",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_asymmetric_encryption",
					"rsa_oaep_2048_sha256",
					"rsa_oaep_3072_sha256",
					"rsa_oaep_4096_sha256",
				},
			},
			{
				Name:       "usage.asymmetric-signing",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_asymmetric_signing",
					"ec_p256_sha256",
					"ec_p384_sha384",
					"rsa_pss_2048_sha256",
					"rsa_pss_3072_sha256",
					"rsa_pss_4096_sha256",
					"rsa_pkcs1_2048_sha256",
					"rsa_pkcs1_3072_sha256",
					"rsa_pkcs1_4096_sha256",
				},
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
			{
				Name:       "origin",
				Short:      `Key origin`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_origin",
					"scaleway_kms",
					"external",
				},
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
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
		Long:      `Retrieve metadata for a specified key using the ` + "`" + `region` + "`" + ` and ` + "`" + `key_id` + "`" + ` parameters.`,
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
				Positional: true,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
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
		Long:      `Modify a key's metadata including name, description and tags, specified by the ` + "`" + `key_id` + "`" + ` and ` + "`" + `region` + "`" + ` parameters.`,
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
				Positional: true,
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
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
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
		Long:      `Permanently delete a key specified by the ` + "`" + `region` + "`" + ` and ` + "`" + `key_id` + "`" + ` parameters. This action is irreversible. Any data encrypted with this key, including data encryption keys, will no longer be decipherable.`,
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
				Positional: true,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
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
		Long:      `Generate a new version of an existing key with new key material. Previous key versions remain usable to decrypt previously encrypted data, but the key's new version will be used for subsequent encryption operations and data key generation.`,
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
				Positional: true,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
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
		Long:      `Apply protection to a given key specified by the ` + "`" + `key_id` + "`" + ` parameter. Applying key protection means that your key can be used and modified, but it cannot be deleted.`,
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
				Positional: true,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
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
				Positional: true,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
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
				Positional: true,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
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
		Long:      `Disable a given key, preventing it to be used for cryptographic operations. Disabling a key renders it unusable. You must specify the ` + "`" + `region` + "`" + ` and ` + "`" + `key_id` + "`" + ` parameters.`,
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
				Positional: true,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
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
		Long:      `Retrieve a list of keys across all Projects in an Organization or within a specific Project. You must specify the ` + "`" + `region` + "`" + `, and either the ` + "`" + `organization_id` + "`" + ` or the ` + "`" + `project_id` + "`" + `.`,
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
				Name:       "usage",
				Short:      `(Optional) Filter keys by usage.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_usage",
					"symmetric_encryption",
					"asymmetric_encryption",
					"asymmetric_signing",
				},
			},
			{
				Name:       "scheduled-for-deletion",
				Short:      `Filter keys based on their deletion status. By default, only keys not scheduled for deletion are returned in the output.`,
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
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
				scw.Region(core.AllLocalities),
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
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
		Short: `Create a data encryption key`,
		Long: `Create a new data encryption key for cryptographic operations outside of Key Manager. The data encryption key is encrypted and must be decrypted using the key you have created in Key Manager.

The data encryption key is returned in plaintext and ciphertext but it should only be stored in its encrypted form (ciphertext). Key Manager does not store your data encryption key. To retrieve your key's plaintext, use the ` + "`" + `Decrypt` + "`" + ` method with your key's ID and ciphertext.`,
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
				Positional: true,
			},
			{
				Name:       "algorithm",
				Short:      `Algorithm with which the data encryption key will be used to encrypt and decrypt arbitrary payloads`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				Default:    core.DefaultValueSetter("aes_256_gcm"),
				EnumValues: []string{
					"unknown_symmetric_encryption",
					"aes_256_gcm",
				},
			},
			{
				Name:       "without-plaintext",
				Short:      `(Optional) Defines whether to return the data encryption key's plaintext in the response object`,
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
			request := args.(*key_manager.GenerateDataKeyRequest)

			client := core.ExtractClient(ctx)
			api := key_manager.NewAPI(client)

			return api.GenerateDataKey(request)
		},
	}
}

func keymanagerKeyEncrypt() *core.Command {
	return &core.Command{
		Short:     `Encrypt a payload`,
		Long:      `Encrypt a payload using an existing key, specified by the ` + "`" + `key_id` + "`" + ` parameter. The maximum payload size that can be encrypted is 64 KB of plaintext.`,
		Namespace: "keymanager",
		Resource:  "key",
		Verb:      "encrypt",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(key_manager.EncryptRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "key-id",
				Short:      `ID of the key to use for encryption`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "plaintext",
				Short:      `Plaintext data to encrypt`,
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
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*key_manager.EncryptRequest)

			client := core.ExtractClient(ctx)
			api := key_manager.NewAPI(client)

			return api.Encrypt(request)
		},
	}
}

func keymanagerKeyDecrypt() *core.Command {
	return &core.Command{
		Short:     `Decrypt an encrypted payload`,
		Long:      `Decrypt an encrypted payload using an existing key, specified by the ` + "`" + `key_id` + "`" + ` parameter. The maximum payload size that can be decrypted is equivalent to the encrypted output of 64 KB of data (around 131 KB).`,
		Namespace: "keymanager",
		Resource:  "key",
		Verb:      "decrypt",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(key_manager.DecryptRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "key-id",
				Short:      `ID of the key to decrypt with`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "ciphertext",
				Short:      `Ciphertext data to decrypt`,
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
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*key_manager.DecryptRequest)

			client := core.ExtractClient(ctx)
			api := key_manager.NewAPI(client)

			return api.Decrypt(request)
		},
	}
}

func keymanagerKeyImportKeyMaterial() *core.Command {
	return &core.Command{
		Short:     `Import key material`,
		Long:      `Import externally generated key material into Key Manager to derive a new cryptographic key. The key's origin must be ` + "`" + `external` + "`" + `.`,
		Namespace: "keymanager",
		Resource:  "key",
		Verb:      "import-key-material",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(key_manager.ImportKeyMaterialRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "key-id",
				Short:      `ID of the key in which to import key material`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "key-material",
				Short:      `The key material The key material is a random sequence of bytes used to derive a cryptographic key.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "salt",
				Short:      `(Optional) Salt value to pass the key derivation function`,
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
			request := args.(*key_manager.ImportKeyMaterialRequest)

			client := core.ExtractClient(ctx)
			api := key_manager.NewAPI(client)

			return api.ImportKeyMaterial(request)
		},
	}
}

func keymanagerKeyDeleteKeyMaterial() *core.Command {
	return &core.Command{
		Short:     `Delete key material`,
		Long:      `Delete previously imported key material. This renders the associated cryptographic key unusable for any operation. The key's origin must be ` + "`" + `external` + "`" + `.`,
		Namespace: "keymanager",
		Resource:  "key",
		Verb:      "delete-key-material",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(key_manager.DeleteKeyMaterialRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "key-id",
				Short:      `ID of the key of which to delete the key material`,
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
			request := args.(*key_manager.DeleteKeyMaterialRequest)

			client := core.ExtractClient(ctx)
			api := key_manager.NewAPI(client)
			e = api.DeleteKeyMaterial(request)
			if e != nil {
				return nil, e
			}

			return &core.SuccessResult{
				Resource: "key",
				Verb:     "delete-key-material",
			}, nil
		},
	}
}
