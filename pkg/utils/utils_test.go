package utils

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
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
