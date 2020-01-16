package core

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/scaleway/scaleway-cli/internal/human"
)

// Command represent a CLI command. From this higher level type we create Cobra command objects.
type Command struct {

	// Namespace is the top level entry point of a command. (e.g scw instance)
	Namespace string

	// Resource is the 2nd command level. Resources are nested in a namespace. (e.g scw instance server)
	Resource string

	// Verb is the 3rd command level. Verbs are nested in a resource. (e.g scw instance server list)
	Verb string

	// Short documentation.
	Short string

	// Long documentation.
	Long string

	// NoClient defines whether the SDK client is not required to run the command.
	NoClient bool

	// Hidden hides the command form usage and auto-complete.
	Hidden bool

	// ArgsType defines the type of argument for this command.
	ArgsType reflect.Type

	// ArgSpecs defines specifications for arguments.
	ArgSpecs ArgSpecs

	// View defines the View for this command.
	// It is used to create the different options for the different Marshalers.
	View *View

	// Examples defines Examples for this command.
	Examples []*Example

	// SeeAlsos presents commands related to this command.
	SeeAlsos []*SeeAlso

	// PreValidateFunc allows to manipulate args before validation
	PreValidateFunc CommandPreValidateFunc

	// ValidateFunc validates a command.
	// If nil, core.DefaultCommandValidateFunc is used by default.
	ValidateFunc CommandValidateFunc

	// Run will be called to execute a command. It will receive a context and parsed argument.
	// Non-nil values returned by this method will be printed out.
	Run CommandRunner

	// WaitFunc will be called if non-nil when the -w (--wait) flag is passed.
	WaitFunc WaitFunc
}

type CommandPreValidateFunc func(ctx context.Context, argsI interface{}) error
type CommandRunner func(ctx context.Context, argsI interface{}) (interface{}, error)
type WaitFunc func(ctx context.Context, argsI, respI interface{}) error

const indexCommandSeparator = "."

func (c *Command) getPath() string {
	path := []string(nil)
	if c.Namespace != "" {
		path = append(path, c.Namespace)
	}
	if c.Resource != "" {
		path = append(path, c.Resource)
	}
	if c.Verb != "" {
		path = append(path, c.Verb)
	}

	return strings.Join(path, indexCommandSeparator)
}

// Commands represent a list of CLI commands, with a index to allow searching.
type Commands struct {
	command      []*Command
	commandIndex map[string]*Command
}

func NewCommands(cmds ...*Command) *Commands {
	c := &Commands{
		command:      []*Command(nil),
		commandIndex: map[string]*Command{},
	}

	for _, cmd := range cmds {
		c.Add(cmd)
	}

	return c
}

// find must take the command path, eg. find("instance","get","server")
func (c *Commands) find(path ...string) (*Command, bool) {
	cmd, exist := c.commandIndex[strings.Join(path, indexCommandSeparator)]
	if exist {
		return cmd, true
	}
	return nil, false
}

func (c *Commands) MustFind(path ...string) *Command {
	cmd, exist := c.find(path...)
	if exist {
		return cmd
	}

	panic(fmt.Errorf("command %v not found", strings.Join(path, " ")))
}

func (c *Commands) Add(cmd *Command) {
	c.command = append(c.command, cmd)
	c.commandIndex[cmd.getPath()] = cmd
}

func (c *Commands) Merge(cmds *Commands) {
	for _, cmd := range cmds.command {
		c.Add(cmd)
	}
}

func (c *Command) getHumanMarshalerOpt() *human.MarshalOpt {
	if c.View != nil {
		return c.View.getHumanMarshalerOpt()
	}
	return nil
}

// seeAlsosAsStr returns all See Alsos as a single string
func (cmd *Command) seeAlsosAsStr() string {
	var seeAlsos []string

	for _, cmdSeeAlso := range cmd.SeeAlsos {
		short := fmt.Sprintf("  # %s", cmdSeeAlso.Short)
		commandStr := fmt.Sprintf("  %s", cmdSeeAlso.Command)

		seeAlsoLines := []string{
			short,
			commandStr,
		}
		seeAlsos = append(seeAlsos, strings.Join(seeAlsoLines, "\n"))
	}

	return strings.Join(seeAlsos, "\n\n")
}
