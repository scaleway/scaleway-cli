package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	log "github.com/Sirupsen/logrus"
)

var cmdTop = &Command{
	Exec:        runTop,
	UsageLine:   "top [OPTIONS] SERVER", // FIXME: add ps options
	Description: "Lookup the running processes of a server",
	Help:        "Lookup the running processes of a server.",
}

func runTop(cmd *Command, args []string) {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "usage: scw %s\n", cmd.UsageLine)
		os.Exit(1)
	}
	serverId := cmd.GetServer(args[0])
	command := "ps"
	server, err := cmd.API.GetServer(serverId)
	if err != nil {
		log.Fatalf("failed to get server information for %s: %s\n", server.Identifier, err)
	}

	execCmd := append(NewSshExecCmd(server.PublicAddress.IP), "--", command)

	log.Debugf("Executing: ssh %s", strings.Join(execCmd, " "))
	out, err := exec.Command("ssh", execCmd...).CombinedOutput()
	fmt.Printf("%s", out)
	if err != nil {
		log.Fatal(err)
	}
}
