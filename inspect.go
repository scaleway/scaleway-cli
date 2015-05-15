package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

var cmdInspect = &Command{
	Exec:        runInspect,
	UsageLine:   "inspect [OPTIONS] IDENTIFIER [IDENTIFIER...]",
	Description: "Inspects servers, images, snapshots and bootscripts.",
	Help:        "Inspects servers, images, snapshots and bootscripts.",
}

type ScalewayResolvedIdentifier struct {
	// Identifiers holds matching identifiers
	Identifiers []ScalewayIdentifier

	// Needle is the criteria used to lookup identifiers
	Needle string
}

// resolveIdentifiers resolves needles provided by the user
func resolveIdentifiers(cmd *Command, needles []string, out chan ScalewayResolvedIdentifier) {
	// first attempt, only lookup from the cache
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
	// fill the cache by fetching from the API and resolve missing identifiers
	if len(unresolved) > 0 {
		var wg sync.WaitGroup
		wg.Add(3)
		go func() {
			cmd.API.GetServers(true, 0)
			wg.Done()
		}()
		go func() {
			cmd.API.GetImages()
			wg.Done()
		}()
		go func() {
			cmd.API.GetSnapshots()
			wg.Done()
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
	close(out)
}

// inspectIdentifiers inspects identifiers concurrently
func inspectIdentifiers(cmd *Command, ci chan ScalewayResolvedIdentifier, cj chan interface{}) {
	var wg sync.WaitGroup
	for {
		idents, ok := <-ci
		if !ok {
			break
		}
		if len(idents.Identifiers) != 1 {
			if len(idents.Identifiers) == 0 {
				fmt.Fprintf(os.Stderr, "Unable to resolve identifier %s\n", idents.Needle)
			} else {
				fmt.Fprintf(os.Stderr, "Too many candidates for %s (%d)\n", idents.Needle, len(idents.Identifiers))
				for _, identifier := range idents.Identifiers {
					// FIXME: also print the name
					fmt.Fprintf(os.Stderr, "%s\n", identifier.Identifier)
				}
			}
		} else {
			ident := idents.Identifiers[0]
			wg.Add(1)
			go func() {
				if ident.Type == IDENTIFIER_SERVER {
					server, err := cmd.API.GetServer(ident.Identifier)
					if err == nil {
						cj <- server
					}
				} else if ident.Type == IDENTIFIER_IMAGE {
					image, err := cmd.API.GetImage(ident.Identifier)
					if err == nil {
						cj <- image
					}
				} else if ident.Type == IDENTIFIER_SNAPSHOT {
					snap, err := cmd.API.GetSnapshot(ident.Identifier)
					if err == nil {
						cj <- snap
					}
				} else if ident.Type == IDENTIFIER_BOOTSCRIPT {
					// FIXME: bootscript
				}
				wg.Done()
			}()
		}
	}
	wg.Wait()
	close(cj)
}

func runInspect(cmd *Command, args []string) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "usage: scw %s\n", cmd.UsageLine)
		os.Exit(1)
	}

	res := "["
	nb_inspected := 0
	ci := make(chan ScalewayResolvedIdentifier)
	cj := make(chan interface{})
	go resolveIdentifiers(cmd, args, ci)
	go inspectIdentifiers(cmd, ci, cj)
	for {
		data, open := <-cj
		if !open {
			break
		}
		data_b, err := json.MarshalIndent(data, "", "  ")
		if err == nil {
			if nb_inspected != 0 {
				res += ",\n"
			}
			res += string(data_b)
			nb_inspected += 1
		}
	}
	res += "]"

	fmt.Fprintf(os.Stdout, "%s\n", res)

	if len(args) != nb_inspected {
		os.Exit(1)
	}
}
