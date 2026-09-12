package core

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testYAMLErrorMessage = "failure"
	testYAMLName         = "example"
)

func TestMarshalGoldenIncludesYAML(t *testing.T) {
	tests := map[string]struct {
		ctx      *CheckFuncCtx
		expected string
	}{
		"stdout": {
			ctx: &CheckFuncCtx{Result: struct {
				Name string `json:"name"`
			}{Name: testYAMLName}},
			expected: "🎲🎲🎲 EXIT CODE: 0 🎲🎲🎲\n" +
				"🟩🟩🟩 JSON STDOUT 🟩🟩🟩\n{\n  \"name\": \"" + testYAMLName + "\"\n}\n" +
				"🟩🟩🟩 YAML STDOUT 🟩🟩🟩\nname: " + testYAMLName + "\n",
		},
		"stderr": {
			ctx: &CheckFuncCtx{Err: errors.New(testYAMLErrorMessage)},
			expected: "🎲🎲🎲 EXIT CODE: 0 🎲🎲🎲\n" +
				"🟥🟥🟥 JSON STDERR 🟥🟥🟥\n{\n  \"error\": \"" + testYAMLErrorMessage + "\"\n}\n" +
				"🟥🟥🟥 YAML STDERR 🟥🟥🟥\nerror: " + testYAMLErrorMessage + "\n",
		},
		"empty CLI error": {
			ctx: &CheckFuncCtx{Err: &CliError{Empty: true}},
			expected: "🎲🎲🎲 EXIT CODE: 0 🎲🎲🎲\n" +
				"🟥🟥🟥 JSON STDERR 🟥🟥🟥\n{}\n" +
				"🟥🟥🟥 YAML STDERR 🟥🟥🟥\n{}\n",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.expected, marshalGolden(t, test.ctx))
		})
	}
}
