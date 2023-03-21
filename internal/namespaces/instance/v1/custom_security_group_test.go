package instance

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/internal/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/lb/v1"
)

func Test_SecureBehindLBSecurityGroup(t *testing.T) {
	commands := GetCommands()
	commands.Merge(lb.GetCommands())
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: commands,
		BeforeFunc: core.BeforeFuncCombine(
			core.ExecStoreBeforeCmd("LB", "scw lb lb create name=foobar description=foobar --wait"),
			createServer("Server"),
			startServer("Server"),
		),
		Cmd: "scw instance security-group secure-behind-lb instance-id={{ .Server.ID }} lb-id={{ .LB.ID }}",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: core.AfterFuncCombine(
			deleteServer("Server"),
			core.ExecAfterCmd("scw lb lb delete {{ .LB.ID }}"),
		),
	}))
}
