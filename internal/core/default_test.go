package core

import (
	"testing"

	"github.com/alecthomas/assert"
)

func Test_ApplyDefaultValues(t *testing.T) {

	type testCase struct {
		argSpecs ArgSpecs
		rawArgs  RawArgs
		expected RawArgs
	}

	run := func(tc *testCase) func(t *testing.T) {
		return func(t *testing.T) {
			result := ApplyDefaultValues(tc.argSpecs, tc.rawArgs)
			assert.Equal(t, tc.expected, result)
		}
	}

	t.Run("Simple", run(&testCase{
		argSpecs: ArgSpecs{
			{Name: "firstname", Default: DefaultValueSetter("john")},
			{Name: "lastname", Default: DefaultValueSetter("doe")},
			{Name: "age", Default: DefaultValueSetter("42")},
		},
		rawArgs:  RawArgs{"age=21"},
		expected: RawArgs{"age=21", "firstname=john", "lastname=doe"},
	}))

	t.Run("Nested", run(&testCase{
		argSpecs: ArgSpecs{
			{Name: "name", Default: DefaultValueSetter("john")},
			{Name: "address.zip", Default: DefaultValueSetter("75008")},
		},
		rawArgs:  RawArgs{"name=paul"},
		expected: RawArgs{"name=paul", "address.zip=75008"},
	}))

	t.Run("Slice", run(&testCase{
		argSpecs: ArgSpecs{
			{Name: "name", Default: DefaultValueSetter("john")},
			{Name: "friends.{index}.name", Default: DefaultValueSetter("bob")},
			{Name: "friends.{index}.age", Default: DefaultValueSetter("42")},
		},
		rawArgs:  RawArgs{"name=paul", "friends.0.name=bob", "friends.1.name=alice"},
		expected: RawArgs{"name=paul", "friends.0.name=bob", "friends.1.name=alice", "friends.0.age=42", "friends.1.age=42"},
	}))

	t.Run("Map", run(&testCase{
		argSpecs: ArgSpecs{
			{Name: "name", Default: DefaultValueSetter("john")},
			{Name: "friends.{key}.age", Default: DefaultValueSetter("42")},
			{Name: "friends.{key}.weight", Default: DefaultValueSetter("80")},
		},
		rawArgs:  RawArgs{"name=paul", "friends.bob.weight=75", "friends.alice.age=21"},
		expected: RawArgs{"name=paul", "friends.bob.weight=75", "friends.alice.age=21", "friends.bob.age=42", "friends.alice.weight=80"},
	}))

	t.Run("Map of slice", run(&testCase{
		argSpecs: ArgSpecs{
			{Name: "map.{key}.{index}.age", Default: DefaultValueSetter("42")},
			{Name: "map.{key}.{index}.weight", Default: DefaultValueSetter("80")},
		},
		rawArgs: RawArgs{
			"map.paul.0.weight=75",
			"map.paul.1.age=42",
			"map.alice.0.weight=60",
			"map.alice.1.age=25",
		},
		expected: RawArgs{
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
