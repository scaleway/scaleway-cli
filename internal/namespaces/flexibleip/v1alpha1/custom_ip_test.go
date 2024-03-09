package flexibleip

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/internal/core"
)

func Test_CreateFlexibleWait(t *testing.T) {
	cmds := GetCommands()

	t.Run("Simple", core.Test(&core.TestConfig{
		Commands:  cmds,
		Cmd:       "scw fip ip create --wait",
		Check:     core.TestCheckGolden(),
		AfterFunc: core.ExecAfterCmd("scw fip ip delete {{ .CmdResult.ID }}"),
	}))
}
