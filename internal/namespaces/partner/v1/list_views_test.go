package partner_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/partner/v1"
)

func Test_ListOrganization(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: partner.GetCommands(),
		Cmd:      "scw partner organization list",
		Check:    core.TestCheckGolden(),
	}))
}
