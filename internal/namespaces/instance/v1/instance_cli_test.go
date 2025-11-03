package instance_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	block "github.com/scaleway/scaleway-cli/v2/internal/namespaces/block/v1alpha1"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/instance/v1"
	"github.com/scaleway/scaleway-cli/v2/internal/testhelpers"
	instanceSDK "github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Server
func Test_ListServer(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: instance.GetCommands(),
		Cmd:      "scw instance server list",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
	}))
}

func Test_GetServer(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands:   instance.GetCommands(),
		BeforeFunc: createServer("Server"),
		Cmd:        "scw instance server get {{ .Server.ID }}",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: deleteServer("Server"),
	}))
}

// Volume
func Test_CreateVolume(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: instance.GetCommands(),
		Cmd:      "scw instance volume create name=test volume-type=l_ssd size=20G",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				require.NotNil(t, ctx.Result)
				assert.Equal(t, "test", ctx.Result.(*instanceSDK.CreateVolumeResponse).Volume.Name)
			},
		),
		AfterFunc: core.ExecAfterCmd("scw instance volume delete {{ .CmdResult.Volume.ID }}"),
	}))

	t.Run("Bad size unit", core.Test(&core.TestConfig{
		Commands: instance.GetCommands(),
		Cmd:      "scw instance volume create name=test volume-type=l_ssd size=20",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(1),
		),
	}))
}

func Test_ServerUpdate(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands:   instance.GetCommands(),
		BeforeFunc: createServer("Server"),
		Cmd:        "scw instance server update {{ .Server.ID }}",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
		),
		AfterFunc: deleteServer("Server"),
	}))

	t.Run("Placement groups", func(t *testing.T) {

		t.Run("No initial placement group & placement-group-id=none", core.Test(&core.TestConfig{
			Commands:   instance.GetCommands(),
			BeforeFunc: createServer("Server"),
			Cmd:        "scw instance server update {{ .Server.ID }} placement-group-id=none",
			Check: core.TestCheckCombine(
				func(t *testing.T, ctx *core.CheckFuncCtx) {
					t.Helper()
					require.NoError(t, ctx.Err)
					assert.Nil(
						t,
						ctx.Result.(*instance.ServerWithWarningsResponse).Server.PlacementGroup,
					)
				},
				core.TestCheckExitCode(0),
			),
			AfterFunc: deleteServer("Server"),
		}))

		t.Run(
			`No initial placement group & placement-group-id=<existing pg id>`,
			core.Test(&core.TestConfig{
				Commands: instance.GetCommands(),
				BeforeFunc: core.BeforeFuncCombine(
					createPlacementGroup("PlacementGroup"),
					createServer("Server"),
				),
				Cmd: `scw instance server update {{ .Server.ID }} placement-group-id={{ .PlacementGroup.ID }}`,
				Check: core.TestCheckCombine(
					func(t *testing.T, ctx *core.CheckFuncCtx) {
						t.Helper()
						require.NoError(t, ctx.Err)
						assert.Equal(t,
							ctx.Meta["PlacementGroup"].(*instanceSDK.PlacementGroup).ID,
							ctx.Result.(*instance.ServerWithWarningsResponse).Server.PlacementGroup.ID,
						)
					},
					core.TestCheckExitCode(0),
				),
				AfterFunc: core.AfterFuncCombine(
					deleteServer("Server"),
					deletePlacementGroup("PlacementGroup"),
				),
			}),
		)

		t.Run(
			`No initial placement group & placement-group-id=<valid, but non existing pg id>`,
			core.Test(&core.TestConfig{
				Commands:   instance.GetCommands(),
				BeforeFunc: createServer("Server"),
				Cmd:        `scw instance server update {{ .Server.ID }} placement-group-id=11111111-1111-1111-1111-111111111111`,
				Check: core.TestCheckCombine(
					core.TestCheckExitCode(1),
					core.TestCheckGolden(),
				),
				AfterFunc: deleteServer("Server"),
			}),
		)

		t.Run(
			`No initial placement group & placement-group-id=<invalid pg id>`,
			core.Test(&core.TestConfig{
				Commands:   instance.GetCommands(),
				BeforeFunc: createServer("Server"),
				Cmd:        `scw instance server update {{ .Server.ID }} placement-group-id=1111111`,
				Check: core.TestCheckCombine(
					core.TestCheckExitCode(1),
					core.TestCheckGolden(),
				),
				AfterFunc: deleteServer("Server"),
			}),
		)

		t.Run(`Initial placement group & placement-group-id=none`, core.Test(&core.TestConfig{
			Commands: instance.GetCommands(),
			BeforeFunc: core.BeforeFuncCombine(
				createPlacementGroup("PlacementGroup"),
				core.ExecStoreBeforeCmd(
					"Server",
					testServerCommand("placement-group-id={{ .PlacementGroup.ID }} stopped=true"),
				),
			),
			Cmd: `scw instance server update {{ .Server.ID }} placement-group-id=none`,
			Check: core.TestCheckCombine(
				func(t *testing.T, ctx *core.CheckFuncCtx) {
					t.Helper()
					require.NoError(t, ctx.Err)
					assert.Nil(
						t,
						ctx.Result.(*instance.ServerWithWarningsResponse).Server.PlacementGroup,
					)
				},
				core.TestCheckExitCode(0),
			),
			AfterFunc: core.AfterFuncCombine(
				deleteServer("Server"),
				deletePlacementGroup("PlacementGroup"),
			),
		}))

		t.Run(
			`Initial placement group & placement-group-id=<current pg id>`,
			core.Test(&core.TestConfig{
				Commands: instance.GetCommands(),
				BeforeFunc: core.BeforeFuncCombine(
					createPlacementGroup("PlacementGroup"),
					core.ExecStoreBeforeCmd(
						"Server",
						testServerCommand("placement-group-id={{ .PlacementGroup.ID }} stopped=true"),
					),
				),
				Cmd: `scw instance server update {{ .Server.ID }} placement-group-id={{ .PlacementGroup.ID }}`,
				Check: core.TestCheckCombine(
					func(t *testing.T, ctx *core.CheckFuncCtx) {
						t.Helper()
						require.NoError(t, ctx.Err)
						assert.Equal(t,
							ctx.Meta["PlacementGroup"].(*instanceSDK.PlacementGroup).ID,
							ctx.Result.(*instance.ServerWithWarningsResponse).Server.PlacementGroup.ID,
						)
					},
					core.TestCheckExitCode(0),
				),
				AfterFunc: core.AfterFuncCombine(
					deleteServer("Server"),
					deletePlacementGroup("PlacementGroup"),
				),
			}),
		)

		t.Run(`Initial placement group & placement-group-id=<new pg id>`, core.Test(&core.TestConfig{
			Commands: instance.GetCommands(),
			BeforeFunc: core.BeforeFuncCombine(
				createPlacementGroup("PlacementGroup1"),
				createPlacementGroup("PlacementGroup2"),
				core.ExecStoreBeforeCmd(
					"Server",
					testServerCommand("placement-group-id={{ .PlacementGroup1.ID }} stopped=true"),
				),
			),
			Cmd: `scw instance server update {{ .Server.ID }} placement-group-id={{ .PlacementGroup2.ID }}`,
			Check: core.TestCheckCombine(
				func(t *testing.T, ctx *core.CheckFuncCtx) {
					t.Helper()
					require.NoError(t, ctx.Err)
					assert.Equal(t,
						ctx.Meta["PlacementGroup2"].(*instanceSDK.PlacementGroup).ID,
						ctx.Result.(*instance.ServerWithWarningsResponse).Server.PlacementGroup.ID,
					)
				},
				core.TestCheckExitCode(0),
			),
			AfterFunc: core.AfterFuncCombine(
				deleteServer("Server"),
				deletePlacementGroup("PlacementGroup1"),
				deletePlacementGroup("PlacementGroup2"),
			),
		}))
	})

	t.Run("IP", func(t *testing.T) {
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
	})

	t.Run("Security Groups", func(t *testing.T) {
		t.Run("Update server security-group-id from server with security-group-id",
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
	})

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

// Snapshot
func Test_SnapshotCreate(t *testing.T) {
	t.Run("simple", core.Test(&core.TestConfig{
		Commands: core.NewCommandsMerge(
			block.GetCommands(),
			instance.GetCommands(),
		),
		BeforeFunc: createServer("Server"),
		Cmd:        `scw block snapshot create volume-id={{ (index .Server.Volumes "0").ID }} -w`,
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: core.AfterFuncCombine(
			core.ExecAfterCmd("scw block snapshot delete {{ .CmdResult.ID }}"),
			deleteServer("Server"),
		),
	}))
}
