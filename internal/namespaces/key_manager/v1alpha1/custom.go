package key_manager

import (
	"context"
	"encoding/base64"

	"github.com/scaleway/scaleway-cli/v2/core"
	key_manager "github.com/scaleway/scaleway-sdk-go/api/key_manager/v1alpha1"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("keymanager").Groups = []string{"security"}

	cmds.MustFind("keymanager", "key", "decrypt").Override(cipherDecrypt)
	cmds.MustFind("keymanager", "key", "encrypt").Override(plaintextEncrypt)

	return cmds
}

func plaintextEncrypt(c *core.Command) *core.Command {
	*c.ArgSpecs.GetByName("plaintext") = core.ArgSpec{
		Name:        "plaintext",
		Short:       "Base64 Plaintext data to encrypt",
		Required:    true,
		CanLoadFile: true,
	}

	c.Interceptor = func(ctx context.Context, argsI any, runner core.CommandRunner) (any, error) {
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
		Short:       "Base64 Ciphertext data to decrypt",
		Required:    true,
		CanLoadFile: true,
	}

	c.Interceptor = func(ctx context.Context, argsI any, runner core.CommandRunner) (any, error) {
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
