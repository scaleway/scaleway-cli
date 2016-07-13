// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package cli

var cmdIPS = &Command{
	Exec:      runIPS,
	UsageLine: "_ips [OPTIONS] [IP_ID [SERVER_ID]]",

	Description: "Interacts with your IPs",
	Hidden:      true,
	Help:        "Interacts with your IPs",
	Examples: `
    $ scw _ips
    $ scw _ips IP_ID
    $ scw _ips --new
    $ scw _ips --attach IP_ID SERVER_ID
    $ scw _ips --delete IP_ID
    $ scw _ips --detach IP_ID
`,
}

func init() {
	cmdIPS.Flag.BoolVar(&ipHelp, []string{"h", "-help"}, false, "Print usage")
	cmdIPS.Flag.BoolVar(&ipNew, []string{"n", "-new"}, false, "Add a new IP")
	cmdIPS.Flag.BoolVar(&ipAttach, []string{"a", "-attach"}, false, "Attach an IP to a server")
	cmdIPS.Flag.BoolVar(&ipDetach, []string{"-detach"}, false, "Detach an IP from a server")
	cmdIPS.Flag.StringVar(&ipDelete, []string{"d", "-delete"}, "", "Detele an IP")
}

var ipHelp bool     // -h, --help flag
var ipNew bool      // -n, --new flag
var ipAttach bool   // -a, --attach flag
var ipDetach bool   // --detach flag
var ipDelete string // -d, --delete flag

func runIPS(cmd *Command, args []string) error {
	if ipHelp {
		return cmd.PrintUsage()
	}
	if ipNew {
		ip, err := cmd.API.NewIP()
		if err != nil {
			return err
		}
		printRawMode(cmd.Streams().Stdout, ip)
		return nil
	}
	if ipDelete != "" {
		return cmd.API.DeleteIP(ipDelete)
	}
	if ipAttach {
		if len(args) != 2 {
			return cmd.PrintShortUsage()
		}
		return cmd.API.AttachIP(args[0], args[1])
	}
	if ipDetach {
		return cmd.API.DetachIP(args[0])
	}
	if len(args) == 1 {
		ip, err := cmd.API.GetIP(args[0])
		if err != nil {
			return err
		}
		printRawMode(cmd.Streams().Stdout, *ip)
		return nil
	}
	ips, err := cmd.API.GetIPS()
	if err != nil {
		return err
	}
	printRawMode(cmd.Streams().Stdout, *ips)
	return nil
}
