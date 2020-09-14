package core

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/scaleway/scaleway-cli/internal/args"
)

type Element struct {
	ID                 int
	Name               string
	ElementsMap        map[string]Element
	ElementsSlice      []Element
	FirstNestedElement *FirstNestedElement
}

type FirstNestedElement struct {
	SecondNestedElement *SecondNestedElement
}

type SecondNestedElement struct {
}

type elementCustom struct {
	*Element
	Short string
}

func Test_DefaultCommandValidateFunc(t *testing.T) {
	type TestCase struct {
		command         *Command
		parsedArguments interface{}
		rawArgs         args.RawArgs
	}

	run := func(testCase TestCase) func(t *testing.T) {
		return func(t *testing.T) {
			err := DefaultCommandValidateFunc()(context.Background(), testCase.command, testCase.parsedArguments, testCase.rawArgs)
			assert.Equal(t, fmt.Errorf("arg validation called"), err)
		}
	}

	t.Run("simple", run(TestCase{
		command: &Command{
			ArgSpecs: ArgSpecs{
				{
					Name: "name",
					ValidateFunc: func(argSpec *ArgSpec, value interface{}) error {
						return fmt.Errorf("arg validation called")
					},
				},
			},
		},
		parsedArguments: &Element{
			Name: "bob",
		},
	}))

	t.Run("map", run(TestCase{
		command: &Command{
			ArgSpecs: ArgSpecs{
				{
					Name: "elements-map.{key}.id",
				},
				{
					Name: "elements-map.{key}.name",
					ValidateFunc: func(argSpec *ArgSpec, value interface{}) error {
						return fmt.Errorf("arg validation called")
					},
				},
			},
		},
		parsedArguments: &Element{
			ElementsMap: map[string]Element{
				"first": {
					ID:   1,
					Name: "first",
				},
				"second": {
					ID:   2,
					Name: "second",
				},
			},
		},
	}))

	t.Run("slice", run(TestCase{
		command: &Command{
			ArgSpecs: ArgSpecs{
				{
					Name: "elements-slice.{index}.id",
				},
				{
					Name: "elements-slice.{index}.name",
					ValidateFunc: func(argSpec *ArgSpec, value interface{}) error {
						return fmt.Errorf("arg validation called")
					},
				},
			},
		},
		parsedArguments: &Element{
			ElementsSlice: []Element{
				{
					ID:   1,
					Name: "first",
				},
				{
					ID:   2,
					Name: "second",
				},
			},
		},
	}))

	t.Run("slice-of-slice", run(TestCase{
		command: &Command{
			ArgSpecs: ArgSpecs{
				{
					Name: "elements-slice.{index}.id",
				},
				{
					Name: "elements-slice.{index}.elements-slice.{index}.name",
					ValidateFunc: func(argSpec *ArgSpec, value interface{}) error {
						return fmt.Errorf("arg validation called")
					},
				},
			},
		},
		parsedArguments: &Element{
			ElementsSlice: []Element{
				{
					ID: 1,
				},
				{
					ElementsSlice: []Element{
						{
							Name: "bob",
						},
					},
				},
			},
		},
	}))

	t.Run("new-field", run(TestCase{
		command: &Command{
			ArgSpecs: ArgSpecs{
				{
					Name: "name",
				},
				{
					Name: "short",
					ValidateFunc: func(argSpec *ArgSpec, value interface{}) error {
						return fmt.Errorf("arg validation called")
					},
				},
			},
		},
		parsedArguments: &elementCustom{
			Short: "bob",
		},
	}))

	t.Run("anonymous-field", run(TestCase{
		command: &Command{
			ArgSpecs: ArgSpecs{
				{
					Name: "short",
				},
				{
					Name: "name",
					ValidateFunc: func(argSpec *ArgSpec, value interface{}) error {
						return fmt.Errorf("arg validation called")
					},
				},
			},
		},
		parsedArguments: &elementCustom{
			Element: &Element{
				Name: "bob",
			},
		},
	}))
}

func Test_DefaultCommandRequiredFunc(t *testing.T) {
	type TestCase struct {
		command         *Command
		parsedArguments interface{}
		rawArgs         args.RawArgs
	}

	runOK := func(testCase TestCase) func(t *testing.T) {
		return func(t *testing.T) {
			err := DefaultCommandValidateFunc()(context.Background(), testCase.command, testCase.parsedArguments, testCase.rawArgs)
			assert.Equal(t, nil, err)
		}
	}

	runErr := func(testCase TestCase, argName string) func(t *testing.T) {
		return func(t *testing.T) {
			err := DefaultCommandValidateFunc()(context.Background(), testCase.command, testCase.parsedArguments, testCase.rawArgs)
			assert.Equal(t, MissingRequiredArgumentError(argName), err)
		}
	}

	t.Run("required-struct", runOK(TestCase{
		command: &Command{
			ArgSpecs: ArgSpecs{
				{
					Name:     "first-nested-element.second-nested-element",
					Required: true,
				},
			},
		},
		rawArgs: []string{"first-nested-element.second-nested-element=test"},
		parsedArguments: &elementCustom{
			Element: &Element{
				Name: "nested",
				FirstNestedElement: &FirstNestedElement{
					SecondNestedElement: &SecondNestedElement{},
				},
			},
		},
	}))

	t.Run("fail-required-struct", runErr(TestCase{
		command: &Command{
			ArgSpecs: ArgSpecs{
				{
					Name:     "first-nested-element.second-nested-element",
					Required: true,
				},
			},
		},
		parsedArguments: &elementCustom{
			Element: &Element{
				Name:               "foo",
				FirstNestedElement: &FirstNestedElement{},
			},
		},
	}, "first-nested-element.second-nested-element"))

	t.Run("required-index", runOK(TestCase{
		command: &Command{
			ArgSpecs: ArgSpecs{
				{
					Name:     "elements-slice.{index}.id",
					Required: true,
				},
			},
		},
		rawArgs: []string{"elements-slice.0.id=1"},
		parsedArguments: &Element{
			ElementsSlice: []Element{
				{
					ID:   0,
					Name: "1",
				},
			},
		},
	}))

	t.Run("fail-required-index", runErr(TestCase{
		command: &Command{
			ArgSpecs: ArgSpecs{
				{
					Name:     "elements-slice.{index}.id",
					Required: true,
				},
			},
		},
		rawArgs: []string{"elements-slice.0.id=1"},
		parsedArguments: &Element{
			ElementsSlice: []Element{
				{
					ID:   0,
					Name: "1",
				},
				{
					ID:   1,
					Name: "0",
				},
			},
		},
	}, "elements-slice.1.id"))
}

func Test_ValidateNoConflict(t *testing.T) {
	type TestCase struct {
		command *Command
		rawArgs args.RawArgs
		arg1    string
		arg2    string
	}

	runOK := func(testCase TestCase) func(t *testing.T) {
		return func(t *testing.T) {
			err := validateNoConflict(testCase.command, testCase.rawArgs)
			assert.Equal(t, nil, err)
		}
	}

	runErr := func(testCase TestCase) func(t *testing.T) {
		return func(t *testing.T) {
			err := validateNoConflict(testCase.command, testCase.rawArgs)
			assert.Equal(t, ArgumentConflictError(testCase.arg1, testCase.arg2), err)
		}
	}

	t.Run("No conflict", runOK(TestCase{
		command: &Command{
			ArgSpecs: ArgSpecs{
				{
					Name:       "a",
					OneOfGroup: "a",
				},
				{
					Name: "b",
				},
			},
		},
		rawArgs: []string{"a=foo", "b=bar"},
	}))

	t.Run("SSH example", runErr(TestCase{
		command: &Command{
			ArgSpecs: ArgSpecs{
				{
					Name:       "ssh-key.{index}",
					OneOfGroup: "ssh",
				},
				{
					Name:       "all-ssh-keys",
					OneOfGroup: "ssh",
				},
			},
		},
		rawArgs: []string{"all-ssh-keys=true", "ssh-key.0=11111111-1111-1111-1111-111111111111"},
		arg1:    "ssh-key.{index}",
		arg2:    "all-ssh-keys",
	}))
}

func Test_ValidateDeprecated(t *testing.T) {
	t.Run("Deprecated", Test(&TestConfig{
		Commands: NewCommands(&Command{
			Namespace:            "plop",
			ArgsType:             reflect.TypeOf(args.RawArgs{}),
			AllowAnonymousClient: true,
			Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
				return &SuccessResult{}, nil
			},
			ArgSpecs: ArgSpecs{
				{
					Name:       "a",
					Deprecated: true,
				},
			},
		}),
		Cmd: "scw plop a=yo",
		Check: TestCheckCombine(
			func(t *testing.T, ctx *CheckFuncCtx) {
				assert.Equal(t, "The argument 'a' is deprecated, more info with: scw plop --help\n", ctx.LogBuffer)
			},
		),
	}))
}
