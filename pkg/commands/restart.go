// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package commands

import (
	"fmt"

	"github.com/scaleway/scaleway-cli/vendor/github.com/Sirupsen/logrus"
)

// RestartArgs are flags for the `RunRestart` function
type RestartArgs struct {
	Servers []string
}

// RunRestart is the handler for 'scw restart'
func RunRestart(ctx CommandContext, args RestartArgs) error {
	hasError := false
	for _, needle := range args.Servers {
		server := ctx.API.GetServerID(needle)
		err := ctx.API.PostServerAction(server, "reboot")
		if err != nil {
			if err.Error() != "server is being stopped or rebooted" {
				logrus.Errorf("failed to restart server %s: %s", server, err)
				hasError = true
			}
		} else {
			fmt.Fprintln(ctx.Stdout, needle)
		}
		if hasError {
			return fmt.Errorf("at least 1 server failed to restart")
		}
	}
	return nil
}
