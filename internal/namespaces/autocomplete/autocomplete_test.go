package autocomplete

import (
	"testing"

	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/stretchr/testify/assert"
)

func Test_Autocomplete(t *testing.T) {
	t.Run("zsh", func(t *testing.T) {
		t.Run("Level1", core.Test(&core.TestConfig{
			Commands: GetCommands(),
			Cmd:      "scw autocomplete complete zsh -- 1",
			Check: core.TestCheckCombine(
				core.TestCheckGolden(),
				core.TestCheckExitCode(0),
			),
		}))
		t.Run("Level2", core.Test(&core.TestConfig{
			Commands: GetCommands(),
			Cmd:      "scw autocomplete complete zsh -- 2 scw",
			Check: core.TestCheckCombine(
				core.TestCheckGolden(),
				core.TestCheckExitCode(0),
			),
		}))
		t.Run("Level3", core.Test(&core.TestConfig{
			Commands: GetCommands(),
			Cmd:      "scw autocomplete complete zsh -- 3 scw autocomplete",
			Check: core.TestCheckCombine(
				core.TestCheckGolden(),
				core.TestCheckExitCode(0),
			),
		}))
	})
	t.Run("fish", func(t *testing.T) {
		t.Run("Level1", core.Test(&core.TestConfig{
			Commands: GetCommands(),
			Cmd:      "scw autocomplete complete fish -- \"\" \"\" \"\"",
			Check: core.TestCheckCombine(
				core.TestCheckGolden(),
				core.TestCheckExitCode(0),
			),
		}))
		t.Run("Level2", core.Test(&core.TestConfig{
			Commands: GetCommands(),
			Cmd:      "scw autocomplete complete fish -- \"\" \"\" \"\" scw",
			Check: core.TestCheckCombine(
				core.TestCheckGolden(),
				core.TestCheckExitCode(0),
			),
		}))
		t.Run("Level3", core.Test(&core.TestConfig{
			Commands: GetCommands(),
			Cmd:      "scw autocomplete complete fish -- \"\" \"\" \"autocomplete\" scw",
			Check: core.TestCheckCombine(
				core.TestCheckGolden(),
				core.TestCheckExitCode(0),
			),
		}))
	})
	t.Run("bash", func(t *testing.T) {
		t.Run("Level1", core.Test(&core.TestConfig{
			Commands: GetCommands(),
			Cmd:      "scw autocomplete complete bash -- \"\" 0 \"\" \"\"",
			Check: core.TestCheckCombine(
				core.TestCheckGolden(),
				core.TestCheckExitCode(0),
			),
		}))
		t.Run("Level2", core.Test(&core.TestConfig{
			Commands: GetCommands(),
			Cmd:      "scw autocomplete complete bash -- \"\" 1 \"scw\" \"\"",
			Check: core.TestCheckCombine(
				core.TestCheckGolden(),
				core.TestCheckExitCode(0),
			),
		}))
		t.Run("Level3", core.Test(&core.TestConfig{
			Commands: GetCommands(),
			Cmd:      "scw autocomplete complete bash -- \"\" 2 \"scw\" \"autocomplete\" \"\"",
			Check: core.TestCheckCombine(
				core.TestCheckGolden(),
				core.TestCheckExitCode(0),
			),
		}))
	})
}

func TestTrimText(t *testing.T) {
	type testCase struct {
		Script       string
		TrimedScript string
	}

	run := func(tc *testCase) func(*testing.T) {
		return func(t *testing.T) {
			result := trimText(tc.Script)
			assert.Equal(t, tc.TrimedScript, result)
		}
	}

	t.Run("empty script", run(&testCase{
		Script:       "",
		TrimedScript: "",
	}))
	t.Run("no indentation on first line", run(&testCase{
		Script: `aaaaa
		bb bb bbbbb
            ccccccccc`,
		TrimedScript: `aaaaa
		bb bb bbbbb
            ccccccccc`,
	}))
	t.Run("trim start, end, left", run(&testCase{
		Script: `
	aaaaa
		b bbbbbbbb
			ccccccccc
	
	`,
		TrimedScript: `aaaaa
	b bbbbbbbb
		ccccccccc`,
	}))
	t.Run("trim mix of spaces and tabs", run(&testCase{
		Script: `
  	  	aaa  aa
  	  		bb bbbbbbb
  	  			ccccccc	cc
  	  	
  	  	`,
		TrimedScript: `aaa  aa
	bb bbbbbbb
		ccccccc	cc`,
	}))
	t.Run("trim bash script", run(&testCase{
		Script: `
			_scw()
			{
				output=$(scw autocomplete complete bash "$COMP_LINE" "$COMP_CWORD" "${COMP_WORDS[@]}")
				COMPREPLY=($output)
				[[ $COMPREPLY == *= ]] && compopt -o nospace
				return
			}
			complete -F _scw scw
		`,
		TrimedScript: `_scw()
{
	output=$(scw autocomplete complete bash "$COMP_LINE" "$COMP_CWORD" "${COMP_WORDS[@]}")
	COMPREPLY=($output)
	[[ $COMPREPLY == *= ]] && compopt -o nospace
	return
}
complete -F _scw scw`,
	}))
}
