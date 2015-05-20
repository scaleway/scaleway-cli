package main

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
)

var cmdTag = &Command{
	Exec:        runTag,
	UsageLine:   "tag [OPTIONS] SNAPSHOT NAME",
	Description: "Tag a snapshot into an image",
	Help:        "Tag a snapshot into an image.",
}

func runTag(cmd *Command, args []string) {
	if len(args) < 2 {
		log.Fatalf("usage: scw %s", cmd.UsageLine)
	}

	snapshotId := cmd.GetSnapshot(args[0])
	snapshot, err := cmd.API.GetSnapshot(snapshotId)
	if err != nil {
		log.Fatalf("Cannot fetch snapshot: %v", err)
	}

	image, err := cmd.API.PostImage(snapshot.Identifier, args[1])
	if err != nil {
		log.Fatalf("Cannot create image: %v", err)
	}
	fmt.Println(image)
}
