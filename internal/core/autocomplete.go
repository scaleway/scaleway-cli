package core

import (
	"context"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/scaleway/scaleway-cli/internal/args"
)

// AutocompleteSuggestions is a list of words to be set to the shell as autocomplete suggestions.
type AutocompleteSuggestions []string

// AutocompleteResponse contains the autocomplete suggestions
type AutocompleteResponse struct {
	Suggestions AutocompleteSuggestions
}

const (
	// positionalValueNodeID are flag values or positional argument.
	// E.g.: `scw test create <value> --flag <value>`
	positionalValueNodeID = "*"
)

type AutoCompleteNodeType uint

const (
	AutoCompleteNodeTypeCommand AutoCompleteNodeType = iota
	AutoCompleteNodeTypeArgument
	AutoCompleteNodeTypeFlag
	AutoCompleteNodeTypeFlagValueConst
	AutoCompleteNodeTypeFlagValueVariable
)

// AutoCompleteArgFunc is the function called to complete arguments values.
// It is retrieved from core.ArgSpec.AutoCompleteFunc.
type AutoCompleteArgFunc func(ctx context.Context, prefix string) AutocompleteSuggestions

// AutoCompleteNode is a node in the AutoComplete Tree.
// An AutoCompleteNode can either represent a command, a subcommand, or a command argument.
type AutoCompleteNode struct {
	Children map[string]*AutoCompleteNode
	Command  *Command
	ArgSpec  *ArgSpec
	Type     AutoCompleteNodeType

	// Name of the current node. Useful for debugging.
	Name string
}

type FlagSpec struct {
	Name             string
	HasVariableValue bool
	EnumValues       []string
}

func (node *AutoCompleteNode) addGlobalFlags() {
	printerTypes := []string{
		PrinterTypeHuman.String(),
		PrinterTypeJSON.String(),
		PrinterTypeYAML.String(),
		PrinterTypeTemplate.String(),
	}

	node.Children["-D"] = NewAutoCompleteFlagNode(node, &FlagSpec{
		Name: "-D",
	})
	node.Children["--debug"] = NewAutoCompleteFlagNode(node, &FlagSpec{
		Name: "--debug",
	})
	node.Children["-h"] = NewAutoCompleteFlagNode(node, &FlagSpec{
		Name: "-h",
	})
	node.Children["--help"] = NewAutoCompleteFlagNode(node, &FlagSpec{
		Name: "--help",
	})
	node.Children["-o"] = NewAutoCompleteFlagNode(node, &FlagSpec{
		Name:       "-o",
		EnumValues: printerTypes,
	})
	node.Children["--output"] = NewAutoCompleteFlagNode(node, &FlagSpec{
		Name:       "--output",
		EnumValues: printerTypes,
	})
	node.Children["-p"] = NewAutoCompleteFlagNode(node, &FlagSpec{
		Name:             "-p",
		HasVariableValue: true,
	})
	node.Children["--profile"] = NewAutoCompleteFlagNode(node, &FlagSpec{
		Name:             "--profile",
		HasVariableValue: true,
	})
}

// newAutoCompleteResponse builds a new AutocompleteResponse
func newAutoCompleteResponse(suggestions []string) *AutocompleteResponse {
	sort.Strings(suggestions)
	return &AutocompleteResponse{
		Suggestions: suggestions,
	}
}

// NewAutoCompleteCommandNode creates a new node corresponding to a command or subcommand.
// These nodes are not necessarily leaf nodes.
func NewAutoCompleteCommandNode() *AutoCompleteNode {
	return &AutoCompleteNode{
		Children: make(map[string]*AutoCompleteNode),
		Type:     AutoCompleteNodeTypeCommand,
	}
}

// NewArgAutoCompleteNode creates a new node corresponding to a command argument.
// These nodes are leaf nodes.
func NewAutoCompleteArgNode(cmd *Command, argSpec *ArgSpec) *AutoCompleteNode {
	return &AutoCompleteNode{
		Children: make(map[string]*AutoCompleteNode),
		ArgSpec:  argSpec,
		Type:     AutoCompleteNodeTypeArgument,
		Command:  cmd,
	}
}

// NewFlagAutoCompleteNode returns a node representing a Flag.
// It creates the children node with possible values if they exist.
// It sets parent.children as children of the lowest nodes:
//  the lowest node is the flag if it has no possible value ;
//  or the lowest nodes are the possible values if the exist.
func NewAutoCompleteFlagNode(parent *AutoCompleteNode, flagSpec *FlagSpec) *AutoCompleteNode {
	node := &AutoCompleteNode{
		Children: make(map[string]*AutoCompleteNode),
		Type:     AutoCompleteNodeTypeFlag,
		Name:     flagSpec.Name,
	}
	if flagSpec.HasVariableValue {
		node.Children[positionalValueNodeID] = &AutoCompleteNode{
			Children: parent.Children,
			Type:     AutoCompleteNodeTypeFlagValueVariable,
		}
	}
	for _, value := range flagSpec.EnumValues {
		node.Children[value] = &AutoCompleteNode{
			Children: parent.Children,
			Type:     AutoCompleteNodeTypeFlagValueConst,
		}
	}
	if len(node.Children) == 0 {
		node.Children = parent.Children
	}
	return node
}

// GetChildOrCreate search a child node by name,
// and either returns it if found
// or create a new child with the given name, and returns it.
func (node *AutoCompleteNode) GetChildOrCreate(name string) *AutoCompleteNode {
	if _, exist := node.Children[name]; !exist {
		node.Children[name] = NewAutoCompleteCommandNode()
	}
	return node.Children[name]
}

// GetChildMatch returns a child for a node if the child exists for this node.
// 3 types of children are supported :
// - command:                             command
// - singular argument name:              argument=
// - plural argument name + alphanumeric: arguments.key1=
func (node *AutoCompleteNode) GetChildMatch(name string) (*AutoCompleteNode, bool) {
	for key, child := range node.Children {
		if key == positionalValueNodeID {
			continue
		}
		key = "^" + key + "$"
		key = strings.ReplaceAll(key, ".", "\\.")
		key = strings.ReplaceAll(key, sliceSchema, "[0-9]+")
		key = strings.ReplaceAll(key, mapSchema, "[0-9a-zA-Z-]+")
		r := regexp.MustCompile(key)
		if r.MatchString(name) {
			return child, true
		}
	}
	return nil, false
}

// isLeafCommand returns true only if n is a node with no child command (namespace, verb, resource) or a positional arg.
// A leaf command can have 2 types of children: arguments or flags
func (node *AutoCompleteNode) isLeafCommand() bool {
	if node.Type != AutoCompleteNodeTypeCommand {
		return false
	}
	for _, child := range node.Children {
		if child.Type == AutoCompleteNodeTypeCommand {
			return false
		}
	}
	return true
}

// BuildAutoCompleteTree builds the autocomplete tree from the commands, subcommands and arguments
func BuildAutoCompleteTree(commands *Commands) *AutoCompleteNode {
	root := NewAutoCompleteCommandNode()
	root.addGlobalFlags()
	for _, cmd := range commands.commands {
		node := root

		// Creates nodes for namespaces, resources, verbs
		for _, part := range []string{cmd.Namespace, cmd.Resource, cmd.Verb} {
			if part != "" {
				node = node.GetChildOrCreate(part)
				node.addGlobalFlags()
			}
		}

		node.Command = cmd

		// We consider ArgSpecs as leaf in the autocomplete tree.
		nonDeprecatedArgs := cmd.ArgSpecs.GetDeprecated(false)
		for _, argSpec := range nonDeprecatedArgs {
			if argSpec.Positional {
				node.Children[positionalValueNodeID] = NewAutoCompleteArgNode(cmd, argSpec)
				continue
			}
			node.Children[argSpec.Name+"="] = NewAutoCompleteArgNode(cmd, argSpec)
		}

		if cmd.WaitFunc != nil {
			node.Children["-w"] = NewAutoCompleteFlagNode(node, &FlagSpec{
				Name: "-w",
			})
			node.Children["--wait"] = NewAutoCompleteFlagNode(node, &FlagSpec{
				Name: "--wait",
			})
		}
	}

	return root
}

// AutoComplete process a command line and returns autocompletion suggestions.
//
// command <flag name>=<flag value beginning><tab> gives no suggestion for now
// eg: scw test flower create name=p -o=jso
func AutoComplete(ctx context.Context, leftWords []string, wordToComplete string, rightWords []string) *AutocompleteResponse {
	commands := ExtractCommands(ctx)

	// Create AutoComplete Tree
	commandTreeRoot := BuildAutoCompleteTree(commands)

	// For each left word that is not a flag nor an argument, we try to go deeper in the autocomplete tree and store the current node in `node`.
	node := commandTreeRoot
	// nodeIndexInWords is the rightmost word index, before the cursor, that contains either a namespace or verb or resource or flag or flag value.
	// see test 'scw test flower delete f'
	nodeIndexInWords := 0

	// We remove command binary name from the left words.
	leftWords = leftWords[1:]

	for i, word := range leftWords {
		children, childrenExists := node.Children[word]
		if !childrenExists {
			children, childrenExists = node.Children[positionalValueNodeID]
		}

		switch {
		case !childrenExists && node.isLeafCommand():
			// word is probably an unknown argument
			// Just skip it

		case !childrenExists:
			// We did not find a child matching exactly the word
			// We did not reach a leaf command, and word is unknown
			return &AutocompleteResponse{}

		case children.Type == AutoCompleteNodeTypeArgument:
			// Do nothing
			// Arguments do not have children: they are not used to go deeper into the tree

		default:
			// word is a namespace or verb or resource or flag or flag value
			node = children
			nodeIndexInWords = i
		}
	}

	// keep track of the completed args
	completedArgs := make(map[string]struct{})

	// keep track of the completed flags
	completedFlags := make(map[string]struct{})

	// We loop through all other words in order to find existing args and flags.
	// When a flag is found it populates `completedFlags`.
	// When an argument is found it populates `completedArgs`.
	for i, word := range append(leftWords, rightWords...) {
		switch {
		// handle --flag=value and --flag
		case isFlag(word):
			completedFlags[wordKey(word)] = struct{}{}

		// handle arg=value
		case isArg(word):
			completedArgs[wordKey(word)+"="] = struct{}{}

		// handle boolean arg
		default:
			if _, exist := node.Children[positionalValueNodeID]; exist && i > nodeIndexInWords {
				completedArgs[word] = struct{}{}
			}
		}
	}

	if isCompletingArgValue(wordToComplete) {
		argName, argValuePrefix := splitArgWord(wordToComplete)
		argNode, exist := node.GetChildMatch(argName)
		if !exist {
			// We try to complete the value of an unknown arg
			return &AutocompleteResponse{}
		}
		suggestions := AutoCompleteArgValue(ctx, argNode.Command, argNode.ArgSpec, argValuePrefix)

		// We need to prefix suggestions with the argName to enable the arg value auto-completion.
		for k, s := range suggestions {
			suggestions[k] = argName + s
		}

		return newAutoCompleteResponse(suggestions)
	}

	// We are trying to complete a node: either a command name or an arg name or a flagname or a positional args
	suggestions := []string(nil)
	for key, child := range node.Children {
		if key == positionalValueNodeID {
			for _, positionalSuggestion := range AutoCompleteArgValue(ctx, child.Command, child.ArgSpec, wordToComplete) {
				if _, exists := completedArgs[positionalSuggestion]; !exists {
					suggestions = append(suggestions, positionalSuggestion)
				}
			}
			continue
		}

		if !hasPrefix(key, wordToComplete) {
			continue
		}

		// if no special keys in the key, we don't need to modify it
		if !strings.Contains(key, sliceSchema) && !strings.Contains(key, mapSchema) {
			if _, exists := completedArgs[key]; exists {
				continue
			}
			if _, exists := completedFlags[key]; exists {
				continue
			}
			if isFlag(key) && wordToComplete == "" {
				// skip autocomplete flag on empty string
				// command: scw <tab>
				// suggestions: instance
				// command: scw -<tab>
				// suggestions: -o
				continue
			}
			suggestions = append(suggestions, key)
			continue
		}

		// we know that we got either a slice, a map, or both
		suggestions = append(suggestions, keySuggestion(key, completedArgs, wordToComplete)...)
	}

	return newAutoCompleteResponse(suggestions)
}

// AutoCompleteArgValue returns suggestions for a (argument name, argument value prefix) pair.
// Priority is given to the AutoCompleteFunc from the ArgSpec, if it is set.
// Otherwise, we use EnumValues from the ArgSpec.
func AutoCompleteArgValue(ctx context.Context, cmd *Command, argSpec *ArgSpec, argValuePrefix string) []string {
	if argSpec == nil {
		return nil
	}
	if argSpec.AutoCompleteFunc != nil {
		return argSpec.AutoCompleteFunc(ctx, argValuePrefix)
	}

	possibleValues := []string(nil)

	if fieldType, err := args.GetArgType(cmd.ArgsType, argSpec.Name); err == nil {
		switch fieldType.Kind() {
		case reflect.Bool:
			possibleValues = []string{"true", "false"}
		}
	}

	if len(argSpec.EnumValues) > 0 {
		possibleValues = argSpec.EnumValues
	}

	suggestions := []string(nil)
	for _, value := range possibleValues {
		if strings.HasPrefix(value, argValuePrefix) {
			suggestions = append(suggestions, value)
		}
	}
	return suggestions
}

func isCompletingArgValue(wordToComplete string) bool {
	wordParts := strings.SplitN(wordToComplete, "=", 2)
	return len(wordParts) == 2
}

func splitArgWord(wordToComplete string) (string, string) {
	wordParts := strings.SplitN(wordToComplete, "=", 2)
	return wordParts[0] + "=", wordParts[1]
}

func wordKey(word string) string {
	return strings.SplitN(word, "=", 2)[0]
}

func isArg(wordToComplete string) bool {
	return strings.Contains(wordToComplete, "=")
}

func isFlag(wordToComplete string) bool {
	return strings.HasPrefix(wordToComplete, "-")
}

// hasPrefix will look if the word to complete prefixes the given key.
// It also handle complexe key schema such as slices and maps. E.g.:
// `security-gr` prefixes `security-group-id`
// `image-ids.0` prefixes `image-ids.{index}`
// `volumes.0.s` prefixes `volumes.{index}.size`
// `ip.fr-par.c` prefixes `ip.{key}.class`
func hasPrefix(key, wordToComplete string) bool {
	switch {
	case strings.HasPrefix(key, wordToComplete):
		return true
	case !strings.Contains(wordToComplete, ".") && (strings.HasPrefix(key, sliceSchema) || strings.HasPrefix(key, mapSchema)):
		return true
	case !strings.Contains(key, ".") || !strings.Contains(wordToComplete, "."):
		return false
	}

	tmp := strings.SplitN(key, ".", 2)
	leftKey, rightKey := tmp[0], tmp[1]
	tmp = strings.SplitN(wordToComplete, ".", 2)
	leftWord, rightWord := tmp[0], tmp[1]
	if leftKey == leftWord || leftKey == sliceSchema || leftKey == mapSchema {
		return hasPrefix(rightKey, rightWord)
	}
	return false
}

// keySuggestion will suggest the next key available for the map (or array) argument.
// Keys are suggested in ascending order arg.0, arg.1, arg.2...
func keySuggestion(key string, completedArg map[string]struct{}, wordToComplete string) []string {
	splitKey := strings.Split(key, ".")
	splitWordToComplete := strings.Split(wordToComplete, ".")

	// let's replace the existing placeholder with already typed values
	for i, k := range splitKey {
		if i >= len(splitWordToComplete) {
			continue
		}
		if k != splitWordToComplete[i] && (k == sliceSchema || k == mapSchema) && splitWordToComplete[i] != "" {
			splitKey[i] = splitWordToComplete[i]
		}
	}
	newKey := strings.Join(splitKey, ".")
	if !strings.Contains(newKey, sliceSchema) && !strings.Contains(newKey, mapSchema) {
		for arg := range completedArg {
			// if the arg is already given, ignore it
			if arg == newKey {
				return []string{}
			}
		}
		return []string{newKey}
	}

	usedIndex := make(map[string]struct{})

	depth := 0
	newKey = strings.ReplaceAll(newKey, sliceSchema, "([0-9]+)")
	newKey = strings.ReplaceAll(newKey, mapSchema, "[0-9a-zA-Z\\-]+")
	r := regexp.MustCompile(newKey)
	for arg := range completedArg {
		matches := r.FindStringSubmatch(arg)
		// the matches will all have the same length
		if len(matches) > 0 {
			depth = len(matches) - 1
			usedIndex[strings.Join(matches[1:], ".")] = struct{}{}
		}
	}
	newKey = strings.ReplaceAll(newKey, "[0-9a-zA-Z\\-]+", mapSchema)

	// let's cut the key on mapSchema
	keyCut := strings.Split(newKey, mapSchema)
	newKey = keyCut[0]

	if depth == 0 {
		return []string{strings.ReplaceAll(newKey, "([0-9]+)", "0")}
	}

	finalKeys := []string{}

	baseIndex := make([]string, depth)
	for i := range baseIndex {
		baseIndex[i] = "0"
	}
	for j := 1; j <= depth; j++ {
		for {
			_, exist := usedIndex[strings.Join(baseIndex, ".")]
			if !exist {
				key := newKey
				for _, v := range baseIndex {
					key = strings.Replace(key, "([0-9]+)", v, 1)
				}
				finalKeys = append(finalKeys, key)
				for k := 1; k <= j; k++ {
					baseIndex[len(baseIndex)-k] = "0"
				}
				break
			}
			if strings.HasSuffix(newKey, ".") {
				key := newKey
				for _, v := range baseIndex {
					key = strings.Replace(key, "([0-9]+)", v, 1)
				}
				finalKeys = append(finalKeys, key)
			}
			newIndex, _ := strconv.Atoi(baseIndex[len(baseIndex)-1])
			baseIndex[len(baseIndex)-1] = strconv.Itoa(newIndex + 1)
		}
		if j != depth {
			newIndex, _ := strconv.Atoi(baseIndex[len(baseIndex)-j-1])
			baseIndex[len(baseIndex)-j-1] = strconv.Itoa(newIndex + 1)
		}
	}
	return finalKeys
}
