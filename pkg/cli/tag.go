// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package cli

import (
	"github.com/scaleway/scaleway-cli/pkg/commands"
	"github.com/scaleway/scaleway-cli/vendor/github.com/Sirupsen/logrus"
)

var cmdTag = &Command{
	Exec:        runTag,
	UsageLine:   "tag [OPTIONS] SNAPSHOT NAME",
	Description: "Tag a snapshot into an image",
	Help:        "Tag a snapshot into an image.",
}

func init() {
	cmdTag.Flag.BoolVar(&tagHelp, []string{"h", "-help"}, false, "Print usage")
}

// Flags
var tagHelp bool // -h, --help flag

func runTag(cmd *Command, rawArgs []string) {
	if tagHelp {
		cmd.PrintUsage()
	}
	if len(rawArgs) != 2 {
		cmd.PrintShortUsage()
	}

	args := commands.TagArgs{
		Snapshot: rawArgs[0],
		Name:     rawArgs[1],
	}
	ctx := cmd.GetContext(rawArgs)
	err := commands.RunTag(ctx, args)
	if err != nil {
		logrus.Fatalf("Cannot execute 'tag': %v", err)
	}
}
