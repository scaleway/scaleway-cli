package core

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/scaleway/scaleway-cli/internal/args"
	"github.com/scaleway/scaleway-cli/internal/interactive"
	"github.com/scaleway/scaleway-sdk-go/logger"
)

const (
	sliceSchema = "{index}"
	mapSchema   = "{key}"
)

// buildUsageArgs builds usage args string.
// This string will be used by cobra usage template.
func buildUsageArgs(ctx context.Context, cmd *Command) string {
	var argsBuffer bytes.Buffer
	tw := tabwriter.NewWriter(&argsBuffer, 0, 0, 3, ' ', 0)

	err := _buildUsageArgs(tw, cmd.ArgSpecs)
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
func _buildUsageArgs(w io.Writer, argSpecs ArgSpecs) error {
	for _, argSpec := range argSpecs {
		argSpecUsageLeftPart := argSpec.Name
		if argSpec.Default != nil {
			_, doc := argSpec.Default()
			argSpecUsageLeftPart = fmt.Sprintf("%s=%s", argSpecUsageLeftPart, doc)
		}
		if !argSpec.Required && !argSpec.Positional {
			argSpecUsageLeftPart = fmt.Sprintf("[%s]", argSpecUsageLeftPart)
		}

		_, err := fmt.Fprintf(w, "  %s\t%s\n", argSpecUsageLeftPart, _buildArgShort(argSpec))
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
		commandLine := ""
		switch {
		case cmdExample.Raw != "":
			commandLine = cmdExample.Raw
			commandLine = strings.Trim(commandLine, "\n")
			commandLine = interactive.RemoveIndent(commandLine)
		case cmdExample.Request != "":
			//  Query and path parameters don't have json tag,
			//  so we need to enforce a JSON tag on every field to make this work.
			var cmdArgs = newObjectWithForcedJSONTags(cmd.ArgsType)
			if err := json.Unmarshal([]byte(cmdExample.Request), cmdArgs); err != nil {
				panic(fmt.Errorf("in command '%s', example '%s': %w", cmd.getPath(), cmdExample.Short, err))
			}
			var cmdArgsAsStrings, err = args.MarshalStruct(cmdArgs)
			positionalArg := cmd.ArgSpecs.GetPositionalArg()
			if positionalArg != nil {
				for i, cmdArg := range cmdArgsAsStrings {
					if !strings.HasPrefix(cmdArg, positionalArg.Prefix()) {
						continue
					}
					cmdArgsAsStrings[i] = strings.TrimLeft(cmdArg, positionalArg.Prefix())
					// Switch the positional args with args at position 0 to make sure it is always at the beginning
					cmdArgsAsStrings[0], cmdArgsAsStrings[i] = cmdArgsAsStrings[i], cmdArgsAsStrings[0]
					break
				}
			}

			if err != nil {
				panic(fmt.Errorf("in command '%s', example '%s': %w", cmd.getPath(), cmdExample.Short, err))
			}

			// Build command line example.
			commandParts := []string{
				binaryName,
				cmd.Namespace,
				cmd.Resource,
				cmd.Verb,
			}
			commandParts = append(commandParts, cmdArgsAsStrings...)
			commandLine = strings.Join(commandParts, " ")
		default:
			panic(fmt.Errorf("in command '%s' invalid example '%s', it should either have a Request or a Raw", cmd.getPath(), cmdExample.Short))
		}

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
