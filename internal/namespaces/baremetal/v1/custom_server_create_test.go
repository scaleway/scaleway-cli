package baremetal

import (
	"testing"

	"github.com/alecthomas/assert"
	"github.com/scaleway/scaleway-cli/internal/core"
	baremetal "github.com/scaleway/scaleway-sdk-go/api/baremetal/v1"
)

// All test below should succeed to create an instance.
func Test_CreateServer(t *testing.T) {
	// Simple use cases
	t.Run("Simple", func(t *testing.T) {
		t.Run("Default", core.Test(&core.TestConfig{
			Commands: GetCommands(),
			Cmd:      "scw baremetal server create zone=nl-ams-1 type=GP-BM2-S -w",
			Check: core.TestCheckCombine(
				core.TestCheckGolden(),
				core.TestCheckExitCode(0),
			),
			AfterFunc: core.ExecAfterCmd("scw baremetal server delete {{ .CmdResult.ID }} zone=nl-ams-1"),
		},
		))

		t.Run("With name", core.Test(&core.TestConfig{
			Commands: GetCommands(),
			Cmd:      "scw baremetal server create name=test-create-server-with-name zone=nl-ams-1 type=GP-BM2-S -w",
			Check: core.TestCheckCombine(
				func(t *testing.T, ctx *core.CheckFuncCtx) {
					assert.Equal(t, "test-create-server-with-name", ctx.Result.(*baremetal.Server).Name)
				},
				core.TestCheckExitCode(0),
			),
			AfterFunc: core.ExecAfterCmd("scw baremetal server delete {{ .CmdResult.ID }} zone=nl-ams-1"),
		}))

		t.Run("Tags", core.Test(&core.TestConfig{
			Commands: GetCommands(),
			Cmd:      "scw baremetal server create tags.0=prod tags.1=blue zone=nl-ams-1 type=GP-BM2-S -w",
			Check: core.TestCheckCombine(
				func(t *testing.T, ctx *core.CheckFuncCtx) {
					assert.Equal(t, "prod", ctx.Result.(*baremetal.Server).Tags[0])
					assert.Equal(t, "blue", ctx.Result.(*baremetal.Server).Tags[1])
				},
				core.TestCheckExitCode(0),
			),
			AfterFunc: core.ExecAfterCmd("scw baremetal server delete {{ .CmdResult.ID }} zone=nl-ams-1"),
		}))
	})
}
