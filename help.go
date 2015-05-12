package main

import (
	"fmt"
	"os"
	"text/template"
)

var cmdHelp = &Command{
	Exec:        nil,
	UsageLine:   "help [COMMAND]",
	Description: "help of the scw command line",
	Help: `
Help prints help information about scw and its commands.

By default, help lists available commands with a short description.
When invoked with a command name, it prints the usage and the help of
the command.
`,
}

func init() {
	// break dependency loop
	cmdHelp.Exec = runHelp
}

var helpTemplate = `Scw is a tool to interact with Scaleway from the command line.

Usage:

	scw command [arguments]

The commands are:

{{range .}} {{.Name | printf "%-12s"}} {{.Description}}
{{end}}
`

var fullHelpTemplate = `{{.Description}}

Usage:

        scw {{.UsageLine}}

{{.Help}}

Options:

{{.Options}}
`

func runHelp(cmd *Command, args []string) {
	if len(args) >= 1 {
		name := args[0]
		for _, cmd := range commands {
			if cmd.Name() == name {
				t := template.New("full")
				template.Must(t.Parse(fullHelpTemplate))
				if err := t.Execute(os.Stdout, cmd); err != nil {
					panic(err)
				}
				return
			}
		}
		fmt.Fprintf(os.Stderr, "Unknown help topic `%s`.  Run 'scw help'.\n", name)
		os.Exit(1)
	} else {
		t := template.New("top")
		template.Must(t.Parse(helpTemplate))
		if err := t.Execute(os.Stdout, commands); err != nil {
			panic(err)
		}
	}
}
