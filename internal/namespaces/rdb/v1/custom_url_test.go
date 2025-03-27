package rdb_test

import (
	"fmt"
	"net"
	"strconv"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/rdb/v1"
	"github.com/stretchr/testify/assert"
)

func Test_UserGetURL(t *testing.T) {
	t.Run("Postgres", core.Test(&core.TestConfig{
		Commands: rdb.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			fetchLatestEngine("PostgreSQL"), createInstance("{{.latestEngine}}"),
		),
		Cmd: "scw rdb user get-url {{ $.Instance.ID }}",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				ip := ctx.Meta["Instance"].(rdb.CreateInstanceResult).Instance.Endpoints[0].IP
				port := ctx.Meta["Instance"].(rdb.CreateInstanceResult).Instance.Endpoints[0].Port
				expected := fmt.Sprintf(
					"postgresql://%s@%s",
					user,
					net.JoinHostPort(ip.String(), strconv.Itoa(int(port))),
				)
				assert.Equal(t, expected, ctx.Result)
			},
		),
		AfterFunc: deleteInstance(),
	}))

	t.Run("MySQL", core.Test(&core.TestConfig{
		Commands: rdb.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			fetchLatestEngine("MySQL"), createInstance("{{.latestEngine}}"),
		),
		Cmd: "scw rdb user get-url {{ $.Instance.ID }}",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				ip := ctx.Meta["Instance"].(rdb.CreateInstanceResult).Instance.Endpoints[0].IP
				port := ctx.Meta["Instance"].(rdb.CreateInstanceResult).Instance.Endpoints[0].Port
				expected := fmt.Sprintf(
					"mysql://%s@%s",
					user,
					net.JoinHostPort(ip.String(), strconv.Itoa(int(port))),
				)
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
			fetchLatestEngine("PostgreSQL"),
			createInstance("{{.latestEngine}}"),
			core.ExecBeforeCmd(
				fmt.Sprintf(
					"scw rdb user create instance-id={{ $.Instance.ID }} name=%s password=%s is-admin=false",
					customUserName,
					customUserPassword,
				),
			),
		),
		Cmd: "scw rdb user get-url {{ $.Instance.ID }} user=" + customUserName,
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				ip := ctx.Meta["Instance"].(rdb.CreateInstanceResult).Instance.Endpoints[0].IP
				port := ctx.Meta["Instance"].(rdb.CreateInstanceResult).Instance.Endpoints[0].Port
				expected := fmt.Sprintf(
					"postgresql://%s@%s",
					customUserName,
					net.JoinHostPort(ip.String(), strconv.Itoa(int(port))),
				)
				assert.Equal(t, expected, ctx.Result)
			},
		),
		AfterFunc: deleteInstance(),
	}))

	t.Run("With custom database", core.Test(&core.TestConfig{
		Commands: rdb.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			fetchLatestEngine("PostgreSQL"),
			createInstance("{{.latestEngine}}"),
			core.ExecBeforeCmd(
				"scw rdb database create instance-id={{ $.Instance.ID }} name="+customDBName,
			),
		),
		Cmd: "scw rdb user get-url {{ $.Instance.ID }} db=" + customDBName,
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				ip := ctx.Meta["Instance"].(rdb.CreateInstanceResult).Instance.Endpoints[0].IP
				port := ctx.Meta["Instance"].(rdb.CreateInstanceResult).Instance.Endpoints[0].Port
				expected := fmt.Sprintf(
					"postgresql://%s@%s/%s",
					user,
					net.JoinHostPort(ip.String(), strconv.Itoa(int(port))),
					customDBName,
				)
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
			fetchLatestEngine("PostgreSQL"), createInstance("{{.latestEngine}}"),
		),
		Cmd: "scw rdb database get-url {{ $.Instance.ID }}",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				ip := ctx.Meta["Instance"].(rdb.CreateInstanceResult).Instance.Endpoints[0].IP
				port := ctx.Meta["Instance"].(rdb.CreateInstanceResult).Instance.Endpoints[0].Port
				expected := fmt.Sprintf(
					"postgresql://%s@%s",
					user,
					net.JoinHostPort(ip.String(), strconv.Itoa(int(port))),
				)
				assert.Equal(t, expected, ctx.Result)
			},
		),
		AfterFunc: deleteInstance(),
	}))
}
