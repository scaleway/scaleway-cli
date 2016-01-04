// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package config

import (
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetConfigFilePath(t *testing.T) {
	Convey("Testing GetConfigFilePath()", t, func() {
		configPath, err := GetConfigFilePath()
		So(err, ShouldBeNil)
		So(configPath, ShouldNotEqual, "")

		homedir, err := GetHomeDir()
		So(err, ShouldBeNil)
		So(strings.Contains(configPath, homedir), ShouldBeTrue)
	})
}

func TestGetHomeDir(t *testing.T) {
	Convey("Testing GetHomeDir()", t, func() {
		homedir, err := GetHomeDir()
		So(err, ShouldBeNil)
		So(homedir, ShouldNotEqual, "")
	})
}
