package rdb

import (
	"fmt"
	"testing"

	"github.com/scaleway/scaleway-cli/internal/core"
)

func Test_GetInstance(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		BeforeFunc: core.ExecStoreBeforeCmd(
			"StartServer",
			fmt.Sprintf("scw rdb instance create node-type=DB-DEV-S is-ha-cluster=false name=%s engine=%s user-name=%s password=%s --wait", name, engine, user, password),
		),
		Cmd:       "scw rdb instance get {{ .StartServer.ID }}",
		Check:     core.TestCheckGolden(),
		AfterFunc: core.ExecAfterCmd("scw rdb instance delete {{ .StartServer.ID }}"),
	}))
}
