package datalab_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/datalab/v1beta1"
)

func Test_ListClusterVersion(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: datalab.GetCommands(),
		Cmd:      "scw datalab cluster-version list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListDatalab(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: datalab.GetCommands(),
		Cmd:      "scw datalab datalab list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListNodeType(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: datalab.GetCommands(),
		Cmd:      "scw datalab node-type list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListNotebookVersion(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: datalab.GetCommands(),
		Cmd:      "scw datalab notebook-version list",
		Check:    core.TestCheckGolden(),
	}))
}
