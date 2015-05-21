package main

import (
	"os"
	"os/exec"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
)

var cmdExec = &Command{
	Exec:        runExec,
	UsageLine:   "exec [OPTIONS] SERVER COMMAND", // FIXME: add [ARGS...] support
	Description: "Run a command on a running server",
	Help:        "Run a command on a running server.",
}

func init() {
	// FIXME: -h
	cmdExec.Flag.BoolVar(&execW, []string{"w", "-wait"}, false, "")
}

// Flags
var execW bool // -w flag

func NewSshExecCmd(ipAddress string, allocateTTY bool) []string {
	execCmd := []string{}

	if os.Getenv("DEBUG") != "1" {
		execCmd = append(execCmd, "-q")
	}

	if os.Getenv("exec_secure") != "1" {
		execCmd = append(execCmd, "-o", "UserKnownHostsFile=/dev/null", "-o", "StrictHostKeyChecking=no")
	}

	execCmd = append(execCmd, "-l", "root", ipAddress)

	if allocateTTY {
		execCmd = append(execCmd, "-t")
	}

	return execCmd
}

func WaitForServerState(api *ScalewayAPI, serverId string, targetState string) (*ScalewayServer, error) {
	var server *ScalewayServer
	var err error

	for {
		server, err = api.GetServer(serverId)
		if err != nil {
			return nil, err
		}
		if server.State == targetState {
			break
		}
		time.Sleep(2)
	}

	return server, nil
}

func runExec(cmd *Command, args []string) {
	if len(args) < 2 {
		log.Fatalf("usage: scw %s", cmd.UsageLine)
	}
	command := args[1]

	serverId := cmd.GetServer(args[0])

	var server *ScalewayServer
	var err error

	if execW {
		// --wait
		server, err = WaitForServerState(cmd.API, serverId, "running")
		if err != nil {
			log.Fatalf("Failed to wait for server to be ready, %v", err)
		}
	} else {
		// no --wait
		server, err := cmd.API.GetServer(serverId)
		if err != nil {
			log.Fatalf("Failed to get server information for %s: %s", server.Identifier, err)
		}
	}

	execCmd := append(NewSshExecCmd(server.PublicAddress.IP, true), "--", command)

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
