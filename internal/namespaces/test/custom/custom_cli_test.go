package test_test

import (
	"testing"

	test "github.com/scaleway/scaleway-cli/v2/internal/namespaces/test/custom"

	"github.com/scaleway/scaleway-cli/v2/internal/core"
)

func Test_AnonymousFields(t *testing.T) {
	t.Run("Usage", core.Test(&core.TestConfig{
		Commands: test.GetCustomCommands(),
		Cmd:      "scw test anonymous-fields -h",
		Check:    core.TestCheckGolden(),
	}))
}
