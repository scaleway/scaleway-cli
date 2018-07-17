// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package cli

import (
	"github.com/scaleway/scaleway-cli/pkg/commands"
)

var cmdS3 = &Command{
	Exec:        runS3,
	UsageLine:   "s3 [OPTIONS]",
	Description: "Access to s3 bucket",
	Help:        "Access to s3 bucket.",
}

func init() {
	cmdS3.Flag.StringVar(&s3Profile, []string{"-profile"}, "", "Specify a profile")
	cmdS3.Flag.BoolVar(&s3Help, []string{"h", "-help"}, false, "Print usage")
}

// Flags
var s3Help bool      // -h, --help flag
var s3Profile string // -p, --profile flag

func runS3(cmd *Command, rawArgs []string) error {
	if s3Help {
		return cmd.PrintUsage()
	}

	args := commands.S3Args{}
	ctx := cmd.GetContext(rawArgs)
	return commands.S3(ctx, args)
}
