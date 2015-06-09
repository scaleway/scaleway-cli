package main

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
)

var cmdRm = &Command{
	Exec:        runRm,
	UsageLine:   "rm [OPTIONS] SERVER [SERVER...]",
	Description: "Remove one or more servers",
	Help:        "Remove one or more servers.",
	Examples: `
    $ scw rm my-stopped-server my-second-stopped-server
    $ scw rm $(scw ps -q)
    $ scw rm $(scw ps | grep mysql | awk '{print $1}')
`,
}

func init() {
	cmdRm.Flag.BoolVar(&rmHelp, []string{"h", "-help"}, false, "Print usage")
}

// Flags
var rmHelp bool // -h, --help flag

func runRm(cmd *Command, args []string) {
	if rmHelp {
		cmd.PrintUsage()
	}
	if len(args) < 1 {
		cmd.PrintShortUsage()
	}

	hasError := false
	for _, needle := range args {
		server := cmd.API.GetServerID(needle)
		err := cmd.API.DeleteServer(server)
		if err != nil {
			log.Errorf("failed to delete server %s: %s", server, err)
			hasError = true
		} else {
			fmt.Println(needle)
		}
	}
	if hasError {
		os.Exit(1)
	}
}
