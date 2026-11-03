package vpcgw_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/vpcgw/v2"
)

func Test_ListGateway(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: vpcgw.GetCommands(),
		Cmd:      "scw vpc-gw gateway list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListGatewayNetwork(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: vpcgw.GetCommands(),
		Cmd:      "scw vpc-gw gateway-network list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListGatewayType(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: vpcgw.GetCommands(),
		Cmd:      "scw vpc-gw gateway-type list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListIp(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: vpcgw.GetCommands(),
		Cmd:      "scw vpc-gw ip list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListPatRule(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: vpcgw.GetCommands(),
		Cmd:      "scw vpc-gw pat-rule list",
		Check:    core.TestCheckGolden(),
	}))
}
