//go:build wasm

package docs

import (
	"context"
	"fmt"

	"github.com/scaleway/scaleway-cli/v2/core"
)

func runBrowser(ctx context.Context, commands *core.Commands) (any, error) {
	return nil, fmt.Errorf("docs browser is not available in WASM builds")
}
