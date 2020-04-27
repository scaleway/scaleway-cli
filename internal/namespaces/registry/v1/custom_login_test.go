package registry

import (
	"io/ioutil"
	"os/exec"
	"strings"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/stretchr/testify/require"
)

func Test_Login(t *testing.T) {
	clientOpts := []scw.ClientOption{
		scw.WithDefaultZone(scw.ZoneFrPar1),
		scw.WithDefaultRegion(scw.RegionFrPar),
		scw.WithAuth("SCWXXXXXXXXXXXXXXXXX", "11111111-1111-1111-1111-111111111111"),
	}

	client, err := scw.NewClient(clientOpts...)
	require.NoError(t, err)

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
		Client: client,
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
		Client: client,
	}))
}
