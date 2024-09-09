package instance_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/instance/v1"
)

func Test_IpCreate(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: instance.GetCommands(),
		Cmd:      "scw instance ip create",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: core.ExecAfterCmd("scw instance ip delete {{ .CmdResult.IP.ID }}"),
	}))
}

func Test_IpDelete(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		BeforeFunc: createIP("Ip"),
		Commands:   instance.GetCommands(),
		Cmd:        "scw instance ip delete {{ .Ip.ID }}",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
	}))
}

func Test_IpGet(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		BeforeFunc: createIP("Ip"),
		Commands:   instance.GetCommands(),
		Cmd:        "scw instance ip get {{ .Ip.ID }}",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: deleteIP("Ip"),
	}))
}

func Test_IpList(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		BeforeFunc: createIP("Ip"),
		Commands:   instance.GetCommands(),
		Cmd:        "scw instance ip list",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: deleteIP("Ip"),
	}))
}

func Test_IpUpdate(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		BeforeFunc: createIP("Ip"),
		Commands:   instance.GetCommands(),
		Cmd:        "scw instance ip update {{ .Ip.ID }}",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: deleteIP("Ip"),
	}))
}
