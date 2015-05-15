package main

import (
	"fmt"
	"os"
)

var cmdRmi = &Command{
	Exec:        runRmi,
	UsageLine:   "rmi [OPTIONS] IMAGE [IMAGE...]",
	Description: "Remove one or more images",
	Help:        "Remove one or more images.",
}

func runRmi(cmd *Command, args []string) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "usage: scw %s\n", cmd.UsageLine)
		os.Exit(1)
	}
	has_error := false
	for _, needle := range args {
		image := cmd.GetImage(needle)
		err := cmd.API.DeleteImage(image)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to delete image %s: %s\n", image, err)
			has_error = true
		} else {
			fmt.Fprintf(os.Stdout, "%s\n", needle)
		}
	}
	if has_error {
		os.Exit(1)
	}
}
