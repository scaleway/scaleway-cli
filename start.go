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
}
