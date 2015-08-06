// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package commands

import (
	"fmt"
	"time"

	log "github.com/scaleway/scaleway-cli/vendor/github.com/Sirupsen/logrus"

	"github.com/scaleway/scaleway-cli/pkg/api"
	"github.com/scaleway/scaleway-cli/pkg/commands/types"
	"github.com/scaleway/scaleway-cli/pkg/utils"
)

var cmdExec = &types.Command{
	Exec:        cmdExecExec,
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

// ExecArgs are flags for the `RunExec` function
type ExecArgs struct {
	Timeout float64
	Wait    bool
	Gateway string
	Server  string
	Command []string
}

func cmdExecExec(cmd *types.Command, rawArgs []string) {
	if execHelp {
		cmd.PrintUsage()
	}
	if len(rawArgs) < 1 {
		cmd.PrintShortUsage()
	}

	args := ExecArgs{
		Timeout: execTimeout,
		Wait:    execW,
		Gateway: execGateway,
		Server:  rawArgs[0],
		Command: rawArgs[1:],
	}
	ctx := cmd.GetContext(rawArgs)
	err := RunExec(ctx, args)
	if err != nil {
		log.Fatalf("Cannot exec 'exec': %v", err)
	}
}

// RunExec is the handler for 'scw exec'
func RunExec(ctx types.CommandContext, args ExecArgs) error {
	serverID := ctx.API.GetServerID(args.Server)

	// Resolve gateway
	if args.Gateway == "" {
		args.Gateway = ctx.Getenv("SCW_GATEWAY")
	}
	var gateway string
	var err error
	if args.Gateway == serverID || args.Gateway == args.Server {
		log.Debugf("The server and the gateway are the same host, using direct access to the server")
		gateway = ""
	} else {
		gateway, err = api.ResolveGateway(ctx.API, args.Gateway)
		if err != nil {
			return fmt.Errorf("Cannot resolve Gateway '%s': %v", args.Gateway, err)
		}
		if gateway != "" {
			log.Debugf("The server will be accessed using the gateway '%s' as a SSH relay", gateway)
		}
	}

	var server *api.ScalewayServer
	if args.Wait {
		// --wait
		log.Debugf("Waiting for server to be ready")
		server, err = api.WaitForServerReady(ctx.API, serverID, gateway)
		if err != nil {
			return fmt.Errorf("Failed to wait for server to be ready, %v", err)
		}
	} else {
		// no --wait
		log.Debugf("scw won't wait for the server to be ready, if it is not, the command will fail")
		server, err = ctx.API.GetServer(serverID)
		if err != nil {
			return fmt.Errorf("Failed to get server information for %s: %v", serverID, err)
		}
	}

	// --timeout
	if args.Timeout > 0 {
		log.Debugf("Setting up a global timeout of %d seconds", args.Timeout)
		// FIXME: avoid use of log.Fatalf here
		go func() {
			time.Sleep(time.Duration(args.Timeout*1000) * time.Millisecond)
			log.Fatalf("Operation timed out")
		}()
	}

	err = utils.SSHExec(server.PublicAddress.IP, server.PrivateIP, args.Command, !args.Wait, gateway)
	if err != nil {
		return fmt.Errorf("Failed to run the command: %v", err)
	}

	log.Debugf("Command successfuly executed")
	return nil
}
