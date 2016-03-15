// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package commands

import (
	"fmt"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func ExampleRunHistory() {
	ctx := testCommandContext()
	args := HistoryArgs{}
	RunHistory(ctx, args)
}

func ExampleRunHistory_complex() {
	ctx := testCommandContext()
	args := HistoryArgs{
		NoTrunc: false,
		Quiet:   false,
		Image:   "",
	}
	RunHistory(ctx, args)
}

func TestRunHistory_realAPI(t *testing.T) {
	ctx := RealAPIContext()
	if ctx == nil {
		t.Skip()
	}
	Convey("Testing RunHistory() on real API", t, func() {
		Convey("ubuntu-wily", func() {
			args := HistoryArgs{
				NoTrunc: false,
				Quiet:   false,
				Image:   "ubuntu-wily",
				Arch:    "arm",
			}

			scopedCtx, scopedStdout, scopedStderr := getScopedCtx(ctx)
			err := RunHistory(*scopedCtx, args)
			So(err, ShouldBeNil)
			So(scopedStderr.String(), ShouldBeEmpty)

			lines := strings.Split(scopedStdout.String(), "\n")
			So(len(lines), ShouldBeGreaterThan, 0)

			firstLine := lines[0]
			colNames := strings.Fields(firstLine)
			So(colNames, ShouldResemble, []string{"IMAGE", "CREATED", "CREATED", "BY", "SIZE"})

			fmt.Println(scopedStdout.String())
		})

		// FIXME: test invalid image
		// FIXME: test image with duplicates cache
		// FIXME: test quiet
		// FIXME: test no-trunc
	})
}
