package instance

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/internal/core"
)

func Test_VolumeTypeList(t *testing.T) {
	t.Run("volume-type list", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw instance volume-type list",
		Check:    core.TestCheckGolden(),
	}))
}
