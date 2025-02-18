package testhelpers

import (
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-sdk-go/api/vpcgw/v1"
)

func CreateGateway(metakey string) core.BeforeFunc {
	return func(ctx *core.BeforeFuncCtx) error {
		res := ctx.ExecuteCmd([]string{"scw", "vpc-gw", "gateway", "create", "--wait"})
		createGatewayResponse := res.(*vpcgw.Gateway)
		ctx.Meta[metakey] = createGatewayResponse

		return nil
	}
}

func CreateGatewayNetwork(metakey string) core.BeforeFunc {
	return core.ExecStoreBeforeCmd(
		"GWNT",
		"scw vpc-gw gateway-network create gateway-id={{ ."+metakey+".ID }} private-network-id={{ .PN.ID }} ipam-config.push-default-route=true --wait",
	)
}

func DeleteGatewayNetwork() core.AfterFunc {
	return core.ExecAfterCmd("scw vpc-gw gateway-network delete {{ .GWNT.ID }} --wait")
}

func DeleteGateway(metakey string) core.AfterFunc {
	return core.ExecAfterCmd("scw vpc-gw gateway delete {{ ." + metakey + ".ID }}")
}

func DeleteIPVpcGw(metakey string) core.AfterFunc {
	return core.ExecAfterCmd("scw vpc-gw ip delete {{ ." + metakey + ".IP.ID }}")
}
