package config_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/config"
)

func Test_ListProfile(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: config.GetCommands(),
		Cmd:      "scw config profile list",
		Check:    core.TestCheckGolden(),
	}))
}
