package baremetal

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/internal/core"
	"github.com/scaleway/scaleway-cli/v2/internal/interactive"
)

func Test_CreateFlexibleIPInteractive(t *testing.T) {
	promptResponse := []string{
		`" "`,
		"description of flexibleIP",
		"tags flexible IP",
	}
	interactive.IsInteractive = true
	t.Run("Simple Interactive", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			createServerAndWaitDefault("Server"),
		),
		Cmd: "scw baremetal server add-flexible-ip {{ .Server.ID }}",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
		),
		AfterFunc: core.AfterFuncCombine(
			deleteServerDefault("Server"),
		),
		PromptResponseMocks: promptResponse,
	}))
}

func Test_CreateFlexibleIP(t *testing.T) {
	interactive.IsInteractive = false
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			createServerAndWaitDefault("Server"),
		),
		Cmd: "scw baremetal server add-flexible-ip {{ .Server.ID }} ip-type=IPv4",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
		),
		AfterFunc: core.AfterFuncCombine(
			deleteServerDefault("Server"),
		),
	}))
}
