package core

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/scaleway/scaleway-cli/internal/args"
)

type testType struct {
	NameID string
	Tag    string
}

type testDate struct {
	Date *time.Time
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
			AllowAnonymousClient: true,
			ArgsType:             reflect.TypeOf(testType{}),
			Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
				return "", nil
			},
		},
		&Command{
			Namespace: "test",
			Resource:  "positional",
			ArgSpecs: ArgSpecs{
				{
					Name:       "name-id",
					Positional: true,
				},
				{
					Name: "tag",
				},
			},
			AllowAnonymousClient: true,
			ArgsType:             reflect.TypeOf(testType{}),
			Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
				return argsI, nil
			},
		},
		&Command{
			Namespace:            "test",
			Resource:             "raw-args",
			ArgsType:             reflect.TypeOf(args.RawArgs{}),
			AllowAnonymousClient: true,
			Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
				res := ""
				rawArgs := *argsI.(*args.RawArgs)
				for i, arg := range rawArgs {
					res += arg
					if i != len(rawArgs)-1 {
						res += " "
					}
				}
				return res, nil
			},
		},
		&Command{
			Namespace:            "test",
			Resource:             "date",
			ArgsType:             reflect.TypeOf(testDate{}),
			AllowAnonymousClient: true,
			Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
				a := argsI.(*testDate)
				return a.Date, nil
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
				Err:  fmt.Errorf("invalid argument 'name_id': arg name must only contain lowercase letters, numbers or dashes"),
				Hint: "Valid arguments are: name-id",
			}),
		),
	}))

	t.Run("value only", Test(&TestConfig{
		Commands: testGetCommands(),
		Cmd:      "scw test ubuntu_focal",
		Check: TestCheckCombine(
			TestCheckExitCode(1),
			TestCheckError(&CliError{
				Err:  fmt.Errorf("invalid argument 'ubuntu_focal': arg name must only contain lowercase letters, numbers or dashes"),
				Hint: "Valid arguments are: name-id",
			}),
		),
	}))

	t.Run("relative date", Test(&TestConfig{
		Commands: testGetCommands(),
		Cmd:      "scw test date date=+3R",
		Check: TestCheckCombine(
			TestCheckExitCode(1),
			TestCheckError(&CliError{
				Message: "could not parse +3R as either an absolute time (RFC3339) nor a relative time (+/-)RFC3339",
				Details: `Absolute time error: parsing time "+3R" as "2006-01-02T15:04:05Z07:00": cannot parse "+3R" as "2006"
Relative time error: unknown unit in duration: "R"
`,
				Err:  fmt.Errorf("date parsing error: +3R"),
				Hint: "Run `scw help date` to learn more about date parsing",
			}),
		),
	}))
}

func Test_RawArgs(t *testing.T) {
	t.Run("Simple", Test(&TestConfig{
		Commands: testGetCommands(),
		Cmd:      "scw test raw-args -- blabla",
		Check: TestCheckCombine(
			TestCheckExitCode(0),
			TestCheckStdout("blabla\n"),
		),
	}))
	t.Run("Multiple", Test(&TestConfig{
		Commands: testGetCommands(),
		Cmd:      "scw test raw-args -- blabla foo bar",
		Check: TestCheckCombine(
			TestCheckExitCode(0),
			TestCheckStdout("blabla foo bar\n"),
		),
	}))
}

func Test_PositionalArg(t *testing.T) {
	t.Run("Error", func(t *testing.T) {
		t.Run("Missing1", Test(&TestConfig{
			Commands: testGetCommands(),
			Cmd:      "scw test positional",
			Check: TestCheckCombine(
				TestCheckExitCode(1),
				TestCheckError(&CliError{
					Err:  fmt.Errorf("a positional argument is required for this command"),
					Hint: "Try running: scw test positional <name-id>",
				}),
			),
		}))

		t.Run("Missing2", Test(&TestConfig{
			Commands: testGetCommands(),
			Cmd:      "scw test positional tag=world",
			Check: TestCheckCombine(
				TestCheckExitCode(1),
				TestCheckError(&CliError{
					Err:  fmt.Errorf("a positional argument is required for this command"),
					Hint: "Try running: scw test positional <name-id> tag=world",
				}),
			),
		}))

		t.Run("Invalid1", Test(&TestConfig{
			Commands: testGetCommands(),
			Cmd:      "scw test positional name-id=plop tag=world",
			Check: TestCheckCombine(
				TestCheckExitCode(1),
				TestCheckError(&CliError{
					Err:  fmt.Errorf("a positional argument is required for this command"),
					Hint: "Try running: scw test positional plop tag=world",
				}),
			),
		}))

		t.Run("Invalid2", Test(&TestConfig{
			Commands: testGetCommands(),
			Cmd:      "scw test positional tag=world name-id=plop",
			Check: TestCheckCombine(
				TestCheckExitCode(1),
				TestCheckError(&CliError{
					Err:  fmt.Errorf("a positional argument is required for this command"),
					Hint: "Try running: scw test positional plop tag=world",
				}),
			),
		}))

		t.Run("Invalid3", Test(&TestConfig{
			Commands: testGetCommands(),
			Cmd:      "scw test positional plop name-id=plop",
			Check: TestCheckCombine(
				TestCheckExitCode(1),
				TestCheckError(&CliError{
					Err:  fmt.Errorf("a positional argument is required for this command"),
					Hint: "Try running: scw test positional plop",
				}),
			),
		}))
	})

	t.Run("simple", Test(&TestConfig{
		Commands: testGetCommands(),
		Cmd:      "scw test positional plop",
		Check:    TestCheckExitCode(0),
	}))

	t.Run("full command", Test(&TestConfig{
		Commands: testGetCommands(),
		Cmd:      "scw test positional plop tag=world",
		Check:    TestCheckExitCode(0),
	}))

	t.Run("full command", Test(&TestConfig{
		Commands: testGetCommands(),
		Cmd:      "scw test positional -h",
		Check: TestCheckCombine(
			TestCheckExitCode(0),
			TestCheckGolden(),
		),
	}))

	t.Run("multi positional", Test(&TestConfig{
		Commands: testGetCommands(),
		Cmd:      "scw test positional tag=tag01 test1 test2",
		Check: TestCheckCombine(
			TestCheckExitCode(0),
			TestCheckGolden(),
		),
	}))

	t.Run("multi positional json", Test(&TestConfig{
		Commands: testGetCommands(),
		Cmd:      "scw test positional -o json tag=tag01 test1 test2",
		Check: TestCheckCombine(
			TestCheckExitCode(0),
			TestCheckGolden(),
		),
	}))
}
