// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package cli

// Command is a Scaleway command
import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/scaleway/scaleway-cli/vendor/github.com/Sirupsen/logrus"
	flag "github.com/scaleway/scaleway-cli/vendor/github.com/docker/docker/pkg/mflag"

	"github.com/scaleway/scaleway-cli/pkg/api"
	"github.com/scaleway/scaleway-cli/pkg/commands"
)

// Command contains everything needed by the cli main loop to calls the workflow, display help and usage, and the context
type Command struct {
	// Exec executes the command
	Exec func(cmd *Command, args []string)

	// Usage is the one-line usage message.
	UsageLine string

	// Description is the description of the command
	Description string

	// Help is the full description of the command
	Help string

	// Examples are some examples of the command
	Examples string

	// Flag is a set of flags specific to this command.
	Flag flag.FlagSet

	// Hidden is a flat to hide command from global help commands listing
	Hidden bool

	// API is the interface used to communicate with Scaleway's API
	API *api.ScalewayAPI
}

// GetContext returns a standard context, with real stdin, stdout, stderr, a configured API and raw arguments
func (c *Command) GetContext(rawArgs []string) commands.CommandContext {
	return commands.CommandContext{
		Stdin:   os.Stdin,
		Stdout:  os.Stdout,
		Stderr:  os.Stderr,
		Env:     os.Environ(),
		RawArgs: rawArgs,
		API:     c.API,
	}
}

// Name returns the command's name
func (c *Command) Name() string {
	name := c.UsageLine
	i := strings.Index(name, " ")
	if i >= 0 {
		name = name[:i]
	}
	return name
}

// PrintUsage prints a full command usage and exits
func (c *Command) PrintUsage() {
	helpMessage, err := commandHelpMessage(c)
	if err != nil {
		logrus.Fatalf("%v", err)
	}
	fmt.Fprintf(os.Stderr, "%s\n", helpMessage)
	os.Exit(1)
}

// PrintShortUsage prints a short command usage and exits
func (c *Command) PrintShortUsage() {
	fmt.Fprintf(os.Stderr, "usage: scw %s. See 'scw %s --help'.\n", c.UsageLine, c.Name())
	os.Exit(1)
}

// Options returns a string describing options of the command
func (c *Command) Options() string {
	var options string
	visitor := func(flag *flag.Flag) {
		name := strings.Join(flag.Names, ", -")
		var optionUsage string
		if flag.DefValue == "" {
			optionUsage = fmt.Sprintf("%s=\"\"", name)
		} else {
			optionUsage = fmt.Sprintf("%s=%s", name, flag.DefValue)
		}
		options += fmt.Sprintf("  -%-20s %s\n", optionUsage, flag.Usage)
	}
	c.Flag.VisitAll(visitor)
	if len(options) == 0 {
		return ""
	}
	return fmt.Sprintf("Options:\n\n%s", options)
}

// ExamplesHelp returns a string describing examples of the command
func (c *Command) ExamplesHelp() string {
	if c.Examples == "" {
		return ""
	}
	return fmt.Sprintf("Examples:\n\n%s", strings.Trim(c.Examples, "\n"))
}

var fullHelpTemplate = `
Usage: scw {{.UsageLine}}

{{.Help}}

{{.Options}}
{{.ExamplesHelp}}
`

func commandHelpMessage(cmd *Command) (string, error) {
	t := template.New("full")
	template.Must(t.Parse(fullHelpTemplate))
	var output bytes.Buffer
	err := t.Execute(&output, cmd)
	if err != nil {
		return "", err
	}
	return strings.Trim(output.String(), "\n"), nil
}
