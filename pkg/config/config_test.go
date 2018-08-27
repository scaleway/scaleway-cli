// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package config

import (
	"os"
	"strings"
	"testing"
	"strconv"
	"math/rand"

	"github.com/scaleway/scaleway-cli/pkg/scwversion"
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

func TestGetConfigFilePathEnv(t *testing.T) {
	Convey("Testing GetConfigFilePath() with env variable", t, func() {
		os.Setenv("SCW_CONFIG_PATH", "./config_testdata1")
		configPath, err := GetConfigFilePath()
		So(err, ShouldBeNil)
		So(configPath, ShouldEqual, "./config_testdata1")
		os.Unsetenv("SCW_CONFIG_PATH")
	})
}

func TestGetConfig(t *testing.T) {
	Convey("Testing GetConfig() with env variable", t, func() {
		os.Setenv("SCW_CONFIG_PATH", "./config_testdata1")
		config, err := GetConfig()
		So(err, ShouldBeNil)
		So(config.Version, ShouldEqual, "test_version")
		So(config.Organization, ShouldEqual, "test_orgID")
		So(config.Token, ShouldEqual, "test_token")
		os.Unsetenv("SCW_CONFIG_PATH")
	})
}

func TestSave(t *testing.T) {
	Convey("Testing SaveConfig() with env variable", t, func() {
		os.Setenv("SCW_CONFIG_PATH", "./config_testdata2")
		randOrg := strconv.FormatInt(rand.Int63(), 16)
		randToken := strconv.FormatInt(rand.Int63(), 16)
		cfg := &Config{
			Organization: strings.Trim(randOrg, "\n"),
			Token:        strings.Trim(randToken, "\n"),
		}
		err := cfg.Save()
		So(err, ShouldBeNil)
		So(cfg.Version, ShouldEqual, scwversion.VERSION)
		So(cfg.Organization, ShouldEqual, randOrg)
		So(cfg.Token, ShouldEqual, randToken)
		os.Unsetenv("SCW_CONFIG_PATH")
	})
}




func TestGetHomeDir(t *testing.T) {
	Convey("Testing GetHomeDir()", t, func() {
		homedir, err := GetHomeDir()
		So(err, ShouldBeNil)
		So(homedir, ShouldNotEqual, "")
	})
}
