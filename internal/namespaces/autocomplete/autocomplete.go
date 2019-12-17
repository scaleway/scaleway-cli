package namespace_autocomplete

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"

	"github.com/scaleway/scaleway-cli/internal/args"
	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-cli/internal/interactive"
)

func GetCommands() *core.Commands {
	return core.NewCommands(
		autocompleteInstallCommand(),
		autocompleteCompleteBashCommand(),
		autocompleteCompleteFishCommand(),
		autocompleteCompleteZshCommand(),
		autocompleteScriptCommand(),
	)
}

type autocompleteScript struct {
	CompleteScript    string
	CompleteFunc      string
	FileName          string
	DestinationFolder map[string]string
	SuccessHelp       string
}

const (
	readWritePermission = 0644
	executePermission   = 0755
)

// autocompleteScripts regroups the autocomplete scripts for the different shells
// The key is the path of the shell.
var autocompleteScripts = map[string]autocompleteScript{
	"bash": {
		// If `scw` is the first word on the command line,
		// after hitting [tab] arguments are sent to `scw autocomplete complete bash`:
		//  - COMP_LINE :  the complete command line
		//  - COMP_CWORD : the index of the word being completed
		//  - COMP_WORDS : the words composing the command line
		//
		// Note that `=` signs are excluding from $COMP_WORDBREAKS. As a result, they are NOT be
		// considered as breaking words and arguments like `image=` will not be split.
		//
		// Then `scw autocomplete complete bash` process the line, and tries to returns suggestions.
		// These scw suggestions are put into `COMPREPLY` which is used by Bash to provides the shell suggestions.
		CompleteFunc: `
			local COMP_WORDBREAKS=${COMP_WORDBREAKS/=/}
			_scw() {
				output=$(scw autocomplete complete bash "$COMP_LINE" "$COMP_CWORD" "${COMP_WORDS[@]}")
				COMPREPLY=($output)
				[[ $COMPREPLY == *= ]] && compopt -o nospace
				return
			}
			complete -F _scw scw
		`,
		CompleteScript: `eval "$(scw autocomplete script shell=bash)"`,
		FileName:       "scw-completion.bash",
		DestinationFolder: map[string]string{
			"darwin": "/usr/local/etc/bash_completion.d/",
		},
		SuccessHelp: "You might want to run 'source ~/.bash_profile'.",
	},
	"fish": {
		// (commandline)                             complete command line
		// (commandline --cursor)                    position of the cursor, as number of chars in the command line
		// (commandline --current-token)             word to complete
		// (commandline --tokenize --cut-at-cursor)  tokenized selection up until the current cursor position
		//                                           formatted as one string-type token per line
		//
		// If files are shown although --no-files is set,
		// it might be because you are using an alias for scw, such as :
		// 		alias scw='go run "$HOME"/scaleway-cli/cmd/scw/main.go'
		// You might want to run 'complete --erase --command go' during development.
		//
		// TODO: send rightWords
		CompleteFunc: `
			complete --erase --command scw;
			complete --command scw --no-files;
			complete --command scw --arguments '(scw autocomplete complete fish (commandline) (commandline --cursor) (commandline --current-token) (commandline --tokenize --cut-at-cursor))';
		`,
		CompleteScript: `eval (scw autocomplete script shell=fish)`,
		FileName:       "scw-completion.fish",
		DestinationFolder: map[string]string{
			"darwin": path.Join(os.Getenv("HOME"), "/.config/fish/completions/"),
		},
		SuccessHelp: "You might want to add 'source ~/.config/fish/completions/scw-completion.fish' to ~/.config/fish/config.fish.",
	},
	"zsh": {
		// If you are using an alias for scw, such as :
		// 		alias scw='go run "$HOME"/scaleway-cli/cmd/scw/main.go'
		// you might want to run 'compdef _scw go' during development.
		CompleteFunc: `
			_scw () {
				# splits $BUFFER, i.e. the complete command line,
				# into shell words using shell parsing rules by Expansion Flag (z) and puts it into an array
				words=("${(z)BUFFER}")

				# If the last char of the line is a space, a last empty word is not added to words.
				# We need to add it manually.
				lastChar="${BUFFER: -1}"
				if [[ $lastChar = *[!\ ]* ]]; then # if $lastChar contains something else than spaces
					: # do nothing
				else
					# words+=('') does not work
					# couldn't find a way to add an empty string to an array
					# we replace 'EMPTY_WORD' by '' later in go code
					words+=('EMPTY_WORD')
				fi
				output=($(scw autocomplete complete zsh $CURSOR $words))
				opts=('-S' ' ')
				if [[ $output == *= ]]; then
					opts=('-S' '')
				fi
				compadd "${opts[@]}" -- "${output[@]}"
			}
			compdef _scw scw
		`,
		CompleteScript: `eval "$(scw autocomplete script shell=zsh)"`,
		FileName:       "_scw",
		DestinationFolder: map[string]string{
			"darwin": path.Join(os.Getenv("HOME"), "/.zsh/completion/"),
		},
		SuccessHelp: "Make sure the completion directory is in your $fpath by adding in ~/.zshrc:\n" +
			"  fpath=(~/.zsh/completion $fpath)\n" +
			"Make sure compinit is loaded or do it by adding in ~/.zshrc:\n" +
			"  autoload -U compinit && compinit",
	},
}

type autocompleteInstallArgs struct {
	Shell string
}

func autocompleteInstallCommand() *core.Command {
	return &core.Command{
		Short:     `Install autocompletion script`,
		Long:      `Install autocompletion script for a given shell and OS.`,
		Namespace: "autocomplete",
		Resource:  "install",
		NoClient:  true,
		ArgSpecs: core.ArgSpecs{
			{
				Name:    "shell",
				Default: core.DefaultValueSetter(os.Getenv("SHELL")),
			},
		},
		ArgsType: reflect.TypeOf(autocompleteInstallArgs{}),
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
			shellName := filepath.Base(argsI.(*autocompleteInstallArgs).Shell)

			// Detect shell and OS
			script, exists := autocompleteScripts[shellName]
			if !exists {
				return nil, unsupportedShellError(shellName)
			}
			destinationFolder, exists := script.DestinationFolder[runtime.GOOS]
			if !exists {
				return nil, unsupportedOsError(runtime.GOOS)
			}

			// Compute destinationPath
			destinationPath := path.Join(destinationFolder, script.FileName)

			// Inform user and ask for confirmation
			_, err := interactive.Printf(trimText(fmt.Sprintf(`
				To enable autocomplete, scw needs to add a script to your file system.
				Shell: %v
				Destination path: %v
				Content: %v`,
				shellName, destinationPath, script.CompleteScript)) + "\n\n")
			if err != nil {
				return nil, err
			}
			continueInstallation, err := interactive.PromptBoolWithConfig(&interactive.PromptBoolConfig{
				Prompt:       fmt.Sprintf("Do you want to install autocomplete for %v?", shellName),
				DefaultValue: true,
			})
			if err != nil {
				return nil, err
			}

			// Early exit if user disagrees
			if !continueInstallation {
				return nil, installationCancelledError(shellName, script.CompleteScript)
			}

			// Create destination folder if it doesn't exist
			os.MkdirAll(destinationFolder, executePermission)

			// Write script
			err = ioutil.WriteFile(destinationPath, []byte(trimText(script.CompleteScript)), readWritePermission)
			if err != nil {
				return nil, err
			}

			return &core.SuccessResult{
				Message: fmt.Sprintf("Autocomplete function for %v installed succesfully.\nCopied %v to %v.\n%v",
					shellName, script.FileName, destinationPath, script.SuccessHelp),
			}, nil
		},
	}
}

func autocompleteCompleteBashCommand() *core.Command {
	return &core.Command{
		Short:     `Autocomplete for Bash`,
		Long:      `Autocomplete for Bash.`,
		Namespace: "autocomplete",
		Resource:  "complete",
		Verb:      "bash",
		// TODO: Switch NoClient to true when cache will be implemented.
		NoClient: false,
		Hidden:   true,
		ArgsType: reflect.TypeOf(args.RawArgs{}),
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
			rawArgs := *argsI.(*args.RawArgs)
			wordIndex, err := strconv.Atoi(rawArgs[1])
			if err != nil {
				return nil, err
			}
			words := rawArgs[2:]
			leftWords := words[:wordIndex]
			wordToComplete := words[wordIndex]
			rightWords := words[wordIndex+1:]

			res := core.AutoComplete(ctx, leftWords, wordToComplete, rightWords)
			return strings.Join(res.Propositions, " "), nil
		},
	}
}

func autocompleteCompleteFishCommand() *core.Command {
	return &core.Command{
		Short:     `Autocomplete for Fish`,
		Long:      `Autocomplete for Fish.`,
		Namespace: "autocomplete",
		Resource:  "complete",
		Verb:      "fish",
		// TODO: Switch NoClient to true when cache will be implemented.
		NoClient: false,
		Hidden:   true,
		ArgsType: reflect.TypeOf(args.RawArgs{}),
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
			rawArgs := *argsI.(*args.RawArgs)
			leftWords := rawArgs[3:]
			wordToComplete := rawArgs[2]

			// TODO: compute rightWords once used by core.AutoComplete()
			// line := rawArgs[0]
			// charIndex, _ := strconv.Atoi(rawArgs[1])
			rightWords := []string(nil)

			res := core.AutoComplete(ctx, leftWords, wordToComplete, rightWords)

			// TODO: decide if we want to add descriptions
			// see https://stackoverflow.com/a/20879411
			// "followed optionally by a tab and a short description."
			return strings.Join(res.Propositions, "\n"), nil
		},
	}
}

func autocompleteCompleteZshCommand() *core.Command {
	return &core.Command{
		Short:     `Autocomplete for Zsh`,
		Long:      `Autocomplete for Zsh.`,
		Namespace: "autocomplete",
		Resource:  "complete",
		Verb:      "zsh",
		// TODO: Switch NoClient to true when cache will be implemented.
		NoClient: false,
		Hidden:   true,
		ArgsType: reflect.TypeOf(args.RawArgs{}),
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
			rawArgs := *argsI.(*args.RawArgs)

			words := rawArgs[1:]
			charIndex, err := strconv.Atoi(rawArgs[0])
			if err != nil {
				return nil, err
			}
			wordIndex := core.WordIndex(charIndex, words)
			leftWords := words[:wordIndex]
			wordToComplete := words[wordIndex]

			// In zsh, couldn't find a way to add an empty string to an array.
			// We added "EMPTY_WORD" instead.
			// "EMPTY_WORD" is replaced by "".
			// see the zsh script, line 106:
			//     words+=('EMPTY_WORD')
			if wordToComplete == "EMPTY_WORD" {
				wordToComplete = ""
			}

			// TODO: compute rightWords once used by core.AutoComplete()
			rightWords := []string(nil)

			res := core.AutoComplete(ctx, leftWords, wordToComplete, rightWords)
			return strings.Join(res.Propositions, " "), nil
		},
	}
}

type autocompleteShowArgs struct {
	Shell string
}

func autocompleteScriptCommand() *core.Command {
	return &core.Command{
		Short:     `Show autocomplete script for current shell`,
		Long:      `Show autocomplete script for current shell.`,
		Namespace: "autocomplete",
		Resource:  "script",
		NoClient:  true,
		ArgSpecs: core.ArgSpecs{
			{
				Name:    "shell",
				Default: core.DefaultValueSetter(os.Getenv("SHELL")),
			},
		},
		ArgsType: reflect.TypeOf(autocompleteShowArgs{}),
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
			shell := filepath.Base(argsI.(*autocompleteShowArgs).Shell)
			script, exists := autocompleteScripts[shell]
			if !exists {
				return nil, unsupportedShellError(shell)
			}
			return trimText(script.CompleteFunc), nil
		},
	}
}

func trimText(str string) string {
	foundFirstNonEmptyLine := false
	strToRemove := ""
	lines := strings.Split(str, "\n")
	for i, line := range lines {
		if !foundFirstNonEmptyLine {
			if len(line) > 0 {
				for _, c := range line {
					if c == ' ' || c == '\t' {
						strToRemove += string(c)
						continue
					}
					break
				}
				foundFirstNonEmptyLine = true
			}
		}
		for _, c := range strToRemove {
			lines[i] = strings.Replace(lines[i], string(c), "", 1)
		}
	}
	lines = removeStartingAndEndingEmptyLines(lines)
	return strings.Join(lines, "\n")
}

func removeStartingAndEndingEmptyLines(lines []string) []string {
	lines = removeStartingEmptyLines(lines)
	lines = reverseLines(lines)
	lines = removeStartingEmptyLines(lines)
	lines = reverseLines(lines)
	return lines
}

func removeStartingEmptyLines(lines []string) []string {
	doAdd := false
	lines2 := []string(nil)
	for _, line := range lines {
		if len(line) > 0 {
			doAdd = true
		}
		if doAdd {
			lines2 = append(lines2, line)
		}
	}
	return lines2
}

func reverseLines(lines []string) []string {
	for i, j := 0, len(lines)-1; i < j; i, j = i+1, j-1 {
		lines[i], lines[j] = lines[j], lines[i]
	}
	return lines
}
