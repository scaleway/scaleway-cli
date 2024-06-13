package vpcgw_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/vpcgw/v1"

	"github.com/scaleway/scaleway-cli/v2/internal/core"
)

func Test_ListGatewayType(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: vpcgw.GetCommands(),
		Cmd:      "scw vpc-gw gateway-type list",
		Check:    core.TestCheckGolden(),
	}))
}
