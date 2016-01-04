// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package commands

import (
	"bytes"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestVersion(t *testing.T) {
	Convey("Testing Version()", t, func() {
		ctx := testCommandContext()
		var buf bytes.Buffer
		ctx.Stdout = &buf

		args := VersionArgs{}

		err := Version(ctx, args)

		So(err, ShouldBeNil)
		So(buf.String(), ShouldContainSubstring, "Client version: ")
		So(buf.String(), ShouldContainSubstring, "Go version (client): ")
		So(buf.String(), ShouldContainSubstring, "Git commit (client): ")
		So(buf.String(), ShouldContainSubstring, "OS/Arch (client): ")

	})
}

func ExampleVersion() {
	ctx := testCommandContext()
	args := VersionArgs{}
	Version(ctx, args)
}
