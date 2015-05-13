package main

import (
	"fmt"
	"os"
)

var cmdRestart = &Command{
	Exec:        runRestart,
	UsageLine:   "restart [OPTIONS] SERVER [SERVER...]",
	Description: "Restart a running server",
	Help:        "Restart a running server.",
}

func runRestart(cmd *Command, args []string) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "usage: scw %s\n", cmd.UsageLine)
		os.Exit(1)
	}
	has_error := false
	for _, needle := range args {
		server := cmd.GetServer(needle)
		err := cmd.API.PostServerAction(server, "reboot")
		if err != nil {
			if err.Error() != "server is being stopped or rebooted" {
				fmt.Fprintf(os.Stderr, "failed to restart server %s: %s\n", server, err)
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
