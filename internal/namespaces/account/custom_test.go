package account

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/scaleway/scaleway-cli/internal/core"
)

func Test_initCommand(t *testing.T) {
	t.Run("simple", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		BeforeFunc: func(ctx *core.BeforeFuncCtx) error {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return err
			}
			path_ := path.Join(homeDir, "/.ssh/id_rsa.pub")
			if _, err := os.Stat(path_); os.IsNotExist(err) {
				content := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC/wF8Q8LjEexuWDc8TfKmWVZ1CiiHK6KvO0E/Rk9+d6ssqbrvtbJWRsJXZFC8+DGWVM0UFFicmOfTwEjDWzuQPFkYhmpXrD1UiLx9Viku1g1qEJgcsyH2uAwwW3OnsH1W44D6Ni/zOzMButFeKZgPeD8H9YNkpbZBZ9QrKFiAhvEyJDYSY0bsbH1/qR5DE+dLNuGlJ/g3kUMVaXSI6dHNcBHTbK0Mse23Uopk2U3BSpvX9JdbcLaYtHDOytwd16rNYui7el3uOmlR8oUpAXkeKQxBPoxgy3qI/P8/l44L9RFpbklkmdiw2ph2ymiSkRSYCWdvEVIK/A+0D8VFjGXOb email"
				_ = os.MkdirAll(path.Join(homeDir, "/.ssh"), 0644)
				err = ioutil.WriteFile(path_, []byte(content), 0644)
				return err
			}
			return nil
		},
		Cmd: `scw account ssh-key init with-ssh-key=true`,
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
	}))
}
