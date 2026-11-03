package datawarehouse_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/datawarehouse/v1beta1"
)

func Test_ListDatabase(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: datawarehouse.GetCommands(),
		Cmd:      "scw datawarehouse database list deployment-id=11111111-1111-1111-1111-111111111111",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListDeployment(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: datawarehouse.GetCommands(),
		Cmd:      "scw datawarehouse deployment list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListPreset(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: datawarehouse.GetCommands(),
		Cmd:      "scw datawarehouse preset list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListUser(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: datawarehouse.GetCommands(),
		Cmd:      "scw datawarehouse user list deployment-id=11111111-1111-1111-1111-111111111111",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListVersion(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: datawarehouse.GetCommands(),
		Cmd:      "scw datawarehouse version list",
		Check:    core.TestCheckGolden(),
	}))
}
