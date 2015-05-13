package main

import (
	"fmt"
	"os"
)

var cmdStart = &Command{
	Exec:        runStart,
	UsageLine:   "start [OPTIONS] SERVER [SERVER...]",
	Description: "Start a stopped server.",
	Help:        "Start a stopped server.",
}

func runStart(cmd *Command, args []string) {
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "usage: scw %s\n", cmd.UsageLine)
		os.Exit(1)
	}
	needle := args[0]
	server := cmd.GetServer(needle)
	err := cmd.API.PostServerAction(server, "poweron")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to start server %s: %s\n", server, err)
		os.Exit(1)
	}
}
