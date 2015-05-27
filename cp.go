package main

import (
	"io"
	"os"
	"os/exec"
	"path"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/pkg/archive"
)

var cmdCp = &Command{
	Exec:        runCp,
	UsageLine:   "cp [OPTIONS] SERVER:PATH HOSTDIR|-",
	Description: "Copy files/folders from a PATH on the server to a HOSTDIR on the host",
	Help:        "Copy files/folders from a PATH on the server to a HOSTDIR on the host\nrunning the the command. Use '-' to write the data as a tar file to STDOUT.",
}

func runCp(cmd *Command, args []string) {
	if len(args) < 1 {
		log.Fatalf("usage: scw %s", cmd.UsageLine)
	}

	hostPath := args[1]

	serverParts := strings.Split(args[0], ":")
	if len(serverParts) != 2 {
		log.Fatalf("usage: scw %s", cmd.UsageLine)
	}

	serverId := cmd.GetServer(serverParts[0])

	server, err := cmd.API.GetServer(serverId)
	if err != nil {
		log.Fatalf("failed to get server information for %s: %s", server.Identifier, err)
	}

	// remoteCommand is executed on the remote server
	// it streams a tarball raw content
	remoteCommand := []string{"tar"}
	remoteCommand = append(remoteCommand, "-C", path.Dir(serverParts[1]))
	if os.Getenv("DEBUG") == "1" {
		remoteCommand = append(remoteCommand, "-v")
	}
	remoteCommand = append(remoteCommand, "-cf", "-")
	remoteCommand = append(remoteCommand, path.Base(serverParts[1]))

	// execCmd contains the ssh connection + the remoteCommand
	execCmd := append(NewSshExecCmd(server.PublicAddress.IP, false, remoteCommand))
	log.Debugf("Executing: ssh %s", strings.Join(execCmd, " "))
	spawn := exec.Command("ssh", execCmd...)

	tarOutputStream, err := spawn.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	tarErrorStream, err := spawn.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}

	err = spawn.Start()
	if err != nil {
		log.Fatalf("Failed to start ssh command: %v", err)
	}

	defer spawn.Wait()

	io.Copy(os.Stderr, tarErrorStream)

	if hostPath == "-" {
		log.Debugf("Writing tarOutputStream(%v) to os.Stdout(%v)", tarOutputStream, os.Stdout)
		written, err := io.Copy(os.Stdout, tarOutputStream)
		log.Debugf("%d bytes written", written)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err = archive.Untar(tarOutputStream, hostPath, &archive.TarOptions{NoLchown: true})
		if err != nil {
			log.Fatalf("Failed to untar the remote archive: %v", err)
		}
	}
}
