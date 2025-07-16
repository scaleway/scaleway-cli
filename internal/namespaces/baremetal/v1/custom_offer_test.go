package baremetal_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/baremetal/v1"
)

func Test_baremetalGetOffer(t *testing.T) {
	offerID := "c5853302-63e4-40c7-a711-4a91629565c8"
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: baremetal.GetCommands(),
		Cmd:      "scw baremetal offer get " + offerID,
		Check:    core.TestCheckGolden(),
	}))
}

func Test_baremetalListOffer(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: baremetal.GetCommands(),
		Cmd:      "scw baremetal offer list",
		Check:    core.TestCheckGolden(),
	}))
}
