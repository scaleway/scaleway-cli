package autocomplete_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/autocomplete"
	"github.com/stretchr/testify/assert"
)

func TestTrimText(t *testing.T) {
	type testCase struct {
		Script        string
		TrimmedScript string
	}

	run := func(tc *testCase) func(*testing.T) {
		return func(t *testing.T) {
			t.Helper()
			result := autocomplete.TrimText(tc.Script)
			assert.Equal(t, tc.TrimmedScript, result)
		}
	}

	t.Run("empty script", run(&testCase{
		Script:        "",
		TrimmedScript: "",
	}))
	t.Run("no indentation on first line", run(&testCase{
		Script: `aaaaa
		bb bb bbbbb
            ccccccccc`,
		TrimmedScript: `aaaaa
		bb bb bbbbb
            ccccccccc`,
	}))
	t.Run("trim start, end, left", run(&testCase{
		Script: `
	aaaaa
		b bbbbbbbb
			ccccccccc
	
	`,
		TrimmedScript: `aaaaa
	b bbbbbbbb
		ccccccccc`,
	}))
	t.Run("trim mix of spaces and tabs", run(&testCase{
		Script: `
  	  	aaa  aa
  	  		bb bbbbbbb
  	  			ccccccc	cc
  	  	
  	  	`,
		TrimmedScript: `aaa  aa
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
		TrimmedScript: `_scw()
{
	output=$(scw autocomplete complete bash "$COMP_LINE" "$COMP_CWORD" "${COMP_WORDS[@]}")
	COMPREPLY=($output)
	[[ $COMPREPLY == *= ]] && compopt -o nospace
	return
}
complete -F _scw scw`,
	}))
}
