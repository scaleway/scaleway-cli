// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package commands

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/scaleway/scaleway-cli/pkg/api"
)

// WaitArgs are flags for the `RunWait` function
type WaitArgs struct {
	Servers []string
}

// RunWait is the handler for 'scw wait'
func RunWait(ctx CommandContext, args WaitArgs) error {
	hasError := false
	for _, needle := range args.Servers {
		serverIdentifier, err := ctx.API.GetServerID(needle)
		if err != nil {
			logrus.Error(err)
			hasError = true
		} else {
			if _, err := api.WaitForServerStopped(ctx.API, serverIdentifier); err != nil {
				logrus.Errorf("failed to wait for server %s: %v", serverIdentifier, err)
				hasError = true
			}
		}
	}

	if hasError {
		return fmt.Errorf("at least 1 server failed to be stopped")
	}
	return nil
}
