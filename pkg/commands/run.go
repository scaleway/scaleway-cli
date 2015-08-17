// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/scaleway/scaleway-cli/pkg/api"
	"github.com/scaleway/scaleway-cli/pkg/utils"
	"github.com/scaleway/scaleway-cli/vendor/github.com/Sirupsen/logrus"
)

// RunArgs are flags for the `Run` function
type RunArgs struct {
	Attach     bool
	Bootscript string
	Command    []string
	Detach     bool
	Gateway    string
	Image      string
	Name       string
	Tags       []string
	Volumes    []string
	AutoRemove bool
	TmpSSHKey  bool
	// DynamicIPRequired
	// Timeout
}

// AddSSHKeyToTags adds the ssh key in the tags
func AddSSHKeyToTags(ctx CommandContext, tags *[]string, image string) error {
	home, err := utils.GetHomeDir()
	if err != nil {
		return fmt.Errorf("unable to find your home %v", err)
	}
	idRsa := filepath.Join(home, ".ssh", "id_rsa")
	if _, err := os.Stat(idRsa); err != nil {
		if os.IsNotExist(err) {
			logrus.Warnln("Unable to find your ~/.ssh/id_rsa")
			logrus.Warnln("Run 'ssh-keygen -t rsa'")
			return nil
		}
	}
	idRsa = strings.Join([]string{idRsa, ".pub"}, "")
	data, err := ioutil.ReadFile(idRsa)
	if err != nil {
		return fmt.Errorf("failed to read %v", err)
	}
	data[7] = '_'
	for i := range data {
		if data[i] == ' ' {
			data = data[:i]
			break
		}
	}
	*tags = append(*tags, strings.Join([]string{"AUTHORIZED_KEY", string(data[:len(data)])}, "="))
	return nil
}

// Run is the handler for 'scw run'
func Run(ctx CommandContext, args RunArgs) error {
	if args.Gateway == "" {
		args.Gateway = ctx.Getenv("SCW_GATEWAY")
	}

	if args.TmpSSHKey {
		err := AddSSHKeyToTags(ctx, &args.Tags, args.Image)
		if err != nil {
			return err
		}
	}
	env := strings.Join(args.Tags, " ")
	volume := strings.Join(args.Volumes, " ")

	// create IMAGE
	logrus.Info("Server creation ...")
	dynamicIPRequired := args.Gateway == ""
	serverID, err := api.CreateServer(ctx.API, args.Image, args.Name, args.Bootscript, env, volume, dynamicIPRequired)
	if err != nil {
		return fmt.Errorf("failed to create server: %v", err)
	}
	logrus.Infof("Server created: %s", serverID)

	if args.AutoRemove {
		defer ctx.API.DeleteServerSafe(serverID)
	}

	// start SERVER
	logrus.Info("Server start requested ...")
	err = api.StartServer(ctx.API, serverID, false)
	if err != nil {
		return fmt.Errorf("failed to start server %s: %v", serverID, err)
	}
	logrus.Info("Server is starting, this may take up to a minute ...")

	if args.Detach {
		fmt.Fprintln(ctx.Stdout, serverID)
		return nil
	}

	if args.Attach {
		// Attach to server serial
		logrus.Info("Attaching to server console ...")
		err = utils.AttachToSerial(serverID, ctx.API.Token, true)
		if err != nil {
			return fmt.Errorf("cannot attach to server serial: %v", err)
		}
	} else {
		// Resolve gateway
		gateway, err := api.ResolveGateway(ctx.API, args.Gateway)
		if err != nil {
			return fmt.Errorf("cannot resolve Gateway '%s': %v", args.Gateway, err)
		}

		// waiting for server to be ready
		logrus.Debug("Waiting for server to be ready")
		// We wait for 30 seconds, which is the minimal amount of time needed by a server to boot
		server, err := api.WaitForServerReady(ctx.API, serverID, gateway)
		if err != nil {
			return fmt.Errorf("cannot get access to server %s: %v", serverID, err)
		}
		logrus.Debugf("SSH server is available: %s:22", server.PublicAddress.IP)
		logrus.Info("Server is ready !")

		// exec -w SERVER COMMAND ARGS...
		if len(args.Command) < 1 {
			logrus.Info("Connecting to server ...")
			err = utils.SSHExec(server.PublicAddress.IP, server.PrivateIP, []string{}, false, gateway)
		} else {
			logrus.Infof("Executing command: %s ...", args.Command)
			err = utils.SSHExec(server.PublicAddress.IP, server.PrivateIP, args.Command, false, gateway)
		}
		if err != nil {
			return fmt.Errorf("command execution failed: %v", err)
		}
		logrus.Info("Command successfuly executed")
	}
	return nil
}
