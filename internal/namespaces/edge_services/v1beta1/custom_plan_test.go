package edge_services_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	edge_services "github.com/scaleway/scaleway-cli/v2/internal/namespaces/edge_services/v1beta1"
)

func Test_PlanList(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: edge_services.GetCommands(),
		Cmd:      "scw edge-services plan list",
		Check:    core.TestCheckGolden(),
	}))
}
