package core

import (
	"fmt"
	"sort"
	"strings"

	"github.com/scaleway/scaleway-cli/v2/internal/alias"
)

// Commands represent a list of CLI commands, with a index to allow searching.
type Commands struct {
	commands     []*Command
	commandIndex map[string]*Command
}

func NewCommands(cmds ...*Command) *Commands {
	c := &Commands{
		commands:     make([]*Command, 0, len(cmds)),
		commandIndex: make(map[string]*Command, len(cmds)),
	}

	for _, cmd := range cmds {
		c.Add(cmd)
	}

	return c
}

func NewCommandsMerge(cmdsList ...*Commands) *Commands {
	cmdCount := 0
	for _, cmds := range cmdsList {
		cmdCount += len(cmds.commands)
	}
	c := &Commands{
		commands:     make([]*Command, 0, cmdCount),
		commandIndex: make(map[string]*Command, cmdCount),
	}
	for _, cmds := range cmdsList {
		for _, cmd := range cmds.commands {
			c.Add(cmd)
		}
	}

	return c
}

func (c *Commands) MustFind(path ...string) *Command {
	cmd, exist := c.find(path...)
	if exist {
		return cmd
	}

	panic(fmt.Errorf("command %v not found", strings.Join(path, " ")))
}

func (c *Commands) Find(path ...string) *Command {
	cmd, exist := c.find(path...)
	if exist {
		return cmd
	}

	return nil
}

func (c *Commands) Remove(namespace, verb string) {
	for i := range c.commands {
		if c.commands[i].Namespace == namespace && c.commands[i].Verb == verb {
			c.commands = append(c.commands[:i], c.commands[i+1:]...)

			return
		}
	}
}

func (c *Commands) RemoveResource(namespace, resource string) {
	for i := range c.commands {
		if c.commands[i].Namespace == namespace && c.commands[i].Resource == resource &&
			c.commands[i].Verb == "" {
			c.commands = append(c.commands[:i], c.commands[i+1:]...)

			return
		}
	}
}

func (c *Commands) Add(cmd *Command) {
	c.commands = append(c.commands, cmd)
	c.commandIndex[cmd.getPath()] = cmd
}

func (c *Commands) Merge(cmds *Commands) {
	for _, cmd := range cmds.commands {
		c.Add(cmd)
	}
}

func (c *Commands) MergeAll(cmds ...*Commands) {
	for _, command := range cmds {
		c.Merge(command)
	}
}

func (c *Commands) GetAll() []*Command {
	return c.commands
}

// GetSortedCommand returns a slice of commands sorted alphabetically
func (c *Commands) GetSortedCommand() []*Command {
	commands := make([]*Command, len(c.commands))
	copy(commands, c.commands)
	sort.Slice(commands, func(i, j int) bool {
		return commands[i].signature() < commands[j].signature()
	})

	return commands
}

func (c *Commands) HasSubCommands(cmd *Command) bool {
	if cmd.Namespace != "" && cmd.Resource != "" && cmd.Verb != "" {
		return false
	}
	if cmd.Namespace == "" && cmd.Resource == "" && cmd.Verb == "" {
		return true
	}
	for _, command := range c.commands {
		if command == cmd {
			continue
		}
		if cmd.Resource == "" && cmd.Namespace == command.Namespace {
			return true
		}
		if cmd.Verb == "" && cmd.Namespace == command.Namespace &&
			cmd.Resource == command.Resource {
			return true
		}
	}

	return false
}

// AliasIsValidCommandChild returns true is alias is a valid child command of given command
// Useful for this case:
// isl => instance server list
// valid child of "instance"
// invalid child of "rdb instance"
func (c *Commands) AliasIsValidCommandChild(command *Command, alias alias.Alias) bool {
	// if alias is of size one, it means it cannot be a child
	if len(alias.Command) == 1 {
		return true
	}

	// if command is verb, it cannot have children
	if command.Verb != "" {
		return true
	}

	// if command is a resource, check command with alias' verb
	if command.Resource != "" {
		return c.Find(command.Namespace, command.Resource, alias.Command[1]) != nil
	}

	// if command is a namespace, check for alias' verb or resource
	if command.Namespace != "" {
		if len(alias.Command) > 2 {
			return c.Find(command.Namespace, alias.Command[1], alias.Command[2]) != nil
		}

		return c.Find(command.Namespace, alias.Command[1]) != nil
	}

	return false
}

// Copy return a copy of all commands
func (c *Commands) Copy() *Commands {
	newCommands := make([]*Command, len(c.commands))
	for i := range c.commands {
		newCommands[i] = c.commands[i].Copy()
	}

	return NewCommands(newCommands...)
}

// find must take the command path, eg. find("instance","get","server")
func (c *Commands) find(path ...string) (*Command, bool) {
	cmd, exist := c.commandIndex[strings.Join(path, indexCommandSeparator)]
	if exist {
		return cmd, true
	}

	return nil, false
}

// addAliases add valid aliases to a command
func (c *Commands) addAliases(command *Command, aliases []alias.Alias) {
	names := make([]string, 0, len(aliases))
	for i := range aliases {
		if c.AliasIsValidCommandChild(command, aliases[i]) && command.MatchAlias(aliases[i]) {
			names = append(names, aliases[i].Name)
		}
	}
	command.Aliases = append(command.Aliases, names...)
}

// applyAliases add resource aliases to each commands
func (c *Commands) applyAliases(config *alias.Config) {
	for _, command := range c.commands {
		aliases := []alias.Alias(nil)
		exists := false
		switch {
		case command.Verb != "":
			aliases, exists = config.ResolveAliasesByFirstWord(command.Verb)
		case command.Resource != "":
			aliases, exists = config.ResolveAliasesByFirstWord(command.Resource)
		case command.Namespace != "":
			aliases, exists = config.ResolveAliasesByFirstWord(command.Namespace)
		}
		if exists {
			c.addAliases(command, aliases)
		}
	}
}
