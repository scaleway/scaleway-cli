package main

import (
	"fmt"
	"os"
	"time"
)

var cmdWait = &Command{
	Exec:        runWait,
	UsageLine:   "wait [OPTIONS] SERVER [SERVER...]",
	Description: "Wait until a running server stops.",
	Help:        "Wait until a running server stops.",
}

func runWait(cmd *Command, args []string) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "usage: scw %s\n", cmd.UsageLine)
		os.Exit(1)
	}
	has_error := false
	for _, needle := range args {
		server_identifier := cmd.GetServer(needle)
		for {
			server, err := cmd.API.GetServer(server_identifier)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to retrieve information from server %s: %s\n", server_identifier, err)
				has_error = true
				break
			}
			if server.State == "stopped" {
				break
			}
			time.Sleep(1)
		}
	}
	if has_error {
		os.Exit(1)
	}
}
