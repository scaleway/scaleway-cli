package autoscaling_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/autoscaling/v1alpha2"
)

func Test_ListAlerts(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: autoscaling.GetCommands(),
		Cmd:      "scw autoscaling alerts list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListLogs(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: autoscaling.GetCommands(),
		Cmd:      "scw autoscaling logs list 11111111-1111-1111-1111-111111111111",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListServers(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: autoscaling.GetCommands(),
		Cmd:      "scw autoscaling servers list 11111111-1111-1111-1111-111111111111",
		Check:    core.TestCheckGolden(),
	}))
}
