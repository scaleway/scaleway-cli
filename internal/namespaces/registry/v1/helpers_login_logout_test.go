package registry

import (
	"fmt"
	"io"

	"github.com/scaleway/scaleway-cli/internal/core"
)

type fakeDockerCommand struct {
	stdout io.Writer
	stderr io.Writer
	args   []string
}

func (f *fakeDockerCommand) Run() error {
	if f.args[0] == "login" {
		fmt.Fprintln(f.stdout, "Login Succeeded")
	} else if f.args[0] == "logout" {
		fmt.Fprintln(f.stdout, "Removing login credentials for rg.fr-par.scw.cloud")
	}
	return nil
}
func (f *fakeDockerCommand) SetStdin(stdin io.Reader) {}

func dockerFakeCommand(name string, stdin io.Reader, stdout io.Writer, stderr io.Writer, args ...string) core.Cmd {
	return &fakeDockerCommand{
		stdout,
		stderr,
		args,
	}
}

type fakePodmanCommand struct {
	stdout io.Writer
	stderr io.Writer
	args   []string
}

func (f *fakePodmanCommand) Run() error {
	if f.args[0] == "login" {
		fmt.Fprintln(f.stdout, "Login Succeeded!")
	} else if f.args[0] == "logout" {
		fmt.Fprintln(f.stdout, "Removed login credentials for rg.fr-par.scw.cloud/testcli")
	}
	return nil
}
func (f *fakePodmanCommand) SetStdin(stdin io.Reader) {}

func podmanFakeCommand(name string, stdin io.Reader, stdout io.Writer, stderr io.Writer, args ...string) core.Cmd {
	return &fakePodmanCommand{
		stdout,
		stderr,
		args,
	}
}
