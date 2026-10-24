package iam

import (
	"errors"
	"fmt"

	"github.com/scaleway/scaleway-cli/v2/core"
)

var (
	ErrCouldNotFindAnSSHKeyAt       = errors.New("could not find an SSH key at")
	ErrInstallationOfSSHKeyCanceled = errors.New("installation of SSH key canceled")
)

func installationCanceled(addKeyInstructions string) *core.CliError {
	return &core.CliError{
		Err:  ErrInstallationOfSSHKeyCanceled,
		Hint: "You can add it later using " + addKeyInstructions,
	}
}

func sshKeyNotFound(filename string, addKeyInstructions string) *core.CliError {
	return &core.CliError{
		Err:  fmt.Errorf("%w %s", ErrCouldNotFindAnSSHKeyAt, filename),
		Hint: "You can add one later using " + addKeyInstructions,
	}
}
