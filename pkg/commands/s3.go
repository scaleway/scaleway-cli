// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package commands

import (
	minio "github.com/minio/mc/cmd"
)

// VersionArgs are flags for the `RunVersion` function
type S3Args struct{}

// Version is the handler for 'scw version'
func S3(ctx CommandContext, args S3Args) error {
	minio.Main()
	return nil
}
