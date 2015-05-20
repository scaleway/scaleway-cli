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
			fmt.Println(needle)
		}
		if has_error {
			os.Exit(1)
		}
	}
}
