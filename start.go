package main

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
)

var cmdStart = &Command{
	Exec:        runStart,
	UsageLine:   "start [OPTIONS] SERVER [SERVER...]",
	Description: "Start a stopped server",
	Help:        "Start a stopped server.",
}

func init() {
	// FIXME: -h
	cmdStart.Flag.BoolVar(&startW, []string{"w", "-wait"}, false, "Synchronous start. Wait for SSH to be ready")
}

// Flags
var startW bool // -w flag

func runStart(cmd *Command, args []string) {
	if len(args) == 0 {
		log.Fatalf("usage: scw %s", cmd.UsageLine)
	}
	has_error := false
	for _, needle := range args {
		server := cmd.GetServer(needle)
		err := cmd.API.PostServerAction(server, "poweron")
		if err != nil {
			if err.Error() != "server should be stopped" {
				log.Errorf("failed to stop server %s: %s", server, err)
				has_error = true
			}
		} else {
			if startW {
				_, err = WaitForServerReady(cmd.API, server)
				if err != nil {
					log.Errorf("Failed to wait for server %s to be ready, %v", needle, err)
					has_error = true
				}
			}

			fmt.Println(needle)
		}
		if has_error {
			os.Exit(1)
		}
	}
}
