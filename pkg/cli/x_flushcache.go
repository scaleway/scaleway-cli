// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package cli

import "fmt"

var cmdFlushCache = &Command{
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

func runFlushCache(cmd *Command, args []string) error {
	if flushCacheHelp {
		return cmd.PrintUsage()
	}
	if len(args) > 0 {
		return cmd.PrintShortUsage()
	}

	err := cmd.API.Cache.Flush()
	if err != nil {
		return fmt.Errorf("Failed to flush the cache")
	}
	fmt.Println("Cache flushed")
	return nil
}
