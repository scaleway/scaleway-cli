package redis_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/redis/v1"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/vpc/v2"
	redisSDK "github.com/scaleway/scaleway-sdk-go/api/redis/v1"
	"github.com/stretchr/testify/assert"
)

const (
	baseCommand = "scw redis cluster create --wait name=%s version=7.2.11 node-type=RED1-micro user-name=admin password=P@sSw0Rd "
	serviceIPsA = "172.16.4.1/22"
	serviceIPsB = "10.16.4.1/22"
	metaNamePNA = "PrivateNetworkA"
	metaNamePNB = "PrivateNetworkB"
)

func Test_Endpoints(t *testing.T) {
	cmds := redis.GetCommands()
	cmds.Merge(vpc.GetCommands())

	t.Run("Single public endpoint", core.Test(&core.TestConfig{
		Commands: redis.GetCommands(),
		BeforeFunc: core.BeforeFuncWhenUpdatingCassette(
			func(_ *core.BeforeFuncCtx) error {
				time.Sleep(1 * time.Minute)

				return nil
			},
		),
		Cmd: fmt.Sprintf(strings.TrimSpace(baseCommand), "1-pub-endpoint"),
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				if ctx.Result != nil {
					endpoints := ctx.Result.(*redisSDK.Cluster).Endpoints
					checkEndpoints(t, endpoints, 1, 0, 0)
				}
			},
		),
		AfterFunc: deleteCluster(),
	}))

	t.Run("Single static private endpoint", core.Test(&core.TestConfig{
		Commands: cmds,
		BeforeFunc: core.BeforeFuncCombine(
			core.ExecStoreBeforeCmd(metaNamePNA, "scw vpc private-network create"),
			core.BeforeFuncWhenUpdatingCassette(
				func(_ *core.BeforeFuncCtx) error {
					time.Sleep(1 * time.Minute)

					return nil
				}),
		),
		Cmd: fmt.Sprintf(baseCommand+
			"endpoints.0.private-network.id={{ .%s.ID }} endpoints.0.private-network.service-ips.0=%s",
			"1-static-priv-endpoint", metaNamePNA, serviceIPsA),
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				if ctx.Result != nil {
					endpoints := ctx.Result.(*redisSDK.Cluster).Endpoints
					checkEndpoints(t, endpoints, 0, 1, 0)
				}
			},
		),
		AfterFunc: core.AfterFuncCombine(
			deleteCluster(),
			deletePrivateNetwork(metaNamePNA),
		),
	}))

	t.Run("Two static private endpoints", core.Test(&core.TestConfig{
		Commands: cmds,
		BeforeFunc: core.BeforeFuncCombine(
			core.ExecStoreBeforeCmd(metaNamePNA, "scw vpc private-network create"),
			core.ExecStoreBeforeCmd(metaNamePNB, "scw vpc private-network create"),
			core.BeforeFuncWhenUpdatingCassette(
				func(_ *core.BeforeFuncCtx) error {
					time.Sleep(1 * time.Minute)

					return nil
				},
			),
		),
		Cmd: fmt.Sprintf(
			"scw redis cluster create --wait name=%s version=7.2.11 node-type=RED1-micro user-name=admin password=P@sSw0Rd --wait "+
				"endpoints.0.private-network.id={{ .%s.ID }} endpoints.0.private-network.service-ips.0=%s "+
				"endpoints.1.private-network.id={{ .%s.ID }} endpoints.1.private-network.service-ips.0=%s",
			"2-static-priv-endpoints",
			metaNamePNA,
			serviceIPsA,
			metaNamePNB,
			serviceIPsB,
		),
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				if ctx.Result != nil {
					endpoints := ctx.Result.(*redisSDK.Cluster).Endpoints
					checkEndpoints(t, endpoints, 0, 2, 0)
				}
			},
		),
		AfterFunc: core.AfterFuncCombine(
			deleteCluster(),
			deletePrivateNetwork(metaNamePNA),
			deletePrivateNetwork(metaNamePNB),
		),
	}))
}

func Test_IpamConfig(t *testing.T) {
	cmds := redis.GetCommands()
	cmds.Merge(vpc.GetCommands())

	t.Run("Single IPAM private endpoint", core.Test(&core.TestConfig{
		Commands: cmds,
		BeforeFunc: core.BeforeFuncCombine(
			core.ExecStoreBeforeCmd(metaNamePNA, "scw vpc private-network create"),
			core.BeforeFuncWhenUpdatingCassette(
				func(_ *core.BeforeFuncCtx) error {
					time.Sleep(1 * time.Minute)

					return nil
				},
			),
		),
		Cmd: fmt.Sprintf(baseCommand+
			"endpoints.0.private-network.enable-ipam=true endpoints.0.private-network.id={{ .%s.ID }}",
			"1-ipam-priv-endpoint", metaNamePNA),
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				if ctx.Result != nil {
					endpoints := ctx.Result.(*redisSDK.Cluster).Endpoints
					checkEndpoints(t, endpoints, 0, 0, 1)
				}
			},
		),
		AfterFunc: core.AfterFuncCombine(
			deleteCluster(),
			deletePrivateNetwork(metaNamePNA),
		),
	}))

	t.Run("Two IPAM private endpoints", core.Test(&core.TestConfig{
		Commands: cmds,
		BeforeFunc: core.BeforeFuncCombine(
			core.ExecStoreBeforeCmd(metaNamePNA, "scw vpc private-network create"),
			core.ExecStoreBeforeCmd(metaNamePNB, "scw vpc private-network create"),
			core.BeforeFuncWhenUpdatingCassette(
				func(_ *core.BeforeFuncCtx) error {
					time.Sleep(1 * time.Minute)

					return nil
				},
			),
		),
		Cmd: fmt.Sprintf(baseCommand+
			"endpoints.0.private-network.enable-ipam=true endpoints.0.private-network.id={{ .%s.ID }} "+
			"endpoints.1.private-network.enable-ipam=true endpoints.1.private-network.id={{ .%s.ID }}",
			"2-ipam-priv-endpoints", metaNamePNA, metaNamePNB),
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				if ctx.Result != nil {
					endpoints := ctx.Result.(*redisSDK.Cluster).Endpoints
					checkEndpoints(t, endpoints, 0, 0, 2)
				}
			},
		),
		AfterFunc: core.AfterFuncCombine(
			deleteCluster(),
			deletePrivateNetwork(metaNamePNA),
			deletePrivateNetwork(metaNamePNB),
		),
	}))

	t.Run("Both IPAM and Static private endpoints", core.Test(&core.TestConfig{
		Commands: cmds,
		BeforeFunc: core.BeforeFuncCombine(
			core.ExecStoreBeforeCmd(metaNamePNA, "scw vpc private-network create"),
			core.ExecStoreBeforeCmd(metaNamePNB, "scw vpc private-network create"),
			core.BeforeFuncWhenUpdatingCassette(
				func(_ *core.BeforeFuncCtx) error {
					time.Sleep(1 * time.Minute)

					return nil
				},
			),
		),
		Cmd: fmt.Sprintf(baseCommand+
			"endpoints.0.private-network.id={{ .%s.ID }} endpoints.0.private-network.enable-ipam=true "+
			"endpoints.1.private-network.id={{ .%s.ID }} endpoints.1.private-network.service-ips.0=%s",
			"1-ipam-1-static-priv-endpoints", metaNamePNA, metaNamePNB, serviceIPsB),
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				if ctx.Result != nil {
					endpoints := ctx.Result.(*redisSDK.Cluster).Endpoints
					checkEndpoints(t, endpoints, 0, 1, 1)
				}
			},
		),
		AfterFunc: core.AfterFuncCombine(
			deleteCluster(),
			deletePrivateNetwork(metaNamePNA),
			deletePrivateNetwork(metaNamePNB),
		),
	}))
}

func Test_EndpointsEdgeCases(t *testing.T) {
	cmds := redis.GetCommands()
	cmds.Merge(vpc.GetCommands())
	expectedError := "You must specify an ipam_config or a service_ips"

	t.Run("Private endpoint with both attributes set", core.Test(&core.TestConfig{
		Commands: cmds,
		BeforeFunc: core.BeforeFuncCombine(
			core.ExecStoreBeforeCmd(metaNamePNA, "scw vpc private-network create"),
			core.BeforeFuncWhenUpdatingCassette(
				func(_ *core.BeforeFuncCtx) error {
					time.Sleep(1 * time.Minute)

					return nil
				},
			),
		),
		Cmd: fmt.Sprintf(baseCommand+
			"endpoints.0.private-network.id={{ .%s.ID }} "+
			"endpoints.0.private-network.enable-ipam=true "+
			"endpoints.0.private-network.service-ips.0=%s",
			"private-endpoint-both", metaNamePNA, serviceIPsA),
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(1),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				assert.Contains(t, ctx.Err.Error(), expectedError)
			},
		),
		AfterFunc: core.AfterFuncCombine(
			deletePrivateNetwork(metaNamePNA),
		),
	}))

	t.Run("Private endpoint with none set", core.Test(&core.TestConfig{
		Commands: cmds,
		BeforeFunc: core.BeforeFuncCombine(
			core.ExecStoreBeforeCmd(metaNamePNA, "scw vpc private-network create"),
			core.BeforeFuncWhenUpdatingCassette(
				func(_ *core.BeforeFuncCtx) error {
					time.Sleep(1 * time.Minute)

					return nil
				},
			),
		),
		Cmd: fmt.Sprintf(baseCommand+
			"endpoints.0.private-network.id={{ .%s.ID }}",
			"private-endpoint-both", metaNamePNA),
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(1),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				assert.Contains(t, ctx.Err.Error(), expectedError)
			},
		),
		AfterFunc: core.AfterFuncCombine(
			deletePrivateNetwork(metaNamePNA),
		),
	}))
}

func deleteCluster() core.AfterFunc {
	return core.ExecAfterCmd("scw redis cluster delete {{ .CmdResult.ID }} --wait")
}

func deletePrivateNetwork(metaName string) core.AfterFunc {
	return core.ExecAfterCmd(fmt.Sprintf("scw vpc private-network delete {{ .%s.ID }}", metaName))
}

func checkEndpoints(
	t *testing.T,
	endpoints []*redisSDK.Endpoint,
	nbExpectedPub, nbExpectedPrivStatic, nbExpectedPrivIpam int,
) {
	t.Helper()
	expectedEndpoints := map[string]int{
		"public":         nbExpectedPub,
		"private-static": nbExpectedPrivStatic,
		"private-ipam":   nbExpectedPrivIpam,
	}
	for _, endpoint := range endpoints {
		switch {
		case endpoint.PrivateNetwork == nil:
			expectedEndpoints["public"]--
		case endpoint.PrivateNetwork.ProvisioningMode == redisSDK.PrivateNetworkProvisioningModeStatic:
			expectedEndpoints["private-static"]--
		case endpoint.PrivateNetwork.ProvisioningMode == redisSDK.PrivateNetworkProvisioningModeIpam:
			expectedEndpoints["private-ipam"]--
		default:
			t.Error("unknown endpoint type")
		}
	}
	ok := true
	for _, nb := range expectedEndpoints {
		if nb != 0 {
			ok = false
		}
	}
	if ok == false {
		nbActualPub := nbExpectedPub - expectedEndpoints["public"]
		nbActualPrivStatic := nbExpectedPrivStatic - expectedEndpoints["private-static"]
		nbActualPrivIpam := nbExpectedPrivIpam - expectedEndpoints["private-ipam"]
		t.Errorf(
			"expected %d public endpoint(s), %d static private endpoint(s) and %d IPAM private endpoint(s), "+
				"got respectively %d, %d and %d",
			nbExpectedPub,
			nbExpectedPrivStatic,
			nbExpectedPrivIpam,
			nbActualPub,
			nbActualPrivStatic,
			nbActualPrivIpam,
		)
	}
}
