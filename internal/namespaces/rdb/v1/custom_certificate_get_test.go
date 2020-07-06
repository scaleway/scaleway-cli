package rdb

import (
	"fmt"
	"testing"

	"github.com/scaleway/scaleway-cli/internal/core"
)

func Test_CertificateGet(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		BeforeFunc: core.ExecStoreBeforeCmd(
			"Instance",
			fmt.Sprintf("scw rdb instance create node-type=db-dev-s is-ha-cluster=false name=%s engine=%s user-name=%s password=%s --wait", name, engine, user, password),
		),
		Cmd:       "scw rdb certificate get {{ .Instance.ID }}",
		Check:     core.TestCheckGolden(),
		AfterFunc: core.ExecAfterCmd("scw rdb instance delete {{ .Instance.ID }}"),
	}))
}
