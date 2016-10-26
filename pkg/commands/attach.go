// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package commands

import "github.com/scaleway/scaleway-cli/pkg/utils"

// AttachArgs are flags for the `RunAttach` function
type AttachArgs struct {
	NoStdin bool
	Server  string
}

// RunAttach is the handler for 'scw attach'
func RunAttach(ctx CommandContext, args AttachArgs) error {
	serverID, err := ctx.API.GetServerID(args.Server)
	if err != nil {
		return err
	}
	_, done, err := utils.AttachToSerial(serverID, ctx.API.Token, ctx.API.ResolveTTYUrl())
	if err != nil {
		return err
	}
	<-done
	return nil
}
