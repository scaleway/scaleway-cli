package core_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/args"
	"github.com/stretchr/testify/assert"
)

func Test_ApplyDefaultValues(t *testing.T) {
	type testCase struct {
		argSpecs core.ArgSpecs
		rawArgs  args.RawArgs
		expected args.RawArgs
	}

	run := func(tc *testCase) func(t *testing.T) {
		return func(t *testing.T) {
			t.Helper()
			result := core.ApplyDefaultValues(t.Context(), tc.argSpecs, tc.rawArgs)
			assert.Equal(t, tc.expected, result)
		}
	}

	t.Run("Simple", run(&testCase{
		argSpecs: core.ArgSpecs{
			{Name: "firstname", Default: core.DefaultValueSetter("john")},
			{Name: "lastname", Default: core.DefaultValueSetter("doe")},
			{Name: "age", Default: core.DefaultValueSetter("42")},
		},
		rawArgs:  args.RawArgs{"age=21"},
		expected: args.RawArgs{"age=21", "firstname=john", "lastname=doe"},
	}))

	t.Run("Nested", run(&testCase{
		argSpecs: core.ArgSpecs{
			{Name: "name", Default: core.DefaultValueSetter("john")},
			{Name: "address.zip", Default: core.DefaultValueSetter("75008")},
		},
		rawArgs:  args.RawArgs{"name=paul"},
		expected: args.RawArgs{"name=paul", "address.zip=75008"},
	}))

	t.Run("Slice", run(&testCase{
		argSpecs: core.ArgSpecs{
			{Name: "name", Default: core.DefaultValueSetter("john")},
			{Name: "friends.{index}.name", Default: core.DefaultValueSetter("bob")},
			{Name: "friends.{index}.age", Default: core.DefaultValueSetter("42")},
		},
		rawArgs: args.RawArgs{"name=paul", "friends.0.name=bob", "friends.1.name=alice"},
		expected: args.RawArgs{
			"name=paul",
			"friends.0.name=bob",
			"friends.1.name=alice",
			"friends.0.age=42",
			"friends.1.age=42",
		},
	}))

	t.Run("Map", run(&testCase{
		argSpecs: core.ArgSpecs{
			{Name: "name", Default: core.DefaultValueSetter("john")},
			{Name: "friends.{key}.age", Default: core.DefaultValueSetter("42")},
			{Name: "friends.{key}.weight", Default: core.DefaultValueSetter("80")},
		},
		rawArgs: args.RawArgs{"name=paul", "friends.bob.weight=75", "friends.alice.age=21"},
		expected: args.RawArgs{
			"name=paul",
			"friends.bob.weight=75",
			"friends.alice.age=21",
			"friends.bob.age=42",
			"friends.alice.weight=80",
		},
	}))

	t.Run("Map of slice", run(&testCase{
		argSpecs: core.ArgSpecs{
			{Name: "map.{key}.{index}.age", Default: core.DefaultValueSetter("42")},
			{Name: "map.{key}.{index}.weight", Default: core.DefaultValueSetter("80")},
		},
		rawArgs: args.RawArgs{
			"map.paul.0.weight=75",
			"map.paul.1.age=42",
			"map.alice.0.weight=60",
			"map.alice.1.age=25",
		},
		expected: args.RawArgs{
			"map.paul.0.weight=75",
			"map.paul.1.age=42",
			"map.alice.0.weight=60",
			"map.alice.1.age=25",
			"map.paul.0.age=42",
			"map.alice.0.age=42",
			"map.paul.1.weight=80",
			"map.alice.1.weight=80",
		},
	}))
}
