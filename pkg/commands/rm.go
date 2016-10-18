// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package commands

import (
	"fmt"

	"github.com/Sirupsen/logrus"
)

// RmArgs are flags for the `RunRm` function
type RmArgs struct {
	Servers []string
	Force   bool
}

// RunRm is the handler for 'scw rm'
func RunRm(ctx CommandContext, args RmArgs) error {
	hasError := false
	for _, needle := range args.Servers {
		server, err := ctx.API.GetServerID(needle)
		if err != nil {
			return err
		}
		if args.Force {
			err = ctx.API.DeleteServerForce(server)
		} else {
			err = ctx.API.DeleteServer(server)
		}
		if err != nil {
			logrus.Errorf("failed to delete server %s: %s", server, err)
			hasError = true
		} else {
			fmt.Fprintln(ctx.Stdout, needle)
		}
	}
	if hasError {
		return fmt.Errorf("at least 1 server failed to be removed")
	}
	return nil
}
