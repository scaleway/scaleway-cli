package scwversion

import (
	"testing"

	. "github.com/scaleway/scaleway-cli/vendor/github.com/smartystreets/goconvey/convey"
)

func TestInit(t *testing.T) {
	Convey("Testing init()", t, func() {
		So(VERSION, ShouldNotEqual, "")
		So(GITCOMMIT, ShouldNotEqual, "")
	})
}
