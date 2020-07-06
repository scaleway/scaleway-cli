package rdb

import (
	"fmt"
	"testing"

	"github.com/scaleway/scaleway-cli/internal/core"
)

const (
	name     = "cli-test"
	user     = "foobar"
	password = "{4xdl*#QOoP+&3XRkGA)]"
	engine   = "PostgreSQL-12"
)

func Test_UpgradeInstance(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		BeforeFunc: core.ExecStoreBeforeCmd(
			"StartServer",
			fmt.Sprintf("scw rdb instance create node-type=db-dev-s is-ha-cluster=false name=%s engine=%s user-name=%s password=%s --wait", name, engine, user, password),
		),
		Cmd:       "scw rdb instance upgrade {{ .StartServer.ID }} node-type=db-dev-m --wait",
		Check:     core.TestCheckGolden(),
		AfterFunc: core.ExecAfterCmd("scw rdb instance delete {{ .StartServer.ID }}"),
	}))
}
