package vpcgw

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/internal/core"
)

func Test_ListGatewayType(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw vpc-gw gateway-type list",
		Check:    core.TestCheckGolden(),
	}))
}
