//go:build wasm

package instance

import (
	"context"
	"errors"

	"github.com/scaleway/scaleway-sdk-go/scw"

	"github.com/scaleway/scaleway-cli/v2/core"
)

func serverConsoleCommand() *core.Command {
	return nil
}

func connectServerConsole(_ context.Context, _ scw.Zone, _ string) error {
	return errors.New("the serial console is not available in this build")
}
