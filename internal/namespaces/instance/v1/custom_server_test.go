package instance_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/interactive"
	block "github.com/scaleway/scaleway-cli/v2/internal/namespaces/block/v1alpha1"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/instance/v1"
	"github.com/scaleway/scaleway-cli/v2/internal/testhelpers"
	blockSDK "github.com/scaleway/scaleway-sdk-go/api/block/v1alpha1"
	instanceSDK "github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ServerVolumeUpdate(t *testing.T) {
	t.Skip("Skipping temporarily")
	t.Run("Attach", func(t *testing.T) {
		t.Run("simple block volume", core.Test(&core.TestConfig{
			Commands: core.NewCommandsMerge(
				block.GetCommands(),
				instance.GetCommands(),
			),
			BeforeFunc: core.BeforeFuncCombine(
				createServer("Server"),
				createSbsVolume("Volume", 10),
			),
			Cmd: "scw instance server attach-volume server-id={{ .Server.ID }} volume-id={{ .Volume.ID }}",
			Check: core.TestCheckCombine(
				func(t *testing.T, ctx *core.CheckFuncCtx) {
					t.Helper()
					require.NoError(t, ctx.Err)
					resp := testhelpers.Value[*instanceSDK.AttachVolumeResponse](t, ctx.Result)
					addedVolume := testhelpers.MapTValue(t, resp.Server.Volumes, "1")
					assert.Equal(
						t,
						instanceSDK.VolumeServerVolumeTypeSbsVolume,
						addedVolume.VolumeType,
					)
				},
				testAttachVolumeServerSBSVolumeSize("0", 10),
				testAttachVolumeServerSBSVolumeSize("1", 10),
			),
			AfterFunc:       deleteServer("Server"),
			DisableParallel: true,
		}))

		t.Run("simple local volume", core.Test(&core.TestConfig{
			Commands: instance.GetCommands(),
			BeforeFunc: core.BeforeFuncCombine(
				createServer("Server"),
				createVolume("Volume", 10, instanceSDK.VolumeVolumeTypeLSSD),
			),
			Cmd: "scw instance server attach-volume server-id={{ .Server.ID }} volume-id={{ .Volume.ID }}",
			Check: core.TestCheckCombine(
				func(t *testing.T, ctx *core.CheckFuncCtx) {
					t.Helper()
					require.NoError(t, ctx.Err)
					resp := testhelpers.Value[*instanceSDK.AttachVolumeResponse](t, ctx.Result)
					addedVolume := testhelpers.MapTValue(t, resp.Server.Volumes, "1")
					assert.Equal(
						t,
						10*scw.GB,
						instance.SizeValue(addedVolume.Size),
						"Size of volume should be 10 GB",
					)
					assert.Equal(t, instanceSDK.VolumeServerVolumeTypeLSSD, addedVolume.VolumeType)
				},
				testAttachVolumeServerSBSVolumeSize("0", 10),
			),
			AfterFunc:       deleteServer("Server"),
			DisableParallel: true,
		}))

		t.Run("invalid volume UUID", core.Test(&core.TestConfig{
			Commands:   instance.GetCommands(),
			BeforeFunc: createServer("Server"),
			Cmd:        "scw instance server attach-volume server-id={{ .Server.ID }} volume-id=11111111-1111-1111-1111-111111111111",
			Check: core.TestCheckCombine(
				core.TestCheckGolden(),
				core.TestCheckExitCode(1),
			),
			AfterFunc:       deleteServer("Server"),
			DisableParallel: true,
		}))
	})
	t.Run("Detach", func(t *testing.T) {
		t.Run("simple block volume", core.Test(&core.TestConfig{
			Commands: core.NewCommandsMerge(
				block.GetCommands(),
				instance.GetCommands(),
			),
			BeforeFunc: core.ExecStoreBeforeCmd(
				"Server",
				testServerCommand("stopped=true image=ubuntu-jammy additional-volumes.0=block:10G"),
			),
			Cmd: `scw instance server detach-volume volume-id={{ (index .Server.Volumes "1").ID }} server-id={{ .Server.ID }}`,
			Check: func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				require.NoError(t, ctx.Err)
				resp := testhelpers.Value[*instanceSDK.DetachServerVolumeResponse](t, ctx.Result)
				assert.NotZero(t, resp.Server.Volumes["0"])
				assert.Nil(t, resp.Server.Volumes["1"])
				assert.Len(
					t,
					ctx.Result.(*instanceSDK.DetachServerVolumeResponse).Server.Volumes,
					1,
				)
			},
			AfterFunc: core.AfterFuncCombine(
				core.ExecAfterCmd(
					`scw block volume wait terminal-status=available {{ (index .Server.Volumes "1").ID }}`,
				),
				core.ExecAfterCmd(`scw block volume delete {{ (index .Server.Volumes "1").ID }}`),
				deleteServer("Server"),
			),
			DisableParallel: true,
		}))

		t.Run("invalid volume UUID", core.Test(&core.TestConfig{
			Commands:   instance.GetCommands(),
			BeforeFunc: createServer("Server"),
			Cmd:        "scw instance server detach-volume volume-id=11111111-1111-1111-1111-111111111111 server-id={{ .Server.ID }}",
			Check: core.TestCheckCombine(
				core.TestCheckGolden(),
				core.TestCheckExitCode(1),
			),
			AfterFunc:       deleteServer("Server"),
			DisableParallel: true,
		}))
	})
}

func Test_ServerUpdateCustom(t *testing.T) {
	// IP cases.
	t.Run("Try to remove ip from server without ip", core.Test(&core.TestConfig{
		Commands:   instance.GetCommands(),
		BeforeFunc: createServer("Server"),
		Cmd:        "scw instance server update {{ .Server.ID }} ip=none",
		Check: core.TestCheckCombine(
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				resp := testhelpers.Value[*instance.ServerWithWarningsResponse](t, ctx.Result)
				assert.Equal(t, (*instanceSDK.ServerIP)(nil), resp.Server.PublicIP)
			},
			core.TestCheckExitCode(0),
		),
		AfterFunc: deleteServer("Server"),
	}))

	t.Run("Update server ip from server without ip", core.Test(&core.TestConfig{
		Commands: instance.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			createServer("Server"),
			createIP("IP"),
		),
		Cmd: "scw instance server update {{ .Server.ID }} ip={{ .IP.Address }}",
		Check: core.TestCheckCombine(
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				ip := testhelpers.MapValue[*instanceSDK.IP](t, ctx.Meta, "IP")
				resp := testhelpers.Value[*instance.ServerWithWarningsResponse](t, ctx.Result)

				assert.NotNil(t, resp.Server)
				assert.NotNil(t, resp.Server.PublicIP)
				assert.Equal(t, ip.Address, resp.Server.PublicIP.Address)
			},
			core.TestCheckExitCode(0),
		),
		AfterFunc: deleteServer("Server"),
	}))

	t.Run("Update server ip from server with ip", core.Test(&core.TestConfig{
		Commands: instance.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			createServer("Server"),
			createIP("IP1"),
			createIP("IP2"),

			// Attach IP1 to Server.
			core.ExecStoreBeforeCmd(
				"UpdatedServer",
				"scw instance server update {{ .Server.ID }} ip={{ .IP1.Address }}",
			),
		),
		Cmd: "scw instance server update {{ .Server.ID }} ip={{ .IP2.Address }}",
		Check: core.TestCheckCombine(
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				// Test that the Server WAS attached to IP1.
				assert.Equal(
					t,
					ctx.Meta["IP1"].(*instanceSDK.IP).Address,
					ctx.Meta["UpdatedServer"].(*instance.ServerWithWarningsResponse).Server.PublicIP.Address,
				)
				// Test that the Server IS attached to IP2.
				assert.Equal(t,
					ctx.Meta["IP2"].(*instanceSDK.IP).Address,
					ctx.Result.(*instance.ServerWithWarningsResponse).Server.PublicIP.Address)
			},
			core.TestCheckExitCode(0),
		),
		AfterFunc: core.AfterFuncCombine(
			deleteServer("Server"),
			deleteIP("IP1"),
		),
	}))

	// Placement group cases.
	t.Run(
		"Update server placement-group-id from server with placement-group-id",
		core.Test(&core.TestConfig{
			Commands: instance.GetCommands(),
			BeforeFunc: core.BeforeFuncCombine(
				createPlacementGroup("PlacementGroup1"),
				createPlacementGroup("PlacementGroup2"),
				core.ExecStoreBeforeCmd(
					"Server",
					testServerCommand(
						"stopped=true image=ubuntu-jammy placement-group-id={{ .PlacementGroup1.ID }}",
					),
				),
			),
			Cmd: "scw instance server update {{ .Server.ID }} placement-group-id={{ .PlacementGroup2.ID }}",
			Check: core.TestCheckCombine(
				core.TestCheckExitCode(0),
				func(t *testing.T, ctx *core.CheckFuncCtx) {
					t.Helper()
					assert.Equal(t,
						ctx.Meta["PlacementGroup2"].(*instanceSDK.PlacementGroup).ID,
						ctx.Result.(*instance.ServerWithWarningsResponse).Server.PlacementGroup.ID)
				},
			),
			AfterFunc: core.AfterFuncCombine(
				deleteServer("Server"),
				deletePlacementGroup("PlacementGroup1"),
				deletePlacementGroup("PlacementGroup2"),
			),
		}),
	)

	// Security group cases.
	t.Run(
		"Update server security-group-id from server with security-group-id",
		core.Test(&core.TestConfig{
			Commands: instance.GetCommands(),
			BeforeFunc: core.BeforeFuncCombine(
				createSecurityGroup("SecurityGroup1"),
				createSecurityGroup("SecurityGroup2"),
				core.ExecStoreBeforeCmd(
					"Server",
					testServerCommand(
						"stopped=true image=ubuntu-jammy security-group-id={{ .SecurityGroup1.ID }}",
					),
				),
			),
			Cmd: "scw instance server update {{ .Server.ID }} security-group-id={{ .SecurityGroup2.ID }}",
			Check: core.TestCheckCombine(
				core.TestCheckExitCode(0),
				func(t *testing.T, ctx *core.CheckFuncCtx) {
					t.Helper()
					assert.Equal(t,
						ctx.Meta["SecurityGroup2"].(*instanceSDK.SecurityGroup).ID,
						ctx.Result.(*instance.ServerWithWarningsResponse).Server.SecurityGroup.ID)
				},
			),
			AfterFunc: core.AfterFuncCombine(
				deleteServer("Server"),
				deleteSecurityGroup("SecurityGroup1"),
				deleteSecurityGroup("SecurityGroup2"),
			),
		}),
	)

	// Volumes cases.
	t.Run("Volumes", func(t *testing.T) {
		t.Run("valid simple block volume", core.Test(&core.TestConfig{
			Commands: core.NewCommandsMerge(
				block.GetCommands(),
				instance.GetCommands(),
			),
			BeforeFunc: core.BeforeFuncCombine(
				createServer("Server"),
				createSbsVolume("Volume", 10),
			),
			Cmd: `scw instance server update {{ .Server.ID }} volume-ids.0={{ (index .Server.Volumes "0").ID }} volume-ids.1={{ .Volume.ID }}`,
			Check: core.TestCheckCombine(
				testServerUpdateServerSBSVolumeSize("0", 10),
				testServerUpdateServerSBSVolumeSize("1", 10),
			),
			AfterFunc: deleteServer("Server"),
		}))

		t.Run("detach all volumes", core.Test(&core.TestConfig{
			Commands: core.NewCommandsMerge(
				block.GetCommands(),
				instance.GetCommands(),
			),
			BeforeFunc: core.ExecStoreBeforeCmd(
				"Server",
				testServerCommand("stopped=true image=ubuntu-jammy additional-volumes.0=block:10G"),
			),
			Cmd: `scw instance server update {{ .Server.ID }} volume-ids=none`,
			Check: func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				require.NoError(t, ctx.Err)
				assert.Empty(t, ctx.Result.(*instance.ServerWithWarningsResponse).Server.Volumes)
			},
			AfterFunc: core.AfterFuncCombine(
				core.ExecAfterCmd(
					`scw block volume wait terminal-status=available {{ (index .Server.Volumes "0").ID }}`,
				),
				core.ExecAfterCmd(`scw block volume delete {{ (index .Server.Volumes "0").ID }}`),
				core.ExecAfterCmd(
					`scw block volume wait terminal-status=available {{ (index .Server.Volumes "1").ID }}`,
				),
				core.ExecAfterCmd(`scw block volume delete {{ (index .Server.Volumes "1").ID }}`),
				deleteServer("Server"),
			),
		}))
	})
}

// These tests needs to be run in sequence
// since they are using the interactive print
func Test_ServerDelete(t *testing.T) {
	interactive.IsInteractive = true

	t.Run("with all volumes", core.Test(&core.TestConfig{
		Commands: instance.GetCommands(),
		BeforeFunc: core.ExecStoreBeforeCmd(
			"Server",
			testServerCommand("stopped=true additional-volumes.0=block:10G"),
		),
		Cmd: `scw instance server delete {{ .Server.ID }} with-ip=true with-volumes=all`,
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		DisableParallel: true,
	}))

	t.Run("only block volumes", core.Test(&core.TestConfig{
		Commands: instance.GetCommands(),
		BeforeFunc: core.ExecStoreBeforeCmd(
			"Server",
			testServerCommand("stopped=true root-volume=l:20G additional-volumes.0=block:10G"),
		),
		Cmd: `scw instance server delete {{ .Server.ID }} with-ip=true with-volumes=block`,
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: core.ExecAfterCmd(
			`scw instance volume delete {{ (index .Server.Volumes "0").ID }}`,
		),
		DisableParallel: true,
	}))

	t.Run("only local volumes", core.Test(&core.TestConfig{
		Commands: core.NewCommandsMerge(
			block.GetCommands(),
			instance.GetCommands(),
		),
		BeforeFunc: core.ExecStoreBeforeCmd(
			"Server",
			testServerCommand("stopped=true root-volume=l:20G additional-volumes.0=block:10G"),
		),
		Cmd: `scw instance server delete {{ .Server.ID }} with-ip=true with-volumes=local`,
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: core.AfterFuncCombine(
			core.ExecAfterCmd(
				`scw block volume wait terminal-status=available {{ (index .Server.Volumes "1").ID }}`,
			),
			core.ExecAfterCmd(`scw block volume delete {{ (index .Server.Volumes "1").ID }}`),
		),
		DisableParallel: true,
	}))

	t.Run("with none volumes", core.Test(&core.TestConfig{
		Commands: core.NewCommandsMerge(
			block.GetCommands(),
			instance.GetCommands(),
		),
		BeforeFunc: core.ExecStoreBeforeCmd(
			"Server",
			testServerCommand("stopped=true root-volume=l:20G additional-volumes.0=block:10G"),
		),
		Cmd: `scw instance server delete {{ .Server.ID }} with-ip=true with-volumes=none`,
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				api := instanceSDK.NewAPI(ctx.Client)
				server := ctx.Meta["Server"].(*instance.ServerWithWarningsResponse).Server
				_, err := api.GetVolume(&instanceSDK.GetVolumeRequest{
					VolumeID: server.Volumes["0"].ID,
				})
				assert.NoError(t, err)
			},
		),
		AfterFunc: core.AfterFuncCombine(
			core.ExecAfterCmd(`scw instance volume delete {{ (index .Server.Volumes "0").ID }}`),
			core.ExecAfterCmd(`scw block volume delete {{ (index .Server.Volumes "1").ID }}`),
		),
		DisableParallel: true,
	}))

	t.Run("with sbs volumes", func(t *testing.T) {
		t.Skip("Skipping 'with sbs volumes' test temporarily")
		core.Test(&core.TestConfig{
			Commands: core.NewCommandsMerge(
				instance.GetCommands(),
				block.GetCommands(),
			),
			BeforeFunc: core.BeforeFuncCombine(
				core.ExecStoreBeforeCmd(
					"BlockVolume",
					"scw block volume create perf-iops=5000 from-empty.size=10G name=cli-test-server-delete-with-sbs-volumes",
				),
				core.ExecStoreBeforeCmd(
					"Server",
					testServerCommand("stopped=true image=ubuntu-jammy"),
				),
				core.ExecBeforeCmd(
					"scw instance server attach-volume server-id={{ .Server.ID }} volume-id={{ .BlockVolume.ID }}",
				),
			),
			Cmd: `scw instance server delete {{ .Server.ID }} with-ip=true with-volumes=all`,
			Check: core.TestCheckCombine(
				core.TestCheckGolden(),
				core.TestCheckExitCode(0),
				func(t *testing.T, ctx *core.CheckFuncCtx) {
					t.Helper()
					api := blockSDK.NewAPI(ctx.Client)
					blockVolume := ctx.Meta["BlockVolume"].(*blockSDK.Volume)
					resp, err := api.GetVolume(&blockSDK.GetVolumeRequest{
						Zone:     blockVolume.Zone,
						VolumeID: blockVolume.ID,
					})
					assert.Error(t, err, "%v", resp)
				},
			),
			DisableParallel: true,
		})
	})

	t.Run("with multiple IPs", core.Test(&core.TestConfig{
		Commands:   instance.GetCommands(),
		BeforeFunc: core.ExecStoreBeforeCmd("Server", testServerCommand("stopped=true ip=both")),
		Cmd:        `scw instance server delete {{ .Server.ID }} with-ip=true with-volumes=all`,
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()

				require.NotNil(t, ctx.Meta["Server"])
				server := ctx.Meta["Server"].(*instance.ServerWithWarningsResponse).Server
				assert.Len(t, server.PublicIPs, 2)
				api := instanceSDK.NewAPI(ctx.Client)
				for _, ip := range server.PublicIPs {
					_, err := api.GetIP(&instanceSDK.GetIPRequest{
						Zone: server.Zone,
						IP:   ip.ID,
					})
					assert.Error(t, err, "expected IP to be deleted")
				}
			},
		),
		DisableParallel: true,
	}))

	interactive.IsInteractive = false
}
