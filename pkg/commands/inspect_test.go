// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package commands

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/scaleway/scaleway-cli/pkg/api"
	. "github.com/smartystreets/goconvey/convey"
)

func ExampleRunInspect() {
	ctx := testCommandContext()
	args := InspectArgs{}
	RunInspect(ctx, args)
}

func ExampleRunInspect_complex() {
	ctx := testCommandContext()
	args := InspectArgs{
		Format:      "",
		Browser:     false,
		Identifiers: []string{},
	}
	RunInspect(ctx, args)
}

func TestRunInspect_realAPI(t *testing.T) {
	ctx := RealAPIContext()
	if ctx == nil {
		t.Skip()
	}
	Convey("Testing RunInspect() on real API", t, func() {
		Convey("image:ubuntu-wily", func() {
			args := InspectArgs{
				Format:      "",
				Browser:     false,
				Identifiers: []string{"image:ubuntu-wily"},
				Arch:        "arm",
			}

			scopedCtx, scopedStdout, scopedStderr := getScopedCtx(ctx)
			err := RunInspect(*scopedCtx, args)
			So(err, ShouldBeNil)
			So(scopedStderr.String(), ShouldBeEmpty)
			fmt.Println(scopedStdout)
			var results []api.ScalewayImage
			err = json.Unmarshal(scopedStdout.Bytes(), &results)
			So(err, ShouldBeNil)
			So(len(results), ShouldEqual, 1)
			So(strings.ToLower(results[0].Name), ShouldContainSubstring, "ubuntu")
			So(strings.ToLower(results[0].Name), ShouldContainSubstring, "wily")

			Convey("-f \"{{.Identifier}}\" image:ubuntu-wily", func() {
				args := InspectArgs{
					Format:      "{{.Identifier}}",
					Browser:     false,
					Identifiers: []string{"image:ubuntu-wily"},
					Arch:        "arm",
				}

				scopedCtx, scopedStdout, scopedStderr := getScopedCtx(ctx)
				err := RunInspect(*scopedCtx, args)
				So(err, ShouldBeNil)
				So(scopedStderr.String(), ShouldBeEmpty)
				uuid := strings.TrimSpace(scopedStdout.String())
				So(results[0].Identifier, ShouldEqual, uuid)
			})
		})
	})
}
