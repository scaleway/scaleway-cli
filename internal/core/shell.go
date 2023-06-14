//go:build !freebsd && !wasm

// shell is disabled on freebsd as current version of github.com/pkg/term@v1.1.0 is not compiling
package core

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/scaleway/scaleway-cli/v2/internal/cache"
	"github.com/scaleway/scaleway-cli/v2/internal/interactive"
	"github.com/scaleway/scaleway-cli/v2/internal/sentry"
	"github.com/spf13/cobra"
)

type Completer struct {
	ctx context.Context
}

type ShellSuggestion struct {
	Text string
	Arg  *ArgSpec
	Cmd  *Command
}

// lastArg returns last element of string or empty string
func lastArg(args []string) string {
	l := len(args)
	if l >= 2 {
		return args[l-1]
	}
	if l == 1 {
		return args[0]
	}
	return ""
}

// firstArg returns first element of list or empty string
func firstArg(args []string) string {
	l := len(args)
	if l >= 1 {
		return args[0]
	}
	return ""
}

// trimLastArg returns all arguments but the last one
// return a nil slice if there is no previous arguments
func trimLastArg(args []string) []string {
	l := len(args)
	if l > 1 {
		return args[:l-1]
	}
	return []string(nil)
}

// trimLastArg returns all arguments but the first one
// return a nil slice if there is no other arguments
func trimFirstArg(args []string) []string {
	l := len(args)
	if l > 1 {
		return args[1:]
	}
	return []string(nil)
}

// argIsOption returns if an argument is an option
func argIsOption(arg string) bool {
	return strings.Contains(arg, "=") || strings.Contains(arg, ".")
}

func argIsPositional(cmd *Command, arg string) bool {
	if cmd.Verb != "" && cmd.Verb == arg {
		return false
	} else if cmd.Resource != "" && cmd.Resource == arg {
		return false
	} else if cmd.Namespace != "" && cmd.Resource == arg {
		return false
	}

	return true
}

// removeOptions removes options from a list of argument
// ex: scw instance create name=myserver
// will be changed to: scw instance server create
func removeOptions(args []string) []string {
	filteredArgs := []string(nil)
	for _, arg := range args {
		if !argIsOption(arg) {
			filteredArgs = append(filteredArgs, arg)
		}
	}
	return filteredArgs
}

// optionToArgSpecName convert option to arg spec name
// from additional-volumes.0=hello to additional-volumes.{index}
// also with multiple indexes pools.0.kubelet-args. to pools.{index}.kubelet-args.{key}
func optionToArgSpecName(option string) string {
	optionName := strings.Split(option, "=")[0]

	if strings.Contains(optionName, ".") {
		// If option is formatted like "additional-volumes.0"
		// it should be converted to "additional-volumes.{index}
		elems := strings.Split(optionName, ".")
		for i := range elems {
			_, err := strconv.Atoi(elems[i])
			if err == nil {
				elems[i] = "{index}"
			}
		}
		if elems[len(elems)-1] == "" {
			elems[len(elems)-1] = "{key}"
		}
		optionName = strings.Join(elems, ".")
	}
	return optionName
}

// getCommand return command object from args and suggest
func getCommand(meta *meta, args []string, suggest string) *Command {
	rawCommand := removeOptions(args)
	suggestIsOption := argIsOption(suggest)

	if !suggestIsOption {
		rawCommand = append(rawCommand, suggest)
	}

	rawCommand = meta.CliConfig.Alias.ResolveAliases(rawCommand)

	// Find the closest command in case there is multiple positional arguments
	for ; len(rawCommand) > 1; rawCommand = rawCommand[:len(rawCommand)-1] {
		command, foundCommand := meta.Commands.find(rawCommand...)
		if foundCommand {
			return command
		}
	}
	return nil
}

// getSuggestDescription will return suggest description
// it will return command description if it is a command
// or option description if suggest is an option of a command
func getSuggestDescription(meta *meta, args []string, suggest string) string {
	command := getCommand(meta, args, suggest)
	if command == nil {
		return "command not found"
	}

	if argIsOption(suggest) {
		option := command.ArgSpecs.GetByName(optionToArgSpecName(suggest))
		if option != nil {
			return option.Short
		}
		return ""
	}

	if argIsPositional(command, suggest) {
		option := command.ArgSpecs.GetPositionalArg()
		if option != nil {
			return option.Short
		}
		return ""
	}

	// Should be a command, just use command short
	return command.Short
}

// sortOptions sorts options, putting required first then order alphabetically
func sortOptions(meta *meta, args []string, toSuggest string, suggestions []string) []string {
	command := getCommand(meta, args, toSuggest)
	if command == nil {
		return suggestions
	}

	argSpecs := []ShellSuggestion(nil)
	for _, suggest := range suggestions {
		argSpec := command.ArgSpecs.GetByName(optionToArgSpecName(suggest))
		argSpecs = append(argSpecs, ShellSuggestion{
			Text: suggest,
			Arg:  argSpec,
		})
	}

	sort.Slice(argSpecs, func(i, j int) bool {
		if argSpecs[i].Arg != nil && argSpecs[j].Arg != nil && argSpecs[i].Arg.Required != argSpecs[j].Arg.Required {
			return argSpecs[i].Arg.Required
		}
		return argSpecs[i].Text < argSpecs[j].Text
	})

	suggests := []string(nil)
	for _, argSpec := range argSpecs {
		suggests = append(suggests, argSpec.Text)
	}

	return suggests
}

// Complete returns the list of suggestion based on prompt content
func (c *Completer) Complete(d prompt.Document) []prompt.Suggest {
	// shell lib can request duplicate Complete request with empty strings as text
	// skipping to avoid cache reset
	if d.Text == "" {
		return nil
	}

	meta := extractMeta(c.ctx)

	argsBeforeCursor := meta.CliConfig.Alias.ResolveAliases(strings.Split(d.TextBeforeCursor(), " "))
	argsAfterCursor := meta.CliConfig.Alias.ResolveAliases(strings.Split(d.TextAfterCursor(), " "))
	currentArg := lastArg(argsBeforeCursor) + firstArg(argsAfterCursor)

	// leftArgs contains all arguments before the one with the cursor
	leftArgs := trimLastArg(argsBeforeCursor)
	// rightWords contains all words after the selected one
	rightWords := trimFirstArg(argsAfterCursor)

	leftWords := append([]string{"scw"}, leftArgs...)

	acr := AutoComplete(c.ctx, leftWords, currentArg, rightWords)

	suggestions := []prompt.Suggest(nil)
	rawSuggestions := []string(acr.Suggestions)

	// if first suggestion is an option, all suggestions should be options
	// we sort them
	if len(rawSuggestions) > 0 && argIsOption(rawSuggestions[0]) {
		rawSuggestions = sortOptions(meta, leftArgs, rawSuggestions[0], rawSuggestions)
	}

	for _, suggest := range rawSuggestions {
		suggestions = append(suggestions, prompt.Suggest{
			Text:        suggest,
			Description: getSuggestDescription(meta, leftArgs, suggest),
		})
	}

	return prompt.FilterHasPrefix(suggestions, currentArg, true)
}

func NewShellCompleter(ctx context.Context) *Completer {
	return &Completer{
		ctx: ctx,
	}
}

// shellExecutor returns the function that will execute command entered in shell
func shellExecutor(rootCmd *cobra.Command, printer *Printer, meta *meta) func(s string) {
	return func(s string) {
		args := strings.Fields(s)

		sentry.AddCommandContext(strings.Join(removeOptions(args), " "))

		rootCmd.SetArgs(meta.CliConfig.Alias.ResolveAliases(args))

		err := rootCmd.Execute()
		if err != nil {
			if _, ok := err.(*interactive.InterruptError); ok {
				return
			}

			printErr := printer.Print(meta.Client, err, nil)
			if printErr != nil {
				_, _ = fmt.Fprintln(os.Stderr, err)
			}

			return
		}

		// command is nil if it does not have a Run function
		// ex: instance -h
		if meta.command == nil {
			return
		}

		autoCompleteCache.Update(meta.command.Namespace)

		printErr := printer.Print(meta.Client, meta.result, meta.command.getHumanMarshalerOpt())
		if printErr != nil {
			_, _ = fmt.Fprintln(os.Stderr, printErr)
		}
	}
}

// Return the shell subcommand
func getShellCommand(rootCmd *cobra.Command) *cobra.Command {
	subcommands := rootCmd.Commands()
	for _, command := range subcommands {
		if command.Name() == "shell" {
			return command
		}
	}
	return nil
}

// RunShell will run an interactive shell that runs cobra commands
func RunShell(ctx context.Context, printer *Printer, meta *meta, rootCmd *cobra.Command, args []string) {
	autoCompleteCache = cache.New()
	completer := NewShellCompleter(ctx)

	shellCobraCommand := getShellCommand(rootCmd)
	shellCobraCommand.InitDefaultHelpFlag()
	_ = shellCobraCommand.ParseFlags(args)
	if isHelp, _ := shellCobraCommand.Flags().GetBool("help"); isHelp {
		shellCobraCommand.HelpFunc()(shellCobraCommand, args)
		return
	}

	// remove shell command so it cannot be called from shell
	rootCmd.RemoveCommand(shellCobraCommand)
	meta.Commands.Remove("shell", "")

	executor := shellExecutor(rootCmd, printer, meta)
	p := prompt.New(
		executor,
		completer.Complete,
		prompt.OptionPrefix(">>> "),
		prompt.OptionSuggestionBGColor(prompt.Purple),
		prompt.OptionSelectedSuggestionBGColor(prompt.Fuchsia),
		prompt.OptionSelectedSuggestionTextColor(prompt.White),
		prompt.OptionDescriptionBGColor(prompt.Purple),
		prompt.OptionSelectedDescriptionBGColor(prompt.Fuchsia),
		prompt.OptionSelectedDescriptionTextColor(prompt.White),
	)
	p.Run()
}
