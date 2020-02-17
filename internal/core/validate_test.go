package core

import (
	"fmt"
	"testing"

	"github.com/alecthomas/assert"
)

func Test_DefaultCommandValidateFunc(t *testing.T) {

	type ReqElement struct {
		ID   int
		Name string
	}

	type Req struct {
		Elements map[string]ReqElement
	}

	commands := NewCommands(
		&Command{
			Namespace: "test",
			Resource:  "flower",
			Verb:      "create",
			ArgSpecs: ArgSpecs{
				{
					Name: "elements.{key}.id",
				},
				{
					Name: "elements.{key}.name",
					ValidateFunc: func(argSpec *ArgSpec, value interface{}) error {
						return fmt.Errorf("arg validation called")
					},
				},
			},
		},
	)

	myMap := map[string]ReqElement{}
	myMap["first"] = ReqElement{
		ID:   1,
		Name: "first",
	}
	myMap["second"] = ReqElement{
		ID:   2,
		Name: "second",
	}
	cmdArgs := &Req{
		Elements: myMap,
	}

	err := DefaultCommandValidateFunc()(commands.MustFind("test", "flower", "create"), cmdArgs)
	assert.Equal(t, fmt.Errorf("arg validation called"), err)
}
