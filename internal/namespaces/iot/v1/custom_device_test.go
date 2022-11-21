package iot

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/internal/core"
)

func Test_CreateDevice(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		BeforeFunc: createHub(),
		Cmd:        "scw iot device create hub-id={{ .Hub.ID }} name=foo",
		Check:      core.TestCheckGolden(),
		AfterFunc: core.AfterFuncCombine(
			core.ExecAfterCmd("scw iot device delete {{ .CmdResult.Device.ID }}"),
			deleteHub(),
		),
	}))
}
