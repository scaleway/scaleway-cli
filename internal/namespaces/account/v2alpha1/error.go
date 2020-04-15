package account

import (
	"fmt"

	"github.com/scaleway/scaleway-cli/internal/core"
)

func installationCanceled(addKeyInstructions string) *core.CliError {
	return &core.CliError{
		Err:  fmt.Errorf("installation of SSH key canceled"),
		Hint: "You can add it later using " + addKeyInstructions,
	}
}

func sshKeyNotFound(filename string, addKeyInstructions string) *core.CliError {
	return &core.CliError{
		Err:  fmt.Errorf("could not find an SSH key at " + filename),
		Hint: "You can add one later using " + addKeyInstructions,
	}
}
