package core

import (
	"context"
	"fmt"
	"reflect"
	"testing"
)

type testType struct {
	Name string
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
