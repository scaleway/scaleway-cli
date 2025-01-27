package webhosting_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	webhosting "github.com/scaleway/scaleway-cli/v2/internal/namespaces/webhosting/v1alpha1"
)

func Test_ListOffer(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: webhosting.GetCommands(),
		Cmd:      "scw webhosting offer list",
		Check:    core.TestCheckGolden(),
	}))
}
