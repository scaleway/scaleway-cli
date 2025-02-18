package iam_test

import (
	"os"
	"path"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	iam "github.com/scaleway/scaleway-cli/v2/internal/namespaces/iam/v1alpha1"
	"github.com/scaleway/scaleway-cli/v2/internal/testhelpers"
	iamsdk "github.com/scaleway/scaleway-sdk-go/api/iam/v1alpha1"
)

func Test_initWithSSHKeyCommand(t *testing.T) {
	tmpDir := os.TempDir()
	key := `ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIBn9mGL7LGZ6/RTIVP7GExiD5gOwgl63MbJGlL7a6U3x foo@foobar.com`
	t.Run("simple", core.Test(&core.TestConfig{
		Commands: iam.GetCommands(),
		BeforeFunc: func(_ *core.BeforeFuncCtx) error {
			pathToPublicKey := path.Join(tmpDir, ".ssh", "id_ed25519.pub")
			_, err := os.Stat(pathToPublicKey)
			if err != nil {
				err := os.MkdirAll(path.Join(tmpDir, ".ssh"), 0o755)
				if err != nil {
					return err
				}
				err = os.WriteFile(pathToPublicKey, []byte(key), 0o644)

				return err
			}

			return err
		},
		Cmd: `scw iam ssh-key init with-ssh-key=true`,
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		OverrideEnv: map[string]string{
			"HOME": tmpDir,
		},
	}))
}

func Test_SSHKeyCreateCommand(t *testing.T) {
	key := `ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBBieay3nO9wViPkuvFVgGGaA1IRlkFrr946yqvg9LxZIRhsnZ61yLCPmIOhvUAZ/gTxZGmhgtMDxkenSUTsG3F0= foobar@foobar`
	t.Run("simple", core.Test(&core.TestConfig{
		Commands: iam.GetCommands(),
		Args: []string{
			"scw", "iam", "ssh-key", "create", "name=foobar", "public-key=" + key,
		},
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: func(ctx *core.AfterFuncCtx) error {
			api := iamsdk.NewAPI(ctx.Client)
			key := testhelpers.Value[*iamsdk.SSHKey](t, ctx.CmdResult)

			return api.DeleteSSHKey(&iamsdk.DeleteSSHKeyRequest{
				SSHKeyID: key.ID,
			})
		},
	}))
}

func Test_SSHKeyRemoveCommand(t *testing.T) {
	key := `ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBGh9rvkJKMu5ljnevB4oRu4i/EnxGS734/UJ6fSDvXGIvT08jIglahc7tM5dvo02abPVXsbiazO25avZZtL6fjo= foobar@foobar`
	t.Run("simple", core.Test(&core.TestConfig{
		Commands:   iam.GetCommands(),
		BeforeFunc: addSSHKey("Key", key),
		Cmd:        "scw iam ssh-key delete {{ .Key.ID }}",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
	}))
}
