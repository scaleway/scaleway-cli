// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package commands

import (
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func ExampleRunCreate() {
	ctx := testCommandContext()
	args := CreateArgs{}
	RunCreate(ctx, args)
}

func ExampleRunCreate_complex() {
	ctx := testCommandContext()
	args := CreateArgs{
		Name:       "test",
		Bootscript: "rescue",
		Tags:       []string{"tag1", "tag2"},
		Volumes:    []string{},
		Image:      "ubuntu-wily",
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

		Convey("--name=unittest-create-standard wily", func() {
			args := CreateArgs{
				Name:           "unittest-create-standard",
				Image:          "wily",
				CommercialType: "X64-2GB",
				IP:             "dynamic",
			}

			scopedCtx, scopedStdout, scopedStderr := getScopedCtx(ctx)
			err := RunCreate(*scopedCtx, args)
			So(err, ShouldBeNil)
			So(scopedStderr.String(), ShouldBeEmpty)

			uuid := strings.TrimSpace(scopedStdout.String())
			So(uuid, shouldBeAnUUID)

			createdUUIDs = append(createdUUIDs, uuid)
		})

		Reset(func() {
			if len(createdUUIDs) > 0 {
				rmCtx, rmStdout, rmStderr := getScopedCtx(ctx)
				rmErr := RunRm(*rmCtx, RmArgs{Servers: createdUUIDs})
				So(rmErr, ShouldBeNil)
				So(rmStderr.String(), ShouldBeEmpty)

				removedUUIDs := strings.Split(strings.TrimSpace(rmStdout.String()), "\n")
				So(removedUUIDs, ShouldResemble, createdUUIDs)
				createdUUIDs = createdUUIDs[:0]
			}
		})
	})
}
