package core

import (
	"context"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/scaleway/scaleway-cli/v2/core/human"
	"github.com/scaleway/scaleway-cli/v2/internal/alias"
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

	// DisableAfterChecks disable checks that run after the command to avoid superfluous message
	DisableAfterChecks bool

	// Hidden hides the command form usage and auto-complete.
	Hidden bool

	// ArgsType defines the type of argument for this command.
	ArgsType reflect.Type

	// ArgSpecs defines specifications for arguments.
	ArgSpecs ArgSpecs

	// AcceptMultiplePositionalArgs defines whether the command can accept multiple positional arguments.
	// If enabled, positional argument is expected to be a list.
	AcceptMultiplePositionalArgs bool

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

	// WebURL will be used as url to open when the --web flag is passed
	// Can contain template of values in request, ex: "url/{{ .Zone }}/{{ .ResourceID }}"
	WebURL string

	// WaitUsage override the usage for the -w (--wait) flag
	WaitUsage string

	// Aliases contains a list of aliases for a command
	Aliases []string
	// cache command path
	path string

	// Groups contains a list of groups IDs
	Groups []string
	//
	Deprecated bool
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
	if c.path != "" {
		return c.path
	}
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

	c.path = strings.Join(path, indexCommandSeparator)

	return c.path
}

func (c *Command) GetCommandLine(binaryName string) string {
	return strings.Trim(
		binaryName+" "+strings.ReplaceAll(c.getPath(), indexCommandSeparator, " "),
		" ",
	)
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
	seeAlsos := make([]string, 0, len(c.SeeAlsos))

	for _, cmdSeeAlso := range c.SeeAlsos {
		short := "  # " + cmdSeeAlso.Short
		commandStr := "  " + cmdSeeAlso.Command

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
	c.Interceptor = CombineCommandInterceptor(interceptors...)
}

// MatchAlias returns true if the alias can be used for this command
func (c *Command) MatchAlias(alias alias.Alias) bool {
	if len(c.ArgSpecs) == 0 {
		// command should be either a namespace or a resource
		// We need to check if child commands match this alias
		return true
	}

	for _, aliasArg := range alias.Args() {
		if c.ArgSpecs.GetByName(aliasArg) == nil {
			return false
		}
	}

	return true
}

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

func (c *Command) getHumanMarshalerOpt() *human.MarshalOpt {
	if c.View != nil {
		return c.View.getHumanMarshalerOpt()
	}

	return nil
}

// get a signature to sort commands
func (c *Command) signature() string {
	return c.Namespace + " " + c.Resource + " " + c.Verb + " " + c.Short
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

// Copy returns a copy of a command
func (c *Command) Copy() *Command {
	newCommand := *c
	newCommand.Aliases = append([]string(nil), c.Aliases...)
	newCommand.Examples = make([]*Example, len(c.Examples))
	for i := range c.Examples {
		e := *c.Examples[i]
		newCommand.Examples[i] = &e
	}
	newCommand.SeeAlsos = make([]*SeeAlso, len(c.SeeAlsos))
	for i := range c.SeeAlsos {
		sa := *c.SeeAlsos[i]
		newCommand.SeeAlsos[i] = &sa
	}

	return &newCommand
}

// Copy return a copy of all commands
func (c *Commands) Copy() *Commands {
	newCommands := make([]*Command, len(c.commands))
	for i := range c.commands {
		newCommands[i] = c.commands[i].Copy()
	}

	return NewCommands(newCommands...)
}
