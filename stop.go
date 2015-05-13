package main

import (
	"fmt"
	"os"
)

var cmdStop = &Command{
	Exec:        runStop,
	UsageLine:   "stop [OPTIONS] SERVER [SERVER...]",
	Description: "Stop a running server",
	Help:        "Stop a running server.",
}

func runStop(cmd *Command, args []string) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "usage: scw %s\n", cmd.UsageLine)
		os.Exit(1)
	}
	has_error := false
	for _, needle := range args {
		server := cmd.GetServer(needle)
		err := cmd.API.PostServerAction(server, "poweroff")
		if err != nil {
			if err.Error() != "server should be running" {
				fmt.Fprintf(os.Stderr, "failed to stop server %s: %s\n", server, err)
				has_error = true
			}
		} else {
			fmt.Fprintf(os.Stdout, "%s\n", needle)
		}
		if has_error {
			os.Exit(1)
		}
	}
}
