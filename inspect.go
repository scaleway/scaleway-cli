package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	log "github.com/Sirupsen/logrus"
)

var cmdInspect = &Command{
	Exec:        runInspect,
	UsageLine:   "inspect [OPTIONS] IDENTIFIER [IDENTIFIER...]",
	Description: "Return low-level information on a server, image, snapshot or bootscript",
	Help:        "Return low-level information on a server, image, snapshot or bootscript.",
	Examples: `
    $ scw inspect a-public-image
    $ scw inspect my-snapshot
    $ scw inspect my-image
    $ scw inspect my-server
    $ scw inspect my-volume
    $ scw inspect my-server | jq '.[0].public_ip.address'
    $ scw inspect $(scw inspect my-image | jq '.[0].root_volume.id')
`,
}

type ScalewayResolvedIdentifier struct {
	// Identifiers holds matching identifiers
	Identifiers []ScalewayIdentifier

	// Needle is the criteria used to lookup identifiers
	Needle string
}

func init() {
	cmdInspect.Flag.BoolVar(&inspectHelp, []string{"h", "-help"}, false, "Print usage")
}

// Flags
var inspectHelp bool // -h, --help flag

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
		wg.Add(4)
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
		go func() {
			cmd.API.GetBootscripts()
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
				log.Errorf("Unable to resolve identifier %s", idents.Needle)
			} else {
				log.Errorf("Too many candidates for %s (%d)", idents.Needle, len(idents.Identifiers))
				for _, identifier := range idents.Identifiers {
					// FIXME: also print the name
					log.Infof("- %s", identifier.Identifier)
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
					bootscript, err := cmd.API.GetBootscript(ident.Identifier)
					if err == nil {
						cj <- bootscript
					}
				}
				wg.Done()
			}()
		}
	}
	wg.Wait()
	close(cj)
}

func runInspect(cmd *Command, args []string) {
	if inspectHelp {
		cmd.PrintUsage()
	}
	if len(args) < 1 {
		cmd.PrintShortUsage()
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

	fmt.Println(res)

	if len(args) != nb_inspected {
		os.Exit(1)
	}
}
