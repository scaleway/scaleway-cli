// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package cli

import "github.com/scaleway/scaleway-cli/pkg/commands"

var cmdExec = &Command{
	Exec:        runExec,
	UsageLine:   "exec [OPTIONS] SERVER [COMMAND] [ARGS...]",
	Description: "Run a command on a running server",
	Help:        "Run a command on a running server.",
	Examples: `
    $ scw exec myserver
    $ scw exec myserver bash
    $ scw exec --gateway=myotherserver myserver bash
    $ scw exec myserver 'tmux a -t joe || tmux new -s joe || bash'
    $ SCW_SECURE_EXEC=1 scw exec myserver bash
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
	cmdExec.Flag.StringVar(&execSSHUser, []string{"-user"}, "root", "Specify SSH user")
	cmdExec.Flag.IntVar(&execSSHPort, []string{"p", "-port"}, 22, "Specify SSH port")
	cmdExec.Flag.BoolVar(&execEnableSSHKeyForwarding, []string{"A"}, false, "Enable SSH keys forwarding")
}

// Flags
var execW bool                      // -w, --wait flag
var execTimeout float64             // -T flag
var execHelp bool                   // -h, --help flag
var execGateway string              // -g, --gateway flag
var execSSHUser string              // --user flag
var execSSHPort int                 // -p, --port flag
var execEnableSSHKeyForwarding bool // -A flag

func runExec(cmd *Command, rawArgs []string) error {
	if execHelp {
		return cmd.PrintUsage()
	}
	if len(rawArgs) < 1 {
		return cmd.PrintShortUsage()
	}

	args := commands.ExecArgs{
		Timeout:                execTimeout,
		Wait:                   execW,
		Gateway:                execGateway,
		Server:                 rawArgs[0],
		Command:                rawArgs[1:],
		SSHUser:                execSSHUser,
		SSHPort:                execSSHPort,
		EnableSSHKeyForwarding: execEnableSSHKeyForwarding,
	}
	ctx := cmd.GetContext(rawArgs)
	return commands.RunExec(ctx, args)
}
