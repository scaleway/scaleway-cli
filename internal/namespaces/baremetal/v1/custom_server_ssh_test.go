package baremetal_test

import (
	"fmt"
	"os/exec"
	"strings"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/baremetal/v1"
	baremetalSDK "github.com/scaleway/scaleway-sdk-go/api/baremetal/v1"
	"github.com/stretchr/testify/assert"
)

// sshOverrideExec asserts the ssh command built by `baremetal server ssh`
// (IPv4 preferred, then first IP, root user on a server without install).
func sshOverrideExec(exitCode int) core.OverrideExecTestFunc {
	return func(ctx *core.ExecFuncCtx, cmd *exec.Cmd) (int, error) {
		server := ctx.Meta["Server"].(*baremetalSDK.Server)

		ip := ""
		for _, serverIP := range server.IPs {
			if serverIP.Version == baremetalSDK.IPVersionIPv4 {
				ip = serverIP.Address.String()

				break
			}
		}

		assert.Equal(ctx.T, fmt.Sprintf("ssh %s -p 22 -l root -t", ip), strings.Join(cmd.Args, " "))

		return exitCode, nil
	}
}

func Test_ServerSSH(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: baremetal.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			selectOffer(zone, offerFilterAvailable),
			createServerAndWait(),
		),
		Cmd:          "scw baremetal server ssh {{ .Server.ID }} zone={{ .Server.Zone }}",
		OverrideExec: sshOverrideExec(0),
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc:       deleteServer(),
		DisableParallel: true,
	}))

	t.Run("With-Exit-Code", core.Test(&core.TestConfig{
		Commands: baremetal.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			selectOffer(zone, offerFilterAvailable),
			createServerAndWait(),
		),
		Cmd:          "scw baremetal server ssh {{ .Server.ID }} zone={{ .Server.Zone }}",
		OverrideExec: sshOverrideExec(130),
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(130),
		),
		AfterFunc:       deleteServer(),
		DisableParallel: true,
	}))

	t.Run("Server-Not-Ready", core.Test(&core.TestConfig{
		Commands: baremetal.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			selectOffer(zone, offerFilterAvailable),
			createServer(),
		),
		Cmd: "scw baremetal server ssh {{ .Server.ID }} zone={{ .Server.Zone }}",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(1),
		),
		AfterFunc:       core.AfterFuncCombine(waitForServer(), deleteServer()),
		DisableParallel: true,
	}))
}
