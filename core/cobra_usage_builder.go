package core

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/scaleway/scaleway-cli/v2/internal/interactive"
	"github.com/scaleway/scaleway-sdk-go/logger"
	"github.com/spf13/cobra"
)

const (
	sliceSchema = "{index}"
	mapSchema   = "{key}"
)

func buildUsageAliases(ctx context.Context, cmd *Command) string {
	var aliasesBuffer bytes.Buffer
	tw := tabwriter.NewWriter(&aliasesBuffer, 0, 0, 2, ' ', 0)

	// Copy and sort alias list
	aliases := make([]string, len(cmd.Aliases))
	copy(aliases, cmd.Aliases)
	sort.Strings(aliases)

	aliasCfg := ExtractAliases(ctx)
	for _, aliasName := range aliases {
		_, _ = fmt.Fprintf(
			tw,
			" %s\t%s\n",
			aliasName,
			strings.Join(aliasCfg.GetAlias(aliasName), " "),
		)
	}
	tw.Flush()

	aliasesStr := strings.TrimSuffix(aliasesBuffer.String(), "\n")

	return aliasesStr
}

// BuildUsageArgs builds usage args string.
// If deprecated is true, true only deprecated argSpecs will be considered.
// This string will be used by cobra usage template.
func BuildUsageArgs(ctx context.Context, cmd *Command, deprecated bool) string {
	var argsBuffer bytes.Buffer
	tw := tabwriter.NewWriter(&argsBuffer, 0, 0, 3, ' ', 0)

	// Filter deprecated argSpecs.
	argSpecs := cmd.ArgSpecs.GetDeprecated(deprecated)

	err := _buildUsageArgs(ctx, tw, argSpecs)
	if err != nil {
		// TODO: decide how to handle this error
		err = fmt.Errorf("building %v: %v", cmd.getPath(), err)
		logger.Debugf("%v", err)
	}
	tw.Flush()

	paramsStr := strings.TrimSuffix(argsBuffer.String(), "\n")

	return paramsStr
}

// _buildUsageArgs builds the arg usage list.
// This should not be called directly.
func _buildUsageArgs(ctx context.Context, w io.Writer, argSpecs ArgSpecs) error {
	inOneOfGroup := false
	lastOneOfGroup := ""

	for _, argSpec := range argSpecs {
		argSpecUsageLeftPart := argSpec.Name
		argSpecUsageRightPart := _buildArgShort(argSpec)

		if argSpec.OneOfGroup != "" {
			if argSpec.OneOfGroup != lastOneOfGroup {
				inOneOfGroup = true
				lastOneOfGroup = argSpec.OneOfGroup
				_, err := fmt.Fprintf(w, "  %s (one of):\n", argSpec.OneOfGroup)
				if err != nil {
					return err
				}
			}
		} else {
			inOneOfGroup = false
			lastOneOfGroup = ""
		}

		if argSpec.Default != nil {
			_, doc := argSpec.Default(ctx)
			argSpecUsageLeftPart = fmt.Sprintf("%s=%s", argSpecUsageLeftPart, doc)
		}
		if !argSpec.Required && !argSpec.Positional {
			argSpecUsageLeftPart = fmt.Sprintf("[%s]", argSpecUsageLeftPart)
		}
		if argSpec.CanLoadFile {
			argSpecUsageRightPart += " (Support file loading with @/path/to/file)"
		}
		if inOneOfGroup {
			argSpecUsageLeftPart = "  " + argSpecUsageLeftPart
		}

		_, err := fmt.Fprintf(w, "  %s\t%s\n", argSpecUsageLeftPart, argSpecUsageRightPart)
		if err != nil {
			return err
		}
	}

	return nil
}

// _buildArgShort builds the arg short string.
// This should not be called directly.
func _buildArgShort(as *ArgSpec) string {
	if len(as.EnumValues) > 0 {
		return fmt.Sprintf("%s (%s)", as.Short, strings.Join(as.EnumValues, " | "))
	}

	return as.Short
}

// buildExamples builds usage examples string.
// This string will be used by cobra usage template.
func buildExamples(binaryName string, cmd *Command) string {
	// Build the examples array.
	examples := make([]string, 0, len(cmd.Examples))

	for _, cmdExample := range cmd.Examples {
		// Build title.
		title := "  " + cmdExample.Short
		commandLine := cmdExample.GetCommandLine(binaryName, cmd)

		commandLine = interactive.Indent(commandLine, 4)
		commandLine = strings.Trim(commandLine, "\n")

		// Add the whole example as a single string.
		exampleLines := []string{
			title,
			commandLine,
		}
		examples = append(examples, strings.Join(exampleLines, "\n"))
	}

	// Return a single string for all examples.
	return strings.Join(examples, "\n\n")
}

// usageFuncBuilder returns the usage function that will be used by cobra to print usage,
// the builder also takes a function that will fill annotations used by the usage template,
// this is done like this to avoid build annotations for each command if not required
func usageFuncBuilder(cmd *cobra.Command, annotationBuilder func()) func(*cobra.Command) error {
	return func(command *cobra.Command) error {
		annotationBuilder()
		// after building annotation we remove this function as we prefer to use default UsageFunc
		cmd.SetUsageFunc(nil)

		return cmd.UsageFunc()(command)
	}
}

func orderCobraCommands(cobraCommands []*cobra.Command) []*cobra.Command {
	commands := make([]*cobra.Command, len(cobraCommands))
	copy(commands, cobraCommands)

	sort.Slice(commands, func(i, j int) bool {
		deprecatedI := commands[i].Deprecated != ""
		deprecatedJ := commands[j].Deprecated != ""
		if deprecatedI == deprecatedJ {
			return commands[i].Use < commands[j].Use
		}

		return !deprecatedI && deprecatedJ
	})

	return commands
}

func orderCobraGroups(cobraGroups []*cobra.Group) []*cobra.Group {
	groups := make([]*cobra.Group, len(cobraGroups))
	copy(groups, cobraGroups)

	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Title < groups[j].Title
	})

	return groups
}

func getCobraCommandsGroups(cobraCommands []*cobra.Command) []*cobra.Group {
	var groups []*cobra.Group
	addedGroups := make(map[string]struct{})

	for _, cobraCommand := range cobraCommands {
		if !cobraCommand.IsAvailableCommand() {
			continue
		}

		for _, group := range cobraCommand.Groups() {
			if _, ok := addedGroups[group.ID]; ok {
				continue
			}

			addedGroups[group.ID] = struct{}{}
			groups = append(groups, group)
		}
	}

	return groups
}
