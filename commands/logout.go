package commands

import (
	"os"

	log "github.com/Sirupsen/logrus"

	types "github.com/scaleway/scaleway-cli/commands/types"
	"github.com/scaleway/scaleway-cli/utils"
)

var cmdLogout = &types.Command{
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

func runLogout(cmd *types.Command, args []string) {
	if logoutHelp {
		cmd.PrintUsage()
	}
	if len(args) != 0 {
		cmd.PrintShortUsage()
	}

	scwrcPath, err := utils.GetConfigFilePath()
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
