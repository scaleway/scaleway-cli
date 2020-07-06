package rdb

import (
	"testing"

	"github.com/scaleway/scaleway-cli/internal/core"
)

func Test_CertificateGet(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		BeforeFunc: core.ExecStoreBeforeCmd(
			"Instance",
			"scw rdb instance create name=foobar engine=PostgreSQL-12 user-name=foobar password=\"12345678pP.8\" node-type=db-dev-s",
		),
		Cmd:   "scw rdb certificate get {{ .Instance.ID }}",
		Check: core.TestCheckGolden(),
	}))
}
