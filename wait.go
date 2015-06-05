package main

import (
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
)

var cmdWait = &Command{
	Exec:        runWait,
	UsageLine:   "wait [OPTIONS] SERVER [SERVER...]",
	Description: "Block until a server stops",
	Help:        "Block until a server stops.",
}

func init() {
	cmdWait.Flag.BoolVar(&waitHelp, []string{"h", "-help"}, false, "Print usage")
}

// Flags
var waitHelp bool // -h, --help flag

func runWait(cmd *Command, args []string) {
	if waitHelp {
		cmd.PrintUsage()
	}
	if len(args) < 1 {
		cmd.PrintShortUsage()
	}

	has_error := false
	for _, needle := range args {
		server_identifier := cmd.GetServer(needle)
		for {
			server, err := cmd.API.GetServer(server_identifier)
			if err != nil {
				log.Errorf("failed to retrieve information from server %s: %s", server_identifier, err)
				has_error = true
				break
			}
			if server.State == "stopped" {
				break
			}
			time.Sleep(1 * time.Second)
		}
	}
	if has_error {
		os.Exit(1)
	}
}
