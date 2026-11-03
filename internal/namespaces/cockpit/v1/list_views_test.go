package cockpit_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/cockpit/v1"
)

func Test_ListContactPoint(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: cockpit.GetCommands(),
		Cmd:      "scw cockpit contact-point list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListDataSource(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: cockpit.GetCommands(),
		Cmd:      "scw cockpit data-source list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListGrafanaUser(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: cockpit.GetCommands(),
		Cmd:      "scw cockpit grafana-user list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListPlan(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: cockpit.GetCommands(),
		Cmd:      "scw cockpit plan list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListProductDashboards(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: cockpit.GetCommands(),
		Cmd:      "scw cockpit product-dashboards list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListToken(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: cockpit.GetCommands(),
		Cmd:      "scw cockpit token list",
		Check:    core.TestCheckGolden(),
	}))
}
