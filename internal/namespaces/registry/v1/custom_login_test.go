package registry_test

import (
	"io"
	"os/exec"
	"strings"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/registry/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Login(t *testing.T) {
	t.Run("docker", core.Test(&core.TestConfig{
		Commands: registry.GetCommands(),
		Cmd:      "scw registry login program=docker",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		OverrideExec: func(ctx *core.ExecFuncCtx, cmd *exec.Cmd) (exitCode int, err error) {
			assert.Equal(
				t,
				"docker login -u scaleway --password-stdin rg.fr-par.scw.cloud",
				strings.Join(cmd.Args, " "),
			)
			stdin, err := io.ReadAll(cmd.Stdin)
			secret, _ := ctx.Client.GetSecretKey()
			require.NoError(t, err)
			assert.Equal(t, secret, string(stdin))

			return 0, nil
		},
	}))
	t.Run("podman", core.Test(&core.TestConfig{
		Commands: registry.GetCommands(),
		Cmd:      "scw registry login program=podman",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		OverrideExec: func(ctx *core.ExecFuncCtx, cmd *exec.Cmd) (exitCode int, err error) {
			assert.Equal(
				t,
				"podman login -u scaleway --password-stdin rg.fr-par.scw.cloud",
				strings.Join(cmd.Args, " "),
			)
			stdin, err := io.ReadAll(cmd.Stdin)
			secret, _ := ctx.Client.GetSecretKey()
			require.NoError(t, err)
			assert.Equal(t, secret, string(stdin))

			return 0, nil
		},
	}))
}
