package main

import (
	"fmt"
	"os"
	"sync"
)

var cmdInspect = &Command{
	Exec:        runInspect,
	UsageLine:   "inspect [OPTIONS] IDENTIFIER [IDENTIFIER...]",
	Description: "Inspects a server, image or a bootscript.",
	Help:        "Inspects a server, image or a bootscript.",
}

type ScalewayResolvedIdentifier struct {
	// Identifiers holds matching identifiers
	Identifiers []ScalewayIdentifier

	// Needle is the criteria used to lookup identifiers
	Needle string
}

// resolveIdentifiers resolves incoming identifiers
func resolveIdentifiers(cmd *Command, needles []string, out chan ScalewayResolvedIdentifier) {
	var unresolved []string
	for _, needle := range needles {
		idents := cmd.API.Cache.LookUpIdentifiers(needle)
		if len(idents) == 0 {
			unresolved = append(unresolved, needle)
		} else {
			out <- ScalewayResolvedIdentifier{
				Identifiers: idents,
				Needle:      needle,
			}
		}
	}
	if len(unresolved) > 0 {
		var wg sync.WaitGroup
		go func() {
			wg.Add(1)
			cmd.API.GetServers(true, 0)
		}()
		go func() {
			wg.Add(1)
			cmd.API.GetImages()
		}()
		wg.Wait()
		for _, needle := range unresolved {
			idents := cmd.API.Cache.LookUpIdentifiers(needle)
			out <- ScalewayResolvedIdentifier{
				Identifiers: idents,
				Needle:      needle,
			}
		}
	}
}

func runInspect(cmd *Command, args []string) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "usage: scw %s\n", cmd.UsageLine)
		os.Exit(1)
	}

	has_error := false

	if has_error {
		os.Exit(1)
	}
}
