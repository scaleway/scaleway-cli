package webhosting

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/internal/core"
)

func Test_ListOffer(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw webhosting offer list",
		Check:    core.TestCheckGolden(),
	}))
}
