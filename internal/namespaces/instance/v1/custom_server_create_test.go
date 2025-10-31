package instance_test

import (
	"fmt"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	block "github.com/scaleway/scaleway-cli/v2/internal/namespaces/block/v1alpha1"
	file "github.com/scaleway/scaleway-cli/v2/internal/namespaces/file/v1alpha1"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/instance/v1"
	"github.com/scaleway/scaleway-cli/v2/internal/testhelpers"
	blockSDK "github.com/scaleway/scaleway-sdk-go/api/block/v1alpha1"
	instanceSDK "github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// deleteServerAfterFunc deletes the created server and its attached volumes and IPs.
func deleteServerAfterFunc() core.AfterFunc {
	return core.ExecAfterCmd(
		"scw instance server delete {{ .CmdResult.ID }} with-volumes=all with-ip=true force-shutdown=true",
	)
}

// All test below should succeed to create an instance.
func Test_CreateServer(t *testing.T) {
	////
	// Simple use cases
	////
	t.Run("Simple", func(t *testing.T) {
		t.Run("Default", core.Test(&core.TestConfig{
			Commands: instance.GetCommands(),
			Cmd:      testServerCommand("image=ubuntu_jammy stopped=true"),
			Check: core.TestCheckCombine(
				core.TestCheckGolden(),
				func(t *testing.T, ctx *core.CheckFuncCtx) {
					t.Helper()
					server := ctx.Result.(*instance.ServerWithWarningsResponse).Server
					assert.NotNil(t, ctx.Result)
					assert.Equal(
						t,
						"Ubuntu 22.04 Jammy Jellyfish",
						server.Image.Name,
					)
				},
				core.TestCheckExitCode(0),
			),
			AfterFunc: deleteServerAfterFunc(),
		}))

		t.Run("GP1-XS", core.Test(&core.TestConfig{
			Commands: instance.GetCommands(),
			Cmd:      "scw instance server create type=GP1-XS image=ubuntu_jammy stopped=true",
			Check: core.TestCheckCombine(
				core.TestCheckExitCode(0),
				func(t *testing.T, ctx *core.CheckFuncCtx) {
					t.Helper()
					server := ctx.Result.(*instance.ServerWithWarningsResponse).Server
					assert.NotNil(t, ctx.Result)
					assert.Equal(t, "GP1-XS", server.CommercialType)
				},
			),
			AfterFunc: deleteServerAfterFunc(),
		}))

		t.Run("With name", core.Test(&core.TestConfig{
			Commands: instance.GetCommands(),
			Cmd:      testServerCommand("image=ubuntu_jammy name=yo stopped=true"),
			Check: core.TestCheckCombine(
				core.TestCheckExitCode(0),
				func(t *testing.T, ctx *core.CheckFuncCtx) {
					t.Helper()
					server := ctx.Result.(*instance.ServerWithWarningsResponse).Server
					assert.NotNil(t, ctx.Result)
					assert.Equal(t, "yo", server.Name)
				},
			),
			AfterFunc: deleteServerAfterFunc(),
		}))

		t.Run("With start", core.Test(&core.TestConfig{
			Commands: instance.GetCommands(),
			Cmd:      testServerCommand("image=ubuntu_jammy -w"),
			Check: core.TestCheckCombine(
				core.TestCheckExitCode(0),
				func(t *testing.T, ctx *core.CheckFuncCtx) {
					t.Helper()
					server := ctx.Result.(*instance.ServerWithWarningsResponse).Server
					assert.NotNil(t, ctx.Result)
					assert.Equal(
						t,
						instanceSDK.ServerStateRunning,
						server.State,
					)
				},
			),
			AfterFunc: deleteServerAfterFunc(),
		}))

		t.Run("Image UUID", core.Test(&core.TestConfig{
			Commands: instance.GetCommands(),
			Cmd:      testServerCommand("image=f974feac-abae-4365-b988-8ec7d1cec10d stopped=true"),
			Check: core.TestCheckCombine(
				func(t *testing.T, ctx *core.CheckFuncCtx) {
					t.Helper()
					server := ctx.Result.(*instance.ServerWithWarningsResponse).Server
					assert.NotNil(t, ctx.Result)
					assert.Equal(
						t,
						"Ubuntu Bionic Beaver",
						server.Image.Name,
					)
				},
				core.TestCheckExitCode(0),
			),
			AfterFunc: deleteServerAfterFunc(),
		}))

		t.Run("Tags", core.Test(&core.TestConfig{
			Commands: instance.GetCommands(),
			Cmd:      testServerCommand("image=ubuntu_jammy tags.0=prod tags.1=blue stopped=true"),
			Check: core.TestCheckCombine(
				core.TestCheckExitCode(0),
				func(t *testing.T, ctx *core.CheckFuncCtx) {
					t.Helper()
					server := ctx.Result.(*instance.ServerWithWarningsResponse).Server
					assert.NotNil(t, ctx.Result)
					assert.Equal(t, "prod", server.Tags[0])
					assert.Equal(t, "blue", server.Tags[1])
				},
			),
			AfterFunc: deleteServerAfterFunc(),
		}))
	})

	////
	// Volume use cases
	////
	t.Run("Volumes", func(t *testing.T) {
		t.Run("valid single local volume", core.Test(&core.TestConfig{
			Commands: instance.GetCommands(),
			Cmd:      testServerCommand("image=ubuntu_bionic root-volume=local:20GB stopped=true"),
			Check: core.TestCheckCombine(
				func(t *testing.T, ctx *core.CheckFuncCtx) {
					t.Helper()
					assert.NotNil(t, ctx.Result)
					server := testhelpers.Value[*instance.ServerWithWarningsResponse](
						t,
						ctx.Result,
					).Server
					volume := testhelpers.MapTValue(t, server.Volumes, "0")
					size := volume.Size
					assert.Equal(
						t,
						20*scw.GB,
						instance.SizeValue(size),
						"Size of volume should be 20 GB",
					)
				},
				core.TestCheckExitCode(0),
			),
			AfterFunc: deleteServerAfterFunc(),
		}))

		t.Run("valid single local snapshot", core.Test(&core.TestConfig{
			Commands: instance.GetCommands(),
			BeforeFunc: core.BeforeFuncCombine(
				core.ExecStoreBeforeCmd(
					"Server",
					testServerCommand("image=ubuntu_bionic root-volume=local:20GB stopped=true"),
				),
				core.ExecStoreBeforeCmd(
					"Snapshot",
					`scw instance snapshot create volume-id={{ (index .Server.Volumes "0").ID }}`,
				),
			),
			Cmd: testServerCommand(
				"image=ubuntu_bionic root-volume=local:{{ .Snapshot.Snapshot.ID }} stopped=true",
			),
			Check: core.TestCheckCombine(
				core.TestCheckExitCode(0),
				func(t *testing.T, ctx *core.CheckFuncCtx) {
					t.Helper()
					assert.NotNil(t, ctx.Result)
					server := testhelpers.Value[*instance.ServerWithWarningsResponse](
						t,
						ctx.Result,
					).Server
					volume := testhelpers.MapTValue(t, server.Volumes, "0")
					size := volume.Size
					assert.Equal(
						t,
						20*scw.GB,
						instance.SizeValue(size),
						"Size of volume should be 20 GB",
					)
				},
			),
			AfterFunc: core.AfterFuncCombine(
				deleteServer("Server"),
				deleteServerAfterFunc(),
				deleteSnapshot("Snapshot"),
			),
		}))

		t.Run("valid single local snapshot without image", core.Test(&core.TestConfig{
			Commands: instance.GetCommands(),
			BeforeFunc: core.BeforeFuncCombine(
				core.ExecStoreBeforeCmd(
					"Server",
					testServerCommand("image=ubuntu_bionic root-volume=local:20GB stopped=true"),
				),
				core.ExecStoreBeforeCmd(
					"Snapshot",
					`scw instance snapshot create volume-id={{ (index .Server.Volumes "0").ID }}`,
				),
			),
			Cmd: testServerCommand(
				"image=none root-volume=local:{{ .Snapshot.Snapshot.ID }} stopped=true",
			),
			Check: core.TestCheckCombine(
				core.TestCheckExitCode(0),
				func(t *testing.T, ctx *core.CheckFuncCtx) {
					t.Helper()
					assert.NotNil(t, ctx.Result)
					server := testhelpers.Value[*instance.ServerWithWarningsResponse](
						t,
						ctx.Result,
					).Server
					volume := testhelpers.MapTValue(t, server.Volumes, "0")
					size := volume.Size
					assert.Equal(
						t,
						20*scw.GB,
						instance.SizeValue(size),
						"Size of volume should be 20 GB",
					)
				},
			),
			AfterFunc: core.AfterFuncCombine(
				deleteServer("Server"),
				deleteServerAfterFunc(),
				deleteSnapshot("Snapshot"),
			),
		}))

		t.Run("valid double local volumes", core.Test(&core.TestConfig{
			Commands: instance.GetCommands(),
			Cmd: testServerCommand(
				"image=ubuntu_bionic root-volume=local:10GB additional-volumes.0=l:10G stopped=true",
			),
			Check: core.TestCheckCombine(
				func(t *testing.T, ctx *core.CheckFuncCtx) {
					t.Helper()
					assert.NotNil(t, ctx.Result)
					server := testhelpers.Value[*instance.ServerWithWarningsResponse](
						t,
						ctx.Result,
					).Server
					size0 := testhelpers.MapTValue(t, server.Volumes, "0").Size
					size1 := testhelpers.MapTValue(t, server.Volumes, "1").Size
					assert.Equal(
						t,
						10*scw.GB,
						instance.SizeValue(size0),
						"Size of volume should be 10 GB",
					)
					assert.Equal(
						t,
						10*scw.GB,
						instance.SizeValue(size1),
						"Size of volume should be 10 GB",
					)
				},
				core.TestCheckExitCode(0),
			),
			AfterFunc: deleteServerAfterFunc(),
		}))

		t.Run("valid double snapshot", core.Test(&core.TestConfig{
			Commands: core.NewCommandsMerge(
				instance.GetCommands(),
				block.GetCommands(),
			),
			BeforeFunc: core.BeforeFuncCombine(
				core.ExecStoreBeforeCmd(
					"Server",
					testServerCommand("image=ubuntu_jammy root-volume=block:20GB stopped=true"),
				),
				core.ExecStoreBeforeCmd(
					"Snapshot",
					`scw block snapshot create volume-id={{ (index .Server.Volumes "0").ID }} -w`,
				),
			),
			Cmd: testServerCommand(
				"image=ubuntu_jammy root-volume=block:{{ .Snapshot.ID }} additional-volumes.0=block:{{ .Snapshot.ID }} stopped=true",
			),
			Check: core.TestCheckCombine(
				core.TestCheckExitCode(0),
				testServerSBSVolumeSize("0", 20),
				testServerSBSVolumeSize("1", 20),
			),
			AfterFunc: core.AfterFuncCombine(
				deleteServer("Server"),
				deleteServerAfterFunc(),
				deleteBlockSnapshot("Snapshot"),
			),
		}))

		t.Run("valid additional block volumes", core.Test(&core.TestConfig{
			Commands: instance.GetCommands(),
			Cmd: testServerCommand(
				"image=ubuntu_jammy additional-volumes.0=b:1G additional-volumes.1=b:5G additional-volumes.2=b:10G stopped=true",
			),
			Check: core.TestCheckCombine(
				core.TestCheckExitCode(0),
				testServerSBSVolumeSize("1", 1),
				testServerSBSVolumeSize("2", 5),
				testServerSBSVolumeSize("3", 10),
			),
			AfterFunc: deleteServerAfterFunc(),
		}))

		t.Run("sbs additional volumes from id", core.Test(&core.TestConfig{
			Commands: core.NewCommandsMerge(
				instance.GetCommands(),
				block.GetCommands(),
			),
			BeforeFunc: core.BeforeFuncCombine(
				createSbsVolume("Volume", 20),
			),
			Cmd: testServerCommand(
				"image=ubuntu_jammy additional-volumes.0={{.Volume.ID}} stopped=true",
			),
			Check: core.TestCheckCombine(
				func(t *testing.T, ctx *core.CheckFuncCtx) {
					t.Helper()
					assert.NotNil(t, ctx.Result)
					server := testhelpers.Value[*instance.ServerWithWarningsResponse](
						t,
						ctx.Result,
					).Server
					volume := testhelpers.MapTValue(t, server.Volumes, "1")
					assert.Equal(t, instanceSDK.VolumeServerVolumeTypeSbsVolume, volume.VolumeType)
				},
				core.TestCheckExitCode(0),
			),
			AfterFunc: core.AfterFuncCombine(
				deleteServerAfterFunc(),
			),
		}))

		t.Run("sbs additional volumes", core.Test(&core.TestConfig{
			Commands: core.NewCommandsMerge(
				instance.GetCommands(),
				block.GetCommands(),
			),
			Cmd: testServerCommand("image=ubuntu_jammy additional-volumes.0=sbs:20G stopped=true"),
			Check: core.TestCheckCombine(
				core.TestCheckExitCode(0),
				func(t *testing.T, ctx *core.CheckFuncCtx) {
					t.Helper()
					assert.NotNil(t, ctx.Result)
					server := testhelpers.Value[*instance.ServerWithWarningsResponse](
						t,
						ctx.Result,
					).Server
					volume := testhelpers.MapTValue(t, server.Volumes, "1")
					assert.Equal(t, instanceSDK.VolumeServerVolumeTypeSbsVolume, volume.VolumeType)
				},
			),
			AfterFunc: core.AfterFuncCombine(
				deleteServerAfterFunc(),
			),
		}))

		t.Run("use sbs root volume", core.Test(&core.TestConfig{
			Commands: core.NewCommandsMerge(
				instance.GetCommands(),
				block.GetCommands(),
			),
			BeforeFunc: core.BeforeFuncCombine(
				createSbsVolume("Volume", 20),
			),
			Cmd: testServerCommand("image=none root-volume={{.Volume.ID}} stopped=true"),
			Check: core.TestCheckCombine(
				core.TestCheckExitCode(0),
				func(t *testing.T, ctx *core.CheckFuncCtx) {
					t.Helper()
					assert.NotNil(t, ctx.Result)
					server := testhelpers.Value[*instance.ServerWithWarningsResponse](
						t,
						ctx.Result,
					).Server
					volume := testhelpers.MapTValue(t, server.Volumes, "0")
					assert.Equal(t, instanceSDK.VolumeServerVolumeTypeSbsVolume, volume.VolumeType)
				},
			),
			AfterFunc: core.AfterFuncCombine(
				deleteServerAfterFunc(),
			),
		}))

		t.Run("create sbs root volume", core.Test(&core.TestConfig{
			Commands: core.NewCommandsMerge(
				instance.GetCommands(),
				block.GetCommands(),
			),
			Cmd: testServerCommand("image=ubuntu_jammy root-volume=sbs:20GB stopped=true"),
			Check: core.TestCheckCombine(
				core.TestCheckExitCode(0),
				func(t *testing.T, ctx *core.CheckFuncCtx) {
					t.Helper()
					assert.NotNil(t, ctx.Result)
					server := testhelpers.Value[*instance.ServerWithWarningsResponse](
						t,
						ctx.Result,
					).Server
					volume := testhelpers.MapTValue(t, server.Volumes, "0")
					assert.Equal(t, instanceSDK.VolumeServerVolumeTypeSbsVolume, volume.VolumeType)
				},
			),
			AfterFunc: core.AfterFuncCombine(
				deleteServerAfterFunc(),
			),
		}))

		t.Run("create sbs root volume with iops", core.Test(&core.TestConfig{
			Commands: core.NewCommandsMerge(
				instance.GetCommands(),
				block.GetCommands(),
			),
			Cmd: testServerCommand(
				"image=ubuntu_jammy root-volume=sbs:20GB:15000 stopped=true --debug",
			),
			Check: core.TestCheckCombine(
				core.TestCheckExitCode(0),
				func(t *testing.T, ctx *core.CheckFuncCtx) {
					t.Helper()
					assert.NotNil(t, ctx.Result)
					server := testhelpers.Value[*instance.ServerWithWarningsResponse](
						t,
						ctx.Result,
					).Server

					rootVolume, rootVolumeExists := server.Volumes["0"]
					assert.True(t, rootVolumeExists)
					assert.Equal(
						t,
						instanceSDK.VolumeServerVolumeTypeSbsVolume,
						rootVolume.VolumeType,
					)

					api := blockSDK.NewAPI(ctx.Client)
					vol, err := api.WaitForVolume(&blockSDK.WaitForVolumeRequest{
						VolumeID:      rootVolume.ID,
						Zone:          rootVolume.Zone,
						RetryInterval: core.DefaultRetryInterval,
					})
					require.NoError(t, err)
					assert.NotNil(t, vol.Specs)
					assert.NotNil(t, vol.Specs.PerfIops)
					assert.Equal(t, uint32(15000), *vol.Specs.PerfIops)
				},
			),
			AfterFunc: core.AfterFuncCombine(
				deleteServerAfterFunc(),
			),
		}))
	})
	////
	// IP use cases
	////
	t.Run("IPs", func(t *testing.T) {
		t.Run("explicit new IP", core.Test(&core.TestConfig{
			Commands: instance.GetCommands(),
			Cmd:      testServerCommand("ip=new stopped=true"),
			Check: core.TestCheckCombine(
				func(t *testing.T, ctx *core.CheckFuncCtx) {
					t.Helper()
					assert.NotNil(t, ctx.Result)
					server := testhelpers.Value[*instance.ServerWithWarningsResponse](
						t,
						ctx.Result,
					).Server
					assert.NotNil(t, server.PublicIP)
					assert.NotEmpty(t, server.PublicIP.Address)
					assert.False(t, server.PublicIP.Dynamic)
				},
				core.TestCheckExitCode(0),
			),
			AfterFunc: deleteServerAfterFunc(),
		}))

		t.Run("run with dynamic IP", core.Test(&core.TestConfig{
			Commands: instance.GetCommands(),
			Cmd:      testServerCommand("ip=dynamic -w"), // dynamic IP is created at runtime
			Check: core.TestCheckCombine(
				func(t *testing.T, ctx *core.CheckFuncCtx) {
					t.Helper()
					require.NoError(t, ctx.Err)
					assert.NotNil(t, ctx.Result)
					server := testhelpers.Value[*instance.ServerWithWarningsResponse](
						t,
						ctx.Result,
					).Server
					assert.NotNil(t, server.PublicIP)
					assert.NotEmpty(t, server.PublicIP.Address)
					assert.True(t, server.DynamicIPRequired)
				},
				core.TestCheckExitCode(0),
			),
			AfterFunc: deleteServerAfterFunc(),
		}))

		t.Run("existing IP", core.Test(&core.TestConfig{
			Commands:   instance.GetCommands(),
			BeforeFunc: createIP("IP"),
			Cmd:        testServerCommand("ip={{ .IP.Address }} stopped=true"),
			Check: core.TestCheckCombine(
				func(t *testing.T, ctx *core.CheckFuncCtx) {
					t.Helper()
					assert.NotNil(t, ctx.Result)
					server := testhelpers.Value[*instance.ServerWithWarningsResponse](
						t,
						ctx.Result,
					).Server
					assert.NotNil(t, server.PublicIP)
					assert.NotEmpty(t, server.PublicIP.Address)
					assert.False(t, server.PublicIP.Dynamic)
				},
				core.TestCheckExitCode(0),
			),
			AfterFunc: deleteServerAfterFunc(),
		}))

		t.Run("existing IP ID", core.Test(&core.TestConfig{
			Commands:   instance.GetCommands(),
			BeforeFunc: createIP("IP"),
			Cmd:        testServerCommand("ip={{ .IP.ID }} stopped=true"),
			Check: core.TestCheckCombine(
				func(t *testing.T, ctx *core.CheckFuncCtx) {
					t.Helper()
					assert.NotNil(t, ctx.Result)
					server := testhelpers.Value[*instance.ServerWithWarningsResponse](
						t,
						ctx.Result,
					).Server
					assert.NotNil(t, server.PublicIP)
					assert.NotEmpty(t, server.PublicIP.Address)
					assert.False(t, server.PublicIP.Dynamic)
				},
				core.TestCheckExitCode(0),
			),
			AfterFunc: deleteServerAfterFunc(),
		}))

		t.Run("with ipv6", core.Test(&core.TestConfig{
			Commands: instance.GetCommands(),
			Cmd: testServerCommand(
				"ip=ipv6 dynamic-ip-required=false -w",
			), // IPv6 is created at runtime
			Check: core.TestCheckCombine(
				func(t *testing.T, ctx *core.CheckFuncCtx) {
					t.Helper()
					require.NotNil(t, ctx.Result, "Server is nil")
					server := testhelpers.Value[*instance.ServerWithWarningsResponse](
						t,
						ctx.Result,
					).Server
					assert.Len(t, server.PublicIPs, 1)
					assert.Equal(t, instanceSDK.ServerIPIPFamilyInet6, server.PublicIPs[0].Family)
				},
				core.TestCheckExitCode(0),
			),
			AfterFunc: deleteServerAfterFunc(),
		}))

		t.Run("with ipv6 and dynamic ip", core.Test(&core.TestConfig{
			Commands: instance.GetCommands(),
			Cmd: testServerCommand(
				"dynamic-ip-required=true ip=ipv6 -w",
			), // IPv6 is created at runtime
			Check: core.TestCheckCombine(
				core.TestCheckExitCode(0),
				func(t *testing.T, ctx *core.CheckFuncCtx) {
					t.Helper()
					assert.NotNil(t, ctx.Result, "server is nil")
					server := testhelpers.Value[*instance.ServerWithWarningsResponse](
						t,
						ctx.Result,
					).Server
					assert.Len(t, server.PublicIPs, 2)
					assert.Equal(t, instanceSDK.ServerIPIPFamilyInet, server.PublicIPs[0].Family)
					assert.True(t, server.PublicIPs[0].Dynamic)
					assert.Equal(t, instanceSDK.ServerIPIPFamilyInet6, server.PublicIPs[1].Family)
				},
			),
			AfterFunc: deleteServerAfterFunc(),
		}))

		t.Run("with ipv6 and ipv4", core.Test(&core.TestConfig{
			Commands: instance.GetCommands(),
			Cmd:      testServerCommand("ip=both -w"), // IPv6 is created at runtime
			Check: core.TestCheckCombine(
				core.TestCheckExitCode(0),
				func(t *testing.T, ctx *core.CheckFuncCtx) {
					t.Helper()
					assert.NotNil(t, ctx.Result, "server is nil")
					server := testhelpers.Value[*instance.ServerWithWarningsResponse](
						t,
						ctx.Result,
					).Server
					assert.Len(t, server.PublicIPs, 2)
					assert.Equal(t, instanceSDK.ServerIPIPFamilyInet, server.PublicIPs[0].Family)
					assert.Equal(t, instanceSDK.ServerIPIPFamilyInet6, server.PublicIPs[1].Family)
				},
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
		Commands: instance.GetCommands(),
		Cmd:      testServerCommand("image=macos"),
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(1),
		),
		DisableParallel: true,
	}))

	t.Run("Error: invalid image UUID", core.Test(&core.TestConfig{
		Commands: instance.GetCommands(),
		Cmd:      testServerCommand("image=7a892c1a-bbdc-491f-9974-4008e3708664"),
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
		Commands: instance.GetCommands(),
		Cmd:      "scw instance server create type=MACBOOK1-S image=ubuntu_jammy",
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
		Commands: instance.GetCommands(),
		Cmd:      testServerCommand("image=ubuntu_jammy root-volume=l:5GB"),
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(1),
		),
		DisableParallel: true,
	}))

	t.Run("Error: invalid total local volumes size: too low 2", core.Test(&core.TestConfig{
		Commands: instance.GetCommands(),
		Cmd: testServerCommand(
			"image=ubuntu_jammy root-volume=l:5GB additional-volumes.0=block:10GB",
		),
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(1),
		),
		DisableParallel: true,
	}))

	t.Run("Error: invalid total local volumes size: too low 3", core.Test(&core.TestConfig{
		Commands:   instance.GetCommands(),
		BeforeFunc: createVolume("Volume", 5, instanceSDK.VolumeVolumeTypeLSSD),
		Cmd:        testServerCommand("image=ubuntu_jammy root-volume={{ .Volume.ID }}"),
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(1),
		),
		AfterFunc:       deleteVolume("Volume"),
		DisableParallel: true,
	}))

	t.Run("Error: invalid total local volumes size: too high 1", core.Test(&core.TestConfig{
		Commands: instance.GetCommands(),
		Cmd: testServerCommand(
			"image=ubuntu_jammy root-volume=local:10GB additional-volumes.0=local:20GB",
		),
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(1),
		),
		DisableParallel: true,
	}))

	t.Run("Error: invalid total local volumes size: too high 2", core.Test(&core.TestConfig{
		Commands: instance.GetCommands(),
		Cmd:      testServerCommand("image=ubuntu_jammy additional-volumes.0=local:30GB"),
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(1),
		),
		DisableParallel: true,
	}))

	t.Run("Error: invalid total local volumes size: too high 3", core.Test(&core.TestConfig{
		Commands:   instance.GetCommands(),
		BeforeFunc: createVolume("Volume", 20, instanceSDK.VolumeVolumeTypeLSSD),
		Cmd: testServerCommand(
			"image=ubuntu_jammy root-volume={{ .Volume.ID }} additional-volumes.0=local:10GB",
		),
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(1),
		),
		AfterFunc:       deleteVolume("Volume"),
		DisableParallel: true,
	}))

	t.Run("Error: invalid root volume size", core.Test(&core.TestConfig{
		Commands: instance.GetCommands(),
		Cmd: testServerCommand(
			"image=ubuntu_jammy root-volume=local:2GB additional-volumes.0=local:18GB",
		),
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(1),
		),
		DisableParallel: true,
	}))

	t.Run("Error: disallow existing root volume ID", core.Test(&core.TestConfig{
		Commands:   instance.GetCommands(),
		BeforeFunc: createVolume("Volume", 20, instanceSDK.VolumeVolumeTypeLSSD),
		Cmd:        testServerCommand("image=ubuntu_jammy root-volume={{ .Volume.ID }}"),
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(1),
		),
		AfterFunc:       deleteVolume("Volume"),
		DisableParallel: true,
	}))

	t.Run("Error: invalid root volume ID", core.Test(&core.TestConfig{
		Commands: instance.GetCommands(),
		Cmd: testServerCommand(
			"image=ubuntu_jammy root-volume=29da9ad9-e759-4a56-82c8-f0607f93055c",
		),
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(1),
		),
		DisableParallel: true,
	}))

	t.Run("Error: already attached additional volume ID", core.Test(&core.TestConfig{
		Commands: instance.GetCommands(),
		BeforeFunc: core.ExecStoreBeforeCmd(
			"Server",
			testServerCommand(
				"name=cli-test image=ubuntu_jammy root-volume=l:10G additional-volumes.0=l:10G stopped=true",
			),
		),
		Cmd: testServerCommand(
			`image=ubuntu_jammy root-volume=l:10G additional-volumes.0={{ (index .Server.Volumes "1").ID }} stopped=true`,
		),
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(1),
		),
		AfterFunc:       deleteServer("Server"),
		DisableParallel: true,
	}))

	t.Run("Error: invalid root volume format", core.Test(&core.TestConfig{
		Commands: instance.GetCommands(),
		Cmd:      testServerCommand("image=ubuntu_jammy root-volume=20GB"),
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(1),
		),
		DisableParallel: true,
	}))

	t.Run("Error: invalid root volume snapshot ID", core.Test(&core.TestConfig{
		Commands: instance.GetCommands(),
		Cmd: testServerCommand(
			"image=ubuntu_jammy root-volume=local:29da9ad9-e759-4a56-82c8-f0607f93055c",
		),
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(1),
		),
		DisableParallel: true,
	}))

	t.Run("Error: invalid additional volume snapshot ID", core.Test(&core.TestConfig{
		Commands: instance.GetCommands(),
		Cmd: testServerCommand(
			"image=ubuntu_jammy additional-volumes.0=block:29da9ad9-e759-4a56-82c8-f0607f93055c",
		),
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
		Commands: instance.GetCommands(),
		Cmd:      testServerCommand("image=ubuntu_jammy ip=23165951-13fd-4a3b-84ed-22c2e96658f2"),
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(1),
		),
	}))

	t.Run("Error: forbidden IP", core.Test(&core.TestConfig{
		Commands: instance.GetCommands(),
		Cmd:      testServerCommand("image=ubuntu_jammy ip=51.15.242.82"),
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(1),
		),
	}))

	t.Run("Error: invalid ip", core.Test(&core.TestConfig{
		Commands: instance.GetCommands(),
		Cmd:      testServerCommand("image=ubuntu_jammy ip=yo"),
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(1),
		),
	}))

	t.Run("Error: image size is incompatible with instance type", core.Test(&core.TestConfig{
		Commands: instance.GetCommands(),
		Cmd:      "scw instance server create image=d4067cdc-dc9d-4810-8a26-0dae51d7df42 type=DEV1-S",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(1),
		),
	}))

	////
	// Windows
	////
	t.Run("Error: ssh key id is required", core.Test(&core.TestConfig{
		Commands: instance.GetCommands(),
		Cmd:      "scw instance server create image=windows_server_2022 type=POP2-2C-8G-WIN",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(1),
		),
	}))
}

func Test_CreateServerScratchStorage(t *testing.T) {
	t.Run("Default scratch storage", core.Test(&core.TestConfig{
		Commands: instance.GetCommands(),
		Cmd:      "scw instance server create type=H100-1-80G image=ubuntu_jammy_gpu_os_12 zone=fr-par-2",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			func(_ *testing.T, ctx *core.CheckFuncCtx) {
				fmt.Println(ctx.LogBuffer)
			},
			core.TestCheckExitCode(0),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				serverResponse, isServerResponse := ctx.Result.(*instance.ServerWithWarningsResponse)
				if !isServerResponse {
					t.Fatalf("Result is not a server")
				}
				server := serverResponse.Server
				additionalVolume, exist := server.Volumes["1"]
				if !exist {
					t.Fatalf("Expected an additional scratch volume, found none")
				}
				assert.Equal(
					t,
					instanceSDK.VolumeServerVolumeTypeScratch,
					additionalVolume.VolumeType,
				)
			},
		),
		AfterFunc: core.ExecAfterCmd(
			"scw instance server delete {{ .CmdResult.ID }} zone=fr-par-2 with-volumes=all with-ip=true force-shutdown=true",
		),
		DisableParallel: true,
	}))
}

func Test_AttachFilesystem(t *testing.T) {
	t.Run("attach filesystem", core.Test(&core.TestConfig{
		Commands: core.NewCommandsMerge(
			instance.GetCommands(),
			file.GetCommands(),
		),
		BeforeFunc: core.BeforeFuncCombine(
			core.ExecStoreBeforeCmd(
				"FileSystem",
				"scw file filesystem create name=instance-fs-cli size=100000000000",
			),
			core.ExecStoreBeforeCmd(
				"Server",
				testServerCommand("stopped=true image=ubuntu-jammy type=POP2-2C-8G"),
			),
		),
		Cmd: "scw instance server attach-filesystem server-id={{ .Server.ID }} filesystem-id={{ .FileSystem.ID }}",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: core.AfterFuncCombine(
			core.ExecAfterCmd(
				"scw instance server detach-filesystem server-id={{ .Server.ID }} filesystem-id={{ .FileSystem.ID }}",
			),
			deleteServer("Server"),
			core.ExecAfterCmd(
				"scw file filesystem delete {{ .FileSystem.ID }}",
			),
		),
		DisableParallel: true,
	}))
}
