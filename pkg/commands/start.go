// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package commands

import (
	"fmt"
	"time"

	"github.com/scaleway/scaleway-cli/pkg/api"
	"github.com/scaleway/scaleway-cli/vendor/github.com/Sirupsen/logrus"
)

// StartArgs are flags for the `RunStart` function
type StartArgs struct {
	Servers []string
	Wait    bool
	Timeout float64
}

// RunStart is the handler for 'scw start'
func RunStart(ctx CommandContext, args StartArgs) error {
	hasError := false
	errChan := make(chan error)
	successChan := make(chan bool)
	remainingItems := len(args.Servers)

	for _, needle := range args.Servers {
		go api.StartServerOnce(ctx.API, needle, args.Wait, successChan, errChan)
	}

	if args.Timeout > 0 {
		go func() {
			time.Sleep(time.Duration(args.Timeout*1000) * time.Millisecond)
			// FIXME: avoid use of fatalf
			logrus.Fatalf("Operation timed out")
		}()
	}

	for {
		select {
		case _ = <-successChan:
			remainingItems--
		case err := <-errChan:
			logrus.Errorf(fmt.Sprintf("%s", err))
			remainingItems--
			hasError = true
		}

		if remainingItems == 0 {
			break
		}
	}
	if hasError {
		return fmt.Errorf("at least 1 server failed to start")
	}
	return nil
}
