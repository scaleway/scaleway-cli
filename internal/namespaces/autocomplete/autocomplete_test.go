package autocomplete

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
