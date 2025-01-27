package iot_test

import (
	"github.com/scaleway/scaleway-cli/v2/core"
)

func createHub() core.BeforeFunc {
	return core.ExecStoreBeforeCmd(
		"Hub",
		"scw iot hub create product-plan=plan_shared --wait",
	)
}

func deleteHub() core.AfterFunc {
	return core.ExecAfterCmd("scw iot hub delete {{ .Hub.ID }}")
}
