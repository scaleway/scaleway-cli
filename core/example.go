package core

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/scaleway/scaleway-cli/v2/internal/args"
	"github.com/scaleway/scaleway-cli/v2/internal/interactive"
)

// Example represents an example for the usage of a CLI command.
type Example struct {
	// Short is the title given to the example.
	Short string

	// ArgsJSON is a JSON encoded representation of the request used in the example. Only one of ArgsJSON or Raw should be provided.
	ArgsJSON string

	// Raw is a raw example. Only one of ArgsJSON or Raw should be provided.
	Raw string
}

func (e *Example) GetCommandLine(binaryName string, cmd *Command) string {
	switch {
	case e.Raw != "":
		res := e.Raw
		res = strings.Trim(res, "\n")
		res = interactive.RemoveIndent(res)

		return res
	case e.ArgsJSON != "":
		//  Query and path parameters don't have json tag,
		//  so we need to enforce a JSON tag on every field to make this work.
		cmdArgs := newObjectWithForcedJSONTags(cmd.ArgsType)
		if err := json.Unmarshal([]byte(e.ArgsJSON), cmdArgs); err != nil {
			panic(fmt.Errorf("in command '%s', example '%s': %w", cmd.getPath(), cmd.Short, err))
		}
		cmdArgsAsStrings, err := args.MarshalStruct(cmdArgs)
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
			panic(fmt.Errorf("in command '%s', example '%s': %w", cmd.getPath(), cmd.Short, err))
		}

		// Build command line example.
		commandParts := []string{
			binaryName,
			cmd.Namespace,
			cmd.Resource,
			cmd.Verb,
		}
		commandParts = append(commandParts, cmdArgsAsStrings...)

		return strings.Join(commandParts, " ")
	default:
		panic(
			fmt.Errorf(
				"in command '%s' invalid example '%s', it should either have a ArgsJSON or a Raw",
				cmd.getPath(),
				cmd.Short,
			),
		)
	}
}
