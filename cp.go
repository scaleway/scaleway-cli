package main

import (
	"fmt"
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

func TarFromSource(api *ScalewayAPI, source string) (*io.ReadCloser, error) {
	var tarOutputStream io.ReadCloser

	// source is a server address + path (scp-like uri)
	if strings.Index(source, ":") > -1 {
		log.Debugf("Creating a tarball remotely and streaming it using SSH")
		serverParts := strings.Split(source, ":")
		if len(serverParts) != 2 {
			return nil, fmt.Errorf("invalid source uri, see 'scw cp -h' for usage")
		}

		serverID := api.GetServerID(serverParts[0])

		server, err := api.GetServer(serverID)
		if err != nil {
			return nil, err
		}

		dir, base := PathToTARPathparts(serverParts[1])
		log.Debugf("Equivalent to 'scp root@%s:%s/%s ...'", server.PublicAddress.IP, dir, base)

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
		spawnSrc := exec.Command("ssh", execCmd...)

		tarOutputStream, err = spawnSrc.StdoutPipe()
		if err != nil {
			return nil, err
		}
		defer tarOutputStream.Close()

		tarErrorStream, err := spawnSrc.StderrPipe()
		if err != nil {
			return nil, err
		}
		defer tarErrorStream.Close()
		io.Copy(os.Stderr, tarErrorStream)

		err = spawnSrc.Start()
		if err != nil {
			return nil, err
		}
		defer spawnSrc.Wait()

		return &tarOutputStream, nil
	}

	// source is stdin
	if source == "-" {
		log.Debugf("Streaming tarball from stdin")
		tarOutputStream = os.Stdin
		defer tarOutputStream.Close()
		return &tarOutputStream, nil
	}

	// source is a path on localhost
	log.Debugf("Taring local path %s", source)
	path, err := filepath.Abs(source)
	if err != nil {
		return nil, err
	}
	path, err = filepath.EvalSymlinks(path)
	if err != nil {
		return nil, err
	}
	log.Debugf("Real local path is %s", path)

	dir, base := PathToTARPathparts(path)

	tarOutputStream, err = archive.TarWithOptions(dir, &archive.TarOptions{
		Compression:  archive.Uncompressed,
		IncludeFiles: []string{base},
	})
	if err != nil {
		return nil, err
	}
	defer tarOutputStream.Close()
	return &tarOutputStream, nil
}

func UntarToDest(api *ScalewayAPI, sourceStream *io.ReadCloser, destination string) error {
	// destination is a server address + path (scp-like uri)
	if strings.Index(destination, ":") > -1 {
		log.Debugf("Streaming using ssh and untaring remotely")
		serverParts := strings.Split(destination, ":")
		if len(serverParts) != 2 {
			return fmt.Errorf("invalid destination uri, see 'scw cp -h' for usage")
		}

		serverID := api.GetServerID(serverParts[0])

		server, err := api.GetServer(serverID)
		if err != nil {
			return err
		}

		// remoteCommand is executed on the remote server
		// it streams a tarball raw content
		remoteCommand := []string{"tar"}
		remoteCommand = append(remoteCommand, "-C", serverParts[1])
		if os.Getenv("DEBUG") == "1" {
			remoteCommand = append(remoteCommand, "-v")
		}
		remoteCommand = append(remoteCommand, "-xf", "-")

		// execCmd contains the ssh connection + the remoteCommand
		execCmd := append(NewSSHExecCmd(server.PublicAddress.IP, false, remoteCommand))
		log.Debugf("Executing: ssh %s", strings.Join(execCmd, " "))
		spawnDst := exec.Command("ssh", execCmd...)

		untarInputStream, err := spawnDst.StdinPipe()
		if err != nil {
			return err
		}
		defer untarInputStream.Close()

		untarErrorStream, err := spawnDst.StderrPipe()
		if err != nil {
			return err
		}
		defer untarErrorStream.Close()

		untarOutputStream, err := spawnDst.StdoutPipe()
		if err != nil {
			return err
		}
		defer untarOutputStream.Close()

		err = spawnDst.Start()
		if err != nil {
			return err
		}
		defer spawnDst.Wait()

		io.Copy(untarInputStream, *sourceStream)
		io.Copy(os.Stderr, untarErrorStream)
		_, err = io.Copy(os.Stdout, untarOutputStream)
		return err
	}

	// destination is stdout
	if destination == "-" { // stdout
		log.Debugf("Writing sourceStream(%v) to os.Stdout(%v)", sourceStream, os.Stdout)
		_, err := io.Copy(os.Stdout, *sourceStream)
		return err
	}

	// destination is a path on localhost
	log.Debugf("Untaring to local path: %s", destination)
	err := archive.Untar(*sourceStream, destination, &archive.TarOptions{NoLchown: true})
	return err
}

func runCp(cmd *Command, args []string) {
	if cpHelp {
		cmd.PrintUsage()
	}
	if len(args) != 2 {
		cmd.PrintShortUsage()
	}

	if strings.Count(args[0], ":") > 1 || strings.Count(args[1], ":") > 1 {
		log.Fatalf("usage: scw %s", cmd.UsageLine)
	}

	sourceStream, err := TarFromSource(cmd.API, args[0])
	if err != nil {
		log.Fatalf("Cannot tar from source '%s': %v", args[0], err)
	}

	err = UntarToDest(cmd.API, sourceStream, args[1])
	if err != nil {
		log.Fatalf("Cannot untar to destionation '%s': %v", args[1], err)
	}
}
