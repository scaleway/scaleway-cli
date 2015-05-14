package main

import (
	"fmt"
	"os"
)

var cmdInspect = &Command{
	Exec:        runInspect,
	UsageLine:   "inspect [OPTIONS] IDENTIFIER [IDENTIFIER...]",
	Description: "Inspects a server, image or a bootscript.",
	Help:        "Inspects a server, image or a bootscript.",
}

func runInspect(cmd *Command, args []string) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "usage: scw %s\n", cmd.UsageLine)
		os.Exit(1)
	}
	has_error := false
	for _, _ = range args {
	}
	if has_error {
		os.Exit(1)
	}
}
