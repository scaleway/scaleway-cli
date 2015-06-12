package commands

import (
	"log"
	"os"
	"text/template"

	types "github.com/scaleway/scaleway-cli/commands/types"
)

// CmdHelp is the 'scw help' command
var CmdHelp = &types.Command{
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
	CmdHelp.Exec = runHelp
	CmdHelp.Flag.BoolVar(&helpHelp, []string{"h", "-help"}, false, "Print usage")
}

// Flags
var helpHelp bool // -h, --help flag

var helpTemplate = `Usage: scw [OPTIONS] COMMAND [arg...]

Interact with Scaleway from the command line.

Options:
 --api-endpoint=APIEndPoint   Set the API endpoint
 -D, --debug=false            Enable debug mode
 -h, --help=false             Print usage
 -v, --version=false          Print version information and quit

Commands:
{{range .}}{{if not .Hidden}}    {{.Name | printf "%-9s"}} {{.Description}}
{{end}}{{end}}
Run 'scw COMMAND --help' for more information on a command.
`

func runHelp(cmd *types.Command, args []string) {
	if waitHelp {
		cmd.PrintUsage()
	}
	if len(args) > 1 {
		cmd.PrintShortUsage()
	}

	if len(args) == 1 {
		name := args[0]
		for _, command := range Commands {
			if command.Name() == name {
				command.PrintUsage()
			}
		}
		log.Fatalf("Unknown help topic `%s`.  Run 'scw help'.", name)
	} else {
		t := template.New("top")
		template.Must(t.Parse(helpTemplate))
		if err := t.Execute(os.Stdout, Commands); err != nil {
			panic(err)
		}
	}
}
