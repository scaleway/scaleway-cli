package main

import (
	"fmt"
	"os"
)

var cmdStart = &Command{
	Exec:        runStart,
	UsageLine:   "start [options] SERVER",
	Description: "Start a Scaleway C1 server",
	Help: `
Start starts a Scaleway C1 server.`,
}

func runStart(cmd *Command, args []string) {
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "usage: scw %s\n", cmd.UsageLine)
		os.Exit(1)
	}

	_, err := GetScalewayAPI()
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to init Scaleway API: %v\n", err)
		os.Exit(1)
	}
}
