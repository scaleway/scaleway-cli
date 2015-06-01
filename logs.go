package main

import log "github.com/Sirupsen/logrus"

var cmdLogs = &Command{
	Exec:        runLogs,
	UsageLine:   "logs [OPTIONS] SERVER",
	Description: "Fetch the logs of a server",
	Help:        "Fetch the logs of a server.",
}

func runLogs(cmd *Command, args []string) {
	if len(args) != 1 {
		log.Fatalf("usage: scw %s", cmd.UsageLine)
	}
	serverId := cmd.GetServer(args[0])
	server, err := cmd.API.GetServer(serverId)
	if err != nil {
		log.Fatalf("Failed to get server information for %s: %v", serverId, err)
	}

	// FIXME: switch to serial history when API is ready

	command := []string{"dmesg"}
	err = serverExec(server.PublicAddress.IP, command)
	if err != nil {
		log.Fatalf("Command execution failed: %v", err)
	}
}
