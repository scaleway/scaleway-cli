// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package commands

import (
	"fmt"
	"strings"

	"github.com/scaleway/scaleway-cli/pkg/api"
)

// CreateArgs are arguments passed to `RunCreate`
type CreateArgs struct {
	Name       string
	Bootscript string
	Tags       []string
	Volumes    []string
	Image      string
}

// RunCreate is the handler for 'scw create'
func RunCreate(ctx CommandContext, args CreateArgs) error {
	env := strings.Join(args.Tags, " ")
	volume := strings.Join(args.Volumes, " ")
	serverID, err := api.CreateServer(ctx.API, args.Image, args.Name, args.Bootscript, env, volume, true)
	if err != nil {
		return err
	}

	fmt.Fprintln(ctx.Stdout, serverID)

	return nil
}
