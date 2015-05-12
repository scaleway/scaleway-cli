// scw interract with Scaleway from the command line
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// Command is a Scaleway command
type Command struct {
	// Exec executes the command
	Exec func(cmd *Command, args []string)

	// Usage is the one-line usage message.
	UsageLine string

	// Description is the description of the command
	Description string

	// Help is the full description of the command
	Help string

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

// Options returns a string describing options of the command
func (c *Command) Options() string {
	var options string
	visitor := func(flag *flag.Flag) {
		options += fmt.Sprintf("  -%-12s %s (%s)\n", flag.Name, flag.Usage, flag.DefValue)
	}
	c.Flag.VisitAll(visitor)
	if len(options) == 0 {
		options = "  no option for this command"
	}
	return options
}

var commands = []*Command{
	cmdHelp,
	cmdLogin,
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		usage()
	}
	name := args[0]
	args = args[1:]

	for _, cmd := range commands {
		if cmd.Name() == name {
			cmd.Flag.SetOutput(ioutil.Discard)
			err := cmd.Flag.Parse(args)
			if err != nil {
				fmt.Fprintf(os.Stderr, "usage: scw %s\n", cmd.UsageLine)
				os.Exit(1)
			}
			cmd.Exec(cmd, cmd.Flag.Args())
			os.Exit(0)
		}
	}

	fmt.Fprintf(os.Stderr, "scw: unknown subcommand %s\nRun 'scw help for usage.\n", name)
	os.Exit(1)
}

func usage() {
	cmdHelp.Exec(cmdHelp, []string{})
	os.Exit(1)
}
