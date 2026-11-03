package secret_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/secret/v1beta1"
)

func Test_ListSecret(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: secret.GetCommands(),
		Cmd:      "scw secret secret list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListVersion(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: secret.GetCommands(),
		Cmd:      "scw secret version list 11111111-1111-1111-1111-111111111111",
		Check:    core.TestCheckGolden(),
	}))
}
