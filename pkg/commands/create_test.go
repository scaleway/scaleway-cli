// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package commands

import (
	"bytes"
	"strings"
	"testing"

	. "github.com/scaleway/scaleway-cli/vendor/github.com/smartystreets/goconvey/convey"
)

func ExampleRunCreate() {
	ctx := ExampleCommandContext()
	args := CreateArgs{}
	RunCreate(ctx, args)
}

func ExampleRunCreate_complex() {
	ctx := ExampleCommandContext()
	args := CreateArgs{
		Name:       "test",
		Bootscript: "rescue",
		Tags:       []string{"tag1", "tag2"},
		Volumes:    []string{},
		Image:      "ubuntu-vivid",
		TmpSSHKey:  false,
	}
	RunCreate(ctx, args)
}

func TestRunCreate_realAPI(t *testing.T) {
	createdUUIDs := []string{}
	ctx := RealAPIContext()
	if ctx == nil {
		t.Skip()
	}

	// FIXME: test empty settings
	// FIXME: test cache duplicates
	// FIXME: test TmpSSHKey
	// FIXME: test Volumes
	// FIXME: test Tags
	// FIXME: test Bootscript

	Convey("Testing RunCreate() on real API", t, func() {
		// Error when image is empty !
		/*
			Convey("no options", func() {
				args := CreateArgs{}
				err := RunCreate(*ctx, args)
				So(err, ShouldBeNil)

				stderr := ctx.Stderr.(*bytes.Buffer).String()
				stdout := ctx.Stdout.(*bytes.Buffer).String()
				fmt.Println(stderr)
				fmt.Println(stdout)
			})
		*/

		Convey("--name=unittest-create-standard ubuntu-vivid", func() {
			args := CreateArgs{
				Name:  "unittest-create-standard",
				Image: "ubuntu-vivid",
			}
			err := RunCreate(*ctx, args)
			So(err, ShouldBeNil)

			stderr := ctx.Stderr.(*bytes.Buffer).String()
			stdout := ctx.Stdout.(*bytes.Buffer).String()
			So(stderr, ShouldBeEmpty)
			So(stdout, shouldBeAnUUID)

			createdUUIDs = append(createdUUIDs, strings.TrimSpace(stdout))
		})

		Reset(func() {
			ctx.Stdout.(*bytes.Buffer).Reset()
			ctx.Stderr.(*bytes.Buffer).Reset()

			if len(createdUUIDs) > 0 {
				err := RunRm(*ctx, RmArgs{
					Servers: createdUUIDs,
				})
				So(err, ShouldBeNil)

				stderr := ctx.Stderr.(*bytes.Buffer).String()
				stdout := ctx.Stdout.(*bytes.Buffer).String()
				So(stderr, ShouldBeEmpty)
				removedUUIDs := strings.Split(strings.TrimSpace(stdout), "\n")
				So(removedUUIDs, ShouldResemble, createdUUIDs)

				createdUUIDs = createdUUIDs[:0]

				ctx.Stdout.(*bytes.Buffer).Reset()
				ctx.Stderr.(*bytes.Buffer).Reset()
			}
		})
	})
}
