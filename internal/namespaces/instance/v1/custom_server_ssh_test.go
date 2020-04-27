package instance

import (
	"fmt"
	"os/exec"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/scaleway/scaleway-cli/internal/core"
)

func Test_ServerSSH(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			createServer("Server"),
			startServer("Server"),
		),
		Cmd: "scw instance server ssh {{ .Server.ID }}",
		OverrideExec: func(ctx *core.ExecFuncCtx, cmd *exec.Cmd) (exitCode int, err error) {
			assert.Equal(ctx.T, ctx.Meta.Tpl("/usr/bin/ssh {{ .Server.PublicIP.Address }} -p 22 -l root -t"), cmd.String())
			_, err = fmt.Fprintln(cmd.Stdout, "This is what SSH command print on stdout")
			assert.NoError(ctx.T, err)
			return 0, nil
		},
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc:       deleteServer("Server"),
		DisableParallel: true,
	}))
}
