package rdb

import (
	"testing"

	"github.com/scaleway/scaleway-cli/internal/core"
)

func Test_CertificateGet(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw rdb certificate get {{ .Instance.ID }}",
		Check:    core.TestCheckGolden(),
	}))
}
