package account

import (
	"os"
	"path"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/internal/core"
)

func Test_initCommand(t *testing.T) {
	tmpDir := os.TempDir()
	t.Run("simple", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		BeforeFunc: func(ctx *core.BeforeFuncCtx) error {
			pathToPublicKey := path.Join(tmpDir, ".ssh", "id_ed25519.pub")
			_, err := os.Stat(pathToPublicKey)
			if err != nil {
				content := "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIBn9mGL7LGZ6/RTIVP7GExiD5gOwgl63MbJGlL7a6U3x foo@foobar.com"
				err := os.MkdirAll(path.Join(tmpDir, ".ssh"), 0755)
				if err != nil {
					return err
				}
				err = os.WriteFile(pathToPublicKey, []byte(content), 0644)
				return err
			}
			return err
		},
		Cmd: `scw account ssh-key init with-ssh-key=true`,
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		OverrideEnv: map[string]string{
			"HOME": tmpDir,
		},
	}))
}
