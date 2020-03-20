package account

import (
	"fmt"

	"github.com/scaleway/scaleway-cli/internal/core"
)

func installationCanceled(addKeyInstructions string) *core.CliError {
	return &core.CliError{
		Err:  fmt.Errorf("installation of SSH key canceled"),
		Hint: "Add it later using " + addKeyInstructions,
	}
}

func sshKeyNotFound(filename string, addKeyInstructions string) *core.CliError {
	return &core.CliError{
		Err:  fmt.Errorf("could not find an ssh key at " + filename),
		Hint: "Add one later using " + addKeyInstructions,
	}
}

func sshKeyAlreadyPresent(shortenedFilename string) *core.CliError {
	return &core.CliError{
		Err: fmt.Errorf("key " + shortenedFilename + " is already present on your scaleway account"),
	}
}
