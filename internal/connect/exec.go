package connect

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	"github.com/scaleway/scaleway-cli/v2/core"
)

// Command prepares a subprocess for ExecCmd.
func Command(name string, args ...string) *exec.Cmd {
	return exec.Command(name, args...) //nolint:gosec
}

// CommandContext prepares a subprocess whose lifetime is bound to ctx.
func CommandContext(ctx context.Context, name string, args ...string) *exec.Cmd {
	return exec.CommandContext(ctx, name, args...) //nolint:gosec
}

// LookPath resolves an executable using the current PATH.
func LookPath(name string) (string, error) {
	return exec.LookPath(name)
}

// SetEnv replaces any existing key in the command's environment with
// key=value, so the child process only sees this session's value.
func SetEnv(cmd *exec.Cmd, key, value string) {
	// When Env is empty exec.Cmd inherits the parent environment, so seed
	// with os.Environ() to preserve that: otherwise the child would only
	// see this one variable (no PATH, HOME, ...).
	src := cmd.Env
	if len(src) == 0 {
		src = os.Environ()
	}
	env := make([]string, 0, len(src)+1)
	for _, e := range src {
		if !strings.HasPrefix(e, key+"=") {
			env = append(env, e)
		}
	}
	env = append(env, key+"="+value)
	cmd.Env = env
}

// ExecCmd runs a child process with the CLI streams and test hook.
// On unix, SIGINT/SIGTERM received while the session is running are forwarded
// to the child process, so an interrupted session (e.g. Ctrl-C during
// `kubectl get pods -w`) terminates the child and returns through the
// command's cleanup (e.g. temporary kubeconfig removal).
// A child killed by a signal exits with 128+signal, the shell convention.
func ExecCmd(ctx context.Context, cmd *exec.Cmd) (exitCode int, err error) {
	return execCmd(ctx, cmd, runChild)
}

// SignalSession runs successive child processes under a single signal subscription.
// Its methods must be called sequentially inside WithSignalSession.
type SignalSession struct {
	signals chan os.Signal
	caught  os.Signal
}

// WithSignalSession keeps interrupts observable between successive commands.
func WithSignalSession(ctx context.Context, run func(*SignalSession) error) error {
	session := listenForSignals()
	defer signal.Stop(session.signals)
	if err := session.Err(ctx); err != nil {
		return err
	}
	if err := run(session); err != nil {
		return err
	}

	return session.Err(ctx)
}

func listenForSignals() *SignalSession {
	s := &SignalSession{signals: make(chan os.Signal, 1)}
	signal.Notify(s.signals, os.Interrupt, syscall.SIGTERM)

	return s
}

// Err reports a signal received during or between commands, or cancellation of ctx.
func (s *SignalSession) Err(ctx context.Context) error {
	if s.caught == nil {
		select {
		case s.caught = <-s.signals:
		default:
			return context.Cause(ctx)
		}
	}
	code := 130
	if s.caught == syscall.SIGTERM {
		code = 143
	}

	return &core.CliError{Err: fmt.Errorf("received %s", s.caught), Empty: true, Code: code}
}

func (s *SignalSession) ExecCmd(ctx context.Context, cmd *exec.Cmd) (int, error) {
	if err := s.Err(ctx); err != nil {
		return 0, err
	}

	return execCmd(ctx, cmd, s.runChild)
}

func (s *SignalSession) runChild(cmd *exec.Cmd) (int, error) {
	if err := cmd.Start(); err != nil {
		return 0, err
	}
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	for {
		select {
		case sig := <-s.signals:
			if s.caught == nil {
				s.caught = sig
			}
			forwardSignal(cmd.Process, sig)
		case err := <-done:
			return childResult(err)
		}
	}
}

func childResult(err error) (int, error) {
	if exitErr, ok := errors.AsType[*exec.ExitError](err); ok {
		if status, ok := exitErr.Sys().(syscall.WaitStatus); ok && status.Signaled() {
			return 128 + int(status.Signal()), nil
		}

		return exitErr.ExitCode(), nil
	}

	return 0, err
}

func execCmd(ctx context.Context, cmd *exec.Cmd, run func(*exec.Cmd) (int, error)) (int, error) {
	meta := core.ExtractMeta(ctx)

	// We do not support override of stdin
	if cmd.Stdin == nil {
		cmd.Stdin = os.Stdin
	}
	cmd.Stdout = core.ExtractStdout(ctx)
	cmd.Stderr = core.ExtractStderr(ctx)

	if meta != nil && meta.OverrideExec != nil {
		return meta.OverrideExec(cmd)
	}

	return run(cmd)
}
