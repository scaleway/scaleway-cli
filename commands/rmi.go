package commands

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"

	types "github.com/scaleway/scaleway-cli/commands/types"
)

var cmdRmi = &types.Command{
	Exec:        runRmi,
	UsageLine:   "rmi [OPTIONS] IMAGE [IMAGE...]",
	Description: "Remove one or more images",
	Help:        "Remove one or more images.",
	Examples: `
    $ scw rmi myimage
    $ scw rmi $(scw images -q)
`,
}

func init() {
	cmdRmi.Flag.BoolVar(&rmiHelp, []string{"h", "-help"}, false, "Print usage")
}

// Flags
var rmiHelp bool // -h, --help flag

func runRmi(cmd *types.Command, args []string) {
	if rmiHelp {
		cmd.PrintUsage()
	}
	if len(args) < 1 {
		cmd.PrintShortUsage()
	}

	if len(args) == 0 {
		log.Fatalf("usage: scw %s", cmd.UsageLine)
	}
	hasError := false
	for _, needle := range args {
		image := cmd.API.GetImageID(needle, true)
		err := cmd.API.DeleteImage(image)
		if err != nil {
			log.Errorf("failed to delete image %s: %s", image, err)
			hasError = true
		} else {
			fmt.Println(needle)
		}
	}
	if hasError {
		os.Exit(1)
	}
}
