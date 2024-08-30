package flexibleip_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/baremetal/v1"
	flexibleip "github.com/scaleway/scaleway-cli/v2/internal/namespaces/flexibleip/v1alpha1"
)

func Test_CreateFlexibleWait(t *testing.T) {
	cmds := flexibleip.GetCommands()
	cmds.Merge(baremetal.GetCommands())

	t.Run("Simple", core.Test(&core.TestConfig{
		Commands:  cmds,
		Cmd:       "scw fip ip create --wait",
		Check:     core.TestCheckGolden(),
		AfterFunc: core.ExecAfterCmd("scw fip ip delete {{ .CmdResult.ID }}"),
	}))
}
