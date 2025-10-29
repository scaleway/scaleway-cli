package rdb_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/rdb/v1"
	"github.com/stretchr/testify/assert"
)

func Test_AddACL(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: rdb.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			fetchLatestEngine("PostgreSQL"),
			createInstance("{{.latestEngine}}"),
		),
		Cmd: "scw rdb acl add 1.2.3.4 instance-id={{ .Instance.ID }} --wait",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				verifyACL(t, ctx, []string{"0.0.0.0/0", "1.2.3.4/32"})
			},
		),
		AfterFunc: deleteInstance(),
	}))

	t.Run("Simple with description", core.Test(&core.TestConfig{
		Commands: rdb.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			fetchLatestEngine("PostgreSQL"),
			createInstance("{{.latestEngine}}"),
		),
		Cmd: "scw rdb acl add 1.2.3.4 instance-id={{ .Instance.ID }} description=some-unique-description --wait",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				verifyACL(t, ctx, []string{"0.0.0.0/0", "1.2.3.4/32"})
			},
		),
		AfterFunc: deleteInstance(),
	}))

	t.Run("Multiple", core.Test(&core.TestConfig{
		Commands: rdb.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			fetchLatestEngine("PostgreSQL"),
			createInstance("{{.latestEngine}}"),
		),
		Cmd: "scw rdb acl add 1.2.3.4 192.168.1.0/30 10.10.10.10 instance-id={{ .Instance.ID }} --wait",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				verifyACL(
					t,
					ctx,
					[]string{"0.0.0.0/0", "1.2.3.4/32", "192.168.1.0/30", "10.10.10.10/32"},
				)
			},
		),
		AfterFunc: deleteInstance(),
	}))

	t.Run("Multiple with description", core.Test(&core.TestConfig{
		Commands: rdb.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			fetchLatestEngine("PostgreSQL"),
			createInstance("{{.latestEngine}}"),
		),
		Cmd: "scw rdb acl add 1.2.3.4 192.168.1.0/30 10.10.10.10 instance-id={{ .Instance.ID }} description=some-unique-description --wait",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				verifyACL(
					t,
					ctx,
					[]string{"0.0.0.0/0", "1.2.3.4/32", "192.168.1.0/30", "10.10.10.10/32"},
				)
			},
		),
		AfterFunc: deleteInstance(),
	}))
}

func Test_DeleteACL(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: rdb.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			fetchLatestEngine("PostgreSQL"),
			createInstance("{{.latestEngine}}"),
		),
		Cmd: "scw rdb acl delete 0.0.0.0/0 instance-id={{ .Instance.ID }} --wait",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				verifyACL(t, ctx, []string{})
			},
		),
		AfterFunc: deleteInstance(),
	}))

	t.Run("Multiple when set", core.Test(&core.TestConfig{
		Commands: rdb.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			fetchLatestEngine("PostgreSQL"),
			createInstance("{{.latestEngine}}"),
			core.ExecBeforeCmd(
				"scw rdb acl add 1.2.3.4 192.168.1.0/32 10.10.10.10 instance-id={{ .Instance.ID }} --wait",
			),
		),
		Cmd: "scw rdb acl delete 1.2.3.4/32 192.168.1.0/32 10.10.10.10/32 instance-id={{ .Instance.ID }} --wait",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				verifyACL(t, ctx, []string{"0.0.0.0/0"})
			},
		),
		AfterFunc: deleteInstance(),
	}))

	t.Run("Multiple when not set", core.Test(&core.TestConfig{
		Commands: rdb.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			fetchLatestEngine("PostgreSQL"),
			createInstance("{{.latestEngine}}"),
			core.ExecBeforeCmd(
				"scw rdb acl add 192.168.1.0/32 instance-id={{ .Instance.ID }} --wait",
			),
		),
		Cmd: "scw rdb acl delete 1.2.3.4/32 192.168.1.0/32 10.10.10.10/32 instance-id={{ .Instance.ID }} --wait",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				verifyACL(t, ctx, []string{"0.0.0.0/0"})
			},
		),
		AfterFunc: deleteInstance(),
	}))
}

func Test_SetACL(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: rdb.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			fetchLatestEngine("PostgreSQL"),
			createInstance("{{.latestEngine}}"),
		),
		Cmd: "scw rdb acl set 1.2.3.4 instance-id={{ .Instance.ID }} descriptions.0=something --wait",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				verifyACL(t, ctx, []string{"1.2.3.4/32"})
				acls := ctx.Result.(*rdb.CustomACLResult).Rules
				assert.Equal(t, "something", acls[0].Description)
			},
		),
		AfterFunc: deleteInstance(),
	}))

	t.Run("Multiple", core.Test(&core.TestConfig{
		Commands: rdb.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			fetchLatestEngine("PostgreSQL"),
			createInstance("{{.latestEngine}}"),
			core.ExecBeforeCmd(
				"scw rdb acl add 1.2.3.4 192.168.1.0/32 10.10.10.10 instance-id={{ .Instance.ID }} --wait",
			),
		),
		Cmd: "scw rdb acl set 1.2.3.4 192.168.1.0/31 11.11.11.11 instance-id={{ .Instance.ID }} descriptions.0=first descriptions.1=second descriptions.2=third --wait",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				verifyACL(t, ctx, []string{"1.2.3.4/32", "192.168.1.0/31", "11.11.11.11/32"})
				acls := ctx.Result.(*rdb.CustomACLResult).Rules
				for _, acl := range acls {
					switch acl.IP.String() {
					case "1.2.3.4/32":
						assert.Equal(t, "first", acl.Description)
					case "192.168.1.0/31":
						assert.Equal(t, "second", acl.Description)
					case "11.11.11.11/32":
						assert.Equal(t, "third", acl.Description)
					}
				}
			},
		),
		AfterFunc: deleteInstance(),
	}))

	t.Run("Multiple with individual descriptions", core.Test(&core.TestConfig{
		Commands: rdb.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			fetchLatestEngine("PostgreSQL"),
			createInstance("{{.latestEngine}}"),
		),
		Cmd: "scw rdb acl add 1.1.1.1 2.2.2.2 3.3.3.3 instance-id={{ .Instance.ID }} descriptions.0=first descriptions.1=second descriptions.2=third --wait",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				verifyACL(t, ctx, []string{"0.0.0.0/0", "1.1.1.1/32", "2.2.2.2/32", "3.3.3.3/32"})
			},
		),
		AfterFunc: deleteInstance(),
	}))

	t.Run("Multiple with partial descriptions", core.Test(&core.TestConfig{
		Commands: rdb.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			fetchLatestEngine("PostgreSQL"),
			createInstance("{{.latestEngine}}"),
		),
		Cmd: "scw rdb acl add 1.1.1.1 2.2.2.2 3.3.3.3 instance-id={{ .Instance.ID }} descriptions.0=first descriptions.1=second descriptions.2=third --wait",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				verifyACL(t, ctx, []string{"0.0.0.0/0", "1.1.1.1/32", "2.2.2.2/32", "3.3.3.3/32"})
			},
		),
		AfterFunc: deleteInstance(),
	}))

	t.Run("Multiple with general description and specific descriptions", core.Test(&core.TestConfig{
		Commands: rdb.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			fetchLatestEngine("PostgreSQL"),
			createInstance("{{.latestEngine}}"),
		),
		Cmd: "scw rdb acl add 1.1.1.1 2.2.2.2 3.3.3.3 instance-id={{ .Instance.ID }} description=default descriptions.0=first descriptions.1=second --wait",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				verifyACL(t, ctx, []string{"0.0.0.0/0", "1.1.1.1/32", "2.2.2.2/32", "3.3.3.3/32"})
			},
		),
		AfterFunc: deleteInstance(),
	}))
}

func verifyACLCustomResponse(t *testing.T, res *rdb.CustomACLResult, expectedRules []string) {
	t.Helper()
	actualRules := res.Rules

	rulesFound := map[string]bool{}
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

func verifyACL(t *testing.T, ctx *core.CheckFuncCtx, expectedRules []string) {
	t.Helper()
	switch res := ctx.Result.(type) {
	case *rdb.CustomACLResult:
		verifyACLCustomResponse(t, res, expectedRules)
	case core.MultiResults:
		verifyACLCustomResponse(t, res[len(res)-1].(*rdb.CustomACLResult), expectedRules)
	default:
		t.Errorf("action is undefined for this type")
	}
}
