package main

import log "github.com/Sirupsen/logrus"

var cmdRename = &Command{
	Exec:        runRename,
	UsageLine:   "rename [OPTIONS] SERVER NEW_NAME",
	Description: "Rename a server",
	Help:        "Rename a server.",
}

func runRename(cmd *Command, args []string) {
	if len(args) < 2 {
		log.Fatalf("usage: scw %s", cmd.UsageLine)
	}

	serverId := cmd.GetServer(args[0])

	var server ScalewayServerPathNameDefinition
	server.Name = args[1]

	err := cmd.API.PatchServerName(serverId, server)
	if err != nil {
		log.Fatalf("Cannot rename server: %v", err)
	}
}
