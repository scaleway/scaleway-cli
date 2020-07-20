package lb

import (
	"github.com/scaleway/scaleway-cli/internal/core"
)

func createLB() core.BeforeFunc {
	return core.ExecStoreBeforeCmd(
		"LB",
		"scw lb lb create name=cli-test description=cli-test --wait",
	)
}

func deleteLB() core.AfterFunc {
	return core.ExecAfterCmd("scw lb lb delete {{ .LB.ID }}")
}
