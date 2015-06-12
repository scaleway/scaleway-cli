package commands

import (
	"fmt"

	log "github.com/Sirupsen/logrus"

	types "github.com/scaleway/scaleway-cli/commands/types"
)

var cmdTag = &types.Command{
	Exec:        runTag,
	UsageLine:   "tag [OPTIONS] SNAPSHOT NAME",
	Description: "Tag a snapshot into an image",
	Help:        "Tag a snapshot into an image.",
}

func init() {
	cmdTag.Flag.BoolVar(&tagHelp, []string{"h", "-help"}, false, "Print usage")
}

// Flags
var tagHelp bool // -h, --help flag

func runTag(cmd *types.Command, args []string) {
	if tagHelp {
		cmd.PrintUsage()
	}
	if len(args) != 2 {
		cmd.PrintShortUsage()
	}

	snapshotID := cmd.API.GetSnapshotID(args[0])
	snapshot, err := cmd.API.GetSnapshot(snapshotID)
	if err != nil {
		log.Fatalf("Cannot fetch snapshot: %v", err)
	}

	image, err := cmd.API.PostImage(snapshot.Identifier, args[1])
	if err != nil {
		log.Fatalf("Cannot create image: %v", err)
	}
	fmt.Println(image)
}
