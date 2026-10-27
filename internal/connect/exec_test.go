package connect_test

import (
	"os/exec"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/internal/connect"
)

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
