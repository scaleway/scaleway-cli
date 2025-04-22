package iam_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/account/v3"
	iam "github.com/scaleway/scaleway-cli/v2/internal/namespaces/iam/v1alpha1"
	"github.com/stretchr/testify/assert"
)

func Test_getPolicyWithRules(t *testing.T) {
	t.Run("simple", core.Test(&core.TestConfig{
		Commands: core.NewCommandsMerge(
			iam.GetCommands(),
			account.GetCommands(),
		),
		BeforeFunc: core.BeforeFuncCombine(
			core.ExecStoreBeforeCmd(
				"Project",
				"scw account project create name=test-cli-get-policy",
			),
			core.ExecStoreBeforeCmd(
				"Policy",
				"scw iam policy create name=test-cli-get-policy no-principal=true rules.0.permission-set-names.0=IPAMReadOnly rules.0.project-ids.0={{ .Project.ID }}",
			),
		),
		Cmd: `scw iam policy get {{ .Policy.ID }}`,
		Check: core.TestCheckCombine(
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				assert.Contains(t, string(ctx.Stdout), "IPAMReadOnly")
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
