package rdb

import (
	"testing"

	"github.com/scaleway/scaleway-cli/internal/core"
)

func Test_ListUser(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		BeforeFunc: createInstance("PostgreSQL-12"),
		Cmd:        "scw rdb user list instance-id={{ .Instance.ID }}",
		Check:      core.TestCheckGolden(),
		AfterFunc:  deleteInstance(),
	}))
}
