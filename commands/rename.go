package commands

import (
	log "github.com/Sirupsen/logrus"

	"github.com/scaleway/scaleway-cli/api"
	types "github.com/scaleway/scaleway-cli/commands/types"
)

var cmdRename = &types.Command{
	Exec:        runRename,
	UsageLine:   "rename [OPTIONS] SERVER NEW_NAME",
	Description: "Rename a server",
	Help:        "Rename a server.",
}

func init() {
	cmdRename.Flag.BoolVar(&renameHelp, []string{"h", "-help"}, false, "Print usage")
}

// Flags
var renameHelp bool // -h, --help flag

func runRename(cmd *types.Command, args []string) {
	if renameHelp {
		cmd.PrintUsage()
	}
	if len(args) != 2 {
		cmd.PrintShortUsage()
	}

	serverID := cmd.API.GetServerID(args[0])

	var server api.ScalewayServerPatchDefinition
	server.Name = &args[1]

	err := cmd.API.PatchServer(serverID, server)
	if err != nil {
		log.Fatalf("Cannot rename server: %v", err)
	} else {
		cmd.API.Cache.InsertServer(serverID, *server.Name)
	}
}
