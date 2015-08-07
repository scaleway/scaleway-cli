// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package commands

import (
	"fmt"
	"runtime"

	"github.com/scaleway/scaleway-cli/pkg/scwversion"
)

// VersionArgs are flags for the `RunVersion` function
type VersionArgs struct{}

// Version is the handler for 'scw version'
func Version(ctx CommandContext, args VersionArgs) error {
	fmt.Fprintf(ctx.Stdout, "Client version: %s\n", scwversion.VERSION)
	fmt.Fprintf(ctx.Stdout, "Go version (client): %s\n", runtime.Version())
	fmt.Fprintf(ctx.Stdout, "Git commit (client): %s\n", scwversion.GITCOMMIT)
	fmt.Fprintf(ctx.Stdout, "OS/Arch (client): %s/%s\n", runtime.GOOS, runtime.GOARCH)
	// FIXME: API version information

	return nil
}
