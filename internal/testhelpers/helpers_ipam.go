package testhelpers

import "github.com/scaleway/scaleway-cli/v2/core"

func CreateIPAMIPFromPrivateNetwork() core.BeforeFunc {
	return core.ExecStoreBeforeCmd(
		"IPAMIP",
		"scw ipam ip create source.private-network-id={{ .PN.ID }}",
	)
}

func DeleteIPAMIP() core.AfterFunc {
	return core.ExecAfterCmd("scw ipam ip delete {{ .IPAMIP.ID }}")
}
