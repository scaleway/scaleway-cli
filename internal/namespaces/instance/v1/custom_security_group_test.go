package instance

import (
	"regexp"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/internal/core"
)

func Test_SecurityGroupGet(t *testing.T) {
	t.Run("Get", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			createSecurityGroup("SecurityGroup"),
		),
		Cmd: "scw instance security-group get {{ .SecurityGroup.ID }}",
		Check: core.TestCheckCombine(
			core.TestCheckGoldenAndReplacePatterns(core.GoldenReplacement{
				// For windows tests
				Pattern:     regexp.MustCompile("notepad|vi"),
				Replacement: "vi",
			}),
			core.TestCheckExitCode(0),
		),
		AfterFunc: deleteSecurityGroup("SecurityGroup"),
	}))
}
