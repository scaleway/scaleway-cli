// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package commands

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/docker/docker/pkg/archive"
	"github.com/scaleway/scaleway-cli/pkg/api"
	"github.com/scaleway/scaleway-cli/pkg/utils"
)

// CpArgs are arguments passed to `RunCp`
type CpArgs struct {
	Gateway     string
	Source      string
	Destination string
	SSHUser     string
	SSHPort     int
}

// RunCp is the handler for 'scw cp'
func RunCp(ctx CommandContext, args CpArgs) error {
	if strings.Count(args.Source, ":") > 1 || strings.Count(args.Destination, ":") > 1 {
		return fmt.Errorf("bad usage, see 'scw help cp'")
	}

	sourceStream, err := TarFromSource(ctx, args.Source, args.Gateway, args.SSHUser, args.SSHPort)
	if err != nil {
		return fmt.Errorf("cannot tar from source '%s': %v", args.Source, err)
	}

	err = UntarToDest(ctx, sourceStream, args.Destination, args.Gateway, args.SSHUser, args.SSHPort)
	if err != nil {
		return fmt.Errorf("cannot untar to destination '%s': %v", args.Destination, err)
	}
	return nil
}

// TarFromSource creates a stream buffer with the tarballed content of the user source
func TarFromSource(ctx CommandContext, source, gateway, user string, port int) (*io.ReadCloser, error) {
	var tarOutputStream io.ReadCloser

	// source is a server address + path (scp-like uri)
	if strings.Contains(source, ":") {
		logrus.Debugf("Creating a tarball remotely and streaming it using SSH")
		serverParts := strings.Split(source, ":")
		if len(serverParts) != 2 {
			return nil, fmt.Errorf("invalid source uri, see 'scw cp -h' for usage")
		}

		serverID, err := ctx.API.GetServerID(serverParts[0])
		if err != nil {
			return nil, err
		}

		server, err := ctx.API.GetServer(serverID)
		if err != nil {
			return nil, err
		}

		dir, base := utils.PathToTARPathparts(serverParts[1])
		logrus.Debugf("Equivalent to 'scp root@%s:%s/%s ...'", server.PublicAddress.IP, dir, base)

		// remoteCommand is executed on the remote server
		// it streams a tarball raw content
		remoteCommand := []string{"tar"}
		remoteCommand = append(remoteCommand, "-C", dir)
		if ctx.Getenv("DEBUG") == "1" {
			remoteCommand = append(remoteCommand, "-v")
		}
		remoteCommand = append(remoteCommand, "-cf", "-")
		remoteCommand = append(remoteCommand, base)

		// Resolve gateway
		if gateway == "" {
			gateway = ctx.Getenv("SCW_GATEWAY")
		}

		if gateway == serverID || gateway == serverParts[0] {
			gateway = ""
		} else {
			gateway, err = api.ResolveGateway(ctx.API, gateway)
			if err != nil {
				return nil, fmt.Errorf("cannot resolve Gateway '%s': %v", gateway, err)
			}
		}

		// execCmd contains the ssh connection + the remoteCommand
		sshCommand := utils.NewSSHExecCmd(server.PublicAddress.IP, server.PrivateIP, user, port, false, remoteCommand, gateway, false)
		logrus.Debugf("Executing: %s", sshCommand)
		spawnSrc := exec.Command("ssh", sshCommand.Slice()[1:]...)

		tarOutputStream, err = spawnSrc.StdoutPipe()
		if err != nil {
			return nil, err
		}

		tarErrorStream, err := spawnSrc.StderrPipe()
		if err != nil {
			return nil, err
		}
		defer tarErrorStream.Close()
		io.Copy(ctx.Stderr, tarErrorStream)

		err = spawnSrc.Start()
		if err != nil {
			return nil, err
		}
		defer spawnSrc.Wait()

		return &tarOutputStream, nil
	}

	// source is stdin
	if source == "-" {
		logrus.Debugf("Streaming tarball from stdin")
		// FIXME: should be ctx.Stdin
		tarOutputStream = os.Stdin
		return &tarOutputStream, nil
	}

	// source is a path on localhost
	logrus.Debugf("Taring local path %s", source)
	path, err := filepath.Abs(source)
	if err != nil {
		return nil, err
	}
	path, err = filepath.EvalSymlinks(path)
	if err != nil {
		return nil, err
	}
	logrus.Debugf("Real local path is %s", path)

	dir, base := utils.PathToTARPathparts(path)

	tarOutputStream, err = archive.TarWithOptions(dir, &archive.TarOptions{
		Compression:  archive.Uncompressed,
		IncludeFiles: []string{base},
	})
	if err != nil {
		return nil, err
	}
	return &tarOutputStream, nil
}

// UntarToDest writes to user destination the streamed tarball in input
func UntarToDest(ctx CommandContext, sourceStream *io.ReadCloser, destination, gateway, user string, port int) error {
	// destination is a server address + path (scp-like uri)
	if strings.Contains(destination, ":") {
		logrus.Debugf("Streaming using ssh and untaring remotely")
		serverParts := strings.Split(destination, ":")
		if len(serverParts) != 2 {
			return fmt.Errorf("invalid destination uri, see 'scw cp -h' for usage")
		}

		serverID, err := ctx.API.GetServerID(serverParts[0])
		if err != nil {
			return err
		}

		server, err := ctx.API.GetServer(serverID)
		if err != nil {
			return err
		}

		// remoteCommand is executed on the remote server
		// it streams a tarball raw content
		remoteCommand := []string{"tar"}
		remoteCommand = append(remoteCommand, "-C", serverParts[1])
		if ctx.Getenv("DEBUG") == "1" {
			remoteCommand = append(remoteCommand, "-v")
		}
		remoteCommand = append(remoteCommand, "-xf", "-")

		// Resolve gateway
		if gateway == "" {
			gateway = ctx.Getenv("SCW_GATEWAY")
		}
		if gateway == serverID || gateway == serverParts[0] {
			gateway = ""
		} else {
			gateway, err = api.ResolveGateway(ctx.API, gateway)
			if err != nil {
				return fmt.Errorf("cannot resolve Gateway '%s': %v", gateway, err)
			}
		}

		// execCmd contains the ssh connection + the remoteCommand
		sshCommand := utils.NewSSHExecCmd(server.PublicAddress.IP, server.PrivateIP, user, port, false, remoteCommand, gateway, false)
		logrus.Debugf("Executing: %s", sshCommand)
		spawnDst := exec.Command("ssh", sshCommand.Slice()[1:]...)

		untarInputStream, err := spawnDst.StdinPipe()
		if err != nil {
			return err
		}
		defer untarInputStream.Close()

		// spawnDst.Stderr = ctx.Stderr
		// spawnDst.Stdout = ctx.Stdout

		err = spawnDst.Start()
		if err != nil {
			return err
		}

		_, err = io.Copy(untarInputStream, *sourceStream)
		return err
	}

	// destination is stdout
	if destination == "-" { // stdout
		logrus.Debugf("Writing sourceStream(%v) to ctx.Stdout(%v)", sourceStream, ctx.Stdout)
		_, err := io.Copy(ctx.Stdout, *sourceStream)
		return err
	}

	// destination is a path on localhost
	logrus.Debugf("Untaring to local path: %s", destination)
	err := archive.Untar(*sourceStream, destination, &archive.TarOptions{NoLchown: true})
	return err
}
