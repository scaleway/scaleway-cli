// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package commands

import (
	"fmt"
	"os"

	"github.com/scaleway/scaleway-cli/pkg/config"
)

// LogoutArgs are flags for the `RunLogout` function
type LogoutArgs struct{}

// RunLogout is the handler for 'scw logout'
func RunLogout(ctx CommandContext, args LogoutArgs) error {
	// FIXME: ask if we need to remove the local ssh key on the account
	scwrcPath, err := config.GetConfigFilePath()
	if err != nil {
		return fmt.Errorf("unable to get scwrc config file path: %v", err)
	}

	if _, err = os.Stat(scwrcPath); err == nil {
		err = os.Remove(scwrcPath)
		if err != nil {
			return fmt.Errorf("unable to remove scwrc config file: %v", err)
		}
	}
	return nil
}
