// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package cli

import (
	"github.com/scaleway/scaleway-cli/pkg/commands"
	"github.com/scaleway/scaleway-cli/vendor/github.com/Sirupsen/logrus"
)

var cmdImages = &Command{
	Exec:        runImages,
	UsageLine:   "images [OPTIONS]",
	Description: "List images",
	Help:        "List images.",
}

func init() {
	cmdImages.Flag.BoolVar(&imagesA, []string{"a", "-all"}, false, "Show all iamges")
	cmdImages.Flag.BoolVar(&imagesNoTrunc, []string{"-no-trunc"}, false, "Don't truncate output")
	cmdImages.Flag.BoolVar(&imagesQ, []string{"q", "-quiet"}, false, "Only show numeric IDs")
	cmdImages.Flag.BoolVar(&imagesHelp, []string{"h", "-help"}, false, "Print usage")
}

// Flags
var imagesA bool       // -a flag
var imagesQ bool       // -q flag
var imagesNoTrunc bool // -no-trunc flag
var imagesHelp bool    // -h, --help flag

func runImages(cmd *Command, rawArgs []string) {
	if imagesHelp {
		cmd.PrintUsage()
	}
	if len(rawArgs) != 0 {
		cmd.PrintShortUsage()
	}

	args := commands.ImagesArgs{
		All:     imagesA,
		Quiet:   imagesQ,
		NoTrunc: imagesNoTrunc,
	}
	ctx := cmd.GetContext(rawArgs)
	err := commands.RunImages(ctx, args)
	if err != nil {
		logrus.Fatalf("Cannot execute 'images': %v", err)
	}
}
