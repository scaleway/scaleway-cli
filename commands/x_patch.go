package commands

import (
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"

	api "github.com/scaleway/scaleway-cli/api"
	types "github.com/scaleway/scaleway-cli/commands/types"
)

var cmdPatch = &types.Command{
	Exec:        runPatch,
	UsageLine:   "_patch [OPTIONS] IDENTIFIER FIELD=VALUE",
	Description: "",
	Hidden:      true,
	Help:        "PATCH an object on the API",
	Examples: `
    $ scw _patch myserver state_detail=booted
`,
}

func init() {
	cmdPatch.Flag.BoolVar(&patchHelp, []string{"h", "-help"}, false, "Print usage")
}

// Flags
var patchHelp bool // -h, --help flag

func runPatch(cmd *types.Command, args []string) {
	if patchHelp {
		cmd.PrintUsage()
	}
	if len(args) != 2 {
		cmd.PrintShortUsage()
	}

	// Parsing FIELD=VALUE
	updateParts := strings.Split(args[1], "=")
	if len(updateParts) != 2 {
		cmd.PrintShortUsage()
	}
	fieldName := updateParts[0]
	newValue := updateParts[1]

	ident := api.GetIdentifier(cmd.API, args[0])
	switch ident.Type {
	case api.IdentifierServer:
		var payload api.ScalewayServerPatchDefinition

		switch fieldName {
		case "state_detail":
			payload.StateDetail = &newValue
		case "name":
			payload.Name = &newValue
			log.Warnf("Use 'scw rename instead'")
		default:
			log.Fatalf("'_patch server %s=' not implemented", fieldName)
		}

		err := cmd.API.PatchServer(ident.Identifier, payload)
		if err != nil {
			log.Fatalf("Cannot rename server: %v", err)
		}
	default:
		log.Fatalf("_patch not implemented for this kind of object")
	}
	fmt.Println(ident.Identifier)
}
