// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package cli

import "fmt"

var cmdCS = &Command{
	Exec:        runCS,
	UsageLine:   "_cs [CONTAINER_NAME]",
	Description: "",
	Hidden:      true,
	Help:        "List containers / datas",
	Examples: `
    $ scw _cs
    $ scw _cs containerName
`,
}

func init() {
	cmdCS.Flag.BoolVar(&csHelp, []string{"h", "-help"}, false, "Print usage")
}

// Flags
var csHelp bool // -h, --help flag

func runCS(cmd *Command, args []string) error {
	if csHelp {
		return cmd.PrintUsage()
	}
	if len(args) > 1 {
		return cmd.PrintShortUsage()
	}
	if len(args) == 0 {
		containers, err := cmd.API.GetContainers()
		if err != nil {
			return fmt.Errorf("Unable to get your containers: %v", err)
		}
		printRawMode(cmd.Streams().Stdout, *containers)
		return nil
	}
	datas, err := cmd.API.GetContainerDatas(args[0])
	if err != nil {
		return fmt.Errorf("Unable to get your data from %s: %v", args[1], err)
	}
	printRawMode(cmd.Streams().Stdout, *datas)
	return nil
}
