package registry

import (
	"io/ioutil"
	"os/exec"
	"strings"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/stretchr/testify/require"
)

func Test_Login(t *testing.T) {
	t.Run("docker", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw registry login program=docker",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		OverrideExec: func(ctx *core.ExecFuncCtx, cmd *exec.Cmd) (exitCode int, err error) {
			assert.Equal(t, "docker login -u scaleway --password-stdin rg.fr-par.scw.cloud", strings.Join(cmd.Args, " "))
			stdin, err := ioutil.ReadAll(cmd.Stdin)
			secret, _ := ctx.Client.GetSecretKey()
			require.NoError(t, err)
			assert.Equal(t, secret, string(stdin))
			return 0, nil
		},
	}))
	t.Run("podman", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw registry login program=podman",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		OverrideExec: func(ctx *core.ExecFuncCtx, cmd *exec.Cmd) (exitCode int, err error) {
			assert.Equal(t, "podman login -u scaleway --password-stdin rg.fr-par.scw.cloud", strings.Join(cmd.Args, " "))
			stdin, err := ioutil.ReadAll(cmd.Stdin)
			secret, _ := ctx.Client.GetSecretKey()
			require.NoError(t, err)
			assert.Equal(t, secret, string(stdin))
			return 0, nil
		},
	}))
}
