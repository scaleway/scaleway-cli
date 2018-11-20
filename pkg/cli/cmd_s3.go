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
	Description: "Access to s3 buckets",
	Help:        "Access to s3 buckets.",
}

func init() {
	cmdS3.Flag.BoolVar(&s3Help, []string{"h", "-help"}, false, "Print usage")
}

// Flags
var s3Help bool // -h, --help flag

func runS3(cmd *Command, rawArgs []string) error {
	var args commands.S3Args
	if s3Help {
		args = commands.S3Args{
			Command: make([]string, 0),
		}
		args.Command = append(args.Command, "--help")
	} else {
		args = commands.S3Args{
			Command: rawArgs,
		}
	}
	ctx := cmd.GetContext(rawArgs)
	return commands.S3(ctx, args)
}
