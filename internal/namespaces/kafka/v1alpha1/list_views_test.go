package kafka_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/kafka/v1alpha1"
)

func Test_ListCluster(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: kafka.GetCommands(),
		Cmd:      "scw kafka cluster list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListNodeType(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: kafka.GetCommands(),
		Cmd:      "scw kafka node-type list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListUsers(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: kafka.GetCommands(),
		Cmd:      "scw kafka users list cluster-id=11111111-1111-1111-1111-111111111111",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListVersion(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: kafka.GetCommands(),
		Cmd:      "scw kafka version list",
		Check:    core.TestCheckGolden(),
	}))
}
