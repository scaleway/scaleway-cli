package main

import (
	"fmt"
	"os"
)

var cmdRm = &Command{
	Exec:        runRm,
	UsageLine:   "rm [OPTIONS] SERVER [SERVER...]",
	Description: "Remove one or more servers",
	Help:        "Remove one or more servers.",
}

func runRm(cmd *Command, args []string) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "usage: scw %s\n", cmd.UsageLine)
		os.Exit(1)
	}
	has_error := false
	for _, needle := range args {
		server := cmd.GetServer(needle)
		err := cmd.API.DeleteServer(server)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to delete server %s: %s\n", server, err)
			has_error = true
		} else {
			fmt.Fprintf(os.Stdout, "%s\n", needle)
		}
	}
	if has_error {
		os.Exit(1)
	}
}
