package main

import (
	"os"
	"text/template"

	log "github.com/Sirupsen/logrus"
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

var helpTemplate = `Usage: scw [OPTIONS] COMMAND [arg...]

Interact with Scaleway from the command line.

Options:
 --api-endpoint=APIEndPoint   Set the API endpoint
 -D, --debug=false            Enable debug mode
 -h, --help=false             Print usage
 -v, --version=false          Print version information and quit

Commands:
{{range .}}    {{.Name | printf "%-9s"}} {{.Description}}
{{end}}
Run 'scw COMMAND --help' for more information on a command.
`

var fullHelpTemplate = `
Usage: scw {{.UsageLine}}

{{.Help}}

{{.Options}}
`

func runHelp(cmd *Command, args []string) {
	if len(args) >= 1 {
		name := args[0]
		for _, cmd := range commands {
			if cmd.Name() == name {
				t := template.New("full")
				template.Must(t.Parse(fullHelpTemplate))
				// FIXME: TrimRight
				if err := t.Execute(os.Stdout, cmd); err != nil {
					panic(err)
				}
				return
			}
		}
		log.Fatalf("Unknown help topic `%s`.  Run 'scw help'.", name)
	} else {
		t := template.New("top")
		template.Must(t.Parse(helpTemplate))
		if err := t.Execute(os.Stdout, commands); err != nil {
			panic(err)
		}
	}
}
