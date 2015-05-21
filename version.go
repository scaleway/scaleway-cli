package main

import (
	"fmt"
	"runtime"

	log "github.com/Sirupsen/logrus"
	"github.com/scaleway/scaleway-cli/scwversion"
)

var cmdVersion = &Command{
	Exec:        runVersion,
	UsageLine:   "version",
	Description: "Show the version information",
	Help:        "Show the version information.",
}

func runVersion(cmd *Command, args []string) {
	if len(args) != 0 {
		log.Fatalf("usage: scw %s", cmd.UsageLine)
	}

	fmt.Printf("Client version: %s\n", scwversion.VERSION)
	fmt.Printf("Go version (client): %s\n", runtime.Version())
	fmt.Printf("Git commit (client): %s\n", scwversion.GITCOMMIT)
	fmt.Printf("OS/Arch (client): %s/%s\n", runtime.GOOS, runtime.GOARCH)
	// FIXME: API version information
}
