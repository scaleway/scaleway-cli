package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/docker/docker/pkg/units"
)

var cmdPs = &Command{
	Exec:        runPs,
	UsageLine:   "ps [OPTIONS]",
	Description: "List servers",
	Help:        "List servers. By default, only running servers are displayed.",
}

func init() {
	// FIXME: -h
	cmdPs.Flag.BoolVar(&psA, "a", false, "Show all servers. Only running servers are shown by default")
	cmdPs.Flag.BoolVar(&psL, "l", false, "Show only the latest created server, include non-running ones")
	cmdPs.Flag.IntVar(&psN, "n", 0, "Show n last created servers, include non-running ones")
	cmdPs.Flag.BoolVar(&psNoTrunc, "no-trunc", false, "Don't truncate output")
	cmdPs.Flag.BoolVar(&psQ, "q", false, "Only display numeric IDs")
}

// Flags
var psA bool       // -a flag
var psL bool       // -l flag
var psQ bool       // -q flag
var psNoTrunc bool // -no-trunc flag
var psN int        // -n flag

// truncIf ensures the input string does not exceed max size if cond is met
func truncIf(str string, max int, cond bool) string {
	if cond && len(str) > max {
		return str[:max]
	}
	return str
}

// wordify convert complex name to a single word without special shell characters
func wordify(str string) string {
	str = regexp.MustCompile(`[^a-zA-Z0-9-]`).ReplaceAllString(str, "_")
	str = regexp.MustCompile(`__+`).ReplaceAllString(str, "_")
	str = strings.Trim(str, "_")
	return str
}

func runPs(cmd *Command, args []string) {
	api, err := GetScalewayAPI()
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to init Scaleway API: %v\n", err)
		os.Exit(1)
	}
	servers, err := api.GetServers()
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to fetch servers from the Scaleway API: %v\n", err)
		os.Exit(1)
	}

	w := tabwriter.NewWriter(os.Stdout, 20, 1, 3, ' ', 0)
	defer w.Flush()
	if !psQ {
		fmt.Fprintf(w, "SERVER ID\tIMAGE\tCOMMAND\tCREATED\tSTATUS\tPORTS\tNAME\n")
	}
	for id, server := range *servers {
		if !(psA || (psN != 0 && id < psN) || server.State == "running") {
			continue
		}

		if psQ {
			fmt.Fprintf(w, "%s\n", server.Identifier)
		} else {
			short_id := truncIf(server.Identifier, 8, !psNoTrunc)
			short_image := truncIf(wordify(server.Image.Name), 25, !psNoTrunc)
			short_name := truncIf(wordify(server.Name), 25, !psNoTrunc)
			creationTime, _ := time.Parse("2006-01-02T15:04:05.000000+00:00", server.CreationDate)
			short_creationDate := units.HumanDuration(time.Now().UTC().Sub(creationTime))
			fmt.Fprintf(w, "%s\t%s\t\t%s\t%s\t\t%s\n", short_id, short_image, short_creationDate, server.State, short_name)
		}

		if psL {
			break
		}
	}
}
