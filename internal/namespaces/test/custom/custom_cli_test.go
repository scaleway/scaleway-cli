package test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/internal/core"
)

func Test_AnonymousFields(t *testing.T) {
	t.Run("Usage", core.Test(&core.TestConfig{
		Commands: GetCustomCommands(),
		Cmd:      "scw test anonymous-fields -h",
		Check:    core.TestCheckGolden(),
	}))
}
