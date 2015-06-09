package main

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
)

var cmdStop = &Command{
	Exec:        runStop,
	UsageLine:   "stop [OPTIONS] SERVER [SERVER...]",
	Description: "Stop a running server",
	Help:        "Stop a running server.",
	Examples: `
    $ scw stop my-running-server my-second-running-server
    $ scw stop -t my-running-server my-second-running-server
    $ scw stop $(scw ps -q)
    $ scw stop $(scw ps | grep mysql | awk '{print $1}')
`,
}

func init() {
	cmdStop.Flag.BoolVar(&stopT, []string{"t", "-terminate"}, false, "Stop and trash a server with its volumes")
	cmdStop.Flag.BoolVar(&stopHelp, []string{"h", "-help"}, false, "Print usage")
}

// Flags
var stopT bool    // -t flag
var stopHelp bool // -h, --help flag

func runStop(cmd *Command, args []string) {
	if stopHelp {
		cmd.PrintUsage()
	}
	if len(args) < 1 {
		cmd.PrintShortUsage()
	}

	hasError := false
	for _, needle := range args {
		server := cmd.API.GetServerID(needle)
		action := "poweroff"
		if stopT {
			action = "terminate"
		}
		err := cmd.API.PostServerAction(server, action)
		if err != nil {
			if err.Error() != "server should be running" && err.Error() != "server is being stopped or rebooted" {
				log.Warningf("failed to stop server %s: %s", server, err)
				hasError = true
			}
		} else {
			fmt.Println(needle)
		}
		if hasError {
			os.Exit(1)
		}
	}
}
