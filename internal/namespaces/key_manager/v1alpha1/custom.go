package key_manager

import (
	"context"
	"encoding/base64"

	"github.com/scaleway/scaleway-cli/v2/internal/core"
	key_manager "github.com/scaleway/scaleway-sdk-go/api/key_manager/v1alpha1"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("keymanager", "key", "decrypt").Override(cipherDecrypt)
	cmds.MustFind("keymanager", "key", "encrypt").Override(plaintextEncrypt)
	cmds.MustFind("keymanager", "key", "create").Override(keyCreate)
	cmds.MustFind("keymanager", "key", "generate-data-key").Override(dataKeyGenerate)

	return cmds
}

func plaintextEncrypt(c *core.Command) *core.Command {
	*c.ArgSpecs.GetByName("plaintext") = core.ArgSpec{
		Name:        "plaintext",
		Short:       "Content of the plaintext in Base64.",
		Required:    true,
		CanLoadFile: true,
	}

	c.Interceptor = func(ctx context.Context, argsI interface{}, runner core.CommandRunner) (interface{}, error) {
		args := argsI.(*key_manager.EncryptRequest)

		p, err := base64.StdEncoding.DecodeString(string(args.Plaintext))
		if err != nil {
			return nil, err
		}

		args.Plaintext = p

		return runner(ctx, args)
	}

	return c
}

func cipherDecrypt(c *core.Command) *core.Command {
	*c.ArgSpecs.GetByName("ciphertext") = core.ArgSpec{
		Name:        "ciphertext",
		Short:       "Content of the ciphertext in Base64.",
		Required:    true,
		CanLoadFile: true,
	}

	c.Interceptor = func(ctx context.Context, argsI interface{}, runner core.CommandRunner) (interface{}, error) {
		args := argsI.(*key_manager.DecryptRequest)

		c, err := base64.StdEncoding.DecodeString(string(args.Ciphertext))
		if err != nil {
			return nil, err
		}

		args.Ciphertext = c

		return runner(ctx, args)
	}

	return c
}

func keyCreate(c *core.Command) *core.Command {
	c.Interceptor = func(ctx context.Context, argsI interface{}, runner core.CommandRunner) (interface{}, error) {
		args := argsI.(*key_manager.CreateKeyRequest)

		if args.Usage == nil {
			defaultUsage := key_manager.KeyAlgorithmSymmetricEncryptionAes256Gcm
			args.Usage = &key_manager.KeyUsage{SymmetricEncryption: &defaultUsage}
		}

		return runner(ctx, args)
	}

	return c
}

func dataKeyGenerate(c *core.Command) *core.Command {
	c.Interceptor = func(ctx context.Context, argsI interface{}, runner core.CommandRunner) (interface{}, error) {
		args := argsI.(*key_manager.GenerateDataKeyRequest)

		if args.Algorithm == "" {
			args.Algorithm = key_manager.DataKeyAlgorithmSymmetricEncryptionAes256Gcm
		}

		return runner(ctx, args)
	}

	return c
}
