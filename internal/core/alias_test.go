package core

import (
	"testing"

	"github.com/alecthomas/assert"
	"github.com/scaleway/scaleway-cli/v2/internal/alias"
)

func TestCommandMatchAlias(t *testing.T) {
	commandWithArg := &Command{
		Namespace: "first",
		Resource:  "command",
		ArgSpecs: ArgSpecs{
			{
				Name: "arg",
			},
		},
	}
	commandWithoutArg := &Command{
		Namespace: "second",
		Resource:  "command",
	}

	testAlias := alias.Alias{
		Name:    "alias",
		Command: []string{"command"},
	}

	assert.True(t, commandWithArg.matchAlias(testAlias))
	assert.True(t, commandWithoutArg.matchAlias(testAlias))

	testAliasWithArg := alias.Alias{
		Name:    "alias",
		Command: []string{"command", "arg=value"},
	}

	assert.True(t, commandWithArg.matchAlias(testAliasWithArg))
	assert.False(t, commandWithoutArg.matchAlias(testAliasWithArg))
}

func TestAliasChildCommand(t *testing.T) {
	namespace := &Command{
		Namespace: "namespace",
	}
	resource := &Command{
		Namespace: "namespace",
		Resource:  "first",
	}

	commands := NewCommands(
		namespace,
		resource,
	)

	validAlias := alias.Alias{
		Name:    "alias",
		Command: []string{"namespace", "first"},
	}

	assert.True(t, commands.aliasIsValidCommandChild(namespace, validAlias))

	invalidAlias := alias.Alias{
		Name:    "alias",
		Command: []string{"namespace", "random"},
	}

	assert.False(t, commands.aliasIsValidCommandChild(namespace, invalidAlias))
}
