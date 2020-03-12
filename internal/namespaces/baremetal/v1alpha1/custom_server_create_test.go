package baremetal

import (
	"testing"

	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
)

// deleteServerAfterFunc deletes the created server and its attached volumes and IPs.
func deleteServerAfterFunc(ctx *core.AfterFuncCtx) error {
	ctx.ExecuteCmd("scw instance server delete with-volumes=all with-ip force-shutdown server-id=" + ctx.CmdResult.(*instance.Server).ID)
	return nil
}

// All test below should succeed to create an instance.
func Test_CreateServer(t *testing.T) {

	////
	// Simple use cases
	////
	t.Run("Simple", func(t *testing.T) {
		t.Run("Default", core.Test(&core.TestConfig{
			Commands: GetCommands(),
			Cmd:      "scw baremetal server create offer-id=964f9b38-577e-470f-a220-7d762f9e8672 name=pouet zone=fr-par-2",
			Check: core.TestCheckCombine(
				core.TestCheckGolden(),
				core.TestCheckExitCode(0),
			),
			AfterFunc: deleteServerAfterFunc,
		}))

		//t.Run("GP1-XS", core.Test(&core.TestConfig{
		//	Commands: GetCommands(),
		//	Cmd:      "scw instance server create type=GP1-XS image=ubuntu_bionic stopped",
		//	Check: core.TestCheckCombine(
		//		func(t *testing.T, ctx *core.CheckFuncCtx) {
		//			assert.Equal(t, "GP1-XS", ctx.Result.(*instance.Server).CommercialType)
		//		},
		//		core.TestCheckExitCode(0),
		//	),
		//	AfterFunc: deleteServerAfterFunc,
		//}))
		//
		//t.Run("With name", core.Test(&core.TestConfig{
		//	Commands: GetCommands(),
		//	Cmd:      "scw instance server create image=ubuntu_bionic name=yo",
		//	Check: core.TestCheckCombine(
		//		func(t *testing.T, ctx *core.CheckFuncCtx) {
		//			assert.Equal(t, "yo", ctx.Result.(*instance.Server).Name)
		//		},
		//		core.TestCheckExitCode(0),
		//	),
		//	AfterFunc: deleteServerAfterFunc,
		//}))
		//
		//t.Run("With start", core.Test(&core.TestConfig{
		//	Commands: GetCommands(),
		//	Cmd:      "scw instance server create image=ubuntu_bionic -w",
		//	Check: core.TestCheckCombine(
		//		func(t *testing.T, ctx *core.CheckFuncCtx) {
		//			assert.Equal(t, instance.ServerStateRunning, ctx.Result.(*instance.Server).State)
		//		},
		//		core.TestCheckExitCode(0),
		//	),
		//	AfterFunc: deleteServerAfterFunc,
		//}))

		//t.Run("Tags", core.Test(&core.TestConfig{
		//	Commands: GetCommands(),
		//	Cmd:      "scw baremetal server create image=ubuntu_bionic tags.0=prod tags.1=blue stopped",
		//	Check: core.TestCheckCombine(
		//		func(t *testing.T, ctx *core.CheckFuncCtx) {
		//			assert.Equal(t, "prod", ctx.Result.(*instance.Server).Tags[0])
		//			assert.Equal(t, "blue", ctx.Result.(*instance.Server).Tags[1])
		//		},
		//		core.TestCheckExitCode(0),
		//	),
		//	AfterFunc: deleteServerAfterFunc,
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
