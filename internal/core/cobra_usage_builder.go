package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"
	"text/tabwriter"

	"github.com/scaleway/scaleway-cli/internal/args"
	"github.com/scaleway/scaleway-cli/internal/interactive"
	"github.com/scaleway/scaleway-sdk-go/logger"
	"github.com/scaleway/scaleway-sdk-go/strcase"
)

const (
	sliceSchema = "{index}"
	mapSchema   = "{key}"
)

// buildUsageArgs builds usage args string.
// This string will be used by cobra usage template.
func buildUsageArgs(cmd *Command) string {
	var argsBuffer bytes.Buffer
	tw := tabwriter.NewWriter(&argsBuffer, 0, 0, 3, ' ', 0)

	err := _buildUsageArgs(tw, cmd.ArgSpecs, cmd.ArgsType, nil)
	if err != nil {
		// TODO: decide how to handle this error
		err = fmt.Errorf("building %v: %v", cmd.getPath(), err)
		logger.Debugf("%v", err)
	}
	tw.Flush()

	paramsStr := strings.TrimSuffix(argsBuffer.String(), "\n")

	return paramsStr
}

// _buildUsageArgs recursively builds the arg usage list.
// This should not be called directly.
func _buildUsageArgs(w io.Writer, argSpecs ArgSpecs, t reflect.Type, parentArgName []string) error {

	// related to protoc_gen_mordor.IsIgnoredFieldType()
	// TODO: make this relation explicit
	// TODO: decide what arguments to ignore
	ignoreKey := false
	if len(parentArgName) == 1 {
		firstWord := parentArgName[0] // pagination args are primary fields
		ignoredArgs := map[string]bool{
			"page":      true,
			"per-page":  true,
			"page-size": true,
		}
		_, ignoreKey = ignoredArgs[firstWord]
	}

	switch {
	case argSpecs == nil:
		return nil

	case ignoreKey:
		return nil

	case argSpecs.GetByName(strcase.ToBashArg(strings.Join(parentArgName, "."))) != nil:
		argName := strings.Join(parentArgName, ".")
		spec := argSpecs.GetByName(argName)
		if spec == nil {
			return nil
		}
		if spec.Default != nil {
			_, doc := spec.Default()
			argName = fmt.Sprintf("%s=%s", argName, doc)
		}
		if !spec.Required {
			argName = fmt.Sprintf("[%s]", argName)
		}
		_, err := fmt.Fprintf(w, "  %s\t%s\n", argName, _buildArgShort(spec))
		return err

	case t.Kind() == reflect.Ptr:
		return _buildUsageArgs(w, argSpecs, t.Elem(), parentArgName)

	case t.Kind() == reflect.Slice:
		return _buildUsageArgs(w, argSpecs, t.Elem(), append(parentArgName, sliceSchema))

	case t.Kind() == reflect.Map:
		return _buildUsageArgs(w, argSpecs, t.Elem(), append(parentArgName, mapSchema))

	case t.Kind() == reflect.Struct:
		for i := 0; i < t.NumField(); i++ {
			err := _buildUsageArgs(w, argSpecs, t.Field(i).Type, append(parentArgName, strcase.ToBashArg(t.Field(i).Name)))
			if err != nil {
				return err
			}
		}
		return nil

	default:
		return fmt.Errorf("'%v' has no usage and is not ignored", parentArgName)
	}
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
func buildExamples(cmd *Command) string {
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
			if err != nil {
				panic(fmt.Errorf("in command '%s', example '%s': %w", cmd.getPath(), cmdExample.Short, err))
			}

			// Build command line example.
			commandParts := []string{
				"scw",
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
