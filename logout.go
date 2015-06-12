package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
)

var cmdLogout = &Command{
	Exec:        runLogout,
	UsageLine:   "logout [OPTIONS]",
	Description: "Log out from the Scaleway API",
	Help:        "Log out from the Scaleway API.",
}

func init() {
	cmdLogout.Flag.BoolVar(&logoutHelp, []string{"h", "-help"}, false, "Print usage")
}

// FLags
var logoutHelp bool // -h, --help flag

func runLogout(cmd *Command, args []string) {
	if logoutHelp {
		cmd.PrintUsage()
	}
	if len(args) != 0 {
		cmd.PrintShortUsage()
	}

	scwrcPath, err := GetConfigFilePath()
	if err != nil {
		log.Fatalf("Unable to get scwrc config file path: %v", err)
	}

	if _, err = os.Stat(scwrcPath); err == nil {
		err = os.Remove(scwrcPath)
		if err != nil {
			log.Fatalf("Unable to remove scwrc config file: %v", err)
		}
	}
}
