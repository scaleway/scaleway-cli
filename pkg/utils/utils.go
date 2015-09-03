// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

// scw helpers

// Package utils contains helpers
package utils

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/scaleway/scaleway-cli/pkg/sshcommand"
	log "github.com/scaleway/scaleway-cli/vendor/github.com/Sirupsen/logrus"
)

// SpawnRedirection is used to redirects the fluxes
type SpawnRedirection struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

// SSHExec executes a command over SSH and redirects file-descriptors
func SSHExec(publicIPAddress string, privateIPAddress string, command []string, checkConnection bool, gateway string) error {
	gatewayUser := "root"
	gatewayIPAddress := gateway
	if strings.Contains(gateway, "@") {
		parts := strings.Split(gatewayIPAddress, "@")
		gatewayUser = parts[0]
		gatewayIPAddress = parts[1]
		gateway = gatewayUser + "@" + gatewayIPAddress
	}

	if publicIPAddress == "" && gatewayIPAddress == "" {
		return errors.New("server does not have public IP")
	}
	if privateIPAddress == "" && gatewayIPAddress != "" {
		return errors.New("server does not have private IP")
	}

	if checkConnection {
		useGateway := gatewayIPAddress != ""
		if useGateway && !IsTCPPortOpen(fmt.Sprintf("%s:22", gatewayIPAddress)) {
			return errors.New("gateway is not available, try again later")
		}
		if !useGateway && !IsTCPPortOpen(fmt.Sprintf("%s:22", publicIPAddress)) {
			return errors.New("server is not ready, try again later")
		}
	}

	sshCommand := NewSSHExecCmd(publicIPAddress, privateIPAddress, true, command, gateway)

	log.Debugf("Executing: %s", sshCommand)

	spawn := exec.Command("ssh", sshCommand.Slice()[1:]...)
	spawn.Stdout = os.Stdout
	spawn.Stdin = os.Stdin
	spawn.Stderr = os.Stderr
	return spawn.Run()
}

// NewSSHExecCmd computes execve compatible arguments to run a command via ssh
func NewSSHExecCmd(publicIPAddress string, privateIPAddress string, allocateTTY bool, command []string, gatewayIPAddress string) *sshcommand.Command {
	quiet := os.Getenv("DEBUG") != "1"
	secureExec := os.Getenv("exec_secure") == "1"
	sshCommand := &sshcommand.Command{
		AllocateTTY:         true,
		Command:             command,
		Host:                publicIPAddress,
		Quiet:               quiet,
		SkipHostKeyChecking: !secureExec,
		User:                "root",
		NoEscapeCommand:     true,
	}
	if gatewayIPAddress != "" {
		sshCommand.Host = privateIPAddress
		sshCommand.Gateway = &sshcommand.Command{
			Host:                gatewayIPAddress,
			SkipHostKeyChecking: !secureExec,
			AllocateTTY:         true,
			Quiet:               quiet,
			User:                "root",
		}
	}

	return sshCommand
}

// GeneratingAnSSHKey generates an SSH key
func GeneratingAnSSHKey(cfg SpawnRedirection, path string, name string) (string, error) {
	args := []string{
		"-t",
		"rsa",
		"-b",
		"4096",
		"-f",
		filepath.Join(path, name),
		"-N",
		"",
		"-C",
		"",
	}
	log.Infof("Executing commands %v", args)
	spawn := exec.Command("ssh-keygen", args...)
	spawn.Stdout = cfg.Stdout
	spawn.Stdin = cfg.Stdin
	spawn.Stderr = cfg.Stderr
	return args[5], spawn.Run()
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

func SSHGetFingerprint(key string) (string, error) {
	tmp, err := ioutil.TempFile("", ".tmp")
	if err != nil {
		return "", fmt.Errorf("Unable to create a tempory file: %v", err)
	}
	defer os.Remove(tmp.Name())
	buff := []byte(key)
	bytesWritten := 0
	for bytesWritten < len(buff) {
		nb, err := tmp.Write(buff[bytesWritten:])
		if err != nil {
			return "", fmt.Errorf("Unable to write: %v", err)
		}
		bytesWritten += nb
	}
	ret, err := exec.Command("ssh-keygen", "-l", "-f", tmp.Name()).Output()
	if err != nil {
		return "", fmt.Errorf("Unable to run ssh-keygen: %v", err)
	}
	return string(ret), nil
}
