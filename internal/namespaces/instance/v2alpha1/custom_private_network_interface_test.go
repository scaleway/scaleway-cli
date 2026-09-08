package instance_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	instanceV1 "github.com/scaleway/scaleway-cli/v2/internal/namespaces/instance/v1"
	instanceV2 "github.com/scaleway/scaleway-cli/v2/internal/namespaces/instance/v2alpha1"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/ipam/v1"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/vpc/v2"
	"github.com/scaleway/scaleway-cli/v2/internal/testhelpers"
	instance "github.com/scaleway/scaleway-sdk-go/api/instance/v2alpha1"
	ipamSDK "github.com/scaleway/scaleway-sdk-go/api/ipam/v1"
	"github.com/stretchr/testify/assert"
)

func Test_PrivateNetworkInterface_Create(t *testing.T) {
	cmds := instanceV2.GetCommands()
	cmds.MergeAll(
		instanceV1.GetCommands(),
		vpc.GetCommands(),
	)

	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: cmds,
		BeforeFunc: core.BeforeFuncCombine(
			testhelpers.CreatePN(),
			createServerV1("Server"),
		),
		Cmd: "scw instance private-network-interface create private-network-id={{ .PN.ID }} server-id={{ .Server.ID }}",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
		),
		AfterFunc: core.AfterFuncCombine(
			deleteServerV1("Server"),
			testhelpers.DeletePN(),
		),
	}))

	cmds.Merge(ipam.GetCommands())

	t.Run("Advanced", core.Test(&core.TestConfig{
		Commands: cmds,
		BeforeFunc: core.BeforeFuncCombine(
			testhelpers.CreatePN(),
			createServerV1("Server"),
			testhelpers.CreateIPAMIPFromPrivateNetwork(),
		),
		Cmd: "scw instance private-network-interface create private-network-id={{ .PN.ID }} server-id={{ .Server.ID }}" +
			" ip-ids.0={{ .IPAMIP.ID }} tags.0=cli tags.1=test",
		Check: core.TestCheckCombine(
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()

				if pni, ok := ctx.Result.(*instance.PrivateNetworkInterface); ok {
					assert.Equal(t, ctx.Meta["IPAMIP"].(*ipamSDK.IP).ID, pni.IPIDs[0])
					if assert.Len(t, pni.Tags, 2) {
						assert.Equal(t, "cli", pni.Tags[0])
						assert.Equal(t, "test", pni.Tags[1])
					}
				} else {
					t.Errorf(
						"expected result of type *instance.PrivateNetworkInterface, got %T",
						ctx.Result,
					)
				}
			},
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
		),
		AfterFunc: core.AfterFuncCombine(
			deleteServerV1("Server"),
			testhelpers.DeletePN(),
		),
	}))

	t.Run("No server", core.Test(&core.TestConfig{
		Commands: cmds,
		BeforeFunc: core.BeforeFuncCombine(
			testhelpers.CreatePN(),
		),
		Cmd: "scw instance private-network-interface create private-network-id={{ .PN.ID }}",
		Check: core.TestCheckCombine(
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()

				if pni, ok := ctx.Result.(*instance.PrivateNetworkInterface); ok {
					assert.Empty(t, pni.ServerID)
				} else {
					t.Errorf(
						"expected result of type *instance.PrivateNetworkInterface, got %T",
						ctx.Result,
					)
				}
			},
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
		),
		AfterFunc: core.AfterFuncCombine(
			testhelpers.DeletePN(),
		),
	}))
}

func Test_PrivateNetworkInterface_Update(t *testing.T) {
	cmds := instanceV2.GetCommands()
	cmds.MergeAll(
		instanceV1.GetCommands(),
		vpc.GetCommands(),
	)

	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: cmds,
		BeforeFunc: core.BeforeFuncCombine(
			testhelpers.CreatePN(),
			createServerV1("Server"),
			createPrivateNetworkInterface("Server", "PN"),
		),
		Cmd: "scw instance private-network-interface update {{ .PNI.ID }} tags.0=cli-test tags.1=update",
		Check: core.TestCheckCombine(
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()

				if pni, ok := ctx.Result.(*instance.PrivateNetworkInterface); ok {
					if assert.Len(t, pni.Tags, 2) {
						assert.Equal(t, "cli-test", pni.Tags[0])
						assert.Equal(t, "update", pni.Tags[1])
					} else {
						t.Errorf(
							"expected result of type *instance.PrivateNetworkInterface, got %T",
							ctx.Result,
						)
					}
				}
			},
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
		),
		AfterFunc: core.AfterFuncCombine(
			deleteServerV1("Server"),
			testhelpers.DeletePN(),
		),
	}))
}

func Test_PrivateNetworkInterface_List(t *testing.T) {
	cmds := instanceV2.GetCommands()
	cmds.MergeAll(
		instanceV1.GetCommands(),
		vpc.GetCommands(),
	)

	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: cmds,
		BeforeFunc: core.BeforeFuncCombine(
			testhelpers.CreatePN(),
			createServerV1("Server"),
			createPrivateNetworkInterface("Server", "PN"),
		),
		Cmd: "scw instance private-network-interface list",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
		),
		AfterFunc: core.AfterFuncCombine(
			deleteServerV1("Server"),
			testhelpers.DeletePN(),
		),
	}))

	t.Run("By server ID", core.Test(&core.TestConfig{
		Commands: cmds,
		BeforeFunc: core.BeforeFuncCombine(
			testhelpers.CreatePN(),
			createServerV1("ServerA"),
			createServerV1("ServerB"),
			createPrivateNetworkInterface("ServerA", "PN"),
		),
		Cmd: "scw instance private-network-interface list server-ids.0={{ .ServerA.ID }}",
		Check: core.TestCheckCombine(
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()

				if resp, ok := ctx.Result.([]*instance.PrivateNetworkInterfaceSummary); ok {
					assert.Len(t, resp, 1)
				} else {
					t.Errorf(
						"expected result of type []*instance.PrivateNetworkInterfaceSummary, got %T",
						ctx.Result,
					)
				}
			},
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
		),
		AfterFunc: core.AfterFuncCombine(
			deleteServerV1("ServerA"),
			deleteServerV1("ServerB"),
			testhelpers.DeletePN(),
		),
	}))
}

func Test_PrivateNetworkInterface_Get(t *testing.T) {
	cmds := instanceV2.GetCommands()
	cmds.MergeAll(
		instanceV1.GetCommands(),
		vpc.GetCommands(),
	)

	t.Run("From ID", core.Test(&core.TestConfig{
		Commands: cmds,
		BeforeFunc: core.BeforeFuncCombine(
			testhelpers.CreatePN(),
			createServerV1("Server"),
			createPrivateNetworkInterface("Server", "PN"),
		),
		Cmd: "scw instance private-network-interface get {{ .NIC.ID }}",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
		),
		AfterFunc: core.AfterFuncCombine(
			deleteServerV1("Server"),
			testhelpers.DeletePN(),
		),
	}))

	t.Run("From MAC address", core.Test(&core.TestConfig{
		Commands: cmds,
		BeforeFunc: core.BeforeFuncCombine(
			testhelpers.CreatePN(),
			createServerV1("Server"),
			createPrivateNetworkInterface("Server", "PN"),
		),
		Cmd: "scw instance private-network-interface get {{ .NIC.MacAddress }} server-id={{ .Server.ID }} private-network-id={{ .PN.ID }}",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
		),
		AfterFunc: core.AfterFuncCombine(
			deleteServerV1("Server"),
			testhelpers.DeletePN(),
		),
	}))
}
