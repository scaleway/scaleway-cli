package core

import (
	"strings"

	cliConfig "github.com/scaleway/scaleway-cli/v2/internal/config"
	"github.com/spf13/cobra"
)

// cobraCompletionOutputFlag returns the default output when completing with cobra,
// as the output flag may be the word being completed (ex: "scw __complete instance server list -o ")
func cobraCompletionOutputFlag(args []string, outputFlag string) string {
	if len(args) > 1 && strings.HasPrefix(args[1], cobra.ShellCompRequestCmd) {
		return cliConfig.DefaultOutput
	}

	return outputFlag
}

// registerCobraFlagCompletions registers cobra completion for global flag values
func registerCobraFlagCompletions(rootCmd *cobra.Command) {
	_ = rootCmd.RegisterFlagCompletionFunc(
		"profile",
		func(cobraCmd *cobra.Command, _ []string, toComplete string) ([]cobra.Completion, cobra.ShellCompDirective) {
			profiles := AutocompleteProfileName()(cobraCmd.Context(), toComplete, nil)

			return profiles, cobra.ShellCompDirectiveNoFileComp
		},
	)
	_ = rootCmd.RegisterFlagCompletionFunc("output", cobra.FixedCompletions([]cobra.Completion{
		string(PrinterTypeHuman),
		string(PrinterTypeJSON),
		string(PrinterTypeYAML),
		string(PrinterTypeTemplate),
	}, cobra.ShellCompDirectiveNoFileComp))
}

// cobraValidArgsFunction returns a cobra completion function for the positional arguments of a command.
// Cobra handles namespaces, resources, verbs and flags, arguments are completed using the scw autocomplete engine.
func cobraValidArgsFunction(cmd *Command) cobra.CompletionFunc {
	return func(cobraCmd *cobra.Command, args []string, toComplete string) ([]cobra.Completion, cobra.ShellCompDirective) {
		ctx := cobraCmd.Context()

		// Keep the anonymous client if the profile cannot be loaded, static completion still works
		_ = createAuthenticatedClient(ctx, cmd)

		leftWords := []string{ExtractBinaryName(ctx)}
		for _, word := range []string{cmd.Namespace, cmd.Resource, cmd.Verb} {
			if word != "" {
				leftWords = append(leftWords, word)
			}
		}
		leftWords = append(leftWords, args...)

		res := AutoComplete(ctx, leftWords, toComplete, nil)

		directive := cobra.ShellCompDirectiveNoFileComp
		completions := make([]cobra.Completion, 0, len(res.Suggestions))
		for _, suggestion := range res.Suggestions {
			if strings.HasPrefix(suggestion, "-") {
				continue
			}
			if strings.HasSuffix(suggestion, "=") {
				directive |= cobra.ShellCompDirectiveNoSpace
			}
			completions = append(completions, suggestion)
		}

		return completions, directive
	}
}
