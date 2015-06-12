package commands

import (
	"fmt"
	"log"
	"sort"
	"strings"

	types "github.com/scaleway/scaleway-cli/commands/types"
	utils "github.com/scaleway/scaleway-cli/utils"
)

var cmdCompletion = &types.Command{
	Exec:        runCompletion,
	UsageLine:   "_completion [OPTIONS] CATEGORY",
	Description: "Completion helper",
	Help:        "Completion helper.",
	Hidden:      true,
	Examples: `
    $ scw _completion servers-all
    $ scw _completion images-all
    $ scw _completion snapshots-all
    $ scw _completion volumes-all
    $ scw _completion bootscripts-all
`,
}

func init() {
	cmdCompletion.Flag.BoolVar(&completionHelp, []string{"h", "-help"}, false, "Print usage")
}

// Flags
var completionHelp bool // -h, --help flag

func runCompletion(cmd *types.Command, args []string) {
	if completionHelp {
		cmd.PrintUsage()
	}
	if len(args) != 1 {
		cmd.PrintShortUsage()
	}

	category := args[0]

	elements := []string{}

	switch category {
	case "servers-all":
		for identifier, name := range cmd.API.Cache.Servers {
			elements = append(elements, identifier, utils.Wordify(name))
		}
	case "images-all":
		for identifier, name := range cmd.API.Cache.Images {
			elements = append(elements, identifier, utils.Wordify(name))
		}
	case "volumes-all":
		for identifier, name := range cmd.API.Cache.Volumes {
			elements = append(elements, identifier, utils.Wordify(name))
		}
	case "snapshots-all":
		for identifier, name := range cmd.API.Cache.Snapshots {
			elements = append(elements, identifier, utils.Wordify(name))
		}
	case "bootscripts-all":
		for identifier, name := range cmd.API.Cache.Bootscripts {
			elements = append(elements, identifier, utils.Wordify(name))
		}
	default:
		log.Fatalf("Unhandled category of completion: %s", category)
	}

	sort.Strings(elements)
	fmt.Println(strings.Join(utils.RemoveDuplicates(elements), "\n"))
}
