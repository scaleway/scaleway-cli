// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package commands

import (
	"os"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/scaleway/scaleway-cli/api"
	types "github.com/scaleway/scaleway-cli/commands/types"
	"github.com/scaleway/scaleway-cli/utils"
)

var cmdExec = &types.Command{
	Exec:        runExec,
	UsageLine:   "exec [OPTIONS] SERVER [COMMAND] [ARGS...]",
	Description: "Run a command on a running server",
	Help:        "Run a command on a running server.",
	Examples: `
    $ scw exec myserver
    $ scw exec myserver bash
    $ scw exec --gateway=myotherserver myserver bash
    $ scw exec myserver 'tmux a -t joe || tmux new -s joe || bash'
    $ exec_secure=1 scw exec myserver bash
    $ scw exec -w $(scw start $(scw create ubuntu-trusty)) bash
    $ scw exec $(scw start -w $(scw create ubuntu-trusty)) bash
    $ scw exec myserver tmux new -d sleep 10
    $ scw exec myserver ls -la | grep password
    $ cat local-file | scw exec myserver 'cat > remote/path'
`,
}

func init() {
	cmdExec.Flag.BoolVar(&execHelp, []string{"h", "-help"}, false, "Print usage")
	cmdExec.Flag.Float64Var(&execTimeout, []string{"T", "-timeout"}, 0, "Set timeout values to seconds")
	cmdExec.Flag.BoolVar(&execW, []string{"w", "-wait"}, false, "Wait for SSH to be ready")
	cmdExec.Flag.StringVar(&execGateway, []string{"g", "-gateway"}, "", "Use a SSH gateway")
}

// Flags
var execW bool          // -w, --wait flag
var execTimeout float64 // -T flag
var execHelp bool       // -h, --help flag
var execGateway string  // -g, --gateway flag

func runExec(cmd *types.Command, args []string) {
	if execHelp {
		cmd.PrintUsage()
	}
	if len(args) < 1 {
		cmd.PrintShortUsage()
	}

	serverID := cmd.API.GetServerID(args[0])

	// Resolve gateway
	if execGateway == "" {
		execGateway = os.Getenv("SCW_GATEWAY")
	}
	var gateway string
	var err error
	if execGateway == serverID || execGateway == args[0] {
		gateway = ""
	} else {
		gateway, err = api.ResolveGateway(cmd.API, execGateway)
		if err != nil {
			log.Fatalf("Cannot resolve Gateway '%s': %v", execGateway, err)
		}
	}

	var server *api.ScalewayServer
	if execW {
		// --wait
		server, err = api.WaitForServerReady(cmd.API, serverID, gateway)
		if err != nil {
			log.Fatalf("Failed to wait for server to be ready, %v", err)
		}
	} else {
		// no --wait
		server, err = cmd.API.GetServer(serverID)
		if err != nil {
			log.Fatalf("Failed to get server information for %s: %v", serverID, err)
		}
	}

	if execTimeout > 0 {
		go func() {
			time.Sleep(time.Duration(execTimeout*1000) * time.Millisecond)
			log.Fatalf("Operation timed out")
		}()
	}

	err = utils.SSHExec(server.PublicAddress.IP, server.PrivateIP, args[1:], !execW, gateway)
	if err != nil {
		log.Fatalf("%v", err)
		os.Exit(1)
	}
	log.Debugf("Command successfuly executed")
}
