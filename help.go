package main

var cmdHelp = &Command{
	Exec:        runHelp,
	UsageLine:   "help [command]",
	Description: "help of the scw command line",
	Help: `
Help prints help information about scw and its commands.

By default, help lists available commands with a short description.
When invoked with a command name, it prints the usage and the help of
the command.
`,
}

func runHelp(cmd *Command, args []string) {
}
