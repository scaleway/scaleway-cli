// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package cli

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
)

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
		for _, container := range containers.Containers {
			fmt.Fprintf(cmd.Streams().Stdout, "s3://%s\n", container.Name)
		}
		return nil
	}
	container := strings.Replace(args[0], "s3://", "", 1)
	datas, err := cmd.API.GetContainerDatas(container)
	if err != nil {
		return fmt.Errorf("Unable to get your data from %s: %v", container, err)
	}
	for _, data := range datas.Container {
		t, err := time.Parse(time.RFC3339, data.LastModified)
		if err != nil {
			return err
		}
		year, month, day := t.Date()
		hour, minute, _ := t.Clock()
		size, err := strconv.Atoi(data.Size)
		if err != nil {
			return err
		}
		fmt.Fprintf(cmd.Streams().Stdout, "%-4d-%02d-%02d %02d:%02d %8s s3://%s/%s\n", year, month, day, hour, minute, humanize.Bytes(uint64(size)), container, data.Name)
	}
	return nil
}
