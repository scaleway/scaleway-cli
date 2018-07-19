// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package cli

import (
	"fmt"
	"reflect"

	"github.com/scaleway/scaleway-cli/pkg/commands"
)

var cmdS3 = &Command{
	Exec:        runS3,
	UsageLine:   "s3 [OPTIONS]",
	Description: "Access to s3 bucket",
	Help:        "Access to s3 bucket.",
}

func init() {
	fmt.Println("reflect =", reflect.TypeOf(cmdS3))
	cmdS3.Flag.StringVar(&s3Profile, []string{"-profile"}, "scw", "Specify a profile")
	cmdS3.Flag.BoolVar(&s3Help, []string{"h", "-help"}, false, "Print usage")
}

// Flags
var s3Help bool      // -h, --help flag
var s3Profile string // -p, --profile flag

func runS3(cmd *Command, rawArgs []string) error {
	if s3Help {
		return cmd.PrintUsage()
	}

	fmt.Println("s3Profile =", s3Profile)
	args := commands.S3Args{}
	ctx := cmd.GetContext(rawArgs)
	return commands.S3(ctx, args)
}
