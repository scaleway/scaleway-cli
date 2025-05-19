package rdb_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/rdb/v1"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/vpc/v2"
	rdbSDK "github.com/scaleway/scaleway-sdk-go/api/rdb/v1"
	"github.com/stretchr/testify/assert"
)

const (
	baseCommand              = "scw rdb instance create node-type=DB-DEV-S is-ha-cluster=false name=%s engine=%s user-name=%s password=%s --wait"
	privateNetworkStaticSpec = " init-endpoints.0.private-network.private-network-id={{ .PN.ID }} init-endpoints.0.private-network.service-ip={{ .IPNet }}"
	privateNetworkIpamSpec   = " init-endpoints.0.private-network.private-network-id={{ .PN.ID }} init-endpoints.0.private-network.enable-ipam=true"
	loadBalancerSpec         = " init-endpoints.1.load-balancer=true"
	publicEndpoint           = "public"
	privateEndpointIpam      = "private IPAM"
	privateEndpointStatic    = "private static"
)

func Test_ListInstance(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: rdb.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			fetchLatestEngine("PostgreSQL"),
			createInstance("{{.latestEngine}}"),
		),
		Cmd:       "scw rdb instance list",
		Check:     core.TestCheckGolden(),
		AfterFunc: deleteInstance(),
	}))
}

func Test_CloneInstance(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: rdb.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			fetchLatestEngine("PostgreSQL"),
			createInstance("{{.latestEngine}}"),
		),
		Cmd:       "scw rdb instance clone {{ .Instance.ID }} node-type=DB-DEV-M name=foobar --wait",
		Check:     core.TestCheckGolden(),
		AfterFunc: deleteInstance(),
	}))
}

func Test_CreateInstance(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: rdb.GetCommands(),
		Cmd:      fmt.Sprintf(baseCommand, name, engine, user, password),
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				instance := ctx.Result.(rdb.CreateInstanceResult).Instance
				checkEndpoints(t, ctx.Client, instance, []string{publicEndpoint})
			},
		),
		AfterFunc: core.ExecAfterCmd("scw rdb instance delete {{ .CmdResult.ID }}"),
	}))

	t.Run("With password generator", core.Test(&core.TestConfig{
		Commands: rdb.GetCommands(),
		Cmd: fmt.Sprintf(
			strings.Replace(baseCommand, "password=%s", "generate-password=true", 1),
			name,
			engine,
			user,
		),
		// do not check the golden as the password generated locally and on CI will necessarily be different
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				instance := ctx.Result.(rdb.CreateInstanceResult).Instance
				checkEndpoints(t, ctx.Client, instance, []string{publicEndpoint})
			},
		),
		AfterFunc: core.ExecAfterCmd("scw rdb instance delete {{ .CmdResult.ID }}"),
	}))
}

func Test_CreateInstanceInitEndpoints(t *testing.T) {
	cmds := rdb.GetCommands()
	cmds.Merge(vpc.GetCommands())

	t.Run("With static private endpoint", core.Test(&core.TestConfig{
		Commands:   cmds,
		BeforeFunc: createPN(),
		Cmd:        fmt.Sprintf(baseCommand+privateNetworkStaticSpec, name, engine, user, password),
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				instance := ctx.Result.(rdb.CreateInstanceResult).Instance
				checkEndpoints(t, ctx.Client, instance, []string{privateEndpointStatic})
			},
		),
		AfterFunc: core.AfterFuncCombine(
			core.ExecAfterCmd("scw rdb instance delete {{ .CmdResult.ID }} --wait"),
			deletePrivateNetwork(),
		),
	}))

	t.Run("With public and static private endpoint", core.Test(&core.TestConfig{
		Commands:   cmds,
		BeforeFunc: createPN(),
		Cmd: fmt.Sprintf(
			baseCommand+privateNetworkStaticSpec+loadBalancerSpec,
			name,
			engine,
			user,
			password,
		),
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				instance := ctx.Result.(rdb.CreateInstanceResult).Instance
				checkEndpoints(
					t,
					ctx.Client,
					instance,
					[]string{publicEndpoint, privateEndpointStatic},
				)
			},
		),
		AfterFunc: core.AfterFuncCombine(
			core.ExecAfterCmd("scw rdb instance delete {{ .CmdResult.ID }} --wait"),
			deletePrivateNetwork(),
		),
	}))

	t.Run("With IPAM private endpoint", core.Test(&core.TestConfig{
		Commands:   cmds,
		BeforeFunc: createPN(),
		Cmd:        fmt.Sprintf(baseCommand+privateNetworkIpamSpec, name, engine, user, password),
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				instance := ctx.Result.(rdb.CreateInstanceResult).Instance
				checkEndpoints(t, ctx.Client, instance, []string{privateEndpointIpam})
			},
		),
		AfterFunc: core.AfterFuncCombine(
			core.ExecAfterCmd("scw rdb instance delete {{ .CmdResult.ID }} --wait"),
			deletePrivateNetwork(),
		),
	}))

	t.Run("With public and IPAM private endpoint", core.Test(&core.TestConfig{
		Commands:   cmds,
		BeforeFunc: createPN(),
		Cmd: fmt.Sprintf(
			baseCommand+privateNetworkIpamSpec+loadBalancerSpec,
			name,
			engine,
			user,
			password,
		),
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				instance := ctx.Result.(rdb.CreateInstanceResult).Instance
				checkEndpoints(
					t,
					ctx.Client,
					instance,
					[]string{publicEndpoint, privateEndpointIpam},
				)
			},
		),
		AfterFunc: core.AfterFuncCombine(
			core.ExecAfterCmd("scw rdb instance delete {{ .CmdResult.ID }} --wait"),
			deletePrivateNetwork(),
		),
	}))
}

func Test_GetInstance(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: rdb.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			fetchLatestEngine("PostgreSQL"),
			createInstance("{{.latestEngine}}"),
		),
		Cmd:       "scw rdb instance get {{ .Instance.ID }}",
		Check:     core.TestCheckGolden(),
		AfterFunc: deleteInstance(),
	}))
}

func Test_UpgradeInstance(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: rdb.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			fetchLatestEngine("PostgreSQL"),
			createInstance("{{.latestEngine}}"),
		),
		Cmd:       "scw rdb instance upgrade {{ .Instance.ID }} node-type=DB-DEV-M --wait",
		Check:     core.TestCheckGolden(),
		AfterFunc: deleteInstance(),
	}))
}

func Test_UpdateInstance(t *testing.T) {
	t.Run("Update instance name", core.Test(&core.TestConfig{
		Commands: rdb.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			fetchLatestEngine("PostgreSQL"),
			createInstance("{{.latestEngine}}"),
		),
		Cmd: "scw rdb instance update {{ .Instance.ID }} name=foo --wait",
		Check: core.TestCheckCombine(
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				assert.Equal(t, "foo", ctx.Result.(*rdbSDK.Instance).Name)
			},
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: deleteInstance(),
	}))

	t.Run("Update instance tags", core.Test(&core.TestConfig{
		Commands: rdb.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			fetchLatestEngine("PostgreSQL"),
			createInstance("{{.latestEngine}}"),
		),
		Cmd: "scw rdb instance update {{ .Instance.ID }} tags.0=a --wait",
		Check: core.TestCheckCombine(
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				assert.Equal(t, "a", ctx.Result.(*rdbSDK.Instance).Tags[0])
			},
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: deleteInstance(),
	}))

	t.Run("Set a timezone", core.Test(&core.TestConfig{
		Commands: rdb.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			fetchLatestEngine("PostgreSQL"),
			createInstance("{{.latestEngine}}"),
		),
		Cmd: "scw rdb instance update {{ .Instance.ID }} settings.0.name=timezone settings.0.value=UTC --wait",
		Check: core.TestCheckCombine(
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				assert.Equal(t, "timezone", ctx.Result.(*rdbSDK.Instance).Settings[5].Name)
				assert.Equal(t, "UTC", ctx.Result.(*rdbSDK.Instance).Settings[5].Value)
			},
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: deleteInstance(),
	}))

	t.Run("Modify default work_mem from 4 to 8 MB", core.Test(&core.TestConfig{
		Commands: rdb.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			fetchLatestEngine("PostgreSQL"),
			createInstance("{{.latestEngine}}"),
		),
		Cmd: "scw rdb instance update {{ .Instance.ID }} settings.0.name=work_mem settings.0.value=8 --wait",
		Check: core.TestCheckCombine(
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				assert.Equal(t, "work_mem", ctx.Result.(*rdbSDK.Instance).Settings[5].Name)
				assert.Equal(t, "8", ctx.Result.(*rdbSDK.Instance).Settings[5].Value)
			},
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: deleteInstance(),
	}))

	t.Run("Modify 3 settings + add a new one", core.Test(&core.TestConfig{
		Commands: rdb.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			fetchLatestEngine("PostgreSQL"),
			createInstance("{{.latestEngine}}"),
			core.ExecBeforeCmd(
				"scw rdb instance update {{ .Instance.ID }} settings.0.name=work_mem settings.0.value=8"+
					" settings.1.name=max_connections settings.1.value=200"+
					" settings.2.name=effective_cache_size settings.2.value=1000"+
					" name=foo1 --wait",
			),
		),
		Cmd: "scw rdb instance update {{ .Instance.ID }} settings.0.name=work_mem settings.0.value=16" +
			" settings.1.name=max_connections settings.1.value=150" +
			" settings.2.name=effective_cache_size settings.2.value=1200" +
			" settings.3.name=maintenance_work_mem settings.3.value=200" +
			" name=foo2 --wait",
		Check: core.TestCheckCombine(
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				assert.Equal(
					t,
					"effective_cache_size",
					ctx.Result.(*rdbSDK.Instance).Settings[0].Name,
				)
				assert.Equal(t, "1200", ctx.Result.(*rdbSDK.Instance).Settings[0].Value)
				assert.Equal(
					t,
					"maintenance_work_mem",
					ctx.Result.(*rdbSDK.Instance).Settings[1].Name,
				)
				assert.Equal(t, "200", ctx.Result.(*rdbSDK.Instance).Settings[1].Value)
				assert.Equal(t, "max_connections", ctx.Result.(*rdbSDK.Instance).Settings[2].Name)
				assert.Equal(t, "150", ctx.Result.(*rdbSDK.Instance).Settings[2].Value)
				assert.Equal(t, "work_mem", ctx.Result.(*rdbSDK.Instance).Settings[5].Name)
				assert.Equal(t, "16", ctx.Result.(*rdbSDK.Instance).Settings[5].Value)
				assert.Equal(t, "foo2", ctx.Result.(*rdbSDK.Instance).Name)
			},
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: deleteInstance(),
	}))
}

func Test_Connect(t *testing.T) {
	t.Run("mysql", core.Test(&core.TestConfig{
		Commands: rdb.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			core.BeforeFuncStoreInMeta("username", user),
			fetchLatestEngine("MySQL"),
			createInstance("{{.latestEngine}}"),
		),
		Cmd: "scw rdb instance connect {{ .Instance.ID }} username={{ .username }}",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		OverrideExec: core.OverrideExecSimple(
			"mysql --host {{ .Instance.Endpoint.IP }} --port {{ .Instance.Endpoint.Port }} --database rdb --user {{ .username }}",
			0,
		),
		AfterFunc: deleteInstance(),
	}))

	t.Run("psql", core.Test(&core.TestConfig{
		Commands: rdb.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			core.BeforeFuncStoreInMeta("username", user),
			fetchLatestEngine("PostgreSQL"), createInstance("{{.latestEngine}}")),
		Cmd: "scw rdb instance connect {{ .Instance.ID }} username={{ .username }}",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		OverrideExec: core.OverrideExecSimple(
			"psql --host {{ .Instance.Endpoint.IP }} --port {{ .Instance.Endpoint.Port }} --username {{ .username }} --dbname rdb",
			0,
		),
		AfterFunc: deleteInstance(),
	}))
	t.Run("psql", core.Test(&core.TestConfig{
		Commands: rdb.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			core.BeforeFuncStoreInMeta("username", user),
			createPN(),
			createInstanceWithPrivateNetworkAndLoadBalancer(),
		),
		Cmd: "scw rdb instance connect {{ .Instance.ID }} username={{ .username }}",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		OverrideExec: core.OverrideExecSimple(
			"psql --host {{ .Instance.Endpoint.IP }} --port {{ .Instance.Endpoint.Port }} --username {{ .username }} --dbname rdb",
			0,
		),
		AfterFunc: deleteInstance(),
	}))
}

func deletePrivateNetwork() core.AfterFunc {
	return core.ExecAfterCmd("scw vpc private-network delete {{ .PN.ID }}")
}
