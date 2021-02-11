package core

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/scaleway/scaleway-cli/internal/interactive"
	"github.com/scaleway/scaleway-sdk-go/logger"
)

const (
	sliceSchema = "{index}"
	mapSchema   = "{key}"
)

// buildUsageArgs builds usage args string.
// If deprecated is true, true only deprecated argSpecs will be considered.
// This string will be used by cobra usage template.
func buildUsageArgs(ctx context.Context, cmd *Command, deprecated bool) string {
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
	for _, argSpec := range argSpecs {
		argSpecUsageLeftPart := argSpec.Name
		argSpecUsageRightPart := _buildArgShort(argSpec)
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
	var examples []string

	for _, cmdExample := range cmd.Examples {
		// Build title.
		title := fmt.Sprintf("  %s", cmdExample.Short)
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
