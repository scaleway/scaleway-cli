package utils

import (
	"sort"
	"strings"
	"testing"

	. "github.com/scaleway/scaleway-cli/vendor/github.com/smartystreets/goconvey/convey"
)

func TestWordify(t *testing.T) {
	Convey("Testing Wordify()", t, func() {
		So(Wordify("Hello World 42 !!"), ShouldEqual, "Hello_World_42")
		So(Wordify("  Hello   World   42   !!  "), ShouldEqual, "Hello_World_42")
		So(Wordify("Hello_World_42"), ShouldEqual, "Hello_World_42")
		So(Wordify(""), ShouldEqual, "")
	})
}

func TestTruncIf(t *testing.T) {
	Convey("Testing TruncIf()", t, func() {
		So(TruncIf("Hello World", 5, false), ShouldEqual, "Hello World")
		So(TruncIf("Hello World", 5, true), ShouldEqual, "Hello")
		So(TruncIf("Hello World", 50, false), ShouldEqual, "Hello World")
		So(TruncIf("Hello World", 50, true), ShouldEqual, "Hello World")
	})
}

func TestPathToTARPathparts(t *testing.T) {
	Convey("Testing PathToTARPathparts()", t, func() {
		dir, base := PathToTARPathparts("/etc/passwd")
		So([]string{"/etc", "passwd"}, ShouldResemble, []string{dir, base})

		dir, base = PathToTARPathparts("/etc")
		So([]string{"/", "etc"}, ShouldResemble, []string{dir, base})

		dir, base = PathToTARPathparts("/etc/")
		So([]string{"/", "etc"}, ShouldResemble, []string{dir, base})

		dir, base = PathToTARPathparts("/long/path/to/file")
		So([]string{"/long/path/to", "file"}, ShouldResemble, []string{dir, base})

		dir, base = PathToTARPathparts("/long/path/to/dir/")
		So([]string{"/long/path/to", "dir"}, ShouldResemble, []string{dir, base})
	})
}

func TestRemoveDuplicates(t *testing.T) {
	Convey("Testing RemoveDuplicates()", t, func() {
		slice := RemoveDuplicates([]string{"a", "b", "c", "a"})
		sort.Strings(slice)
		So(slice, ShouldResemble, []string{"a", "b", "c"})

		slice = RemoveDuplicates([]string{"a", "b", "c", "a"})
		sort.Strings(slice)
		So(slice, ShouldResemble, []string{"a", "b", "c"})

		slice = RemoveDuplicates([]string{"a", "b", "c", "a", "a", "b", "d"})
		sort.Strings(slice)
		So(slice, ShouldResemble, []string{"a", "b", "c", "d"})

		slice = RemoveDuplicates([]string{"a", "b", "c", "a", ""})
		sort.Strings(slice)
		So(slice, ShouldResemble, []string{"", "a", "b", "c"})
	})
}

func TestGetHomeDir(t *testing.T) {
	Convey("Testing GetHomeDir()", t, func() {
		homedir, err := GetHomeDir()
		So(err, ShouldBeNil)
		So(homedir, ShouldNotEqual, "")
	})
}

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
