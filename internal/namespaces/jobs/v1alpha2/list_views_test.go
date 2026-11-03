package jobs_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/jobs/v1alpha2"
)

func Test_ListDefinition(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: jobs.GetCommands(),
		Cmd:      "scw jobs definition list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListRun(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: jobs.GetCommands(),
		Cmd:      "scw jobs run list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListTrigger(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: jobs.GetCommands(),
		Cmd:      "scw jobs trigger list job-definition-id=11111111-1111-1111-1111-111111111111",
		Check:    core.TestCheckGolden(),
	}))
}
