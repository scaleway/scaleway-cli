package s2s_vpn_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/s2s_vpn/v1alpha1"
)

func Test_ListConnection(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: s2s_vpn.GetCommands(),
		Cmd:      "scw s2s-vpn connection list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListCustomerGateway(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: s2s_vpn.GetCommands(),
		Cmd:      "scw s2s-vpn customer-gateway list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListRoutingPolicy(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: s2s_vpn.GetCommands(),
		Cmd:      "scw s2s-vpn routing-policy list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListVpnGateway(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: s2s_vpn.GetCommands(),
		Cmd:      "scw s2s-vpn vpn-gateway list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListVpnGatewayType(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: s2s_vpn.GetCommands(),
		Cmd:      "scw s2s-vpn vpn-gateway-type list",
		Check:    core.TestCheckGolden(),
	}))
}
