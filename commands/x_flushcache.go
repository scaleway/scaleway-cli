// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package commands

import (
	"fmt"
	"log"

	types "github.com/scaleway/scaleway-cli/commands/types"
)

var cmdFlushCache = &types.Command{
	Exec:        runFlushCache,
	UsageLine:   "_flush-cache [OPTIONS]",
	Description: "",
	Hidden:      true,
	Help:        "Flush cache",
}

func init() {
	cmdFlushCache.Flag.BoolVar(&flushCacheHelp, []string{"h", "-help"}, false, "Print usage")
}

// Flags
var flushCacheHelp bool // -h, --help flag

func runFlushCache(cmd *types.Command, args []string) {
	if flushCacheHelp {
		cmd.PrintUsage()
	}
	if len(args) > 0 {
		cmd.PrintShortUsage()
	}

	err := cmd.API.Cache.Flush()
	if err != nil {
		log.Fatal("Failed to flush the cache")
	}

	fmt.Println("Cache flushed")
}
