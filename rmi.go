package main

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
)

var cmdRmi = &Command{
	Exec:        runRmi,
	UsageLine:   "rmi [OPTIONS] IMAGE [IMAGE...]",
	Description: "Remove one or more images",
	Help:        "Remove one or more images.",
}

func runRmi(cmd *Command, args []string) {
	if len(args) == 0 {
		log.Fatalf("usage: scw %s", cmd.UsageLine)
	}
	has_error := false
	for _, needle := range args {
		image := cmd.GetImage(needle)
		err := cmd.API.DeleteImage(image)
		if err != nil {
			log.Errorf("failed to delete image %s: %s", image, err)
			has_error = true
		} else {
			fmt.Println(needle)
		}
	}
	if has_error {
		os.Exit(1)
	}
}
