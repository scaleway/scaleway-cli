package vpcgw_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/vpc/v2"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/vpcgw/v1"
	"github.com/scaleway/scaleway-cli/v2/internal/testhelpers"
)

func Test_vpcGwGatewayNetworkGet(t *testing.T) {
	cmds := vpcgw.GetCommands()
	cmds.Merge(vpc.GetCommands())

	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: cmds,
		BeforeFunc: core.BeforeFuncCombine(
			testhelpers.CreatePN(),
			testhelpers.CreateGateway("GW"),
			testhelpers.CreateGatewayNetwork("GW"),
		),
		Cmd:   "scw vpc-gw gateway-network get {{ .GWNT.ID }}",
		Check: core.TestCheckGolden(),
		AfterFunc: core.AfterFuncCombine(
			testhelpers.DeleteGatewayNetwork(),
			testhelpers.DeleteGateway("GW"),
			testhelpers.DeleteIPVpcGw("GW"),
		),
	}))
}
