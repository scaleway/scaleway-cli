package main

import (
	"fmt"
	"os"
	"text/tabwriter"
)

var cmdPs = &Command{
	Exec:        runPs,
	UsageLine:   "ps [options]",
	Description: "Ps lists Scaleway servers",
	Help: `
Ps lists Scaleway servers. By default, only running servers are displayed.`,
}

func init() {
	cmdPs.Flag.BoolVar(&psA, "a", false, "show all servers. only running servers are shown by default")
	cmdPs.Flag.BoolVar(&psL, "l", false, "show only the latest created server, include non-running ones")
	cmdPs.Flag.BoolVar(&psNoTrunc, "no-trunc", false, "don't truncate output")
	cmdPs.Flag.IntVar(&psN, "n", 0, "show n last created servers, include non-running ones")
	cmdPs.Flag.BoolVar(&psQ, "q", false, "only display numeric IDs")
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
			fmt.Fprintf(w, "%s\n", truncIf(server.Identifier, 8, !psNoTrunc))
		} else {
			short_id := truncIf(server.Identifier, 8, !psNoTrunc)
			short_image := truncIf(server.Image.Name, 20, !psNoTrunc)
			short_name := truncIf(server.Name, 20, !psNoTrunc)
			fmt.Fprintf(w, "%s\t%s\t\t%s\t%s\t\t%s\n", short_id, short_image, server.CreationDate, server.State, short_name)
		}

		if psL {
			break
		}
	}
}
