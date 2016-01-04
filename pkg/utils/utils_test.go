package utils

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
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

func TestGeneratingAnSSHKey(t *testing.T) {
	Convey("Testing GeneratingAnSSHKey()", t, func() {
		streams := SpawnRedirection{
			Stdout: os.Stdout,
			Stderr: os.Stderr,
			Stdin:  os.Stdin,
		}

		tmpDir, err := ioutil.TempDir("/tmp", "scaleway-test")
		So(err, ShouldBeNil)

		tmpFile, err := ioutil.TempFile(tmpDir, "ssh-key")
		So(err, ShouldBeNil)

		err = os.Remove(tmpFile.Name())
		So(err, ShouldBeNil)

		filePath, err := GeneratingAnSSHKey(streams, tmpDir, filepath.Base(tmpFile.Name()))
		So(err, ShouldBeNil)
		So(filePath, ShouldEqual, tmpFile.Name())

		os.Remove(tmpFile.Name())
	})
}
