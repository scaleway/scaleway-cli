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
	UsageLine:   "cp [OPTIONS] SERVER:PATH|HOSTPATH|- SERVER:PATH|HOSTPATH|-",
	Description: "Copy files/folders from a PATH on the server to a HOSTDIR on the host",
	Help:        "Copy files/folders from a PATH on the server to a HOSTDIR on the host\nrunning the command. Use '-' to write the data as a tar file to STDOUT.",
	Examples: `
    $ scw cp path/to/my/local/file myserver:path
    $ scw cp myserver:path path/to/my/local/file
    $ scw cp myserver:path myserver2:path
`,
}

func init() {
	cmdCp.Flag.BoolVar(&cpHelp, []string{"h", "-help"}, false, "Print usage")
}

// Flags
var cpHelp bool // -h, --help flag

func runCp(cmd *Command, args []string) {
	if cpHelp {
		cmd.PrintUsage()
	}
	if len(args) != 2 {
		cmd.PrintShortUsage()
	}

	var tarOutputStream io.ReadCloser

	// source
	source := args[0]
	if strings.Index(source, ":") > -1 { // source server address
		serverParts := strings.Split(args[0], ":")
		if len(serverParts) != 2 {
			log.Fatalf("usage: scw %s", cmd.UsageLine)
		}

		serverID := cmd.API.GetServerID(serverParts[0])

		server, err := cmd.API.GetServer(serverID)
		if err != nil {
			log.Fatalf("Failed to get server information for %s: %v", serverID, err)
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
		execCmd := append(NewSSHExecCmd(server.PublicAddress.IP, false, remoteCommand))
		log.Debugf("Executing: ssh %s", strings.Join(execCmd, " "))
		spawn := exec.Command("ssh", execCmd...)

		tarOutputStream, err = spawn.StdoutPipe()
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
	} else if source == "-" { // stdin
		log.Fatalf("'scw cp - ...' is not yet implemented")
	} else { // source host path
		log.Fatalf("'scw cp HOSTPATH ...' is not yet implemented")
	}

	// destination
	destination := args[1]
	if strings.Index(destination, ":") > -1 { // destination server address
		log.Fatalf("'scw cp ... SERVER' is not yet implemented")
	} else if destination == "-" { // stdout
		log.Debugf("Writing tarOutputStream(%v) to os.Stdout(%v)", tarOutputStream, os.Stdout)
		written, err := io.Copy(os.Stdout, tarOutputStream)
		log.Debugf("%d bytes written", written)
		if err != nil {
			log.Fatal(err)
		}

	} else { // destination host path
		err := archive.Untar(tarOutputStream, destination, &archive.TarOptions{NoLchown: true})
		if err != nil {
			log.Fatalf("Failed to untar the remote archive: %v", err)
		}
	}
}
