package instance

import (
	"testing"

	"github.com/scaleway/scaleway-cli/internal/core"
)

func deleteServerAfterFunc(ctx *core.AfterFuncCtx) error {
	// Get ID of the created server.
	serverID, err := ctx.ExtractResourceID()
	if err != nil {
		return err
	}

	// Delete the test volume.
	ctx.ExecuteCmd("scw instance server delete server-id=" + serverID)
	return nil
}

// All test below should succeed to create an instance.
func Test_CreateServer(t *testing.T) {

	////
	// Image
	////
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands:  GetCommands(),
		Cmd:       "scw instance server create image=ubuntu-bionic",
		AfterFunc: deleteServerAfterFunc,
		Check:     core.TestCheckGolden(),
	}))

	// TODO: add all success cases
}

// None of the tests below should succeed to create an instance.
func Test_CreateServerErrors(t *testing.T) {

	////
	// Image errors
	////
	t.Run("Error: missing image label", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw instance server create",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(1),
		),
	}))

	t.Run("Error: invalid image label", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw instance server create image=macos",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(1),
		),
	}))

	t.Run("Error: invalid image UUID", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw instance server create image=7a892c1a-bbdc-491f-9974-4008e3708664",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(1),
		),
	}))

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

	////
	// Volume errors
	////
	t.Run("Error: invalid total local volumes size: too low 1", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw instance server create image=ubuntu_bionic root-volume=l:10GB",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(1),
		),
	}))

	t.Run("Error: invalid total local volumes size: too low 2", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw instance server create image=ubuntu_bionic root-volume=l:10GB additional-volumes.0=block:10GB",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(1),
		),
	}))

	t.Run("Error: invalid total local volumes size: too high 1", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw instance server create image=ubuntu_bionic root-volume=local:10GB additional-volumes.0=local:20GB",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(1),
		),
	}))

	t.Run("Error: invalid total local volumes size: too high 2", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw instance server create image=ubuntu_bionic additional-volumes.0=local:10GB",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(1),
		),
	}))

	t.Run("Error: invalid total local volumes size: too high 3", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		BeforeFunc: func(ctx *core.BeforeFuncCtx) error {
			ctx.Meta["response"] = ctx.ExecuteCmd("scw instance volume create name=cli-test size=20G volume-type=l_ssd")
			return nil
		},
		Cmd: "scw instance server create image=ubuntu_bionic root-volume={{ .response.volume.id }} additional-volumes.0=local:10GB",
		AfterFunc: func(ctx *core.AfterFuncCtx) error {
			ctx.ExecuteCmd("scw instance volume delete volume-id={{ .response.volume.id }}")
			return nil
		},
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(1),
		),
	}))

	t.Run("Error: invalid root volume size", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw instance server create image=ubuntu_bionic root-volume=local:2GB additional-volumes.0=local:18GB",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(1),
		),
	}))

	t.Run("Error: invalid root volume type", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw instance server create image=ubuntu_bionic root-volume=block:20GB",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(1),
		),
	}))

	t.Run("Error: disallow existing root volume ID", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		BeforeFunc: func(ctx *core.BeforeFuncCtx) error {
			ctx.Meta["response"] = ctx.ExecuteCmd("scw instance volume create name=cli-test size=20G volume-type=l_ssd")
			return nil
		},
		Cmd: "scw instance server create image=ubuntu_bionic root-volume={{ .response.volume.id }}",
		AfterFunc: func(ctx *core.AfterFuncCtx) error {
			ctx.ExecuteCmd("scw instance volume delete volume-id={{ .response.volume.id }}")
			return nil
		},
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(1),
		),
	}))

	t.Run("Error: invalid root volume ID", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw instance server create image=ubuntu_bionic root-volume=29da9ad9-e759-4a56-82c8-f0607f93055c",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(1),
		),
	}))

	t.Run("Error: already attached additional volume ID", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		BeforeFunc: func(ctx *core.BeforeFuncCtx) error {
			ctx.Meta["server"] = ctx.ExecuteCmd("scw instance server create name=cli-test image=ubuntu_bionic root-volume=l:10G additional-volumes.0=l:10G")
			return nil
		},
		Cmd: `scw instance server create image=ubuntu_bionic root-volume=l:10G additional-volumes.0={{ (index .server.volumes "1").id }}`,
		AfterFunc: func(ctx *core.AfterFuncCtx) error {
			ctx.ExecuteCmd("scw instance server delete server-id={{ .server.id }} delete-volumes delete-ip")
			return nil
		},
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(1),
		),
	}))

	t.Run("Error: invalid root volume format", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw instance server create image=ubuntu_bionic root-volume=20GB",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(1),
		),
	}))

	////
	// TODO: IP errors
	////
}
