package testhelpers

import (
	"fmt"

	"github.com/scaleway/scaleway-cli/v2/core"
)

func CreateTemplate(metaKey string, args string) core.BeforeFunc {
	return core.ExecStoreBeforeCmd(
		metaKey,
		fmt.Sprintf("scw instance template create name=%s %s", core.GetRandomName("tpl"), args),
	)
}

func DeleteTemplate(metaKey string) core.AfterFunc {
	return func(ctx *core.AfterFuncCtx) error {
		return core.ExecAfterCmd("scw instance template delete {{ ." + metaKey + ".ID }}")(ctx)
	}
}
