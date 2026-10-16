package core

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"time"
)

type OverrideExecFunc func(cmd *exec.Cmd) (exitCode int, err error)

func defaultOverrideExec(cmd *exec.Cmd) (exitCode int, err error) {
	err = cmd.Run()
	if exitErr, isExitErr := err.(*exec.ExitError); isExitErr {
		return exitErr.ExitCode(), nil
	}

	return 0, err
}

func ExecCmd(ctx context.Context, cmd *exec.Cmd) (exitCode int, err error) {
	meta := ExtractMeta(ctx)

	// We do not support override of stdin
	if cmd.Stdin == nil {
		cmd.Stdin = os.Stdin
	}

	cmd.Stdout = meta.stdout
	cmd.Stderr = meta.stderr

	return meta.OverrideExec(cmd)
}

// CheckReachable dials the TCP endpoint with a short timeout so connect
// commands fail fast on an unreachable endpoint instead of hanging in the client.
func CheckReachable(host string, port uint32) error {
	conn, err := (&net.Dialer{Timeout: 5 * time.Second}).Dial(
		"tcp",
		net.JoinHostPort(host, strconv.FormatUint(uint64(port), 10)),
	)
	if err != nil {
		return fmt.Errorf(
			"cannot reach %s:%d from this machine (%v); check your network/firewall, "+
				"or use --private-network if a private endpoint is attached",
			host,
			port,
			err,
		)
	}
	conn.Close()

	return nil
}
