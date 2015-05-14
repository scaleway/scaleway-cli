package main

import (
	"fmt"
	"os"
	"os/exec"
)

var cmdExec = &Command{
	Exec:        runExec,
	UsageLine:   "exec [OPTIONS] SERVER COMMAND [ARGS...]",
	Description: "Run a command on a running server",
	Help:        "Run a command on a running server.",
}

func runExec(cmd *Command, args []string) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "usage: scw %s\n", cmd.UsageLine)
		os.Exit(1)
	}
	serverId := cmd.GetServer(args[0])
	server, err := cmd.API.GetServer(serverId)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get server information for %s: %s\n", server, err)
	}
	cmds := exec.Command("ssh", server.PublicAddress.IP)
	cmds.Stdout = os.Stdout
	cmds.Stdin = os.Stdin
	cmds.Stderr = os.Stderr
	err = cmds.Run()
	if err != nil {
		os.Exit(1)
	}
}
