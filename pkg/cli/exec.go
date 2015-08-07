// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package cli

import (
	"github.com/scaleway/scaleway-cli/pkg/commands"
	"github.com/scaleway/scaleway-cli/vendor/github.com/Sirupsen/logrus"
)

var cmdExec = &Command{
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

func cmdExecExec(cmd *Command, rawArgs []string) {
	if execHelp {
		cmd.PrintUsage()
	}
	if len(rawArgs) < 1 {
		cmd.PrintShortUsage()
	}

	args := commands.ExecArgs{
		Timeout: execTimeout,
		Wait:    execW,
		Gateway: execGateway,
		Server:  rawArgs[0],
		Command: rawArgs[1:],
	}
	ctx := cmd.GetContext(rawArgs)
	err := commands.RunExec(ctx, args)
	if err != nil {
		logrus.Fatalf("Cannot exec 'exec': %v", err)
	}
}
