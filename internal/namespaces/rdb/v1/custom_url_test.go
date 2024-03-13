package rdb_test

import (
	"fmt"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/rdb/v1"

	"github.com/scaleway/scaleway-cli/v2/internal/core"
	"github.com/stretchr/testify/assert"
)

func Test_UserGetURL(t *testing.T) {
	t.Run("Postgres", core.Test(&core.TestConfig{
		Commands: rdb.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			createInstance("PostgreSQL-12"),
		),
		Cmd: "scw rdb user get-url {{ $.Instance.ID }}",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				ip := ctx.Meta["Instance"].(rdb.CreateInstanceResult).Instance.Endpoints[0].IP
				port := ctx.Meta["Instance"].(rdb.CreateInstanceResult).Instance.Endpoints[0].Port
				expected := fmt.Sprintf("postgresql://%s@%s:%d", user, ip, port)
				assert.Equal(t, expected, ctx.Result)
			},
		),
		AfterFunc: deleteInstance(),
	}))

	t.Run("MySQL", core.Test(&core.TestConfig{
		Commands: rdb.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			createInstance("MySQL-8"),
		),
		Cmd: "scw rdb user get-url {{ $.Instance.ID }}",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				ip := ctx.Meta["Instance"].(rdb.CreateInstanceResult).Instance.Endpoints[0].IP
				port := ctx.Meta["Instance"].(rdb.CreateInstanceResult).Instance.Endpoints[0].Port
				expected := fmt.Sprintf("mysql://%s@%s:%d", user, ip, port)
				assert.Equal(t, expected, ctx.Result)
			},
		),
		AfterFunc: deleteInstance(),
	}))

	customUserName := "custom-user"
	customUserPassword := "23uv5g%dwYIpb"
	customDBName := "custom-db"

	t.Run("With custom user", core.Test(&core.TestConfig{
		Commands: rdb.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			createInstance("PostgreSQL-12"),
			core.ExecBeforeCmd(fmt.Sprintf("scw rdb user create instance-id={{ $.Instance.ID }} name=%s password=%s is-admin=false", customUserName, customUserPassword)),
		),
		Cmd: fmt.Sprintf("scw rdb user get-url {{ $.Instance.ID }} user=%s", customUserName),
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				ip := ctx.Meta["Instance"].(rdb.CreateInstanceResult).Instance.Endpoints[0].IP
				port := ctx.Meta["Instance"].(rdb.CreateInstanceResult).Instance.Endpoints[0].Port
				expected := fmt.Sprintf("postgresql://%s@%s:%d", customUserName, ip, port)
				assert.Equal(t, expected, ctx.Result)
			},
		),
		AfterFunc: deleteInstance(),
	}))

	t.Run("With custom database", core.Test(&core.TestConfig{
		Commands: rdb.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			createInstance("PostgreSQL-12"),
			core.ExecBeforeCmd(fmt.Sprintf("scw rdb database create instance-id={{ $.Instance.ID }} name=%s", customDBName)),
		),
		Cmd: fmt.Sprintf("scw rdb user get-url {{ $.Instance.ID }} db=%s", customDBName),
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				ip := ctx.Meta["Instance"].(rdb.CreateInstanceResult).Instance.Endpoints[0].IP
				port := ctx.Meta["Instance"].(rdb.CreateInstanceResult).Instance.Endpoints[0].Port
				expected := fmt.Sprintf("postgresql://%s@%s:%d/%s", user, ip, port, customDBName)
				assert.Equal(t, expected, ctx.Result)
			},
		),
		AfterFunc: deleteInstance(),
	}))
}

func Test_DatabaseGetURL(t *testing.T) {
	t.Run("Postgres", core.Test(&core.TestConfig{
		Commands: rdb.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			createInstance("PostgreSQL-12"),
		),
		Cmd: "scw rdb database get-url {{ $.Instance.ID }}",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				ip := ctx.Meta["Instance"].(rdb.CreateInstanceResult).Instance.Endpoints[0].IP
				port := ctx.Meta["Instance"].(rdb.CreateInstanceResult).Instance.Endpoints[0].Port
				expected := fmt.Sprintf("postgresql://%s@%s:%d", user, ip, port)
				assert.Equal(t, expected, ctx.Result)
			},
		),
		AfterFunc: deleteInstance(),
	}))
}
