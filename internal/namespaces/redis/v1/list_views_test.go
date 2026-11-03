package redis_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/redis/v1"
)

func Test_ListCluster(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: redis.GetCommands(),
		Cmd:      "scw redis cluster list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListNodeType(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: redis.GetCommands(),
		Cmd:      "scw redis node-type list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListVersion(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: redis.GetCommands(),
		Cmd:      "scw redis version list",
		Check:    core.TestCheckGolden(),
	}))
}
