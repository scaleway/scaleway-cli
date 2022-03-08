package instance

import (
	"testing"

	"github.com/alecthomas/assert"
	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// deleteServerAfterFunc deletes the created server and its attached volumes and IPs.
func deleteServerAfterFunc() core.AfterFunc {
	return core.ExecAfterCmd("scw instance server delete {{ .CmdResult.ID }} with-volumes=all with-ip=true force-shutdown=true")
}

// All test below should succeed to create an instance.
func Test_CreateServer(t *testing.T) {
	////
	// Simple use cases
	////
	t.Run("Simple", func(t *testing.T) {
		t.Run("Default", core.Test(&core.TestConfig{
			Commands: GetCommands(),
			Cmd:      "scw instance server create image=ubuntu_bionic stopped=true",
			Check: core.TestCheckCombine(
				core.TestCheckGolden(),
				func(t *testing.T, ctx *core.CheckFuncCtx) {
					assert.Equal(t, "Ubuntu 18.04 Bionic Beaver", ctx.Result.(*instance.Server).Image.Name)
				},
				core.TestCheckExitCode(0),
			),
			AfterFunc: deleteServerAfterFunc(),
		}))

		t.Run("GP1-XS", core.Test(&core.TestConfig{
			Commands: GetCommands(),
			Cmd:      "scw instance server create type=GP1-XS image=ubuntu_bionic stopped=true",
			Check: core.TestCheckCombine(
				func(t *testing.T, ctx *core.CheckFuncCtx) {
					assert.Equal(t, "GP1-XS", ctx.Result.(*instance.Server).CommercialType)
				},
				core.TestCheckExitCode(0),
			),
			AfterFunc: deleteServerAfterFunc(),
		}))

		t.Run("With name", core.Test(&core.TestConfig{
			Commands: GetCommands(),
			Cmd:      "scw instance server create image=ubuntu_bionic name=yo stopped=true",
			Check: core.TestCheckCombine(
				func(t *testing.T, ctx *core.CheckFuncCtx) {
					assert.Equal(t, "yo", ctx.Result.(*instance.Server).Name)
				},
				core.TestCheckExitCode(0),
			),
			AfterFunc: deleteServerAfterFunc(),
		}))

		t.Run("With start", core.Test(&core.TestConfig{
			Commands: GetCommands(),
			Cmd:      "scw instance server create image=ubuntu_bionic -w",
			Check: core.TestCheckCombine(
				func(t *testing.T, ctx *core.CheckFuncCtx) {
					assert.Equal(t, instance.ServerStateRunning, ctx.Result.(*instance.Server).State)
				},
				core.TestCheckExitCode(0),
			),
			AfterFunc: deleteServerAfterFunc(),
		}))

		t.Run("With bootscript", core.Test(&core.TestConfig{
			Commands: GetCommands(),
			Cmd:      "scw instance server create image=ubuntu_bionic bootscript-id=eb760e3c-30d8-49a3-b3ad-ad10c3aa440b stopped=true",
			Check: core.TestCheckCombine(
				func(t *testing.T, ctx *core.CheckFuncCtx) {
					assert.Equal(t, "eb760e3c-30d8-49a3-b3ad-ad10c3aa440b", ctx.Result.(*instance.Server).Bootscript.ID)
					assert.Equal(t, instance.BootTypeBootscript, ctx.Result.(*instance.Server).BootType)
				},
				core.TestCheckExitCode(0),
			),
			AfterFunc: deleteServerAfterFunc(),
		}))

		t.Run("Image UUID", core.Test(&core.TestConfig{
			Commands: GetCommands(),
			Cmd:      "scw instance server create image=f974feac-abae-4365-b988-8ec7d1cec10d stopped=true",
			Check: core.TestCheckCombine(
				func(t *testing.T, ctx *core.CheckFuncCtx) {
					assert.Equal(t, "Ubuntu Bionic Beaver", ctx.Result.(*instance.Server).Image.Name)
				},
				core.TestCheckExitCode(0),
			),
			AfterFunc: deleteServerAfterFunc(),
		}))

		t.Run("Tags", core.Test(&core.TestConfig{
			Commands: GetCommands(),
			Cmd:      "scw instance server create image=ubuntu_bionic tags.0=prod tags.1=blue stopped=true",
			Check: core.TestCheckCombine(
				func(t *testing.T, ctx *core.CheckFuncCtx) {
					assert.Equal(t, "prod", ctx.Result.(*instance.Server).Tags[0])
					assert.Equal(t, "blue", ctx.Result.(*instance.Server).Tags[1])
				},
				core.TestCheckExitCode(0),
			),
			AfterFunc: deleteServerAfterFunc(),
		}))
	})

	////
	// Volume use cases
	////
	t.Run("Volumes", func(t *testing.T) {
		t.Run("valid single local volume", core.Test(&core.TestConfig{
			Commands: GetCommands(),
			Cmd:      "scw instance server create image=ubuntu_bionic root-volume=local:20GB stopped=true",
			Check: core.TestCheckCombine(
				func(t *testing.T, ctx *core.CheckFuncCtx) {
					assert.Equal(t, 20*scw.GB, ctx.Result.(*instance.Server).Volumes["0"].Size)
				},
				core.TestCheckExitCode(0),
			),
			AfterFunc: deleteServerAfterFunc(),
		}))

		t.Run("valid double local volumes", core.Test(&core.TestConfig{
			Commands: GetCommands(),
			Cmd:      "scw instance server create image=ubuntu_bionic root-volume=local:10GB additional-volumes.0=l:10G stopped=true",
			Check: core.TestCheckCombine(
				func(t *testing.T, ctx *core.CheckFuncCtx) {
					assert.Equal(t, 10*scw.GB, ctx.Result.(*instance.Server).Volumes["0"].Size)
					assert.Equal(t, 10*scw.GB, ctx.Result.(*instance.Server).Volumes["1"].Size)
				},
				core.TestCheckExitCode(0),
			),
			AfterFunc: deleteServerAfterFunc(),
		}))

		t.Run("valid additional block volumes", core.Test(&core.TestConfig{
			Commands: GetCommands(),
			Cmd:      "scw instance server create image=ubuntu_bionic additional-volumes.0=b:1G additional-volumes.1=b:5G additional-volumes.2=b:10G stopped=true",
			Check: core.TestCheckCombine(
				func(t *testing.T, ctx *core.CheckFuncCtx) {
					assert.Equal(t, 1*scw.GB, ctx.Result.(*instance.Server).Volumes["1"].Size)
					assert.Equal(t, 5*scw.GB, ctx.Result.(*instance.Server).Volumes["2"].Size)
					assert.Equal(t, 10*scw.GB, ctx.Result.(*instance.Server).Volumes["3"].Size)
				},
				core.TestCheckExitCode(0),
			),
			AfterFunc: deleteServerAfterFunc(),
		}))
	})
	////
	// IP use cases
	////
	t.Run("IPs", func(t *testing.T) {
		t.Run("explicit new IP", core.Test(&core.TestConfig{
			Commands: GetCommands(),
			Cmd:      "scw instance server create image=ubuntu_bionic ip=new stopped=true",
			Check: core.TestCheckCombine(
				func(t *testing.T, ctx *core.CheckFuncCtx) {
					assert.NotEmpty(t, ctx.Result.(*instance.Server).PublicIP.Address)
					assert.Equal(t, false, ctx.Result.(*instance.Server).PublicIP.Dynamic)
				},
				core.TestCheckExitCode(0),
			),
			AfterFunc: deleteServerAfterFunc(),
		}))

		t.Run("run with dynamic IP", core.Test(&core.TestConfig{
			Commands: GetCommands(),
			Cmd:      "scw instance server create image=ubuntu_bionic ip=dynamic -w", // dynamic IP is created at runtime
			Check: core.TestCheckCombine(
				func(t *testing.T, ctx *core.CheckFuncCtx) {
					assert.NoError(t, ctx.Err)
					assert.NotEmpty(t, ctx.Result.(*instance.Server).PublicIP.Address)
					assert.Equal(t, true, ctx.Result.(*instance.Server).DynamicIPRequired)
				},
				core.TestCheckExitCode(0),
			),
			AfterFunc: deleteServerAfterFunc(),
		}))

		t.Run("existing IP", core.Test(&core.TestConfig{
			Commands:   GetCommands(),
			BeforeFunc: createIP("IP"),
			Cmd:        "scw instance server create image=ubuntu_bionic ip={{ .IP.Address }} stopped=true",
			Check: core.TestCheckCombine(
				func(t *testing.T, ctx *core.CheckFuncCtx) {
					assert.NotEmpty(t, ctx.Result.(*instance.Server).PublicIP.Address)
					assert.Equal(t, false, ctx.Result.(*instance.Server).PublicIP.Dynamic)
				},
				core.TestCheckExitCode(0),
			),
			AfterFunc: deleteServerAfterFunc(),
		}))

		t.Run("existing IP ID", core.Test(&core.TestConfig{
			Commands:   GetCommands(),
			BeforeFunc: createIP("IP"),
			Cmd:        "scw instance server create image=ubuntu_bionic ip={{ .IP.ID }} stopped=true",
			Check: core.TestCheckCombine(
				func(t *testing.T, ctx *core.CheckFuncCtx) {
					assert.NotEmpty(t, ctx.Result.(*instance.Server).PublicIP.Address)
					assert.Equal(t, false, ctx.Result.(*instance.Server).PublicIP.Dynamic)
				},
				core.TestCheckExitCode(0),
			),
			AfterFunc: deleteServerAfterFunc(),
		}))

		t.Run("with ipv6", core.Test(&core.TestConfig{
			Commands: GetCommands(),
			Cmd:      "scw instance server create image=ubuntu_bionic ipv6=true -w", // IPv6 is created at runtime
			Check: core.TestCheckCombine(
				func(t *testing.T, ctx *core.CheckFuncCtx) {
					assert.NotEmpty(t, ctx.Result.(*instance.Server).IPv6.Address)
				},
				core.TestCheckExitCode(0),
			),
			AfterFunc: deleteServerAfterFunc(),
		}))
	})
}

// None of the tests below should succeed to create an instance.
// These tests need to be run in sequence since they are having warnings in the stderr
// and these warnings can be captured by other tests
func Test_CreateServerErrors(t *testing.T) {
	////
	// Image errors
	////
	t.Run("Error: invalid image label", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw instance server create image=macos",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(1),
		),
		DisableParallel: true,
	}))

	t.Run("Error: invalid image UUID", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw instance server create image=7a892c1a-bbdc-491f-9974-4008e3708664",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(1),
		),
		DisableParallel: true,
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
		DisableParallel: true,
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
		DisableParallel: true,
	}))

	t.Run("Error: invalid total local volumes size: too low 2", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw instance server create image=ubuntu_bionic root-volume=l:10GB additional-volumes.0=block:10GB",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(1),
		),
		DisableParallel: true,
	}))

	t.Run("Error: invalid total local volumes size: too low 3", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw instance server create image=ubuntu_bionic root-volume=block:20GB",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(1),
		),
		DisableParallel: true,
	}))

	t.Run("Error: invalid total local volumes size: too high 1", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw instance server create image=ubuntu_bionic root-volume=local:10GB additional-volumes.0=local:20GB",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(1),
		),
		DisableParallel: true,
	}))

	t.Run("Error: invalid total local volumes size: too high 2", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw instance server create image=ubuntu_bionic additional-volumes.0=local:10GB",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(1),
		),
		DisableParallel: true,
	}))

	t.Run("Error: invalid total local volumes size: too high 3", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		BeforeFunc: createVolume("Volume", 20, instance.VolumeVolumeTypeLSSD),
		Cmd:        "scw instance server create image=ubuntu_bionic root-volume={{ .Volume.ID }} additional-volumes.0=local:10GB",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(1),
		),
		AfterFunc:       deleteVolume("Volume"),
		DisableParallel: true,
	}))

	t.Run("Error: invalid root volume size", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw instance server create image=ubuntu_bionic root-volume=local:2GB additional-volumes.0=local:18GB",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(1),
		),
		DisableParallel: true,
	}))

	t.Run("Error: disallow existing root volume ID", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		BeforeFunc: createVolume("Volume", 20, instance.VolumeVolumeTypeLSSD),
		Cmd:        "scw instance server create image=ubuntu_bionic root-volume={{ .Volume.ID }}",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(1),
		),
		AfterFunc:       deleteVolume("Volume"),
		DisableParallel: true,
	}))

	t.Run("Error: invalid root volume ID", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw instance server create image=ubuntu_bionic root-volume=29da9ad9-e759-4a56-82c8-f0607f93055c",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(1),
		),
		DisableParallel: true,
	}))

	t.Run("Error: already attached additional volume ID", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		BeforeFunc: core.ExecStoreBeforeCmd("Server", "scw instance server create name=cli-test image=ubuntu_bionic root-volume=l:10G additional-volumes.0=l:10G stopped=true"),
		Cmd:        `scw instance server create image=ubuntu_bionic root-volume=l:10G additional-volumes.0={{ (index .Server.Volumes "1").ID }} stopped=true`,
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(1),
		),
		AfterFunc:       deleteServer("Server"),
		DisableParallel: true,
	}))

	t.Run("Error: invalid root volume format", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw instance server create image=ubuntu_bionic root-volume=20GB",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(1),
		),
		DisableParallel: true,
	}))

	////
	// IP errors
	////
	t.Run("Error: not found ip ID", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw instance server create image=ubuntu_bionic ip=23165951-13fd-4a3b-84ed-22c2e96658f2",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(1),
		),
	}))

	t.Run("Error: forbidden IP", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw instance server create image=ubuntu_bionic ip=51.15.242.82",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(1),
		),
	}))

	t.Run("Error: invalid ip", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw instance server create image=ubuntu_bionic ip=yo",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(1),
		),
	}))

	t.Run("Error: image size is incompatible with instance type", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw instance server create image=d4067cdc-dc9d-4810-8a26-0dae51d7df42 type=DEV1-S",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(1),
		),
	}))
}
