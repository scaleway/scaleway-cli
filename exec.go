package main

import (
	"os"
	"os/exec"
	"strings"

	log "github.com/Sirupsen/logrus"
)

var cmdExec = &Command{
	Exec:        runExec,
	UsageLine:   "exec [OPTIONS] SERVER COMMAND", // FIXME: add [ARGS...] support
	Description: "Run a command on a running server",
	Help:        "Run a command on a running server.",
}

func NewSshExecCmd(ipAddress string) []string {
	execCmd := []string{}

	if os.Getenv("DEBUG") != "1" {
		execCmd = append(execCmd, "-q")
	}

	if os.Getenv("exec_secure") != "1" {
		execCmd = append(execCmd, "-o", "UserKnownHostsFile=/dev/null", "-o", "StrictHostKeyChecking=no")
	}

	execCmd = append(execCmd, "-l", "root", ipAddress, "-t")
	return execCmd
}

func runExec(cmd *Command, args []string) {
	if len(args) < 2 {
		log.Fatalf("usage: scw %s", cmd.UsageLine)
	}
	serverId := cmd.GetServer(args[0])
	command := args[1]
	server, err := cmd.API.GetServer(serverId)
	if err != nil {
		log.Fatalf("failed to get server information for %s: %s", server.Identifier, err)
	}

	execCmd := append(NewSshExecCmd(server.PublicAddress.IP), "--", command)

	log.Debugf("Executing: ssh %s", strings.Join(execCmd, " "))
	spawn := exec.Command("ssh", execCmd...)
	spawn.Stdout = os.Stdout
	spawn.Stdin = os.Stdin
	spawn.Stderr = os.Stderr
	err = spawn.Run()
	if err != nil {
		os.Exit(1)
	}
}
