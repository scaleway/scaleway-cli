package main

import (
	"os"
	"os/exec"
	"strings"

	log "github.com/Sirupsen/logrus"
)

var cmdKill = &Command{
	Exec:        runKill,
	UsageLine:   "kill [OPTIONS] SERVER",
	Description: "Kill a running server",
	Help:        "Kill a running server.",
}

// FIXME: add --signal option

func runKill(cmd *Command, args []string) {
	if len(args) < 1 {
		log.Fatalf("usage: scw %s", cmd.UsageLine)
	}
	serverId := cmd.GetServer(args[0])
	command := "halt"
	server, err := cmd.API.GetServer(serverId)
	if err != nil {
		log.Fatalf("failed to get server information for %s: %s", server.Identifier, err)
	}

	execCmd := append(NewSshExecCmd(server.PublicAddress.IP, true), "--", command)

	log.Debugf("Executing: ssh %s", strings.Join(execCmd, " "))

	spawn := exec.Command("ssh", execCmd...)
	spawn.Stdout = os.Stdout
	spawn.Stdin = os.Stdin
	spawn.Stderr = os.Stderr
	err = spawn.Run()
	if err != nil {
		log.Fatal(err)
	}
}
