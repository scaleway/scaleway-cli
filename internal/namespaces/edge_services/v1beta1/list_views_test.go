package edge_services_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/edge_services/v1beta1"
)

func Test_ListBackendStage(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: edge_services.GetCommands(),
		Cmd:      "scw edge-services backend-stage list pipeline-id=11111111-1111-1111-1111-111111111111",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListCacheStage(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: edge_services.GetCommands(),
		Cmd:      "scw edge-services cache-stage list pipeline-id=11111111-1111-1111-1111-111111111111",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListDnsStage(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: edge_services.GetCommands(),
		Cmd:      "scw edge-services dns-stage list pipeline-id=11111111-1111-1111-1111-111111111111",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListPipeline(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: edge_services.GetCommands(),
		Cmd:      "scw edge-services pipeline list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListPlan(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: edge_services.GetCommands(),
		Cmd:      "scw edge-services plan list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListPurgeRequest(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: edge_services.GetCommands(),
		Cmd:      "scw edge-services purge-request list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListRouteStage(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: edge_services.GetCommands(),
		Cmd:      "scw edge-services route-stage list pipeline-id=11111111-1111-1111-1111-111111111111",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListTlsStage(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: edge_services.GetCommands(),
		Cmd:      "scw edge-services tls-stage list pipeline-id=11111111-1111-1111-1111-111111111111",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListWafStage(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: edge_services.GetCommands(),
		Cmd:      "scw edge-services waf-stage list pipeline-id=11111111-1111-1111-1111-111111111111",
		Check:    core.TestCheckGolden(),
	}))
}
