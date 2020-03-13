package baremetal

import (
	"testing"

	"github.com/alecthomas/assert"
	"github.com/scaleway/scaleway-cli/internal/core"
	baremetal "github.com/scaleway/scaleway-sdk-go/api/baremetal/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// deleteServerAfterFunc deletes the created server and its attached volumes and IPs.
func deleteServerAfterFunc(ctx *core.AfterFuncCtx) error {
	ctx.ExecuteCmd("scw baremetal server delete server-id=" + ctx.CmdResult.(*baremetal.Server).ID)
	return nil
}

// All test below should succeed to create an instance.
func Test_CreateServer(t *testing.T) {
	// Simple use cases
	t.Run("Simple", func(t *testing.T) {
		t.Run("Default", core.Test(&core.TestConfig{
			Commands: GetCommands(),
			Cmd:      "scw baremetal server create",
			Check: core.TestCheckCombine(
				core.TestCheckGolden(),
				core.TestCheckExitCode(0),
			),
			AfterFunc:   deleteServerAfterFunc,
			DefaultZone: scw.ZoneFrPar2,
		}))

		t.Run("With name", core.Test(&core.TestConfig{
			Commands: GetCommands(),
			Cmd:      "scw baremetal server create name=yo zone=fr-par-2",
			Check: core.TestCheckCombine(
				func(t *testing.T, ctx *core.CheckFuncCtx) {
					assert.Equal(t, "yo", ctx.Result.(*baremetal.Server).Name)
				},
				core.TestCheckExitCode(0),
			),
			AfterFunc: deleteServerAfterFunc,
		}))

		t.Run("Tags", core.Test(&core.TestConfig{
			Commands: GetCommands(),
			Cmd:      "scw baremetal server create tags.0=prod tags.1=blue",
			Check: core.TestCheckCombine(
				func(t *testing.T, ctx *core.CheckFuncCtx) {
					assert.Equal(t, "prod", ctx.Result.(*baremetal.Server).Tags[0])
					assert.Equal(t, "blue", ctx.Result.(*baremetal.Server).Tags[1])
				},
				core.TestCheckExitCode(0),
			),
			AfterFunc: deleteServerAfterFunc,
		}))

		//t.Run("HC-BM1-L", core.Test(&core.TestConfig{
		//	Commands: GetCommands(),
		//	Cmd:      "scw baremetal server create type=HC-BM1-L",
		//	Check: core.TestCheckCombine(
		//		func(t *testing.T, ctx *core.CheckFuncCtx) {
		//			assert.Equal(t, "HC-BM1-L", ctx.Result.(*baremetal.Server).CommercialType)
		//		},
		//		core.TestCheckExitCode(0),
		//	),
		//	AfterFunc:   deleteServerAfterFunc,
		//	DefaultZone: scw.ZoneFrPar2,
		//}))
	})
}

// None of the tests below should succeed to create an instance.
func Test_CreateServerErrors(t *testing.T) {
	////
	// Instance type errors
	////
	t.Run("Error: invalid instance type", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw instance server create type=MACBOOK1-S image=ubuntu_bionic",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(1),
		),
	}))
}
