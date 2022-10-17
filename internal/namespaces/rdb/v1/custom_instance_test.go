package rdb

import (
	"fmt"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/scaleway/scaleway-cli/v2/internal/core"
	"github.com/scaleway/scaleway-sdk-go/api/rdb/v1"
)

func Test_ListInstance(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		BeforeFunc: createInstance("PostgreSQL-12"),
		Cmd:        "scw rdb instance list",
		Check:      core.TestCheckGolden(),
		AfterFunc:  deleteInstance(),
	}))
}

func Test_CloneInstance(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		BeforeFunc: createInstance("PostgreSQL-12"),
		Cmd:        "scw rdb instance clone {{ .Instance.ID }} node-type=DB-DEV-M name=foobar --wait",
		Check:      core.TestCheckGolden(),
		AfterFunc:  deleteInstance(),
	}))
}

func Test_CreateInstance(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands:  GetCommands(),
		Cmd:       fmt.Sprintf("scw rdb instance create node-type=DB-DEV-S is-ha-cluster=false name=%s engine=%s user-name=%s password=%s --wait", name, engine, user, password),
		Check:     core.TestCheckGolden(),
		AfterFunc: core.ExecAfterCmd("scw rdb instance delete {{ .CmdResult.ID }}"),
	}))
}

func Test_GetInstance(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		BeforeFunc: createInstance("PostgreSQL-12"),
		Cmd:        "scw rdb instance get {{ .Instance.ID }}",
		Check:      core.TestCheckGolden(),
		AfterFunc:  deleteInstance(),
	}))
}

func Test_UpgradeInstance(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		BeforeFunc: createInstance("PostgreSQL-12"),
		Cmd:        "scw rdb instance upgrade {{ .Instance.ID }} node-type=DB-DEV-M --wait",
		Check:      core.TestCheckGolden(),
		AfterFunc:  deleteInstance(),
	}))
}

func Test_UpdateInstance(t *testing.T) {
	t.Run("Update instance name", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		BeforeFunc: createInstance("PostgreSQL-12"),
		Cmd:        "scw rdb instance update {{ .Instance.ID }} name=foo --wait",
		Check: core.TestCheckCombine(
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				assert.Equal(t, "foo", ctx.Result.(*rdb.Instance).Name)
			},
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: deleteInstance(),
	}))

	t.Run("Update instance tags", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		BeforeFunc: createInstance("PostgreSQL-12"),
		Cmd:        "scw rdb instance update {{ .Instance.ID }} tags.0=a --wait",
		Check: core.TestCheckCombine(
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				assert.Equal(t, "a", ctx.Result.(*rdb.Instance).Tags[0])
			},
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: deleteInstance(),
	}))

	t.Run("Set a timezone", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		BeforeFunc: createInstance("PostgreSQL-12"),
		Cmd:        "scw rdb instance update {{ .Instance.ID }} settings.0.name=timezone settings.0.value=UTC --wait",
		Check: core.TestCheckCombine(
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				assert.Equal(t, "timezone", ctx.Result.(*rdb.Instance).Settings[6].Name)
				assert.Equal(t, "UTC", ctx.Result.(*rdb.Instance).Settings[6].Value)
			},
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: deleteInstance(),
	}))

	t.Run("Modify default work_mem from 4 to 8 MB", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		BeforeFunc: createInstance("PostgreSQL-12"),
		Cmd:        "scw rdb instance update {{ .Instance.ID }} settings.0.name=work_mem settings.0.value=8 --wait",
		Check: core.TestCheckCombine(
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				assert.Equal(t, "work_mem", ctx.Result.(*rdb.Instance).Settings[0].Name)
				assert.Equal(t, "8", ctx.Result.(*rdb.Instance).Settings[0].Value)
			},
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: deleteInstance(),
	}))

	t.Run("Modify 3 settings + add a new one", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			createInstance("PostgreSQL-12"),
			core.ExecBeforeCmd("scw rdb instance update {{ .Instance.ID }} settings.0.name=work_mem settings.0.value=8"+
				" settings.1.name=max_connections settings.1.value=200"+
				" settings.2.name=effective_cache_size settings.2.value=1000"+
				" name=foo1 --wait"),
		),
		Cmd: "scw rdb instance update {{ .Instance.ID }} settings.0.name=work_mem settings.0.value=16" +
			" settings.1.name=max_connections settings.1.value=150" +
			" settings.2.name=effective_cache_size settings.2.value=1200" +
			" settings.3.name=maintenance_work_mem settings.3.value=200" +
			" name=foo2 --wait",
		Check: core.TestCheckCombine(
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				assert.Equal(t, "work_mem", ctx.Result.(*rdb.Instance).Settings[0].Name)
				assert.Equal(t, "16", ctx.Result.(*rdb.Instance).Settings[0].Value)
				assert.Equal(t, "max_connections", ctx.Result.(*rdb.Instance).Settings[1].Name)
				assert.Equal(t, "150", ctx.Result.(*rdb.Instance).Settings[1].Value)
				assert.Equal(t, "effective_cache_size", ctx.Result.(*rdb.Instance).Settings[2].Name)
				assert.Equal(t, "1200", ctx.Result.(*rdb.Instance).Settings[2].Value)
				assert.Equal(t, "maintenance_work_mem", ctx.Result.(*rdb.Instance).Settings[3].Name)
				assert.Equal(t, "200", ctx.Result.(*rdb.Instance).Settings[3].Value)
				assert.Equal(t, "foo2", ctx.Result.(*rdb.Instance).Name)
			},
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: deleteInstance(),
	}))
}

func Test_Connect(t *testing.T) {
	t.Run("mysql", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			core.BeforeFuncStoreInMeta("username", user),
			createInstance("MySQL-8"),
		),
		Cmd: "scw rdb instance connect {{ .Instance.ID }} username={{ .username }}",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		OverrideExec: core.OverrideExecSimple("mysql --host {{ .Instance.Endpoint.IP }} --port {{ .Instance.Endpoint.Port }} --database rdb --user {{ .username }}", 0),
		AfterFunc:    deleteInstance(),
	}))

	t.Run("psql", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			core.BeforeFuncStoreInMeta("username", user),
			createInstance("PostgreSQL-12"),
		),
		Cmd: "scw rdb instance connect {{ .Instance.ID }} username={{ .username }}",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		OverrideExec: core.OverrideExecSimple("psql --host {{ .Instance.Endpoint.IP }} --port {{ .Instance.Endpoint.Port }} --username {{ .username }} --dbname rdb", 0),
		AfterFunc:    deleteInstance(),
	}))
}
