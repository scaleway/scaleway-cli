package ipam

import "github.com/scaleway/scaleway-cli/v2/core"

func CreateIPFromPrivateNetwork() core.BeforeFunc {
	return core.ExecStoreBeforeCmd(
		"IPAMIP",
		"scw ipam ip create source.private-network-id={{ .PN.ID }}",
	)
}

func DeleteIP() core.AfterFunc {
	return core.ExecAfterCmd("scw ipam ip delete {{ .IPAMIP.ID }}")
}
