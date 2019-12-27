package util

import (
	"bytes"
	"context"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/scaleway/scaleway-sdk-go/strcase"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/internal/args"
	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-cli/internal/tabwriter"
	"github.com/scaleway/scaleway-cli/internal/terminal"
	"github.com/scaleway/scaleway-sdk-go/logger"
)

func GetCommands() *core.Commands {
	return core.NewCommands(
		utilReport(),
	)
}

type NoCommandShortError struct {
	Command *core.Command
}

func (e *NoCommandShortError) Error() string {
	return fmt.Sprintf("[NO-SHORT] no short for '%s'", getCommandPath(e.Command))
}

type NoCommandLongError struct {
	Command *core.Command
}

func (e *NoCommandLongError) Error() string {
	return fmt.Sprintf("[NO-LONG] no long for '%s'", getCommandPath(e.Command))
}

type NoArgSpecShortError struct {
	Command *core.Command
	ArgSpec *core.ArgSpec
}

func (e *NoArgSpecShortError) Error() string {
	return fmt.Sprintf("[NO-ARG-SHORT] no short for '%s.%s'", getCommandPath(e.Command), e.ArgSpec.Name)
}

type ArgumentMissingFromUpdateError struct {
	Command *core.Command
	ArgSpec *core.ArgSpec
}

func (e *ArgumentMissingFromUpdateError) Error() string {
	path := e.Command.Namespace + "." + e.Command.Resource
	return fmt.Sprintf("[MISS-FROM-UPDATE] argument '%s' is used in '%v.create' but not in '%v.update' ", e.ArgSpec.Name, path, path)
}

type ArgumentMissingFromCreateError struct {
	Command *core.Command
	ArgSpec *core.ArgSpec
}

func (e *ArgumentMissingFromCreateError) Error() string {
	path := e.Command.Namespace + "." + e.Command.Resource
	return fmt.Sprintf("[MISS-FROM-CREATE] argument '%s' is used in '%v.update' but not in '%v.create' ", e.ArgSpec.Name, path, path)
}

type MissingArgumentInArgSpecsError struct {
	Command       *core.Command
	ParentArgName string
}

func (e *MissingArgumentInArgSpecsError) Error() string {
	return fmt.Sprintf("[MISSING-ARGUMENT] for '%s', parameter '%s' is defined in Request, is not ignored, and has no matching argument in ArgSpecs",
		getCommandPath(e.Command),
		e.ParentArgName)
}

// Represents a group of command having the same namespace and resource
type Group struct {
	namespace     string
	resource      string
	verbsCommands map[string]*core.Command
	args          map[string]bool
}

func utilReport() *core.Command {
	return &core.Command{
		Short:     ``,
		Long:      ``,
		Namespace: "util",
		Resource:  "report",
		ArgsType:  reflect.TypeOf(args.RawArgs{}),
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {

			// Slice of all commands
			commands := core.ExtractCommands(ctx)

			// Strings to be returned by this command
			strs := []string{}

			// Build errors for Longs and Short fields
			errors := []error{}
			for _, command := range commands.GetAll() {
				if command.Short == "" {
					errors = append(errors, &NoCommandShortError{
						Command: command,
					})
				}
				if command.Long == "" {
					errors = append(errors, &NoCommandLongError{
						Command: command,
					})
				}
				for _, argSpec := range command.ArgSpecs {
					if argSpec.Short == "" {
						errors = append(errors, &NoArgSpecShortError{
							Command: command,
							ArgSpec: argSpec,
						})
					}
				}
			}

			//
			for _, command := range commands.GetAll() {
				checkUsageArgsErrors := checkArgSpecsMatchRequests(command.ArgSpecs, command.ArgsType, command, nil, nil)
				errors = append(errors, checkUsageArgsErrors...)
			}

			// Append errors to string result
			strs = append(strs, "# Errors")
			for _, e := range errors {
				strs = append(strs, fmt.Sprintf("%v", e))
			}
			strs = append(strs, "")

			// Build groups
			groups := map[string]*Group{}
			for _, command := range commands.GetAll() {
				if command.Namespace != "" && command.Resource != "" && command.Verb != "" {
					path := command.Namespace + " " + command.Resource
					_, exists := groups[path]
					if !exists {
						groups[path] = &Group{
							namespace:     command.Namespace,
							resource:      command.Resource,
							verbsCommands: make(map[string]*core.Command),
							args:          make(map[string]bool),
						}
					}
					groups[path].verbsCommands[command.Verb] = command
					for _, arg := range command.ArgSpecs {
						groups[path].args[arg.Name] = true
					}
				}
			}

			// Append legend to read Arguments section
			strs = append(strs, "# Arguments")
			strs = append(strs, "██:               required argument")
			strs = append(strs, "░░:               optional argument")
			strs = append(strs, "[used-in-create]: argument is used in create, but not present in update")
			strs = append(strs, "[used-in-update]: argument is used in update, but not present in create")
			strs = append(strs, "")

			// Build & Append Arguments section
			for _, groupKey := range mapSortedKeys(groups) {
				group := groups[groupKey]
				if group.resource == "image" && group.namespace == "instance" {
					logger.Debugf("BOB2 %v", group.args)
				}
				strs = append(strs, "## scw "+group.namespace+" "+group.resource)
				grid := [][]string{}
				grid = append(grid, []string{})
				grid[0] = append(grid[0], "verbs/args")
				grid[0] = append(grid[0], mapSortedKeys(group.args)...)

				for i, verbKey := range mapSortedKeys(group.verbsCommands) {
					command := group.verbsCommands[verbKey]
					grid = append(grid, []string{})
					grid[i+1] = append(grid[i+1], verbKey)
					for _, argKey := range mapSortedKeys(group.args) {
						warnings := []string{}
						char := ""
						argSpec := command.ArgSpecs.GetByName(argKey)

						if argSpec == nil && verbKey == "update" && group.verbsCommands["create"].ArgSpecs.GetByName(argKey) != nil {
							warnings = append(warnings, "used-in-create")
							errors = append(errors, &ArgumentMissingFromUpdateError{
								Command: command,
								ArgSpec: argSpec,
							})

						}
						if _, updateExist := group.verbsCommands["update"]; updateExist {
							if argSpec == nil && verbKey == "create" && group.verbsCommands["update"].ArgSpecs.GetByName(argKey) != nil && argKey != group.verbsCommands["create"].Resource+"-id" {
								warnings = append(warnings, "used-in-update")
								errors = append(errors, &ArgumentMissingFromCreateError{
									Command: command,
									ArgSpec: argSpec,
								})
							}
						}
						if argSpec != nil {
							if argSpec.Required {
								char = "█"
								if len(warnings) > 0 {
									warnings = append(warnings, "required")
								}
							} else {
								char = "░"
								if len(warnings) > 0 {
									warnings = append(warnings, "optional")
								}
							}
						}
						str_ := ""
						if len(warnings) == 0 {
							for range argKey {
								str_ += char
							}
						} else {
							str_ = fmt.Sprintf("%v", warnings)
						}
						grid[i+1] = append(grid[i+1], str_)
					}
				}
				gridStr, _ := formatGrid(grid)
				strs = append(strs, gridStr)
				strs = append(strs, "")
			}

			// Build Errors summary
			errorsCountByType := make(map[string]int)
			for _, e := range errors {
				typeName := reflect.TypeOf(e).Elem().Name()
				if _, exists := errorsCountByType[typeName]; !exists {
					errorsCountByType[typeName] = 0
				}
				errorsCountByType[typeName] = errorsCountByType[typeName] + 1
			}

			// Append Errors summary
			strs = append(strs, "# Errors summary")
			for _, key := range mapSortedKeys(errorsCountByType) {
				count := errorsCountByType[key]
				strs = append(strs, fmt.Sprintf(" - %v: %v", key, count))
			}

			return strings.Join(strs, "\n"), nil
		},
	}
}

// copied from mordor/scaleway-cli/internal/human/marshal.go
// Padding between column
const colPadding = 2

// copied from mordor/scaleway-cli/internal/human/marshal.go
func formatGrid(grid [][]string) (string, error) {
	buffer := bytes.Buffer{}
	maxCols := computeMaxCols(grid)
	w := tabwriter.NewWriter(&buffer, 5, 1, colPadding, ' ', tabwriter.ANSIGraphicsRendition)
	for _, line := range grid {
		fmt.Fprintln(w, strings.Join(line[:maxCols], "\t"))
	}
	w.Flush()
	return strings.TrimSpace(buffer.String()), nil
}

// copied from mordor/scaleway-cli/internal/human/marshal.go
// computeMaxCols calculates how many row we can fit in terminal width.
func computeMaxCols(grid [][]string) int {
	maxCols := len(grid[0])
	// If we are not writing to Stdout or through a tty Stdout, returns max length
	if color.NoColor {
		return maxCols
	}
	width := terminal.GetWidth()
	colMaxSize := make([]int, len(grid[0]))
	for i := 0; i < len(grid); i++ {
		lineSize := 0
		for j := 0; j < maxCols; j++ {
			size := len(grid[i][j]) + colPadding
			if size >= colMaxSize[j] {
				colMaxSize[j] = size
			}
			lineSize += colMaxSize[j]
			if lineSize > width {
				maxCols = j
			}
		}
	}
	return maxCols
}

// copied from mordor/src/envoy-gen/utils.go
func mapSortedKeys(m interface{}) []string {
	v := reflect.ValueOf(m)
	mapKey := v.MapKeys()

	keys := make([]string, 0, v.Len())
	for _, k := range mapKey {
		keys = append(keys, k.String())
	}
	sort.Strings(keys)
	return keys
}

// copied from mordor/scaleway-cli/internal/core/command.go Command.getPath()
func getCommandPath(c *core.Command) string {
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
	return strings.Join(path, ".")
}

// copied from mordor/scaleway-cli/internal/core/cobra_usage_builder.go
const (
	sliceSchema = "{idx}"
	mapSchema   = "{key}"
)

// copied and modified from mordor/scaleway-cli/internal/core/cobra_usage_builder.go _buildUsageArgs()
func checkArgSpecsMatchRequests(argSpecs core.ArgSpecs, t reflect.Type, c *core.Command, parentArgName []string, errors []error) []error {

	if errors == nil {
		errors = []error{}
	}

	// related to protoc_gen_mordor.IsIgnoredFieldType()
	// TODO: make this relation explicit
	// TODO: decide what arguments to ignore
	ignoreKey := false
	if len(parentArgName) > 0 {
		lastKey := parentArgName[len(parentArgName)-1]
		ignoredKeys := map[string]bool{
			//"page":      true,
			//"per-page":  true,
			//"zone":      true,
			//"region":    true,
			//"page-size": true,
			//"organization": true,
		}
		_, ignoreKey = ignoredKeys[lastKey]
	}

	switch {
	case argSpecs == nil:
		return errors

	case ignoreKey:
		return errors

	case argSpecs.GetByName(strcase.ToBashArg(strings.Join(parentArgName, "."))) != nil:
		return errors

	case t.Kind() == reflect.Ptr:
		return checkArgSpecsMatchRequests(argSpecs, t.Elem(), c, parentArgName, errors)

	case t.Kind() == reflect.Slice:
		return checkArgSpecsMatchRequests(argSpecs, t.Elem(), c, append(parentArgName, sliceSchema), errors)

	case t.Kind() == reflect.Map:
		return checkArgSpecsMatchRequests(argSpecs, t.Elem(), c, append(parentArgName, mapSchema), errors)

	case t.Kind() == reflect.Struct:
		for i := 0; i < t.NumField(); i++ {
			errors = checkArgSpecsMatchRequests(argSpecs, t.Field(i).Type, c, append(parentArgName, strcase.ToBashArg(t.Field(i).Name)), errors)
		}
		return errors

	default:
		return append(errors, &MissingArgumentInArgSpecsError{
			Command:       c,
			ParentArgName: strings.Join(parentArgName, "."),
		})
	}
}
