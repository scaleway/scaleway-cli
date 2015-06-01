package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
)

var cmdExec = &Command{
	Exec:        runExec,
	UsageLine:   "exec [OPTIONS] SERVER COMMAND [ARGS...]", // FIXME: add [ARGS...] support
	Description: "Run a command on a running server",
	Help:        "Run a command on a running server.",
}

func init() {
	// FIXME: -h
	cmdExec.Flag.BoolVar(&execW, []string{"w", "-wait"}, false, "Wait for SSH to be ready")
}

// Flags
var execW bool // -w flag

func NewSshExecCmd(ipAddress string, allocateTTY bool, command []string) []string {
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

	execCmd = append(execCmd, "--", "/bin/sh", "-e")

	if os.Getenv("DEBUG") == "1" {
		execCmd = append(execCmd, "-x")
	}

	execCmd = append(execCmd, "-c")

	execCmd = append(execCmd, fmt.Sprintf("%q", strings.Join(command, " ")))

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
		time.Sleep(1 * time.Second)
	}

	return server, nil
}

func WaitForTcpPortOpen(dest string) error {
	for {
		conn, err := net.Dial("tcp", dest)
		if err == nil {
			defer conn.Close()
			break
		}
		time.Sleep(1 * time.Second)
	}
	return nil
}

func WaitForServerReady(api *ScalewayAPI, serverId string) (*ScalewayServer, error) {
	server, err := WaitForServerState(api, serverId, "running")
	if err != nil {
		return nil, err
	}

	dest := fmt.Sprintf("%s:22", server.PublicAddress.IP)

	err = WaitForTcpPortOpen(dest)
	if err != nil {
		return nil, err
	}

	return server, nil
}

func runExec(cmd *Command, args []string) {
	if len(args) < 2 {
		log.Fatalf("usage: scw %s", cmd.UsageLine)
	}

	serverId := cmd.GetServer(args[0])

	var server *ScalewayServer
	var err error
	if execW {
		// --wait
		server, err = WaitForServerReady(cmd.API, serverId)
		if err != nil {
			log.Fatalf("Failed to wait for server to be ready, %v", err)
		}
	} else {
		// no --wait
		server, err = cmd.API.GetServer(serverId)
		if err != nil {
			log.Fatalf("Failed to get server information for %s: %s", server.Identifier, err)
		}
	}

	execCmd := append(NewSshExecCmd(server.PublicAddress.IP, true, args[1:]))

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
