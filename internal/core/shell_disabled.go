//go:build freebsd
// +build freebsd

package core

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func RunShell(ctx context.Context, printer *Printer, meta *meta, rootCmd *cobra.Command, args []string) {
	err := printer.Print(fmt.Errorf("shell is currently disabled on freebsd"), nil)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}
}
