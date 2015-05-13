package main

import (
	"fmt"
	"os"
	"runtime"
)

var cmdVersion = &Command{
	Exec:        runVersion,
	UsageLine:   "version",
	Description: "Show the version information",
	Help:        "Show the version information.",
}

func runVersion(cmd *Command, args []string) {
	if len(args) != 0 {
		fmt.Fprintf(os.Stderr, "usage: scw %s\n", cmd.UsageLine)
		os.Exit(1)
	}

	_, err := GetScalewayAPI()
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to init Scaleway API: %v\n", err)
		os.Exit(1)
	}

	// FIXME: fmt.Printf("Client version: %s\n", "FIXME")
	// FIXME: fmt.Printf("Client SDK version: %s\n", "FIXME")
	fmt.Printf("Go version (client): %s\n", runtime.Version())
	// FIXME: fmt.Printf("Git commit (client): %s\n", "FIXME")
	fmt.Printf("OS/Arch (client): %s/%s\n", runtime.GOOS, runtime.GOARCH)
	// FIXME: API version information
}
