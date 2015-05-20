package main

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
)

var cmdCommit = &Command{
	Exec:        runCommit,
	UsageLine:   "commit [OPTIONS] SERVER [NAME]",
	Description: "Create a new snapshot from a server's volume",
	Help:        "Create a new snapshot from a server's volume.",
}

func runCommit(cmd *Command, args []string) {
	if len(args) == 0 {
		log.Fatalf("usage: scw %s", cmd.UsageLine)
	}

	serverId := cmd.GetServer(args[0])
	server, err := cmd.API.GetServer(serverId)
	if err != nil {
		log.Fatalf("Cannot fetch server: %v", err)
	}
	var volume ScalewayVolume = server.Volumes["0"]
	var name string
	if len(args) > 1 {
		name = args[1]
	} else {
		name = volume.Name + "-snapshot"
	}
	snapshot, err := cmd.API.PostSnapshot(volume.Identifier, name)
	if err != nil {
		log.Fatalf("Cannot create snapshot: %v", err)
	}
	fmt.Println(snapshot)
}
