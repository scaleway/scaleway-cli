package main

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
)

var cmdRestart = &Command{
	Exec:        runRestart,
	UsageLine:   "restart [OPTIONS] SERVER [SERVER...]",
	Description: "Restart a running server",
	Help:        "Restart a running server.",
}

func runRestart(cmd *Command, args []string) {
	if len(args) == 0 {
		log.Fatalf("usage: scw %s", cmd.UsageLine)
	}
	has_error := false
	for _, needle := range args {
		server := cmd.GetServer(needle)
		err := cmd.API.PostServerAction(server, "reboot")
		if err != nil {
			if err.Error() != "server is being stopped or rebooted" {
				log.Errorf("failed to restart server %s: %s", server, err)
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
