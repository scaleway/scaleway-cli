package cli

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRunRun(t *testing.T) {
	Convey("Testing runRun", t, func() {
		cmdRun.Flag.Parse([]string{"-a", "-d"})
		rawArgs := []string{"my-server"}
		So(runRun(cmdRun, rawArgs), ShouldResemble, fmt.Errorf("conflicting options: -a and -d"))
		runAttachFlag = false
		runDetachFlag = false

		cmdRun.Flag.Parse([]string{"-a"})
		rawArgs = []string{"my-server", "my-command"}
		So(runRun(cmdRun, rawArgs), ShouldResemble, fmt.Errorf("conflicting options: -a and COMMAND"))
		runAttachFlag = false

		cmdRun.Flag.Parse([]string{"-d"})
		rawArgs = []string{"my-server", "my-command"}
		So(runRun(cmdRun, rawArgs), ShouldResemble, fmt.Errorf("conflicting options: -d and COMMAND"))
		runDetachFlag = false

		cmdRun.Flag.Parse([]string{"-d", "--rm"})
		rawArgs = []string{"my-server"}
		So(runRun(cmdRun, rawArgs), ShouldResemble, fmt.Errorf("conflicting options: --detach and --rm"))
	})
}
