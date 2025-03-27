package feedback_test

import (
	"os/exec"
	"runtime"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/feedback"
	"github.com/stretchr/testify/assert"
)

func Test_FeedbackBug(t *testing.T) {
	t.Run("simple", core.Test(&core.TestConfig{
		Commands: feedback.GetCommands(),
		Cmd:      "scw feedback bug",
		OverrideExec: func(_ *core.ExecFuncCtx, cmd *exec.Cmd) (exitCode int, err error) {
			var observed string
			switch runtime.GOOS {
			case feedback.Windows:
				// 0: "rundll32", 1: "url.dll,FileProtocolHandler" 2: url
				observed = cmd.Args[2]
			default:
				observed = cmd.Args[1]
			}
			assert.Equal(
				t,
				"https://github.com/scaleway/scaleway-cli/issues/new?body=%0A%23%23+Description%3A%0A%0A%23%23+How+to+reproduce%3A%0A%0A%23%23%23+Command+attempted%0A%0A%23%23%23+Expected+Behavior%0A%0A%23%23%23+Actual+Behavior%0A%0A%23%23+More+info%0A%0A%23%23+Version%0A%0AVersion++++++++++0.0.0%26%2343%3Btest%0ABuildDate++++++++unknown%0AGoVersion++++++++runtime.Version%28%29%0AGitBranch++++++++unknown%0AGitCommit++++++++unknown%0AGoArch+++++++++++runtime.GOARCH%0AGoOS+++++++++++++runtime.GOOS%0AUserAgentPrefix++scaleway-cli%0A&issueTemplate=bug_report.md&labels=bug",
				observed,
			)

			return 0, nil
		},
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
	}))
}

func Test_FeedbackFeature(t *testing.T) {
	t.Run("simple", core.Test(&core.TestConfig{
		Commands: feedback.GetCommands(),
		Cmd:      "scw feedback feature",
		OverrideExec: func(_ *core.ExecFuncCtx, cmd *exec.Cmd) (exitCode int, err error) {
			var observed string
			switch runtime.GOOS {
			case feedback.Windows:
				// 0: "rundll32", 1: "url.dll,FileProtocolHandler" 2: url
				observed = cmd.Args[2]
			default:
				observed = cmd.Args[1]
			}
			assert.Equal(
				t,
				"https://github.com/scaleway/scaleway-cli/issues/new?body=%0A%23%23+Description%0A%0A%23%23+How+this+functionality+would+be+exposed%0A%0A%23%23+References%0A%0A%23%23+Version%0A%0AVersion++++++++++0.0.0%26%2343%3Btest%0ABuildDate++++++++unknown%0AGoVersion++++++++runtime.Version%28%29%0AGitBranch++++++++unknown%0AGitCommit++++++++unknown%0AGoArch+++++++++++runtime.GOARCH%0AGoOS+++++++++++++runtime.GOOS%0AUserAgentPrefix++scaleway-cli%0A&issueTemplate=feature_request.md&labels=enhancement",
				observed,
			)

			return 0, nil
		},
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
	}))
}
