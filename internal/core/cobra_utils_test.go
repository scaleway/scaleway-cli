package core

import (
	"context"
	"fmt"
	"reflect"
	"testing"
)

type testType struct {
	NameID string
	Tag    string
}

func testGetCommands() *Commands {
	return NewCommands(
		&Command{
			Namespace: "test",
			ArgSpecs: ArgSpecs{
				{
					Name: "name-id",
				},
			},
			ArgsType: reflect.TypeOf(testType{}),
			Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
				return "", nil
			},
		},
		&Command{
			Namespace: "test-positional",
			ArgSpecs: ArgSpecs{
				{
					Name:       "name-id",
					Positional: true,
				},
				{
					Name: "tag",
				},
			},
			ArgsType: reflect.TypeOf(testType{}),
			Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
				return "", nil
			},
		},
	)
}

func Test_handleUnmarshalErrors(t *testing.T) {
	t.Run("underscore", Test(&TestConfig{
		Commands: testGetCommands(),
		Cmd:      "scw test name_id",
		Check: TestCheckCombine(
			TestCheckExitCode(1),
			TestCheckError(&CliError{
				Err:  fmt.Errorf("unknown argument 'name_id'"),
				Hint: fmt.Sprintf("Valid arguments are: name-id"),
			}),
		),
	}))

	t.Run("value only", Test(&TestConfig{
		Commands: testGetCommands(),
		Cmd:      "scw test ubuntu-bionic",
		Check: TestCheckCombine(
			TestCheckExitCode(1),
			TestCheckError(&CliError{
				Err:  fmt.Errorf("unknown argument 'ubuntu-bionic'"),
				Hint: fmt.Sprintf("Valid arguments are: name-id"),
			}),
		),
	}))
}

func Test_PositionalArg(t *testing.T) {
	t.Run("Error", func(t *testing.T) {
		t.Run("Missing1", Test(&TestConfig{
			Commands: testGetCommands(),
			Cmd:      "scw test-positional",
			Check: TestCheckCombine(
				TestCheckExitCode(1),
				TestCheckError(&CliError{
					Err:  fmt.Errorf("a positional argument is required for this command"),
					Hint: "Try running 'scw test-positional <name-id>'.",
				}),
			),
		}))

		t.Run("Missing2", Test(&TestConfig{
			Commands: testGetCommands(),
			Cmd:      "scw test-positional tag=world",
			Check: TestCheckCombine(
				TestCheckExitCode(1),
				TestCheckError(&CliError{
					Err:  fmt.Errorf("a positional argument is required for this command"),
					Hint: "Try running 'scw test-positional <name-id> tag=world'.",
				}),
			),
		}))

		t.Run("Invalid1", Test(&TestConfig{
			Commands: testGetCommands(),
			Cmd:      "scw test-positional name-id=plop tag=world",
			Check: TestCheckCombine(
				TestCheckExitCode(1),
				TestCheckError(&CliError{
					Err:  fmt.Errorf("a positional argument is required for this command"),
					Hint: "Try running 'scw test-positional plop tag=world'.",
				}),
			),
		}))

		t.Run("Invalid2", Test(&TestConfig{
			Commands: testGetCommands(),
			Cmd:      "scw test-positional tag=world name-id=plop",
			Check: TestCheckCombine(
				TestCheckExitCode(1),
				TestCheckError(&CliError{
					Err:  fmt.Errorf("a positional argument is required for this command"),
					Hint: fmt.Sprintf("Try running 'scw test-positional plop tag=world'."),
				}),
			),
		}))

		t.Run("Invalid3", Test(&TestConfig{
			Commands: testGetCommands(),
			Cmd:      "scw test-positional plop name-id=plop",
			Check: TestCheckCombine(
				TestCheckExitCode(1),
				TestCheckError(&CliError{
					Err:  fmt.Errorf("a positional argument is required for this command"),
					Hint: fmt.Sprintf("Try running 'scw test-positional plop'."),
				}),
			),
		}))
	})

	t.Run("simple", Test(&TestConfig{
		Commands: testGetCommands(),
		Cmd:      "scw test-positional plop",
		Check:    TestCheckExitCode(0),
	}))

	t.Run("full command", Test(&TestConfig{
		Commands: testGetCommands(),
		Cmd:      "scw test-positional plop tag=world",
		Check:    TestCheckExitCode(0),
	}))
}
