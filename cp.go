package main

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
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
    $ scw cp myserver:path/to/file path/to/my/local/dir
    $ scw cp myserver:path/to/file myserver2:path/to/dir
    $ scw cp myserver:path/to/file - > myserver-pathtofile-backup.tar
    $ scw cp myserver:path/to/file - | tar -tvf -
    $ scw cp path/to/my/local/dir  myserver:path
    $ scw cp myserver:path/to/dir  path/to/my/local/dir
    $ scw cp myserver:path/to/dir  myserver2:path/to/dir
    $ scw cp myserver:path/to/dir  - > myserver-pathtodir-backup.tar
    $ scw cp myserver:path/to/dir  - | tar -tvf -
    $ cat archive.tar | scw cp - myserver:/path
    $ tar -cvf - . | scw cp - myserver:path
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
	var tarErrorStream io.ReadCloser

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

		dir, base := PathToTARPathparts(serverParts[1])

		// remoteCommand is executed on the remote server
		// it streams a tarball raw content
		remoteCommand := []string{"tar"}
		remoteCommand = append(remoteCommand, "-C", dir)
		if os.Getenv("DEBUG") == "1" {
			remoteCommand = append(remoteCommand, "-v")
		}
		remoteCommand = append(remoteCommand, "-cf", "-")
		remoteCommand = append(remoteCommand, base)

		// execCmd contains the ssh connection + the remoteCommand
		execCmd := append(NewSSHExecCmd(server.PublicAddress.IP, false, remoteCommand))
		log.Debugf("Executing: ssh %s", strings.Join(execCmd, " "))
		spawn := exec.Command("ssh", execCmd...)

		tarOutputStream, err = spawn.StdoutPipe()
		if err != nil {
			log.Fatal(err)
		}
		tarErrorStream, err = spawn.StderrPipe()
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
		tarOutputStream = os.Stdin
	} else { // source host path
		log.Debugf("Creating tarball of local path %s", source)
		path, err := filepath.Abs(source)
		if err != nil {
			log.Fatalf("Cannot tar local path: %v", err)
		}
		path, err = filepath.EvalSymlinks(path)
		if err != nil {
			log.Fatalf("Cannot tar local path: %v", err)
		}
		log.Debugf("Real local path is %s", path)

		dir, base := PathToTARPathparts(path)

		tarOutputStream, err = archive.TarWithOptions(dir, &archive.TarOptions{
			Compression:  archive.Uncompressed,
			IncludeFiles: []string{base},
		})
		if err != nil {
			log.Fatalf("Cannot tar local path: %v", err)
		}
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
