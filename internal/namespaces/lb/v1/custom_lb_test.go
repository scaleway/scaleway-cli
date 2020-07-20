package lb

import (
	"testing"

	"github.com/scaleway/scaleway-cli/internal/core"
)

func Test_ListLB(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		BeforeFunc: createLB(),
		Cmd:        "scw lb lb list",
		Check:      core.TestCheckGolden(),
		AfterFunc:  deleteLB(),
	}))
}

func Test_CreateLB(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands:  GetCommands(),
		Cmd:       "scw lb lb create name=foobar description=foobar --wait",
		Check:     core.TestCheckGolden(),
		AfterFunc: core.ExecAfterCmd("scw lb lb delete {{ .CmdResult.ID }}"),
	}))
}

func Test_GetLB(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		BeforeFunc: createLB(),
		Cmd:        "scw lb lb get {{ .LB.ID }}",
		Check:      core.TestCheckGolden(),
		AfterFunc:  deleteLB(),
	}))
}

func Test_WaitLB(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		BeforeFunc: core.ExecStoreBeforeCmd(
			"LB",
			"scw lb lb create name=cli-test description=cli-test",
		),
		Cmd:       "scw lb lb wait {{ .LB.ID }}",
		Check:     core.TestCheckGolden(),
		AfterFunc: deleteLB(),
	}))
}
