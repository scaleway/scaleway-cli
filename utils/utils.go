// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

// scw helpers
package utils

import (
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
)

// SSHExec executes a command over SSH and redirects file-descriptors
func SSHExec(ipAddress string, command []string, checkConnection bool) error {
	if ipAddress == "" {
		return errors.New("server does not have public IP")
	}

	if checkConnection {
		if !IsTCPPortOpen(fmt.Sprintf("%s:22", ipAddress)) {
			return errors.New("server is not ready, try again later")
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

	if len(command) > 0 {
		execCmd = append(execCmd, "--", "/bin/sh", "-e")

		if os.Getenv("DEBUG") == "1" {
			execCmd = append(execCmd, "-x")
		}

		execCmd = append(execCmd, "-c")

		execCmd = append(execCmd, fmt.Sprintf("%q", strings.Join(command, " ")))
	}

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

// TruncIf ensures the input string does not exceed max size if cond is met
func TruncIf(str string, max int, cond bool) string {
	if cond && len(str) > max {
		return str[:max]
	}
	return str
}

// Wordify convert complex name to a single word without special shell characters
func Wordify(str string) string {
	str = regexp.MustCompile(`[^a-zA-Z0-9-]`).ReplaceAllString(str, "_")
	str = regexp.MustCompile(`__+`).ReplaceAllString(str, "_")
	str = strings.Trim(str, "_")
	return str
}

// PathToTARPathparts returns the two parts of a unix path
func PathToTARPathparts(fullPath string) (string, string) {
	fullPath = strings.TrimRight(fullPath, "/")
	return path.Dir(fullPath), path.Base(fullPath)
}

// RemoveDuplicates transforms an array into a unique array
func RemoveDuplicates(elements []string) []string {
	encountered := map[string]bool{}

	// Create a map of all unique elements.
	for v := range elements {
		encountered[elements[v]] = true
	}

	// Place all keys from the map into a slice.
	result := []string{}
	for key := range encountered {
		result = append(result, key)
	}
	return result
}

// GetConfigFilePath returns the path to the Scaleway CLI config file
func GetConfigFilePath() (string, error) {
	homeDir := os.Getenv("HOME") // *nix
	if homeDir == "" {           // Windows
		homeDir = os.Getenv("USERPROFILE")
	}
	if homeDir == "" {
		return "", errors.New("user home directory not found")
	}

	return filepath.Join(homeDir, ".scwrc"), nil
}

const termjsBin string = "termjs-cli"

// AttachToSerial tries to connect to server serial using 'term.js-cli' and fallback with a help message
func AttachToSerial(serverID string, apiToken string, attachStdin bool) error {
	termjsURL := fmt.Sprintf("https://tty.cloud.online.net?server_id=%s&type=serial&auth_token=%s", serverID, apiToken)

	args := []string{}
	if !attachStdin {
		args = append(args, "--no-stdin")
	}
	args = append(args, termjsURL)
	log.Debugf("Executing: %s %v", termjsBin, args)
	// FIXME: check if termjs-cli is installed
	spawn := exec.Command(termjsBin, args...)
	spawn.Stdout = os.Stdout
	spawn.Stdin = os.Stdin
	spawn.Stderr = os.Stderr
	err := spawn.Run()
	if err != nil {
		log.Warnf(`
You need to install '%s' from https://github.com/moul/term.js-cli

    npm install -g term.js-cli

However, you can access your serial using a web browser:

    %s

`, termjsBin, termjsURL)
		return err
	}
	return nil
}
