package main

import (
	"fmt"
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

func runExec(cmd *Command, args []string) {
	if len(args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: scw %s\n", cmd.UsageLine)
		os.Exit(1)
	}
	serverId := cmd.GetServer(args[0])
	command := args[1]
	server, err := cmd.API.GetServer(serverId)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get server information for %s: %s\n", server, err)
	}
	execCmd := []string{"-l", "root", server.PublicAddress.IP, "-t", "--", command}
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
