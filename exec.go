package main

import (
	"errors"
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
	UsageLine:   "exec [OPTIONS] SERVER COMMAND [ARGS...]",
	Description: "Run a command on a running server",
	Help:        "Run a command on a running server.",
	Examples: `
    $ scw exec myserver bash
    $ scw exec myserver 'tmux a -t joe || tmux new -s joe || bash'
    $ exec_secure=1 scw exec myserver bash
    $ scw exec -w $(scw start $(scw create ubuntu-trusty)) bash
    $ scw exec $(scw start -w $(scw create ubuntu-trusty)) bash
    $ scw exec myserver tmux new -d sleep 10
    $ scw exec myserver ls -la | grep password
`,
}

func init() {
	cmdExec.Flag.BoolVar(&execHelp, []string{"h", "-help"}, false, "Print usage")
	cmdExec.Flag.Float64Var(&execTimeout, []string{"T", "-timeout"}, 0, "Set timeout values to seconds")
	cmdExec.Flag.BoolVar(&execW, []string{"w", "-wait"}, false, "Wait for SSH to be ready")
}

// Flags
var execW bool          // -w, --wait flag
var execTimeout float64 // -T flag
var execHelp bool       // -h, --help flag

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

func IsTcpPortOpen(dest string) bool {
	conn, err := net.Dial("tcp", dest)
	if err == nil {
		defer conn.Close()
	}
	return err == nil
}

func WaitForTcpPortOpen(dest string) error {
	for {
		if IsTcpPortOpen(dest) {
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

func serverExec(server *ScalewayServer, command []string, checkConnection bool) error {
	ipAddress := server.PublicAddress.IP
	if ipAddress == "" {
		return errors.New("Server does not have public IP")
	}

	if checkConnection {
		if !IsTcpPortOpen(fmt.Sprintf("%s:22", ipAddress)) {
			return errors.New("Server is not ready, try again later.")
		}
	}

	execCmd := append(NewSshExecCmd(ipAddress, true, command))

	log.Debugf("Executing: ssh %s", strings.Join(execCmd, " "))
	spawn := exec.Command("ssh", execCmd...)
	spawn.Stdout = os.Stdout
	spawn.Stdin = os.Stdin
	spawn.Stderr = os.Stderr
	return spawn.Run()
}

func runExec(cmd *Command, args []string) {
	if execHelp {
		cmd.PrintUsage()
	}
	if len(args) < 2 {
		cmd.PrintShortUsage()
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
			log.Fatalf("Failed to get server information for %s: %v", serverId, err)
		}
	}

	if execTimeout > 0 {
		go func() {
			time.Sleep(time.Duration(execTimeout*1000) * time.Millisecond)
			log.Fatalf("Operation timed out")
		}()
	}

	err = serverExec(server, args[1:], !execW)
	if err != nil {
		log.Fatalf("%v", err)
		os.Exit(1)
	}
	log.Debugf("Command successfuly executed")
}
