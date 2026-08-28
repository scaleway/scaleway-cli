package autoscaling_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	autoscaling "github.com/scaleway/scaleway-cli/v2/internal/namespaces/autoscaling/v1alpha2"
	instanceCLIV2 "github.com/scaleway/scaleway-cli/v2/internal/namespaces/instance/v2alpha1"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/lb/v1"
	"github.com/scaleway/scaleway-cli/v2/internal/testhelpers"
	autoscalingSDK "github.com/scaleway/scaleway-sdk-go/api/autoscaling/v1alpha2"
	instanceSDK "github.com/scaleway/scaleway-sdk-go/api/instance/v2alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/stretchr/testify/assert"
)

func deleteGroup() core.AfterFunc {
	return func(ctx *core.AfterFuncCtx) error {
		return core.ExecAfterCmd("scw autoscaling group delete {{ .CmdResult.ID }}")(ctx)
	}
}

func Test_Group(t *testing.T) {
	cmds := autoscaling.GetCommands()
	cmds.Merge(instanceCLIV2.GetCommands())

	t.Run("Create OK", core.Test(&core.TestConfig{
		Commands:   cmds,
		BeforeFunc: testhelpers.CreateTemplate("Template", "server-type=PRO2-S"),
		Cmd: "scw autoscaling group create scaling-policy-spec.minimum-size=1 scaling-policy-spec.maximum-size=5" +
			" tags.0=cli-test-asg scaling-policy-spec.scale-out-cooldown=3m scaling-policy-spec.scale-in-step=2" +
			" scaling-policy-spec.memory-target.target-avg-percent=80 template-id={{ .Template.ID }}",
		Check: core.TestCheckCombine(
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()

				if group, ok := ctx.Result.(*autoscalingSDK.Group); !ok {
					t.Errorf("expected result of type *autoscalingSDK.Group, got %T", ctx.Result)
				} else {
					assert.Equal(t, uint32(5), group.ScalingPolicy.MaximumSize)
					assert.Equal(t, uint32(1), group.ScalingPolicy.MinimumSize)
					assert.Len(t, group.Tags, 1)
					assert.Equal(t, "cli-test-asg", group.Tags[0])
					assert.Equal(
						t,
						scw.Duration{Seconds: 180},
						*group.ScalingPolicy.ScaleOutCooldown,
					)
					assert.Equal(t, uint32(2), group.ScalingPolicy.ScaleInStep)
					assert.Equal(t, uint32(80), group.ScalingPolicy.MemoryTarget.TargetAvgPercent)
				}
			},
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: core.AfterFuncCombine(
			deleteGroup(),
			testhelpers.DeleteTemplate("Template"),
		),
	}))

	t.Run("Create fail: missing size", core.Test(&core.TestConfig{
		Commands:   cmds,
		BeforeFunc: testhelpers.CreateTemplate("Template", "server-type=PRO2-S"),
		Cmd:        "scw autoscaling group create template-id={{ .Template.ID }}",
		Check: core.TestCheckCombine(
			core.TestCheckStderrContains(`- 'scaling_policy_spec' is required`),
			core.TestCheckExitCode(1),
		),
		AfterFunc: core.AfterFuncCombine(
			testhelpers.DeleteTemplate("Template"),
		),
	}))

	t.Run("Create fail: missing target", core.Test(&core.TestConfig{
		Commands:   cmds,
		BeforeFunc: testhelpers.CreateTemplate("Template", "server-type=PRO2-S"),
		Cmd: "scw autoscaling group create template-id={{ .Template.ID }}" +
			" scaling-policy-spec.minimum-size=1 scaling-policy-spec.maximum-size=5",
		Check: core.TestCheckCombine(
			core.TestCheckStderrContains(
				`At least one argument from the 'scaling-policy-target' group is required`,
			),
			core.TestCheckExitCode(1),
		),
		AfterFunc: core.AfterFuncCombine(
			testhelpers.DeleteTemplate("Template"),
		),
	}))

	cmds.Merge(lb.GetCommands())

	t.Run("Update", core.Test(&core.TestConfig{
		Commands: cmds,
		BeforeFunc: core.BeforeFuncCombine(
			testhelpers.CreateTemplate("Tmpl1", "server-type=PRO2-S"),
			testhelpers.CreateTemplate("Tmpl2", "server-type=PRO2-L"),
			core.ExecStoreBeforeCmd(
				"Group",
				"scw autoscaling group create scaling-policy-spec.minimum-size=0 scaling-policy-spec.maximum-size=12"+
					" scaling-policy-spec.fixed-size.size=3 template-id={{ .Tmpl1.ID }}",
			),
		),
		Cmd: "scw autoscaling group update {{ .Group.ID }} scaling-policy-spec.cpu-target.target-avg-percent=70" +
			" scaling-policy-spec.maximum-size=25 tags.0=tags tags.1=were tags.2=updated" +
			" scaling-policy-spec.scale-out-step=5 template-id={{ .Tmpl2.ID }}",
		Check: core.TestCheckCombine(
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()

				if group, ok := ctx.Result.(*autoscalingSDK.Group); !ok {
					t.Errorf("expected result of type *autoscalingSDK.Group, got %T", ctx.Result)
				} else {
					tmpl2ID := ctx.Meta["Tmpl2"].(*instanceSDK.Template).ID

					assert.Equal(t, uint32(70), group.ScalingPolicy.CPUTarget.TargetAvgPercent)
					assert.Equal(t, uint32(25), group.ScalingPolicy.MaximumSize)
					assert.Equal(t, uint32(0), group.ScalingPolicy.MinimumSize)
					assert.Len(t, group.Tags, 3)
					assert.Equal(t, "tags", group.Tags[0])
					assert.Equal(t, "were", group.Tags[1])
					assert.Equal(t, "updated", group.Tags[2])
					assert.Equal(t, uint32(5), group.ScalingPolicy.ScaleOutStep)
					assert.Equal(t, tmpl2ID, group.TemplateID)
				}
			},
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: core.AfterFuncCombine(
			deleteGroup(),
			testhelpers.DeleteTemplate("Tmpl1"),
			testhelpers.DeleteTemplate("Tmpl2"),
		),
	}))
}
