// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package cli

import "github.com/scaleway/scaleway-cli/pkg/commands"

var cmdTag = &Command{
	Exec:        runTag,
	UsageLine:   "tag [OPTIONS] SNAPSHOT NAME",
	Description: "Tag a snapshot into an image",
	Help:        "Tag a snapshot into an image.",
}

func init() {
	cmdTag.Flag.BoolVar(&tagHelp, []string{"h", "-help"}, false, "Print usage")
	cmdTag.Flag.StringVar(&tagBootscript, []string{"-bootscript"}, "", "Assign bootscript")
	cmdTag.Flag.StringVar(&tagArch, []string{"-arch"}, "arm", "Image architecture arm, x86_64")
}

// Flags
var tagHelp bool         // -h, --help flag
var tagBootscript string // --bootscript flag
var tagArch string       // --arch flag

func runTag(cmd *Command, rawArgs []string) error {
	if tagHelp {
		return cmd.PrintUsage()
	}
	if len(rawArgs) != 2 {
		return cmd.PrintShortUsage()
	}

	args := commands.TagArgs{
		Snapshot:   rawArgs[0],
		Name:       rawArgs[1],
		Bootscript: tagBootscript,
		Arch:       tagArch,
	}
	ctx := cmd.GetContext(rawArgs)
	return commands.RunTag(ctx, args)
}
