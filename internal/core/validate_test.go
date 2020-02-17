package core

import (
	"fmt"
	"testing"

	"github.com/alecthomas/assert"
)

type Element struct {
	ID            int
	Name          string
	ElementsMap   map[string]Element
	ElementsSlice []Element
}

type elementCustom struct {
	*Element
	Short string
}

func Test_DefaultCommandValidateFunc(t *testing.T) {

	type TestCase struct {
		command         *Command
		parsedArguments interface{}
	}

	run := func(testCase TestCase) func(t *testing.T) {
		return func(t *testing.T) {
			err := DefaultCommandValidateFunc()(testCase.command, testCase.parsedArguments)
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
