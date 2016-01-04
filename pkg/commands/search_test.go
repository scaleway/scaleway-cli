// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package commands

import (
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func ExampleRunSearch() {
	ctx := testCommandContext()
	args := SearchArgs{}
	RunSearch(ctx, args)
}

func ExampleRunSearch_complex() {
	ctx := testCommandContext()
	args := SearchArgs{
		Term:    "",
		NoTrunc: false,
	}
	RunSearch(ctx, args)
}

func TestRunSearch_realAPI(t *testing.T) {
	ctx := RealAPIContext()
	if ctx == nil {
		t.Skip()
	}
	Convey("Testing RunSearch() on real API", t, func() {
		Convey("ubuntu", func() {
			args := SearchArgs{
				Term:    "ubuntu",
				NoTrunc: false,
			}

			scopedCtx, scopedStdout, scopedStderr := getScopedCtx(ctx)
			err := RunSearch(*scopedCtx, args)
			So(err, ShouldBeNil)
			So(scopedStderr.String(), ShouldBeEmpty)

			lines := strings.Split(scopedStdout.String(), "\n")
			So(len(lines), ShouldBeGreaterThan, 0)

			firstLine := lines[0]
			colNames := strings.Fields(firstLine)
			So(colNames, ShouldResemble, []string{"NAME", "DESCRIPTION", "STARS", "OFFICIAL", "AUTOMATED"})
		})

		// FIXME: test invalid word
		// FIXME: test no-trunc
	})
}
