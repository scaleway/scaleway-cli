// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package config

import (
	"math/rand"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/scaleway/scaleway-cli/pkg/scwversion"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGetConfigFilePath(t *testing.T) {
	Convey("Testing GetConfigFilePath()", t, func() {
		os.Unsetenv("SCW_CONFIG_PATH")
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
	Convey("Testing GetConfig() with and without env variable", t, func() {
		rand.Seed(time.Now().UTC().UnixNano())
		randOrg := strconv.FormatInt(rand.Int63(), 16)
		randToken := strconv.FormatInt(rand.Int63(), 16)
		cfg := &Config{
			Organization: strings.Trim(randOrg, "\n"),
			Token:        strings.Trim(randToken, "\n"),
		}
		os.Setenv("SCW_CONFIG_PATH", "./config_testdata1")
		err := cfg.Save("")
		So(err, ShouldBeNil)
		cfg, err = GetConfig("./config_testdata1")
		So(cfg.Organization, ShouldEqual, randOrg)
		So(cfg.Token, ShouldEqual, randToken)
		os.Unsetenv("SCW_CONFIG_PATH")
		cfg, err = GetConfig("./config_testdata1")
		So(err, ShouldBeNil)
		So(cfg.Organization, ShouldEqual, randOrg)
		So(cfg.Token, ShouldEqual, randToken)
		os.Setenv("SCW_ORGANIZATION", randOrg)
		os.Setenv("SCW_TOKEN", randToken)
		cfg, err = GetConfig("")
		So(err, ShouldBeNil)
		So(cfg.Organization, ShouldEqual, randOrg)
		So(cfg.Token, ShouldEqual, randToken)
		os.Unsetenv("SCW_ORGANIZATION")
		os.Unsetenv("SCW_TOKEN")
	})
}

func TestSave(t *testing.T) {
	Convey("Testing SaveConfig() with and without env variable", t, func() {
		os.Setenv("SCW_CONFIG_PATH", "./config_testdata2")
		rand.Seed(time.Now().UTC().UnixNano())
		randOrg := strconv.FormatInt(rand.Int63(), 16)
		randToken := strconv.FormatInt(rand.Int63(), 16)
		cfg := &Config{
			Organization: strings.Trim(randOrg, "\n"),
			Token:        strings.Trim(randToken, "\n"),
		}
		err := cfg.Save("")
		So(err, ShouldBeNil)
		cfg, err = GetConfig("")
		So(err, ShouldBeNil)
		So(cfg.Version, ShouldEqual, scwversion.VERSION)
		So(cfg.Organization, ShouldEqual, randOrg)
		So(cfg.Token, ShouldEqual, randToken)
		os.Unsetenv("SCW_CONFIG_PATH")

		randOrg = strconv.FormatInt(rand.Int63(), 16)
		randToken = strconv.FormatInt(rand.Int63(), 16)
		cfg = &Config{
			Organization: strings.Trim(randOrg, "\n"),
			Token:        strings.Trim(randToken, "\n"),
		}
		err = cfg.Save("./config_testdata2")
		So(err, ShouldBeNil)
		cfg, err = GetConfig("./config_testdata2")
		So(err, ShouldBeNil)
		So(cfg.Version, ShouldEqual, scwversion.VERSION)
		So(cfg.Organization, ShouldEqual, randOrg)
		So(cfg.Token, ShouldEqual, randToken)

	})
}

func TestGetHomeDir(t *testing.T) {
	Convey("Testing GetHomeDir()", t, func() {
		homedir, err := GetHomeDir()
		So(err, ShouldBeNil)
		So(homedir, ShouldNotEqual, "")
	})
}
