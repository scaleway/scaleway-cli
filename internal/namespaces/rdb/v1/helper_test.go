package rdb

import (
	"fmt"

	"github.com/scaleway/scaleway-cli/internal/core"
)

const (
	name     = "cli-test"
	user     = "foobar"
	password = "{4xdl*#QOoP+&3XRkGA)]"
	engine   = "PostgreSQL-12"
)

func createInstance(engine string) core.BeforeFunc {
	return core.ExecStoreBeforeCmd(
		"Instance",
		fmt.Sprintf("scw rdb instance create node-type=DB-DEV-S is-ha-cluster=false name=%s engine=%s user-name=%s password=%s --wait", name, engine, user, password),
	)
}

func deleteInstance() core.AfterFunc {
	return core.ExecAfterCmd("scw rdb instance delete {{ .Instance.ID }}")
}
