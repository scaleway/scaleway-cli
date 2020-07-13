package core

import (
	"context"
	"fmt"
	"reflect"
	"sort"
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

	// AllowAnonymousClient defines whether the SDK client can run the command without be authenticated.
	AllowAnonymousClient bool

	// DisableTelemetry disable telemetry for the command.
	DisableTelemetry bool

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

	// Interceptor are middleware func that can intercept context and args before they are sent to Run
	// You can combine multiple CommandInterceptor using AddInterceptors method.
	Interceptor CommandInterceptor

	// Run will be called to execute a command. It will receive a context and parsed argument.
	// Non-nil values returned by this method will be printed out.
	Run CommandRunner

	// WaitFunc will be called if non-nil when the -w (--wait) flag is passed.
	WaitFunc WaitFunc
}

// CommandPreValidateFunc allows to manipulate args before validation.
type CommandPreValidateFunc func(ctx context.Context, argsI interface{}) error

// CommandInterceptor allow to intercept and manipulate a runner arguments and return value.
// It can for example be used to change arguments type or catch runner errors.
type CommandInterceptor func(ctx context.Context, argsI interface{}, runner CommandRunner) (interface{}, error)

// CommandRunner returns the command response or an error.
type CommandRunner func(ctx context.Context, argsI interface{}) (interface{}, error)

// WaitFunc returns the updated response (respI if unchanged) or an error.
type WaitFunc func(ctx context.Context, argsI, respI interface{}) (interface{}, error)

const indexCommandSeparator = "."

// Override replaces or mutates the Command via a builder function.
func (c *Command) Override(builder func(command *Command) *Command) {
	// Assign the value in case the builder creates a new Command object.
	*c = *builder(c)
}

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

func (c *Command) GetCommandLine(binaryName string) string {
	return strings.Trim(binaryName+" "+strings.ReplaceAll(c.getPath(), indexCommandSeparator, " "), " ")
}

func (c *Command) GetUsage(binaryName string, commands *Commands) string {
	parts := []string{
		c.GetCommandLine(binaryName),
	}

	if commands.HasSubCommands(c) {
		parts = append(parts, "<command>")
	}
	if positionalArg := c.ArgSpecs.GetPositionalArg(); positionalArg != nil {
		parts = append(parts, "<"+positionalArg.Name+" ...>")
	}
	if len(c.ArgSpecs) > 0 {
		parts = append(parts, "[arg=value ...]")
	}
	return strings.Join(parts, " ")
}

// seeAlsosAsStr returns all See Alsos as a single string
func (c *Command) seeAlsosAsStr() string {
	var seeAlsos []string

	for _, cmdSeeAlso := range c.SeeAlsos {
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

// AddInterceptors add one or multiple interceptors to a command.
// These new interceptors will be added after the already present interceptors (if any).
func (c *Command) AddInterceptors(interceptors ...CommandInterceptor) {
	interceptors = append([]CommandInterceptor{c.Interceptor}, interceptors...)
	c.Interceptor = combineCommandInterceptor(interceptors...)
}

// Commands represent a list of CLI commands, with a index to allow searching.
type Commands struct {
	commands     []*Command
	commandIndex map[string]*Command
}

func NewCommands(cmds ...*Command) *Commands {
	c := &Commands{
		commands:     []*Command(nil),
		commandIndex: map[string]*Command{},
	}

	for _, cmd := range cmds {
		c.Add(cmd)
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

func (c *Commands) Add(cmd *Command) {
	c.commands = append(c.commands, cmd)
	c.commandIndex[cmd.getPath()] = cmd
}

func (c *Commands) Merge(cmds *Commands) {
	for _, cmd := range cmds.commands {
		c.Add(cmd)
	}
}

func (c *Commands) GetAll() []*Command {
	return c.commands
}

// find must take the command path, eg. find("instance","get","server")
func (c *Commands) find(path ...string) (*Command, bool) {
	cmd, exist := c.commandIndex[strings.Join(path, indexCommandSeparator)]
	if exist {
		return cmd, true
	}
	return nil, false
}

// GetSortedCommand returns a slice of commands sorted alphabetically
func (c *Commands) GetSortedCommand() []*Command {
	commands := make([]*Command, len(c.commands))
	copy(commands, c.commands)

	sort.Slice(commands, func(i, j int) bool {
		return fmt.Sprintf("%s %s %s", commands[i].Namespace, commands[i].Resource, commands[i].Verb) < fmt.Sprintf("%s %s %s", commands[j].Namespace, commands[j].Resource, commands[j].Verb)
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
		if cmd.Verb == "" && cmd.Namespace == command.Namespace && cmd.Resource == command.Resource {
			return true
		}
	}
	return false
}

func (c *Command) getHumanMarshalerOpt() *human.MarshalOpt {
	if c.View != nil {
		return c.View.getHumanMarshalerOpt()
	}
	return nil
}
