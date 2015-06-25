// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package commands

import (
	"os"
	"time"

	log "github.com/Sirupsen/logrus"

	api "github.com/scaleway/scaleway-cli/api"
	types "github.com/scaleway/scaleway-cli/commands/types"
	"github.com/scaleway/scaleway-cli/utils"
)

var cmdRun = &types.Command{
	Exec:        runRun,
	UsageLine:   "run [OPTIONS] IMAGE [COMMAND] [ARG...]",
	Description: "Run a command in a new server",
	Help:        "Run a command in a new server.",
	Examples: `
    $ scw run ubuntu-trusty
    $ scw run --name=mydocker docker docker run moul/nyancat:armhf
    $ scw run --bootscript=3.2.34 --env="boot=live rescue_image=http://j.mp/scaleway-ubuntu-trusty-tarball" 50GB bash
    $ scw run attach alpine
`,
}

func init() {
	cmdRun.Flag.StringVar(&runCreateName, []string{"-name"}, "", "Assign a name")
	cmdRun.Flag.StringVar(&runCreateBootscript, []string{"-bootscript"}, "", "Assign a bootscript")
	cmdRun.Flag.StringVar(&runCreateEnv, []string{"e", "-env"}, "", "Provide metadata tags passed to initrd (i.e., boot=resue INITRD_DEBUG=1)")
	cmdRun.Flag.StringVar(&runCreateVolume, []string{"v", "-volume"}, "", "Attach additional volume (i.e., 50G)")
	cmdRun.Flag.BoolVar(&runHelpFlag, []string{"h", "-help"}, false, "Print usage")
	cmdRun.Flag.BoolVar(&runAttachFlag, []string{"a", "-attach"}, false, "Attach to serial console")
	// FIXME: handle start --timeout
}

// Flags
var runCreateName string       // --name flag
var runCreateBootscript string // --bootscript flag
var runCreateEnv string        // -e, --env flag
var runCreateVolume string     // -v, --volume flag
var runHelpFlag bool           // -h, --help flag
var runAttachFlag bool         // -a, --attach flag

func runRun(cmd *types.Command, args []string) {
	if runHelpFlag {
		cmd.PrintUsage()
	}
	if len(args) < 1 {
		cmd.PrintShortUsage()
	}
	if runAttachFlag && len(args) > 1 {
		log.Fatalf("Cannot use '--attach' and 'COMMAND [ARG...]' at the same time. See 'scw run --help'")
	}

	// create IMAGE
	log.Debugf("Creating a new server")
	serverID, err := api.CreateServer(cmd.API, args[0], runCreateName, runCreateBootscript, runCreateEnv, runCreateVolume)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}
	log.Debugf("Created server: %s", serverID)

	// start SERVER
	log.Debugf("Starting server")
	err = api.StartServer(cmd.API, serverID, false)
	if err != nil {
		log.Fatalf("Failed to start server %s: %v", serverID, err)
	}
	log.Debugf("Server is booting")

	if runAttachFlag {
		// Attach to server serial
		log.Debugf("Attaching to server console")
		err = utils.AttachToSerial(serverID, cmd.API.Token, true)
		if err != nil {
			log.Fatalf("Cannot attach to server serial: %v", err)
		}
	} else {
		// waiting for server to be ready
		log.Debugf("Waiting for server to be ready")
		// We wait for 30 seconds, which is the minimal amount of time needed by a server to boot
		time.Sleep(30 * time.Second)
		server, err := api.WaitForServerReady(cmd.API, serverID)
		if err != nil {
			log.Fatalf("Cannot get access to server %s: %v", serverID, err)
		}
		log.Debugf("Server is ready: %s", server.PublicAddress.IP)

		// exec -w SERVER COMMAND ARGS...
		log.Debugf("Executing command")
		if len(args) < 2 {
			err = utils.SSHExec(server.PublicAddress.IP, []string{}, false)
		} else {
			err = utils.SSHExec(server.PublicAddress.IP, args[1:], false)
		}
		if err != nil {
			log.Debugf("Command execution failed: %v", err)
			os.Exit(1)
		}
		log.Debugf("Command successfuly executed")
	}
}
