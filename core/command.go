package core

import (
	"context"
	"reflect"
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
type CommandPreValidateFunc func(ctx context.Context, argsI any) error

// CommandInterceptor allow to intercept and manipulate a runner arguments and return value.
// It can for example be used to change arguments type or catch runner errors.
type CommandInterceptor func(ctx context.Context, argsI any, runner CommandRunner) (any, error)

// CommandRunner returns the command response or an error.
type CommandRunner func(ctx context.Context, argsI any) (any, error)

// WaitFunc returns the updated response (respI if unchanged) or an error.
type WaitFunc func(ctx context.Context, argsI, respI any) (any, error)

const indexCommandSeparator = "."

// Override replaces or mutates the Command via a builder function.
func (c *Command) Override(builder func(command *Command) *Command) {
	// Assign the value in case the builder creates a new Command object.
	*c = *builder(c)
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

func (c *Command) DebugString() string {
	return c.getPath()
}

// get a signature to sort commands
func (c *Command) signature() string {
	return c.Namespace + " " + c.Resource + " " + c.Verb + " " + c.Short
}

func (c *Command) getHumanMarshalerOpt() *human.MarshalOpt {
	if c.View != nil {
		return c.View.getHumanMarshalerOpt()
	}

	return nil
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
