package main

import (
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
)

func sshExec(ipAddress string, command []string, checkConnection bool) error {
	if ipAddress == "" {
		return errors.New("Server does not have public IP")
	}

	if checkConnection {
		if !IsTCPPortOpen(fmt.Sprintf("%s:22", ipAddress)) {
			return errors.New("Server is not ready, try again later.")
		}
	}

	execCmd := append(NewSSHExecCmd(ipAddress, true, command))

	log.Debugf("Executing: ssh %s", strings.Join(execCmd, " "))
	spawn := exec.Command("ssh", execCmd...)
	spawn.Stdout = os.Stdout
	spawn.Stdin = os.Stdin
	spawn.Stderr = os.Stderr
	return spawn.Run()
}

// NewSSHExecCmd computes execve compatible arguments to run a command via ssh
func NewSSHExecCmd(ipAddress string, allocateTTY bool, command []string) []string {
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

// WaitForTCPPortOpen calls IsTCPPortOpen in a loop
func WaitForTCPPortOpen(dest string) error {
	for {
		if IsTCPPortOpen(dest) {
			break
		}
		time.Sleep(1 * time.Second)
	}
	return nil
}

// IsTCPPortOpen returns true if a TCP communication with "host:port" can be initialized
func IsTCPPortOpen(dest string) bool {
	conn, err := net.DialTimeout("tcp", dest, time.Duration(2000)*time.Millisecond)
	if err == nil {
		defer conn.Close()
	}
	return err == nil
}

// truncIf ensures the input string does not exceed max size if cond is met
func truncIf(str string, max int, cond bool) string {
	if cond && len(str) > max {
		return str[:max]
	}
	return str
}

// wordify convert complex name to a single word without special shell characters
func wordify(str string) string {
	str = regexp.MustCompile(`[^a-zA-Z0-9-]`).ReplaceAllString(str, "_")
	str = regexp.MustCompile(`__+`).ReplaceAllString(str, "_")
	str = strings.Trim(str, "_")
	return str
}
