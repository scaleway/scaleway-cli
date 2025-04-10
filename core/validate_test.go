package core_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/args"
	"github.com/stretchr/testify/assert"
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

type SecondNestedElement struct{}

type elementCustom struct {
	*Element
	Short string
}

func Test_DefaultCommandValidateFunc(t *testing.T) {
	type TestCase struct {
		command         *core.Command
		parsedArguments interface{}
		rawArgs         args.RawArgs
	}

	run := func(testCase TestCase) func(t *testing.T) {
		return func(t *testing.T) {
			t.Helper()
			err := core.DefaultCommandValidateFunc()(
				t.Context(),
				testCase.command,
				testCase.parsedArguments,
				testCase.rawArgs,
			)
			assert.Equal(t, errors.New("arg validation called"), err)
		}
	}

	t.Run("simple", run(TestCase{
		command: &core.Command{
			ArgSpecs: core.ArgSpecs{
				{
					Name: "name",
					ValidateFunc: func(_ *core.ArgSpec, _ interface{}) error {
						return errors.New("arg validation called")
					},
				},
			},
		},
		parsedArguments: &Element{
			Name: "bob",
		},
	}))

	t.Run("map", run(TestCase{
		command: &core.Command{
			ArgSpecs: core.ArgSpecs{
				{
					Name: "elements-map.{key}.id",
				},
				{
					Name: "elements-map.{key}.name",
					ValidateFunc: func(_ *core.ArgSpec, _ interface{}) error {
						return errors.New("arg validation called")
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
		command: &core.Command{
			ArgSpecs: core.ArgSpecs{
				{
					Name: "elements-slice.{index}.id",
				},
				{
					Name: "elements-slice.{index}.name",
					ValidateFunc: func(_ *core.ArgSpec, _ interface{}) error {
						return errors.New("arg validation called")
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
		command: &core.Command{
			ArgSpecs: core.ArgSpecs{
				{
					Name: "elements-slice.{index}.id",
				},
				{
					Name: "elements-slice.{index}.elements-slice.{index}.name",
					ValidateFunc: func(_ *core.ArgSpec, _ interface{}) error {
						return errors.New("arg validation called")
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
		command: &core.Command{
			ArgSpecs: core.ArgSpecs{
				{
					Name: "name",
				},
				{
					Name: "short",
					ValidateFunc: func(_ *core.ArgSpec, _ interface{}) error {
						return errors.New("arg validation called")
					},
				},
			},
		},
		parsedArguments: &elementCustom{
			Short: "bob",
		},
	}))

	t.Run("anonymous-field", run(TestCase{
		command: &core.Command{
			ArgSpecs: core.ArgSpecs{
				{
					Name: "short",
				},
				{
					Name: "name",
					ValidateFunc: func(_ *core.ArgSpec, _ interface{}) error {
						return errors.New("arg validation called")
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
		command         *core.Command
		parsedArguments interface{}
		rawArgs         args.RawArgs
	}

	runOK := func(testCase TestCase) func(t *testing.T) {
		return func(t *testing.T) {
			t.Helper()
			err := core.DefaultCommandValidateFunc()(
				t.Context(),
				testCase.command,
				testCase.parsedArguments,
				testCase.rawArgs,
			)
			assert.NoError(t, err)
		}
	}

	runErr := func(testCase TestCase, argName string) func(t *testing.T) {
		return func(t *testing.T) {
			t.Helper()
			err := core.DefaultCommandValidateFunc()(
				t.Context(),
				testCase.command,
				testCase.parsedArguments,
				testCase.rawArgs,
			)
			assert.Equal(t, core.MissingRequiredArgumentError(argName), err)
		}
	}

	t.Run("required-struct", runOK(TestCase{
		command: &core.Command{
			ArgSpecs: core.ArgSpecs{
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
		command: &core.Command{
			ArgSpecs: core.ArgSpecs{
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
		command: &core.Command{
			ArgSpecs: core.ArgSpecs{
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
		command: &core.Command{
			ArgSpecs: core.ArgSpecs{
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
		command *core.Command
		rawArgs args.RawArgs
		arg1    string
		arg2    string
	}

	runOK := func(testCase TestCase) func(t *testing.T) {
		return func(t *testing.T) {
			t.Helper()
			err := core.ValidateNoConflict(testCase.command, testCase.rawArgs)
			assert.NoError(t, err)
		}
	}

	runErr := func(testCase TestCase) func(t *testing.T) {
		return func(t *testing.T) {
			t.Helper()
			err := core.ValidateNoConflict(testCase.command, testCase.rawArgs)
			assert.Equal(t, core.ArgumentConflictError(testCase.arg1, testCase.arg2), err)
		}
	}

	t.Run("No conflict", runOK(TestCase{
		command: &core.Command{
			ArgSpecs: core.ArgSpecs{
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
		command: &core.Command{
			ArgSpecs: core.ArgSpecs{
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
	t.Run("Deprecated", core.Test(&core.TestConfig{
		Commands: core.NewCommands(&core.Command{
			Namespace:            "plop",
			ArgsType:             reflect.TypeOf(args.RawArgs{}),
			AllowAnonymousClient: true,
			Run: func(_ context.Context, _ interface{}) (i interface{}, e error) {
				return &core.SuccessResult{}, nil
			},
			ArgSpecs: core.ArgSpecs{
				{
					Name:       "a",
					Deprecated: true,
				},
			},
		}),
		Cmd: "scw plop a=yo",
		Check: core.TestCheckCombine(
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				assert.Equal(
					t,
					"The argument 'a' is deprecated, more info with: scw plop --help\n",
					ctx.LogBuffer,
				)
			},
		),
	}))
}

func TestNewOneOfGroupManager(t *testing.T) {
	type TestCase struct {
		command                *core.Command
		expectedGroups         map[string][]string
		expectedRequiredGroups map[string]bool
	}

	tests := []struct {
		name     string
		testCase TestCase
		testFunc func(*testing.T, TestCase)
	}{
		{
			name: "Basic OneOf Groups",
			testCase: TestCase{
				command: &core.Command{
					ArgSpecs: core.ArgSpecs{
						{Name: "a", OneOfGroup: "group1"},
						{Name: "b", OneOfGroup: "group1"},
					},
				},
				expectedGroups:         map[string][]string{"group1": {"a", "b"}},
				expectedRequiredGroups: map[string]bool{},
			},
			testFunc: func(t *testing.T, tc TestCase) {
				t.Helper()
				manager := core.NewOneOfGroupManager(tc.command)
				assert.Equal(t, tc.expectedGroups, manager.Groups)
				assert.Equal(t, tc.expectedRequiredGroups, manager.RequiredGroups)
			},
		},
		{
			name: "With Required Group",
			testCase: TestCase{
				command: &core.Command{
					ArgSpecs: core.ArgSpecs{
						{Name: "a", OneOfGroup: "group1", Required: true},
						{Name: "b", OneOfGroup: "group1", Required: true},
					},
				},
				expectedGroups:         map[string][]string{"group1": {"a", "b"}},
				expectedRequiredGroups: map[string]bool{"group1": true},
			},
			testFunc: func(t *testing.T, tc TestCase) {
				t.Helper()
				manager := core.NewOneOfGroupManager(tc.command)
				assert.Equal(t, tc.expectedGroups, manager.Groups)
				assert.Equal(t, tc.expectedRequiredGroups, manager.RequiredGroups)
			},
		},
		{
			name: "With two Group no required",
			testCase: TestCase{
				command: &core.Command{
					ArgSpecs: core.ArgSpecs{
						{Name: "a", OneOfGroup: "group1"},
						{Name: "b", OneOfGroup: "group1"},
						{Name: "c", OneOfGroup: "group2"},
						{Name: "d", OneOfGroup: "group2"},
					},
				},
				expectedGroups: map[string][]string{
					"group1": {"a", "b"},
					"group2": {"c", "d"},
				},
				expectedRequiredGroups: map[string]bool{},
			},
			testFunc: func(t *testing.T, tc TestCase) {
				t.Helper()
				manager := core.NewOneOfGroupManager(tc.command)
				assert.Equal(t, tc.expectedGroups, manager.Groups)
				assert.Equal(t, tc.expectedRequiredGroups, manager.RequiredGroups)
			},
		},
		{
			name: "With two Group with one required",
			testCase: TestCase{
				command: &core.Command{
					ArgSpecs: core.ArgSpecs{
						{Name: "a", OneOfGroup: "group1", Required: true},
						{Name: "b", OneOfGroup: "group1", Required: true},
						{Name: "c", OneOfGroup: "group2"},
						{Name: "d", OneOfGroup: "group2"},
					},
				},
				expectedGroups: map[string][]string{
					"group1": {"a", "b"},
					"group2": {"c", "d"},
				},
				expectedRequiredGroups: map[string]bool{
					"group1": true,
				},
			},
			testFunc: func(t *testing.T, tc TestCase) {
				t.Helper()
				manager := core.NewOneOfGroupManager(tc.command)
				assert.Equal(t, tc.expectedGroups, manager.Groups)
				assert.Equal(t, tc.expectedRequiredGroups, manager.RequiredGroups)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t, tt.testCase)
		})
	}
}

func TestValidateRequiredOneOfGroups(t *testing.T) {
	tests := []struct {
		name          string
		setupManager  func() *core.OneOfGroupManager
		rawArgs       args.RawArgs
		expectedError string
		ArgsType      interface{}
	}{
		{
			name: "Required group satisfied with first argument",
			setupManager: func() *core.OneOfGroupManager {
				return &core.OneOfGroupManager{
					Groups:         map[string][]string{"group1": {"a", "b"}},
					RequiredGroups: map[string]bool{"group1": true},
				}
			},
			rawArgs: []string{"a=true"},
			ArgsType: struct {
				A bool
				B bool
			}{},
			expectedError: "",
		},
		{
			name: "Required group satisfied with second argument",
			setupManager: func() *core.OneOfGroupManager {
				return &core.OneOfGroupManager{
					Groups:         map[string][]string{"group1": {"a", "b"}},
					RequiredGroups: map[string]bool{"group1": true},
				}
			},
			rawArgs: []string{"b=true"},
			ArgsType: struct {
				A bool
				B bool
			}{},
			expectedError: "",
		},
		{
			name: "Required group not satisfied",
			setupManager: func() *core.OneOfGroupManager {
				return &core.OneOfGroupManager{
					Groups:         map[string][]string{"group1": {"a", "b"}},
					RequiredGroups: map[string]bool{"group1": true},
				}
			},
			rawArgs: []string{"c=true"},
			ArgsType: struct {
				A bool
				B bool
				C bool
			}{},
			expectedError: "at least one argument from the 'group1' group is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := tt.setupManager()
			err := manager.ValidateRequiredOneOfGroups(tt.rawArgs, tt.ArgsType)

			if tt.expectedError == "" {
				assert.NoError(t, err, "Expected no error, got %v", err)
			} else {
				assert.EqualError(t, err, tt.expectedError, "Expected error message '%s', got '%v'", tt.expectedError, err)
			}
		})
	}
}

func TestValidateUniqueOneOfGroups(t *testing.T) {
	tests := []struct {
		name          string
		setupManager  func() *core.OneOfGroupManager
		rawArgs       args.RawArgs
		expectedError string
		ArgsType      interface{}
	}{
		{
			name: "Required group satisfied with first argument",
			setupManager: func() *core.OneOfGroupManager {
				return &core.OneOfGroupManager{
					Groups: map[string][]string{"group1": {"a", "b"}},
				}
			},
			rawArgs: []string{"A=true"},
			ArgsType: struct {
				A bool
				B bool
			}{},
			expectedError: "",
		},
		{
			name: "No arguments passed",
			setupManager: func() *core.OneOfGroupManager {
				return &core.OneOfGroupManager{
					Groups: map[string][]string{"group1": {"a", "b"}},
				}
			},
			rawArgs: []string{},
			ArgsType: struct {
				A bool
				B bool
			}{},
			expectedError: "",
		},
		{
			name: "Multiple groups, all satisfied",
			setupManager: func() *core.OneOfGroupManager {
				return &core.OneOfGroupManager{
					Groups: map[string][]string{
						"group1": {"a", "b"},
						"group2": {"c", "d"},
					},
				}
			},
			rawArgs: []string{"a=true", "c=true"},
			ArgsType: struct {
				A string
				B string
				C string
				D string
			}{},
			expectedError: "",
		},
		{
			name: "Multiple groups, one satisfied",
			setupManager: func() *core.OneOfGroupManager {
				return &core.OneOfGroupManager{
					Groups: map[string][]string{
						"group1": {"a", "b"},
						"group2": {"c", "d"},
					},
				}
			},
			rawArgs: []string{"a=true"},
			ArgsType: struct {
				A string
				B string
				C string
				D string
			}{},
			expectedError: "",
		},
		{
			name: "Multiple groups, not exclusive argument for groups 2",
			setupManager: func() *core.OneOfGroupManager {
				return &core.OneOfGroupManager{
					Groups: map[string][]string{
						"group1": {"a", "b"},
						"group2": {"c", "d"},
					},
				}
			},
			rawArgs: []string{"a=true", "c=true", "d=true"},
			ArgsType: struct {
				A string
				B string
				C string
				D string
			}{},
			expectedError: "arguments 'c' and 'd' are mutually exclusive",
		},
		{
			name: "Multiple groups, not exclusive argument for groups 1",
			setupManager: func() *core.OneOfGroupManager {
				return &core.OneOfGroupManager{
					Groups: map[string][]string{
						"group1": {"a", "b"},
						"group2": {"c", "d"},
					},
				}
			},
			rawArgs: []string{"a=true", "b=true", "c=true"},
			ArgsType: struct {
				A string
				B string
				C string
				D string
			}{},
			expectedError: "arguments 'a' and 'b' are mutually exclusive",
		},
		{
			name: "One group, not exclusive argument for groups 1",
			setupManager: func() *core.OneOfGroupManager {
				return &core.OneOfGroupManager{
					Groups: map[string][]string{
						"group1": {"a", "b", "c", "d"},
					},
				}
			},
			rawArgs: []string{"a=true", "d=true"},
			ArgsType: struct {
				A string
				B string
				C string
				D string
			}{},
			expectedError: "arguments 'a' and 'd' are mutually exclusive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := tt.setupManager()
			err := manager.ValidateUniqueOneOfGroups(tt.rawArgs, tt.ArgsType)
			if tt.expectedError == "" {
				assert.NoError(t, err, "Expected no error, got %v", err)
			} else {
				assert.EqualError(t, err, tt.expectedError, "Expected error message '%s', got '%v'", tt.expectedError, err)
			}
		})
	}
}

func Test_ValidateOneOf(t *testing.T) {
	t.Run("Simple one-of validation check", core.Test(&core.TestConfig{
		Commands: core.NewCommands(&core.Command{
			Namespace:            "oneof",
			ArgsType:             reflect.TypeOf(args.RawArgs{}),
			AllowAnonymousClient: true,
			Run: func(_ context.Context, _ interface{}) (i interface{}, e error) {
				return &core.SuccessResult{}, nil
			},
			ArgSpecs: core.ArgSpecs{
				{
					Name:       "a",
					OneOfGroup: "groups1",
				},
				{
					Name:       "b",
					OneOfGroup: "groups1",
				},
			},
		}),
		Cmd: "scw oneof a=yo",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
		),
	}))
	t.Run("Required argument group check passes", func(t *testing.T) {
		core.Test(&core.TestConfig{
			Commands: core.NewCommands(&core.Command{
				Namespace:            "oneof",
				ArgsType:             reflect.TypeOf(args.RawArgs{}),
				AllowAnonymousClient: true,
				Run: func(_ context.Context, _ interface{}) (i interface{}, e error) {
					return &core.SuccessResult{}, nil
				},
				ArgSpecs: core.ArgSpecs{
					{
						Name:       "a",
						OneOfGroup: "group1",
						Required:   true,
					},
					{
						Name:       "b",
						OneOfGroup: "group1",
						Required:   true,
					},
				},
			}),
			Cmd: "scw oneof b=yo",
			Check: core.TestCheckCombine(
				core.TestCheckExitCode(0),
			),
		})(t)
	})

	t.Run("Fail when required group is missing", func(t *testing.T) {
		core.Test(&core.TestConfig{
			Commands: core.NewCommands(&core.Command{
				Namespace:            "oneof",
				ArgsType:             reflect.TypeOf(args.RawArgs{}),
				AllowAnonymousClient: true,
				Run: func(_ context.Context, _ interface{}) (i interface{}, e error) {
					return &core.SuccessResult{}, nil
				},
				ArgSpecs: core.ArgSpecs{
					{
						Name:       "a",
						OneOfGroup: "group1",
						Required:   true,
					},
					{
						Name:       "b",
						OneOfGroup: "group1",
						Required:   true,
					},
				},
			}),
			Cmd: "scw oneof c=yo",
			Check: core.TestCheckCombine(
				core.TestCheckExitCode(1),
				core.TestCheckError(
					errors.New("at least one argument from the 'group1' group is required"),
				),
			),
		})(t)
	})

	t.Run("Check for mutual exclusivity in arguments", func(t *testing.T) {
		core.Test(&core.TestConfig{
			Commands: core.NewCommands(&core.Command{
				Namespace: "oneof",
				ArgsType: reflect.TypeOf(struct {
					A string
					B string
				}{}), AllowAnonymousClient: true,
				Run: func(_ context.Context, _ interface{}) (i interface{}, e error) {
					return &core.SuccessResult{}, nil
				},
				ArgSpecs: core.ArgSpecs{
					{
						Name:       "a",
						OneOfGroup: "group1",
					},
					{
						Name:       "b",
						OneOfGroup: "group1",
					},
				},
			}),
			Cmd: "scw oneof a=yo b=no",
			Check: core.TestCheckCombine(
				core.TestCheckExitCode(1),
				core.TestCheckError(errors.New("arguments 'a' and 'b' are mutually exclusive")),
			),
		})(t)
	})

	t.Run("Three arguments' mutual exclusivity test", func(t *testing.T) {
		core.Test(&core.TestConfig{
			Commands: core.NewCommands(&core.Command{
				Namespace: "oneof",
				ArgsType: reflect.TypeOf(struct {
					A string
					B string
					C string
				}{}), AllowAnonymousClient: true,
				Run: func(_ context.Context, _ interface{}) (i interface{}, e error) {
					return &core.SuccessResult{}, nil
				},
				ArgSpecs: core.ArgSpecs{
					{
						Name:       "a",
						OneOfGroup: "group1",
					},
					{
						Name:       "b",
						OneOfGroup: "group1",
					},
					{
						Name:       "c",
						OneOfGroup: "group1",
					},
				},
			}),
			Cmd: "scw oneof a=yo c=no",
			Check: core.TestCheckCombine(
				core.TestCheckExitCode(1),
				core.TestCheckError(errors.New("arguments 'a' and 'c' are mutually exclusive")),
			),
		})(t)
	})

	t.Run("Indexed arguments' exclusivity check", func(t *testing.T) {
		core.Test(&core.TestConfig{
			Commands: core.NewCommands(&core.Command{
				Namespace: "oneof",
				ArgsType: reflect.TypeOf(struct {
					SSHKey     []string
					AllSSHKeys bool
				}{}),
				AllowAnonymousClient: true,
				Run: func(_ context.Context, _ interface{}) (i interface{}, e error) {
					return &core.SuccessResult{}, nil
				},
				ArgSpecs: core.ArgSpecs{
					{
						Name:       "ssh-key.{index}",
						OneOfGroup: "ssh",
						Required:   true,
					},
					{
						Name:       "all-ssh-keys",
						OneOfGroup: "ssh",
						Required:   true,
					},
				},
			}),
			Cmd: "scw oneof all-ssh-keys=true ssh-key.0=11111111-1111-1111-1111-111111111111",
			Check: core.TestCheckCombine(
				core.TestCheckExitCode(1),
				core.TestCheckError(
					errors.New("arguments 'ssh-key.0' and 'all-ssh-keys' are mutually exclusive"),
				),
			),
		})(t)
	})

	t.Run("Passing an indexed argument test", func(t *testing.T) {
		core.Test(&core.TestConfig{
			Commands: core.NewCommands(&core.Command{
				Namespace: "oneof",
				ArgsType: reflect.TypeOf(struct {
					SSHKey     []string
					AllSSHKeys bool
				}{}),
				AllowAnonymousClient: true,
				Run: func(_ context.Context, _ interface{}) (i interface{}, e error) {
					return &core.SuccessResult{}, nil
				},
				ArgSpecs: core.ArgSpecs{
					{
						Name:       "ssh-key.{index}",
						OneOfGroup: "ssh",
					},
					{
						Name:       "all-ssh-keys",
						OneOfGroup: "ssh",
					},
				},
			}),
			Cmd: "scw oneof ssh-key.0=11111111-1111-1111-1111-111111111111",
			Check: core.TestCheckCombine(
				core.TestCheckExitCode(0),
			),
		})(t)
	})

	t.Run("Required indexed argument satisfies condition", func(t *testing.T) {
		core.Test(&core.TestConfig{
			Commands: core.NewCommands(&core.Command{
				Namespace: "oneof",
				ArgsType: reflect.TypeOf(struct {
					SSHKey     []string
					AllSSHKeys bool
				}{}),
				AllowAnonymousClient: true,
				Run: func(_ context.Context, _ interface{}) (i interface{}, e error) {
					return &core.SuccessResult{}, nil
				},
				ArgSpecs: core.ArgSpecs{
					{
						Name:       "ssh-key.{index}",
						OneOfGroup: "ssh",
						Required:   true,
					},
					{
						Name:       "all-ssh-keys",
						OneOfGroup: "ssh",
						Required:   true,
					},
				},
			}),
			Cmd: "scw oneof ssh-key.0=11111111-1111-1111-1111-111111111111",
			Check: core.TestCheckCombine(
				core.TestCheckExitCode(0),
			),
		})(t)
	})

	t.Run("Exclusive all SSH keys and indexed key test", func(t *testing.T) {
		core.Test(&core.TestConfig{
			Commands: core.NewCommands(&core.Command{
				Namespace: "oneof",
				ArgsType: reflect.TypeOf(struct {
					SSHKey     []string
					AllSSHKeys bool
				}{}),
				AllowAnonymousClient: true,
				Run: func(_ context.Context, _ interface{}) (i interface{}, e error) {
					return &core.SuccessResult{}, nil
				},
				ArgSpecs: core.ArgSpecs{
					{
						Name:       "ssh-key.{index}",
						OneOfGroup: "ssh",
						Required:   true,
					},
					{
						Name:       "all-ssh-keys",
						OneOfGroup: "ssh",
						Required:   true,
					},
				},
			}),
			Cmd: "scw oneof all-ssh-keys=true",
			Check: core.TestCheckCombine(
				core.TestCheckExitCode(0),
			),
		})(t)
	})

	t.Run("Ungrouped argument with unsatisfied group fails", func(t *testing.T) {
		core.Test(&core.TestConfig{
			Commands: core.NewCommands(&core.Command{
				Namespace: "oneof",
				ArgsType: reflect.TypeOf(struct {
					SSHKey     []string
					AllSSHKeys bool
					Arg        bool
				}{}),
				AllowAnonymousClient: true,
				Run: func(_ context.Context, _ interface{}) (i interface{}, e error) {
					return &core.SuccessResult{}, nil
				},
				ArgSpecs: core.ArgSpecs{
					{
						Name:       "ssh-key.{index}",
						OneOfGroup: "ssh",
						Required:   true,
					},
					{
						Name:       "all-ssh-keys",
						OneOfGroup: "ssh",
						Required:   true,
					},
					{
						Name: "arg",
					},
				},
			}),
			Cmd: "scw oneof arg=true",
			Check: core.TestCheckCombine(
				core.TestCheckExitCode(1),
				core.TestCheckError(
					errors.New("at least one argument from the 'ssh' group is required"),
				),
			),
		})(t)
	})
}
