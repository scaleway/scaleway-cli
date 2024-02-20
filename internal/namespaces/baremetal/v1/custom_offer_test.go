package baremetal

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/internal/core"
)

func Test_baremetalGetOffer(t *testing.T) {
	offerID := "c5853302-63e4-40c7-a711-4a91629565c8"
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw baremetal offer get " + offerID,
		Check:    core.TestCheckGolden(),
	}))
}
