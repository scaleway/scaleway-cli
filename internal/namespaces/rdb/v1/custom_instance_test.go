package rdb

import (
	"fmt"
	"testing"

	"github.com/scaleway/scaleway-cli/internal/core"
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
