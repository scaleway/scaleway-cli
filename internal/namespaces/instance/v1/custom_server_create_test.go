package instance

import (
	"testing"

	"github.com/alecthomas/assert"
	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
)

func deleteServerAfterFunc(ctx *core.AfterFuncCtx) error {
	// Delete volumes, ips and the server
	ctx.ExecuteCmd("scw instance server delete delete-volumes delete-ip force-shutdown server-id=" + ctx.CmdResult.(*instance.Server).ID)
	return nil
}

// All test below should succeed to create an instance.
func Test_CreateServer(t *testing.T) {
	////
	// Simple use cases
	////
	t.Run("Simple", func(t *testing.T) {
		t.Run("Default", core.Test(&core.TestConfig{
			Commands:  GetCommands(),
			Cmd:       "scw instance server create image=ubuntu-bionic",
			AfterFunc: deleteServerAfterFunc,
			Check: func(t *testing.T, ctx *core.CheckFuncCtx) {
				assert.Equal(t, "Ubuntu Bionic Beaver", ctx.Result.(*instance.Server).Image.Name)
			},
		}))

		t.Run("GP1-XS", core.Test(&core.TestConfig{
			Commands:  GetCommands(),
			Cmd:       "scw instance server create type=GP1-XS image=ubuntu-bionic",
			AfterFunc: deleteServerAfterFunc,
			Check: func(t *testing.T, ctx *core.CheckFuncCtx) {
				assert.Equal(t, "GP1-XS", ctx.Result.(*instance.Server).CommercialType)
			},
		}))

		t.Run("With name", core.Test(&core.TestConfig{
			Commands:  GetCommands(),
			Cmd:       "scw instance server create image=ubuntu-bionic name=yo",
			AfterFunc: deleteServerAfterFunc,
			Check: func(t *testing.T, ctx *core.CheckFuncCtx) {
				assert.Equal(t, "yo", ctx.Result.(*instance.Server).Name)
			},
		}))

		t.Run("With bootscript", core.Test(&core.TestConfig{
			Commands:  GetCommands(),
			Cmd:       "scw instance server create image=ubuntu-bionic bootscript-id=eb760e3c-30d8-49a3-b3ad-ad10c3aa440b",
			AfterFunc: deleteServerAfterFunc,
			Check: func(t *testing.T, ctx *core.CheckFuncCtx) {
				assert.Equal(t, "eb760e3c-30d8-49a3-b3ad-ad10c3aa440b", ctx.Result.(*instance.Server).Bootscript.ID)
			},
		}))

		t.Run("With start", core.Test(&core.TestConfig{
			Commands:  GetCommands(),
			Cmd:       "scw instance server create image=ubuntu-bionic start -w",
			AfterFunc: deleteServerAfterFunc,
			Check: func(t *testing.T, ctx *core.CheckFuncCtx) {
				assert.Equal(t, instance.ServerStateRunning, ctx.Result.(*instance.Server).State)
			},
		}))
	})

	// TODO: finish to add all success cases
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
			ctx.Meta["Response"] = ctx.ExecuteCmd("scw instance volume create name=cli-test size=20G volume-type=l_ssd")
			return nil
		},
		Cmd: "scw instance server create image=ubuntu_bionic root-volume={{ .Response.Volume.ID }} additional-volumes.0=local:10GB",
		AfterFunc: func(ctx *core.AfterFuncCtx) error {
			ctx.ExecuteCmd("scw instance volume delete volume-id={{ .Response.Volume.ID }}")
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
			ctx.Meta["Response"] = ctx.ExecuteCmd("scw instance volume create name=cli-test size=20G volume-type=l_ssd")
			return nil
		},
		Cmd: "scw instance server create image=ubuntu_bionic root-volume={{ .Response.Volume.ID }}",
		AfterFunc: func(ctx *core.AfterFuncCtx) error {
			ctx.ExecuteCmd("scw instance volume delete volume-id={{ .Response.Volume.ID }}")
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
			ctx.Meta["Server"] = ctx.ExecuteCmd("scw instance server create name=cli-test image=ubuntu_bionic root-volume=l:10G additional-volumes.0=l:10G")
			return nil
		},
		Cmd: `scw instance server create image=ubuntu_bionic root-volume=l:10G additional-volumes.0={{ (index .Server.Volumes "1").ID }}`,
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
	// IP errors
	////
	t.Run("Error: not found ip ID", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw instance server create image=ubuntu-bionic ip=23165951-13fd-4a3b-84ed-22c2e96658f2",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(1),
		),
	}))

	t.Run("Error: forbidden IP", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw instance server create image=ubuntu-bionic ip=51.15.242.82",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(1),
		),
	}))

	t.Run("Error: invalid ip", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw instance server create image=ubuntu-bionic ip=yo",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(1),
		),
	}))
}
