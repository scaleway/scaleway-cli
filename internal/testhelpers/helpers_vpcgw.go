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
		"scw vpc-gw gateway-network create gateway-id={{ ."+metakey+".ID }} private-network-id={{ .PN.ID }} --wait",
	)
}

func CreateGatewayNetworkDHCP(metakey string) core.BeforeFunc {
	return core.ExecStoreBeforeCmd(
		"GWNT",
		"scw vpc-gw gateway-network create gateway-id={{ ."+metakey+".ID }} private-network-id={{ .PN.ID }} enable-dhcp=true dhcp-id={{ .DHCP.ID }} --wait",
	)
}

func CreateDHCP() core.BeforeFunc {
	return core.ExecStoreBeforeCmd(
		"DHCP",
		"scw vpc-gw dhcp create subnet=192.168.1.0/24 enable-dynamic=true",
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

func DeleteDHCP() core.AfterFunc {
	return core.ExecAfterCmd("scw vpc-gw dhcp delete {{ .DHCP.ID }}")
}
