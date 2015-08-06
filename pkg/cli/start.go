// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package cli

import (
	"fmt"
	"os"
	"time"

	log "github.com/scaleway/scaleway-cli/vendor/github.com/Sirupsen/logrus"

	"github.com/scaleway/scaleway-cli/pkg/api"
)

var cmdStart = &Command{
	Exec:        runStart,
	UsageLine:   "start [OPTIONS] SERVER [SERVER...]",
	Description: "Start a stopped server",
	Help:        "Start a stopped server.",
}

func init() {
	cmdStart.Flag.BoolVar(&startW, []string{"w", "-wait"}, false, "Synchronous start. Wait for SSH to be ready")
	cmdStart.Flag.Float64Var(&startTimeout, []string{"T", "-timeout"}, 0, "Set timeout values to seconds")
	cmdStart.Flag.BoolVar(&startHelp, []string{"h", "-help"}, false, "Print usage")
}

// Flags
var startW bool          // -w flag
var startTimeout float64 // -T flag
var startHelp bool       // -h, --help flag

func runStart(cmd *Command, args []string) {
	if startHelp {
		cmd.PrintUsage()
	}
	if len(args) < 1 {
		cmd.PrintShortUsage()
	}

	hasError := false
	errChan := make(chan error)
	successChan := make(chan bool)
	remainingItems := len(args)

	for i := range args {
		needle := args[i]
		go api.StartServerOnce(cmd.API, needle, startW, successChan, errChan)
	}

	if startTimeout > 0 {
		go func() {
			time.Sleep(time.Duration(startTimeout*1000) * time.Millisecond)
			log.Fatalf("Operation timed out")
		}()
	}

	for {
		select {
		case _ = <-successChan:
			remainingItems--
		case err := <-errChan:
			log.Errorf(fmt.Sprintf("%s", err))
			remainingItems--
			hasError = true
		}

		if remainingItems == 0 {
			break
		}
	}
	if hasError {
		os.Exit(1)
	}
}
