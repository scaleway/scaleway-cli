// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

// Package utils contains logquiet
package utils

import (
	"fmt"
	"os"
)

// LogQuiet Displays info if quiet is activated
func LogQuiet(str string) {
	if os.Getenv("QUIET") == "" {
		fmt.Fprintf(os.Stderr, "%s", str)
	}
}
