package core

import (
	"strings"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
)

func init() {
	// we disable cobra command sorting to position important commands at the top when looking at the usage.
	cobra.EnableCommandSorting = false
	cobra.AddTemplateFunc("removeClientFlags", removeClientFlags)
}

// cobraBuilder will transform a []*Command to a valid Cobra root command.
// Cobra root command is a tree data struct. During the build process we
// use an index to attache leaf command to their parent.
type cobraBuilder struct {
	commands []*Command
	meta     *meta
}

// build creates the cobra root command.
func (b *cobraBuilder) build() *cobra.Command {
	index := map[string]*cobra.Command{}
	commandsIndex := map[string]*Command{}

	rootCmd := &cobra.Command{
		Use: "scw",

		// Do not display error with cobra, we handle it in bootstrap.
		SilenceErrors: true,

		// Do not display usage on error.
		SilenceUsage: true,
	}
	rootCmd.SetOut(b.meta.stderr)

	for _, cmd := range b.commands {
		// If namespace command has not yet been created. We create an empty cobra command to allow leaf to be attached.
		if _, namespaceExist := index[cmd.Namespace]; !namespaceExist {
			cobraCmd := &cobra.Command{Use: cmd.Namespace}
			index[cmd.Namespace] = cobraCmd
			commandsIndex[cmd.Namespace] = cmd
			rootCmd.AddCommand(cobraCmd)
		}

		// If Resource is empty, the command represent a namespace directly.
		if cmd.Resource == "" {
			continue
		}

		// Same as namespace but with resource.
		resourceKey := strings.Join([]string{cmd.Namespace, cmd.Resource}, indexCommandSeparator)
		if _, resourceExist := index[resourceKey]; !resourceExist {
			cobraCmd := &cobra.Command{Use: cmd.Resource}
			index[resourceKey] = cobraCmd
			commandsIndex[resourceKey] = cmd
			index[cmd.Namespace].AddCommand(cobraCmd)
		}

		if cmd.Verb == "" {
			continue
		}

		// Same as namespace but with verbs.
		verbKey := strings.Join([]string{cmd.Namespace, cmd.Resource, cmd.Verb}, indexCommandSeparator)
		if _, verbExist := index[verbKey]; !verbExist {
			cobraCmd := &cobra.Command{Use: cmd.Verb}
			index[verbKey] = cobraCmd
			commandsIndex[verbKey] = cmd
			index[resourceKey].AddCommand(cobraCmd)
		}
	}

	for k := range index {
		b.hydrateCobra(index[k], commandsIndex[k])
	}

	return rootCmd
}

// hydrateCobra hydrates a cobra command from a *Command.
// Field like Short, Long will be copied over.
// More complex field like PreRun or Run will also be generated if needed.
func (b *cobraBuilder) hydrateCobra(cobraCmd *cobra.Command, cmd *Command) {
	cobraCmd.Short = cmd.Short
	cobraCmd.Long = cmd.Long
	cobraCmd.Hidden = cmd.Hidden

	cobraCmd.SetUsageTemplate(usageTemplate)

	if cobraCmd.Annotations == nil {
		cobraCmd.Annotations = make(map[string]string)
	}

	cobraCmd.Annotations["NoClient"] = "false"
	if cmd.NoClient {
		cobraCmd.Annotations["NoClient"] = "true"
	}

	if cmd.ArgsType != nil {
		cobraCmd.Annotations["UsageArgs"] = buildUsageArgs(cmd)
	}

	if cmd.Examples != nil {
		cobraCmd.Annotations["Examples"] = buildExamples(cmd)
	}

	if cmd.SeeAlsos != nil {
		cobraCmd.Annotations["SeeAlsos"] = cmd.seeAlsosAsStr()
	}

	cobraCmd.PreRunE = cobraPreRunInitMeta(newMetaContext(b.meta), cmd)

	if cmd.Run != nil {
		cobraCmd.RunE = cobraRun(newMetaContext(b.meta), cmd)
	}

	if cmd.WaitFunc != nil {
		cobraCmd.PersistentFlags().BoolP("wait", "w", false, "wait until the "+cmd.Resource+" is ready")
	}

	cobraCmd.Annotations["CommandUsage"] = strings.ReplaceAll(cobraCmd.CommandPath(), "scw", "scw [global-flags]")
	if cobraCmd.HasAvailableSubCommands() || len(cobraCmd.Commands()) > 0 {
		cobraCmd.Annotations["CommandUsage"] += " <command>"
	}
	if cobraCmd.HasAvailableLocalFlags() || cobraCmd.HasAvailableFlags() || cobraCmd.LocalFlags() != nil {
		cobraCmd.Annotations["CommandUsage"] += " [flags]"
	}
	if len(cmd.ArgSpecs) > 0 {
		cobraCmd.Annotations["CommandUsage"] += " [arg=value ...]"
	}
}

const usageTemplate = `USAGE:
  {{.Annotations.CommandUsage}}{{if gt (len .Aliases) 0}}

ALIASES:
  {{.NameAndAliases}}{{end}}{{if .Annotations.Examples}}

EXAMPLES:
{{.Annotations.Examples}}{{end}}{{if .Annotations.UsageArgs}}

ARGS:
{{.Annotations.UsageArgs}}{{end}}{{if .HasAvailableSubCommands}}

AVAILABLE COMMANDS:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{if .Short}}{{.Short}}{{end}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

FLAGS:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

GLOBAL FLAGS:
{{ (removeClientFlags .InheritedFlags .Annotations.NoClient).FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .Annotations.SeeAlsos}}

SEE ALSO:
{{.Annotations.SeeAlsos}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`

func removeClientFlags(flags *flag.FlagSet, noClient string) *flag.FlagSet {
	if noClient == "true" {
		panicOnError(flags.MarkHidden("access-key"))
		panicOnError(flags.MarkHidden("secret-key"))
	}
	return flags
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}
