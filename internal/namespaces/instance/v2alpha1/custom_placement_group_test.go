package instance_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	instanceV1 "github.com/scaleway/scaleway-cli/v2/internal/namespaces/instance/v1"
	instance "github.com/scaleway/scaleway-cli/v2/internal/namespaces/instance/v2alpha1"
	instanceSDK "github.com/scaleway/scaleway-sdk-go/api/instance/v2alpha1"
	"github.com/stretchr/testify/assert"
)

func Test_PlacementGroup_Get(t *testing.T) {
	cmds := instance.GetCommands()

	t.Run("Without servers", core.Test(&core.TestConfig{
		BeforeFunc: core.BeforeFuncCombine(
			core.ExecStoreBeforeCmd(
				"PlacementGroup",
				"scw instance placement-group create",
			),
		),
		Commands: cmds,
		Cmd:      "scw instance placement-group get {{ .PlacementGroup.ID }}",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
		),
		AfterFunc: core.AfterFuncCombine(
			core.ExecAfterCmd(
				"scw instance placement-group delete {{ .PlacementGroup.ID }}",
			),
		),
	}))

	cmds.Merge(instanceV1.GetCommands())

	t.Run("With servers", core.Test(&core.TestConfig{
		BeforeFunc: core.BeforeFuncCombine(
			core.ExecStoreBeforeCmd(
				"PlacementGroup",
				"scw instance placement-group create"),
			core.ExecStoreBeforeCmd(
				"ServerA",
				"scw instance server create type=DEV1-S image=ubuntu_jammy ip=none stopped=true placement-group-id={{ .PlacementGroup.ID }}",
			),
		),
		Commands: cmds,
		Cmd:      "scw instance placement-group get {{ .PlacementGroup.ID }}",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
		),
		AfterFunc: core.AfterFuncCombine(
			core.ExecAfterCmd(
				"scw instance placement-group delete {{ .PlacementGroup.ID }}",
			),
			core.ExecAfterCmd("scw instance server delete {{ .ServerA.ID }}"),
		),
	}))
}

func Test_PlacementGroup_List(t *testing.T) {
	t.Run("simple", core.Test(&core.TestConfig{
		BeforeFunc: core.BeforeFuncCombine(
			core.ExecStoreBeforeCmd(
				"PG1",
				"scw instance placement-group create policy-type=max_availability",
			),
			core.ExecStoreBeforeCmd(
				"PG2",
				"scw instance placement-group create policy-type=low_latency",
			),
		),
		Commands: instance.GetCommands(),
		Cmd:      "scw instance placement-group list",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
		),
		AfterFunc: core.AfterFuncCombine(
			core.ExecAfterCmd("scw instance placement-group delete {{ .PG1.ID }}"),
			core.ExecAfterCmd("scw instance placement-group delete {{ .PG2.ID }}"),
		),
	}))
}

func Test_PlacementGroup_Update(t *testing.T) {
	t.Run("simple", core.Test(&core.TestConfig{
		BeforeFunc: core.BeforeFuncCombine(
			core.ExecStoreBeforeCmd(
				"PlacementGroup",
				"scw instance placement-group create policy-type=max_availability tags.0=test",
			),
		),
		Commands: instance.GetCommands(),
		Cmd: "scw instance placement-group update {{ .PlacementGroup.ID }} name=cli-test-pg-update" +
			" policy-type=low_latency tags.0=test-update tags.1=cli",
		Check: core.TestCheckCombine(
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				if pg, ok := ctx.Result.(*instanceSDK.PlacementGroup); !ok {
					t.Errorf("result was not *instance.PlacementGroup but %T", ctx.Result)
				} else {
					assert.Equal(
						t,
						pg.PolicyType,
						instanceSDK.PlacementGroupPolicyType("low_latency"),
					)
					assert.Equal(t, "cli-test-pg-update", pg.Name)
					if !assert.Len(t, pg.Tags, 2) {
						t.Fatal()
					}
					assert.Equal(t, "test-update", pg.Tags[0])
					assert.Equal(t, "cli", pg.Tags[1])
				}
			},
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
		),
		AfterFunc: core.AfterFuncCombine(
			core.ExecAfterCmd(
				"scw instance placement-group delete {{ .PlacementGroup.ID }}",
			),
		),
	}))
}
