package test_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	test "github.com/scaleway/scaleway-cli/v2/internal/namespaces/test/custom"
)

func Test_AnonymousFields(t *testing.T) {
	t.Run("Usage", core.Test(&core.TestConfig{
		Commands: test.GetCustomCommands(),
		Cmd:      "scw test anonymous-fields -h",
		Check:    core.TestCheckGolden(),
	}))
}
