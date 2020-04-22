// Copyright (c) 2015 Manfred Touron

// +build !windows

package console

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sys/unix"
)

func notifySignalSIGWINCH(c chan<- os.Signal) {
	signal.Notify(c, syscall.SIGWINCH)
}

func resetSignalSIGWINCH() {
	signal.Reset(syscall.SIGWINCH)
}

func syscallTIOCGWINSZ() ([]byte, error) {
	ws, err := unix.IoctlGetWinsize(0, unix.TIOCGWINSZ)
	if err != nil {
		return nil, fmt.Errorf("ioctl error: %v", err)
	}
	tws := winsize{Rows: ws.Row, Columns: ws.Col}
	b, err := json.Marshal(tws)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal error: %v", err)
	}
	return b, err
}
