package main

import log "github.com/Sirupsen/logrus"

var cmdRename = &Command{
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

func runRename(cmd *Command, args []string) {
	if renameHelp {
		cmd.PrintUsage()
	}
	if len(args) != 2 {
		cmd.PrintShortUsage()
	}

	serverID := cmd.API.GetServerID(args[0])

	var server ScalewayServerPatchNameDefinition
	server.Name = args[1]

	err := cmd.API.PatchServerName(serverID, server)
	if err != nil {
		log.Fatalf("Cannot rename server: %v", err)
	} else {
		cmd.API.Cache.InsertServer(serverID, server.Name)
	}
}
