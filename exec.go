package main

import (
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
)

var cmdExec = &Command{
	Exec:        runExec,
	UsageLine:   "exec [OPTIONS] SERVER COMMAND [ARGS...]",
	Description: "Run a command on a running server",
	Help:        "Run a command on a running server.",
	Examples: `
    $ scw exec myserver bash
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
}

// Flags
var execW bool          // -w, --wait flag
var execTimeout float64 // -T flag
var execHelp bool       // -h, --help flag

func runExec(cmd *Command, args []string) {
	if execHelp {
		cmd.PrintUsage()
	}
	if len(args) < 2 {
		cmd.PrintShortUsage()
	}

	serverID := cmd.API.GetServerID(args[0])

	var server *ScalewayServer
	var err error
	if execW {
		// --wait
		server, err = WaitForServerReady(cmd.API, serverID)
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

	err = sshExec(server.PublicAddress.IP, args[1:], !execW)
	if err != nil {
		log.Fatalf("%v", err)
		os.Exit(1)
	}
	log.Debugf("Command successfuly executed")
}
