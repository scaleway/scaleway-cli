package iam_test

import (
	"errors"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/account/v3"
	iam "github.com/scaleway/scaleway-cli/v2/internal/namespaces/iam/v1alpha1"
	"github.com/scaleway/scaleway-cli/v2/internal/testhelpers"
	"github.com/stretchr/testify/assert"
)

func Test_createRule(t *testing.T) {
	t.Run("simple", core.Test(&core.TestConfig{
		Commands: core.NewCommandsMerge(
			iam.GetCommands(),
			account.GetCommands(),
		),
		BeforeFunc: core.BeforeFuncCombine(
			core.ExecStoreBeforeCmd(
				"Project",
				"scw account project create name=test-cli-iam-create-rule",
			),
			core.ExecStoreBeforeCmd(
				"Policy",
				"scw iam policy create name=test-cli-iam-create-rule no-principal=true rules.0.permission-set-names.0=IPAMReadOnly rules.0.project-ids.0={{ .Project.ID }}",
			),
		),
		Cmd: `scw iam rule create {{ .Policy.ID }} permission-set-names.0=VPCReadOnly project-ids.0={{ .Project.ID }}`,
		Check: core.TestCheckCombine(
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				assert.Contains(t, string(ctx.Stdout), "IPAMReadOnly")
				assert.Contains(t, string(ctx.Stdout), "VPCReadOnly")
			},
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: core.AfterFuncCombine(
			core.ExecAfterCmd("scw iam policy delete {{ .Policy.ID }}"),
			core.ExecAfterCmd("scw account project delete project-id={{ .Project.ID }}"),
		),
	}))
}

func Test_deleteRule(t *testing.T) {
	t.Run("simple", core.Test(&core.TestConfig{
		Commands: core.NewCommandsMerge(
			iam.GetCommands(),
			account.GetCommands(),
		),
		BeforeFunc: core.BeforeFuncCombine(
			core.ExecStoreBeforeCmd(
				"Project",
				"scw account project create name=test-cli-iam-delete-rule",
			),
			core.ExecStoreBeforeCmd(
				"Policy",
				"scw iam policy create name=test-cli-iam-delete-rule no-principal=true rules.0.permission-set-names.0=IPAMReadOnly rules.0.project-ids.0={{ .Project.ID }} rules.1.permission-set-names.0=VPCReadOnly rules.1.project-ids.0={{ .Project.ID }}",
			),
			core.ExecStoreBeforeCmd("Policy", "scw iam policy get {{ .Policy.ID }}"),
			func(ctx *core.BeforeFuncCtx) error {
				// Get first Rule ID
				policy := testhelpers.MapValue[*iam.PolicyGetInterceptorResponse](
					t,
					ctx.Meta,
					"Policy",
				)
				if len(policy.Rules) != 2 {
					return errors.New("expected two rules in policy")
				}
				ctx.Meta["Rule"] = policy.Rules[0]

				return nil
			},
		),
		Cmd: `scw iam rule delete {{ .Policy.ID }} rule-id={{ .Rule.ID }}`,
		Check: core.TestCheckCombine(
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				assert.NotContains(t, string(ctx.Stdout), "IPAMReadOnly")
				assert.Contains(t, string(ctx.Stdout), "VPCReadOnly")
			},
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: core.AfterFuncCombine(
			core.ExecAfterCmd("scw iam policy delete {{ .Policy.ID }}"),
			core.ExecAfterCmd("scw account project delete project-id={{ .Project.ID }}"),
		),
	}))
}
