package vpc_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/vpc/v2"
)

func Test_ListIngressRule(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: vpc.GetCommands(),
		Cmd:      "scw vpc ingress-rule list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListPrivateNetwork(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: vpc.GetCommands(),
		Cmd:      "scw vpc private-network list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListRoute(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: vpc.GetCommands(),
		Cmd:      "scw vpc route list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListVpc(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: vpc.GetCommands(),
		Cmd:      "scw vpc vpc list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListVpcConnector(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: vpc.GetCommands(),
		Cmd:      "scw vpc vpc-connector list",
		Check:    core.TestCheckGolden(),
	}))
}
