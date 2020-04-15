package account

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/scaleway/scaleway-cli/internal/core"
	account "github.com/scaleway/scaleway-sdk-go/api/account/v2alpha1"
)

func Test_initCommand(t *testing.T) {
	tmpDir := os.TempDir()
	t.Run("simple", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		BeforeFunc: func(ctx *core.BeforeFuncCtx) error {
			pathToPublicKey := path.Join(tmpDir, ".ssh", "id_rsa.pub")
			_, err := os.Stat(pathToPublicKey)
			if err != nil {
				content := "ssh-rsa AAAAB3NzaC1yc2EAAAEDAQABAAABAQC/wF8Q8LjEexuWDc8TfKmWVZ1CiiHK6KvO0E/Rk9+d6ssqbrvtbJWRsJXZFC8+DGWVM0UFFicmOfTwEjDWzuQPFkYhmpXrD1UiLx9Viku1g1qEJgcsyH2uAwwW3OnsH1W44D6Ni/zOzMButFeKZgPeD8H9YNkpbZBZ9QrKFiAhvEyJDYSY0bsbH1/qR5DE+dLNuGlJ/g3kUMVaXSI6dHNcBHTbK0Mse23Uopk2U3BSpvX9JdbcLaYtHDOytwd16rNYui7el3uOmlR8oUpAXkeKQxBPoxgy3qI/P8/l44L9RFpbklkmdiw2ph2ymiSkRSYCWdvEVIK/A+0D8VFjGXOb email"
				err := os.MkdirAll(path.Join(tmpDir, ".ssh"), 0755)
				if err != nil {
					return err
				}
				err = ioutil.WriteFile(pathToPublicKey, []byte(content), 0644)
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

func Test_SSHKeyAddCommand(t *testing.T) {
	key := `ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBBieay3nO9wViPkuvFVgGGaA1IRlkFrr946yqvg9LxZIRhsnZ61yLCPmIOhvUAZ/gTxZGmhgtMDxkenSUTsG3F0= foobar@foobar`
	t.Run("simple", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Args: []string{
			"scw", "account", "ssh-key", "add", "name=foobar", "public-key=" + key,
		},
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: func(ctx *core.AfterFuncCtx) error {
			api := account.NewAPI(ctx.Client)
			return api.DeleteSSHKey(&account.DeleteSSHKeyRequest{
				SSHKeyID: ctx.CmdResult.(*account.SSHKey).ID,
			})
		},
	}))
}

func Test_SSHKeyRemoveCommand(t *testing.T) {
	key := `ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBGh9rvkJKMu5ljnevB4oRu4i/EnxGS734/UJ6fSDvXGIvT08jIglahc7tM5dvo02abPVXsbiazO25avZZtL6fjo= foobar@foobar`
	t.Run("simple", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		BeforeFunc: addSSHKey("Key", key),
		Cmd:        "scw account ssh-key remove {{ .Key.ID }}",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
	}))
}
