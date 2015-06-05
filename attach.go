package main

import (
	"fmt"
	"os"
	"os/exec"

	log "github.com/Sirupsen/logrus"
)

var cmdAttach = &Command{
	Exec:        runAttach,
	UsageLine:   "attach [OPTIONS] SERVER",
	Description: "Attach to a server serial console",
	Help:        "Attach to a running server serial console.",
	Examples: `
    $ scw attach my-running-server
    $ scw attach $(scw start my-stopped-server)
    $ scw attach $(scw start $(scw create ubuntu-vivid))
`,
}

func init() {
	cmdAttach.Flag.BoolVar(&attachHelp, []string{"h", "-help"}, false, "Print usage")
}

// Flags
var attachHelp bool // -h, --help flag

const termjsBin string = "termjs-cli"

func runAttach(cmd *Command, args []string) {
	if attachHelp {
		cmd.PrintUsage()
	}
	if len(args) != 1 {
		cmd.PrintShortUsage()
	}

	serverID := cmd.GetServer(args[0])

	termjsURL := fmt.Sprintf("https://tty.cloud.online.net?server_id=%s&type=serial&auth_token=%s", serverID, cmd.API.Token)

	log.Debugf("Executing: %s %s", termjsBin, termjsURL)
	// FIXME: check if termjs-cli is installed
	spawn := exec.Command(termjsBin, termjsURL)
	spawn.Stdout = os.Stdout
	spawn.Stdin = os.Stdin
	spawn.Stderr = os.Stderr
	err := spawn.Run()
	if err != nil {
		log.Warnf("%v", err)
		fmt.Fprintf(os.Stderr, `
You need to install '%s' from https://github.com/moul/term.js-cli

    npm install -g term.js-cli

However, you can access your serial using a web browser:

    %s

`, termjsBin, termjsURL)
		os.Exit(1)
	}
}
