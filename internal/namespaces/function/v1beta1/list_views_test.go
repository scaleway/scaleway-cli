package function_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/function/v1beta1"
)

func Test_ListCron(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: function.GetCommands(),
		Cmd:      "scw function cron list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListDomain(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: function.GetCommands(),
		Cmd:      "scw function domain list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListFunction(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: function.GetCommands(),
		Cmd:      "scw function function list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListToken(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: function.GetCommands(),
		Cmd:      "scw function token list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListTrigger(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: function.GetCommands(),
		Cmd:      "scw function trigger list",
		Check:    core.TestCheckGolden(),
	}))
}
