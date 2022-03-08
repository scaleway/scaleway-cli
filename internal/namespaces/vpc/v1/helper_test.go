package vpc

import (
	"github.com/scaleway/scaleway-cli/internal/core"
)

func createInstance() core.BeforeFunc {
	return core.ExecStoreBeforeCmd(
		"Instance",
		"scw instance server create stopped=true image=ubuntu_focal",
	)
}

func deleteInstance() core.AfterFunc {
	return core.ExecAfterCmd("scw instance server delete {{ .Instance.ID }}")
}

func createPN() core.BeforeFunc {
	return core.ExecStoreBeforeCmd(
		"PN",
		"scw vpc private-network create",
	)
}

func deletePN() core.AfterFunc {
	return core.ExecAfterCmd("scw vpc private-network delete {{ .PN.ID }}")
}

func createNIC() core.BeforeFunc {
	return core.ExecStoreBeforeCmd(
		"NIC",
		"scw instance private-nic create server-id={{ .Instance.ID }} private-network-id={{ .PN.ID }}",
	)
}
