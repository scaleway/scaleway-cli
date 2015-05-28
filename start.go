package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
)

var cmdStart = &Command{
	Exec:        runStart,
	UsageLine:   "start [OPTIONS] SERVER [SERVER...]",
	Description: "Start a stopped server",
	Help:        "Start a stopped server.",
}

func init() {
	// FIXME: -h
	cmdStart.Flag.BoolVar(&startW, []string{"w", "-wait"}, false, "Synchronous start. Wait for SSH to be ready")
	cmdStart.Flag.Float64Var(&startTimeout, []string{"T", "-timeout"}, 0, "Set timeout values to seconds")
}

// Flags
var startW bool          // -w flag
var startTimeout float64 // -T flag

func startOnce(cmd *Command, needle string, successChan chan bool, errChan chan error) {
	server := cmd.GetServer(needle)

	err := cmd.API.PostServerAction(server, "poweron")
	if err != nil {
		if err.Error() != "server should be stopped" {
			errChan <- errors.New(fmt.Sprintf("failed to stop server %s: %v", server, err))
			return
		}
	} else {
		if startW {
			_, err = WaitForServerReady(cmd.API, server)
			if err != nil {
				errChan <- errors.New(fmt.Sprintf("Failed to wait for server %s to be ready, %v", needle, err))
				return
			}
		}
		fmt.Println(needle)
	}
	successChan <- true
}

func runStart(cmd *Command, args []string) {
	if len(args) == 0 {
		log.Fatalf("usage: scw %s", cmd.UsageLine)
	}

	hasError := false
	errChan := make(chan error)
	successChan := make(chan bool)
	remainingItems := len(args)

	for _, needle := range args {
		go startOnce(cmd, needle, successChan, errChan)
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
