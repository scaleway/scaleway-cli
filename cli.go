// scw interract with Scaleway from the command line
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

// Command is a Scaleway command
type Command struct {
	// Exec executes the command
	Exec func(args []string)

	// Usage is the one-line usage message.
	UsageLine string

	// Description is the description of the command
	Description string

	// Flag is a set of flags specific to this command.
	Flag flag.FlagSet

	// CustomFlags indicates that the command will do its own
	// flag parsing.
	CustomFlags bool
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

// Usage prints a usage on stderr then exits
func (c *Command) Usage() {
	fmt.Fprintf(os.Stderr, "usage: %s\n\n", c.Usage)
	fmt.Fprintf(os.Stderr, "%s\n", c.Description)
	os.Exit(1)
}

var usageTemplate = `Scw is a tool to interract with Scaleway from the command line.

Usage:

	scw command [arguments]

`

// usage prints the usage then exits
func usage() {
	fmt.Fprintf(os.Stderr, usageTemplate)
	os.Exit(1)
}

var commands = []*Command{}

func main() {
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		usage()
	}

	name := args[0]
	for _, cmd := range commands {
		if cmd.Name() == name {
			args := args[1:]
			cmd.Exec(args)
			os.Exit(0)
		}
	}

	fmt.Fprintf(os.Stderr, "scw: unknown subcommand %s\nRun 'scw help for usage.\n", name)
	os.Exit(1)
}
