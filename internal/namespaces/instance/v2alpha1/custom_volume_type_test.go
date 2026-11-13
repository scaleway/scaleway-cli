package instance_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/instance/v2alpha1"
)

func Test_VolumeType_List(t *testing.T) {
	t.Run("volume-type list", core.Test(&core.TestConfig{
		Commands: instance.GetCommands(),
		Cmd:      "scw instance volume-type list",
		Check:    core.TestCheckGolden(),
	}))
}
