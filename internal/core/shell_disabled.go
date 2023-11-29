//go:build freebsd || wasm

package core

import (
	"context"
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"
)

func RunShell(ctx context.Context, printer *Printer, meta *meta, rootCmd *cobra.Command, args []string) {
	err := printer.Print(fmt.Errorf("shell is currently disabled on %s/%s", runtime.GOARCH, runtime.GOOS), nil)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}
}
