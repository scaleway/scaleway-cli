package rdb_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/rdb/v1"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/vpc/v2"
	"github.com/scaleway/scaleway-sdk-go/api/ipam/v1"
	rdbSDK "github.com/scaleway/scaleway-sdk-go/api/rdb/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func Test_EndpointCreate(t *testing.T) {
	cmds := rdb.GetCommands()
	cmds.Merge(vpc.GetCommands())

	t.Run("Public", core.Test(&core.TestConfig{
		Commands: cmds,
		BeforeFunc: core.BeforeFuncCombine(
			createPN(),
			createInstanceWithPrivateNetwork(),
		),
		Cmd: "scw rdb endpoint create {{ .Instance.ID }} load-balancer=true --wait",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				instance := ctx.Meta["Instance"].(rdb.CreateInstanceResult).Instance
				checkEndpoints(
					t,
					ctx.Client,
					instance,
					[]string{privateEndpointStatic, publicEndpoint},
				)
			},
		),
		AfterFunc: deleteInstance(),
	}))

	t.Run("Private static", core.Test(&core.TestConfig{
		Commands: cmds,
		BeforeFunc: core.BeforeFuncCombine(
			createPN(),
			fetchLatestEngine("PostgreSQL"), createInstance("{{.latestEngine}}"),
		),
		Cmd: "scw rdb endpoint create {{ .Instance.ID }} private-network.private-network-id={{ .PN.ID }} private-network.service-ip={{ .IPNet }} --wait",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				instance := ctx.Meta["Instance"].(rdb.CreateInstanceResult).Instance
				checkEndpoints(
					t,
					ctx.Client,
					instance,
					[]string{privateEndpointStatic, publicEndpoint},
				)
			},
		),
		AfterFunc: core.AfterFuncCombine(
			deleteInstanceAndWait(),
			deletePrivateNetwork(),
		),
	}))

	t.Run("Private IPAM", core.Test(&core.TestConfig{
		Commands: cmds,
		BeforeFunc: core.BeforeFuncCombine(
			createPN(),
			fetchLatestEngine("PostgreSQL"), createInstance("{{.latestEngine}}"),
		),
		Cmd: "scw rdb endpoint create {{ .Instance.ID }} private-network.private-network-id={{ .PN.ID }} private-network.enable-ipam=true --wait",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				instance := ctx.Meta["Instance"].(rdb.CreateInstanceResult).Instance
				checkEndpoints(
					t,
					ctx.Client,
					instance,
					[]string{privateEndpointIpam, publicEndpoint},
				)
			},
		),
		AfterFunc: core.AfterFuncCombine(
			deleteInstanceAndWait(),
			deletePrivateNetwork(),
		),
	}))
}

func Test_EndpointDelete(t *testing.T) {
	cmds := rdb.GetCommands()
	cmds.Merge(vpc.GetCommands())

	t.Run("Public", core.Test(&core.TestConfig{
		Commands: cmds,
		BeforeFunc: core.BeforeFuncCombine(
			createPN(),
			createInstanceWithPrivateNetworkAndLoadBalancer(),
			listEndpointsInMeta(),
		),
		Cmd: "scw rdb endpoint delete {{ .PublicEndpoint.ID }} instance-id={{ .Instance.ID }} --wait",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				instance := ctx.Meta["Instance"].(rdb.CreateInstanceResult).Instance
				checkEndpoints(t, ctx.Client, instance, []string{privateEndpointStatic})
			},
		),
		AfterFunc: core.AfterFuncCombine(
			deleteInstanceAndWait(),
			deletePrivateNetwork(),
		),
	}))

	t.Run("Private", core.Test(&core.TestConfig{
		Commands: cmds,
		BeforeFunc: core.BeforeFuncCombine(
			createPN(),
			createInstanceWithPrivateNetworkAndLoadBalancer(),
			listEndpointsInMeta(),
		),
		Cmd: "scw rdb endpoint delete {{ .PrivateEndpoint.ID }} instance-id={{ .Instance.ID }} --wait",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				instance := ctx.Meta["Instance"].(rdb.CreateInstanceResult).Instance
				checkEndpoints(t, ctx.Client, instance, []string{publicEndpoint})
			},
		),
		AfterFunc: core.AfterFuncCombine(
			deleteInstanceAndWait(),
			deletePrivateNetwork(),
		),
	}))

	t.Run("All", core.Test(&core.TestConfig{
		Commands: cmds,
		BeforeFunc: core.BeforeFuncCombine(
			fetchLatestEngine("PostgreSQL"), createInstance("{{.latestEngine}}"),
			listEndpointsInMeta(),
		),
		Cmd: "scw rdb endpoint delete {{ .PublicEndpoint.ID }} instance-id={{ .Instance.ID }} --wait",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				instance := ctx.Meta["Instance"].(rdb.CreateInstanceResult).Instance
				checkEndpoints(t, ctx.Client, instance, []string{})
			},
		),
		AfterFunc: core.AfterFuncCombine(
			deleteInstance(),
		),
	}))
}

func Test_EndpointGet(t *testing.T) {
	cmds := rdb.GetCommands()
	cmds.Merge(vpc.GetCommands())

	t.Run("Public", core.Test(&core.TestConfig{
		Commands: cmds,
		BeforeFunc: core.BeforeFuncCombine(
			createPN(),
			createInstanceWithPrivateNetworkAndLoadBalancer(),
			listEndpointsInMeta(),
		),
		Cmd: "scw rdb endpoint get {{ .PublicEndpoint.ID }}",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				instance := ctx.Meta["Instance"].(rdb.CreateInstanceResult).Instance
				checkEndpoints(
					t,
					ctx.Client,
					instance,
					[]string{publicEndpoint, privateEndpointStatic},
				)
			},
		),
		AfterFunc: core.AfterFuncCombine(
			deleteInstanceAndWait(),
			deletePrivateNetwork(),
		),
	}))

	t.Run("Private", core.Test(&core.TestConfig{
		Commands: cmds,
		BeforeFunc: core.BeforeFuncCombine(
			createPN(),
			createInstanceWithPrivateNetworkAndLoadBalancer(),
			listEndpointsInMeta(),
		),
		Cmd: "scw rdb endpoint get {{ .PrivateEndpoint.ID }}",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				instance := ctx.Meta["Instance"].(rdb.CreateInstanceResult).Instance
				checkEndpoints(
					t,
					ctx.Client,
					instance,
					[]string{publicEndpoint, privateEndpointStatic},
				)
			},
		),
		AfterFunc: core.AfterFuncCombine(
			deleteInstanceAndWait(),
			deletePrivateNetwork(),
		),
	}))
}

func Test_EndpointList(t *testing.T) {
	cmds := rdb.GetCommands()
	cmds.Merge(vpc.GetCommands())

	t.Run("Multiple endpoints", core.Test(&core.TestConfig{
		Commands: cmds,
		BeforeFunc: core.BeforeFuncCombine(
			createPN(),
			createInstanceWithPrivateNetworkAndLoadBalancer(),
		),
		Cmd:   "scw rdb endpoint list {{ .Instance.ID }}",
		Check: core.TestCheckGolden(),
		AfterFunc: core.AfterFuncCombine(
			deleteInstanceAndWait(),
			deletePrivateNetwork(),
		),
	}))
}

func checkEndpoints(
	t *testing.T,
	client *scw.Client,
	instance *rdbSDK.Instance,
	expected []string,
) {
	t.Helper()
	rdbAPI := rdbSDK.NewAPI(client)
	ipamAPI := ipam.NewAPI(client)
	foundEndpoints := map[string]bool{}

	// First we need to update the instance as the information comes from the test's meta and may be outdated
	instanceUpdated, err := rdbAPI.GetInstance(&rdbSDK.GetInstanceRequest{
		Region:     instance.Region,
		InstanceID: instance.ID,
	})
	if err != nil {
		t.Errorf("could not get instance %s", instance.ID)
	}
	instance = instanceUpdated

	for _, endpoint := range instance.Endpoints {
		if endpoint.LoadBalancer != nil {
			foundEndpoints[publicEndpoint] = true
		}
		if endpoint.PrivateNetwork != nil {
			ips, err := ipamAPI.ListIPs(&ipam.ListIPsRequest{
				Region:       instance.Region,
				ResourceID:   &instance.ID,
				ResourceType: "rdb_instance",
				IsIPv6:       scw.BoolPtr(false),
			}, scw.WithAllPages())
			if err != nil {
				t.Errorf("could not list IPs: %v", err)
			}
			switch ips.TotalCount {
			case 1:
				foundEndpoints[privateEndpointIpam] = true
			case 0:
				foundEndpoints[privateEndpointStatic] = true
			default:
				t.Errorf("expected no more than 1 IP for instance, got %d", ips.TotalCount)
			}
		}
	}

	// Check that every expected endpoint got found
	for _, e := range expected {
		_, ok := foundEndpoints[e]
		if !ok {
			t.Errorf("expected a %s endpoint but got none", e)
		}
		delete(foundEndpoints, e)
	}
	// Check that no unexpected endpoint was found
	if len(foundEndpoints) > 0 {
		for e := range foundEndpoints {
			t.Errorf("found a %s endpoint when none was expected", e)
		}
	}
}
