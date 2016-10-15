// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package commands

import (
	"fmt"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/scaleway/scaleway-cli/pkg/api"
)

// StartArgs are flags for the `RunStart` function
type StartArgs struct {
	Servers  []string
	Wait     bool
	Timeout  float64
	SetState string
}

// RunStart is the handler for 'scw start'
func RunStart(ctx CommandContext, args StartArgs) error {
	hasError := false
	errChan := make(chan error)
	successChan := make(chan string)
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
		case name := <-successChan:
			fmt.Println(name)
			remainingItems--
		case err := <-errChan:
			logrus.Errorf("%s", err)
			remainingItems--
			hasError = true
		}

		if remainingItems == 0 {
			break
		}
	}
	if args.SetState != "" {
		var wg sync.WaitGroup

		for _, needle := range args.Servers {
			wg.Add(1)
			go func(n string) {
				defer wg.Done()
				serverID, err := ctx.API.GetServerID(n)
				if err != nil {
					return
				}

				for {
					server, err := ctx.API.GetServer(serverID)
					if err != nil {
						logrus.Errorf("%s", err)
						return
					}
					if server.StateDetail == "kernel-started" {
						err = ctx.API.PatchServer(serverID, api.ScalewayServerPatchDefinition{
							StateDetail: &args.SetState,
						})
						if err != nil {
							logrus.Errorf("%s", err)
						}
						return
					}
					time.Sleep(1 * time.Second)
				}
			}(needle)
		}
		wg.Wait()
	}
	if hasError {
		return fmt.Errorf("at least 1 server failed to start")
	}
	return nil
}
