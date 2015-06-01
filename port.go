package main

import log "github.com/Sirupsen/logrus"

var cmdPort = &Command{
	Exec:        runPort,
	UsageLine:   "port [OPTIONS] SERVER [PRIVATE_PORT[/PROTO]]",
	Description: "Lookup the public-facing port that is NAT-ed to PRIVATE_PORT",
	Help:        "List port mappings for the SERVER, or lookup the public-facing port that is NAT-ed to the PRIVATE_PORT",
}

func runPort(cmd *Command, args []string) {
	if len(args) < 1 {
		log.Fatalf("usage: scw %s", cmd.UsageLine)
	}
	serverId := cmd.GetServer(args[0])
	server, err := cmd.API.GetServer(serverId)
	if err != nil {
		log.Fatalf("Failed to get server information for %s: %v", serverId, err)
	}

	command := []string{"netstat -lutn 2>/dev/null | grep LISTEN"}
	err = serverExec(server.PublicAddress.IP, command)
	if err != nil {
		log.Fatalf("Command execution failed: %v", err)
	}
}
