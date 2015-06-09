package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"text/template"

	log "github.com/Sirupsen/logrus"
)

var cmdInspect = &Command{
	Exec:        runInspect,
	UsageLine:   "inspect [OPTIONS] IDENTIFIER [IDENTIFIER...]",
	Description: "Return low-level information on a server, image, snapshot, volume or bootscript",
	Help:        "Return low-level information on a server, image, snapshot, volume or bootscript.",
	Examples: `
    $ scw inspect a-public-image
    $ scw inspect my-snapshot
    $ scw inspect my-volume
    $ scw inspect my-image
    $ scw inspect my-server
    $ scw inspect my-volume
    $ scw inspect my-server | jq '.[0].public_ip.address'
    $ scw inspect $(scw inspect my-image | jq '.[0].root_volume.id')
    $ scw inspect -f "{{ .PublicAddress.IP }}" my-server
`,
}

// ScalewayResolvedIdentifier represents a list of matching identifier for a specifier pattern
type ScalewayResolvedIdentifier struct {
	// Identifiers holds matching identifiers
	Identifiers []ScalewayIdentifier

	// Needle is the criteria used to lookup identifiers
	Needle string
}

func init() {
	cmdInspect.Flag.BoolVar(&inspectHelp, []string{"h", "-help"}, false, "Print usage")
	cmdInspect.Flag.StringVar(&inspectFormat, []string{"f", "-format"}, "", "Format the output using the given go template.")
}

// Flags
var inspectFormat string // -f, --format flat
var inspectHelp bool     // -h, --help flag

// resolveIdentifiers resolves needles provided by the user
func resolveIdentifiers(api *ScalewayAPI, needles []string, out chan ScalewayResolvedIdentifier) {
	// first attempt, only lookup from the cache
	var unresolved []string
	for _, needle := range needles {
		idents := api.Cache.LookUpIdentifiers(needle)
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
		wg.Add(5)
		go func() {
			api.GetServers(true, 0)
			wg.Done()
		}()
		go func() {
			api.GetImages()
			wg.Done()
		}()
		go func() {
			api.GetSnapshots()
			wg.Done()
		}()
		go func() {
			api.GetVolumes()
			wg.Done()
		}()
		go func() {
			api.GetBootscripts()
			wg.Done()
		}()
		wg.Wait()
		for _, needle := range unresolved {
			idents := api.Cache.LookUpIdentifiers(needle)
			out <- ScalewayResolvedIdentifier{
				Identifiers: idents,
				Needle:      needle,
			}
		}
	}
	close(out)
}

// inspectIdentifiers inspects identifiers concurrently
func inspectIdentifiers(api *ScalewayAPI, ci chan ScalewayResolvedIdentifier, cj chan interface{}) {
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
				if ident.Type == IdentifierServer {
					server, err := api.GetServer(ident.Identifier)
					if err == nil {
						cj <- server
					}
				} else if ident.Type == IdentifierImage {
					image, err := api.GetImage(ident.Identifier)
					if err == nil {
						cj <- image
					}
				} else if ident.Type == IdentifierSnapshot {
					snap, err := api.GetSnapshot(ident.Identifier)
					if err == nil {
						cj <- snap
					}
				} else if ident.Type == IdentifierVolume {
					snap, err := api.GetVolume(ident.Identifier)
					if err == nil {
						cj <- snap
					}
				} else if ident.Type == IdentifierBootscript {
					bootscript, err := api.GetBootscript(ident.Identifier)
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

var funcMap = template.FuncMap{
	"json": func(v interface{}) string {
		a, _ := json.Marshal(v)
		return string(a)
	},
}

func runInspect(cmd *Command, args []string) {
	if inspectHelp {
		cmd.PrintUsage()
	}
	if len(args) < 1 {
		cmd.PrintShortUsage()
	}

	res := "["
	nbInspected := 0
	ci := make(chan ScalewayResolvedIdentifier)
	cj := make(chan interface{})
	go resolveIdentifiers(cmd.API, args, ci)
	go inspectIdentifiers(cmd.API, ci, cj)
	for {
		data, open := <-cj
		if !open {
			break
		}
		if inspectFormat == "" {
			dataB, err := json.MarshalIndent(data, "", "  ")
			if err == nil {
				if nbInspected != 0 {
					res += ",\n"
				}
				res += string(dataB)
				nbInspected++
			}
		} else {
			tmpl, err := template.New("").Funcs(funcMap).Parse(inspectFormat)
			if err != nil {
				log.Fatalf("Format parsing error: %v", err)
			}

			err = tmpl.Execute(os.Stdout, data)
			if err != nil {
				log.Fatalf("Format execution error: %v", err)
			}
			fmt.Fprint(os.Stdout, "\n")
		}
	}
	res += "]"

	if inspectFormat == "" {
		fmt.Println(res)
	}

	if len(args) != nbInspected {
		os.Exit(1)
	}
}
