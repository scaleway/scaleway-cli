package autocomplete_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/autocomplete"
	"github.com/stretchr/testify/assert"
)

func TestFishSuggestions(t *testing.T) {
	tests := map[string]struct {
		Response *core.AutocompleteResponse
		Expected string
	}{
		"no description": {
			Response: &core.AutocompleteResponse{
				Suggestions: core.AutocompleteSuggestions{"id1", "id2"},
			},
			Expected: "id1\nid2",
		},
		"with description": {
			Response: &core.AutocompleteResponse{
				Suggestions:  core.AutocompleteSuggestions{"id1", "id2"},
				Descriptions: map[string]string{"id1": "flower1"},
			},
			Expected: "id1\tflower1\nid2",
		},
		"arg value suggestion": {
			Response: &core.AutocompleteResponse{
				Suggestions:  core.AutocompleteSuggestions{"flower-id=id1"},
				Descriptions: map[string]string{"flower-id=id1": "flower1"},
			},
			Expected: "flower-id=id1\tflower1",
		},
		"sanitized description": {
			Response: &core.AutocompleteResponse{
				Suggestions:  core.AutocompleteSuggestions{"id1"},
				Descriptions: map[string]string{"id1": "flower\t1\n2"},
			},
			Expected: "id1\tflower 1 2",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.Expected, autocomplete.FishSuggestions(tc.Response))
		})
	}
}

func TestZshSuggestions(t *testing.T) {
	tests := map[string]struct {
		Response *core.AutocompleteResponse
		Expected string
	}{
		"no description": {
			Response: &core.AutocompleteResponse{
				Suggestions: core.AutocompleteSuggestions{"id1", "id2"},
			},
			Expected: "id1\nid2",
		},
		"with description": {
			Response: &core.AutocompleteResponse{
				Suggestions:  core.AutocompleteSuggestions{"id1", "id2"},
				Descriptions: map[string]string{"id1": "flower1"},
			},
			Expected: "id1:flower1\nid2",
		},
		"arg value suggestion": {
			Response: &core.AutocompleteResponse{
				Suggestions:  core.AutocompleteSuggestions{"flower-id=id1"},
				Descriptions: map[string]string{"flower-id=id1": "flower1"},
			},
			Expected: "flower-id=id1:flower1",
		},
		"sanitized description": {
			Response: &core.AutocompleteResponse{
				Suggestions:  core.AutocompleteSuggestions{"id1"},
				Descriptions: map[string]string{"id1": "flower\t1\n2"},
			},
			Expected: "id1:flower 1 2",
		},
		"escaped value": {
			Response: &core.AutocompleteResponse{
				Suggestions:  core.AutocompleteSuggestions{"id:1"},
				Descriptions: map[string]string{"id:1": "flower1"},
			},
			Expected: `id\:1:flower1`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.Expected, autocomplete.ZshSuggestions(tc.Response))
		})
	}
}

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
