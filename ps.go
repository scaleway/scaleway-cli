package main

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
	cmdPs.Flag.BoolVar(&psNoTrunc, "-no-trunc", false, "don't truncate output")
	cmdPs.Flag.IntVar(&psN, "n", 0, "show n last created servers, include non-running ones")
	cmdPs.Flag.BoolVar(&psQ, "q", false, "only display numeric IDs")
}

// Flags
var psA bool       // -a flag
var psL bool       // -l flag
var psQ bool       // -q flag
var psNoTrunc bool // --no-trunc flag
var psN int        // -n flag

func runPs(cmd *Command, args []string) {
}
