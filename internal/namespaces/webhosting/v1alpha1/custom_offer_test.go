package webhosting_test

import (
	"testing"

	webhosting "github.com/scaleway/scaleway-cli/v2/internal/namespaces/webhosting/v1alpha1"

	"github.com/scaleway/scaleway-cli/v2/internal/core"
)

func Test_ListOffer(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: webhosting.GetCommands(),
		Cmd:      "scw webhosting offer list",
		Check:    core.TestCheckGolden(),
	}))
}
