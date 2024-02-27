package rdb

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/internal/core"
)

func Test_AddACL(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		BeforeFunc: createInstance("PostgreSQL-12"),
		Cmd:        "scw rdb acl add 1.2.3.4 instance-id={{ .Instance.ID }} --wait",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				verifyACL(ctx, t, []string{"0.0.0.0/0", "1.2.3.4/32"})
			},
		),
		AfterFunc: deleteInstance(),
	}))

	t.Run("Multiple", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		BeforeFunc: createInstance("PostgreSQL-12"),
		Cmd:        "scw rdb acl add 1.2.3.4 192.168.1.0/30 10.10.10.10 instance-id={{ .Instance.ID }} --wait",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				verifyACL(ctx, t, []string{"0.0.0.0/0", "1.2.3.4/32", "192.168.1.0/30", "10.10.10.10/32"})
			},
		),
		AfterFunc: deleteInstance(),
	}))
}

func Test_DeleteACL(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		BeforeFunc: createInstance("PostgreSQL-12"),
		Cmd:        "scw rdb acl delete 0.0.0.0/0 instance-id={{ .Instance.ID }} --wait",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				verifyACL(ctx, t, []string{})
			},
		),
		AfterFunc: deleteInstance(),
	}))

	t.Run("Multiple when set", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			createInstance("PostgreSQL-12"),
			core.ExecBeforeCmd("scw rdb acl add 1.2.3.4 192.168.1.0/32 10.10.10.10 instance-id={{ .Instance.ID }} --wait"),
		),
		Cmd: "scw rdb acl delete 1.2.3.4/32 192.168.1.0/32 10.10.10.10/32 instance-id={{ .Instance.ID }} --wait",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				verifyACL(ctx, t, []string{"0.0.0.0/0"})
			},
		),
		AfterFunc: deleteInstance(),
	}))

	t.Run("Multiple when not set", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			createInstance("PostgreSQL-12"),
			core.ExecBeforeCmd("scw rdb acl add 192.168.1.0/32 instance-id={{ .Instance.ID }} --wait"),
		),
		Cmd: "scw rdb acl delete 1.2.3.4/32 192.168.1.0/32 10.10.10.10/32 instance-id={{ .Instance.ID }} --wait",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				verifyACL(ctx, t, []string{"0.0.0.0/0"})
			},
		),
		AfterFunc: deleteInstance(),
	}))
}

func Test_SetACL(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		BeforeFunc: createInstance("PostgreSQL-12"),
		Cmd:        "scw rdb acl set rules.0.ip=1.2.3.4 instance-id={{ .Instance.ID }} --wait",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				verifyACL(ctx, t, []string{"1.2.3.4/32"})
			},
		),
		AfterFunc: deleteInstance(),
	}))

	t.Run("Multiple", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			createInstance("PostgreSQL-12"),
			core.ExecBeforeCmd("scw rdb acl add 1.2.3.4 192.168.1.0/32 10.10.10.10 instance-id={{ .Instance.ID }} --wait"),
		),
		Cmd: "scw rdb acl set rules.0.ip=1.2.3.4 rules.1.ip=192.168.1.0/31 rules.2.ip=11.11.11.11 instance-id={{ .Instance.ID }} --wait",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				verifyACL(ctx, t, []string{"1.2.3.4/32", "192.168.1.0/31", "11.11.11.11/32"})
			},
		),
		AfterFunc: deleteInstance(),
	}))
}

func verifyACLCustomResponse(res *rdbACLCustomResult, t *testing.T, expectedRules []string) {
	actualRules := res.Rules

	var rulesFound = map[string]bool{}
	for _, expectedRule := range expectedRules {
		rulesFound[expectedRule] = false
	}

	for _, actualRule := range actualRules {
		_, ok := rulesFound[actualRule.IP.String()]
		if !ok {
			t.Errorf("found rule for %s when none was expected", actualRule.IP.String())
		} else {
			rulesFound[actualRule.IP.String()] = true
		}
	}

	for rule, found := range rulesFound {
		if found == false {
			t.Errorf("expected rule for %s, got none", rule)
		}
	}
}

func verifyACL(ctx *core.CheckFuncCtx, t *testing.T, expectedRules []string) {
	switch res := ctx.Result.(type) {
	case *rdbACLCustomResult:
		verifyACLCustomResponse(res, t, expectedRules)
	case core.MultiResults:
		verifyACLCustomResponse(res[len(res)-1].(*rdbACLCustomResult), t, expectedRules)
	default:
		t.Errorf("action is undefined for this type")
	}
}
