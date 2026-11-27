package connect_test

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/internal/connect"
)

func TestCommandCanceledBeforeExec(t *testing.T) {
	t.Parallel()

	binary, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(t.Context())
	cmd := connect.CommandContext(ctx, binary, "-test.run=^$")
	cancel()

	_, err = connect.ExecCmd(ctx, cmd)
	if !errors.Is(err, context.Canceled) || cmd.Process != nil {
		t.Fatalf("process = %v, error = %v; want nil, context.Canceled", cmd.Process, err)
	}
}

func TestIsInterrupted(t *testing.T) {
	t.Parallel()

	if connect.IsInterrupted(errors.New("ssh failed")) {
		t.Fatal("plain error must not be interrupted")
	}
	ie := &connect.InterruptedError{Signal: os.Interrupt, Code: 130}
	if !connect.IsInterrupted(ie) {
		t.Fatal("InterruptedError must match")
	}
}

func TestSetEnv(t *testing.T) {
	t.Parallel()

	cmd := exec.Command("true")
	connect.SetEnv(cmd, "FOO", "bar")

	if v, ok := cmdEnvValue(cmd.Env, "FOO"); !ok || v != "bar" {
		t.Errorf("FOO = %q (found %v), want %q", v, ok, "bar")
	}
	// The parent environment must be inherited: exec.Cmd with an empty Env
	// inherits it, so SetEnv must not replace it with just FOO.
	if _, ok := cmdEnvValue(cmd.Env, "PATH"); !ok {
		t.Error("PATH missing from cmd.Env, parent environment was dropped")
	}
}

func cmdEnvValue(env []string, key string) (string, bool) {
	for _, e := range env {
		if len(e) > len(key)+1 && e[:len(key)+1] == key+"=" {
			return e[len(key)+1:], true
		}
	}

	return "", false
}
