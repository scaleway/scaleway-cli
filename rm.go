package main

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
)

var cmdRm = &Command{
	Exec:        runRm,
	UsageLine:   "rm [OPTIONS] SERVER [SERVER...]",
	Description: "Remove one or more servers",
	Help:        "Remove one or more servers.",
}

func runRm(cmd *Command, args []string) {
	if len(args) == 0 {
		log.Fatalf("usage: scw %s", cmd.UsageLine)
	}
	has_error := false
	for _, needle := range args {
		server := cmd.GetServer(needle)
		err := cmd.API.DeleteServer(server)
		if err != nil {
			log.Errorf("failed to delete server %s: %s", server, err)
			has_error = true
		} else {
			fmt.Println(needle)
		}
	}
	if has_error {
		os.Exit(1)
	}
}
