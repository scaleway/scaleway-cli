package baremetal_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/interactive"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/baremetal/v1"
	flexibleip "github.com/scaleway/scaleway-cli/v2/internal/namespaces/flexibleip/v1alpha1"
)

func Test_CreateFlexibleIPInteractive(t *testing.T) {
	promptResponse := []string{
		`" "`,
	}
	interactive.IsInteractive = true
	cmds := baremetal.GetCommands()
	cmds.Merge(flexibleip.GetCommands())
	t.Run("Simple Interactive", core.Test(&core.TestConfig{
		Commands: cmds,
		BeforeFunc: core.BeforeFuncCombine(
			checkStockOffer(),
			createServerAndWait(),
		),
		Cmd: "scw baremetal server add-flexible-ip {{ .Server.ID }} zone=" + zone,
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
		),
		AfterFunc: core.AfterFuncCombine(
			deleteServer("Server"),
			core.ExecAfterCmd("scw fip ip delete {{ .CmdResult.ID }} zone="+zone),
		),
		PromptResponseMocks: promptResponse,
	}))
}

func Test_CreateFlexibleIP(t *testing.T) {
	interactive.IsInteractive = false
	cmds := baremetal.GetCommands()
	cmds.Merge(flexibleip.GetCommands())
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: cmds,
		BeforeFunc: core.BeforeFuncCombine(
			checkStockOffer(),
			createServerAndWait(),
		),
		Cmd: "scw baremetal server add-flexible-ip {{ .Server.ID }} ip-type=IPv4 zone=" + zone,
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
		),
		AfterFunc: core.AfterFuncCombine(
			deleteServer("Server"),
			core.ExecAfterCmd("scw fip ip delete {{ .CmdResult.ID }} zone="+zone),
		),
	}))
}
